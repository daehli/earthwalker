// handlers in this file serve the actual Earthwalker game
// similar to former challenge/server.go
package handlers

import (
	"log"
	"math"
	"net/http"
	"time"

	"github.com/golang/geo/s2"
	"github.com/pkg/errors"
	"gitlab.com/glatteis/earthwalker/player"
	"gitlab.com/glatteis/earthwalker/streetviewserver"
)

type ChallengeHandler struct {
}

// ServeChallenge serves a challenge to the user (using the /game?c= url).
func (handler *ChallengeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	challengeKey, ok := r.URL.Query()["c"]
	// This is probably what they call "user error"
	if !ok || len(challengeKey) == 0 {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	actualKey := challengeKey[0]

	session, err := player.GetSessionFromCookie(r)
	if err == player.ErrPlayerSessionNotFound || err == player.ErrPlayerSessionDoesNotExist {
		http.Redirect(w, r, "/set_nickname?c="+actualKey, http.StatusFound)
		return
	} else if err != nil {
		log.Println(err)
		http.Error(w, "Some internal error occured, sorry!", http.StatusUnprocessableEntity)
		return
	}

	foundChallenge, err := GetChallenge(actualKey)
	if err == ErrChallengeNotFound {
		http.Error(w, "this challenge does not exist!", http.StatusNotFound)
		return
	} else if err != nil {
		log.Println(err)
		http.Error(w, "there was some kind of internal error, sorry!", http.StatusUnprocessableEntity)
		return
	}

	if session.GameID != actualKey {
		newSession := player.NewSession()
		newSession.Nickname = session.Nickname
		session = newSession
		err := player.StorePlayerSession(session)
		if err != nil {
			log.Println(err)
			http.Error(w, "there was some kind of internal error, sorry!", http.StatusUnprocessableEntity)
			return
		}
		player.SetSessionCookie(session, w)
	}

	session.GameID = foundChallenge.UniqueIdentifier
	round := session.Round()

	if round > len(foundChallenge.Places) {
		http.Redirect(w, r, "/summary", http.StatusFound)
		return
	}

	if session.TimeStarted == nil {
		now := time.Now()
		session.TimeStarted = &now
	}

	err = player.StorePlayerSession(session)
	if err != nil {
		log.Println(errors.Wrap(err, "could not save a session"))
		http.Error(w, "could not save your session ", http.StatusInternalServerError)
		return
	}

	streetviewserver.ServeLocation(foundChallenge.Places[round-1], w, r)
}

// == Scoring ========

const earthRadius = 6371000
const maxScore = 5000
const graceDistance = 25 // perfect scores will be given within this distance (meters)
// effectively, score will be divided by decayBase every decayDistance meters
const decayBase = 2
const decayDistance = 1070000

// CalculateScoreAndDistance calculates the score and distance for a guessed location (and its actual location).
func CalculateScoreAndDistance(actualLocation s2.LatLng, guessLocation s2.LatLng) (int, float64) {
	distance := actualLocation.Distance(guessLocation).Radians() * earthRadius
	if distance < graceDistance {
		return maxScore, distance
	}
	factor := math.Pow(decayBase, -float64(distance)/decayDistance)
	points := int(factor * maxScore)
	return points, distance
}
