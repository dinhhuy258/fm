package actions

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
	IndexExpression string
}

// BashOutputMessage contains output from bash command execution
type BashOutputMessage struct {
	Output string
	Silent bool
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

// ErrorMessage contains an error
type ErrorMessage struct {
	Err error
}

// SetInputBufferMessage sets the input buffer value
type SetInputBufferMessage struct {
	Value string
}

// UpdateInputBufferFromKeyMessage updates input buffer from last key press
type UpdateInputBufferFromKeyMessage struct{}

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

// SortType represents sorting options.
type SortType string

const (
	SortTypeName      SortType = "name"
	SortTypeSize      SortType = "size"
	SortTypeDate      SortType = "date"
	SortTypeExtension SortType = "extension"
	SortTypeDirFirst  SortType = "dir_first"
	SortTypeReverse   SortType = "reverse"
)

// SortingMessage handles sorting actions
type SortingMessage struct {
	SortType SortType
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
