package config

// getDefaultConfig returns the default configuration for the application.
func getDefaultConfig() *Config {
	return &Config{
		General: &GeneralConfig{
			LogInfoUI: &UIConfig{
				Prefix: "[Info] ",
				Suffix: "",
				Style: &StyleConfig{
					Fg: "green",
				},
			},
			LogWarningUI: &UIConfig{
				Prefix: "[Warning] ",
				Suffix: "",
				Style: &StyleConfig{
					Fg: "yellow",
				},
			},
			LogErrorUI: &UIConfig{
				Prefix: "[Error] ",
				Suffix: "",
				Style: &StyleConfig{
					Fg: "red",
				},
			},
			ExplorerTable: &ExplorerTableConfig{
				DefaultUI: &DefaultUIConfig{
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
					Style: &StyleConfig{
						Fg: "white",
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
				FirstEntryPrefix: "├─",
				EntryPrefix:      "├─",
				LastEntryPrefix:  "└─",
			},
			ShowHidden: false,
		},
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
		Modes: &ModesConfig{
			Builtins: builtinModeConfigs,
			Customs:  map[string]*ModeConfig{},
		},
	}
}
