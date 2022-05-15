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
			keyBindings: &KeyBindings{
				OnAlphabet: &Action{
					Help: "mark save",
					Commands: []*command.Command{
						{
							Func: command.MarkSave,
						},
					},
				},
				OnKeys: map[string]*Action{
					"esc": {
						Help: "cancel",
						Commands: []*command.Command{
							{
								Func: command.PopMode,
							},
						},
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
			keyBindings: &KeyBindings{
				OnAlphabet: &Action{
					Help: "mark load",
					Commands: []*command.Command{
						{
							Func: command.MarkLoad,
						},
					},
				},
				OnKeys: map[string]*Action{
					"esc": {
						Help: "cancel",
						Commands: []*command.Command{
							{
								Func: command.PopMode,
							},
						},
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
	helps := make([]*Help, 0, len(m.keyBindings.OnKeys)+len(app.Marks))

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
