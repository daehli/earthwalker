<script>
    import {onMount} from 'svelte';
    import { loc } from './stores.js';

    const challengeCookieName = "earthwalker_lastChallenge";
    const resultCookiePrefix = "earthwalker_lastResult_";

    let ewapi = new EarthwalkerAPI();
    let challengeID;
    let nickname = "";

    onMount(async () => {
        challengeID = getChallengeID();
    });

    // TODO: this duplicates a function in CreateChallenge.
    //       consider consolidating.
    async function handleFormSubmit() {
        let challengeResultID = await submitNewChallengeResult();
        // set the generated challenge as the current challenge
        document.cookie = challengeCookieName + "=" + challengeID + ";path=/;max-age=172800";
        // set the generated ChallengeResult as the current ChallengeResult
        // for the Challenge with challengeID
        document.cookie = resultCookiePrefix + challengeID + "=" + challengeResultID + ";path=/;max-age=172800";
        window.location.replace("/play");
    }

    // TODO: this duplicates a function in CreateChallenge.
    //       consolidate to api lib
    async function submitNewChallengeResult() {
        let challengeResult = {
            ChallengeID: challengeID,
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
        <p>Challenge ID: <code>{challengeID}</code></p>
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