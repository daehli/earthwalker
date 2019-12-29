package challenge

import (
	"github.com/pkg/errors"
	"gitlab.com/glatteis/earthwalker/player"
	"gitlab.com/glatteis/earthwalker/streetviewserver"
	"log"
	"net/http"
	"time"
)

// ServeChallenge serves a challenge to the user (using the /game?c= url).
func ServeChallenge(w http.ResponseWriter, r *http.Request) {
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

// WriteNicknameAndSession writes a nickname and a session if the session
// does not exist yet, otherwise writes the nickname to the session.
// Only returns an error if it is exceptional.
func WriteNicknameAndSession(w http.ResponseWriter, r *http.Request, nickname string) error {
	session, err := player.GetSessionFromCookie(r)

	var writeSession bool
	if err != nil {
		if err != player.ErrPlayerSessionNotFound {
			return err
		}
		session = player.NewSession()
		writeSession = true
	}

	var writeSessionCookie bool
	if session.Nickname != nickname {
		session.Nickname = nickname
		writeSessionCookie = true
	}

	if writeSession {
		err := player.StorePlayerSession(session)
		if err != nil {
			return err
		}
	}

	if writeSessionCookie {
		player.SetSessionCookie(session, w)
	}

	return nil
}
