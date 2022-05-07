package app

import (
	"github.com/dinhhuy258/fm/pkg/app/command"
)

type DeleteMode struct {
	*Mode
}

func (*DeleteMode) GetName() string {
	return "delete"
}

func createDeleteMode() *DeleteMode {
	return &DeleteMode{
		&Mode{
			KeyBindings: &KeyBindings{
				OnKeys: map[string]*Action{
					"d": {
						Help: "delete",
						Commands: []*command.Command{
							{
								Func: command.DeleteCurrent,
							},
						},
					},
					"s": {
						Help: "delete selections",
						Commands: []*command.Command{
							{
								Func: command.DeleteSelections,
							},
						},
					},
					"esc": {
						Help: "cancel",
						Commands: []*command.Command{
							{
								Func: command.PopMode,
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
