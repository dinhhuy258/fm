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

// UIConfig represents the config for UI
type UIConfig struct {
	Prefix         string       `yaml:"prefix"`
	Suffix         string       `yaml:"suffix"`
	FileStyle      *StyleConfig `yaml:"fileStyle"`
	DirectoryStyle *StyleConfig `yaml:"directoryStyle"`
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

	ui.FileStyle = ui.FileStyle.merge(other.FileStyle)
	ui.DirectoryStyle = ui.DirectoryStyle.merge(other.DirectoryStyle)

	return &ui
}

// LogUIConfig represents the config for logging UI.
type LogUIConfig struct {
	Prefix string       `yaml:"prefix"`
	Suffix string       `yaml:"suffix"`
	Style  *StyleConfig `yaml:"style"`
}

// merge user config with default config.
func (luc LogUIConfig) merge(other *LogUIConfig) *LogUIConfig {
	if other == nil {
		return &luc
	}

	if other.Prefix != "" {
		luc.Prefix = other.Prefix
	}

	if other.Suffix != "" {
		luc.Suffix = other.Suffix
	}

	luc.Style = luc.Style.merge(other.Style)

	return &luc
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

	ethc.Style = ethc.Style.merge(other.Style)

	return &ethc
}

// ExplorerTableConfig represents the config for the explorer table.
type ExplorerTableConfig struct {
	IndexHeader       *ExplorerTableHeaderConfig `yaml:"indexHeader"`
	NameHeader        *ExplorerTableHeaderConfig `yaml:"nameHeader"`
	PermissionsHeader *ExplorerTableHeaderConfig `yaml:"permissionsHeader"`
	SizeHeader        *ExplorerTableHeaderConfig `yaml:"sizeHeader"`
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

	return &etc
}

// GeneralConfig represents the general config for the application.
type GeneralConfig struct {
	SelectionUI      *UIConfig `yaml:"selectionUi"`
	FocusUI          *UIConfig `yaml:"focusUi"`
	DefaultUI        *UIConfig `yaml:"defaultUi"`
	FocusSelectionUI *UIConfig `yaml:"focusSelectionUi"`

	LogInfoUI    *LogUIConfig `yaml:"logInfoUi"`
	LogWarningUI *LogUIConfig `yaml:"logWarningUi"`
	LogErrorUI   *LogUIConfig `yaml:"logErrorUi"`

	ExplorerTable *ExplorerTableConfig `yaml:"explorerTable"`

	ShowHidden bool `yaml:"showHidden"`
}

// merge user config with default config.
func (gc GeneralConfig) merge(other *GeneralConfig) *GeneralConfig {
	if other == nil {
		return &gc
	}

	gc.DefaultUI.merge(other.DefaultUI)
	gc.FocusUI.merge(other.FocusUI)
	gc.SelectionUI.merge(other.SelectionUI)
	gc.FocusSelectionUI.merge(other.FocusSelectionUI)

	gc.LogInfoUI.merge(other.LogInfoUI)
	gc.LogWarningUI.merge(other.LogWarningUI)
	gc.LogErrorUI.merge(other.LogErrorUI)

	gc.ExplorerTable.merge(other.ExplorerTable)

	gc.ShowHidden = other.ShowHidden

	return &gc
}

// Config represents the config for the application.
type Config struct {
	General    *GeneralConfig `yaml:"general"`
	PathPrefix string         `yaml:"pathPrefix"`
	PathSuffix string         `yaml:"pathSuffix"`

	NodeTypesConfig    *NodeTypesConfig       `yaml:"nodeTypes"`
	CustomModeConfigs  map[string]*ModeConfig `yaml:"customModeConfigs"`
	BuiltinModeConfigs map[string]*ModeConfig `yaml:"builtinModeConfigs"`
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

		mergeUserConfig(userConfig)
	}

	return nil
}

// mergeUserNodeTypesConfig merges the user node types config.
func mergeUserNodeTypesConfig(userNodeTypesConfig *NodeTypesConfig) {
	if userNodeTypesConfig == nil {
		return
	}

	AppConfig.NodeTypesConfig.File = AppConfig.NodeTypesConfig.File.merge(userNodeTypesConfig.File)
	AppConfig.NodeTypesConfig.Directory = AppConfig.NodeTypesConfig.Directory.merge(userNodeTypesConfig.Directory)

	if userNodeTypesConfig.Extensions != nil {
		AppConfig.NodeTypesConfig.Extensions = map[string]*NodeTypeConfig{}

		for ext, ntc := range userNodeTypesConfig.Extensions {
			AppConfig.NodeTypesConfig.Extensions[ext] = AppConfig.NodeTypesConfig.File.merge(ntc)
		}
	}
}

// mergeUserModeConfig merges the user mode config.
func mergeUserModeConfig(customModeConfigs map[string]*ModeConfig, builtinModeConfigs map[string]*ModeConfig) {
	for name, mode := range customModeConfigs {
		mode.Name = name
	}

	AppConfig.CustomModeConfigs = customModeConfigs

	for builtinUserConfigName, builtinUserConfig := range builtinModeConfigs {
		builtinMode, hasBuiltinConfig := AppConfig.BuiltinModeConfigs[builtinUserConfigName]

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

// mergeUserConfig merges the user config.
func mergeUserConfig(userConfig *Config) {
	AppConfig.General = AppConfig.General.merge(userConfig.General)

	if userConfig.PathPrefix != "" {
		AppConfig.PathPrefix = userConfig.PathPrefix
	}

	if userConfig.PathSuffix != "" {
		AppConfig.PathSuffix = userConfig.PathSuffix
	}

	mergeUserModeConfig(userConfig.CustomModeConfigs, userConfig.BuiltinModeConfigs)
	mergeUserNodeTypesConfig(userConfig.NodeTypesConfig)
}
