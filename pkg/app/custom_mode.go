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

	for key, commandConfig := range keyBindings.OnKeys {
		customMode.KeyBindings.OnKeys[key] = &Action{
			Commands: []*command.Command{
				toCommand(commandConfig),
			},
		}
		customMode.helps = append(customMode.helps, &Help{
			Key: key,
			Msg: commandConfig.Help,
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
