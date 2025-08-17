package config

// getDefaultConfig returns the default configuration for the application.
func getDefaultConfig() *Config {
	return &Config{
		General: &GeneralConfig{
			FrameUI: &FrameUI{
				SelFrameColor: "green",
				FrameColor:    "white",
			},
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
					Percentage: 15,
				},
				NameHeader: &ExplorerTableHeaderConfig{
					Name:       "┌──── name",
					Percentage: 85,
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
		NodeTypes: &NodeTypesConfig{
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
			Extensions: getExtensionsNodeTypeConfig(),
			Specials:   getSpecialsNodeTypeConfig(),
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
