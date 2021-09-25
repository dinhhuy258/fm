package mode

import (
	"errors"

	"github.com/dinhhuy258/fm/pkg/message"
)

var (
	ErrModeNotFound = errors.New("mode not found")
	ErrEmptyModes   = errors.New("empty modes")
)

type Action struct {
	Help     string
	Messages []message.Message
}

type KeyBindings struct {
	OnKeys     map[string]*Action
	OnAlphabet *Action
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
