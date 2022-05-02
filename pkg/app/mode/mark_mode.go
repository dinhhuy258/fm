package mode

import (
	message2 "github.com/dinhhuy258/fm/pkg/app/message"
)

func createMarkSaveMode() *Mode {
	return &Mode{
		Name: "mark save",
		KeyBindings: &KeyBindings{
			OnAlphabet: &Action{
				Help: "mark save",
				Messages: []message2.Message{
					{
						Func: message2.MarkSave,
					},
				},
			},
			OnKeys: map[string]*Action{
				"esc": {
					Help: "cancel",
					Messages: []message2.Message{
						{
							Func: message2.PopMode,
						},
					},
				},
			},
		},
	}
}

func createMarkLoadMode() *Mode {
	return &Mode{
		Name: "mark load",
		KeyBindings: &KeyBindings{
			OnAlphabet: &Action{
				Help: "mark load",
				Messages: []message2.Message{
					{
						Func: message2.MarkLoad,
					},
				},
			},
			OnKeys: map[string]*Action{
				"esc": {
					Help: "cancel",
					Messages: []message2.Message{
						{
							Func: message2.PopMode,
						},
					},
				},
			},
		},
	}
}
