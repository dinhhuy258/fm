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

// GeneralConfig represents the general config for the application.
type GeneralConfig struct {
	SelectionUI      *UIConfig `yaml:"selectionUi"`
	FocusUI          *UIConfig `yaml:"focusUi"`
	DefaultUI        *UIConfig `yaml:"defaultUi"`
	FocusSelectionUI *UIConfig `yaml:"focusSelectionUi"`
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

	return &gc
}

// Config represents the config for the application.
type Config struct {
	General            *GeneralConfig         `yaml:"general"`
	PathPrefix         string                 `yaml:"pathPrefix"`
	PathSuffix         string                 `yaml:"pathSuffix"`
	LogErrorColor      string                 `yaml:"logErrorColor"`
	LogWarningColor    string                 `yaml:"logWarningColor"`
	LogInfoColor       string                 `yaml:"logInfoColor"`
	ShowHidden         bool                   `yaml:"showHidden"`
	IndexHeader        string                 `yaml:"indexHeader"`
	IndexPercentage    int                    `yaml:"indexPercentage"`
	PathHeader         string                 `yaml:"pathHeader"`
	PathPercentage     int                    `yaml:"pathPercentage"`
	FileModeHeader     string                 `yaml:"fileModeHeader"`
	FileModePercentage int                    `yaml:"fileModePercentage"`
	SizeHeader         string                 `yaml:"sizeHeader"`
	SizePercentage     int                    `yaml:"sizePercentage"`
	LogErrorFormat     string                 `yaml:"logErrorFormat"`
	LogWarningFormat   string                 `yaml:"logWarningFormat"`
	LogInfoFormat      string                 `yaml:"logInfoFormat"`
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

	if userConfig.LogErrorColor != "" {
		AppConfig.LogErrorColor = userConfig.LogErrorColor
	}

	if userConfig.LogWarningColor != "" {
		AppConfig.LogWarningColor = userConfig.LogWarningColor
	}

	if userConfig.LogInfoColor != "" {
		AppConfig.LogInfoColor = userConfig.LogInfoColor
	}

	if userConfig.IndexHeader != "" {
		AppConfig.IndexHeader = userConfig.IndexHeader
	}

	if userConfig.IndexPercentage != 0 {
		AppConfig.IndexPercentage = userConfig.IndexPercentage
	}

	if userConfig.PathHeader != "" {
		AppConfig.PathHeader = userConfig.PathHeader
	}

	if userConfig.PathPercentage != 0 {
		AppConfig.PathPercentage = userConfig.PathPercentage
	}

	if userConfig.FileModeHeader != "" {
		AppConfig.FileModeHeader = userConfig.FileModeHeader
	}

	if userConfig.FileModePercentage != 0 {
		AppConfig.FileModePercentage = userConfig.FileModePercentage
	}

	if userConfig.SizeHeader != "" {
		AppConfig.SizeHeader = userConfig.SizeHeader
	}

	if userConfig.SizePercentage != 0 {
		AppConfig.SizePercentage = userConfig.SizePercentage
	}

	if userConfig.PathPrefix != "" {
		AppConfig.PathPrefix = userConfig.PathPrefix
	}

	if userConfig.PathSuffix != "" {
		AppConfig.PathSuffix = userConfig.PathSuffix
	}

	if userConfig.LogErrorFormat != "" {
		AppConfig.LogErrorFormat = userConfig.LogErrorFormat
	}

	if userConfig.LogWarningFormat != "" {
		AppConfig.LogWarningFormat = userConfig.LogWarningFormat
	}

	if userConfig.LogInfoFormat != "" {
		AppConfig.LogInfoFormat = userConfig.LogInfoFormat
	}

	AppConfig.ShowHidden = userConfig.ShowHidden

	mergeUserModeConfig(userConfig.CustomModeConfigs, userConfig.BuiltinModeConfigs)
	mergeUserNodeTypesConfig(userConfig.NodeTypesConfig)
}
