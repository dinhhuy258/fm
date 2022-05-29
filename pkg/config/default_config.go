package config

// getDefaultConfig returns the default configuration for the application.
func getDefaultConfig() *Config {
	return &Config{
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
		FocusPrefix:        "▸[",
		FocusSuffix:        "]",
		FocusBg:            "black",
		FocusFg:            "blue",
		SelectionPrefix:    " {",
		SelectionSuffix:    "}",
		SelectionColor:     "green",
		LogErrorFormat:     "[ERROR] ",
		LogErrorColor:      "red",
		LogWarningFormat:   "[WARNING] ",
		LogWarningColor:    "yellow",
		LogInfoFormat:      "[INFO] ",
		LogInfoColor:       "green",
		NodeTypesConfig: &NodeTypesConfig{
			File: &NodeTypeConfig{
				Color: "white",
				Icon:  "",
			},
			Directory: &NodeTypeConfig{
				Color: "cyan",
				Icon:  "",
			},
			Extensions: map[string]*NodeTypeConfig{},
		},
		BuiltinModeConfigs: builtinModeConfigs,
		CustomModeConfigs:  map[string]*ModeConfig{},
	}
}
