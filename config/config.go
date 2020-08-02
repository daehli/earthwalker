// Package config handles the config.toml file and the environment variables
package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/BurntSushi/toml"
	"gitlab.com/glatteis/earthwalker/domain"
)

// Read a Config from environment variables and TOML file, and return it
func Read() (domain.Config, error) {
	// defaults
	appPath := AppPath()
	conf := domain.Config{
		ConfigPath:           getEnv("EARTHWALKER_CONFIG_PATH", appPath+"/config.toml"),
		StaticPath:           appPath,
		DBPath:               appPath + "/badger",
		Port:                 "8080",
		TileServerURL:        "https://tiles.wmflabs.org/osm/{z}/{x}/{y}.png",
		NoLabelTileServerURL: "https://tiles.wmflabs.org/osm-no-labels/{z}/{x}/{y}.png",
	}

	// TOML
	tomlData, err := ioutil.ReadFile(conf.ConfigPath)
	if err != nil {
		log.Printf("Error reading/no config file at '%s', using defaults.\n", conf.ConfigPath)
	}
	if err := toml.Unmarshal(tomlData, &conf); err != nil {
		return conf, fmt.Errorf("error parsing TOML config file: %v", err)
	}

	// env vars
	conf.Port = getEnv("EARTHWALKER_PORT", conf.Port)
	conf.DBPath = getEnv("EARTHWALKER_DB_PATH", conf.DBPath)
	conf.StaticPath = getEnv("EARTHWALKER_STATIC_PATH", conf.StaticPath)

	return conf, nil
}

func getEnv(key string, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && len(v) > 0 {
		return v
	}
	return fallback
}

// AppPath gets the executable's path.
func AppPath() string {
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal("App path not accessible!")
	}
	return path.Dir(appPath)
}
