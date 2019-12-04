package challenge

import (
	"gitlab.com/glatteis/earthwalker/player"
	"gitlab.com/glatteis/earthwalker/streetviewserver"
	"log"
	"net/http"
)

// ServeChallenge serves a challenge to the user (using the /game?c= url).
func ServeChallenge(w http.ResponseWriter, r *http.Request) {
	challengeKey, ok := r.URL.Query()["c"]
	// This is probably what they call "user error"
	if !ok || len(challengeKey) == 0 {
		http.Redirect(w, r, "/", 302)
		return
	}
	actualKey := challengeKey[0]

	var cookieNotPresent bool
	session, err := player.GetSessionFromCookie(r)
	if err != nil || session.GameID != actualKey {
		if err != player.PlayerSessionNotFoundError {
			log.Println(err)
		}
		session = player.NewSession()
		cookieNotPresent = true
	}

	foundChallenge, err := GetChallenge(actualKey)
	if err == ChallengeNotFoundError {
		w.Write([]byte("this challenge does not exist!"))
		w.WriteHeader(404)
		return
	} else if err != nil {
		log.Println(err)
		w.Write([]byte("there was some kind of internal error, sorry!"))
		w.WriteHeader(500)
		return
	}

	session.GameID = foundChallenge.UniqueIdentifier
	round := session.Round()

	err = player.StorePlayerSession(session)
	if err != nil {
		log.Println("Could not save a session: " + err.Error())
		w.Write([]byte("Could not save your session: " + err.Error()))
		return
	}
	if cookieNotPresent {
		log.Println("Setting a new cookie for", session)
		player.SetSessionCookie(session, w)
	}

	streetviewserver.ServeLocation(foundChallenge.Places[round-1], w, r)
}
