package app

import (
	"github.com/dinhhuy258/fm/pkg/app/command"
)

type DefaultMode struct {
	*Mode
}

func (*DefaultMode) GetName() string {
	return "default"
}

func createDefaultMode() *DefaultMode {
	return &DefaultMode{
		&Mode{
			KeyBindings: &KeyBindings{
				OnKeys: map[string]*Action{
					"j": {
						Help: "down",
						Commands: []*command.Command{
							{

								Func: command.FocusNext,
							},
						},
					},
					"k": {
						Help: "up",
						Commands: []*command.Command{
							{
								Func: command.FocusPrevious,
							},
						},
					},
					"l": {
						Help: "enter",
						Commands: []*command.Command{
							{
								Func: command.Enter,
							},
						},
					},
					"h": {
						Help: "back",
						Commands: []*command.Command{
							{
								Func: command.Back,
							},
						},
					},
					"m": {
						Help: "mark save",
						Commands: []*command.Command{
							{
								Func: command.SwitchMode,
								Args: []interface{}{"mark-save"},
							},
						},
					},
					"`": {
						Help: "mark load",
						Commands: []*command.Command{
							{
								Func: command.SwitchMode,
								Args: []interface{}{"mark-load"},
							},
						},
					},
					"d": {
						Help: "delete",
						Commands: []*command.Command{
							{
								Func: command.SwitchMode,
								Args: []interface{}{"delete"},
							},
						},
					},
					"p": {
						Help: "copy",
						Commands: []*command.Command{
							{
								Func: command.PasteSelections,
								Args: []interface{}{"copy"},
							},
						},
					},
					"x": {
						Help: "cut",
						Commands: []*command.Command{
							{
								Func: command.PasteSelections,
								Args: []interface{}{"cut"},
							},
						},
					},
					"n": {
						Help: "new",
						Commands: []*command.Command{
							{
								Func: command.NewFile,
							},
						},
					},
					"ctrl+i": {
						Help: "next visited path",
						Commands: []*command.Command{
							{
								Func: command.NextVisitedPath,
							},
						},
					},
					"ctrl+o": {
						Help: "last visited path",
						Commands: []*command.Command{
							{
								Func: command.LastVisitedPath,
							},
						},
					},
					"ctrl+r": {
						Help: "refresh",
						Commands: []*command.Command{
							{
								Func: command.Refresh,
							},
						},
					},
					"space": {
						Help: "toggle selection",
						Commands: []*command.Command{
							{
								Func: command.ToggleSelection,
							},
						},
					},
					"ctrl+space": {
						Help: "clear selection",
						Commands: []*command.Command{
							{
								Func: command.ClearSelection,
							},
						},
					},
					".": {
						Help: "toggle hidden",
						Commands: []*command.Command{
							{
								Func: command.ToggleHidden,
							},
						},
					},
					// TODO: Move to configuration
					"g": {
						Help: "go to",
						Commands: []*command.Command{
							{
								Func: command.SwitchMode,
								Args: []interface{}{"go-to"},
							},
						},
					},
					"q": {
						Help: "quit",
						Commands: []*command.Command{
							{
								Func: command.Quit,
							},
						},
					},
				},
			},
		},
	}
}
