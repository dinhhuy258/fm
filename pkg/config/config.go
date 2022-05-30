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

// NodeTypeConfig represents the config for the node type (file/directory).
type NodeTypeConfig struct {
	Color string `yaml:"color"`
	Icon  string `yaml:"icon"`
}

type NodeTypesConfig struct {
	File       *NodeTypeConfig            `yaml:"file"`
	Directory  *NodeTypeConfig            `yaml:"directory"`
	Extensions map[string]*NodeTypeConfig `yaml:"extensions"`
}

// Config represents the config for the application.
type Config struct {
	SelectionColor     string                 `yaml:"selectionColor"`
	LogErrorColor      string                 `yaml:"logErrorColor"`
	LogWarningColor    string                 `yaml:"logWarningColor"`
	LogInfoColor       string                 `yaml:"logInfoColor"`
	FocusBg            string                 `yaml:"focusBg"`
	FocusFg            string                 `yaml:"focusFg"`
	ShowHidden         bool                   `yaml:"showHidden"`
	IndexHeader        string                 `yaml:"indexHeader"`
	IndexPercentage    int                    `yaml:"indexPercentage"`
	PathHeader         string                 `yaml:"pathHeader"`
	PathPercentage     int                    `yaml:"pathPercentage"`
	FileModeHeader     string                 `yaml:"fileModeHeader"`
	FileModePercentage int                    `yaml:"fileModePercentage"`
	SizeHeader         string                 `yaml:"sizeHeader"`
	SizePercentage     int                    `yaml:"sizePercentage"`
	PathPrefix         string                 `yaml:"pathPrefix"`
	PathSuffix         string                 `yaml:"pathSuffix"`
	FocusPrefix        string                 `yaml:"focusPrefix"`
	FocusSuffix        string                 `yaml:"focusSuffix"`
	SelectionPrefix    string                 `yaml:"selectionPrefix"`
	SelectionSuffix    string                 `yaml:"selectionSuffix"`
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

	if userNodeTypesConfig.File != nil {
		if AppConfig.NodeTypesConfig.File.Color == "" {
			AppConfig.NodeTypesConfig.File.Color = userNodeTypesConfig.File.Color
		}

		if AppConfig.NodeTypesConfig.File.Icon == "" {
			AppConfig.NodeTypesConfig.File.Icon = userNodeTypesConfig.File.Icon
		}
	}

	if userNodeTypesConfig.Directory != nil {
		if AppConfig.NodeTypesConfig.Directory.Color == "" {
			AppConfig.NodeTypesConfig.Directory.Color = userNodeTypesConfig.Directory.Color
		}

		if AppConfig.NodeTypesConfig.Directory.Icon == "" {
			AppConfig.NodeTypesConfig.Directory.Icon = userNodeTypesConfig.Directory.Icon
		}
	}

	// Currently, the default extension node type is not configurable.
	// We can assign it to the user config if it is set.
	if userNodeTypesConfig.Extensions != nil {
		AppConfig.NodeTypesConfig.Extensions = map[string]*NodeTypeConfig{}

		for ext, ntc := range userNodeTypesConfig.Extensions {
			if ntc.Color == "" {
				ntc.Color = AppConfig.NodeTypesConfig.File.Color
			}
			if ntc.Icon == "" {
				ntc.Icon = AppConfig.NodeTypesConfig.File.Icon
			}

			AppConfig.NodeTypesConfig.Extensions[ext] = ntc
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
	if userConfig.SelectionColor != "" {
		AppConfig.SelectionColor = userConfig.SelectionColor
	}

	if userConfig.LogErrorColor != "" {
		AppConfig.LogErrorColor = userConfig.LogErrorColor
	}

	if userConfig.LogWarningColor != "" {
		AppConfig.LogWarningColor = userConfig.LogWarningColor
	}

	if userConfig.LogInfoColor != "" {
		AppConfig.LogInfoColor = userConfig.LogInfoColor
	}

	if userConfig.FocusBg != "" {
		AppConfig.FocusBg = userConfig.FocusBg
	}

	if userConfig.FocusFg != "" {
		AppConfig.FocusFg = userConfig.FocusFg
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

	if userConfig.FocusPrefix != "" {
		AppConfig.FocusPrefix = userConfig.FocusPrefix
	}

	if userConfig.FocusSuffix != "" {
		AppConfig.FocusSuffix = userConfig.FocusSuffix
	}

	if userConfig.SelectionPrefix != "" {
		AppConfig.SelectionPrefix = userConfig.SelectionPrefix
	}

	if userConfig.SelectionSuffix != "" {
		AppConfig.SelectionSuffix = userConfig.SelectionSuffix
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
