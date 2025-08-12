package actions

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/fm/pkg/pipe"
)

// ActionHandler handles execution of config messages
type ActionHandler struct {
	modeManager ModeManager
	pipe        *pipe.Pipe
}

// ModeManager interface for mode management operations
type ModeManager interface {
	SwitchToMode(mode string) error
	GetPreviousMode() string
}

// NewActionHandler creates a new action handler
func NewActionHandler(modeManager ModeManager, pipe *pipe.Pipe) *ActionHandler {
	return &ActionHandler{
		modeManager: modeManager,
		pipe:        pipe,
	}
}

// ExecuteMessages executes a list of messages from config
func (ah *ActionHandler) ExecuteMessages(messages []*config.MessageConfig, currentPath string, inputBuffer string) []tea.Cmd {
	var cmds []tea.Cmd

	for _, message := range messages {
		if cmd := ah.executeMessage(message, currentPath, inputBuffer); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return cmds
}

// ExecuteMessage executes a single message (public method for pipe messages)
func (ah *ActionHandler) ExecuteMessage(message *config.MessageConfig, currentPath string, inputBuffer string) tea.Cmd {
	return ah.executeMessage(message, currentPath, inputBuffer)
}

// executeMessage executes a single message and returns a tea.Cmd if needed
func (ah *ActionHandler) executeMessage(message *config.MessageConfig, currentPath string, inputBuffer string) tea.Cmd {
	switch message.Name {
	// Core messages
	case "SwitchMode":
		return ah.executeSwitchMode(message)
	case "Quit":
		return ah.executeQuit(message)
	case "Null":
		// Null operation - do nothing
		return nil

	// Navigation messages
	case "ChangeDirectory":
		return ah.executeChangeDirectory(message)
	case "FocusPath":
		return ah.executeFocusPath(message)
	case "FocusByIndex":
		return ah.executeFocusByIndex(message)
	case "FocusNext":
		return ah.executeFocusNext(message)
	case "FocusPrevious":
		return ah.executeFocusPrevious(message)
	case "FocusFirst":
		return ah.executeFocusFirst(message)
	case "FocusLast":
		return ah.executeFocusLast(message)
	case "Enter":
		return ah.executeEnter(message)
	case "Back":
		return ah.executeBack(message)

	// Selection messages
	case "ToggleSelection":
		return ah.executeToggleSelection(message)
	case "ClearSelection":
		return ah.executeClearSelection(message)
	case "SelectAll":
		return ah.executeSelectAll(message)
	case "ToggleSelectionByPath":
		return ah.executeToggleSelectionByPath(message)

	// Sorting messages
	case "SortByName":
		return ah.executeSortByName(message)
	case "SortBySize":
		return ah.executeSortBySize(message)
	case "SortByDateModified":
		return ah.executeSortByDateModified(message)
	case "SortByExtension":
		return ah.executeSortByExtension(message)
	case "SortByDirFirst":
		return ah.executeSortByDirFirst(message)
	case "ReverseSort":
		return ah.executeReverseSort(message)

	// Bash execution
	case "BashExec":
		return ah.executeBashExec(message, currentPath, inputBuffer, false)
	case "BashExecSilently":
		return ah.executeBashExec(message, currentPath, inputBuffer, true)

	// Input and logging
	case "SetInputBuffer":
		return ah.executeSetInputBuffer(message)
	case "UpdateInputBufferFromKey":
		return ah.executeUpdateInputBufferFromKey(message)
	case "LogSuccess":
		return ah.executeLog(message, "success")
	case "LogError":
		return ah.executeLog(message, "error")
	case "LogInfo":
		return ah.executeLog(message, "info")
	case "LogWarning":
		return ah.executeLog(message, "warning")

	// UI control
	case "ToggleHidden":
		return ah.executeToggleHidden(message)
	case "Refresh":
		return ah.executeRefresh(message)
	case "ClearLog":
		return ah.executeClearLog(message)
	case "ToggleLog":
		return ah.executeToggleLog(message)

	default:
		// Log unknown message types
		return func() tea.Msg {
			return LogMessage{
				Level:   "warning",
				Message: fmt.Sprintf("Unknown message type: %s", message.Name),
			}
		}
	}
}