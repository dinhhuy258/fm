package actions

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/fm/pkg/pipe"
)

// ActionHandlerFunc defines the signature for action handlers
type ActionHandlerFunc func(message *config.MessageConfig, currentPath string, inputBuffer string) tea.Cmd

// ActionHandler handles execution of config messages
type ActionHandler struct {
	pipe      *pipe.Pipe
	actionMap map[string]ActionHandlerFunc
}

// wrapSimple creates a wrapper function for handlers that only need the message
func (ah *ActionHandler) wrapSimple(handler func(*config.MessageConfig) tea.Cmd) ActionHandlerFunc {
	return func(message *config.MessageConfig, currentPath string, inputBuffer string) tea.Cmd {
		return handler(message)
	}
}

// wrapBashExec creates a wrapper function for bash execution with specific silent flag
func (ah *ActionHandler) wrapBashExec(silent bool) ActionHandlerFunc {
	return func(message *config.MessageConfig, currentPath string, inputBuffer string) tea.Cmd {
		return ah.executeBashExec(message, silent)
	}
}

// NewActionHandler creates a new action handler
func NewActionHandler(pipe *pipe.Pipe) *ActionHandler {
	ah := &ActionHandler{
		pipe: pipe,
	}
	ah.initActionMap()

	return ah
}

// initActionMap initializes the action handler map
func (ah *ActionHandler) initActionMap() {
	ah.actionMap = map[string]ActionHandlerFunc{
		// Core messages
		"SwitchMode": ah.wrapSimple(ah.executeSwitchMode),
		"Quit":       ah.wrapSimple(ah.executeQuit),
		"Null":       ah.executeNull,

		// Navigation messages
		"ChangeDirectory": ah.wrapSimple(ah.executeChangeDirectory),
		"FocusPath":       ah.wrapSimple(ah.executeFocusPath),
		"FocusByIndex":    ah.wrapSimple(ah.executeFocusByIndex),
		"FocusNext":       ah.wrapSimple(ah.executeFocusNext),
		"FocusPrevious":   ah.wrapSimple(ah.executeFocusPrevious),
		"FocusFirst":      ah.wrapSimple(ah.executeFocusFirst),
		"FocusLast":       ah.wrapSimple(ah.executeFocusLast),
		"Enter":           ah.wrapSimple(ah.executeEnter),
		"Back":            ah.wrapSimple(ah.executeBack),

		// Selection messages
		"ToggleSelection":       ah.wrapSimple(ah.executeToggleSelection),
		"ClearSelection":        ah.wrapSimple(ah.executeClearSelection),
		"SelectAll":             ah.wrapSimple(ah.executeSelectAll),
		"ToggleSelectionByPath": ah.wrapSimple(ah.executeToggleSelectionByPath),

		// Sorting messages
		"SortByName":         ah.wrapSimple(ah.executeSortByName),
		"SortBySize":         ah.wrapSimple(ah.executeSortBySize),
		"SortByDateModified": ah.wrapSimple(ah.executeSortByDateModified),
		"SortByExtension":    ah.wrapSimple(ah.executeSortByExtension),
		"SortByDirFirst":     ah.wrapSimple(ah.executeSortByDirFirst),
		"ReverseSort":        ah.wrapSimple(ah.executeReverseSort),

		// Bash execution
		"BashExec":         ah.wrapBashExec(false),
		"BashExecSilently": ah.wrapBashExec(true),

		// Input and logging
		"SetInputBuffer":           ah.wrapSimple(ah.executeSetInputBuffer),
		"UpdateInputBufferFromKey": ah.wrapSimple(ah.executeUpdateInputBufferFromKey),
		"LogSuccess":               ah.wrapSimple(ah.executeLogSuccess),
		"LogError":                 ah.wrapSimple(ah.executeLogError),
		"LogInfo":                  ah.wrapSimple(ah.executeLogInfo),
		"LogWarning":               ah.wrapSimple(ah.executeLogWarning),

		// UI control
		"ToggleHidden": ah.wrapSimple(ah.executeToggleHidden),
		"Refresh":      ah.wrapSimple(ah.executeRefresh),
		"ClearLog":     ah.wrapSimple(ah.executeClearLog),
		"ToggleLog":    ah.wrapSimple(ah.executeToggleLog),
	}
}

// executeNull handles null operations
func (ah *ActionHandler) executeNull(message *config.MessageConfig, currentPath string, inputBuffer string) tea.Cmd {
	// Null operation - do nothing
	return nil
}

// ExecuteMessages executes a list of messages from config
func (ah *ActionHandler) ExecuteMessages(
	messages []*config.MessageConfig,
	currentPath string,
	inputBuffer string,
) []tea.Cmd {
	var cmds []tea.Cmd

	for _, message := range messages {
		if cmd := ah.ExecuteMessage(message, currentPath, inputBuffer); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return cmds
}

// ExecuteMessage executes a single message (public method for pipe messages)
func (ah *ActionHandler) ExecuteMessage(
	message *config.MessageConfig,
	currentPath string,
	inputBuffer string,
) tea.Cmd {
	if handler, exists := ah.actionMap[message.Name]; exists {
		return handler(message, currentPath, inputBuffer)
	}

	return func() tea.Msg {
		return LogMessage{
			Level:   LogLevelWarning,
			Message: fmt.Sprintf("Unknown message type: %s", message.Name),
		}
	}
}
