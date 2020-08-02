<script>
    import { onMount } from 'svelte';
    import { loc, ewapi, globalMap, globalResult } from './stores.js';
    import MapInfo from './components/MapInfo.svelte';

    // TODO: FIXME: clean up/use async properly
    // TODO: improve page design, looks rather cluttered right now

    const challengeCookieName = "earthwalker_lastChallenge";
    const resultCookiePrefix = "earthwalker_lastResult_";

    let statusText = "Twiddling thumbs...";

    let streetViewService = new google.maps.StreetViewService();

    const popTIFLoc = "/public/nasa_pop_data.tif";
    let popTIF = undefined;

    let challengeID;
    let numSVReqs = 0;
    let numFound = 0;
    $: foundProportion = numFound / numSVReqs;
    let foundCoords = [];
    let done = false;

    // DOM elements
    let submitButton;
    // bindings
    let nickname = "";

    onMount(async () => {
        statusText = "Getting Map settings...";
        let mapid = getURLParam("mapid");
        if (mapid && (!$globalMap || $globalMap.MapID !== mapid)) {
            $globalMap = await $ewapi.getMap(mapid);
        } else if (!mapid) {
            alert("No Map ID in URL!");
            statusText = "Please make sure there is a valid Map ID in the URL."
            return;
        }

        statusText = "Looking up population density data...";
        popTIF = await loadGeoTIF(popTIFLoc);

        statusText = "Fetching panoramas...";
        foundCoords = await fetchPanos(
            streetViewService, 
            $globalMap, 
            popTIF, 
            (panoWasFound) => {
                if (panoWasFound) {
                    numFound++;
                }
                numSVReqs++;
            }
        );

        statusText = "Sending Challenge to server..."
        challengeID = await submitNewChallenge();

        statusText = "Done!";
        done = true;
    });

    async function submitNewChallenge() {
        let convertedCoords = foundCoords.map((coord, i) => ({RoundNum: i, Location: {Lat: coord.lat(), Lng: coord.lng()}}));
        let challenge = {
            MapID: $globalMap.MapID,
            Places: convertedCoords
        };
        let data = await $ewapi.postChallenge(challenge);
        return data.ChallengeID;
    }

    async function handleFormSubmit() {
        // TODO: duplicates code in Join
        let challengeResultID = await submitNewChallengeResult();
        $globalResult = await $ewapi.getResult(challengeResultID);
        // set the generated challenge as the current challenge
        document.cookie = challengeCookieName + "=" + challengeID + ";path=/;max-age=172800";
        // set the generated ChallengeResult as the current ChallengeResult
        // for the Challenge with challengeID
        document.cookie = resultCookiePrefix + challengeID + "=" + challengeResultID + ";path=/;max-age=172800";
        window.location.replace("/play");
    }

    async function submitNewChallengeResult() {
        let challengeResult = {
            ChallengeID: challengeID,
            Nickname: nickname,
        };
        let data = await $ewapi.postResult(challengeResult);
        return data.ChallengeResultID;
    }
</script>

<main>
    <form on:submit|preventDefault={handleFormSubmit} class="container">
        <br>
        <h2>New Challenge - {$globalMap ? $globalMap.Name : "Loading..."}</h2>
        <MapInfo/>
        <br>
        <h4 class="text-center" id="status">{statusText}</h4>
        <div action="" method="post">
            <div class="progress">
                <div 
                    style={"width: " + ($globalMap ? (100 * numFound) / $globalMap.NumRounds : 0) + "%;"} 
                    class="progress-bar" id="loading-progress" role="progressbar">
                    {numFound + "/" + ($globalMap ? $globalMap.NumRounds : 0)}
                </div>
            </div>
            <small class="text-muted">
                Earthwalker is getting random locations from Google Street View.
                This happens in your browser, because there is only an API in JavaScript for this.
                Yes, that is kind of silly, I'm sorry.
            </small>
            <br/>
            <small class="text-muted">
                {numFound} locations found after {numSVReqs} StreetView API requests ({Math.round(foundProportion * 100)}% success rate).
            </small>
            {#if foundProportion && foundProportion < 0.3}
                <small class="text-muted">
                    To find locations more quickly, you may want to consider reducing the specificity of your map settings.
                </small>
            {/if}
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
                    <button type="button" id="copy-game-link" class="btn btn-primary" on:click={(e) => {showChallengeLinkPrompt(challengeID);}}>
                        Copy link to this game
                    </button>
                {/if}
            </div>
        </div>
    </form>
</main>
