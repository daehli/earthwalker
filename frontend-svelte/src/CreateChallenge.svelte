<script>
    import {onMount} from 'svelte';
    import { loc } from './stores.js';
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
    const popTIFLoc = "/public/nasa_pop_data.tif";
    // fetchPano will query the streetview API this many times before giving up
    const MAX_REQS = 30;

    const challengeCookieName = "earthwalker_lastChallenge";
    const resultCookiePrefix = "earthwalker_lastResult_";

    let ewapi = new EarthwalkerAPI();

    let statusText = "Twiddling thumbs...";

    let streetViewService = new google.maps.StreetViewService();

    let mapSettings = undefined;
    let popTIF = undefined;

    let mapID;
    let challengeID;
    let numFound = 0;
    let foundCoords = [];
    let done = false;

    // DOM elements
    let submitButton;
    // bindings
    let nickname = "";

    onMount(async () => {
        statusText = "Looking up population density data...";
        popTIF = await loadGeoTIF(popTIFLoc);
        console.log("TIF loaded"); // TODO: remove debug
        statusText = "Getting Map settings...";
        mapID = getURLParam("mapid");
        console.log("map id parsed from url params"); // TODO: remove debug
        mapSettings = await ewapi.getMap(mapID);
        console.log("map settings fetched from server"); // TODO: remove debug
        statusText = "Fetching panoramas...";
        console.log(mapSettings);
        foundCoords = await fetchPanos(streetViewService, mapSettings);
        challengeID = await submitNewChallenge();
        statusText = "Done!";
        done = true;
    });

    async function handleFormSubmit() {
        let challengeResultID = await submitNewChallengeResult();
        // set the generated challenge as the current challenge
        document.cookie = challengeCookieName + "=" + challengeID + ";path=/;max-age=172800";
        // set the generated ChallengeResult as the current ChallengeResult
        // for the Challenge with challengeID
        document.cookie = resultCookiePrefix + challengeID + "=" + challengeResultID + ";path=/;max-age=172800";
        window.location.replace("/play");
    }

    // TODO: remove debug
    function printCoords() {
        foundCoords.forEach((coord) => {console.log(coord.lat().toString() + ", " + coord.lng().toString())});
    }

    async function submitNewChallenge() {
        let convertedCoords = foundCoords.map((coord, i) => ({RoundNum: i, Location: {Lat: coord.lat(), Lng: coord.lng()}}));
        let challenge = {
            MapID: mapID,
            Places: convertedCoords
        };
        let data = await ewapi.postChallenge(challenge);
        return data.ChallengeID;
    }

    async function submitNewChallengeResult() {
        let challengeResult = {
            ChallengeID: challengeID,
            Nickname: nickname,
        };
        let data = await ewapi.postResult(challengeResult);
        return data.ChallengeResultID;
    }

    // == POPULATION DENSITY ========
    // TODO: can we find another way to do population density?
    //       This TIF requires a 6.5mb binary load.
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
        let foundLatLngs = await Promise.all(promises);
        return foundLatLngs;
    }

    async function fetchPano(svService, settings) {
        let source = settings.Source == 1 ? google.maps.StreetViewSource.OUTDOOR : google.maps.StreetViewSource.DEFAULT;
        let randomLatLng;
        let foundLatLng = null;
        for (let iters = 0; iters < MAX_REQS; iters++) {
            randomLatLng = await getRandomConstrainedLatLng(settings.Polygon, settings.MinDensity, settings.MaxDensity);
            foundLatLng = await new Promise((resolve, reject) => {
                streetViewService.getPanorama({
                    location: randomLatLng,
                    preference: SV_PREF,
                    radius: PANO_SEARCH_RADIUS,
                    source: source,
                }, (result, status) => {resolve(handlePanoResponse(result, status));});
            });
            if (foundLatLng) {
                numFound++;
                return foundLatLng
            }
        }
        
        function handlePanoResponse(result, status, foundLatLng) {
            if (status == google.maps.StreetViewStatus.OK && resultPanoIsGood(result, settings)) {
                return result.location.latLng;
            } else {
                console.log("Failed to get location; api request: " + status.toString() + "\n");
            }
        }
        // TODO: FIXME: display message to user when this happens
        //       maybe suggest creating a less specific map or allow
        //       them to try to fetch panos again.
        alert("Too many requests without a good streetview pano!  Reload the page to try again or create a map with fewer restrictions.");
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
        <h4 class="text-center" id="status">{statusText}</h4>
        <div action="" method="post">
            <div class="progress">
                <div 
                    style={"width: " + (mapSettings && mapSettings.NumRounds ? ((100 * numFound) / mapSettings.NumRounds) : 0) + "%;"} 
                    class="progress-bar" id="loading-progress" role="progressbar">
                    {numFound + "/" + (mapSettings && mapSettings.NumRounds ? mapSettings.NumRounds : 0)}
                </div>
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
                    <input required type="text" class="form-control" id="Nickname" bind:value={nickname}/>
                </div>
            </div>
            <div>
                <button bind:this={submitButton} id="submit-button" class="btn btn-primary" style="color: #fff;" disabled={!done || !nickname}>Start Challenge</button>
                {#if {done}}
                    <button id="copy-game-link" class="btn btn-primary" on:click={() => showChallengeLinkPrompt(challengeID)}>
                        Copy link to this game
                    </button>
                {/if}
            </div>
        </div>
    </form>
</main>