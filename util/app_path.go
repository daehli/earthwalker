// Package util contains some utilities.
package util

import (
	"log"
	"os"
	"path"
)

// AppPath gets the executable's path.
func AppPath() string {
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal("App path not accessible!")
	}
	return path.Dir(appPath)
}
