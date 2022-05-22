package app

import (
	"github.com/dinhhuy258/fm/pkg/app/command"
	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/fm/pkg/key"
)

type DefaultMode struct {
	*Mode
}

func (*DefaultMode) GetName() string {
	return "default"
}

func createDefaultMode() *DefaultMode {
	defaultMode := &DefaultMode{
		&Mode{
			keyBindings: &KeyBindings{
				OnKeys: map[key.Key]*Action{
					key.GetKey("j"): {
						Help: "down",
						Commands: []*command.Command{
							{
								Func: command.FocusNext,
							},
						},
					},
					key.GetKey("k"): {
						Help: "up",
						Commands: []*command.Command{
							{
								Func: command.FocusPrevious,
							},
						},
					},
					key.GetKey("l"): {
						Help: "enter",
						Commands: []*command.Command{
							{
								Func: command.Enter,
							},
						},
					},
					key.GetKey("h"): {
						Help: "back",
						Commands: []*command.Command{
							{
								Func: command.Back,
							},
						},
					},
					key.GetKey("m"): {
						Help: "mark save",
						Commands: []*command.Command{
							{
								Func: command.SwitchMode,
								Args: []interface{}{"mark-save"},
							},
						},
					},
					key.GetKey("`"): {
						Help: "mark load",
						Commands: []*command.Command{
							{
								Func: command.SwitchMode,
								Args: []interface{}{"mark-load"},
							},
						},
					},
					key.GetKey("d"): {
						Help: "delete",
						Commands: []*command.Command{
							{
								Func: command.SwitchMode,
								Args: []interface{}{"delete"},
							},
						},
					},
					key.GetKey("p"): {
						Help: "copy",
						Commands: []*command.Command{
							{
								Func: command.PasteSelections,
								Args: []interface{}{"copy"},
							},
						},
					},
					key.GetKey("x"): {
						Help: "cut",
						Commands: []*command.Command{
							{
								Func: command.PasteSelections,
								Args: []interface{}{"cut"},
							},
						},
					},
					key.GetKey("r"): {
						Help: "rename",
						Commands: []*command.Command{
							{
								Func: command.Rename,
							},
						},
					},
					key.GetKey("ctrl+r"): {
						Help: "refresh",
						Commands: []*command.Command{
							{
								Func: command.Refresh,
							},
						},
					},
					key.GetKey("space"): {
						Help: "toggle selection",
						Commands: []*command.Command{
							{
								Func: command.ToggleSelection,
							},
						},
					},
					key.GetKey("ctrl+space"): {
						Help: "clear selection",
						Commands: []*command.Command{
							{
								Func: command.ClearSelection,
							},
						},
					},
					key.GetKey("."): {
						Help: "toggle hidden",
						Commands: []*command.Command{
							{
								Func: command.ToggleHidden,
							},
						},
					},
					key.GetKey("q"): {
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

	defaultModeConfig := config.AppConfig.DefaultModeConfig

	for k, actionConfig := range defaultModeConfig.KeyBindings.OnKeys {
		key := key.GetKey(k)

		defaultMode.keyBindings.OnKeys[key] = &Action{
			Commands: []*command.Command{},
		}

		for _, commandConfig := range actionConfig.Commands {
			defaultMode.keyBindings.OnKeys[key].Commands = append(
				defaultMode.keyBindings.OnKeys[key].Commands,
				toCommand(commandConfig),
			)
		}

		defaultMode.keyBindings.OnKeys[key].Help = actionConfig.Help
	}

	return defaultMode
}
