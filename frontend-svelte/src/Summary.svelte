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
    let allResults = [];
    let map;

    let guessLocs;
    let actualLocs;
    let scoreDists = [];

    // leaflet
    let scoreMap;
    let scoreMapPolyGroup;
    let scoreMapGuessGroup;

    onMount(async () => {
        challenge = await ewapi.getChallenge(challengeID);
        allResults = await ewapi.getAllResults(challengeID);
        map = await ewapi.getMap(challenge.MapID);

        actualLocs = challenge.Places.map((place) => place.Location);
        allResults.forEach(r => {
            r.scoreDists = r.Guesses.map((guess, i) => calcScoreDistance(guess.Location.Lat, guess.Location.Lng, actualLocs[i].Lat, actualLocs[i].Lng, map.GraceDistance, map.Area));
            r.scoreDists = r.scoreDists.concat(Array(map.NumRounds - r.scoreDists.length).fill([0, 0]));
            r.totalScore = r.scoreDists.reduce((acc, val) => acc + val[0], 0);
            r.totalDist = r.scoreDists.reduce((acc, val) => acc + val[1], 0)
        });
        result = allResults.find(r => r.ChallengeResultID === challengeResultID);
        allResults.sort((a, b) => b.totalScore - a.totalScore);
        allResults = allResults;

        setupScoreMap();
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
        if (map.Polygon) {
            showPolygonOnMap(scoreMapPolyGroup, map.Polygon);
        }
        scoreMapGuessGroup = L.layerGroup().addTo(scoreMap);
        showGuessesOnMap();
    }

    function showGuessesOnMap() {
        scoreMapGuessGroup.clearLayers();
        result.Guesses.forEach((guess, i) => {
            showGuessOnMap(scoreMapGuessGroup, guess.Location, actualLocs[i], i, result.Nickname, result.Icon);
        });
    }

    function showChallengeLinkPrompt() {
        let link = window.location.origin + "/play?id=" + challengeID;
	    window.prompt("Copy to clipboard: Ctrl+C, Enter", link);
    }

    function switchToResult(index) {
        result = allResults[index];
        showGuessesOnMap();
    }
</script>

<style>
    #leaderboard tr {
        cursor: pointer;
    }
</style>

<!-- This prevents users who haven't finished the challenge from viewing
     TODO: cleaner protection for this page -->
{#if result && map && result.Guesses.length == map.NumRounds}
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

    <div style="margin-top: 2em; text-align: center;">
	<h3>{result && result.Nickname ? result.Nickname + "\'s" : "Your"} scores:</h3>
	<table class="table table-striped">
		<thead>
		<th scope="col">Round</th>
		<th scope="col">Points</th>
		<th scope="col">Distance Off</th>
		</thead>
		<tbody>
        {#if result && result.scoreDists}
            {#each result.scoreDists as scoreDist, i}
                <tr scope="row">
                    <td>{i + 1}</td>
                    <td>{scoreDist[0]}</td>
                    <td>{distString(scoreDist[1])}</td>
                </tr>
            {/each}
        {/if}
		</tbody>
	</table>
    </div>

    <div id="leaderboard" style="margin-top: 2em; text-align: center;">
		<h3>Leaderboard</h3>
        <table class="table table-striped">
            <thead>
            <th scope="col">Icon</th>
            <th scope="col">Nickname</th>
            <th scope="col">Number of Points</th>
            <th scope="col">Total Distance Off</th>
            </thead>
            <tbody>
                {#each allResults as curResult, i}
                    {#if curResult.Guesses.length == map.NumRounds}
                        <tr scope="row" on:click={() => {switchToResult(i);}}>
                            <td><img style="height: 20px;" src={svgIcon("?", curResult && curResult.Icon ? curResult.Icon : 0)}/></td>
                            <td>{curResult.Nickname}</td>
                            <td>{curResult.totalScore}</td>
                            <td>{distString(curResult.totalDist)}</td>
                        </tr>
                    {/if}
                {/each}
            </tbody>
        </table>
    </div>
</div>
{:else}
    <div class="text-center">
        <h3>Loading...</h3>
        <h3>You must finish the game to view this page.</h3>
    </div>
{/if}