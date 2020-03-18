// Package database opens and closes the database, and holds its object.
package database

import (
	"github.com/dgraph-io/badger"
	"gitlab.com/glatteis/earthwalker/util"
	"log"
)

var database *badger.DB

func init() {
	var err error
	database, err = badger.Open(badger.DefaultOptions(util.AppPath() + "/badger/"))
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
