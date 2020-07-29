<script>
    import { onMount } from 'svelte';
    import { loc, ewapi, globalMap, globalChallenge, globalResult } from './stores.js';
    import LeafletGuessesMap from './components/LeafletGuessesMap.svelte';
    import Leaderboard from './components/Leaderboard.svelte';

    // data
    let allResults = [];
    let result;
    let displayedResult;

    // reactive
    let curRound = 0;
    $: [score, distance] = result ? calcScoreDistance(result.Guesses[curRound], $globalChallenge.Places[curRound], $globalMap.GraceDistance, $globalMap.Area) : [0, 0];

    async function fetchData() {
        allResults = await $ewapi.getAllResults($globalChallenge.ChallengeID);
        allResults.forEach(r => {
            r.scoreDists = r.Guesses.map((guess, i) => calcScoreDistance(guess, $globalChallenge.Places[i], $globalMap.GraceDistance, $globalMap.Area));
            r.scoreDists = r.scoreDists.concat(Array($globalMap.NumRounds - r.scoreDists.length).fill([0, 0]));
        });
        result = allResults.find(r => r.ChallengeResultID === $globalResult.ChallengeResultID);
        curRound = result.Guesses.length - 1;
        allResults.sort((a, b) => b.scoreDists[curRound][0] - a.scoreDists[curRound][0]);
        allResults = allResults;
        displayedResult = result;
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
    {#await fetchData()}
        <h2>Loading...</h2>
    {:then}
        <LeafletGuessesMap displayedResult={displayedResult} showAll={false}/>
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
                    <Leaderboard bind:displayedResult={displayedResult} {allResults} {curRound}/>
                </div>
                <p class="text-muted">Reload the page to see other player's scores once they finish this round.</p>
                {#if $globalMap.NumRounds && result && result.Guesses && result.Guesses.length == $globalMap.NumRounds}
                    <button type="button" class="btn btn-primary" on:click={() => {$loc = "/summary";}}>Go to summary</button>
                {:else}
                    <button type="button" class="btn btn-primary" on:click={() => {window.location.replace("/play?id="+$globalChallenge.ChallengeID);}}>Continue to next round</button>
                {/if}
            </div>
        </div>
    {/await}
</main>