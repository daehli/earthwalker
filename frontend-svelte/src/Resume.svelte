<script>
    // TODO: consolidate cookie management to separate script
    const challengeCookieName = "earthwalker_lastChallenge"
    const resultCookiePrefix = "earthwalker_lastResult_"

    // ID of most recently created or played challenge
    let lastChallengeID;
    // ID of most recent result for challenge with ID lastChallengeID (from cookie)
    let lastResultID;

    // TODO: this reactive statement doesn't actually react
    $: parseCookies(document.cookie)

    function parseCookies(cookieStr) {
        console.log("re-parsing cookies");
        let cookies = cookieStr.split("; ");
        let lastChallengeCookie = cookies.find(row => row.startsWith(challengeCookieName));
        if (lastChallengeCookie) {
            lastChallengeID = lastChallengeCookie.split('=')[1];
            let lastResultCookie = cookies.find(row => row.startsWith(resultCookiePrefix + lastChallengeID));
            if (lastResultCookie) {
                lastResultID = lastResultCookie.split('=')[1];
            } else {
                lastResultID = "";
            }
        } else {
            lastChallengeID = "";
            lastResultID = "";
        }
    }
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
    <a href="/createmap" class="btn btn-primary">New Map</a>
</main>