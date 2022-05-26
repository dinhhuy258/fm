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
		CustomModeConfigs:  []ModeConfig{},
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
