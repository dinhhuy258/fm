package config

// MessageConfig represents the config for the message.
type MessageConfig struct {
	Name string   `yaml:"name"`
	Args []string `yaml:"args"`
}

// ActionConfig represents the config for the action.
type ActionConfig struct {
	Help     string           `yaml:"help"`
	Messages []*MessageConfig `yaml:"messages"`
}

// KeyBindingsConfig represents the config for the key bindings.
type KeyBindingsConfig struct {
	OnKeys  map[string]*ActionConfig `yaml:"onKeys"`
	Default *ActionConfig            `yaml:"default"`
}

// ModeConfig represents the config for the mode.
type ModeConfig struct {
	Name        string            `yaml:"-"`
	KeyBindings KeyBindingsConfig `yaml:"keyBindings"`
}

// StyleConfig represents the config for style
type StyleConfig struct {
	Fg          string   `yaml:"fg"`
	Bg          string   `yaml:"bg"`
	Decorations []string `yaml:"decorations"`
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
	Icon  string       `yaml:"icon"`
	Style *StyleConfig `yaml:"style"`
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
	File       *NodeTypeConfig            `yaml:"file"`
	Directory  *NodeTypeConfig            `yaml:"directory"`
	Extensions map[string]*NodeTypeConfig `yaml:"extensions"`
}

func (ntc NodeTypesConfig) merge(other *NodeTypesConfig) *NodeTypesConfig {
	if other == nil {
		return &ntc
	}

	ntc.File = ntc.File.merge(other.File)
	ntc.Directory = ntc.Directory.merge(other.Directory)

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
	Prefix string       `yaml:"prefix"`
	Suffix string       `yaml:"suffix"`
	Style  *StyleConfig `yaml:"style"`
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
	Prefix         string       `yaml:"prefix"`
	Suffix         string       `yaml:"suffix"`
	FileStyle      *StyleConfig `yaml:"fileStyle"`
	DirectoryStyle *StyleConfig `yaml:"directoryStyle"`
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
	Name       string       `yaml:"name"`
	Percentage int          `yaml:"percentage"`
	Style      *StyleConfig `yaml:"style"`
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
	IndexHeader       *ExplorerTableHeaderConfig `yaml:"indexHeader"`
	NameHeader        *ExplorerTableHeaderConfig `yaml:"nameHeader"`
	PermissionsHeader *ExplorerTableHeaderConfig `yaml:"permissionsHeader"`
	SizeHeader        *ExplorerTableHeaderConfig `yaml:"sizeHeader"`

	DefaultUI        *DefaultUIConfig `yaml:"defaultUi"`
	FocusUI          *UIConfig        `yaml:"focusUi"`
	SelectionUI      *UIConfig        `yaml:"selectionUi"`
	FocusSelectionUI *UIConfig        `yaml:"focusSelectionUi"`

	FirstEntryPrefix string `yaml:"firstEntryPrefix"`
	EntryPrefix      string `yaml:"entryPrefix"`
	LastEntryPrefix  string `yaml:"lastEntryPrefix"`
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

// GeneralConfig represents the general config for the application.
type GeneralConfig struct {
	LogInfoUI    *UIConfig `yaml:"logInfoUi"`
	LogWarningUI *UIConfig `yaml:"logWarningUi"`
	LogErrorUI   *UIConfig `yaml:"logErrorUi"`

	ExplorerTable *ExplorerTableConfig `yaml:"explorerTable"`

	ShowHidden bool `yaml:"showHidden"`
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

	gc.ShowHidden = other.ShowHidden

	return &gc
}

// ModesConfig represents the config for the custom and builtin modes.
type ModesConfig struct {
	Customs  map[string]*ModeConfig `yaml:"customs"`
	Builtins map[string]*ModeConfig `yaml:"builtins"`
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
	General         *GeneralConfig   `yaml:"general"`
	Modes           *ModesConfig     `yaml:"modes"`
	NodeTypesConfig *NodeTypesConfig `yaml:"nodeTypes"`
}

var AppConfig *Config

// LoadConfig loads the config from config file and default config then merges them.
func LoadConfig() error {
	configFilePath := getConfigFilePath()
	AppConfig = getDefaultConfig()

	if configFilePath.IsPresent() {
		userConfig, err := loadConfigFromFile(*configFilePath.Get())
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
