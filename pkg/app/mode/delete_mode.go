package mode

import (
	message2 "github.com/dinhhuy258/fm/pkg/app/message"
)

func createDeleteMode() *Mode {
	return &Mode{
		Name: "delete",
		KeyBindings: &KeyBindings{
			OnKeys: map[string]*Action{
				"d": {
					Help: "delete",
					Messages: []message2.Message{
						{
							Func: message2.DeleteCurrent,
						},
					},
				},
				"s": {
					Help: "delete selections",
					Messages: []message2.Message{
						{
							Func: message2.DeleteSelections,
						},
					},
				},
				"esc": {
					Help: "cancel",
					Messages: []message2.Message{
						{
							Func: message2.PopMode,
						},
					},
				},
				"q": {
					Help: "quit",
					Messages: []message2.Message{
						{
							Func: message2.Quit,
						},
					},
				},
			},
		},
	}
}
