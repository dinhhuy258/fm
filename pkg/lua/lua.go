package lua

import (
	"github.com/dinhhuy258/fm/pkg/config"
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
	defaultConfig := config.GetDefaultConfig()
	fmConfig := l.state.NewUserData()
	fmConfig.Value = defaultConfig

	l.state.SetGlobal("fm", fmConfig)
	err := l.state.DoFile(fileConfig)
	if err != nil {
		return err
	}

	return nil
}
