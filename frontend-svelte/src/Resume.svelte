<script>
    import {onMount} from 'svelte';
    import { loc } from './stores.js';

    export let curChallenge, curResult;
</script>

<style>
    main {
        margin: 2em;
    }
</style>

<main>
    {#if curChallenge && curResult}
        <a href={"/play?id=" + curChallenge.ChallengeID} class="btn btn-primary">Resume Game</a>
        <p>Challenge ID: <code>{curChallenge.ChallengeID}</code>, Result ID: <code>{curResult.ChallengeResultID}</code></p>
        <hr/>
    {:else}
        <p>No game in progress.</p>
    {/if}
    <p on:click={() => {$loc = "/createmap";}} class="btn btn-primary">New Map</p>
</main>