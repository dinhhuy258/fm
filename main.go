package main

import (
	"fmt"

	"github.com/dinhhuy258/fm/pkg/config"
	lua "github.com/yuin/gopher-lua"
)

var version = "unversioned"

func main() {
	l := lua.NewState()
	defer l.Close()

	l.SetGlobal("version", lua.LString(version))
	version := l.GetGlobal("version")
	fmt.Println(version.String())
	l.DoString(`print("Hello World!")`)

	config.LoadConfig()
	// config.AppConfig
	Execute(config.AppConfig, "")
}

// method to execute lua config file and get the configuration object
func Execute(c *config.Config, file string) {
	l := lua.NewState()
	defer l.Close()
	// configUserData := lua.LUserData{Value: c}

	configUserData := l.NewUserData()
	configUserData.Value = c

	l.SetGlobal("config", configUserData)
	l.DoString(`print("Hello World!")`)

	// Read the config object from the lua state
	config := l.GetGlobal("config")
	golangObject := Convert(config.(*lua.LUserData))

	fmt.Println(golangObject.General.ExplorerTable.DefaultUI.DirectoryStyle.Fg)
	// if config.Type() != lua.LTUserData {
	// 	panic("config is not a userdata")
	// }
	// c = config.(*lua.LUserData).Value.(*config.Config)
}

// Convert lua.LUserData object to *config.Config
func Convert(l *lua.LUserData) *config.Config {
	return l.Value.(*config.Config)
}

