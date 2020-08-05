<script>
    import { globalMap } from '../stores.js';
</script>

<main>
    {#if $globalMap}
        <p>
            This Map consists of {$globalMap.NumRounds} rounds, 
            {#if $globalMap.TimeLimit > 0}
                each with a time limit of {Math.floor($globalMap.TimeLimit / 60)}:{Math.floor($globalMap.TimeLimit % 60).toString().padStart(2, '0')}.  
            {:else}
                with no time limit.  
            {/if}
            Full points will be awarded for guesses within {$globalMap.GraceDistance} meters.
        </p>
    {:else}
        <p>Loading Map...</p>
    {/if}
</main>
