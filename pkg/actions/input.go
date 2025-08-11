package actions

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dinhhuy258/fm/pkg/config"
)

// executeSetInputBuffer handles input buffer manipulation and shows input
func (ah *ActionHandler) executeSetInputBuffer(message *config.MessageConfig) tea.Cmd {
	return func() tea.Msg {
		value := ""
		if len(message.Args) > 0 {
			value = message.Args[0]
		}

		return SetInputBufferMessage{
			Value:     value,
			ShowInput: true, // Always show input when setting buffer
		}
	}
}

// executeUpdateInputBufferFromKey handles input buffer updates from key
func (ah *ActionHandler) executeUpdateInputBufferFromKey(message *config.MessageConfig) tea.Cmd {
	return func() tea.Msg {
		return UpdateInputBufferFromKeyMessage{}
	}
}

// executeLog handles logging messages
func (ah *ActionHandler) executeLog(message *config.MessageConfig, level string) tea.Cmd {
	return func() tea.Msg {
		if len(message.Args) == 0 {
			return LogMessage{
				Level:   level,
				Message: fmt.Sprintf("Empty %s message", level),
			}
		}

		return LogMessage{
			Level:   level,
			Message: message.Args[0],
		}
	}
}
