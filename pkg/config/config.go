package config

type CommandConfig struct {
	Name string
	Args []interface{}
}

type ActionConfig struct {
	Help     string
	Commands []*CommandConfig
}

type KeyBindingsConfig struct {
	OnKeys  map[string]*ActionConfig
	Default *ActionConfig
}

type ModeConfig struct {
	Name        string
	KeyBindings KeyBindingsConfig
}

type Config struct {
	SelectionColor     string
	DirectoryColor     string
	SizeStyle          string
	LogErrorColor      string
	LogWarningColor    string
	LogInfoColor       string
	FocusBg            string
	FocusFg            string
	ShowHidden         bool
	IndexHeader        string
	IndexPercentage    int
	PathHeader         string
	PathPercentage     int
	SizeHeader         string
	SizePercentage     int
	PathPrefix         string
	PathSuffix         string
	FocusPrefix        string
	FocusSuffix        string
	SelectionPrefix    string
	SelectionSuffix    string
	FolderIcon         string
	FileIcon           string
	LogErrorFormat     string
	LogWarningFormat   string
	LogInfoFormat      string
	CustomModeConfigs  []ModeConfig
	BuiltinModeConfigs []ModeConfig
	DefaultModeConfig  ModeConfig
}

var AppConfig *Config

func LoadConfig() {
	AppConfig = &Config{
		ShowHidden:         false,
		IndexHeader:        "index",
		IndexPercentage:    10,
		PathHeader:         "╭──── path",
		PathPercentage:     70,
		SizeHeader:         "size",
		SizePercentage:     20,
		PathPrefix:         "├─",
		PathSuffix:         "╰─",
		FocusPrefix:        "▸[",
		FocusSuffix:        "]",
		FocusBg:            "black",
		FocusFg:            "blue",
		SelectionPrefix:    "{",
		SelectionSuffix:    "}",
		SelectionColor:     "green",
		FolderIcon:         "",
		FileIcon:           "",
		DirectoryColor:     "cyan",
		SizeStyle:          "white",
		LogErrorFormat:     "[ERROR] ",
		LogErrorColor:      "red",
		LogWarningFormat:   "[WARNING] ",
		LogWarningColor:    "yellow",
		LogInfoFormat:      "[INFO] ",
		LogInfoColor:       "green",
		BuiltinModeConfigs: builtinModeConfigs,
		CustomModeConfigs: []ModeConfig{
			{
				Name: "go-to",
				KeyBindings: KeyBindingsConfig{
					OnKeys: map[string]*ActionConfig{
						"~": {
							Help: "Home",
							Commands: []*CommandConfig{
								{
									Name: "ChangeDirectory",
									Args: []interface{}{"/Users/dinhhuy258"},
								},
								{
									Name: "PopMode",
								},
							},
						},
						"d": {
							Help: "Downloads",
							Commands: []*CommandConfig{
								{
									Name: "ChangeDirectory",
									Args: []interface{}{"/Users/dinhhuy258/Downloads"},
								},
								{
									Name: "PopMode",
								},
							},
						},
						"D": {
							Help: "Documents",
							Commands: []*CommandConfig{
								{
									Name: "ChangeDirectory",
									Args: []interface{}{"/Users/dinhhuy258/Documents"},
								},
								{
									Name: "PopMode",
								},
							},
						},
						"w": {
							Help: "Workspace",
							Commands: []*CommandConfig{
								{
									Name: "ChangeDirectory",
									Args: []interface{}{"/Users/dinhhuy258/Workspace"},
								},
								{
									Name: "PopMode",
								},
							},
						},
						"h": {
							Help: "Desktop",
							Commands: []*CommandConfig{
								{
									Name: "ChangeDirectory",
									Args: []interface{}{"/Users/dinhhuy258/Desktop"},
								},
								{
									Name: "PopMode",
								},
							},
						},
						"g": {
							Help: "focus first",
							Commands: []*CommandConfig{
								{
									Name: "FocusFirst",
								},
								{
									Name: "PopMode",
								},
							},
						},
						"q": {
							Help: "quit",
							Commands: []*CommandConfig{
								{
									Name: "Quit",
								},
							},
						},
						"esc": {
							Help: "cancel",
							Commands: []*CommandConfig{
								{
									Name: "PopMode",
								},
							},
						},
					},
				},
			},
		},
		DefaultModeConfig: ModeConfig{
			Name: "default",
			KeyBindings: KeyBindingsConfig{
				OnKeys: map[string]*ActionConfig{
					"g": {
						Help: "go to",
						Commands: []*CommandConfig{
							{
								Name: "SwitchMode",
								Args: []interface{}{"go-to"},
							},
						},
					},
					"G": {
						Help: "focus last",
						Commands: []*CommandConfig{
							{
								Name: "FocusLast",
							},
						},
					},
					"n": {
						Help: "new file",
						Commands: []*CommandConfig{
							{
								Name: "SwitchMode",
								Args: []interface{}{"new-file"},
							},
							{
								Name: "SetInputBuffer",
								Args: []interface{}{""},
							},
						},
					},
					"/": {
						Help: "search",
						Commands: []*CommandConfig{
							{
								Name: "SwitchMode",
								Args: []interface{}{"search"},
							},
							{
								Name: "SetInputBuffer",
								Args: []interface{}{""},
							},
						},
					},
				},
			},
		},
	}
}
