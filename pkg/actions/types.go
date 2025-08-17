package actions

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/dinhhuy258/fm/pkg/types"
)

// Message types for the executor

// ModeChangedMessage indicates a mode change occurred
type ModeChangedMessage struct {
	Mode string
}

// FocusPathMessage requests focusing on a specific path
type FocusPathMessage struct {
	Path string
}

// FocusByIndexMessage requests focusing by index
type FocusByIndexMessage struct {
	Index int
}

// LogLevel represents different log levels
type LogLevel string

const (
	LogLevelInfo    LogLevel = "info"
	LogLevelWarning LogLevel = "warning"
	LogLevelError   LogLevel = "error"
	LogLevelSuccess LogLevel = "success"
)

// String returns the string representation of the log level
func (l LogLevel) String() string {
	return string(l)
}

// LogMessage contains a log message
type LogMessage struct {
	Level   LogLevel
	Message string
}

// SetInputBufferMessage sets the input buffer value
type SetInputBufferMessage struct {
	Value string
}

// UpdateInputBufferFromKeyMessage updates input buffer from last key press
type UpdateInputBufferFromKeyMessage struct {
	Key tea.KeyMsg
}

// NavigationAction represents navigation actions.
type NavigationAction string

const (
	NavigationActionChangeDirectory NavigationAction = "change_directory"
	NavigationActionNext            NavigationAction = "next"
	NavigationActionPrevious        NavigationAction = "previous"
	NavigationActionFirst           NavigationAction = "first"
	NavigationActionLast            NavigationAction = "last"
	NavigationActionEnter           NavigationAction = "enter"
	NavigationActionBack            NavigationAction = "back"
)

// NavigationMessage handles navigation actions
type NavigationMessage struct {
	Action NavigationAction
	Path   string // Used with "change_directory" action
}

// SelectionAction represents selection actions.
type SelectionAction string

const (
	SelectionActionToggle SelectionAction = "toggle"
	SelectionActionClear  SelectionAction = "clear"
	SelectionActionAll    SelectionAction = "all"
)

// SelectionMessage handles selection actions
type SelectionMessage struct {
	Action SelectionAction
}

// UIAction represents UI control actions.
type UIAction string

const (
	UIActionToggleHidden UIAction = "toggle_hidden"
	UIActionRefresh      UIAction = "refresh"
)

// UIMessage handles UI control actions
type UIMessage struct {
	Action UIAction
}

// Sorting action constants that map to the unified sort types
const (
	SortTypeName      = string(types.SortTypeName)
	SortTypeSize      = string(types.SortTypeSize)
	SortTypeDate      = string(types.SortTypeDate)
	SortTypeExtension = string(types.SortTypeExtension)
	SortTypeDirFirst  = string(types.SortTypeDirFirst)
	SortTypeReverse   = "reverse" // Special action for toggling reverse
)

// SortingMessage handles sorting actions
type SortingMessage struct {
	SortType string // Can be a types.SortType value or "reverse"
}

// ToggleSelectionByPathMessage handles toggling selection by path
type ToggleSelectionByPathMessage struct {
	Path string
}

// BashExecMessage triggers bash command execution
type BashExecMessage struct {
	Script string
}

// BashExecSilentlyMessage triggers bash command execution silently
type BashExecSilentlyMessage struct {
	Script string
}

// ChangeDirectoryMessage handles directory change requests
type ChangeDirectoryMessage struct {
	Path string
}