package config

import (
	"os"
	"path/filepath"

	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/optional"
	"gopkg.in/yaml.v3"
)

// AppDir is the name of the directory where the config file is stored.
const AppDir = "fm"

// ConfigFileName is the name of the config file that gets created.
const ConfigFileName = "config.yml"

// getConfigFilePath returns the user config file
func getConfigFilePath() optional.Optional[string] {
	configDir := os.Getenv("XDG_CONFIG_HOME")

	if configDir == "" {
		homeDir := os.Getenv("HOME")

		if homeDir == "" {
			return optional.NewEmpty[string]()
		}

		configDir = homeDir + "/.config"
	}

	configFilePath := filepath.Join(configDir, AppDir, ConfigFileName)
	if fs.IsFileExists(configFilePath) {
		return optional.New(configFilePath)
	}

	return optional.NewEmpty[string]()
}

// loadConfigFromFile loads the config file from the given path.
func loadConfigFromFile(path string) (*Config, error) {
	config := &Config{}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal((data), &config)

	return config, err
}
