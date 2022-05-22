package app

import (
	"github.com/dinhhuy258/fm/pkg/app/command"
	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/fm/pkg/key"
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
				OnKeys:  map[key.Key]*Action{},
				Default: nil,
			},
		},
		helps: []*Help{},
	}

	for k, actionConfig := range keyBindings.OnKeys {
		key := key.GetKey(k)

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

	if keyBindings.Default != nil {
		customMode.keyBindings.Default = &Action{
			Commands: []*command.Command{},
		}

		for _, commandConfig := range keyBindings.Default.Commands {
			customMode.keyBindings.Default.Commands = append(
				customMode.keyBindings.Default.Commands,
				toCommand(commandConfig),
			)
		}
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
	case "FocusNext":
		return &command.Command{
			Func: command.FocusNext,
		}
	case "FocusPrevious":
		return &command.Command{
			Func: command.FocusPrevious,
		}
	case "FocusFirst":
		return &command.Command{
			Func: command.FocusFirst,
		}
	case "FocusLast":
		return &command.Command{
			Func: command.FocusLast,
		}
	case "Enter":
		return &command.Command{
			Func: command.Enter,
		}
	case "Back":
		return &command.Command{
			Func: command.Back,
		}
	case "ChangeDirectory":
		return &command.Command{
			Func: command.ChangeDirectory,
			Args: commandConfig.Args,
		}
	case "PasteSelections":
		return &command.Command{
			Func: command.ChangeDirectory,
			Args: commandConfig.Args,
		}
	case "UpdateInputBufferFromKey":
		return &command.Command{
			Func: command.UpdateInputBufferFromKey,
		}
	case "SetInputBuffer":
		return &command.Command{
			Func: command.SetInputBuffer,
			Args: commandConfig.Args,
		}
	case "NewFileFromInput":
		return &command.Command{
			Func: command.NewFileFromInput,
		}
	case "DeleteCurrent":
		return &command.Command{
			Func: command.DeleteCurrent,
		}
	case "DeleteSelections":
		return &command.Command{
			Func: command.DeleteSelections,
		}
	case "SearchFromInput":
		return &command.Command{
			Func: command.SearchFromInput,
		}
	case "Refresh":
		return &command.Command{
			Func: command.Refresh,
		}
	case "ToggleSelection":
		return &command.Command{
			Func: command.ToggleSelection,
		}
	case "ClearSelection":
		return &command.Command{
			Func: command.ClearSelection,
		}
	case "MarkSave":
		return &command.Command{
			Func: command.MarkSave,
		}
	case "MarkLoad":
		return &command.Command{
			Func: command.MarkLoad,
		}
	}

	return nil
}
