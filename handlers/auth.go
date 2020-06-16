// handles session management
package handlers

import (
	"bytes"
	"encoding/gob"
	"errors"
	"math/rand"
	"net/http"
	"time"

	"github.com/dgraph-io/badger"
)

// ErrPlayerSessionNotFound is the error that occurs when no player session is found.
var ErrPlayerSessionNotFound = errors.New("no player session found")

// SetSessionCookie sets the session cookie of a session into the browser.
func SetSessionCookie(session Session, w http.ResponseWriter) {
	c := http.Cookie{
		Name:   "earthwalker-session",
		Value:  session.UniqueIdentifier,
		MaxAge: int((24 * time.Hour).Seconds()),
	}
	http.SetCookie(w, &c)
}

// GetSessionFromCookie retrieves the cookie from a session
func GetSessionFromCookie(r *http.Request) (Session, error) {
	var cookie *http.Cookie
	for _, c := range r.Cookies() {
		if c.Name == "earthwalker-session" {
			cookie = c
		}
	}
	if cookie == nil {
		return Session{}, ErrPlayerSessionNotFound
	}

	session, err := LoadPlayerSession(cookie.Value)
	if err != nil {
		return Session{}, err
	}

	return session, nil
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
