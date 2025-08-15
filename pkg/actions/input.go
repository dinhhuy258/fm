package actions

import (
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
			Value: value,
		}
	}
}

// executeUpdateInputBufferFromKey handles input buffer updates from key
func (ah *ActionHandler) executeUpdateInputBufferFromKey(message *config.MessageConfig) tea.Cmd {
	return func() tea.Msg {
		return UpdateInputBufferFromKeyMessage{}
	}
}

// executeLogSuccess handles success logging messages
func (ah *ActionHandler) executeLogSuccess(message *config.MessageConfig) tea.Cmd {
	return func() tea.Msg {
		return LogMessage{
			Level:   LogLevelSuccess,
			Message: message.Args[0],
		}
	}
}

// executeLogError handles error logging messages
func (ah *ActionHandler) executeLogError(message *config.MessageConfig) tea.Cmd {
	return func() tea.Msg {
		return LogMessage{
			Level:   LogLevelError,
			Message: message.Args[0],
		}
	}
}

// executeLogInfo handles info logging messages
func (ah *ActionHandler) executeLogInfo(message *config.MessageConfig) tea.Cmd {
	return func() tea.Msg {
		return LogMessage{
			Level:   LogLevelInfo,
			Message: message.Args[0],
		}
	}
}

// executeLogWarning handles warning logging messages
func (ah *ActionHandler) executeLogWarning(message *config.MessageConfig) tea.Cmd {
	return func() tea.Msg {
		return LogMessage{
			Level:   LogLevelWarning,
			Message: message.Args[0],
		}
	}
}