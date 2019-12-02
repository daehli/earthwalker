package challenge

import (
	"bytes"
	"encoding/gob"
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
