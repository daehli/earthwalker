<script>
    export let displayedResult = null;
    export let allResults;
    export let curRound;

    $: console.log(allResults);
</script>

<style>
    :global(.leaderboard tr.clickable) {
        cursor: pointer;
    }

    :global(.leaderboard tr.highlight) {
        background-color: lightblue !important;
    }
</style>

<table class="table table-striped leaderboard">
    <thead>
    <th scope="col">Icon</th>
    <th scope="col">Nickname</th>
    <th scope="col">Points</th>
    <th scope="col">Distance Off</th>
    </thead>
    <tbody>
        {#each allResults as curResult, i}
            {#if curResult.Guesses.length > curRound}
                <tr 
                    scope="row" 
                    class={displayedResult ? displayedResult.ChallengeResultID === curResult.ChallengeResultID ? "clickable highlight" : "clickable" : ""} 
                    on:click={() => {displayedResult = allResults[i];}}
                >
                    <td><img style="height: 20px;" src={svgIcon("?", curResult && curResult.Icon ? curResult.Icon : 0)}/></td>
                    <td>{curResult.Nickname}</td>
                    <td>{curResult.totalScore || (curResult.scoreDists ? curResult.scoreDists[curRound][0] : 0)}</td>
                    <td>{distString(curResult.totalDist || (curResult.scoreDists ? curResult.scoreDists[curRound][1] : 0))}</td>
                </tr>
            {/if}
        {/each}
    </tbody>
</table>