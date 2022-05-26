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
	Name        string            `yaml:"name"`
	KeyBindings KeyBindingsConfig `yaml:"keyBindings"`
}

type Config struct {
	SelectionColor     string       `yaml:"selectionColor"`
	DirectoryColor     string       `yaml:"directoryColor"`
	SizeStyle          string       `yaml:"sizeStyle"`
	LogErrorColor      string       `yaml:"logErrorColor"`
	LogWarningColor    string       `yaml:"logWarningColor"`
	LogInfoColor       string       `yaml:"logInfoColor"`
	FocusBg            string       `yaml:"focusBg"`
	FocusFg            string       `yaml:"focusFg"`
	ShowHidden         bool         `yaml:"showHidden"`
	IndexHeader        string       `yaml:"indexHeader"`
	IndexPercentage    int          `yaml:"indexPercentage"`
	PathHeader         string       `yaml:"pathHeader"`
	PathPercentage     int          `yaml:"pathPercentage"`
	SizeHeader         string       `yaml:"sizeHeader"`
	SizePercentage     int          `yaml:"sizePercentage"`
	PathPrefix         string       `yaml:"pathPrefix"`
	PathSuffix         string       `yaml:"pathSuffix"`
	FocusPrefix        string       `yaml:"focusPrefix"`
	FocusSuffix        string       `yaml:"focusSuffix"`
	SelectionPrefix    string       `yaml:"selectionPrefix"`
	SelectionSuffix    string       `yaml:"selectionSuffix"`
	FolderIcon         string       `yaml:"folderIcon"`
	FileIcon           string       `yaml:"fileIcon"`
	LogErrorFormat     string       `yaml:"logErrorFormat"`
	LogWarningFormat   string       `yaml:"logWarningFormat"`
	LogInfoFormat      string       `yaml:"logInfoFormat"`
	CustomModeConfigs  []ModeConfig `yaml:"customModeConfigs"`
	BuiltinModeConfigs []ModeConfig `yaml:"builtinModeConfigs"`
	DefaultModeConfig  ModeConfig   `yaml:"defaultModeConfig"`
}

var AppConfig *Config

func mergeConfig(config *Config) {
	if config.SelectionColor != "" {
		AppConfig.SelectionColor = config.SelectionColor
	}

	if config.DirectoryColor != "" {
		AppConfig.DirectoryColor = config.DirectoryColor
	}

	if config.SizeStyle != "" {
		AppConfig.SizeStyle = config.SizeStyle
	}

	if config.LogErrorColor != "" {
		AppConfig.LogErrorColor = config.LogErrorColor
	}

	if config.LogWarningColor != "" {
		AppConfig.LogWarningColor = config.LogWarningColor
	}

	if config.LogInfoColor != "" {
		AppConfig.LogInfoColor = config.LogInfoColor
	}

	if config.FocusBg != "" {
		AppConfig.FocusBg = config.FocusBg
	}

	if config.FocusFg != "" {
		AppConfig.FocusFg = config.FocusFg
	}

	if config.IndexHeader != "" {
		AppConfig.IndexHeader = config.IndexHeader
	}

	if config.IndexPercentage != 0 {
		AppConfig.IndexPercentage = config.IndexPercentage
	}

	if config.PathHeader != "" {
		AppConfig.PathHeader = config.PathHeader
	}

	if config.PathPercentage != 0 {
		AppConfig.PathPercentage = config.PathPercentage
	}

	if config.SizeHeader != "" {
		AppConfig.SizeHeader = config.SizeHeader
	}

	if config.SizePercentage != 0 {
		AppConfig.SizePercentage = config.SizePercentage
	}

	if config.PathPrefix != "" {
		AppConfig.PathPrefix = config.PathPrefix
	}

	if config.PathSuffix != "" {
		AppConfig.PathSuffix = config.PathSuffix
	}

	if config.FocusPrefix != "" {
		AppConfig.FocusPrefix = config.FocusPrefix
	}

	if config.FocusSuffix != "" {
		AppConfig.FocusSuffix = config.FocusSuffix
	}

	if config.SelectionPrefix != "" {
		AppConfig.SelectionPrefix = config.SelectionPrefix
	}

	if config.SelectionSuffix != "" {
		AppConfig.SelectionSuffix = config.SelectionSuffix
	}

	if config.FolderIcon != "" {
		AppConfig.FolderIcon = config.FolderIcon
	}

	if config.FileIcon != "" {
		AppConfig.FileIcon = config.FileIcon
	}

	if config.LogErrorFormat != "" {
		AppConfig.LogErrorFormat = config.LogErrorFormat
	}

	if config.LogWarningFormat != "" {
		AppConfig.LogWarningFormat = config.LogWarningFormat
	}

	if config.LogInfoFormat != "" {
		AppConfig.LogInfoFormat = config.LogInfoFormat
	}

	AppConfig.ShowHidden = config.ShowHidden
	AppConfig.CustomModeConfigs = config.CustomModeConfigs
}

func LoadConfig() error {
	configFilePath, err := getConfigFileOrCreateIfMissing()
	if err != nil {
		return err
	}

	AppConfig = getDefaultConfig()

	config, err := loadConfigFromFile(*configFilePath)
	if err != nil {
		return err
	}

	mergeConfig(config)

	return nil
}
