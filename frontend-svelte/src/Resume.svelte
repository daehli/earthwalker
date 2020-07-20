<script>
    import { onMount } from 'svelte';
    import { loc, ewapi, globalMap, globalChallenge, globalResult } from './stores.js';

    // these functions are in every component, because it's easier that way
    // TODO: FIXME: a cleaner way with no race conditions.  
    //       Derived stores with promises/callbacks?
    async function setResultChallengeMap(resultID) {
        $globalResult = await $ewapi.getResult(resultID);
        if (!$globalChallenge || $globalResult.ChallengeID !== $globalChallenge.ChallengeID) {
            return setChallengeMap($globalResult.ChallengeID);
        }
    }

    async function setChallengeMap(challengeID) {
        $globalChallenge = await $ewapi.getChallenge(challengeID);
        if (!$globalMap || $globalChallenge.MapID !== $globalMap.MapID) {
            $globalMap = await $ewapi.getMap($globalChallenge.MapID);
        }
    }

    onMount(async () => {
        if (!$globalResult) {
            let challengeID = getChallengeID();
            if (challengeID) {
                await setResultChallengeMap(getChallengeResultID(challengeID));
            }
        }
    });
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