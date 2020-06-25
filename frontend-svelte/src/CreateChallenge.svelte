<script>
    // TODO: svelteify this file
    // TODO: Better organization of this file + additional documentation
    //       The flow is pretty confusing right now.
    // In the meantime, here's what happens in this script:
    //     * On DOM load, fetch the population density TIF and Map (id in URL),
    //       then based on that information, populate foundCoords with panos
    //       from the streetview API
    //     * Once we have mapSettings.NumRounds panos in foundCoords, automatically
    //       submit a POST request to the server with a new Challenge containing
    //       those coords.  The server responds with the ID of the new Challenge,
    //       which we store.  Then, enable the "Start Challenge" button.
    //     * When the user presses the button, we read their Nickname from the form
    //       and send it along with the saved ChallengeID as a new ChallengeResult.
    //       The server again responds with the ChallengeResult's ID.  We set two
    //       cookies storing the ChallengeID and ChallengeResultID, then redirect
    //       to /play.
    // TODO: FIXME: clean up/use async properly
    // TODO: This script currently assumes that mapSettings is good and valid.
    //       However, there's no validation in place to make sure that is the
    //       case.  We should either handle missing/bad information here, or
    //       implement a "trial run" fetchPanos at the end of the map creation
    //       process.

    // search radius in meters - using 500 (formerly 50,000) causes more NO_RESULTS
    // responses, but the API also takes much less time to fulfill the requests.
    // It also means we can use StreetViewPreference.BEST without so many duplicate
    // responses.
    const PANO_SEARCH_RADIUS = 500;
    // NEAREST or BEST.  BEST seems to give more actual streetview results (rather
    // than third party photospheres) so I'm going with that.
    const SV_PREF = google.maps.StreetViewPreference.BEST;
    // discard polar panos, they're usually garbage
    const LAT_LIMIT = 85;
    const mapURL = "/api/maps/"
    const popTIFLoc = "/static/nasa_pop_data.tif"

    const challengeCookieName = "earthwalker_lastChallenge"
    const resultCookiePrefix = "earthwalker_lastResult_"

    let statusReadout;

    let streetViewService = new google.maps.StreetViewService();

    let mapSettings = undefined;
    let popTIF = undefined;

    let mapID;
    let challengeID;
    let foundCoords = [];

    document.addEventListener('DOMContentLoaded', async (event) => {
        statusReadout = document.getElementById("status");
        statusReadout.textContent = "Looking up population density data...";
        popTIF = await loadGeoTIF(popTIFLoc);
        statusReadout.textContent = "Getting Map settings...";
        mapSettings = await fetchMapSettings(mapURL);
        statusReadout.textContent = "Fetching panoramas...";
        console.log(mapSettings);
        fetchPanos(streetViewService, mapSettings);
    });

    async function updateUI(numFound, numRounds) {
        let bar = document.getElementById("loading-progress")
        bar.setAttribute("style", "width: " + ((100 * numFound) / numRounds) + "%;");
        bar.textContent = numFound.toString() + "/" + numRounds.toString();
        if (numFound == numRounds) {
            challengeID = await submitNewChallenge();
            let challengeLink = window.location.origin + "/challenge?id=" + challengeID
            // TODO: nicer challenge link readout
            document.getElementById("status").textContent = "Done! Challenge Link: " + challengeLink;
            document.getElementById("submit-button").disabled = false;
        }
    }

    async function handleFormSubmit() {
        console.log("Form submitted!");
        let challengeResultID = await submitNewChallengeResult();
        // set the generated challenge as the current challenge
        document.cookie = challengeCookieName + ":" + challengeID + ";path=/;max-age=172800";
        // set the generated ChallengeResult as the current ChallengeResult
        // for the Challenge with challengeID
        document.cookie = resultCookiePrefix + challengeID + ":" + challengeResultID + ";path=/;max-age=172800";
        //window.location.replace("/play");
    }

    // TODO: remove debug
    function printCoords() {
        foundCoords.forEach((coord) => {console.log(coord.lat().toString() + ", " + coord.lng().toString())});
    }

    async function submitNewChallenge() {
        let convertedCoords = foundCoords.map((coord, i) => ({RoundNum: i, Location: {Lat: coord.lat(), Lng: coord.lng()}}));
        let challenge = JSON.stringify({
            MapID: mapID,
            Places: convertedCoords
        });
        let response = await fetch("api/challenges", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: challenge,
        });
        let data = await response.json();
        return data.ChallengeID;
    }

    async function submitNewChallengeResult() {
        let challengeResult = JSON.stringify({
            ChallengeID: challengeID,
            Nickname: document.getElementById("Nickname").value,
        });
        let response = await fetch("api/results", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: challengeResult,
        });
        let data = await response.json();
        return data.ChallengeResultID;
    }

    // == MAP SETTINGS ========
    async function fetchMapSettings(url) {
        let params = new URLSearchParams(window.location.search)
        // TODO: consider having default map settings if there's no ID
        if (!params.has("mapid")) {
            alert("URL has no map ID!");
            return;
        }
        mapID = params.get("mapid");
        let response = await fetch(url+mapID);
        return response.json();
    }

    // == POPULATION DENSITY ========
    async function loadGeoTIF(loc) {
        const response = await fetch("/public/nasa_pop_data.tif");
        const arrayBuffer = await response.arrayBuffer();
        return await GeoTIFF.fromArrayBuffer(arrayBuffer);
    }
    // get normalized (0.0 - 1.0) population density at lat, lng
    async function getLocationPopulation(lat, lng) {
        const delta = 0.1;
        // TODO: consider passing popTIF as arg
        let value = await popTIF.readRasters({
            bbox: [lng, lat, lng + 10 * delta, lat + 10 * delta],
            resX: delta,
            resY: delta,
        });
        let actualValue = value[0][0];
        // 255 means ocean
        if (actualValue == 255) {
            actualValue = 0;
        }
        return actualValue / 255;
    }

    // == GET PANOS ========
    async function fetchPanos(svService, settings) {
        const promises = [];
        for (let i = 0; i < settings.NumRounds; i++) {
            promises.push(fetchPano(svService, settings));
        }
        return Promise.all(promises);
    }

    async function fetchPano(svService, settings) {
        let randomLatLng = await getRandomConstrainedLatLng(settings.Polygon, settings.MinDensity, settings.MaxDensity);
        
        async function handlePanoResponse(result, status) {
            if (status == google.maps.StreetViewStatus.OK && resultPanoIsGood(result, settings)) {
                foundCoords.push(result.location.latLng);
                updateUI(foundCoords.length, settings.NumRounds);
            } else {
                console.log("Failed to get location; api request: " + status.toString() + "\n");
                fetchPano(svService, settings);
            }
        }

        let source = settings.Source == 1 ? google.maps.StreetViewSource.OUTDOOR : google.maps.StreetViewSource.DEFAULT;
        streetViewService.getPanorama({
            location: randomLatLng,
            preference: SV_PREF,
            radius: PANO_SEARCH_RADIUS,
            source: source,
        }, handlePanoResponse);
    }

    // returns whether result (pano) meets the requirements of mapInfo
    function resultPanoIsGood(result, settings) {
        if (result.location.latLng.lat() > LAT_LIMIT || result.location.latLng.lat() < -1 * LAT_LIMIT) {return false;}

        if (settings.Copyright == 1 && !result.copyright.includes("Google")) {
            return false;
        }
        if (settings.Copyright == 2 && result.copyright.includes("Google")) {
            return false;
        }

        if (settings.Connectedness == 1 && result.links.length == 0) {
            return false;
        }
        if (settings.Connectedness == 2 && result.links.length > 0) {
            return false;
        }

        let locationTurfPoint = turf.point([result.location.latLng.lng(), result.location.latLng.lat()]);
        if (settings.Polygon != null && !turf.booleanPointInPolygon(locationTurfPoint, settings.Polygon)) {
            return false;
        }

        // TODO: duplicate checking that doesn't rely on globals
        //       and doesn't have a race condition
        for (let i = 0; i < foundCoords.length; i++) {
            if (foundCoords[i].equals(result.location.latLng)) {
                console.log("duplicate!");
                return false;
            }
        }

        return true;
    }

    // get a random google.maps.LatLng within the specified polygon and with
    // a population density in the specified range
    async function getRandomConstrainedLatLng(polygon, minDensity, maxDensity) {
        // TODO: function assignment as control flow is heinous
        let getRandomLngLatInBounds;
        let pointInPolygon;
        if (polygon == null) {
            getRandomLngLatInBounds = getRandomLngLat;
            pointInPolygon = (_) => true;
        } else {
            let bounds = turf.bbox(polygon);
            getRandomLngLatInBounds = function() {
                let randomLng = (Math.random() * (bounds[2] - bounds[0]) + bounds[0]);
                let randomLat = (Math.random() * (bounds[3] - bounds[1]) + bounds[1]);
                return [randomLng, randomLat];
            }
            pointInPolygon = function(lnglat) {
                return turf.booleanPointInPolygon(turf.point(lnglat), polygon);
            }
        }

        async function popDensityInLimits(lnglat) {
            let density = (await getLocationPopulation(lnglat[1], lnglat[0])) * 100;
            return density <= maxDensity && density >= minDensity;
        }
        
        let lnglat;
        do {
            lnglat = getRandomLngLatInBounds();
        } while (!pointInPolygon(lnglat) || !(await popDensityInLimits(lnglat)));
        return new google.maps.LatLng(lnglat[1], lnglat[0]);
    }

    // get a random google.maps.LatLng, anywhere
    function getRandomLngLat() {
        let randomLng = (Math.random() * 360 - 180);
        let randomLat = (Math.random() * 180 - 90);
        return [randomLng, randomLat];
    }
</script>

<main>
    <form on:submit|preventDefault={handleFormSubmit} class="container">

        <br>

        <h2>New Challenge</h2>

        <br>

        <h4 class="text-center" id="status">Twiddling thumbs...</h4>

        <div action="" method="post">
            <div class="progress">
                <div class="progress-bar" id="loading-progress" role="progressbar"></div>
            </div>
            <small class="text-muted">
                Earthwalker is getting random locations from Google Street View.
                This happens in your browser, because there is only an API in JavaScript for this.
                Yes, that is kind of silly, I'm sorry.
            </small>

            <!-- TODO: show map settings -->

            <br/>
            <br/>

            <div class="form-group">
                <div class="input-group">
                    <div class="input-group-prepend">
                        <div class="input-group-text">Your Nickname</div>
                    </div>
                    <input required type="text" class="form-control" id="Nickname"/>
                </div>
            </div>

            <button id="submit-button" class="btn btn-primary" style="margin-bottom: 2em; color: #fff;" disabled>Start Challenge</button>

        </div>
    </form>
</main>