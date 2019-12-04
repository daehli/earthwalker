package challenge

import (
	"bytes"
	"encoding/gob"
	"errors"
	"github.com/dgraph-io/badger"
	"github.com/golang/geo/s2"
	"gitlab.com/glatteis/earthwalker/database"
	"math/rand"
)

// A Challenge represents a number of places along with all kinds of associated data.
type Challenge struct {
	Places           []s2.LatLng
	UniqueIdentifier string
	Settings         ChallengeSettings
}

// The ChallengeSettings contain user-configurable options about the game.
type ChallengeSettings struct {
	NumRounds      int
	LabeledMinimap bool
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func NewChallenge(places []s2.LatLng) (Challenge, error) {
	challenge := Challenge{
		Places:           places,
		UniqueIdentifier: randSeq(5),
		// TODO make channel settings configurable
		Settings: ChallengeSettings{
			LabeledMinimap: true,
			NumRounds:      5,
		},
	}

	err := database.GetDB().Update(func(txn *badger.Txn) error {
		var buffer bytes.Buffer
		gob.NewEncoder(&buffer).Encode(challenge)
		return txn.Set([]byte("challenge-"+challenge.UniqueIdentifier), buffer.Bytes())
	})

	if err != nil {
		return Challenge{}, err
	}
	return challenge, nil
}

// The ChallengeNotFoundError is the error that is returned by GetChallenge when no challenge
// of that id is present.
var ChallengeNotFoundError = errors.New("challenge not found!")

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
		return Challenge{}, ChallengeNotFoundError
	} else if err != nil {
		return Challenge{}, err
	}

	var foundChallenge Challenge
	gob.NewDecoder(bytes.NewBuffer(challengeBytes)).Decode(&foundChallenge)

	return foundChallenge, nil
}
