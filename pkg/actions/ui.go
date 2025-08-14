package actions

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dinhhuy258/fm/pkg/config"
)

// executeSwitchMode handles mode switching
func (ah *ActionHandler) executeSwitchMode(message *config.MessageConfig) tea.Cmd {
	return func() tea.Msg {
		if len(message.Args) == 0 {
			return ErrorMessage{Err: fmt.Errorf("SwitchMode requires a mode argument")}
		}

		targetMode := message.Args[0]

		return ModeChangedMessage{
			NewMode: targetMode,
		}
	}
}

// executeToggleHidden handles toggling hidden file visibility
func (ah *ActionHandler) executeToggleHidden(_ *config.MessageConfig) tea.Cmd {
	return func() tea.Msg {
		return UIMessage{Action: UIActionToggleHidden}
	}
}

// executeRefresh handles refreshing the current directory
func (ah *ActionHandler) executeRefresh(_ *config.MessageConfig) tea.Cmd {
	return func() tea.Msg {
		return UIMessage{Action: UIActionRefresh}
	}
}

// executeClearLog handles clearing the log
func (ah *ActionHandler) executeClearLog(_ *config.MessageConfig) tea.Cmd {
	return func() tea.Msg {
		return UIMessage{Action: UIActionClearLog}
	}
}

// executeToggleLog handles toggling log view visibility
func (ah *ActionHandler) executeToggleLog(_ *config.MessageConfig) tea.Cmd {
	return func() tea.Msg {
		return UIMessage{Action: UIActionToggleLog}
	}
}

// executeQuit handles application quit
func (ah *ActionHandler) executeQuit(_ *config.MessageConfig) tea.Cmd {
	return tea.Quit
}
