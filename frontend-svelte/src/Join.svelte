<script>
    import { onMount } from 'svelte';
    import { loc, ewapi, globalChallenge, globalResult } from './stores.js';

    const challengeCookieName = "earthwalker_lastChallenge";
    const resultCookiePrefix = "earthwalker_lastResult_";

    let nickname = "";

     onMount(async () => {
         $globalChallenge = await $ewapi.getChallenge(getChallengeID());
     });

    // TODO: this duplicates a function in CreateChallenge.
    //       consider consolidating.
    async function handleFormSubmit() {
        $globalResult = await $ewapi.getResult(await submitNewChallengeResult());
        // set the generated challenge as the current challenge
        document.cookie = challengeCookieName + "=" + $globalChallenge.ChallengeID + ";path=/;max-age=172800";
        // set the generated ChallengeResult as the current ChallengeResult
        // for the Challenge with challengeID
        document.cookie = resultCookiePrefix + $globalChallenge.ChallengeID + "=" + $globalResult.ChallengeResultID + ";path=/;max-age=172800";
        window.location.replace("/play");
    }

    // TODO: this duplicates a function in CreateChallenge.
    //       consolidate to api lib
    async function submitNewChallengeResult() {
        let challengeResult = {
            ChallengeID: $globalChallenge.ChallengeID,
            Nickname: nickname,
        };
        let data = await $ewapi.postResult(challengeResult);
        return data.ChallengeResultID;
    }

</script>

<main>
    <form on:submit|preventDefault={handleFormSubmit} class="container">
        <br>
        <h2>Join Challenge</h2>
        <p>Challenge ID: <code>{$globalChallenge ? $globalChallenge.ChallengeID : "Loading..."}</code></p>
        <div action="">
            <!-- TODO: show map settings -->
            <div class="form-group">
                <div class="input-group">
                    <div class="input-group-prepend">
                        <div class="input-group-text">Your Nickname</div>
                    </div>
                    <input bind:value={nickname} required type="text" class="form-control" id="Nickname"/>
                </div>
            </div>

            <button id="submit-button" class="btn btn-primary" style="margin-bottom: 2em; color: #fff;">Start Challenge</button>

        </div>
    </form>
</main>