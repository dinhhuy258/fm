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

func LoadConfig() error {
	configFilePath, err := getConfigFileOrCreateIfMissing()
	if err != nil {
		return err
	}

	AppConfig, err = loadConfigFromFile(*configFilePath)

	return err
}
