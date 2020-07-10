<script>
    import {onMount} from 'svelte';
    import { loc } from './stores.js';

    // ID of most recently created or played challenge
    let lastChallengeID;
    // ID of most recent result for challenge with ID lastChallengeID (from cookie)
    let lastResultID;

    onMount(async () => {
        lastChallengeID = getChallengeID();
        lastResultID = getChallengeResultID(lastChallengeID);
    });
</script>

<style>
    main {
        margin: 2em;
    }
</style>

<main>
    {#if lastResultID}
        <a href={"/play?id=" + lastResultID} class="btn btn-primary">Resume Game</a>
        <p>Challenge ID: <code>{lastChallengeID}</code>, Result ID: <code>{lastResultID}</code></p>
        <hr/>
    {:else}
        <p>No game in progress.</p>
    {/if}
    <p on:click={() => {$loc = "/createmap";}} class="btn btn-primary">New Map</p>
</main>