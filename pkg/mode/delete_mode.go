package mode

import "github.com/dinhhuy258/fm/pkg/message"

func createDeleteMode() *Mode {
	return &Mode{
		Name: "delete",
		KeyBindings: &KeyBindings{
			OnKeys: map[string]*Action{
				"d": {
					Help: "delete",
					Messages: []message.Message{
						{
							Func: message.DeleteCurrent,
						},
					},
				},
				"s": {
					Help: "delete selections",
					Messages: []message.Message{
						{
							Func: message.DeleteSelections,
						},
					},
				},
				"esc": {
					Help: "cancel",
					Messages: []message.Message{
						{
							Func: message.PopMode,
						},
					},
				},
				"q": {
					Help: "quit",
					Messages: []message.Message{
						{
							Func: message.Quit,
						},
					},
				},
			},
		},
	}
}

