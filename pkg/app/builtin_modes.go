package app

func createDeleteMode() *Mode {
	return &Mode{
		name: "delete",
		keyBindings: &KeyBindings{
			onKeys: map[string]*Action{
				"d": {
					help: "delete",
					messages: []Message{
						{
							f: deleteCurrent,
						},
					},
				},
				"s": {
					help: "delete selections",
					messages: []Message{
						{
							f: deleteSelections,
						},
					},
				},
				"esc": {
					help: "cancel",
					messages: []Message{
						{
							f: popMode,
						},
					},
				},
				"q": {
					help: "quit",
					messages: []Message{
						{
							f: quit,
						},
					},
				},
			},
		},
	}
}

func createDefaultMode() *Mode {
	return &Mode{
		name: "default",
		keyBindings: &KeyBindings{
			onKeys: map[string]*Action{
				"j": {
					help: "down",
					messages: []Message{
						{
							f: focusNext,
						},
					},
				},
				"k": {
					help: "up",
					messages: []Message{
						{
							f: focusPrevious,
						},
					},
				},
				"l": {
					help: "enter",
					messages: []Message{
						{
							f: enter,
						},
					},
				},
				"h": {
					help: "back",
					messages: []Message{
						{
							f: back,
						},
					},
				},
				"d": {
					help: "delete",
					messages: []Message{
						{
							f:    switchMode,
							args: []interface{}{"delete"},
						},
					},
				},
				"ctrl+i": {
					help: "next visited path",
					messages: []Message{
						{
							f: nextVisitedPath,
						},
					},
				},
				"ctrl+o": {
					help: "last visited path",
					messages: []Message{
						{
							f: lastVisitedPath,
						},
					},
				},
				"ctrl+r": {
					help: "refresh",
					messages: []Message{
						{
							f: refresh,
						},
					},
				},
				"space": {
					help: "toggle selection",
					messages: []Message{
						{
							f: toggleSelection,
						},
					},
				},
				"ctrl+space": {
					help: "clear selection",
					messages: []Message{
						{
							f: clearSelection,
						},
					},
				},
				"q": {
					help: "quit",
					messages: []Message{
						{
							f: quit,
						},
					},
				},
			},
		},
	}
}
