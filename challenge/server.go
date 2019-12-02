package challenge

import (
	"bytes"
	"encoding/gob"
	"github.com/dgraph-io/badger"
	"gitlab.com/glatteis/earthwalker/database"
	"gitlab.com/glatteis/earthwalker/player"
	"gitlab.com/glatteis/earthwalker/streetviewserver"
	"log"
	"net/http"
	"strconv"
)

// ServeChallenge serves a challenge to the user (using the /game?c= url).
func ServeChallenge(w http.ResponseWriter, r *http.Request) {
	challengeKey, ok := r.URL.Query()["c"]
	// This is probably what they call "user error"
	if !ok || len(challengeKey) == 0 {
		http.Redirect(w, r, "/", 302)
		return
	}

	round, ok := r.URL.Query()["round"]
	var roundAsString string
	if !ok || len(round) == 0 {
		roundAsString = "1"
	} else {
		roundAsString = round[0]
	}

	roundAsInt, err := strconv.Atoi(roundAsString)
	if err != nil {
		w.Write([]byte("Round must be an integer!"))
		return
	}

	actualKey := challengeKey[0]
	var challengeBytes []byte

	err = database.GetDB().Update(func(txn *badger.Txn) error {
		result, err := txn.Get([]byte("challenge-" + actualKey))
		if err != nil {
			return err
		}

		var res []byte
		err = result.Value(func(val []byte) error {
			res = append([]byte{}, val...)
			return nil
		})

		if err != nil {
			return err
		}

		challengeBytes = res
		return nil
	})

	if err == badger.ErrKeyNotFound {
		w.Write([]byte("Sorry, that challenge does not exist!"))
		return
	} else if err != nil {
		log.Println(err)
		w.Write([]byte("Internal error, sorry! Please contact an administrator or something."))
		return
	}

	var foundChallenge Challenge
	gob.NewDecoder(bytes.NewBuffer(challengeBytes)).Decode(&foundChallenge)

	if roundAsInt < 1 || roundAsInt > foundChallenge.Settings.NumRounds {
		w.Write([]byte("Round number not in range!"))
		return
	}

	var cookieNotPresent bool
	session, err := player.GetSessionFromCookie(r)
	if err != nil {
		session = player.NewSession()
		cookieNotPresent = true
	}
	session.CurrentGameID = foundChallenge.UniqueIdentifier
	session.CurrentRound = roundAsInt
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

	streetviewserver.ServeLocation(foundChallenge.Places[roundAsInt-1], w, r)
}
