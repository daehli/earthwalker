package badgerdb

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/dgraph-io/badger"
	"gitlab.com/glatteis/earthwalker/domain"
)

// TODO: make store and get more symmetrical?
func storeStruct(db *badger.DB, key string, t interface{}) error {
	err := db.Update(func(txn *badger.Txn) error {
		var buffer bytes.Buffer
		gob.Register(map[string]interface{}{})
		gob.Register([]interface{}{})
		err := gob.NewEncoder(&buffer).Encode(t)
		if err != nil {
			return err
		}
		return txn.Set([]byte(key), buffer.Bytes())
	})
	return err
}

func getBytes(db *badger.DB, key string) ([]byte, error) {
	var byteSlice []byte
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			byteSlice = append([]byte{}, val...)
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})
	return byteSlice, err
}

// Init opens and returns a badger database connection
// don't forget to close it
func Init(path string) (*badger.DB, error) {
	db, err := badger.Open(badger.DefaultOptions(path))
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Close closes the given badger database connection
// (provided so you don't have to import badger just to do this)
func Close(db *badger.DB) {
	db.Close()
}

// MapStore badger implementation (see domain)
type MapStore struct {
	DB *badger.DB
}

const mapPrefix = "map-"

// Insert a domain.Map into store's badger db
func (store MapStore) Insert(m domain.Map) error {
	err := storeStruct(store.DB, mapPrefix+m.MapID, m)
	if err != nil {
		return fmt.Errorf("failed to write map to badger DB: %v", err)
	}
	return nil
}

// Get a domain.Map with the given mapID from store's badger db
// TODO: reduce code repetition in Get methods
func (store MapStore) Get(mapID string) (domain.Map, error) {
	mapBytes, err := getBytes(store.DB, mapPrefix+mapID)
	if err != nil || len(mapBytes) == 0 {
		return domain.Map{}, fmt.Errorf("failed to read map from badger DB: %v", err)
	}

	var foundMap domain.Map
	gob.Register(foundMap)
	err = gob.NewDecoder(bytes.NewBuffer(mapBytes)).Decode(&foundMap)
	if err != nil {
		return domain.Map{}, fmt.Errorf("failed to decode map from bytes: %v", err)
	}
	return foundMap, nil
}

// ChallengeStore badger implementation (see domain)
type ChallengeStore struct {
	DB *badger.DB
}

const challengePrefix = "challenge-"

// Insert a domain.Challenge into store's badger db
func (store ChallengeStore) Insert(c domain.Challenge) error {
	err := storeStruct(store.DB, challengePrefix+c.ChallengeID, c)
	if err != nil {
		return fmt.Errorf("failed to write challenge to badger DB: %v", err)
	}
	return nil
}

// Get a domain.Challenge with the given challengeID from store's badger db
func (store ChallengeStore) Get(challengeID string) (domain.Challenge, error) {
	var challengeBytes []byte
	challengeBytes, err := getBytes(store.DB, challengePrefix+challengeID)
	if err != nil {
		return domain.Challenge{}, fmt.Errorf("failed to read challenge from badger DB: %v", err)
	}

	var foundChallenge domain.Challenge
	err = gob.NewDecoder(bytes.NewBuffer(challengeBytes)).Decode(&foundChallenge)
	if err != nil {
		return domain.Challenge{}, fmt.Errorf("failed to decode challenge from bytes: %v", err)
	}
	return foundChallenge, nil
}

// note: no ChallengePlaceStore implementation,
// because we just store the entire Challenge as a blob
// one will probably be necessary for relational databases
// (which don't take well to arbitrary length fields)

// ChallengeResultStore badger implementation (see domain)
type ChallengeResultStore struct {
	DB *badger.DB
}

const challengeResultPrefix = "result-"

// Insert a domain.ChallengeResult into store's badger db
func (store ChallengeResultStore) Insert(r domain.ChallengeResult) error {
	err := storeStruct(store.DB, challengeResultPrefix+r.ChallengeResultID, r)
	if err != nil {
		return fmt.Errorf("failed to write challenge result to badger DB: %v", err)
	}
	return nil
}

// Get a domain.ChallengeResult with the given challengeResultID from store's badger db
func (store ChallengeResultStore) Get(challengeResultID string) (domain.ChallengeResult, error) {
	resultBytes, err := getBytes(store.DB, challengeResultPrefix+challengeResultID)
	if err != nil {
		return domain.ChallengeResult{}, fmt.Errorf("failed to read result from badger DB: %v", err)
	}

	var foundResult domain.ChallengeResult
	err = gob.NewDecoder(bytes.NewBuffer(resultBytes)).Decode(&foundResult)
	if err != nil {
		return domain.ChallengeResult{}, fmt.Errorf("failed to decode result from bytes: %v", err)
	}
	return foundResult, nil
}

// GetAll is not implemented
// TODO: implement this if necessary, otherwise consider removing it from the interface
func (store ChallengeResultStore) GetAll(challengeID string) ([]domain.ChallengeResult, error) {
	return make([]domain.ChallengeResult, 0), fmt.Errorf("ChallengeResultStore.GetAll is not implemented")
}

// note: as above, no GuessStore implementation,
// because we just store the entire ChallengeResult as a blob
