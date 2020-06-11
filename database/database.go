// Package database opens and closes the database, and holds its object.
package database

import (
	"log"
	"os"

	"github.com/dgraph-io/badger"
	"gitlab.com/glatteis/earthwalker/util"
)

var database *badger.DB

func init() {
	path, err := getDBPath()
	if err != nil {
		log.Fatal(err)
	}
	database, err = badger.Open(badger.DefaultOptions(path))
	if err != nil {
		log.Fatal(err)
	}
}

func getDBPath() (string, error) {
	path := ""
	pathSuffix := os.Getenv("EARTHWALKER_DB_PATH")
	pathRel := os.Getenv("EARTHWALKER_DB_PATH_REL")
	if pathRel == "cwd" {
		cwd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		path = cwd + pathSuffix
	} else if pathRel == "absolute" {
		path = pathSuffix
	} else { // default: relative to executable
		path = util.AppPath() + pathSuffix
	}
	return path, nil
}

// CloseDB closes the datbase. This is there to be deferred in the main function.
func CloseDB() {
	database.Close()
}

// GetDB gets the database.
func GetDB() *badger.DB {
	return database
}
