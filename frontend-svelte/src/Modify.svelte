<script>
    // Be warned, traveller. You are entering the domain of some very dodgy javascript
    // hacks. Maybe that is what you like. If so, please look around.

    let challengeID;
    let challenge;
    let map;
    let challengeResultID;
    let challengeResult;
    let tileServerURL;

    const challengeCookieName = "earthwalker_lastChallenge";
    const resultCookiePrefix = "earthwalker_lastResult_";

    async function injectStylesheet() {
        parseCookies();
        challengeResult = await fetchChallengeResult();
        if (!challengeResult.Guesses) {
            challengeResult.Guesses = [];
        }
        console.log("Challenge result: "); // TODO: remove debug
        console.log(challengeResult); // TODO: remove debug
        challenge = await fetchChallenge();
        console.log("Challenge: "); // TODO: remove debug
        console.log(challenge); // TODO: remove debug
        map = await fetchMap(challenge.MapID);
        console.log("Map: "); // TODO: remove debug
        console.log(map); // TODO: remove debug
        tileServerURL = await fetchTileServerURL(map.ShowLabels);
        console.log("Tile server: "); // TODO: remove debug
        console.log(tileServerURL); // TODO: remove debug

        // This MutationObserver always resets the title to earthwalker.
        let interval = setInterval(function() {
            try {
                new MutationObserver(function(mutations) {
                    if (document.title != "earthwalker") {
                        document.title = "earthwalker";
                    }
                }).observe(
                    document.querySelector('title'),
                    { childList: true }
                );
                clearInterval(interval);
            } catch (e) {
                // Title element is not there yet.
                // Wait for the next poll...
            }
        }, 50);

        createMinimap();
    }

    async function fetchTileServerURL(showLabels) {
        if (showLabels) {
            let response = await fetch("/api/config/tileserver");
            let data = await response.json();
            return data.tileserver;
        } else  {
            let response = await fetch("/api/config/nolabeltileserver");
            let data = await response.json();
            return data.nolabeltileserver;
        }
    }

    // TODO: duplicates function in CreateChallenge
    //       consolidate to API lib
    async function fetchMap(mapID) {
        let response = await fetch("/api/maps/"+mapID);
        return response.json();
    }

    async function fetchChallenge() {
        let response = await fetch("/api/challenges/"+challengeID);
        return response.json();
    }

    async function fetchChallengeResult() {
        let response = await fetch("/api/results/"+challengeResultID);
        return response.json();
    }

    function parseCookies() {
        let params = new URLSearchParams(window.location.search);
        let cookies = document.cookie.split("; ");
        if (params.has("id")) {
            challengeID = params.get("id");
        } else {
            let lastChallengeCookie = cookies.find(row => row.startsWith(challengeCookieName));
            if (lastChallengeCookie) {
                challengeID = lastChallengeCookie.split('=')[1];
            } else {
                alert("Could not determine challenge ID!");
            }
        }
        if (challengeID) {
            let lastResultCookie = cookies.find(row => row.startsWith(resultCookiePrefix + challengeID));
            if (lastResultCookie) {
                challengeResultID = lastResultCookie.split('=')[1];
            } else {
                challengeResultID = "";
            }
        } else {
            challengeID = "";
            challengeResultID = "";
        }
    }

    window.onload = injectStylesheet;
    // Sometimes, the google scripts crash on startup. Just reload the page if that happens.
    window.onerror = function(e) {
        if (e.includes("Timer")) {
            location.reload(false);
        }
    };

    let replaceStateLocal = history.replaceState;
    history.replaceState = function() {
    }

    let pushStateLocal = history.pushState;
    history.pushState = function() {
    }

    let leafletMap = null;
    let hasGuessed = false;

    let marker = null;
    function makeGuess(latlng) {
        console.log(latlng);
        if (hasGuessed) {
            return;
        }
        hasGuessed = true;
        latlng = latlng.wrap();
        let guess = {
            ChallengeResultID: challengeResultID,
            RoundNum: challengeResult.Guesses.length,
            Location: {Lat: latlng.lat, Lng: latlng.lng},
        };
        console.log(guess);
        fetch("/api/guesses", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(guess),
        }).then((response) => {
            if (response.ok) {
                window.location.replace("/scores");
            } else {
                alert("Failed to submit guess?!");
            }
        });
    }

    // The leaflet minimap!
    function createMinimap() {
        let floatingContainer = document.getElementById("leaflet-container");

        leafletMap = L.map("leaflet-map");

        // Load marker if it was previously stored (see reload button)
        let oldMarker = null;
        try {
            console.log(sessionStorage.getItem("lastMarker"));
            oldMarker = JSON.parse(sessionStorage.getItem("lastMarker"));
        } finally {
            if (oldMarker != null && oldMarker.gameID == challengeResultID && oldMarker.roundNumber == challengeResult.Guesses.length) {
                marker = L.marker(L.latLng(oldMarker.lat, oldMarker.lng));
                marker.addTo(leafletMap);
            }
            sessionStorage.removeItem("lastMarker");
        }

        // Zoom out map
        setTimeout(function() {
            if (oldMarker == null) {
                leafletMap.setView([0.0, 0.0], .1);
            } else {
                leafletMap.setView([oldMarker.lat, oldMarker.lng], 1);
            }
        }, 100);

        // only show text labels on minimap if the user wishes so
        let labelsEnabled = false; // TODO: FIXME: set from config
        // TODO: FIXME: set from config
        L.tileLayer(tileServerURL, {
            attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OSM</a> contributors, <a href="https://wikitech.wikimedia.org/wiki/Wikitech:Cloud_Services_Terms_of_use">Wikimedia Cloud Services</a>'
        }).addTo(leafletMap);

        function onMapClick(event) {
            if (marker != null) {
                leafletMap.removeControl(marker);
            }
            marker = L.marker(event.latlng);
            marker.addTo(leafletMap);
            guessButton.className = guessButton.className.replace("disabled", "");
        }

        leafletMap.on("click", onMapClick);

        setTimeout(function() {
            leafletMap.invalidateSize();
        }, 100);

        // TODO: can we move the compass without doing this?
        // Move the compass from inside the google code to the compass container.
        let compassContainer = document.getElementById("compass-container");
        let compass = document.getElementById("compass");
        compass.parentNode.removeChild(compass);
        compassContainer.appendChild(compass);
        
        // score, round number, and timer
        // TODO: implement this
        if (map.TimeLimit > 0) {
            let remainingTime = map.TimeLimit;
            setTimer = function() {
                minutes = Math.floor(remainingTime / 60);
                seconds = Math.floor(remainingTime % 60).toString();
                while (seconds.length < 2) seconds = "0" + seconds;
                let remainingTimeInfo = "Time: " + minutes + ":" + seconds;
                roundInfoSpan.innerHTML = getScoreInfo() + "<br/>" + getRoundInfo() + "<br/>" + remainingTimeInfo + "<br/>";
            }
            let interval = setInterval(function() {
                remainingTime -= 1;
                if (remainingTime <= 0) {
                    if (marker == null) {
                        makeGuess(L.latLng(0, 0));
                    } else {
                        makeGuess(marker.getLatLng());
                    }
                    clearInterval(interval);
                }
                setTimer();
            }, 1000);
        }
    }

    let sizes = [
        [150, 150],
        [300, 300],
        [500, 500],
        [800, 800],
    ];

    function scaleMap(bigger) {
        let map = document.getElementById("leaflet-container");

        let size = [map.scrollWidth, map.scrollHeight];
        let nextSize = null;

        let index = -1;
        for (el in sizes) {
            index++;
            if (sizes[el][0] == size[0]) {
                break;
            }
        }

        if (bigger) {
            index++;
        } else {
            index--;
        }

        if (index < 0) {
            index = 0;
        }
        if (index > sizes.length) {
            index = sizes.length;
        }

        map.style.width = sizes[index][0] + "px";
        map.style.height = sizes[index][1] + "px";

        leafletMap.invalidateSize();
    }
</script>

<style>
    #leaflet-map {
    }
</style>


<div id="leaflet-container">
    <div id="navigation-bar" class="btn-group btn-group-sm">
        <button class="btn btn-light" on:click={() => {scaleMap(true)}}>⬉</button>
        <button class="btn btn-light" on:click={() => {scaleMap(false)}}>⬊</button>
    </div>
    <button 
        class="btn btn-primary btn-sm float-right disabled" 
        on:click={() => {
            if (marker == null) {
                alert("You have to add a marker first! Do this by clicking the map.");
                return;
            }
            // Post data back to earthwalker.
            makeGuess(marker.getLatLng());
        }}>
        Guess!
    </button>
    <div id="leaflet-map"></div>
</div>
<div id="compass-container"></div>
<div id="round-info-container">
    <span class="round-info-span align-middle">
        {"Round: " + (challengeResult && map ? (challengeResult.Guesses.length + 1) + "/" + map.NumRounds : "loading...")}
        <br/>
        {"Total points: not implemented"}
        <br/>
    </span>
    <button 
        class="btn btn-light mx-sm-2 align-middle"
        on:click={() => {
            if (marker != null) {
                // Put marker into session storage
                console.log(marker.getLatLng());
                sessionStorage.setItem("lastMarker", JSON.stringify({
                    "lat": marker.getLatLng().lat,
                    "lng": marker.getLatLng().lng,
                    "roundNumber": challengeResult.Guesses.length,
                    "gameID": challengeResultID,
                }));
            }
            // https://www.phpied.com/files/location-location/location-location.html
            location = location;
        }}>go to start</button>
</div>