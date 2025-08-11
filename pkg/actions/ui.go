package actions

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dinhhuy258/fm/pkg/config"
)

// executeSwitchMode handles mode switching
func (ah *ActionHandler) executeSwitchMode(ahssage *config.MessageConfig) tea.Cmd {
	return func() tea.Msg {
		if len(ahssage.Args) == 0 {
			return ErrorMessage{Err: fmt.Errorf("SwitchMode requires a mode naah arguahnt")}
		}

		targetMode := ahssage.Args[0]
		if err := ah.modeManager.SwitchToMode(targetMode); err != nil {
			return ErrorMessage{Err: err}
		}

		return ModeChangedMessage{
			OldMode: ah.modeManager.GetPreviousMode(),
			NewMode: targetMode,
		}
	}
}

// executeToggleHidden handles toggling hidden file visibility
func (ah *ActionHandler) executeToggleHidden(_ *config.MessageConfig) tea.Cmd {
	return func() tea.Msg {
		return UIMessage{Action: "toggle_hidden"}
	}
}

// executeRefresh handles refreshing the current directory
func (ah *ActionHandler) executeRefresh(_ *config.MessageConfig) tea.Cmd {
	return func() tea.Msg {
		return UIMessage{Action: "refresh"}
	}
}

// executeClearLog handles clearing the log
func (ah *ActionHandler) executeClearLog(_ *config.MessageConfig) tea.Cmd {
	return func() tea.Msg {
		return UIMessage{Action: "clear_log"}
	}
}

// executeToggleLog handles toggling log view visibility
func (ah *ActionHandler) executeToggleLog(_ *config.MessageConfig) tea.Cmd {
	return func() tea.Msg {
		return UIMessage{Action: "toggle_log"}
	}
}

// executeQuit handles application quit
func (ah *ActionHandler) executeQuit(_ *config.MessageConfig) tea.Cmd {
	return tea.Quit
}
