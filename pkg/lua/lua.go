package lua

import (
	lua "github.com/yuin/gopher-lua"
)

// Lua represent object stores lua state
type Lua struct {
	state *lua.LState
}

// NewLua create new Lua object instance
func NewLua() *Lua {
	return &Lua{
		state: lua.NewState(),
	}
}

// Close lua object state
func (l *Lua) Close() {
	l.state.Close()
}

// LoadConfig load config from the given file
func (l *Lua) LoadConfig(fileConfig string) error {
	// defaultConfig := config.GetDefaultConfig()
	fmConfig := l.state.NewTable()
	generalConfig := l.state.NewTable()

	generalConfig.RawSetString("log_info_ui", l.state.NewTable())
	generalConfig.RawSetString("log_warning_ui", l.state.NewTable())
	generalConfig.RawSetString("log_error_ui", l.state.NewTable())
	generalConfig.RawSetString("explorer", l.state.NewTable())

	fmConfig.RawSetString("general", generalConfig)


	l.state.SetGlobal("fm", fmConfig)

	err := l.state.DoFile("/Users/dinhhuy258/Workspace/fm/init.lua")
	if err != nil {
		return err
	}

	luaConfig := l.state.GetGlobal("fm")
	print(luaConfig)
	// l.state.GetTable()

	return nil
}
