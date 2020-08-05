<script>
    import { loc, globalChallenge, globalResult } from './stores.js';
</script>

<main>
    {#if $globalChallenge && $globalResult}
        <a href={"/play?id=" + $globalChallenge.ChallengeID} class="btn btn-primary">Resume Game</a>
        <p>Challenge ID: <code>{$globalChallenge.ChallengeID}</code>, Result ID: <code>{$globalResult.ChallengeResultID}</code></p>
        <hr/>
    {:else}
        <p>No game in progress.</p>
    {/if}
    <p on:click={() => {$loc = "/createmap";}} class="btn btn-primary">New Map</p>
</main>