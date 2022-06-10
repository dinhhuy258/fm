package config

import (
	"os"
	"path/filepath"

	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/optional"
	"github.com/yuin/gluamapper"
	gopher_lua "github.com/yuin/gopher-lua"
)

// AppDir is the name of the directory where the config file is stored.
const AppDir = "fm"

// ConfigFileName is the name of the config file that gets created.
const ConfigFileName = "config.lua"

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
func loadConfigFromFile(path string, luaState *gopher_lua.LState) (*Config, error) {
	defaultConfigTbl := AppConfig.ToLuaTable(luaState)
	luaState.SetGlobal("fm", defaultConfigTbl)

	if err := luaState.DoFile(path); err != nil {
		return nil, err
	}

	var config Config

	mapper := gluamapper.NewMapper(gluamapper.Option{
		NameFunc: func(s string) string {
			return s
		},
		TagName: "mapper",
	})

	fmConfig, _ := luaState.GetGlobal("fm").(*gopher_lua.LTable)
	if err := mapper.Map(fmConfig, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
