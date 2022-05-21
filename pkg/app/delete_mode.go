package app

import (
	"github.com/dinhhuy258/fm/pkg/app/command"
	"github.com/dinhhuy258/fm/pkg/key"
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
				OnKeys: map[key.Key]*Action{
					key.GetKey("d"): {
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
					key.GetKey("s"): {
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
					key.GetKey("esc"): {
						Help: "cancel",
						Commands: []*command.Command{
							{
								Func: command.PopMode,
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
}
