package app

import (
	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/fm/pkg/key"
	"github.com/dinhhuy258/fm/pkg/message"
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
			Messages: []*message.Message{},
		}

		for _, messageConfig := range actionConfig.Messages {
			customMode.keyBindings.OnKeys[key].Messages = append(
				customMode.keyBindings.OnKeys[key].Messages,
				toMessage(messageConfig),
			)
		}

		customMode.helps = append(customMode.helps, &Help{
			Key: key,
			Msg: actionConfig.Help,
		})
	}

	if keyBindings.Default != nil {
		customMode.keyBindings.Default = &Action{
			Messages: []*message.Message{},
		}

		for _, messageConfig := range keyBindings.Default.Messages {
			customMode.keyBindings.Default.Messages = append(
				customMode.keyBindings.Default.Messages,
				toMessage(messageConfig),
			)
		}
	}

	return &customMode
}

// TODO: Find a better way to convert string to message
func toMessage(messageConfig *config.MessageConfig) *message.Message {
	messageName := messageConfig.Name

	switch messageName {
	case "Quit":
		return &message.Message{
			Func: message.Quit,
		}
	case "SwitchMode":
		return &message.Message{
			Func: message.SwitchMode,
			Args: messageConfig.Args,
		}
	case "PopMode":
		return &message.Message{
			Func: message.PopMode,
		}
	case "FocusNext":
		return &message.Message{
			Func: message.FocusNext,
		}
	case "FocusPrevious":
		return &message.Message{
			Func: message.FocusPrevious,
		}
	case "FocusFirst":
		return &message.Message{
			Func: message.FocusFirst,
		}
	case "FocusLast":
		return &message.Message{
			Func: message.FocusLast,
		}
	case "Enter":
		return &message.Message{
			Func: message.Enter,
		}
	case "Back":
		return &message.Message{
			Func: message.Back,
		}
	case "ChangeDirectory":
		return &message.Message{
			Func: message.ChangeDirectory,
			Args: messageConfig.Args,
		}
	case "PasteSelections":
		return &message.Message{
			Func: message.ChangeDirectory,
			Args: messageConfig.Args,
		}
	case "UpdateInputBufferFromKey":
		return &message.Message{
			Func: message.UpdateInputBufferFromKey,
		}
	case "SetInputBuffer":
		return &message.Message{
			Func: message.SetInputBuffer,
			Args: messageConfig.Args,
		}
	case "NewFileFromInput":
		return &message.Message{
			Func: message.NewFileFromInput,
		}
	case "DeleteCurrent":
		return &message.Message{
			Func: message.DeleteCurrent,
		}
	case "DeleteSelections":
		return &message.Message{
			Func: message.DeleteSelections,
		}
	case "SearchFromInput":
		return &message.Message{
			Func: message.SearchFromInput,
		}
	case "Refresh":
		return &message.Message{
			Func: message.Refresh,
		}
	case "ToggleSelection":
		return &message.Message{
			Func: message.ToggleSelection,
		}
	case "ClearSelection":
		return &message.Message{
			Func: message.ClearSelection,
		}
	case "MarkSave":
		return &message.Message{
			Func: message.MarkSave,
		}
	case "MarkLoad":
		return &message.Message{
			Func: message.MarkLoad,
		}
	}

	return nil
}
