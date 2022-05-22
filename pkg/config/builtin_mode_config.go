package config

var newFileModeConfig = ModeConfig{
	Name: "new-file",
	KeyBindings: KeyBindingsConfig{
		OnKeys: map[string]*ActionConfig{
			"ctrl+c": {
				Help: "quit",
				Commands: []*CommandConfig{
					{
						Name: "Quit",
					},
				},
			},
			"enter": {
				Help: "new file",
				Commands: []*CommandConfig{
					{
						Name: "NewFileFromInput",
					},
					{
						Name: "PopMode",
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
		Default: &ActionConfig{
			Commands: []*CommandConfig{
				{
					Name: "UpdateInputBufferFromKey",
				},
			},
		},
	},
}

var deleteModeConfig = ModeConfig{
	Name: "delete",
	KeyBindings: KeyBindingsConfig{
		OnKeys: map[string]*ActionConfig{
			"ctrl+c": {
				Help: "quit",
				Commands: []*CommandConfig{
					{
						Name: "Quit",
					},
				},
			},
			"d": {
				Help: "delete current",
				Commands: []*CommandConfig{
					{
						Name: "SetInputBuffer",
						Args: []interface{}{"Do you want to delete this file? (y/n) "},
					},
					{
						Name: "SwitchMode",
						Args: []interface{}{"delete-current"},
					},
				},
			},
			"s": {
				Help: "delete selections",
				Commands: []*CommandConfig{
					{
						Name: "SetInputBuffer",
						Args: []interface{}{"Do you want to delete selected files? (y/n) "},
					},
					{
						Name: "SwitchMode",
						Args: []interface{}{"delete-selections"},
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
}

var deleteCurrentModeConfig = ModeConfig{
	Name: "delete-current",
	KeyBindings: KeyBindingsConfig{
		OnKeys: map[string]*ActionConfig{
			"ctrl+c": {
				Help: "quit",
				Commands: []*CommandConfig{
					{
						Name: "Quit",
					},
				},
			},
			"y": {
				Help: "delete",
				Commands: []*CommandConfig{
					{
						Name: "DeleteCurrent",
					},
					{
						Name: "PopMode",
					},
					{
						Name: "PopMode",
					},
				},
			},
		},
		Default: &ActionConfig{
			Help: "cancel",
			Commands: []*CommandConfig{
				{
					Name: "PopMode",
				},
				{
					Name: "PopMode",
				},
			},
		},
	},
}

var deleteSelectionsModeConfig = ModeConfig{
	Name: "delete-selections",
	KeyBindings: KeyBindingsConfig{
		OnKeys: map[string]*ActionConfig{
			"ctrl+c": {
				Help: "quit",
				Commands: []*CommandConfig{
					{
						Name: "Quit",
					},
				},
			},
			"y": {
				Help: "delete selections",
				Commands: []*CommandConfig{
					{
						Name: "DeleteSelections",
					},
					{
						Name: "PopMode",
					},
					{
						Name: "PopMode",
					},
				},
			},
		},
		Default: &ActionConfig{
			Help: "cancel",
			Commands: []*CommandConfig{
				{
					Name: "PopMode",
				},
				{
					Name: "PopMode",
				},
			},
		},
	},
}

var builtinModeConfigs = []ModeConfig{
	newFileModeConfig,
	deleteModeConfig,
	deleteCurrentModeConfig,
	deleteSelectionsModeConfig,
}
