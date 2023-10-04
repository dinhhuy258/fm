package config

import (
	gopher_lua "github.com/yuin/gopher-lua"

	"github.com/dinhhuy258/fm/pkg/config/lua"
)

// MessageConfig represents the config for the message.
type MessageConfig struct {
	Name string   `mapper:"name"`
	Args []string `mapper:"args"`
}

// toLuaTable convert to LuaTable object
func (mc *MessageConfig) toLuaTable(luaState *gopher_lua.LState) *gopher_lua.LTable {
	tbl := luaState.NewTable()
	tbl.RawSetString("name", gopher_lua.LString(mc.Name))

	argsTbl := luaState.NewTable()
	for _, arg := range mc.Args {
		argsTbl.Append(gopher_lua.LString(arg))
	}

	tbl.RawSetString("args", argsTbl)

	return tbl
}

// ActionConfig represents the config for the action.
type ActionConfig struct {
	Help     string           `mapper:"help"`
	Messages []*MessageConfig `mapper:"messages"`
}

// toLuaTable convert to LuaTable object
func (ac *ActionConfig) toLuaTable(luaState *gopher_lua.LState) *gopher_lua.LTable {
	tbl := luaState.NewTable()
	tbl.RawSetString("help", gopher_lua.LString(ac.Help))

	msgTbl := luaState.NewTable()
	for _, msg := range ac.Messages {
		msgTbl.Append(msg.toLuaTable(luaState))
	}

	tbl.RawSetString("messages", msgTbl)

	return tbl
}

// KeyBindingsConfig represents the config for the key bindings.
type KeyBindingsConfig struct {
	OnKeys   map[string]*ActionConfig `mapper:"on_keys"`
	OnNumber *ActionConfig            `mapper:"on_number"`
	Default  *ActionConfig            `mapper:"default"`
}

// toLuaTable convert to LuaTable object
func (kbc *KeyBindingsConfig) toLuaTable(luaState *gopher_lua.LState) *gopher_lua.LTable {
	tbl := luaState.NewTable()

	onKeyTbl := luaState.NewTable()
	for key, actionConfig := range kbc.OnKeys {
		onKeyTbl.RawSetString(key, actionConfig.toLuaTable(luaState))
	}

	tbl.RawSetString("on_keys", onKeyTbl)

	if kbc.OnNumber != nil {
		tbl.RawSetString("on_number", kbc.OnNumber.toLuaTable(luaState))
	} else {
		tbl.RawSetString("on_number", gopher_lua.LNil)
	}

	if kbc.Default != nil {
		tbl.RawSetString("default", kbc.Default.toLuaTable(luaState))
	} else {
		tbl.RawSetString("default", gopher_lua.LNil)
	}

	return tbl
}

// ModeConfig represents the config for the mode.
type ModeConfig struct {
	Name        string            `mapper:"name"`
	KeyBindings KeyBindingsConfig `mapper:"key_bindings"`
}

// toLuaTable convert to LuaTable object
func (mc *ModeConfig) toLuaTable(luaState *gopher_lua.LState) *gopher_lua.LTable {
	tbl := luaState.NewTable()

	tbl.RawSetString("name", gopher_lua.LString(mc.Name))
	tbl.RawSetString("key_bindings", mc.KeyBindings.toLuaTable(luaState))

	return tbl
}

// StyleConfig represents the config for style
type StyleConfig struct {
	Fg          string   `mapper:"fg"`
	Bg          string   `mapper:"bg"`
	Decorations []string `mapper:"decorations"`
}

// toLuaTable convert to LuaTable object
func (sc *StyleConfig) toLuaTable(luaState *gopher_lua.LState) *gopher_lua.LTable {
	tbl := luaState.NewTable()

	tbl.RawSetString("fg", gopher_lua.LString(sc.Fg))
	tbl.RawSetString("bg", gopher_lua.LString(sc.Bg))

	decorationTbl := luaState.NewTable()
	for _, decoration := range sc.Decorations {
		decorationTbl.Append(gopher_lua.LString(decoration))
	}

	tbl.RawSetString("decorations", decorationTbl)

	return tbl
}

// NodeTypeConfig represents the config for the node type (file/directory).
type NodeTypeConfig struct {
	Icon  string       `mapper:"icon"`
	Style *StyleConfig `mapper:"style"`
}

// toLuaTable convert to LuaTable object
func (ntc *NodeTypeConfig) toLuaTable(luaState *gopher_lua.LState) *gopher_lua.LTable {
	tbl := luaState.NewTable()

	tbl.RawSetString("icon", gopher_lua.LString(ntc.Icon))

	if ntc.Style != nil {
		tbl.RawSetString("style", ntc.Style.toLuaTable(luaState))
	} else {
		tbl.RawSetString("style", gopher_lua.LNil)
	}

	return tbl
}

// NodeTypesConfig represents the config for node types
type NodeTypesConfig struct {
	File             *NodeTypeConfig            `mapper:"file"`
	Directory        *NodeTypeConfig            `mapper:"directory"`
	FileSymlink      *NodeTypeConfig            `mapper:"file_symlink"`
	DirectorySymlink *NodeTypeConfig            `mapper:"directory_symlink"`
	Extensions       map[string]*NodeTypeConfig `mapper:"extensions"`
	Specials         map[string]*NodeTypeConfig `mapper:"specials"`
}

// toLuaTable convert to LuaTable object
func (ntc *NodeTypesConfig) toLuaTable(luaState *gopher_lua.LState) *gopher_lua.LTable {
	tbl := luaState.NewTable()

	if ntc.File != nil {
		tbl.RawSetString("file", ntc.File.toLuaTable(luaState))
	} else {
		tbl.RawSetString("file", gopher_lua.LNil)
	}

	if ntc.Directory != nil {
		tbl.RawSetString("directory", ntc.Directory.toLuaTable(luaState))
	} else {
		tbl.RawSetString("directory", gopher_lua.LNil)
	}

	if ntc.FileSymlink != nil {
		tbl.RawSetString("file_symlink", ntc.FileSymlink.toLuaTable(luaState))
	} else {
		tbl.RawSetString("file_symlink", gopher_lua.LNil)
	}

	if ntc.DirectorySymlink != nil {
		tbl.RawSetString("directory_symlink", ntc.DirectorySymlink.toLuaTable(luaState))
	} else {
		tbl.RawSetString("directory_symlink", gopher_lua.LNil)
	}

	extensionTbl := luaState.NewTable()
	for ext, extConfig := range ntc.Extensions {
		extensionTbl.RawSetString(ext, extConfig.toLuaTable(luaState))
	}

	tbl.RawSetString("extensions", extensionTbl)

	specialsTbl := luaState.NewTable()
	for fileName, specialConfig := range ntc.Specials {
		specialsTbl.RawSetString(fileName, specialConfig.toLuaTable(luaState))
	}

	tbl.RawSetString("specials", specialsTbl)

	return tbl
}

type UIConfig struct {
	Prefix string       `mapper:"prefix"`
	Suffix string       `mapper:"suffix"`
	Style  *StyleConfig `mapper:"style"`
}

// toLuaTable convert to LuaTable object
func (ui *UIConfig) toLuaTable(luaState *gopher_lua.LState) *gopher_lua.LTable {
	tbl := luaState.NewTable()

	tbl.RawSetString("prefix", gopher_lua.LString(ui.Prefix))
	tbl.RawSetString("suffix", gopher_lua.LString(ui.Suffix))

	if ui.Style != nil {
		tbl.RawSetString("style", ui.Style.toLuaTable(luaState))
	} else {
		tbl.RawSetString("style", gopher_lua.LNil)
	}

	return tbl
}

// DefaultUIConfig represents the config for UI
type DefaultUIConfig struct {
	Prefix         string       `mapper:"prefix"`
	Suffix         string       `mapper:"suffix"`
	FileStyle      *StyleConfig `mapper:"file_style"`
	DirectoryStyle *StyleConfig `mapper:"directory_style"`
}

// toLuaTable convert to LuaTable object
func (ui *DefaultUIConfig) toLuaTable(luaState *gopher_lua.LState) *gopher_lua.LTable {
	tbl := luaState.NewTable()

	tbl.RawSetString("prefix", gopher_lua.LString(ui.Prefix))
	tbl.RawSetString("suffix", gopher_lua.LString(ui.Suffix))

	if ui.FileStyle != nil {
		tbl.RawSetString("file_style", ui.FileStyle.toLuaTable(luaState))
	} else {
		tbl.RawSetString("file_style", gopher_lua.LNil)
	}

	if ui.DirectoryStyle != nil {
		tbl.RawSetString("directory_style", ui.DirectoryStyle.toLuaTable(luaState))
	} else {
		tbl.RawSetString("directory_style", gopher_lua.LNil)
	}

	return tbl
}

// ExplorerTableHeaderConfig represents the config for the explorer table header.
type ExplorerTableHeaderConfig struct {
	Name       string       `mapper:"name"`
	Percentage int          `mapper:"percentage"`
	Style      *StyleConfig `mapper:"style"`
}

// toLuaTable convert to LuaTable object
func (ethc *ExplorerTableHeaderConfig) toLuaTable(luaState *gopher_lua.LState) *gopher_lua.LTable {
	tbl := luaState.NewTable()

	tbl.RawSetString("name", gopher_lua.LString(ethc.Name))
	tbl.RawSetString("percentage", gopher_lua.LNumber(ethc.Percentage))

	if ethc.Style != nil {
		tbl.RawSetString("style", ethc.Style.toLuaTable(luaState))
	} else {
		tbl.RawSetString("style", gopher_lua.LNil)
	}

	return tbl
}

// ExplorerTableConfig represents the config for the explorer table.
type ExplorerTableConfig struct {
	IndexHeader       *ExplorerTableHeaderConfig `mapper:"index_header"`
	NameHeader        *ExplorerTableHeaderConfig `mapper:"name_header"`
	PermissionsHeader *ExplorerTableHeaderConfig `mapper:"permissions_header"`
	SizeHeader        *ExplorerTableHeaderConfig `mapper:"size_header"`

	DefaultUI        *DefaultUIConfig `mapper:"default_ui"`
	FocusUI          *UIConfig        `mapper:"focus_ui"`
	SelectionUI      *UIConfig        `mapper:"selection_ui"`
	FocusSelectionUI *UIConfig        `mapper:"focus_selection_ui"`

	FirstEntryPrefix string `mapper:"first_entry_prefix"`
	EntryPrefix      string `mapper:"entry_prefix"`
	LastEntryPrefix  string `mapper:"last_entry_prefix"`
}

// toLuaTable convert to LuaTable object
func (etc *ExplorerTableConfig) toLuaTable(luaState *gopher_lua.LState) *gopher_lua.LTable {
	tbl := luaState.NewTable()

	if etc.IndexHeader != nil {
		tbl.RawSetString("index_header", etc.IndexHeader.toLuaTable(luaState))
	} else {
		tbl.RawSetString("index_header", gopher_lua.LNil)
	}

	if etc.NameHeader != nil {
		tbl.RawSetString("name_header", etc.NameHeader.toLuaTable(luaState))
	} else {
		tbl.RawSetString("name_header", gopher_lua.LNil)
	}

	if etc.PermissionsHeader != nil {
		tbl.RawSetString("permissions_header", etc.PermissionsHeader.toLuaTable(luaState))
	} else {
		tbl.RawSetString("permissions_header", gopher_lua.LNil)
	}

	if etc.SizeHeader != nil {
		tbl.RawSetString("size_header", etc.SizeHeader.toLuaTable(luaState))
	} else {
		tbl.RawSetString("size_header", gopher_lua.LNil)
	}

	if etc.DefaultUI != nil {
		tbl.RawSetString("default_ui", etc.DefaultUI.toLuaTable(luaState))
	} else {
		tbl.RawSetString("default_ui", gopher_lua.LNil)
	}

	if etc.FocusUI != nil {
		tbl.RawSetString("focus_ui", etc.FocusUI.toLuaTable(luaState))
	} else {
		tbl.RawSetString("focus_ui", gopher_lua.LNil)
	}

	if etc.SelectionUI != nil {
		tbl.RawSetString("selection_ui", etc.SelectionUI.toLuaTable(luaState))
	} else {
		tbl.RawSetString("selection_ui", gopher_lua.LNil)
	}

	if etc.FocusSelectionUI != nil {
		tbl.RawSetString("focus_selection_ui", etc.FocusSelectionUI.toLuaTable(luaState))
	} else {
		tbl.RawSetString("focus_selection_ui", gopher_lua.LNil)
	}

	tbl.RawSetString("first_entry_prefix", gopher_lua.LString(etc.FirstEntryPrefix))
	tbl.RawSetString("entry_prefix", gopher_lua.LString(etc.EntryPrefix))
	tbl.RawSetString("last_entry_prefix", gopher_lua.LString(etc.LastEntryPrefix))

	return tbl
}

// SortingConfig represents the config for sorting
type SortingConfig struct {
	SortType         string `mapper:"sort_type"`
	Reverse          *bool  `mapper:"reverse"`
	IgnoreCase       *bool  `mapper:"ignore_case"`
	IgnoreDiacritics *bool  `mapper:"ignore_diacritics"`
}

// toLuaTable convert to LuaTable object
func (sc *SortingConfig) toLuaTable(luaState *gopher_lua.LState) *gopher_lua.LTable {
	tbl := luaState.NewTable()

	tbl.RawSetString("sort_type", gopher_lua.LString(sc.SortType))

	if sc.Reverse != nil {
		tbl.RawSetString("reverse", gopher_lua.LBool(*sc.Reverse))
	} else {
		tbl.RawSetString("reverse", gopher_lua.LNil)
	}

	if sc.IgnoreCase != nil {
		tbl.RawSetString("ignore_case", gopher_lua.LBool(*sc.IgnoreCase))
	} else {
		tbl.RawSetString("ignore_case", gopher_lua.LNil)
	}

	if sc.IgnoreDiacritics != nil {
		tbl.RawSetString("ignore_diacritics", gopher_lua.LBool(*sc.IgnoreDiacritics))
	} else {
		tbl.RawSetString("ignore_diacritics", gopher_lua.LNil)
	}

	return tbl
}

// FrameUI represents config for frame ui
type FrameUI struct {
	SelFrameColor string `mapper:"sel_frame_color"`
	FrameColor    string `mapper:"frame_color"`
}

// toLuaTable convert to LuaTable object
func (fu *FrameUI) toLuaTable(luaState *gopher_lua.LState) *gopher_lua.LTable {
	tbl := luaState.NewTable()

	tbl.RawSetString("sel_frame_color", gopher_lua.LString(fu.SelFrameColor))
	tbl.RawSetString("frame_color", gopher_lua.LString(fu.FrameColor))

	return tbl
}

// GeneralConfig represents the general config for the application.
type GeneralConfig struct {
	FrameUI *FrameUI `mapper:"frame_ui"`

	LogInfoUI    *UIConfig `mapper:"log_info_ui"`
	LogWarningUI *UIConfig `mapper:"log_warning_ui"`
	LogErrorUI   *UIConfig `mapper:"log_error_ui"`

	ExplorerTable *ExplorerTableConfig `mapper:"explorer_table"`

	Sorting    *SortingConfig `mapper:"sorting"`
	ShowHidden bool           `mapper:"show_hidden"`
}

// toLuaTable convert to LuaTable object
func (gc *GeneralConfig) toLuaTable(luaState *gopher_lua.LState) *gopher_lua.LTable {
	tbl := luaState.NewTable()

	if gc.FrameUI != nil {
		tbl.RawSetString("frame_ui", gc.FrameUI.toLuaTable(luaState))
	} else {
		tbl.RawSetString("frame_ui", gopher_lua.LNil)
	}

	if gc.LogInfoUI != nil {
		tbl.RawSetString("log_info_ui", gc.LogInfoUI.toLuaTable(luaState))
	} else {
		tbl.RawSetString("log_info_ui", gopher_lua.LNil)
	}

	if gc.LogWarningUI != nil {
		tbl.RawSetString("log_warning_ui", gc.LogWarningUI.toLuaTable(luaState))
	} else {
		tbl.RawSetString("log_warning_ui", gopher_lua.LNil)
	}

	if gc.LogErrorUI != nil {
		tbl.RawSetString("log_error_ui", gc.LogErrorUI.toLuaTable(luaState))
	} else {
		tbl.RawSetString("log_error_ui", gopher_lua.LNil)
	}

	if gc.ExplorerTable != nil {
		tbl.RawSetString("explorer_table", gc.ExplorerTable.toLuaTable(luaState))
	} else {
		tbl.RawSetString("explorer_table", gopher_lua.LNil)
	}

	if gc.Sorting != nil {
		tbl.RawSetString("sorting", gc.Sorting.toLuaTable(luaState))
	} else {
		tbl.RawSetString("sorting", gopher_lua.LNil)
	}

	tbl.RawSetString("show_hidden", gopher_lua.LBool(gc.ShowHidden))

	return tbl
}

// ModesConfig represents the config for the custom and builtin modes.
type ModesConfig struct {
	Customs  map[string]*ModeConfig `mapper:"customs"`
	Builtins map[string]*ModeConfig `mapper:"builtins"`
}

// toLuaTable convert to LuaTable object
func (m *ModesConfig) toLuaTable(luaState *gopher_lua.LState) *gopher_lua.LTable {
	tbl := luaState.NewTable()

	customTbl := luaState.NewTable()
	for name, modeConfig := range m.Customs {
		customTbl.RawSetString(name, modeConfig.toLuaTable(luaState))
	}

	tbl.RawSetString("customs", customTbl)

	builtinTbl := luaState.NewTable()
	for name, modeConfig := range m.Builtins {
		builtinTbl.RawSetString(name, modeConfig.toLuaTable(luaState))
	}

	tbl.RawSetString("builtins", builtinTbl)

	return tbl
}

// Config represents the config for the application.
type Config struct {
	General   *GeneralConfig   `mapper:"general"`
	Modes     *ModesConfig     `mapper:"modes"`
	NodeTypes *NodeTypesConfig `mapper:"node_types"`
}

// toLuaTable convert to LuaTable object
func (c *Config) toLuaTable(luaState *gopher_lua.LState) *gopher_lua.LTable {
	tbl := luaState.NewTable()

	if c.General != nil {
		tbl.RawSetString("general", c.General.toLuaTable(luaState))
	} else {
		tbl.RawSetString("general", gopher_lua.LNil)
	}

	if c.Modes != nil {
		tbl.RawSetString("modes", c.Modes.toLuaTable(luaState))
	} else {
		tbl.RawSetString("modes", gopher_lua.LNil)
	}

	if c.NodeTypes != nil {
		tbl.RawSetString("node_types", c.NodeTypes.toLuaTable(luaState))
	} else {
		tbl.RawSetString("node_types", gopher_lua.LNil)
	}

	return tbl
}

var AppConfig *Config

// LoadConfig loads the config from config file and default config then merges them.
func LoadConfig(lua *lua.Lua) error {
	configFilePath := getConfigFilePath()

	if configFilePath.IsPresent() {
		userConfig, err := loadConfigFromFile(*configFilePath.Get(), lua.GetState())
		if err != nil {
			return err
		}

		AppConfig = userConfig
	} else {
		AppConfig = getDefaultConfig()
	}

	return nil
}
