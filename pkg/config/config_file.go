package config

import (
	"os"
	"path/filepath"

	"github.com/yuin/gluamapper"
	gopher_lua "github.com/yuin/gopher-lua"

	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/types"
)

// AppDir is the name of the directory where the config file is stored.
const AppDir = "fm"

// ConfigFileName is the name of the config file that gets created.
const ConfigFileName = "config.lua"

// getConfigFilePath returns the user config file
func getConfigFilePath() types.Optional[string] {
	configDir := os.Getenv("XDG_CONFIG_HOME")

	if configDir == "" {
		homeDir := os.Getenv("HOME")

		if homeDir == "" {
			return types.NewEmptyOptional[string]()
		}

		configDir = homeDir + "/.config"
	}

	configFilePath := filepath.Join(configDir, AppDir, ConfigFileName)
	if fs.IsPathExists(configFilePath) {
		return types.NewOptional(configFilePath)
	}

	return types.NewEmptyOptional[string]()
}

// loadConfigFromFile loads the config file from the given path.
func loadConfigFromFile(path string, luaState *gopher_lua.LState) (*Config, error) {
	defaultConfigTbl := getDefaultConfig().toLuaTable(luaState)
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
