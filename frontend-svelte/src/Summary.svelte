<script>
    // TODO: most of this script is duplicated in Scores.svelte.
    //       (also a bit in Modify.svelte)
    //       consolidate.
        import {onMount} from 'svelte';

    let ewapi = new EarthwalkerAPI();
    let challengeID = getChallengeID();
    let challengeResultID = getChallengeResultID(challengeID);
    let challenge;
    let result;
    let map;

    let guessLocs;
    let actualLocs;
    let scoreDists = [];

    // leaflet
    let scoreMap;
    let scoreMapPolyGroup;

    onMount(async () => {
        challenge = await ewapi.getChallenge(challengeID);
        result = await ewapi.getResult(challengeResultID);
        map = await ewapi.getMap(challenge.MapID);
        // TODO: FIXME: this code assumes Guesses and challenge.Places are 
        //              ordered, which the API does not guarantee
        guessLocs = result.Guesses.map((guess) => guess.Location);
        actualLocs = challenge.Places.map((place) => place.Location);
        
        setupScoreMap();
        // TODO: FIXME: this code assumes Guesses and challenge.Places are 
        //              ordered, which the API does not guarantee
        scoreDists = guessLocs.map((guessLoc, i) => 
            calcScoreDistance(guessLoc.Lat, guessLoc.Lng, actualLocs[i].Lat, actualLocs[i].Lng, map.GraceDistance, map.Area));
        console.log(scoreDists);
    });

    async function setupScoreMap() {
        let tileServer = (await ewapi.getTileServer()).tileserver;
        console.log(tileServer);
        scoreMap = L.map("score-map");
        console.log(scoreMap);
        scoreMap.setView([0.0, 0.0], 1);
        L.tileLayer(tileServer, {
            attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OSM</a> contributors, <a href="https://wikitech.wikimedia.org/wiki/Wikitech:Cloud_Services_Terms_of_use">Wikimedia Cloud Servides</a>'
        }).addTo(scoreMap);
        scoreMapPolyGroup = L.layerGroup().addTo(scoreMap);
        showPolygonOnMap();
        showGuessesOnMap();
    }

    function showPolygonOnMap() {
        scoreMapPolyGroup.clearLayers();
        let map_poly = L.geoJSON(map.Polygon).addTo(scoreMapPolyGroup);
        scoreMap.fitBounds(map_poly.getBounds());
    }

    function showGuessesOnMap() {
        guessLocs.forEach((loc, i) => {
            showGuessOnMap(loc, actualLocs[i]);
        });
    }

    // TODO: show results from other users 
    //       (kv db not really suited to this, maybe switch to relational)
    function showGuessOnMap(guessLoc, actualLoc) {
        let polyline = L.polyline([[guessLoc.Lat, guessLoc.Lng], [actualLoc.Lat, actualLoc.Lng]], {color: '#007bff'}).addTo(scoreMap);
        L.marker([guessLoc.Lat, guessLoc.Lng], {
            title: result.Nickname,
            icon: makeIcon("question" + result.Icon + ".png"),
        }).addTo(scoreMap).openPopup();
        L.marker([actualLoc.Lat, actualLoc.Lng], {
            title: "Actual Position",
            icon: makeIcon("answer.png"),
        }).addTo(scoreMap).openPopup();
    }

    let makeIcon = function(name) {
        return L.icon({
        iconUrl: "public/icons/" + name,
        iconSize: [50/2, 82/2],
        iconAnchor: [25/2, 82/2],
        shadowUrl: "public/leaflet/images/marker-shadow.png",
        shadowSize: [41, 41],
        shadowAnchor: [12, 41]
        });
    };

    function showChallengeLinkPrompt() {
        let link = window.location.origin + "/play?id=" + challengeID;
	    window.prompt("Copy to clipboard: Ctrl+C, Enter", link);
    }
</script>


<div id="score-map" style="width: 100%; height: 50vh;"></div>

<div class="container">
    <br>
    <div class="row">
        <div class="col text-center">
            <button id="copy-game-link" class="btn btn-primary" on:click={showChallengeLinkPrompt}>
                Copy link to this game
            </button>
        </div>
    </div>

    <div style="margin-top: 10%; text-align: center;">
	<h3>Your scores:</h3>
	<table class="table table-striped">
		<thead>
		<th scope="col">Round</th>
		<th scope="col">Points</th>
		<th scope="col">Distance Off</th>
		</thead>
		<tbody>
        {#if result && result.Guesses}
            {#each scoreDists as scoreDist, i}
                <tr scope="row">
                    <td>{i + 1}</td>
                    <td>{scoreDist[0]}</td>
                    <td>{scoreDist[1]}</td>
                </tr>
            {/each}
        {/if}
		</tbody>
	</table>
    </div>

    <!-- TODO: implement challengeResult rankings
    <div style="margin-top: 10%; text-align: center;">
		<h3>Leaderboard</h3>
	<table class="table table-striped">
	    <thead>
		<th scope="col">Icon</th>
		<th scope="col">Nickname</th>
		<th scope="col">Number of Points</th>
		<th scope="col">Total Distance Off</th>
	    </thead>
	    <tbody>
	      
	    </tbody>
	</table>
    </div>-->
</div>