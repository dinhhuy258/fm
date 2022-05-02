package mode

import (
	"github.com/dinhhuy258/fm/pkg/app/command"
)

func createMarkSaveMode() *Mode {
	return &Mode{
		Name: "mark save",
		KeyBindings: &KeyBindings{
			OnAlphabet: &command.Command{
				Help: "mark save",
				Func: command.MarkSave,
			},
			OnKeys: map[string]*command.Command{
				"esc": {
					Help: "cancel",
					Func: command.PopMode,
				},
			},
		},
	}
}

func createMarkLoadMode() *Mode {
	return &Mode{
		Name: "mark load",
		KeyBindings: &KeyBindings{
			OnAlphabet: &command.Command{
				Help: "mark load",
				Func: command.MarkLoad,
			},
			OnKeys: map[string]*command.Command{
				"esc": {
					Help: "cancel",
					Func: command.PopMode,
				},
			},
		},
	}
}
