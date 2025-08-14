package actions

// Message types for the executor

// ModeChangedMessage indicates a mode change occurred
type ModeChangedMessage struct {
	OldMode string
	NewMode string
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

// LogMessage contains a log message
type LogMessage struct {
	Level   string // "info", "warning", "error", "success"
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

// NavigationMessage handles navigation actions
type NavigationMessage struct {
	Action string // "next", "previous", "first", "last", "enter", "back", "change_directory"
	Path   string // Used with "change_directory" action
}

// SelectionMessage handles selection actions
type SelectionMessage struct {
	Action string // "toggle", "clear", "all"
}

// UIMessage handles UI control actions
type UIMessage struct {
	Action string // "toggle_hidden", "refresh", "clear_log", "toggle_log"
}

// SortingMessage handles sorting actions
type SortingMessage struct {
	SortType string // "name", "size", "date", "extension", "dir_first", "reverse"
}

// ToggleSelectionByPathMessage handles toggling selection by path
type ToggleSelectionByPathMessage struct {
	Path string
}

// WriteSelectionsMessage requests writing selections to pipe before bash execution
type WriteSelectionsMessage struct {
	Script        string
	CurrentPath   string
	InputBuffer   string
	Silent        bool
	SelectionPath string
}

// InteractiveBashMessage handles interactive bash execution with TUI suspension
type InteractiveBashMessage struct {
	Script      string
	Environment []string
	WorkingDir  string
}