package app

import "errors"

var (
	ErrModeNotFound = errors.New("mode not found")
	ErrEmptyModes   = errors.New("empty modes")
)

type Message struct {
	f    func(app *App, params ...interface{}) error
	args []interface{}
}

type Action struct {
	help     string
	messages []Message
}

type KeyBindings struct {
	onKeys map[string]*Action
}

type Mode struct {
	name        string
	keyBindings *KeyBindings
}

type Modes struct {
	modes        []*Mode
	builtinModes map[string]*Mode
}

func NewModes() *Modes {
	builtinModes := make(map[string]*Mode)
	builtinModes["default"] = createDefaultMode()
	builtinModes["delete"] = createDeleteMode()

	return &Modes{
		modes:        make([]*Mode, 0, 5),
		builtinModes: builtinModes,
	}
}

func (m *Modes) Push(mode string) error {
	builtinMode, hasMode := m.builtinModes[mode]
	if !hasMode {
		return ErrModeNotFound
	}

	m.modes = append(m.modes, builtinMode)

	return nil
}

func (m *Modes) Pop() error {
	if len(m.modes) <= 1 {
		return ErrEmptyModes
	}

	m.modes = m.modes[:len(m.modes)-1]

	return nil
}

func (m *Modes) Peek() *Mode {
	return m.modes[len(m.modes)-1]
}
