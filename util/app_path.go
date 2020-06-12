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

// StaticPath returns the path to parent directory of 'static' and 'templates'
// By default (no env variables) this is equivalent to AppPath()
// (moving these files is left to the user)
func StaticPath() string {
	staticPath := os.Getenv("EARTHWALKER_STATIC_PATH")
	if staticPath == "" {
		return AppPath()
	}
	return staticPath
}
