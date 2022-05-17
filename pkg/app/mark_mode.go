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
						{
							Func: command.PopMode,
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

	marks map[string]string
}

func createMarkLoadMode(marks map[string]string) *MarkLoadMode {
	return &MarkLoadMode{
		Mode: &Mode{
			keyBindings: &KeyBindings{
				OnAlphabet: &Action{
					Help: "mark load",
					Commands: []*command.Command{
						{
							Func: command.MarkLoad,
						},
						{
							Func: command.PopMode,
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
		marks: marks,
	}
}

func (*MarkLoadMode) GetName() string {
	return "mark-load"
}

func (m *MarkLoadMode) GetHelp() []*Help {
	helps := make([]*Help, 0, len(m.keyBindings.OnKeys)+len(m.marks))

	for key, mark := range m.marks {
		helps = append(helps, &Help{
			Key: key,
			Msg: mark,
		})
	}

	for key, command := range m.keyBindings.OnKeys {
		helps = append(helps, &Help{
			Key: key,
			Msg: command.Help,
		})
	}

	return helps
}
