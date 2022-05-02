package mode

import (
	message2 "github.com/dinhhuy258/fm/pkg/app/message"
)

func createDefaultMode() *Mode {
	return &Mode{
		Name: "default",
		KeyBindings: &KeyBindings{
			OnKeys: map[string]*Action{
				"j": {
					Help: "down",
					Messages: []message2.Message{
						{
							Func: message2.FocusNext,
						},
					},
				},
				"k": {
					Help: "up",
					Messages: []message2.Message{
						{
							Func: message2.FocusPrevious,
						},
					},
				},
				"l": {
					Help: "enter",
					Messages: []message2.Message{
						{
							Func: message2.Enter,
						},
					},
				},
				"h": {
					Help: "back",
					Messages: []message2.Message{
						{
							Func: message2.Back,
						},
					},
				},
				"m": {
					Help: "mark save",
					Messages: []message2.Message{
						{
							Func: message2.SwitchMode,
							Args: []interface{}{"mark-save"},
						},
					},
				},
				"`": {
					Help: "mark load",
					Messages: []message2.Message{
						{
							Func: message2.SwitchMode,
							Args: []interface{}{"mark-load"},
						},
					},
				},
				"d": {
					Help: "delete",
					Messages: []message2.Message{
						{
							Func: message2.SwitchMode,
							Args: []interface{}{"delete"},
						},
					},
				},
				"p": {
					Help: "copy",
					Messages: []message2.Message{
						{
							Func: message2.PasteSelections,
							Args: []interface{}{"copy"},
						},
					},
				},
				"x": {
					Help: "cut",
					Messages: []message2.Message{
						{
							Func: message2.PasteSelections,
							Args: []interface{}{"cut"},
						},
					},
				},
				"n": {
					Help: "new",
					Messages: []message2.Message{
						{
							Func: message2.NewFile,
						},
					},
				},
				"ctrl+i": {
					Help: "next visited path",
					Messages: []message2.Message{
						{
							Func: message2.NextVisitedPath,
						},
					},
				},
				"ctrl+o": {
					Help: "last visited path",
					Messages: []message2.Message{
						{
							Func: message2.LastVisitedPath,
						},
					},
				},
				"ctrl+r": {
					Help: "refresh",
					Messages: []message2.Message{
						{
							Func: message2.Refresh,
						},
					},
				},
				"space": {
					Help: "toggle selection",
					Messages: []message2.Message{
						{
							Func: message2.ToggleSelection,
						},
					},
				},
				"ctrl+space": {
					Help: "clear selection",
					Messages: []message2.Message{
						{
							Func: message2.ClearSelection,
						},
					},
				},
				".": {
					Help: "toggle hidden",
					Messages: []message2.Message{
						{
							Func: message2.ToggleHidden,
						},
					},
				},
				"q": {
					Help: "quit",
					Messages: []message2.Message{
						{
							Func: message2.Quit,
						},
					},
				},
			},
		},
	}
}
