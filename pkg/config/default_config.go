package config

// getDefaultConfig returns the default configuration for the application.
func getDefaultConfig() *Config {
	return &Config{
		General: &GeneralConfig{
			DefaultUI: &UIConfig{
				Prefix: "  ",
				Suffix: "",
			},
			FocusUI: &UIConfig{
				Prefix: "▸[",
				Suffix: "]",
				Style: &StyleConfig{
					Decorations: []string{
						"bold",
					},
				},
			},
			SelectionUI: &UIConfig{
				Prefix: " {",
				Suffix: "}",
				Style: &StyleConfig{
					Fg: "green",
				},
			},
			FocusSelectionUI: &UIConfig{
				Prefix: "▸[",
				Suffix: "]",
				Style: &StyleConfig{
					Fg: "green",
					Decorations: []string{
						"bold",
					},
				},
			},
		},
		ShowHidden:         false,
		IndexHeader:        "index",
		IndexPercentage:    10,
		PathHeader:         "┌──── path",
		PathPercentage:     65,
		FileModeHeader:     "permissions",
		FileModePercentage: 15,
		SizeHeader:         "size",
		SizePercentage:     10,
		PathPrefix:         "├─",
		PathSuffix:         "└─",
		LogErrorFormat:     "[ERROR] ",
		LogErrorColor:      "red",
		LogWarningFormat:   "[WARNING] ",
		LogWarningColor:    "yellow",
		LogInfoFormat:      "[INFO] ",
		LogInfoColor:       "green",
		NodeTypesConfig: &NodeTypesConfig{
			File: &NodeTypeConfig{
				Style: &StyleConfig{
					Fg: "white",
				},
				Icon: "",
			},
			Directory: &NodeTypeConfig{
				Style: &StyleConfig{
					Fg: "cyan",
				},
				Icon: "",
			},
			Extensions: map[string]*NodeTypeConfig{},
		},
		BuiltinModeConfigs: builtinModeConfigs,
		CustomModeConfigs:  map[string]*ModeConfig{},
	}
}
