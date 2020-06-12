// Package player handles player objects and serves stuff based on that.
package player

import (
	"bytes"
	"encoding/gob"
	"errors"
	"github.com/dgraph-io/badger"
	"gitlab.com/glatteis/earthwalker/database"
	"math/rand"
	"net/http"
	"time"
)

// A Session stores the data of a player
type Session struct {
	// UniqueIdentifier is the session identifier stored in the key.
	UniqueIdentifier string
	// The Nickname the player gives themselves.
	Nickname string
	// GameID is game identifier the player might be currently in.
	GameID string
	// Points is the number of points of rounds that the player has already completed.
	Points []int
	// Distances are the respective distances as floats.
	Distances []float64
	// GuessedPositions are the guessed positions as float64 tuples ([lat, lng])
	GuessedPositions [][]float64
	// TimeStarted is the time that a player started a specific round from the challenge.
	// If the player hasn't started the round yet, this will be nil.
	TimeStarted *time.Time
	// IconColor is the player's icon color (see static/icons).
	IconColor int
}

// Round returns the round the player is currently in
func (p Session) Round() int {
	return 1 + len(p.Points)
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// NewSession creates a new player session
func NewSession() Session {
	return Session{
		UniqueIdentifier: randSeq(10),
		IconColor:        (rand.Int() % 200) + 1,
	}
}

// StorePlayerSession stores a player session in the database.
func StorePlayerSession(session Session) error {
	err := database.GetDB().Update(func(txn *badger.Txn) error {
		var buffer bytes.Buffer
		gob.NewEncoder(&buffer).Encode(session)
		return txn.Set([]byte("session-"+session.UniqueIdentifier), buffer.Bytes())
	})

	if err != nil {
		return err
	}
	return nil
}

// RemovePlayerSession removes a player session in the database.
func RemovePlayerSession(session Session) error {
	err := database.GetDB().Update(func(txn *badger.Txn) error {
		var buffer bytes.Buffer
		gob.NewEncoder(&buffer).Encode(session)
		return txn.Delete([]byte("session-" + session.UniqueIdentifier))
	})

	if err != nil {
		return err
	}
	return nil
}

// ErrPlayerSessionDoesNotExist is the error that is thrown when a player does not exist
var ErrPlayerSessionDoesNotExist = errors.New("this player does not exist")

// LoadPlayerSession loads a player session from the database.
// You should probably use GetSessionFromCookie in cookies.go.
func LoadPlayerSession(id string) (Session, error) {
	var playerBytes []byte

	err := database.GetDB().Update(func(txn *badger.Txn) error {
		result, err := txn.Get([]byte("session-" + id))
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

		playerBytes = res
		return nil
	})

	if err == badger.ErrKeyNotFound {
		return Session{}, ErrPlayerSessionDoesNotExist
	} else if err != nil {
		return Session{}, err
	}

	var foundSession Session
	gob.NewDecoder(bytes.NewBuffer(playerBytes)).Decode(&foundSession)

	return foundSession, nil
}

// WriteNicknameAndSession writes a nickname and a session if the session
// does not exist yet, otherwise writes the nickname to the session.
// Only returns an error if it is exceptional.
func WriteNicknameAndSession(w http.ResponseWriter, r *http.Request, nickname string) error {
	session, err := GetSessionFromCookie(r)

	var writeSession bool
	if err != nil {
		if err != ErrPlayerSessionNotFound && err != ErrPlayerSessionDoesNotExist {
			return err
		}
		session = NewSession()
		writeSession = true
	}

	if session.Nickname != nickname {
		session.Nickname = nickname
		writeSession = true
	}

	if writeSession {
		err := StorePlayerSession(session)
		if err != nil {
			return err
		}
	}

	SetSessionCookie(session, w)

	return nil
}
