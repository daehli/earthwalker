<script>
    import {onMount} from 'svelte';
    import { loc } from './stores.js';

    // data
    let ewapi = new EarthwalkerAPI();
    let challengeID = getChallengeID();
    let challengeResultID = getChallengeResultID(challengeID);
    let challenge;
    let allResults = [];
    let result;
    let map;

    // leaflet
    let scoreMap;
    let scoreMapPolyGroup;

    // reactive
    let dataLoaded = false;
    let scoreMapLoaded = false;
    let curRound = 0;
    $: lastGuess = dataLoaded ? result.Guesses[curRound].Location : [0, 0];
    $: lastActual = dataLoaded ? challenge.Places[curRound].Location : [0, 0];
    $: [score, distance] = dataLoaded ? calcScoreDistance(lastGuess.Lat, lastGuess.Lng, lastActual.Lat, lastActual.Lng, map.GraceDistance, map.Area) : [0, 0];
    $: if (scoreMapLoaded) {showGuessOnMap(scoreMap, lastGuess, lastActual, curRound, result.Nickname, result.Icon, true);}

    onMount(async () => {
        allResults = await ewapi.getAllResults(challengeID);
        challenge = await ewapi.getChallenge(challengeID);
        map = await ewapi.getMap(challenge.MapID);
        allResults.forEach(r => {
            r.scoreDists = r.Guesses.map((guess, i) => calcScoreDistance(guess.Location.Lat, guess.Location.Lng, challenge.Places[i].Location.Lat, challenge.Places[i].Location.Lng, map.GraceDistance, map.Area));
            r.scoreDists = r.scoreDists.concat(Array(map.NumRounds - r.scoreDists.length).fill([0, 0]));
        });
        result = allResults.find(r => r.ChallengeResultID === challengeResultID);
        curRound = result.Guesses.length - 1;
        allResults.sort((a, b) => b.scoreDists[curRound][0] - a.scoreDists[curRound][0]);
        allResults = allResults;
        dataLoaded = true;
        setupScoreMap();
    });

    async function setupScoreMap() {
        let tileServer = (await ewapi.getTileServer()).tileserver;
        scoreMap = L.map("score-map");
        scoreMap.setView([0.0, 0.0], 1);
        L.tileLayer(tileServer, {
            attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OSM</a> contributors, <a href="https://wikitech.wikimedia.org/wiki/Wikitech:Cloud_Services_Terms_of_use">Wikimedia Cloud Servides</a>'
        }).addTo(scoreMap);
        scoreMapPolyGroup = L.layerGroup().addTo(scoreMap);
        if (map.Polygon) {
            showPolygonOnMap(scoreMapPolyGroup, map.Polygon);
        }
        scoreMapLoaded = true;
    }
</script>

<style>
    #nopoints {
        text-align: center;
        width: 100%;
        margin: auto;
    }
</style>

<main>
    <div id="score-map" style="width: 100%; height: 50vh;"></div>
    <div class="container">
        <div style="margin-top: 2em; text-align: center;">
            <p class="text-center">
                You were {distString(distance)} from the correct position. Your marker is 
                {#if result && result.Guesses && result.Icon}
                <img 
                    style="height: 40px;" 
                    alt={result && result.Nickname ? result.Nickname : "Your Icon"}
                    src={svgIcon((curRound + 1).toString(), result.Icon)}/>
                {/if}
            </p>
            <div class="progress" style="height: 40px;">
                <div 
                    class="progress-bar" 
                    role="progressbar" 
                    style="width: {score / 50}%;" 
                    aria-valuenow={score.toString()} 
                    aria-valuemin="0" 
                    aria-valuemax="5000">
                    {score ? (score + "/5000 points") : ""}
                </div>
                {#if !score}
                    <div id="nopoints">Sorry, you didn't receive any points for this round.</div>
                {/if}
            </div>
            <div id="leaderboard" style="margin-top: 2em; text-align: center;">
                <h3>Leaderboard</h3>
                <table class="table table-striped">
                    <thead>
                    <th scope="col">Icon</th>
                    <th scope="col">Nickname</th>
                    <th scope="col">Points</th>
                    <th scope="col">Distance Off</th>
                    </thead>
                    <tbody>
                        {#each allResults as curResult, i}
                            {#if curResult.Guesses.length > curRound}
                                <tr scope="row">
                                    <td><img style="height: 20px;" src={svgIcon("?", curResult && curResult.Icon ? curResult.Icon : 0)}/></td>
                                    <td>{curResult.Nickname}</td>
                                    <td>{curResult.scoreDists ? curResult.scoreDists[curRound][0] : 0}</td>
                                    <td>{distString(curResult.scoreDists ?curResult.scoreDists[curRound][1] : 0)}</td>
                                </tr>
                            {/if}
                        {/each}
                    </tbody>
                </table>
            </div>
            <p class="text-muted">Reload the page to see other player's scores once they finish this round.</p>
            {#if map && map.NumRounds && result && result.Guesses && result.Guesses.length == map.NumRounds}
                <button type="button" class="btn btn-primary" on:click={() => {$loc = "/summary";}}>Go to summary</button>
            {:else}
                <button type="button" class="btn btn-primary" on:click={() => {window.location.replace("/play");}}>Continue to next round</button>
            {/if}
        </div>
    </div>
</main>