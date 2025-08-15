package actions

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/dinhhuy258/fm/pkg/config"
)

// executeBashExec handles bash command execution
func (ah *ActionHandler) executeBashExec(
	message *config.MessageConfig,
	silent bool,
) tea.Cmd {
	return func() tea.Msg {
		script := message.Args[0]
		if silent {
			return BashExecSilentlyMessage{Script: script}
		}

		return BashExecMessage{Script: script}
	}
}
