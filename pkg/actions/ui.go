package actions

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/dinhhuy258/fm/pkg/config"
)

// executeSwitchMode handles mode switching
func (ah *ActionHandler) executeSwitchMode(message *config.MessageConfig) tea.Cmd {
	return func() tea.Msg {
		return ModeChangedMessage{
			Mode: message.Args[0],
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

// executeQuit handles application quit
func (ah *ActionHandler) executeQuit(_ *config.MessageConfig) tea.Cmd {
	return tea.Quit
}