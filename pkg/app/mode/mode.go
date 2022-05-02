package mode

import (
	"errors"

	"github.com/dinhhuy258/fm/pkg/app/command"
	"github.com/dinhhuy258/fm/pkg/app/context"
)

var (
	ErrModeNotFound = errors.New("mode not found")
	ErrEmptyModes   = errors.New("empty modes")
)

type KeyBindings struct {
	OnKeys     map[string]*command.Command
	OnAlphabet *command.Command
}

type Mode struct {
	Name        string
	KeyBindings *KeyBindings
}

type Modes struct {
	Modes        []*Mode
	BuiltinModes map[string]*Mode
}

func NewModes() *Modes {
	builtinModes := make(map[string]*Mode)
	builtinModes["default"] = createDefaultMode()
	builtinModes["delete"] = createDeleteMode()
	builtinModes["mark-save"] = createMarkSaveMode()
	builtinModes["mark-load"] = createMarkLoadMode()

	return &Modes{
		Modes:        make([]*Mode, 0, 5),
		BuiltinModes: builtinModes,
	}
}

func (m *Modes) Push(mode string) error {
	builtinMode, hasMode := m.BuiltinModes[mode]
	if !hasMode {
		return ErrModeNotFound
	}

	m.Modes = append(m.Modes, builtinMode)

	return nil
}

func (m *Modes) Pop() error {
	if len(m.Modes) <= 1 {
		return ErrEmptyModes
	}

	m.Modes = m.Modes[:len(m.Modes)-1]

	return nil
}

func (m *Modes) Peek() *Mode {
	return m.Modes[len(m.Modes)-1]
}

func (m *Mode) GetHelp(state *context.State) ([]string, []string) {
	if m.Name == "mark load" {
		keys := make([]string, 0, len(m.KeyBindings.OnKeys)+len(state.Marks))
		helps := make([]string, 0, len(m.KeyBindings.OnKeys)+len(state.Marks))

		for k, m := range state.Marks {
			keys = append(keys, k)
			helps = append(helps, m)
		}

		for k, a := range m.KeyBindings.OnKeys {
			keys = append(keys, k)
			helps = append(helps, a.Help)
		}

		return keys, helps
	}

	keys := make([]string, 0, len(m.KeyBindings.OnKeys)+1)
	helps := make([]string, 0, len(m.KeyBindings.OnKeys)+1)
	keybindings := m.KeyBindings

	if keybindings.OnAlphabet != nil {
		keys = append(keys, "alphabet")
		helps = append(helps, keybindings.OnAlphabet.Help)
	}

	for k, a := range keybindings.OnKeys {
		keys = append(keys, k)
		helps = append(helps, a.Help)
	}

	return keys, helps
}
