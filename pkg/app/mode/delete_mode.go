package mode

import (
	"github.com/dinhhuy258/fm/pkg/app/command"
)

type DeleteMode struct {
	*Mode
}

func (_ *DeleteMode) GetName() string {
	return "delete"
}

func createDeleteMode() *DeleteMode {
	return &DeleteMode{
		&Mode{
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
		},
	}
}
