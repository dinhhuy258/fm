package app

import (
	"errors"
	"log"

	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/fm/pkg/key"
	"github.com/dinhhuy258/fm/pkg/msg"
)

var (
	ErrModeNotFound = errors.New("mode not found")
	ErrEmptyModes   = errors.New("empty modes")
)

// Action represents an action.
type Action struct {
	messages []*msg.Message
}

// Help represents a help for the mode.
type Help struct {
	key key.Key
	msg string
}

// KeyBindings represents a key bindings config for the mode.
type KeyBindings struct {
	onKeys        map[key.Key]*Action
	defaultAction *Action
}

// Mode represents a mode.
type Mode struct {
	name        string
	keyBindings *KeyBindings
	helps       []*Help
}

// GetName returns the name of the mode.
func (m *Mode) GetName() string {
	return m.name
}

// GetKeyBindings returns the key bindings of the mode.
func (m *Mode) GetKeyBindings() *KeyBindings {
	return m.keyBindings
}

// GetHelp returns the help for the mode.
func (m *Mode) GetHelp() []*Help {
	return m.helps
}

// Modes contains a list of modes.
type Modes struct {
	// The current mode.
	modes []*Mode
	// The list of builtin modes.
	builtinModes map[string]*Mode
	// The list of user config modes.
	customModes map[string]*Mode
	// The callback function when the mode changes.
	onModeChange func(*Mode)
}

// CreateModes creates modes from config.
func CreateModes(onModeChange func(*Mode)) *Modes {
	builtinModes := make(map[string]*Mode)
	for _, builtinMode := range config.AppConfig.BuiltinModeConfigs {
		builtinModes[builtinMode.Name] = createMode(builtinMode.Name, builtinMode.KeyBindings)
	}

	customModes := make(map[string]*Mode)
	for _, customMode := range config.AppConfig.CustomModeConfigs {
		customModes[customMode.Name] = createMode(customMode.Name, customMode.KeyBindings)
	}

	return &Modes{
		modes:        make([]*Mode, 0, 5),
		builtinModes: builtinModes,
		customModes:  customModes,
		onModeChange: onModeChange,
	}
}

// Push pushes a mode to the mode stack.
func (m *Modes) Push(mode string) error {
	if builtinMode, hasBuiltinMode := m.builtinModes[mode]; hasBuiltinMode {
		m.modes = append(m.modes, builtinMode)
		m.onModeChange(builtinMode)

		return nil
	}

	if customMode, hasCustomMode := m.customModes[mode]; hasCustomMode {
		m.modes = append(m.modes, customMode)
		m.onModeChange(customMode)

		return nil
	}

	return ErrModeNotFound
}

// Pop pops a mode from the mode stack.
func (m *Modes) Pop() error {
	if len(m.modes) <= 1 {
		return ErrEmptyModes
	}

	m.modes = m.modes[:len(m.modes)-1]
	m.onModeChange(m.modes[len(m.modes)-1])

	return nil
}

// Peek returns the current mode.
func (m *Modes) Peek() *Mode {
	return m.modes[len(m.modes)-1]
}

// createMode creates a mode from config.
func createMode(name string, keyBindings config.KeyBindingsConfig) *Mode {
	mode := Mode{
		name: name,
		keyBindings: &KeyBindings{
			onKeys:        map[key.Key]*Action{},
			defaultAction: nil,
		},
		helps: []*Help{},
	}

	for k, actionConfig := range keyBindings.OnKeys {
		key := key.GetKey(k)

		mode.keyBindings.onKeys[key] = &Action{
			messages: []*msg.Message{},
		}

		for _, messageConfig := range actionConfig.Messages {
			message, err := msg.NewMessage(messageConfig.Name, messageConfig.Args...)
			if err != nil {
				log.Fatalf("message not found: %s", messageConfig.Name)
			}

			mode.keyBindings.onKeys[key].messages = append(
				mode.keyBindings.onKeys[key].messages,
				message,
			)
		}

		mode.helps = append(mode.helps, &Help{
			key: key,
			msg: actionConfig.Help,
		})
	}

	if keyBindings.Default != nil {
		mode.keyBindings.defaultAction = &Action{
			messages: []*msg.Message{},
		}

		for _, messageConfig := range keyBindings.Default.Messages {
			message, err := msg.NewMessage(messageConfig.Name, messageConfig.Args...)
			if err != nil {
				log.Fatalf("message not found: %s", messageConfig.Name)
			}

			mode.keyBindings.defaultAction.messages = append(
				mode.keyBindings.defaultAction.messages,
				message,
			)
		}
	}

	return &mode
}
