package mode

import (
	"github.com/dinhhuy258/fm/pkg/app/command"
)

func createDeleteMode() *Mode {
	return &Mode{
		Name: "delete",
		KeyBindings: &KeyBindings{
			OnKeys: map[string]*Action{
				"d": {
					Help: "delete",
					Commands: []command.Command{
						{
							Func: command.DeleteCurrent,
						},
					},
				},
				"s": {
					Help: "delete selections",
					Commands: []command.Command{
						{
							Func: command.DeleteSelections,
						},
					},
				},
				"esc": {
					Help: "cancel",
					Commands: []command.Command{
						{
							Func: command.PopMode,
						},
					},
				},
				"q": {
					Help: "quit",
					Commands: []command.Command{
						{
							Func: command.Quit,
						},
					},
				},
			},
		},
	}
}
