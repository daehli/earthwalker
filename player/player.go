package player

import (
	"bytes"
	"encoding/gob"
	"errors"
	"github.com/dgraph-io/badger"
	"gitlab.com/glatteis/earthwalker/database"
	"math/rand"
	"time"
)

type PlayerSession struct {
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
	// GuessedPositions are the guessed positions as float64 tuples
	GuessedPositions [][]float64
	// TimeLeft is the time the player had left when earthwalker last checked.
	TimeLeft time.Time
}

func (p PlayerSession) Round() int {
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

func NewSession() PlayerSession {
	return PlayerSession{
		UniqueIdentifier: randSeq(10),
	}
}

// StorePlayerSession stores a player session in the database.
func StorePlayerSession(session PlayerSession) error {
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
func RemovePlayerSession(session PlayerSession) error {
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

// loadPlayerSession loads a player session from the database.
// Is private, because you should use GetSessionFromCookie in cookies.go.
func loadPlayerSession(id string) (PlayerSession, error) {
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
		return PlayerSession{}, errors.New("this player does not exist")
	} else if err != nil {
		return PlayerSession{}, err
	}

	var foundSession PlayerSession
	gob.NewDecoder(bytes.NewBuffer(playerBytes)).Decode(&foundSession)

	return foundSession, nil
}
