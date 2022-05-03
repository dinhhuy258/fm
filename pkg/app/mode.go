package app

import (
	"errors"

	"github.com/dinhhuy258/fm/pkg/app/command"
)

var (
	ErrModeNotFound = errors.New("mode not found")
	ErrEmptyModes   = errors.New("empty modes")
)

type KeyBindings struct {
	OnKeys     map[string]*command.Command
	OnAlphabet *command.Command
}

type IMode interface {
	GetName() string
	GetKeyBindings() *KeyBindings
	GetHelp(state *State) ([]string, []string)
}

type Mode struct {
	IMode
	KeyBindings *KeyBindings
}

func (m *Mode) GetKeyBindings() *KeyBindings {
	return m.KeyBindings
}

type Modes struct {
	Modes        []IMode
	BuiltinModes map[string]IMode
}

func NewModes() *Modes {
	builtinModes := make(map[string]IMode)
	builtinModes["default"] = createDefaultMode()
	builtinModes["delete"] = createDeleteMode()
	builtinModes["mark-save"] = createMarkSaveMode()
	builtinModes["mark-load"] = createMarkLoadMode()

	return &Modes{
		Modes:        make([]IMode, 0, 5),
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

func (m *Modes) Peek() IMode {
	return m.Modes[len(m.Modes)-1]
}

func (m *Mode) GetHelp(*State) ([]string, []string) {
	keys := make([]string, 0, len(m.GetKeyBindings().OnKeys)+1)
	helps := make([]string, 0, len(m.GetKeyBindings().OnKeys)+1)
	keybindings := m.GetKeyBindings()

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
