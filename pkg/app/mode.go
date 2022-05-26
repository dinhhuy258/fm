package app

import (
	"errors"

	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/fm/pkg/key"
	"github.com/dinhhuy258/fm/pkg/message"
)

var (
	ErrModeNotFound = errors.New("mode not found")
	ErrEmptyModes   = errors.New("empty modes")
)

type Action struct {
	Help     string
	Messages []*message.Message
}

type Help struct {
	Key key.Key
	Msg string
}

type KeyBindings struct {
	OnKeys     map[key.Key]*Action
	OnAlphabet *Action
	Default    *Action
}

// TODO: Remove interface???
type IMode interface {
	GetName() string
	GetKeyBindings() *KeyBindings
	GetHelp() []*Help
}

type Mode struct {
	IMode
	keyBindings *KeyBindings
}

func (m *Mode) GetKeyBindings() *KeyBindings {
	return m.keyBindings
}

type Modes struct {
	Modes        []IMode
	BuiltinModes map[string]IMode
	CustomModes  map[string]IMode
}

func CreateAllModes(marks map[string]string) *Modes {
	builtinModes := make(map[string]IMode)

	for _, builtinMode := range config.AppConfig.BuiltinModeConfigs {
		builtinModes[builtinMode.Name] = createCustomMode(builtinMode.Name, builtinMode.KeyBindings)
	}

	customModes := make(map[string]IMode)
	for _, customMode := range config.AppConfig.CustomModeConfigs {
		customModes[customMode.Name] = createCustomMode(customMode.Name, customMode.KeyBindings)
	}

	return &Modes{
		Modes:        make([]IMode, 0, 5),
		BuiltinModes: builtinModes,
		CustomModes:  customModes,
	}
}

func (m *Modes) Push(mode string) error {
	if builtinMode, hasBuiltinMode := m.BuiltinModes[mode]; hasBuiltinMode {
		m.Modes = append(m.Modes, builtinMode)

		return nil
	}

	if customMode, hasCustomMode := m.CustomModes[mode]; hasCustomMode {
		m.Modes = append(m.Modes, customMode)

		return nil
	}

	return ErrModeNotFound
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

func (m *Mode) GetHelp() []*Help {
	helps := make([]*Help, 0, len(m.keyBindings.OnKeys)+1)
	keybindings := m.keyBindings

	if keybindings.OnAlphabet != nil {
		helps = append(helps, &Help{
			Key: "alphabet",
			Msg: keybindings.OnAlphabet.Help,
		})
	}

	for key, message := range keybindings.OnKeys {
		helps = append(helps, &Help{
			Key: key,
			Msg: message.Help,
		})
	}

	return helps
}
