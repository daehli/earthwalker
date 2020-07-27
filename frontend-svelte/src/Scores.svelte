<script>
    import { onMount } from 'svelte';
    import { loc, ewapi, globalMap, globalChallenge, globalResult } from './stores.js';
    import LeafletGuessesMap from './LeafletGuessesMap.svelte';

    // data
    let allResults = [];
    let result;

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
        <LeafletGuessesMap displayedResult={result} showAll={false}/>
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
                {#if $globalMap.NumRounds && result && result.Guesses && result.Guesses.length == $globalMap.NumRounds}
                    <button type="button" class="btn btn-primary" on:click={() => {$loc = "/summary";}}>Go to summary</button>
                {:else}
                    <button type="button" class="btn btn-primary" on:click={() => {window.location.replace("/play");}}>Continue to next round</button>
                {/if}
            </div>
        </div>
    {/await}
</main>