package app

import (
	"errors"
	"log"

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

type Mode struct {
	*Mode
	name        string
	keyBindings *KeyBindings
	helps       []*Help
}

func (m *Mode) GetName() string {
	return m.name
}

func (m *Mode) GetKeyBindings() *KeyBindings {
	return m.keyBindings
}

func (m *Mode) GetHelp() []*Help {
	return m.helps
}

type Modes struct {
	Modes        []*Mode
	BuiltinModes map[string]*Mode
	CustomModes  map[string]*Mode
}

func CreateAllModes(marks map[string]string) *Modes {
	builtinModes := make(map[string]*Mode)

	for _, builtinMode := range config.AppConfig.BuiltinModeConfigs {
		builtinModes[builtinMode.Name] = createMode(builtinMode.Name, builtinMode.KeyBindings)
	}

	customModes := make(map[string]*Mode)
	for _, customMode := range config.AppConfig.CustomModeConfigs {
		customModes[customMode.Name] = createMode(customMode.Name, customMode.KeyBindings)
	}

	return &Modes{
		Modes:        make([]*Mode, 0, 5),
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

func (m *Modes) Peek() *Mode {
	return m.Modes[len(m.Modes)-1]
}

func createMode(name string, keyBindings config.KeyBindingsConfig) *Mode {
	mode := Mode{
		name: name,
		keyBindings: &KeyBindings{
			OnKeys:  map[key.Key]*Action{},
			Default: nil,
		},
		helps: []*Help{},
	}

	for k, actionConfig := range keyBindings.OnKeys {
		key := key.GetKey(k)

		mode.keyBindings.OnKeys[key] = &Action{
			Messages: []*message.Message{},
		}

		for _, messageConfig := range actionConfig.Messages {
			message, err := message.NewMessage(messageConfig.Name, messageConfig.Args...)
			if err != nil {
				log.Fatalf("message not found: %s", messageConfig.Name)
			}

			mode.keyBindings.OnKeys[key].Messages = append(
				mode.keyBindings.OnKeys[key].Messages,
				message,
			)
		}

		mode.helps = append(mode.helps, &Help{
			Key: key,
			Msg: actionConfig.Help,
		})
	}

	if keyBindings.Default != nil {
		mode.keyBindings.Default = &Action{
			Messages: []*message.Message{},
		}

		for _, messageConfig := range keyBindings.Default.Messages {
			message, err := message.NewMessage(messageConfig.Name, messageConfig.Args...)
			if err != nil {
				log.Fatalf("message not found: %s", messageConfig.Name)
			}

			mode.keyBindings.Default.Messages = append(
				mode.keyBindings.Default.Messages,
				message,
			)
		}
	}

	return &mode
}
