package actions

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/dinhhuy258/fm/pkg/config"
)

// executeBashExec handles bash command execution
func (ah *ActionHandler) executeBashExec(
	message *config.MessageConfig,
) tea.Cmd {
	return func() tea.Msg {
		script := message.Args[0]

		return BashExecMessage{
			Script: script,
		}
	}
}

// executeBashExec handles bash command execution
func (ah *ActionHandler) executeBashExecSilently(
	message *config.MessageConfig,
) tea.Cmd {
	return func() tea.Msg {
		script := message.Args[0]

		return BashExecSilentlyMessage{
			Script: script,
		}
	}
}