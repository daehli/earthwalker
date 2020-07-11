// This file contains code useful across the application, including wrappers
// for the database API.
// TODO: consider dividing this into multiple files.

// == common functions ========

const challengeCookieName = "earthwalker_lastChallenge";
const resultCookiePrefix = "earthwalker_lastResult_";

// getChallengeID from the URL (key: "id"), else get the value of cookie
// lastChallenge, else null
function getChallengeID() {
    let id = getURLParam("id");
    if (id) {
        return id;
    }
    return getCookieValue(challengeCookieName);
}

// getChallengeResultID from cookie resultCookiePrefix+challengeID, else null
function getChallengeResultID(challengeID) {
    return getCookieValue(resultCookiePrefix+challengeID);
}

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

function showChallengeLinkPrompt(challengeID) {
    let link = window.location.origin + "/play?id=" + challengeID;
    window.prompt("Copy to clipboard: Ctrl+C, Enter", link);
}


// == Leaflet Map ========

// 0 <= hue int < 360
function showGuessOnMap(map, guess, actual, roundNum, nickname, hue, focus=false) {
    let polyline = L.polyline([[guess.Location.Lat, guess.Location.Lng], [actual.Location.Lat, actual.Location.Lng]], {color: '#007bff'}).addTo(map);
    L.marker([guess.Location.Lat, guess.Location.Lng], {
        title: nickname,
        icon: makeIcon(roundNum + 1, hue),
    }).addTo(map).openPopup();
    L.marker([actual.Location.Lat, actual.Location.Lng], {
        title: "Actual Position",
        icon: makeIcon("!", hue),
    }).addTo(map).openPopup();
    if (focus) {
        map.fitBounds(polyline.getBounds(), {padding: [20, 20]});
    }
}

let makeIcon = function(text, hue) {
    return L.icon({
    iconUrl: svgIcon(text, hue),
    iconSize: [48, 48],
    iconAnchor: [24, 44],
    shadowUrl: "public/leaflet/images/marker-shadow.png",
    shadowSize: [41, 41],
    shadowAnchor: [12, 41]
    });
};

function svgIcon(text, hue) {
    return `data:image/svg+xml,
    <svg xmlns="http://www.w3.org/2000/svg" height="48px" viewBox="0 0 24 24" width="48px">
        <path fill="hsl(${hue}, 90%, 40%)" stroke="black" stroke-width="0.5px" d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0"/>
        <text x="50%" y="50%" dominant-baseline="middle" text-anchor="middle" font-family="'sans-serif'" fill="white" font-size="0.8em">
            ${text}
        </text>
    </svg>`
}

function showPolygonOnMap(layer, polygon) {
    layer.clearLayers();
    return L.geoJSON(polygon).addTo(layer);
}


// == Scoring ========
// TODO: tweak scoring consts

// distances in meters
const earthRadius = 6371009;
const earthArea = 510066000000000
const earthSqrt = 22584640;
const maxScore = 5000;
// score is divided by decayBase every halfDistance meters (if area=earthArea)
const decayBase = 2;
const halfDistance = 1000000;

// [score, distance] given location of guess and pano, graceDistance, and Polygon area
function calcScoreDistance(guess, actual, graceDistance=0, area=earthArea) {
    // TODO: cleaner handling of maps with no Polygon (maybe give maps area earthArea on creation?)
    if (area == 0) {
        area = earthArea;
    }
    // consider the guess invalid and return a score of zero
    if (Math.abs(guess.Location.Lat > 90)) {
        return 0
    }
    let guessPoint = turf.point([guess.Location.Lng, guess.Location.Lat]);
    let actualPoint =  turf.point([actual.Location.Lng, actual.Location.Lat]);
    let distance = turf.distance(guessPoint, actualPoint, {units: "kilometers"}) * 1000.0;
    if (distance < graceDistance) {
        return [maxScore, distance];
    }
    let relativeArea = Math.sqrt(area) / earthSqrt;
    let factor = Math.pow(decayBase, -1 * (distance - graceDistance) / (halfDistance * relativeArea));
    return [Math.round(factor * maxScore), distance];
}

// totalScore given _ordered_ arrays of {Lat, Lng}.
// actuals must be at least as long as guesses
function calcTotalScore(guesses, actuals, graceDistance=0, area=earthArea) {
    let totalScore = 0;
    let currentScore;
    guesses.forEach((guess, i) => {
        [currentScore, _] = calcScoreDistance(guess, actuals[i], graceDistance, area);
        totalScore += currentScore;
    });
    return totalScore;
}

// returns a prettified distance given float meters
function distString(meters) {
    if (meters < 1000) {
        return (meters).toFixed(1) + " m";
    }
    return (meters / 1000).toFixed(1) + " km";
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

function orderRounds(arrWithRoundNums) {
    return arrWithRoundNums.sort((a, b) => a.RoundNum - b.RoundNum);
}

// methods return promises
class EarthwalkerAPI {
    constructor(baseURL="") {
        this.configURL = baseURL + "/api/config";
        this.mapsURL = baseURL + "/api/maps";
        this.challengesURL = baseURL + "/api/challenges";
        this.resultsURL = baseURL + "/api/results";
        this.allResultsURL = baseURL + "/api/results/all";
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

    async getChallenge(challengeID) {
        let challenge = await getObject(this.challengesURL+"/"+challengeID);
        if (challenge.Places) {
            challenge.Places = orderRounds(challenge.Places);
        } else {
            challenge.Places = [];
        }
        return challenge;
    }

    postChallenge(challenge) {
        return postObject(this.challengesURL, challenge);
    }

    async getResult(resultID) {
        let result = await getObject(this.resultsURL+"/"+resultID);
        if (result.Guesses) {
            result.Guesses = orderRounds(result.Guesses);
        } else {
            result.Guesses = [];
        }
        return result;
    }

    async getAllResults(challengeID) {
        let results = await getObject(this.allResultsURL+"/"+challengeID);
        results.forEach(result => {
            if (result.Guesses) {
                result.Guesses = orderRounds(result.Guesses);
            } else {
                result.Guesses = [];
            }
        });
        return results;
    }

    postResult(result) {
        return postObject(this.resultsURL, result);
    }

    postGuess(guess) {
        return postObject(this.guessesURL, guess);
    }
}