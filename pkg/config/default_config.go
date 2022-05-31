package config

// getDefaultConfig returns the default configuration for the application.
func getDefaultConfig() *Config {
	return &Config{
		General: &GeneralConfig{
			DefaultUI: &UIConfig{
				Prefix: "  ",
				Suffix: "",
				FileStyle: &StyleConfig{
					Fg: "white",
				},
				DirectoryStyle: &StyleConfig{
					Fg: "cyan",
				},
			},
			FocusUI: &UIConfig{
				Prefix: "▸[",
				Suffix: "]",
				FileStyle: &StyleConfig{
					Fg: "white",
					Decorations: []string{
						"bold",
					},
				},
				DirectoryStyle: &StyleConfig{
					Fg: "cyan",
					Decorations: []string{
						"bold",
					},
				},
			},
			SelectionUI: &UIConfig{
				Prefix: " {",
				Suffix: "}",
				FileStyle: &StyleConfig{
					Fg: "green",
				},
				DirectoryStyle: &StyleConfig{
					Fg: "green",
				},
			},
			FocusSelectionUI: &UIConfig{
				Prefix: "▸[",
				Suffix: "]",
				FileStyle: &StyleConfig{
					Fg: "green",
					Decorations: []string{
						"bold",
					},
				},
				DirectoryStyle: &StyleConfig{
					Fg: "green",
					Decorations: []string{
						"bold",
					},
				},
			},
			LogInfoUI: &LogUIConfig{
				Prefix: "[Info] ",
				Suffix: "",
				Style: &StyleConfig{
					Fg: "green",
				},
			},
			LogWarningUI: &LogUIConfig{
				Prefix: "[Warning] ",
				Suffix: "",
				Style: &StyleConfig{
					Fg: "yellow",
				},
			},
			LogErrorUI: &LogUIConfig{
				Prefix: "[Error] ",
				Suffix: "",
				Style: &StyleConfig{
					Fg: "red",
				},
			},
			ExplorerTable: &ExplorerTableConfig{
				IndexHeader: &ExplorerTableHeaderConfig{
					Name:       "index",
					Percentage: 10,
				},
				NameHeader: &ExplorerTableHeaderConfig{
					Name:       "┌──── name",
					Percentage: 65,
				},
				PermissionsHeader: &ExplorerTableHeaderConfig{
					Name:       "permissions",
					Percentage: 15,
				},
				SizeHeader: &ExplorerTableHeaderConfig{
					Name:       "size",
					Percentage: 10,
				},
			},
			ShowHidden: false,
		},
		PathPrefix: "├─",
		PathSuffix: "└─",
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
