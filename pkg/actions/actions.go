package actions

import (
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dinhhuy258/fm/pkg/config"
)

// ActionHandlerFunc defines the signature for action handlers
type ActionHandlerFunc func(message *config.MessageConfig, originalKey tea.KeyMsg) tea.Cmd

// ActionHandler handles execution of config messages
type ActionHandler struct {
	actionMap map[string]ActionHandlerFunc
}

// NewActionHandler creates a new action handler
func NewActionHandler() *ActionHandler {
	ah := &ActionHandler{}
	ah.initActionMap()

	return ah
}

// initActionMap initializes the action handler map
func (ah *ActionHandler) initActionMap() {
	ah.actionMap = map[string]ActionHandlerFunc{
		// Core messages
		"SwitchMode": func(message *config.MessageConfig, _ tea.KeyMsg) tea.Cmd {
			return func() tea.Msg {
				return ModeChangedMessage{Mode: message.Args[0]}
			}
		},
		"Quit": func(_ *config.MessageConfig, _ tea.KeyMsg) tea.Cmd {
			return tea.Quit
		},
		"Null": func(_ *config.MessageConfig, _ tea.KeyMsg) tea.Cmd {
			return nil
		},

		// Navigation messages
		"ChangeDirectory": func(message *config.MessageConfig, _ tea.KeyMsg) tea.Cmd {
			return func() tea.Msg {
				return ChangeDirectoryMessage{Path: message.Args[0]}
			}
		},
		"FocusPath": func(message *config.MessageConfig, _ tea.KeyMsg) tea.Cmd {
			return func() tea.Msg {
				return FocusPathMessage{Path: message.Args[0]}
			}
		},
		"FocusByIndex": func(message *config.MessageConfig, _ tea.KeyMsg) tea.Cmd {
			return func() tea.Msg {
				index, _ := strconv.Atoi(message.Args[0])

				return FocusByIndexMessage{Index: index}
			}
		},
		"FocusNext": func(_ *config.MessageConfig, _ tea.KeyMsg) tea.Cmd {
			return func() tea.Msg {
				return NavigationMessage{Action: NavigationActionNext}
			}
		},
		"FocusPrevious": func(_ *config.MessageConfig, _ tea.KeyMsg) tea.Cmd {
			return func() tea.Msg {
				return NavigationMessage{Action: NavigationActionPrevious}
			}
		},
		"FocusFirst": func(_ *config.MessageConfig, _ tea.KeyMsg) tea.Cmd {
			return func() tea.Msg {
				return NavigationMessage{Action: NavigationActionFirst}
			}
		},
		"FocusLast": func(_ *config.MessageConfig, _ tea.KeyMsg) tea.Cmd {
			return func() tea.Msg {
				return NavigationMessage{Action: NavigationActionLast}
			}
		},
		"Enter": func(_ *config.MessageConfig, _ tea.KeyMsg) tea.Cmd {
			return func() tea.Msg {
				return NavigationMessage{Action: NavigationActionEnter}
			}
		},
		"Back": func(_ *config.MessageConfig, _ tea.KeyMsg) tea.Cmd {
			return func() tea.Msg {
				return NavigationMessage{Action: NavigationActionBack}
			}
		},

		// Selection messages
		"ToggleSelection": func(_ *config.MessageConfig, _ tea.KeyMsg) tea.Cmd {
			return func() tea.Msg {
				return SelectionMessage{Action: SelectionActionToggle}
			}
		},
		"ClearSelection": func(_ *config.MessageConfig, _ tea.KeyMsg) tea.Cmd {
			return func() tea.Msg {
				return SelectionMessage{Action: SelectionActionClear}
			}
		},
		"SelectAll": func(_ *config.MessageConfig, _ tea.KeyMsg) tea.Cmd {
			return func() tea.Msg {
				return SelectionMessage{Action: SelectionActionAll}
			}
		},
		"ToggleSelectionByPath": func(message *config.MessageConfig, _ tea.KeyMsg) tea.Cmd {
			return func() tea.Msg {
				return ToggleSelectionByPathMessage{Path: message.Args[0]}
			}
		},

		// Sorting messages
		"SortByName": func(_ *config.MessageConfig, _ tea.KeyMsg) tea.Cmd {
			return func() tea.Msg {
				return SortingMessage{SortType: SortTypeName}
			}
		},
		"SortBySize": func(_ *config.MessageConfig, _ tea.KeyMsg) tea.Cmd {
			return func() tea.Msg {
				return SortingMessage{SortType: SortTypeSize}
			}
		},
		"SortByDateModified": func(_ *config.MessageConfig, _ tea.KeyMsg) tea.Cmd {
			return func() tea.Msg {
				return SortingMessage{SortType: SortTypeDate}
			}
		},
		"SortByExtension": func(_ *config.MessageConfig, _ tea.KeyMsg) tea.Cmd {
			return func() tea.Msg {
				return SortingMessage{SortType: SortTypeExtension}
			}
		},
		"SortByDirFirst": func(_ *config.MessageConfig, _ tea.KeyMsg) tea.Cmd {
			return func() tea.Msg {
				return SortingMessage{SortType: SortTypeDirFirst}
			}
		},
		"ReverseSort": func(_ *config.MessageConfig, _ tea.KeyMsg) tea.Cmd {
			return func() tea.Msg {
				return SortingMessage{SortType: SortTypeReverse}
			}
		},

		// Bash execution
		"BashExec": func(message *config.MessageConfig, _ tea.KeyMsg) tea.Cmd {
			return func() tea.Msg {
				return BashExecMessage{Script: message.Args[0]}
			}
		},
		"BashExecSilently": func(message *config.MessageConfig, _ tea.KeyMsg) tea.Cmd {
			return func() tea.Msg {
				return BashExecSilentlyMessage{Script: message.Args[0]}
			}
		},

		// Input and logging
		"SetInputBuffer": func(message *config.MessageConfig, _ tea.KeyMsg) tea.Cmd {
			return func() tea.Msg {
				return SetInputBufferMessage{Value: message.Args[0]}
			}
		},
		"UpdateInputBufferFromKey": func(_ *config.MessageConfig, originalKey tea.KeyMsg) tea.Cmd {
			return func() tea.Msg {
				return UpdateInputBufferFromKeyMessage{Key: originalKey}
			}
		},
		"LogSuccess": func(message *config.MessageConfig, _ tea.KeyMsg) tea.Cmd {
			return func() tea.Msg {
				return LogMessage{Level: LogLevelSuccess, Message: message.Args[0]}
			}
		},
		"LogError": func(message *config.MessageConfig, _ tea.KeyMsg) tea.Cmd {
			return func() tea.Msg {
				return LogMessage{Level: LogLevelError, Message: message.Args[0]}
			}
		},
		"LogInfo": func(message *config.MessageConfig, _ tea.KeyMsg) tea.Cmd {
			return func() tea.Msg {
				return LogMessage{Level: LogLevelInfo, Message: message.Args[0]}
			}
		},
		"LogWarning": func(message *config.MessageConfig, _ tea.KeyMsg) tea.Cmd {
			return func() tea.Msg {
				return LogMessage{Level: LogLevelWarning, Message: message.Args[0]}
			}
		},

		// UI control
		"ToggleHidden": func(_ *config.MessageConfig, _ tea.KeyMsg) tea.Cmd {
			return func() tea.Msg {
				return UIMessage{Action: UIActionToggleHidden}
			}
		},
		"Refresh": func(_ *config.MessageConfig, _ tea.KeyMsg) tea.Cmd {
			return func() tea.Msg {
				return UIMessage{Action: UIActionRefresh}
			}
		},
	}
}

// ExecuteMessages executes a list of messages from config sequentially
func (ah *ActionHandler) ExecuteMessages(
	messages []*config.MessageConfig,
	originalKey tea.KeyMsg,
) tea.Cmd {
	var cmds []tea.Cmd

	for _, message := range messages {
		if cmd := ah.ExecuteMessage(message, originalKey); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	if len(cmds) == 0 {
		return nil
	}

	return tea.Sequence(cmds...)
}

// ExecuteMessage executes a single message
func (ah *ActionHandler) ExecuteMessage(
	message *config.MessageConfig,
	originalKey tea.KeyMsg,
) tea.Cmd {
	if action, exists := ah.actionMap[message.Name]; exists {
		return action(message, originalKey)
	}

	return func() tea.Msg {
		return LogMessage{
			Level:   LogLevelWarning,
			Message: fmt.Sprintf("Unknown message type: %s", message.Name),
		}
	}
}
