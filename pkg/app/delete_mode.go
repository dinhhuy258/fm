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
			keyBindings: &KeyBindings{
				OnKeys: map[string]*Action{
					"d": {
						Help: "delete",
						Commands: []*command.Command{
							{
								Func: command.DeleteCurrent,
							},
							{
								Func: command.PopMode,
							},
						},
					},
					"s": {
						Help: "delete selections",
						Commands: []*command.Command{
							{
								Func: command.DeleteSelections,
							},
							{
								Func: command.PopMode,
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
