package badgerdb

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strconv"

	"github.com/dgraph-io/badger"
	"gitlab.com/glatteis/earthwalker/domain"
)

// TODO: simpler method, this logic was lifted straight from challenge.go
func setBytesFunc(key string, t interface{}) func(*badger.Txn) error {
	return func(txn *badger.Txn) error {
		var buffer bytes.Buffer
		gob.NewEncoder(&buffer).Encode(t)
		return txn.Set([]byte(key), buffer.Bytes())
	}
}

func getBytesFunc(key string, bytes []byte) func(*badger.Txn) error {
	return func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			bytes = append([]byte{}, val...)
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	}
}

func Init(path string) (*badger.DB, error) {
	db, err := badger.Open(badger.DefaultOptions(path))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Close(db *badger.DB) {
	db.Close()
}

type MapStore struct {
	DB *badger.DB
}

func (store MapStore) Insert(m domain.Map) error {
	err := store.DB.Update(setBytesFunc("map-"+m.MapID, m))
	if err != nil {
		return fmt.Errorf("Failed to write map to badger DB: %v\n", err)
	}
	return nil
}

// TODO: reduce code repetition in Get methods
func (store MapStore) Get(mapID string) (domain.Map, error) {
	var mapBytes []byte
	err := store.DB.View(getBytesFunc("map-"+mapID, mapBytes))
	if err != nil {
		return domain.Map{}, fmt.Errorf("Failed to read map from badger DB: %v\n", err)
	}

	var foundMap domain.Map
	err = gob.NewDecoder(bytes.NewBuffer(mapBytes)).Decode(&foundMap)
	if err != nil {
		return domain.Map{}, fmt.Errorf("Failed to decode map from bytes: %v\n", err)
	}
	return foundMap, nil
}

type ChallengeStore struct {
	DB *badger.DB
}

func (store ChallengeStore) Insert(c domain.Challenge) error {
	err := store.DB.Update(setBytesFunc("challenge-"+c.ChallengeID, c))
	if err != nil {
		return fmt.Errorf("Failed to write challenge to badger DB: %v\n", err)
	}
	return nil
}

func (store ChallengeStore) Get(challengeID string) (domain.Challenge, error) {
	var challengeBytes []byte
	err := store.DB.View(getBytesFunc("challenge-"+challengeID, challengeBytes))
	if err != nil {
		return domain.Challenge{}, fmt.Errorf("Failed to read challenge from badger DB: %v\n", err)
	}

	var foundChallenge domain.Challenge
	err = gob.NewDecoder(bytes.NewBuffer(challengeBytes)).Decode(&foundChallenge)
	if err != nil {
		return domain.Challenge{}, fmt.Errorf("Failed to decode challenge from bytes: %v\n", err)
	}
	return foundChallenge, nil
}

type ChallengePlaceStore struct {
	DB *badger.DB
}

func (store ChallengePlaceStore) Insert(p domain.ChallengePlace) error {
	err := store.DB.Update(setBytesFunc("challengePlace-"+p.ChallengeID+strconv.Itoa(p.RoundNum), p))
	if err != nil {
		return fmt.Errorf("Failed to write ChallengePlace to badger DB: %v\n", err)
	}
	return nil
}
