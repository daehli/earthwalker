package badgerdb

import (
	"bytes"
	"encoding/gob"
	"fmt"

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
