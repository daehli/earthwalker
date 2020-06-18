// Package config handles the config.toml file and the environment variables
package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"gitlab.com/glatteis/earthwalker/util"
)

// Config holds server-wide settings
type Config struct {
	ConfigPath           string
	StaticPath           string
	DBPath               string
	Port                 string
	TileServerURL        string
	NoLabelTileServerURL string
}

// Read a Config from environment variables and TOML file, and return it
func Read() (Config, error) {
	conf := Config{
		ConfigPath: getEnv("EARTHWALKER_CONFIG_PATH", "config.toml"),
		StaticPath: getEnv("EARTHWALKER_STATIC_PATH", util.AppPath()),
		DBPath:     getDBPath(),
		Port:       getEnv("EARTHWALKER_PORT", "8080"),
	}

	tomlData, err := ioutil.ReadFile(conf.ConfigPath)
	if err != nil {
		log.Printf("Error reading/no config file at '%s', using defaults.\n", conf.ConfigPath)
		conf.TileServerURL = "https://tiles.wmflabs.org/osm/{z}/{x}/{y}.png"
		conf.NoLabelTileServerURL = "https://tiles.wmflabs.org/osm-no-labels/{z}/{x}/{y}.png"
	}
	if err := toml.Unmarshal(tomlData, &conf); err != nil {
		return conf, fmt.Errorf("error parsing TOML config file: %v", err)
	}
	return conf, nil
}

func getEnv(key string, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}

func getDBPath() string {
	path := ""
	pathSuffix := getEnv(os.Getenv("EARTHWALKER_DB_PATH"), "/badger/")
	pathRel := os.Getenv("EARTHWALKER_DB_PATH_REL")
	if pathRel == "cwd" {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		path = cwd + pathSuffix
	} else if pathRel == "absolute" {
		path = pathSuffix
	} else { // default: relative to executable
		path = util.AppPath() + pathSuffix
	}
	return path
}
