package config

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var errHomeDirNotFound = errors.New("neither $XDG_CONFIG_HOME nor $HOME are defined")

// AppDir is the name of the directory where the config file is stored.
const AppDir = "fm"

// ConfigFileName is the name of the config file that gets created.
const ConfigFileName = "config.yml"

// getDefaultConfigYamlContents returns the default config file contents as a string.
func getDefaultConfigYamlContents() string {
	defaultConfig := getDefaultConfig()

	yaml, _ := yaml.Marshal(defaultConfig)

	return string(yaml)
}

// writeDefaultConfigContents writes the default config file contents to the given file.
func writeDefaultConfigContents(newConfigFile *os.File) error {
	_, err := newConfigFile.WriteString(getDefaultConfigYamlContents())
	if err != nil {
		return err
	}

	return nil
}

// createConfigFileIfMissing creates the config file if it doesn't exist
func createConfigFileIfMissing(configFilePath string) error {
	if _, err := os.Stat(configFilePath); errors.Is(err, os.ErrNotExist) {
		newConfigFile, err := os.OpenFile(configFilePath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0o666)
		if err != nil {
			return err
		}
		defer newConfigFile.Close()

		return writeDefaultConfigContents(newConfigFile)
	}

	return nil
}

// getConfigFileOrCreateIfMissing returns the config file or creates it if it doesn't exist.
func getConfigFileOrCreateIfMissing() (*string, error) {
	var err error

	configDir := os.Getenv("XDG_CONFIG_HOME")

	if configDir == "" {
		homeDir := os.Getenv("HOME")

		if homeDir == "" {
			return nil, errHomeDirNotFound
		}

		configDir = homeDir + "/.config"
	}

	prsConfigDir := filepath.Join(configDir, AppDir)
	err = os.MkdirAll(prsConfigDir, os.ModePerm)

	if err != nil {
		return nil, err
	}

	configFilePath := filepath.Join(prsConfigDir, ConfigFileName)
	err = createConfigFileIfMissing(configFilePath)

	if err != nil {
		return nil, err
	}

	return &configFilePath, nil
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
