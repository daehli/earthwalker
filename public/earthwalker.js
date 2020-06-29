// == common functions ========

const challengeCookieName = "earthwalker_lastChallenge";
const resultCookiePrefix = "earthwalker_lastResult_";

// return value of url param with key, else null
function getURLParam(key) {
    let params = new URLSearchParams(window.location.search)
    if (!params.has(key)) {
        return;
    }
    return params.get(key);
}

// getCookieValue with specified name, else null
function getCookieValue(name) {
    let cookies = document.cookie.split("; ");
    let cookie = cookies.find(row => row.startsWith(name));
    if (cookie) {
        return cookie.split('=')[1];
    }
    return null;
}

// == JS API layer ========

// helpers

// gets object from the given URL, else null
async function getObject(url) {
    let response = await fetch(url);
    if (response.ok) {
        return response.json();
    }
    return null
}

// posts object to the given URL, returns response object else null
async function postObject(url, object) {
    let response = await fetch(url, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(object),
    });
    if (response.ok) {
        return response.json();
    }
    return null
}

// methods return promises
class EarthwalkerAPI {
    constructor(baseURL="") {
        this.configURL = baseURL + "/api/config";
        this.mapsURL = baseURL + "/api/maps";
        this.challengesURL = baseURL + "/api/challenges";
        this.resultsURL = baseURL + "/api/results";
        this.guessesURL = baseURL + "/api/guesses";
    }

    // get tile server url (as object) from server, nolabel if specified
    getTileServer(labeled=true) {
        return getObject(this.configURL + (labeled ? "/tileserver" : "/nolabeltileserver"))
    }

    // get map object from server by id
    getMap(mapID) {
        return getObject(this.mapsURL+"/"+mapID);
    }

    // post new map object to server
    postMap(map) {
        return postObject(this.mapsURL, map);
    }

    getChallenge(challengeID) {
        return getObject(this.challengesURL+"/"+challengeID);
    }

    postChallenge(challenge) {
        return postObject(this.challengesURL, challenge);
    }

    getResult(resultID) {
        return getObject(this.resultsURL+"/"+resultID);
    }

    postResult(result) {
        return postObject(this.resultsURL, result);
    }

    postGuess(guess) {
        return postObject(this.guessesURL, guess);
    }

    // I shoved some extra stuff in here because I don't understand javascript imports

    // getChallengeID from the URL (key: "id"), else get the value of cookie
    // lastChallenge, else null
    getChallengeID() {
        let id = getURLParam("id");
        if (id) {
            return id;
        }
        return getCookieValue(challengeCookieName);
    }

    // getChallengeResultID from cookie resultCookiePrefix+challengeID, else null
    getChallengeResultID(challengeID) {
        return getCookieValue(resultCookiePrefix+challengeID);
    }
}