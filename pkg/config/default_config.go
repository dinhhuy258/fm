package config

func getDefaultConfig() *Config {
	return &Config{
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
									Args: []string{"/Users/dinhhuy258"},
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
									Args: []string{"/Users/dinhhuy258/Downloads"},
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
									Args: []string{"/Users/dinhhuy258/Documents"},
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
									Args: []string{"/Users/dinhhuy258/Workspace"},
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
									Args: []string{"/Users/dinhhuy258/Desktop"},
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
								Args: []string{"go-to"},
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
								Args: []string{"new-file"},
							},
							{
								Name: "SetInputBuffer",
								Args: []string{""},
							},
						},
					},
					"/": {
						Help: "search",
						Commands: []*CommandConfig{
							{
								Name: "SwitchMode",
								Args: []string{"search"},
							},
							{
								Name: "SetInputBuffer",
								Args: []string{""},
							},
						},
					},
				},
			},
		},
	}
}
