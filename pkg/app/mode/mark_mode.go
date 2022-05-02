package mode

import (
	"github.com/dinhhuy258/fm/pkg/app/command"
)

func createMarkSaveMode() *Mode {
	return &Mode{
		Name: "mark save",
		KeyBindings: &KeyBindings{
			OnAlphabet: &Action{
				Help: "mark save",
				Commands: []command.Command{
					{
						Func: command.MarkSave,
					},
				},
			},
			OnKeys: map[string]*Action{
				"esc": {
					Help: "cancel",
					Commands: []command.Command{
						{
							Func: command.PopMode,
						},
					},
				},
			},
		},
	}
}

func createMarkLoadMode() *Mode {
	return &Mode{
		Name: "mark load",
		KeyBindings: &KeyBindings{
			OnAlphabet: &Action{
				Help: "mark load",
				Commands: []command.Command{
					{
						Func: command.MarkLoad,
					},
				},
			},
			OnKeys: map[string]*Action{
				"esc": {
					Help: "cancel",
					Commands: []command.Command{
						{
							Func: command.PopMode,
						},
					},
				},
			},
		},
	}
}
