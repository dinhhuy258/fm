package mode

import "github.com/dinhhuy258/fm/pkg/message"

func createDeleteMode() *Mode {
	return &Mode{
		Name: "delete",
		KeyBindings: &KeyBindings{
			OnKeys: map[string]*Action{
				"d": {
					Help: "delete",
					Messages: []message.Message{
						{
							Func: message.DeleteCurrent,
						},
					},
				},
				"s": {
					Help: "delete selections",
					Messages: []message.Message{
						{
							Func: message.DeleteSelections,
						},
					},
				},
				"esc": {
					Help: "cancel",
					Messages: []message.Message{
						{
							Func: message.PopMode,
						},
					},
				},
				"q": {
					Help: "quit",
					Messages: []message.Message{
						{
							Func: message.Quit,
						},
					},
				},
			},
		},
	}
}

func createDefaultMode() *Mode {
	return &Mode{
		Name: "default",
		KeyBindings: &KeyBindings{
			OnKeys: map[string]*Action{
				"j": {
					Help: "down",
					Messages: []message.Message{
						{
							Func: message.FocusNext,
						},
					},
				},
				"k": {
					Help: "up",
					Messages: []message.Message{
						{
							Func: message.FocusPrevious,
						},
					},
				},
				"l": {
					Help: "enter",
					Messages: []message.Message{
						{
							Func: message.Enter,
						},
					},
				},
				"h": {
					Help: "back",
					Messages: []message.Message{
						{
							Func: message.Back,
						},
					},
				},
				"d": {
					Help: "delete",
					Messages: []message.Message{
						{
							Func: message.SwitchMode,
							Args: []interface{}{"delete"},
						},
					},
				},
				"p": {
					Help: "copy",
					Messages: []message.Message{
						{
							Func: message.CopySelections,
						},
					},
				},
				"ctrl+i": {
					Help: "next visited path",
					Messages: []message.Message{
						{
							Func: message.NextVisitedPath,
						},
					},
				},
				"ctrl+o": {
					Help: "last visited path",
					Messages: []message.Message{
						{
							Func: message.LastVisitedPath,
						},
					},
				},
				"ctrl+r": {
					Help: "refresh",
					Messages: []message.Message{
						{
							Func: message.Refresh,
						},
					},
				},
				"space": {
					Help: "toggle selection",
					Messages: []message.Message{
						{
							Func: message.ToggleSelection,
						},
					},
				},
				"ctrl+space": {
					Help: "clear selection",
					Messages: []message.Message{
						{
							Func: message.ClearSelection,
						},
					},
				},
				".": {
					Help: "toggle hidden",
					Messages: []message.Message{
						{
							Func: message.ToggleHidden,
						},
					},
				},
				"q": {
					Help: "quit",
					Messages: []message.Message{
						{
							Func: message.Quit,
						},
					},
				},
			},
		},
	}
}
