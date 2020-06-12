// Package database opens and closes the database, and holds its object.
package database

import (
	"log"

	"github.com/dgraph-io/badger"
	"gitlab.com/glatteis/earthwalker/config"
)

var database *badger.DB

func init() {
	path := config.Env.EarthwalkerDBPath
	var err error
	database, err = badger.Open(badger.DefaultOptions(path))
	if err != nil {
		log.Fatal(err)
	}
}

// CloseDB closes the datbase. This is there to be deferred in the main function.
func CloseDB() {
	database.Close()
}

// GetDB gets the database.
func GetDB() *badger.DB {
	return database
}
