<script>
    import {onMount} from 'svelte';

    const challengeCookieName = "earthwalker_lastChallenge";
    const resultCookiePrefix = "earthwalker_lastResult_";

    onMount(async () => {
        let params = new URLSearchParams(window.location.search)
        // TODO: consider having default map settings if there's no ID
        if (!params.has("id")) {
            alert("URL has no challenge ID!");
            return;
        }
        challengeID = params.get("id");
    });

    let challengeID;

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
        let challengeResult = JSON.stringify({
            ChallengeID: challengeID,
            Nickname: document.getElementById("Nickname").value,
        });
        let response = await fetch("api/results", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: challengeResult,
        });
        let data = await response.json();
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
                    <input required type="text" class="form-control" id="Nickname"/>
                </div>
            </div>

            <button id="submit-button" class="btn btn-primary" style="margin-bottom: 2em; color: #fff;">Start Challenge</button>

        </div>
    </form>
</main>