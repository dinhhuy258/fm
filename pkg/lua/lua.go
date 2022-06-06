package lua

import (
	"github.com/dinhhuy258/fm/pkg/config"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
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


	defaultConfig := config.GetDefaultConfig()
	tbl := luar.New(l.state, defaultConfig)


	l.state.SetGlobal("fm", tbl.(*lua.LUserData).Metatable)

	err := l.state.DoFile("/Users/dinhhuy258/Workspace/fm/init.lua")
	if err != nil {
		return err
	}

	// luaConfig := convertToTable(l.state.GetGlobal("fm"))
	// luaConfig.ForEach(func(l1, l2 lua.LValue) {
	// 	print(l1.String())
	// })
	// l.state.GetTable()

	print(tbl)


	return nil
}

// Convert LValue to LTable
func convertToTable(lValue lua.LValue) *lua.LTable {
	if lValue.Type() == lua.LTTable {
		return lValue.(*lua.LTable)
	}

	return nil
}

// Convert map[string]interface{} to LTable
func (l *Lua) convertToLTable(m map[string]interface{}) *lua.LTable {
	lTable := l.state.NewTable()

	for key, value := range m {
		switch value.(type) {
		case map[string]interface{}:
			lTable.RawSetString(key, l.convertToLTable(value.(map[string]interface{})))
		case int:
			lTable.RawSetString(key, lua.LNumber(value.(int)))
		case string:
			lTable.RawSetString(key, lua.LString(value.(string)))
		case bool:
			lTable.RawSetString(key, lua.LBool(value.(bool)))
		}
	}

	return lTable
}

// Convert LTable to map[string]interface{}
func convertToMap(lTable *lua.LTable) map[string]interface{} {
	result := make(map[string]interface{})

	lTable.ForEach(func(key, value lua.LValue) {
		switch value.Type() {
		case lua.LTTable:
			result[key.String()] = convertToMap(value.(*lua.LTable))
		default:
			result[key.String()] = value
		}
	})

	return result
}
