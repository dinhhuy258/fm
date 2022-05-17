package app

import (
	"github.com/dinhhuy258/fm/pkg/app/command"
	"github.com/dinhhuy258/fm/pkg/config"
)

type CustomMode struct {
	*Mode
	name  string
	helps []*Help
}

func (m *CustomMode) GetName() string {
	return m.name
}

func (m *CustomMode) GetHelp() []*Help {
	return m.helps
}

func createCustomMode(name string, keyBindings config.KeyBindingsConfig) *CustomMode {
	customMode := CustomMode{
		name: name,
		Mode: &Mode{
			keyBindings: &KeyBindings{
				OnKeys: map[string]*Action{},
			},
		},
		helps: []*Help{},
	}

	for key, actionConfig := range keyBindings.OnKeys {
		customMode.keyBindings.OnKeys[key] = &Action{
			Commands: []*command.Command{},
		}

		for _, commandConfig := range actionConfig.Commands {
			customMode.keyBindings.OnKeys[key].Commands = append(
				customMode.keyBindings.OnKeys[key].Commands,
				toCommand(commandConfig),
			)
		}

		customMode.helps = append(customMode.helps, &Help{
			Key: key,
			Msg: actionConfig.Help,
		})
	}

	return &customMode
}

// TODO: Find a better way to convert string to command
func toCommand(commandConfig *config.CommandConfig) *command.Command {
	commandString := commandConfig.Name

	switch commandString {
	case "Quit":
		return &command.Command{
			Func: command.Quit,
		}
	case "SwitchMode":
		return &command.Command{
			Func: command.SwitchMode,
			Args: commandConfig.Args,
		}
	case "PopMode":
		return &command.Command{
			Func: command.PopMode,
		}
	case "ChangeDirectory":
		return &command.Command{
			Func: command.ChangeDirectory,
			Args: commandConfig.Args,
		}
	}

	return nil
}
