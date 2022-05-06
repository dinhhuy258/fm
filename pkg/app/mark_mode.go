package app

import (
	"github.com/dinhhuy258/fm/pkg/app/command"
)

type MarkSaveMode struct {
	*Mode
}

func (*MarkSaveMode) GetName() string {
	return "mark-save"
}

func createMarkSaveMode() *MarkSaveMode {
	return &MarkSaveMode{
		Mode: &Mode{
			KeyBindings: &KeyBindings{
				OnAlphabet: &command.Command{
					Help: "mark save",
					Func: command.MarkSave,
				},
				OnKeys: map[string]*command.Command{
					"esc": {
						Help: "cancel",
						Func: command.PopMode,
					},
				},
			},
		},
	}
}

type MarkLoadMode struct {
	*Mode
}

func createMarkLoadMode() *MarkLoadMode {
	return &MarkLoadMode{
		Mode: &Mode{
			KeyBindings: &KeyBindings{
				OnAlphabet: &command.Command{
					Help: "mark load",
					Func: command.MarkLoad,
				},
				OnKeys: map[string]*command.Command{
					"esc": {
						Help: "cancel",
						Func: command.PopMode,
					},
				},
			},
		},
	}
}

func (*MarkLoadMode) GetName() string {
	return "mark-load"
}

func (m *MarkLoadMode) GetHelp(app *App) []*Help {
	helps := make([]*Help, 0, len(m.KeyBindings.OnKeys)+len(app.Marks))

	for key, mark := range app.Marks {
		helps = append(helps, &Help{
			Key: key,
			Msg: mark,
		})
	}

	for key, command := range m.GetKeyBindings().OnKeys {
		helps = append(helps, &Help{
			Key: key,
			Msg: command.Help,
		})
	}

	return helps
}
