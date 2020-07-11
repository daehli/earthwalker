<script>
    import {onMount} from 'svelte';
    import { loc } from './stores.js';

    export let ewapi, curChallenge, curResult;

    const challengeCookieName = "earthwalker_lastChallenge";
    const resultCookiePrefix = "earthwalker_lastResult_";

    let nickname = "";

    // TODO: this duplicates a function in CreateChallenge.
    //       consider consolidating.
    async function handleFormSubmit() {
        curResult = await ewapi.getResult(await submitNewChallengeResult());
        // set the generated challenge as the current challenge
        document.cookie = challengeCookieName + "=" + curChallenge.ChallengeID + ";path=/;max-age=172800";
        // set the generated ChallengeResult as the current ChallengeResult
        // for the Challenge with challengeID
        document.cookie = resultCookiePrefix + curChallenge.ChallengeID + "=" + curResult.ChallengeResultID + ";path=/;max-age=172800";
        window.location.replace("/play");
    }

    // TODO: this duplicates a function in CreateChallenge.
    //       consolidate to api lib
    async function submitNewChallengeResult() {
        let challengeResult = {
            ChallengeID: curChallenge.ChallengeID,
            Nickname: nickname,
        };
        let data = await ewapi.postResult(challengeResult);
        return data.ChallengeResultID;
    }

</script>

<main>
    <form on:submit|preventDefault={handleFormSubmit} class="container">
        <br>
        <h2>Join Challenge</h2>
        <p>Challenge ID: <code>{curChallenge.ChallengeID}</code></p>
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