package mode

import (
	"github.com/dinhhuy258/fm/pkg/app/command"
	"github.com/dinhhuy258/fm/pkg/app/context"
)

type MarkSaveMode struct {
	*Mode
}

func (_ *MarkSaveMode) GetName() string {
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

func (_ *MarkLoadMode) GetName() string {
	return "mark-load"
}

func (m *MarkLoadMode) GetHelp(state *context.State) ([]string, []string) {
	keys := make([]string, 0, len(m.GetKeyBindings().OnKeys)+len(state.Marks))
	helps := make([]string, 0, len(m.GetKeyBindings().OnKeys)+len(state.Marks))

	for k, m := range state.Marks {
		keys = append(keys, k)
		helps = append(helps, m)
	}

	for k, a := range m.GetKeyBindings().OnKeys {
		keys = append(keys, k)
		helps = append(helps, a.Help)
	}

	return keys, helps
}
