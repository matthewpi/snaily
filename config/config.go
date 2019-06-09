package config

import (
	"encoding/json"
	"github.com/matthewpi/snaily/logger"
	"os"
	"runtime"
)

// Config .
type Config struct {
	Build struct {
		Name      string `json:"name"`
		Version   string `json:"version"`
		Branch    string `json:"branch"`
		Commit    string `json:"commit"`
		Date      string `json:"date"`
		GoVersion string `json:"goVersion"`
	} `json:"-"`

	Backend struct {
		Redis struct {
			URI      string `json:"uri"`
			Password string `json:"password"`
			Database int    `json:"database"`
		} `json:"redis"`
	} `json:"backend"`

	API struct {
		Hypixel struct {
			Key string `json:"key"`
		} `json:"hypixel"`

		Steam struct {
			Key string `json:"key"`
		} `json:"steam"`

		Youtube struct {
			Key string `json:"key"`
		} `json:"youtube"`
	} `json:"api"`

	Discord struct {
		Token  string `json:"token"`
		Prefix string `json:"prefix"`

		Status struct {
			Active bool   `json:"active"`
			Name   string `json:"name"`
			Type   string `json:"type"`
		} `json:"status"`

		Guilds map[string]struct {
			Channels struct {
				Messages    string `json:"messages"`
				Punishments string `json:"punishments"`
			} `json:"channels"`

			Roles map[string]string `json:"roles"`
		} `json:"guilds"`
	}

	Filter struct {
		Active bool     `json:"active"`
		Words  []string `json:"words"`
	} `json:"filter"`
}

var config *Config

// Get returns the loaded config object.
func Get() *Config {
	return config
}

// Load loads the configuration from the disk.
func Load(name string, version string, branch string, commit string, date string) error {
	file, err := os.Open(".env/config.json")

	if err != nil {
		return err
	}

	parser := json.NewDecoder(file)
	err = parser.Decode(&config)
	if err != nil {
		return err
	}

	config.Build.Name = name
	config.Build.Version = version
	config.Build.Branch = branch
	config.Build.Commit = commit
	config.Build.Date = date
	config.Build.GoVersion = runtime.Version()

	err = file.Close()
	if err != nil {
		logger.Fatalw("[Preflight] Failed to close configuration file.", logger.Err(err))
	}

	return nil
}
