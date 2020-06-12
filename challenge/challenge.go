// Package challenge handles challenge datatypes, and some serving based on those.
package challenge

import (
	"bytes"
	"encoding/gob"
	"errors"
	"github.com/dgraph-io/badger"
	"github.com/golang/geo/s2"
	"gitlab.com/glatteis/earthwalker/database"
	"math/rand"
	"time"
)

// A Challenge represents a number of places along with all kinds of associated data.
type Challenge struct {
	Places           []s2.LatLng
	Guesses          [][]Guess
	UniqueIdentifier string
	SummaryPassword  string // The "password" behind the link for the summary.
	Settings         Settings
}

// The Settings contain user-configurable options about the game.
type Settings struct {
	NumRounds      int
	TimerDuration  *time.Duration // the timer duration per round. Nil if there is no timer set.
	ConnectedOnly  bool           // include only connected panoramas
	LabeledMinimap bool
}

// A Guess contains a guess on a round from someone.
// This already contains the user's nickname so that the user
// doesn't have to be looked up from the database every time.
type Guess struct {
	GuessLocation  s2.LatLng
	PlayerID       string
	PlayerNickname string
	PlayerColor    int
}

// NewChallenge creates a new challenge with the parameters and stores it.
func NewChallenge(places []s2.LatLng, settings Settings) (Challenge, error) {
	randSeq := func(n int) string {
		var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
		b := make([]rune, n)
		for i := range b {
			b[i] = letters[rand.Intn(len(letters))]
		}
		return string(b)
	}

	challenge := Challenge{
		Places:           places,
		UniqueIdentifier: randSeq(5),
		Guesses:          make([][]Guess, len(places)),
		SummaryPassword:  randSeq(3),
		Settings:         settings,
	}

	err := StoreChallenge(challenge)

	if err != nil {
		return Challenge{}, err
	}
	return challenge, nil
}

// StoreChallenge stores changes to a challenge in the database.
func StoreChallenge(challenge Challenge) error {
	return database.GetDB().Update(func(txn *badger.Txn) error {
		var buffer bytes.Buffer
		gob.NewEncoder(&buffer).Encode(challenge)
		return txn.Set([]byte("challenge-"+challenge.UniqueIdentifier), buffer.Bytes())
	})
}

// ErrChallengeNotFound is the error that is returned by GetChallenge when no challenge
// of that id is present.
var ErrChallengeNotFound = errors.New("challenge not found")

// GetChallenge loads a challenge from an id.
func GetChallenge(id string) (Challenge, error) {
	var challengeBytes []byte

	err := database.GetDB().Update(func(txn *badger.Txn) error {
		result, err := txn.Get([]byte("challenge-" + id))
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
		return Challenge{}, ErrChallengeNotFound
	} else if err != nil {
		return Challenge{}, err
	}

	var foundChallenge Challenge
	gob.NewDecoder(bytes.NewBuffer(challengeBytes)).Decode(&foundChallenge)

	return foundChallenge, nil
}
