package config

type CommandConfig struct {
	Name string   `yaml:"name"`
	Args []string `yaml:"args"`
}

type ActionConfig struct {
	Help     string           `yaml:"help"`
	Commands []*CommandConfig `yaml:"commands"`
}

type KeyBindingsConfig struct {
	OnKeys  map[string]*ActionConfig `yaml:"onKeys"`
	Default *ActionConfig            `yaml:"default"`
}

type ModeConfig struct {
	Name        string            `yaml:"-"`
	KeyBindings KeyBindingsConfig `yaml:"keyBindings"`
}

type Config struct {
	SelectionColor     string                 `yaml:"selectionColor"`
	DirectoryColor     string                 `yaml:"directoryColor"`
	SizeStyle          string                 `yaml:"sizeStyle"`
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
	SizeHeader         string                 `yaml:"sizeHeader"`
	SizePercentage     int                    `yaml:"sizePercentage"`
	PathPrefix         string                 `yaml:"pathPrefix"`
	PathSuffix         string                 `yaml:"pathSuffix"`
	FocusPrefix        string                 `yaml:"focusPrefix"`
	FocusSuffix        string                 `yaml:"focusSuffix"`
	SelectionPrefix    string                 `yaml:"selectionPrefix"`
	SelectionSuffix    string                 `yaml:"selectionSuffix"`
	FolderIcon         string                 `yaml:"folderIcon"`
	FileIcon           string                 `yaml:"fileIcon"`
	LogErrorFormat     string                 `yaml:"logErrorFormat"`
	LogWarningFormat   string                 `yaml:"logWarningFormat"`
	LogInfoFormat      string                 `yaml:"logInfoFormat"`
	CustomModeConfigs  map[string]*ModeConfig `yaml:"customModeConfigs"`
	BuiltinModeConfigs map[string]*ModeConfig `yaml:"builtinModeConfigs"`
}

var AppConfig *Config

func mergeUserModeConfig(userConfig *Config) {
	for name, mode := range userConfig.CustomModeConfigs {
		mode.Name = name
	}

	AppConfig.CustomModeConfigs = userConfig.CustomModeConfigs

	for builtinUserConfigName, builtinUserConfig := range userConfig.BuiltinModeConfigs {
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

func mergeUserConfig(userConfig *Config) {
	if userConfig.SelectionColor != "" {
		AppConfig.SelectionColor = userConfig.SelectionColor
	}

	if userConfig.DirectoryColor != "" {
		AppConfig.DirectoryColor = userConfig.DirectoryColor
	}

	if userConfig.SizeStyle != "" {
		AppConfig.SizeStyle = userConfig.SizeStyle
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

	if userConfig.FolderIcon != "" {
		AppConfig.FolderIcon = userConfig.FolderIcon
	}

	if userConfig.FileIcon != "" {
		AppConfig.FileIcon = userConfig.FileIcon
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

	mergeUserModeConfig(userConfig)
}

func LoadConfig() error {
	//TODO: Consider to remove code to create config file on missing
	configFilePath, err := getConfigFileOrCreateIfMissing()
	if err != nil {
		return err
	}

	AppConfig = getDefaultConfig()

	userConfig, err := loadConfigFromFile(*configFilePath)
	if err != nil {
		return err
	}

	mergeUserConfig(userConfig)

	return nil
}
