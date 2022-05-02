package mode

import (
	"github.com/dinhhuy258/fm/pkg/app/command"
)

type DefaultMode struct {
	*Mode
}

func createDefaultMode() *DefaultMode {
	return &DefaultMode{
		&Mode{
			Name: "default",
			KeyBindings: &KeyBindings{
				OnKeys: map[string]*command.Command{
					"j": {
						Help: "down",
						Func: command.FocusNext,
					},
					"k": {
						Help: "up",
						Func: command.FocusPrevious,
					},
					"l": {
						Help: "enter",
						Func: command.Enter,
					},
					"h": {
						Help: "back",
						Func: command.Back,
					},
					"m": {
						Help: "mark save",
						Func: command.SwitchMode,
						Args: []interface{}{"mark-save"},
					},
					"`": {
						Help: "mark load",
						Func: command.SwitchMode,
						Args: []interface{}{"mark-load"},
					},
					"d": {
						Help: "delete",
						Func: command.SwitchMode,
						Args: []interface{}{"delete"},
					},
					"p": {
						Help: "copy",
						Func: command.PasteSelections,
						Args: []interface{}{"copy"},
					},
					"x": {
						Help: "cut",
						Func: command.PasteSelections,
						Args: []interface{}{"cut"},
					},
					"n": {
						Help: "new",
						Func: command.NewFile,
					},
					"ctrl+i": {
						Help: "next visited path",
						Func: command.NextVisitedPath,
					},
					"ctrl+o": {
						Help: "last visited path",
						Func: command.LastVisitedPath,
					},
					"ctrl+r": {
						Help: "refresh",
						Func: command.Refresh,
					},
					"space": {
						Help: "toggle selection",
						Func: command.ToggleSelection,
					},
					"ctrl+space": {
						Help: "clear selection",
						Func: command.ClearSelection,
					},
					".": {
						Help: "toggle hidden",
						Func: command.ToggleHidden,
					},
					"q": {
						Help: "quit",
						Func: command.Quit,
					},
				},
			},
		},
	}
}
