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

func (m *CustomMode) GetHelp(app *App) []*Help {
	return m.helps
}

func createCustomMode(name string, keyBindings config.KeyBindingsConfig) *CustomMode {
	customMode := CustomMode{
		name: name,
		Mode: &Mode{
			KeyBindings: &KeyBindings{
				OnKeys: map[string]*Action{},
			},
		},
		helps: []*Help{},
	}

	for key, actionConfig := range keyBindings.OnKeys {
		customMode.KeyBindings.OnKeys[key] = &Action{
			Commands: []*command.Command{},
		}

		for _, commandConfig := range actionConfig.Commands {
			customMode.KeyBindings.OnKeys[key].Commands =
				append(customMode.KeyBindings.OnKeys[key].Commands, toCommand(commandConfig))
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
