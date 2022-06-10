package config

// GetDefaultConfig returns the default configuration for the application.
func GetDefaultConfig() *Config {
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
			Sorting: &SortingConfig{
				Reverse:          newBool(false),
				SortType:         "dirFirst",
				IgnoreCase:       newBool(true),
				IgnoreDiacritics: newBool(true),
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
			FileSymlink: &NodeTypeConfig{
				Style: &StyleConfig{
					Fg: "white",
				},
				Icon: "",
			},
			DirectorySymlink: &NodeTypeConfig{
				Style: &StyleConfig{
					Fg: "cyan",
				},
				Icon: "",
			},
			Extensions: extensionsNodeTypeConfig,
		},
		Modes: &ModesConfig{
			Builtins: builtinModeConfigs,
			Customs:  map[string]*ModeConfig{},
		},
	}
}

// newBool return the pointer of bool val
func newBool(val bool) *bool {
	return &val
}
