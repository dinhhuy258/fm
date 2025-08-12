package actions

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dinhhuy258/fm/pkg/config"
)

// executeChangeDirectory handles directory changes
func (ah *ActionHandler) executeChangeDirectory(message *config.MessageConfig) tea.Cmd {
	return func() tea.Msg {
		if len(message.Args) == 0 {
			return ErrorMessage{Err: fmt.Errorf("ChangeDirectory requires a path argument")}
		}

		targetPath := message.Args[0]

		// Expand home directory and environment variables
		if strings.HasPrefix(targetPath, "~") {
			home, err := os.UserHomeDir()
			if err != nil {
				return ErrorMessage{Err: fmt.Errorf("failed to get home directory: %w", err)}
			}
			if targetPath == "~" {
				targetPath = home
			} else if strings.HasPrefix(targetPath, "~/") {
				targetPath = filepath.Join(home, targetPath[2:])
			}
		}

		// Expand environment variables
		targetPath = os.ExpandEnv(targetPath)

		// Check if directory exists
		if _, err := os.Stat(targetPath); os.IsNotExist(err) {
			return ErrorMessage{Err: fmt.Errorf("directory does not exist: %s", targetPath)}
		}

		// Return a message to load the directory
		return NavigationMessage{Action: "change_directory", Path: targetPath}
	}
}

// executeFocusPath handles focusing on a specific path
func (ah *ActionHandler) executeFocusPath(message *config.MessageConfig) tea.Cmd {
	return func() tea.Msg {
		if len(message.Args) == 0 {
			return ErrorMessage{Err: fmt.Errorf("FocusPath requires a path argument")}
		}

		targetPath := message.Args[0]

		return FocusPathMessage{Path: targetPath}
	}
}

// executeFocusByIndex handles focusing by index
func (ah *ActionHandler) executeFocusByIndex(message *config.MessageConfig) tea.Cmd {
	return func() tea.Msg {
		if len(message.Args) == 0 {
			return ErrorMessage{Err: fmt.Errorf("FocusByIndex requires an index argument")}
		}

		// Parse index (could be a number or expression)
		indexStr := message.Args[0]

		return FocusByIndexMessage{IndexExpression: indexStr}
	}
}

// executeFocusNext handles moving focus to next item
func (ah *ActionHandler) executeFocusNext(_ *config.MessageConfig) tea.Cmd {
	return func() tea.Msg {
		return NavigationMessage{Action: "next"}
	}
}

// executeFocusPrevious handles moving focus to previous item
func (ah *ActionHandler) executeFocusPrevious(_ *config.MessageConfig) tea.Cmd {
	return func() tea.Msg {
		return NavigationMessage{Action: "previous"}
	}
}

// executeFocusFirst handles moving focus to first item
func (ah *ActionHandler) executeFocusFirst(_ *config.MessageConfig) tea.Cmd {
	return func() tea.Msg {
		return NavigationMessage{Action: "first"}
	}
}

// executeFocusLast handles moving focus to last item
func (ah *ActionHandler) executeFocusLast(message *config.MessageConfig) tea.Cmd {
	return func() tea.Msg {
		return NavigationMessage{Action: "last"}
	}
}

// executeEnter handles entering directories or opening files
func (ah *ActionHandler) executeEnter(message *config.MessageConfig) tea.Cmd {
	return func() tea.Msg {
		return NavigationMessage{Action: "enter"}
	}
}

// executeBack handles going to parent directory
func (ah *ActionHandler) executeBack(message *config.MessageConfig) tea.Cmd {
	return func() tea.Msg {
		return NavigationMessage{Action: "back"}
	}
}