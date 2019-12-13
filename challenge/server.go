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

	session, err := player.GetSessionFromCookie(r)
	if err == player.PlayerSessionNotFoundError {
		http.Redirect(w, r, "/set_nickname?c="+actualKey, 302)
		return
	} else if err != nil {
		log.Println(err)
		w.Write([]byte("Some internal error occured, sorry!"))
		w.WriteHeader(500)
		return
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

	if session.GameID != actualKey {
		newSession := player.NewSession()
		newSession.Nickname = session.Nickname
		session = newSession
		err := player.StorePlayerSession(session)
		if err != nil {
			log.Println(err)
			w.Write([]byte("there was some kind of internal error, sorry!"))
			w.WriteHeader(500)
			return
		}
		player.SetSessionCookie(session, w)
	}

	session.GameID = foundChallenge.UniqueIdentifier
	round := session.Round()

	if round > len(foundChallenge.Places) {
		http.Redirect(w, r, "/summary", 302)
		return
	}

	err = player.StorePlayerSession(session)
	if err != nil {
		log.Println("Could not save a session: " + err.Error())
		w.Write([]byte("Could not save your session: " + err.Error()))
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
		if err != player.PlayerSessionNotFoundError {
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
