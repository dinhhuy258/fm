package actions

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dinhhuy258/fm/pkg/config"
)

// executeToggleSelection handles toggling selection of current item
func (ah *ActionHandler) executeToggleSelection(_ *config.MessageConfig) tea.Cmd {
	return func() tea.Msg {
		return SelectionMessage{Action: "toggle"}
	}
}

// executeClearSelection handles clearing all selections
func (ah *ActionHandler) executeClearSelection(_ *config.MessageConfig) tea.Cmd {
	return func() tea.Msg {
		return SelectionMessage{Action: "clear"}
	}
}

// executeSelectAll handles selecting all items
func (ah *ActionHandler) executeSelectAll(_ *config.MessageConfig) tea.Cmd {
	return func() tea.Msg {
		return SelectionMessage{Action: "all"}
	}
}

// executeToggleSelectionByPath handles toggling selection by path
func (ah *ActionHandler) executeToggleSelectionByPath(message *config.MessageConfig) tea.Cmd {
	return func() tea.Msg {
		if len(message.Args) == 0 {
			return ErrorMessage{Err: fmt.Errorf("ToggleSelectionByPath requires a path argument")}
		}

		path := message.Args[0]

		return ToggleSelectionByPathMessage{Path: path}
	}
}
