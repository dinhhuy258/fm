package config

var defaultModeConfig = ModeConfig{
	Name: "default",
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
			"j": {
				Help: "down",
				Commands: []*CommandConfig{
					{
						Name: "FocusNext",
					},
				},
			},
			"k": {
				Help: "up",
				Commands: []*CommandConfig{
					{
						Name: "FocusPrevious",
					},
				},
			},
			"l": {
				Help: "enter",
				Commands: []*CommandConfig{
					{
						Name: "Enter",
					},
				},
			},
			"h": {
				Help: "back",
				Commands: []*CommandConfig{
					{
						Name: "Back",
					},
				},
			},
			"m": {
				Help: "mark save",
				Commands: []*CommandConfig{
					{
						Name: "SwitchMode",
						Args: []interface{}{"mark-save"},
					},
				},
			},
			"`": {
				Help: "mark load",
				Commands: []*CommandConfig{
					{
						Name: "SwitchMode",
						Args: []interface{}{"mark-load"},
					},
				},
			},
			"d": {
				Help: "delete",
				Commands: []*CommandConfig{
					{
						Name: "SwitchMode",
						Args: []interface{}{"delete"},
					},
				},
			},
			"p": {
				Help: "copy",
				Commands: []*CommandConfig{
					{
						Name: "PasteSelections",
						Args: []interface{}{"copy"},
					},
				},
			},
			"x": {
				Help: "cut",
				Commands: []*CommandConfig{
					{
						Name: "PasteSelections",
						Args: []interface{}{"cut"},
					},
				},
			},
			"ctrl+r": {
				Help: "refresh",
				Commands: []*CommandConfig{
					{
						Name: "Refresh",
					},
				},
			},
			"space": {
				Help: "toggle selection",
				Commands: []*CommandConfig{
					{
						Name: "ToggleSelection",
					},
				},
			},
			"ctrl+space": {
				Help: "clear selection",
				Commands: []*CommandConfig{
					{
						Name: "ClearSelection",
					},
				},
			},
			".": {
				Help: "toggle hidden",
				Commands: []*CommandConfig{
					{
						Name: "ToggleHidden",
					},
				},
			},
		},
	},
}

var markSaveModeConfig = ModeConfig{
	Name: "mark-save",
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
			"esc": {
				Help: "cancel",
				Commands: []*CommandConfig{
					{
						Name: "FocusNext",
					},
				},
			},
		},
		Default: &ActionConfig{
			Commands: []*CommandConfig{
				{
					Name: "MarkSave",
				},
				{
					Name: "PopMode",
				},
			},
		},
	},
}

var markLoadModeConfig = ModeConfig{
	Name: "mark-load",
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
			"esc": {
				Help: "cancel",
				Commands: []*CommandConfig{
					{
						Name: "FocusNext",
					},
				},
			},
		},
		Default: &ActionConfig{
			Commands: []*CommandConfig{
				{
					Name: "MarkLoad",
				},
				{
					Name: "PopMode",
				},
			},
		},
	},
}

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

var searchModeConfig = ModeConfig{
	Name: "search",
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
				Help: "search",
				Commands: []*CommandConfig{
					{
						Name: "SearchFromInput",
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
	defaultModeConfig,
	markSaveModeConfig,
	markLoadModeConfig,
	newFileModeConfig,
	searchModeConfig,
	deleteModeConfig,
	deleteCurrentModeConfig,
	deleteSelectionsModeConfig,
}
