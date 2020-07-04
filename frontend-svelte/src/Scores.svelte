<script>
    import {onMount} from 'svelte';

    let ewapi = new EarthwalkerAPI();
    let challengeID = getChallengeID();
    let challengeResultID = getChallengeResultID(challengeID);
    let challenge;
    let result;
    let map;

    let lastGuess;
    let lastActual;
    let score = 0;
    let distance = 0;

    // leaflet
    let scoreMap;
    let scoreMapPolyGroup;

    onMount(async () => {
        challenge = await ewapi.getChallenge(challengeID);
        result = await ewapi.getResult(challengeResultID);
        map = await ewapi.getMap(challenge.MapID);
        console.log(map);
        // TODO: FIXME: this code assumes Guesses and challenge.Places are 
        //              ordered, which the API does not guarantee
        lastGuess = result.Guesses[result.Guesses.length - 1].Location;
        lastActual = challenge.Places[result.Guesses.length - 1].Location;
        [score, distance] = calcScoreDistance(lastGuess.Lat, lastGuess.Lng, lastActual.Lat, lastActual.Lng, map.GraceDistance, map.Area);
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
            showPolygonOnMap();
        }
        showGuessOnMap(scoreMap, lastGuess, lastActual, result.Guesses.length - 1, result.Nickname, result.Icon, true);
        showLastGuessOnMap();
    }

    // TODO: this function is duplicated in CreateMap.  Consolidate.
    function showPolygonOnMap() {
        scoreMapPolyGroup.clearLayers();
        let map_poly = L.geoJSON(map.Polygon).addTo(scoreMapPolyGroup);
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
        <div style="margin-top: 10%; text-align: center;">
            <p class="text-center">
                You were {distString(distance)} from the correct position. Your marker is 
                <img 
                    style="height: 40px;" 
                    alt={result && result.Nickname ? result.Nickname : "Your Icon"}
                    src={svgIcon("?", result && result.Icon ? result.Icon : 0)}/>
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
            <p class="text-muted"><strike>Reload the page to see other player's scores once they finish this round.</strike></p>
            <p class="text-muted">Other player scores not yet implemented.</p>
            {#if map && map.NumRounds && result && result.Guesses && result.Guesses.length == map.NumRounds}
                <button type="button" class="btn btn-primary" onclick="window.location.href = '/summary'">Go to summary</button>
            {:else}
                <button type="button" class="btn btn-primary" onclick="window.location.href = '/play'">Continue to next round</button>
            {/if}
        </div>
    </div>
</main>