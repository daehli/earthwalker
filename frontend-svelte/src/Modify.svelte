<script>
    import { onMount } from 'svelte';

    let ewapi = new EarthwalkerAPI();
    // data fetched from server
    let challengeID;
    let challenge;
    let map;
    let challengeResultID;
    let challengeResult;
    let tileServerURL;
    // timer
    let timeRemaining = 0;
    // map sizing
    const mapSizes = [
        [150, 150],
        [300, 300],
        [500, 500],
        [800, 800],
    ];
    let curMapSize = 1;
    // sets title to "earthwalker"
    let titleInterval;
    // decrements timeRemaining once per second
    let timerInterval;

    // DOM elements
    let floatingContainer;
    let guessButton;

    let leafletMap = null;
    let hasGuessed = false;
    let marker = null;

    // I assume these are part of the streetview sorcery, I'm not messing with them
    let replaceStateLocal = history.replaceState;
    history.replaceState = function() {
    }

    let pushStateLocal = history.pushState;
    history.pushState = function() {
    }

    onMount(async () => {
        challengeID = getChallengeID();
        if (!challengeID) {
            alert("Could not determine challenge ID!");
        }
        challengeResultID = getChallengeResultID(challengeID);
        if (challengeResultID) {
            challengeResult = await ewapi.getResult(challengeResultID);
            if (!challengeResult.Guesses) {
                challengeResult.Guesses = [];
            }
        } else {
            alert("Could not determine result ID! (How'd you get here?)");
        }
        
        challenge = await ewapi.getChallenge(challengeID);
        map = await ewapi.getMap(challenge.MapID);
        tileServerURL = (await ewapi.getTileServer(map.ShowLabels)).tileserver;
        
        titleInterval = setInterval(setTitle, 100);

        createMinimap();
    });

    // Sometimes, the google scripts crash on startup. Just reload the page if that happens.
    window.onerror = function(e) {
        if (e.includes("Timer")) {
            location.reload(false);
        }
    };

    function setTitle() {
        try {
            new MutationObserver(function(mutations) {
                if (document.title != "Earthwalker") {
                    document.title = "Earthwalker";
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
    }

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
        ewapi.postGuess(guess).then((response) => {
            if (response) {
                window.location.replace("/scores");
            } else {
                alert("Failed to submit guess?!");
            }
        });
    }

    // The leaflet minimap!
    function createMinimap() {
        leafletMap = L.map("leaflet-map");

        // Load marker if it was previously stored (see reload button)
        let oldMarker = null;
        try {
            oldMarker = JSON.parse(sessionStorage.getItem("lastMarker"));
        } finally {
            if (oldMarker != null && oldMarker.gameID == challengeResultID && oldMarker.roundNumber == challengeResult.Guesses.length) {
                marker = L.marker(L.latLng(oldMarker.lat, oldMarker.lng));
                marker.addTo(leafletMap);
                guessButton.className = guessButton.className.replace("disabled", "");
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
        // TODO: can we use an absolute timer instead of this interval?
        if (map.TimeLimit > 0) {
            timeRemaining = map.TimeLimit;
            timerInterval = setInterval(function() {
                timeRemaining -= 1;
                if (timeRemaining <= 0) {
                    if (marker == null) {
                        makeGuess(L.latLng(0, 0));
                    } else {
                        makeGuess(marker.getLatLng());
                    }
                    clearInterval(timeRemaining);
                }
            }, 1000);
        }
    }

    function scaleMap(bigger) {
        if (bigger) {
            if (curMapSize < mapSizes.length - 1) {
                curMapSize++;
            }
        } else {
            if (curMapSize > 0) {
                curMapSize--;
            }
        }

        floatingContainer.style.width = mapSizes[curMapSize][0] + "px";
        floatingContainer.style.height = mapSizes[curMapSize][1] + "px";

        leafletMap.invalidateSize();
    }
</script>

<style>
</style>

<div bind:this={floatingContainer} id="leaflet-container">
    <div id="navigation-bar" class="btn-group btn-group-sm">
        <button class="btn btn-light" on:click={() => {scaleMap(true)}}>⬉</button>
        <button class="btn btn-light" on:click={() => {scaleMap(false)}}>⬊</button>
    </div>
    <button 
        bind:this={guessButton}
        class="btn btn-primary btn-sm float-right disabled" 
        on:click={() => {
            if (marker == null) {
                alert("You have to add a marker first! Do this by clicking the map.");
                return;
            }
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
        {#if timeRemaining > 0}
            Time: {Math.floor(timeRemaining / 60).toString()}:{Math.floor(timeRemaining % 60).toString().padStart(2, '0')}
            <br/>
        {/if}
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