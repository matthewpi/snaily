package config

import (
	"encoding/json"
	"os"
	"runtime"
)

// Config .
type Config struct {
	Environment string `json:"environment"`

	Build struct {
		Name      string `json:"name"`
		Version   string `json:"version"`
		Branch    string `json:"branch"`
		Commit    string `json:"commit"`
		Date      string `json:"date"`
		GoVersion string `json:"goVersion"`
	} `json:"-"`

	Backend struct {
		MongoDB struct {
			URI      string `json:"uri"`
			Database string `json:"database"`
		} `json:"mongodb"`

		Redis struct {
			URI      string `json:"uri"`
			Password string `json:"password"`
			Database int    `json:"database"`
		} `json:"redis"`
	} `json:"backend"`

	Discord struct {
		Token  string `json:"token"`
		Prefix string `json:"prefix"`

		Channels struct {
			Punishments string `json:"punishments"`
			Messages    string `json:"messages"`
		} `json:"channels"`

		Roles struct {
			Enhanced string `json:"enhanced"`
			Boombox  string `json:"boombox"`
			Muted    string `json:"muted"`
		} `json:"roles"`
	} `json:"discord"`

	Steam struct {
		Key string `json:"key"`
	} `json:"steam"`
}

var config *Config

// Get returns the loaded config object.
func Get() *Config {
	return config
}

// IsProduction is self explanatory..
func IsProduction() bool {
	return config.Environment == "production"
}

// Load loads the configuration from the disk.
func Load(name string, version string, branch string, commit string, date string) error {
	file, err := os.Open(".env/config.json")
	defer file.Close()

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

	if len(config.Backend.MongoDB.URI) < 1 {
		config.Backend.MongoDB.URI = "mongodb://127.0.0.1:27017"
	}

	if len(config.Backend.MongoDB.Database) < 1 {
		config.Backend.MongoDB.Database = "stacktracefun"
	}

	return nil
}
