package app

import (
	"github.com/dinhhuy258/fm/pkg/app/command"
	"github.com/dinhhuy258/fm/pkg/config"
)

type CustomMode struct {
	*Mode
	name string
}

func (mode *CustomMode) GetName() string {
	return mode.name
}

func createCustomMode(name string, keyBindings config.KeyBindingsConfig) *CustomMode {
	customMode := CustomMode{
		name: name,
		Mode: &Mode{
			KeyBindings: &KeyBindings{
				OnKeys: map[string]*command.Command{},
			},
		},
	}

	for key, commandConfig := range keyBindings.OnKeys {
		customMode.KeyBindings.OnKeys[key] = toCommand(commandConfig)
	}

	return &customMode
}

func (mode *CustomMode) GetHelp(app *App) ([]string, []string) {
	keys := make([]string, 0)
	helps := make([]string, 0)

	return keys, helps
}

// TODO: Find a better way to convert string to command
func toCommand(commandConfig *config.CommandConfig) *command.Command {
	commandString := commandConfig.Name

	switch commandString {
	case "PopMode":
		return &command.Command{
			Help: "cancel",
			Func: command.PopMode,
		}
	case "ChangeDirectory":
		return &command.Command{
			Help: "change directory",
			Func: command.ChangeDirectory,
			Args: commandConfig.Args,
		}
	}

	return nil
}
