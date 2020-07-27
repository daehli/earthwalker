<script>
    // TODO: most of this script is duplicated in Scores.svelte.
    //       (also a bit in Modify.svelte)
    //       consolidate.
    import {onMount} from 'svelte';
    import { loc, ewapi, globalMap, globalChallenge, globalResult } from './stores.js';
    import LeafletGuessesMap from './LeafletGuessesMap.svelte';

    let displayedResult;
    let allResults = [];

    let guessLocs;
    let actualLocs;
    let scoreDists = [];

    // leaflet
    let scoreMap;
    let scoreMapPolyGroup;
    let scoreMapGuessGroup;

    async function fetchData() {
        allResults = await $ewapi.getAllResults($globalChallenge.ChallengeID);
        allResults.forEach(r => {
            r.scoreDists = r.Guesses.map((guess, i) => calcScoreDistance(guess, $globalChallenge.Places[i], $globalMap.GraceDistance, $globalMap.Area));
            r.scoreDists = r.scoreDists.concat(Array($globalMap.NumRounds - r.scoreDists.length).fill([0, 0]));
            r.totalScore = r.scoreDists.reduce((acc, val) => acc + val[0], 0);
            r.totalDist = r.scoreDists.reduce((acc, val) => acc + val[1], 0)
        });
        displayedResult = allResults.find(r => r.ChallengeResultID === $globalResult.ChallengeResultID);
        allResults.sort((a, b) => b.totalScore - a.totalScore);
        allResults = allResults;
    }
</script>

<style>
    #leaderboard tr {
        cursor: pointer;
    }
</style>

<!-- This prevents users who haven't finished the challenge from viewing
     TODO: cleaner protection for this page -->
{#if $globalResult.Guesses && $globalMap.NumRounds && $globalResult.Guesses.length == $globalMap.NumRounds}
    {#await fetchData()}
        <h2>Loading...</h2>
    {:then}
        <LeafletGuessesMap {displayedResult} showAll={true}/>

        <div class="container">
            <br>
            <div class="row">
                <div class="col text-center">
                    <button type="button" id="copy-game-link" class="btn btn-primary" on:click={() => showChallengeLinkPrompt($globalChallenge.ChallengeID)}>
                        Copy link to this game
                    </button>
                </div>
            </div>

            <div style="margin-top: 2em; text-align: center;">
                <h3>{displayedResult && displayedResult.Nickname ? displayedResult.Nickname + "\'s" : "Your"} scores:</h3>
                <table class="table table-striped">
                    <thead>
                    <th scope="col">Round</th>
                    <th scope="col">Points</th>
                    <th scope="col">Distance Off</th>
                    </thead>
                    <tbody>
                    {#if displayedResult && displayedResult.scoreDists}
                        {#each displayedResult.scoreDists as scoreDist, i}
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
                        {#each allResults as result, i}
                            {#if result.Guesses.length == $globalMap.NumRounds}
                                <tr scope="row" on:click={() => {displayedResult = allResults[i];}}>
                                    <td><img style="height: 20px;" src={svgIcon("?", result && result.Icon ? result.Icon : 0)}/></td>
                                    <td>{result.Nickname}</td>
                                    <td>{result.totalScore}</td>
                                    <td>{distString(result.totalDist)}</td>
                                </tr>
                            {/if}
                        {/each}
                    </tbody>
                </table>
            </div>
        </div>
    {/await}
{:else}
    <div class="text-center">
        <h2>You must finish the game to view this page.</h2>
    </div>
{/if}