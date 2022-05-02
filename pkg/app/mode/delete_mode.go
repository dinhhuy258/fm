package mode

import (
	"github.com/dinhhuy258/fm/pkg/app/command"
)

func createDeleteMode() *Mode {
	return &Mode{
		Name: "delete",
		KeyBindings: &KeyBindings{
			OnKeys: map[string]*command.Command{
				"d": {
					Help: "delete",
					Func: command.DeleteCurrent,
				},
				"s": {
					Help: "delete selections",
					Func: command.DeleteSelections,
				},
				"esc": {
					Help: "cancel",
					Func: command.PopMode,
				},
				"q": {
					Help: "quit",
					Func: command.Quit,
				},
			},
		},
	}
}
