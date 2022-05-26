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
		CustomModeConfigs:  map[string]*ModeConfig{},
	}
}
