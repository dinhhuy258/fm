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

// GetState returns the lua state object
func (l *Lua) GetState() *lua.LState {
	return l.state
}

// Close lua object state
func (l *Lua) Close() {
	l.state.Close()
}
