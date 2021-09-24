package mode

import "github.com/dinhhuy258/fm/pkg/message"

func createMarkSaveMode() *Mode {
	return &Mode{
		Name: "mark save",
		KeyBindings: &KeyBindings{
			OnAlphabet: &Action{
				Help: "mark save",
				Messages: []message.Message{
					{
						Func: message.MarkSave,
					},
				},
			},
			OnKeys: map[string]*Action{
				"esc": {
					Help: "cancel",
					Messages: []message.Message{
						{
							Func: message.PopMode,
						},
					},
				},
			},
		},
	}
}
