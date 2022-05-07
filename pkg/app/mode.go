package app

import (
	"errors"

	"github.com/dinhhuy258/fm/pkg/app/command"
	"github.com/dinhhuy258/fm/pkg/config"
)

var (
	ErrModeNotFound = errors.New("mode not found")
	ErrEmptyModes   = errors.New("empty modes")
)

type Action struct {
	Help     string
	Commands []*command.Command
}

type Help struct {
	Key string
	Msg string
}

type KeyBindings struct {
	OnKeys     map[string]*Action
	OnAlphabet *Action
}

type IMode interface {
	GetName() string
	GetKeyBindings() *KeyBindings
	GetHelp(app *App) []*Help
	OnModeStarted(app *App)
}

type Mode struct {
	IMode
	KeyBindings *KeyBindings
}

// TODO Considering remove this method
func (m *Mode) GetKeyBindings() *KeyBindings {
	return m.KeyBindings
}

func (m *Mode) OnModeStarted(*App) {
}

type Modes struct {
	Modes        []IMode
	BuiltinModes map[string]IMode
	CustomModes  map[string]IMode
}

func NewModes() *Modes {
	builtinModes := make(map[string]IMode)
	builtinModes["default"] = createDefaultMode()
	builtinModes["delete"] = createDeleteMode()
	builtinModes["mark-save"] = createMarkSaveMode()
	builtinModes["mark-load"] = createMarkLoadMode()
	builtinModes["search"] = createSearchMode()

	customModes := make(map[string]IMode)
	for _, customMode := range config.AppConfig.ModeConfigs {
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

func (m *Mode) GetHelp(*App) []*Help {
	helps := make([]*Help, 0, len(m.KeyBindings.OnKeys)+1)
	keybindings := m.KeyBindings

	if keybindings.OnAlphabet != nil {
		helps = append(helps, &Help{
			Key: "alphabet",
			Msg: keybindings.OnAlphabet.Help,
		})
	}

	for key, command := range keybindings.OnKeys {
		helps = append(helps, &Help{
			Key: key,
			Msg: command.Help,
		})
	}

	return helps
}
