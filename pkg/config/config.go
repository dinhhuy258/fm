package config

import (
	"github.com/dinhhuy258/fm/pkg/lua"
	gopher_lua "github.com/yuin/gopher-lua"
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
	OnKeys  map[string]*ActionConfig `mapper:"onKeys"`
	Default *ActionConfig            `mapper:"default"`
}

// toLuaTable convert to LuaTable object
func (kbc *KeyBindingsConfig) toLuaTable(luaState *gopher_lua.LState) *gopher_lua.LTable {
	tbl := luaState.NewTable()

	onKeyTbl := luaState.NewTable()
	for key, actionConfig := range kbc.OnKeys {
		onKeyTbl.RawSetString(key, actionConfig.toLuaTable(luaState))
	}
	tbl.RawSetString("onKeys", onKeyTbl)

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
	KeyBindings KeyBindingsConfig `mapper:"keyBindings"`
}

// toLuaTable convert to LuaTable object
func (mc *ModeConfig) toLuaTable(luaState *gopher_lua.LState) *gopher_lua.LTable {
	tbl := luaState.NewTable()

	tbl.RawSetString("name", gopher_lua.LString(mc.Name))
	tbl.RawSetString("keyBindings", mc.KeyBindings.toLuaTable(luaState))

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

// merge user config with default config.
func (sc StyleConfig) merge(other *StyleConfig) *StyleConfig {
	if other == nil {
		return &sc
	}

	if other.Fg != "" {
		sc.Fg = other.Fg
	}

	if other.Bg != "" {
		sc.Bg = other.Bg
	}

	if other.Decorations != nil {
		sc.Decorations = other.Decorations
	}

	return &sc
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

// merge user config with default config.
func (ntc NodeTypeConfig) merge(other *NodeTypeConfig) *NodeTypeConfig {
	if other == nil {
		return &ntc
	}

	if other.Icon != "" {
		ntc.Icon = other.Icon
	}

	ntc.Style = ntc.Style.merge(other.Style)

	return &ntc
}

// NodeTypesConfig represents the config for node types
type NodeTypesConfig struct {
	File             *NodeTypeConfig            `mapper:"file"`
	Directory        *NodeTypeConfig            `mapper:"directory"`
	FileSymlink      *NodeTypeConfig            `mapper:"fileSymlink"`
	DirectorySymlink *NodeTypeConfig            `mapper:"directorySymlink"`
	Extensions       map[string]*NodeTypeConfig `mapper:"extensions"`
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
		tbl.RawSetString("fileSymlink", ntc.FileSymlink.toLuaTable(luaState))
	} else {
		tbl.RawSetString("fileSymlink", gopher_lua.LNil)
	}

	if ntc.DirectorySymlink != nil {
		tbl.RawSetString("directorySymlink", ntc.DirectorySymlink.toLuaTable(luaState))
	} else {
		tbl.RawSetString("directorySymlink", gopher_lua.LNil)
	}

	extensionTbl := luaState.NewTable()
	for ext, extConfig := range ntc.Extensions {
		extensionTbl.RawSetString(ext, extConfig.toLuaTable(luaState))
	}
	tbl.RawSetString("extensions", extensionTbl)

	return tbl
}

func (ntc NodeTypesConfig) merge(other *NodeTypesConfig) *NodeTypesConfig {
	if other == nil {
		return &ntc
	}

	ntc.File = ntc.File.merge(other.File)
	ntc.Directory = ntc.Directory.merge(other.Directory)
	ntc.FileSymlink = ntc.FileSymlink.merge(other.FileSymlink)
	ntc.DirectorySymlink = ntc.DirectorySymlink.merge(other.DirectorySymlink)

	if other.Extensions != nil {
		ntc.Extensions = map[string]*NodeTypeConfig{}

		for ext, extConfig := range other.Extensions {
			ntc.Extensions[ext] = ntc.File.merge(extConfig)
		}
	}

	return &ntc
}

// UIConfig represents the config for UI
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

// merge user config with default config.
func (ui UIConfig) merge(other *UIConfig) *UIConfig {
	if other == nil {
		return &ui
	}

	if other.Prefix != "" {
		ui.Prefix = other.Prefix
	}

	if other.Suffix != "" {
		ui.Suffix = other.Suffix
	}

	ui.Style = ui.Style.merge(other.Style)

	return &ui
}

// DefaultUIConfig represents the config for UI
type DefaultUIConfig struct {
	Prefix         string       `mapper:"prefix"`
	Suffix         string       `mapper:"suffix"`
	FileStyle      *StyleConfig `mapper:"fileStyle"`
	DirectoryStyle *StyleConfig `mapper:"directoryStyle"`
}

// toLuaTable convert to LuaTable object
func (ui *DefaultUIConfig) toLuaTable(luaState *gopher_lua.LState) *gopher_lua.LTable {
	tbl := luaState.NewTable()

	tbl.RawSetString("prefix", gopher_lua.LString(ui.Prefix))
	tbl.RawSetString("suffix", gopher_lua.LString(ui.Suffix))

	if ui.FileStyle != nil {
		tbl.RawSetString("fileStyle", ui.FileStyle.toLuaTable(luaState))
	} else {
		tbl.RawSetString("fileStyle", gopher_lua.LNil)
	}

	if ui.DirectoryStyle != nil {
		tbl.RawSetString("directoryStyle", ui.DirectoryStyle.toLuaTable(luaState))
	} else {
		tbl.RawSetString("directoryStyle", gopher_lua.LNil)
	}

	return tbl
}

// merge user config with default config.
func (ui DefaultUIConfig) merge(other *DefaultUIConfig) *DefaultUIConfig {
	if other == nil {
		return &ui
	}

	if other.Prefix != "" {
		ui.Prefix = other.Prefix
	}

	if other.Suffix != "" {
		ui.Suffix = other.Suffix
	}

	ui.FileStyle = ui.FileStyle.merge(other.FileStyle)
	ui.DirectoryStyle = ui.DirectoryStyle.merge(other.DirectoryStyle)

	return &ui
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

// merge user config with default config.
func (ethc ExplorerTableHeaderConfig) merge(other *ExplorerTableHeaderConfig) *ExplorerTableHeaderConfig {
	if other == nil {
		return &ethc
	}

	if other.Name != "" {
		ethc.Name = other.Name
	}

	if other.Percentage != 0 {
		ethc.Percentage = other.Percentage
	}

	// By dèault the style is null
	ethc.Style = other.Style

	return &ethc
}

// ExplorerTableConfig represents the config for the explorer table.
type ExplorerTableConfig struct {
	IndexHeader       *ExplorerTableHeaderConfig `mapper:"indexHeader"`
	NameHeader        *ExplorerTableHeaderConfig `mapper:"nameHeader"`
	PermissionsHeader *ExplorerTableHeaderConfig `mapper:"permissionsHeader"`
	SizeHeader        *ExplorerTableHeaderConfig `mapper:"sizeHeader"`

	DefaultUI        *DefaultUIConfig `mapper:"defaultUi"`
	FocusUI          *UIConfig        `mapper:"focusUi"`
	SelectionUI      *UIConfig        `mapper:"selectionUi"`
	FocusSelectionUI *UIConfig        `mapper:"focusSelectionUi"`

	FirstEntryPrefix string `mapper:"firstEntryPrefix"`
	EntryPrefix      string `mapper:"entryPrefix"`
	LastEntryPrefix  string `mapper:"lastEntryPrefix"`
}

// toLuaTable convert to LuaTable object
func (etc *ExplorerTableConfig) toLuaTable(luaState *gopher_lua.LState) *gopher_lua.LTable {
	tbl := luaState.NewTable()

	if etc.IndexHeader != nil {
		tbl.RawSetString("indexHeader", etc.IndexHeader.toLuaTable(luaState))
	} else {
		tbl.RawSetString("indexHeader", gopher_lua.LNil)
	}

	if etc.NameHeader != nil {
		tbl.RawSetString("nameHeader", etc.NameHeader.toLuaTable(luaState))
	} else {
		tbl.RawSetString("nameHeader", gopher_lua.LNil)
	}

	if etc.PermissionsHeader != nil {
		tbl.RawSetString("permissionsHeader", etc.PermissionsHeader.toLuaTable(luaState))
	} else {
		tbl.RawSetString("permissionsHeader", gopher_lua.LNil)
	}

	if etc.SizeHeader != nil {
		tbl.RawSetString("sizeHeader", etc.SizeHeader.toLuaTable(luaState))
	} else {
		tbl.RawSetString("sizeHeader", gopher_lua.LNil)
	}

	if etc.DefaultUI != nil {
		tbl.RawSetString("defaultUi", etc.DefaultUI.toLuaTable(luaState))
	} else {
		tbl.RawSetString("defaultUi", gopher_lua.LNil)
	}

	if etc.FocusUI != nil {
		tbl.RawSetString("focusUi", etc.FocusUI.toLuaTable(luaState))
	} else {
		tbl.RawSetString("focusUi", gopher_lua.LNil)
	}

	if etc.SelectionUI != nil {
		tbl.RawSetString("selectionUi", etc.SelectionUI.toLuaTable(luaState))
	} else {
		tbl.RawSetString("selectionUi", gopher_lua.LNil)
	}

	if etc.FocusSelectionUI != nil {
		tbl.RawSetString("focusSelectionUi", etc.FocusSelectionUI.toLuaTable(luaState))
	} else {
		tbl.RawSetString("focusSelectionUi", gopher_lua.LNil)
	}

	tbl.RawSetString("firstEntryPrefix", gopher_lua.LString(etc.FirstEntryPrefix))
	tbl.RawSetString("entryPrefix", gopher_lua.LString(etc.EntryPrefix))
	tbl.RawSetString("lastEntryPrefix", gopher_lua.LString(etc.LastEntryPrefix))

	return tbl
}

// merge user config with default config.
func (etc ExplorerTableConfig) merge(other *ExplorerTableConfig) *ExplorerTableConfig {
	if other == nil {
		return &etc
	}

	etc.IndexHeader = etc.IndexHeader.merge(other.IndexHeader)
	etc.NameHeader = etc.NameHeader.merge(other.NameHeader)
	etc.PermissionsHeader = etc.PermissionsHeader.merge(other.PermissionsHeader)
	etc.SizeHeader = etc.SizeHeader.merge(other.SizeHeader)

	etc.DefaultUI = etc.DefaultUI.merge(other.DefaultUI)
	etc.FocusUI = etc.FocusUI.merge(other.FocusUI)
	etc.SelectionUI = etc.SelectionUI.merge(other.SelectionUI)
	etc.FocusSelectionUI = etc.FocusSelectionUI.merge(other.FocusSelectionUI)

	if other.FirstEntryPrefix != "" {
		etc.FirstEntryPrefix = other.FirstEntryPrefix
	}

	if other.EntryPrefix != "" {
		etc.EntryPrefix = other.EntryPrefix
	}

	if other.LastEntryPrefix != "" {
		etc.LastEntryPrefix = other.LastEntryPrefix
	}

	return &etc
}

// SortingConfig represents the config for sorting
type SortingConfig struct {
	SortType         string `mapper:"sortType"`
	Reverse          *bool  `mapper:"reverse"`
	IgnoreCase       *bool  `mapper:"ignoreCase"`
	IgnoreDiacritics *bool  `mapper:"ignoreDiacritics"`
}

// toLuaTable convert to LuaTable object
func (sc *SortingConfig) toLuaTable(luaState *gopher_lua.LState) *gopher_lua.LTable {
	tbl := luaState.NewTable()

	tbl.RawSetString("sortType", gopher_lua.LString(sc.SortType))

	if sc.Reverse != nil {
		tbl.RawSetString("reverse", gopher_lua.LBool(*sc.Reverse))
	} else {
		tbl.RawSetString("reverse", gopher_lua.LNil)
	}

	if sc.IgnoreCase != nil {
		tbl.RawSetString("ignoreCase", gopher_lua.LBool(*sc.IgnoreCase))
	} else {
		tbl.RawSetString("ignoreCase", gopher_lua.LNil)
	}

	if sc.IgnoreDiacritics != nil {
		tbl.RawSetString("ignoreDiacritics", gopher_lua.LBool(*sc.IgnoreDiacritics))
	} else {
		tbl.RawSetString("ignoreDiacritics", gopher_lua.LNil)
	}

	return tbl
}

// merge user config with default config.
func (sc SortingConfig) merge(other *SortingConfig) *SortingConfig {
	if other == nil {
		return &sc
	}

	if other.SortType != "" {
		sc.SortType = other.SortType
	}

	if other.Reverse != nil {
		sc.Reverse = other.Reverse
	}

	if other.IgnoreCase != nil {
		sc.IgnoreCase = other.IgnoreCase
	}

	if other.IgnoreDiacritics != nil {
		sc.IgnoreDiacritics = other.IgnoreDiacritics
	}

	return &sc
}

// GeneralConfig represents the general config for the application.
type GeneralConfig struct {
	LogInfoUI    *UIConfig `mapper:"logInfoUi"`
	LogWarningUI *UIConfig `mapper:"logWarningUi"`
	LogErrorUI   *UIConfig `mapper:"logErrorUi"`

	ExplorerTable *ExplorerTableConfig `mapper:"explorerTable"`

	Sorting    *SortingConfig `mapper:"sorting"`
	ShowHidden bool           `mapper:"showHidden"`
}

// toLuaTable convert to LuaTable object
func (gc *GeneralConfig) toLuaTable(luaState *gopher_lua.LState) *gopher_lua.LTable {
	tbl := luaState.NewTable()

	if gc.LogInfoUI != nil {
		tbl.RawSetString("logInfoUi", gc.LogInfoUI.toLuaTable(luaState))
	} else {
		tbl.RawSetString("logInfoUi", gopher_lua.LNil)
	}

	if gc.LogWarningUI != nil {
		tbl.RawSetString("logWarningUi", gc.LogWarningUI.toLuaTable(luaState))
	} else {
		tbl.RawSetString("logWarningUi", gopher_lua.LNil)
	}

	if gc.LogErrorUI != nil {
		tbl.RawSetString("logErrorUi", gc.LogErrorUI.toLuaTable(luaState))
	} else {
		tbl.RawSetString("logErrorUi", gopher_lua.LNil)
	}

	if gc.ExplorerTable != nil {
		tbl.RawSetString("explorerTable", gc.ExplorerTable.toLuaTable(luaState))
	} else {
		tbl.RawSetString("explorerTable", gopher_lua.LNil)
	}

	if gc.Sorting != nil {
		tbl.RawSetString("sorting", gc.Sorting.toLuaTable(luaState))
	} else {
		tbl.RawSetString("sorting", gopher_lua.LNil)
	}

	tbl.RawSetString("showHidden", gopher_lua.LBool(gc.ShowHidden))

	return tbl
}

// merge user config with default config.
func (gc GeneralConfig) merge(other *GeneralConfig) *GeneralConfig {
	if other == nil {
		return &gc
	}

	gc.ExplorerTable = gc.ExplorerTable.merge(other.ExplorerTable)

	gc.LogInfoUI = gc.LogInfoUI.merge(other.LogInfoUI)
	gc.LogWarningUI = gc.LogWarningUI.merge(other.LogWarningUI)
	gc.LogErrorUI = gc.LogErrorUI.merge(other.LogErrorUI)

	gc.Sorting = gc.Sorting.merge(other.Sorting)

	gc.ShowHidden = other.ShowHidden

	return &gc
}

// ModesConfig represents the config for the custom and builtin modes.
type ModesConfig struct {
	Customs  map[string]*ModeConfig `mapper:"customs"`
	Builtins map[string]*ModeConfig `mapper:"builtins"`
}

// toLuaTable convert to LuaTable object
func (mc *ModesConfig) toLuaTable(luaState *gopher_lua.LState) *gopher_lua.LTable {
	tbl := luaState.NewTable()

	customTbl := luaState.NewTable()
	for name, modeConfig := range mc.Customs {
		customTbl.RawSetString(name, modeConfig.toLuaTable(luaState))
	}
	tbl.RawSetString("customs", customTbl)

	builtinTbl := luaState.NewTable()
	for name, modeConfig := range mc.Builtins {
		builtinTbl.RawSetString(name, modeConfig.toLuaTable(luaState))
	}
	tbl.RawSetString("builtins", builtinTbl)

	return tbl
}

func (m ModesConfig) merge(other *ModesConfig) *ModesConfig {
	if other == nil {
		return &m
	}

	if other.Customs != nil {
		for name, mode := range other.Customs {
			mode.Name = name
		}

		m.Customs = other.Customs
	}

	if other.Builtins != nil {
		for builtinModeName, builtinUserConfig := range other.Builtins {
			builtinMode, hasBuiltinConfig := m.Builtins[builtinModeName]
			if !hasBuiltinConfig {
				continue
			}

			for key, action := range builtinUserConfig.KeyBindings.OnKeys {
				builtinMode.KeyBindings.OnKeys[key] = action
			}

			if builtinUserConfig.KeyBindings.Default != nil {
				builtinMode.KeyBindings.Default = builtinUserConfig.KeyBindings.Default
			}
		}
	}

	return &m
}

// Config represents the config for the application.
type Config struct {
	General         *GeneralConfig   `mapper:"general"`
	Modes           *ModesConfig     `mapper:"modes"`
	NodeTypesConfig *NodeTypesConfig `mapper:"nodeTypes"`
}

// ToLuaTable convert to LuaTable object
func (c *Config) ToLuaTable(luaState *gopher_lua.LState) *gopher_lua.LTable {
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

	if c.NodeTypesConfig != nil {
		tbl.RawSetString("nodeTypes", c.NodeTypesConfig.toLuaTable(luaState))
	} else {
		tbl.RawSetString("nodeTypes", gopher_lua.LNil)
	}

	return tbl
}

var AppConfig *Config

// LoadConfig loads the config from config file and default config then merges them.
func LoadConfig(lua *lua.Lua) error {
	configFilePath := getConfigFilePath()
	AppConfig = GetDefaultConfig()

	if configFilePath.IsPresent() {
		userConfig, err := loadConfigFromFile(*configFilePath.Get(), lua.GetState())
		if err != nil {
			return err
		}

		// Merge user config with default config.
		AppConfig.General = AppConfig.General.merge(userConfig.General)
		AppConfig.NodeTypesConfig = AppConfig.NodeTypesConfig.merge(userConfig.NodeTypesConfig)
		AppConfig.Modes = AppConfig.Modes.merge(userConfig.Modes)
	}

	return nil
}
