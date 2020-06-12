// Package config handles the config.toml file and the environment variables
package config

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"gitlab.com/glatteis/earthwalker/util"
)

// EnvType represents the environment file contents
type EnvType struct {
	// EarthwalkerStaticPath is documented in start.sh.sample
	EarthwalkerStaticPath string
	// EarthwalkerDBPath is documented in start.sh.sample
	EarthwalkerDBPath string
	// EarthwalkerPort is documented in start.sh.sample
	EarthwalkerPort string
	// EarthwalkerConfigPath is documented in start.sh.sample
	EarthwalkerConfigPath string
}

var Env EnvType

// FileType represents the config file contents
type FileType struct {
	// TileServerURL is the url of the tile server
	TileServerURL string
	// NoLabelTileServerURL is the url of the tile server when no labels are activated
	NoLabelTileServerURL string
}

var File FileType

func orDefault(envInput string, def string) string {
	if envInput == "" {
		return def
	}
	return envInput
}

func getDBPath() string {
	path := ""
	pathSuffix := orDefault(os.Getenv("EARTHWALKER_DB_PATH"), "/badger/")
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

func init() {
	// Initialize Env
	Env = EnvType{
		EarthwalkerStaticPath: orDefault(os.Getenv("EARTHWALKER_STATIC_PATH"), util.AppPath()),
		EarthwalkerDBPath:     getDBPath(),
		EarthwalkerPort:       orDefault(os.Getenv("EARTHWALKER_PORT"), ""), // Default is intentionally "" s.t. the command line port gets respected
		EarthwalkerConfigPath: orDefault(os.Getenv("EARTHWALKER_CONFIG_PATH"), "config.toml"),
	}

	// Initialize Config File
	tomlData, err := ioutil.ReadFile(Env.EarthwalkerConfigPath)
	if err != nil {
		log.Println("defaulting to default config.")
		File = FileType{
			TileServerURL:        "https://tiles.wmflabs.org/osm/{z}/{x}/{y}.png",
			NoLabelTileServerURL: "https://tiles.wmflabs.org/osm-no-labels/{z}/{x}/{y}.png",
		}
	}
	if err := toml.Unmarshal(tomlData, &File); err != nil {
		log.Fatal(err)
	}
}
