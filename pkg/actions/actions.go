package actions

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dinhhuy258/fm/pkg/config"
)

// ActionHandlerFunc defines the signature for action handlers
type ActionHandlerFunc func(message *config.MessageConfig, currentPath string, inputBuffer string) tea.Cmd

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
		"SwitchMode": func(message *config.MessageConfig, _ string, _ string) tea.Cmd {
			return func() tea.Msg {
				return ModeChangedMessage{Mode: message.Args[0]}
			}
		},
		"Quit": func(_ *config.MessageConfig, _ string, _ string) tea.Cmd {
			return tea.Quit
		},
		"Null": func(_ *config.MessageConfig, _ string, _ string) tea.Cmd {
			return nil
		},

		// Navigation messages
		"ChangeDirectory": func(message *config.MessageConfig, _ string, _ string) tea.Cmd {
			return func() tea.Msg {
				if len(message.Args) == 0 {
					return ErrorMessage{Err: fmt.Errorf("ChangeDirectory requires a path argument")}
				}

				targetPath := message.Args[0]

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

				targetPath = os.ExpandEnv(targetPath)

				if _, err := os.Stat(targetPath); os.IsNotExist(err) {
					return ErrorMessage{Err: fmt.Errorf("directory does not exist: %s", targetPath)}
				}

				return NavigationMessage{Action: NavigationActionChangeDirectory, Path: targetPath}
			}
		},
		"FocusPath": func(message *config.MessageConfig, _ string, _ string) tea.Cmd {
			return func() tea.Msg {
				if len(message.Args) == 0 {
					return ErrorMessage{Err: fmt.Errorf("FocusPath requires a path argument")}
				}

				return FocusPathMessage{Path: message.Args[0]}
			}
		},
		"FocusByIndex": func(message *config.MessageConfig, _ string, _ string) tea.Cmd {
			return func() tea.Msg {
				if len(message.Args) == 0 {
					return ErrorMessage{Err: fmt.Errorf("FocusByIndex requires an index argument")}
				}

				return FocusByIndexMessage{IndexExpression: message.Args[0]}
			}
		},
		"FocusNext": func(_ *config.MessageConfig, _ string, _ string) tea.Cmd {
			return func() tea.Msg {
				return NavigationMessage{Action: NavigationActionNext}
			}
		},
		"FocusPrevious": func(_ *config.MessageConfig, _ string, _ string) tea.Cmd {
			return func() tea.Msg {
				return NavigationMessage{Action: NavigationActionPrevious}
			}
		},
		"FocusFirst": func(_ *config.MessageConfig, _ string, _ string) tea.Cmd {
			return func() tea.Msg {
				return NavigationMessage{Action: NavigationActionFirst}
			}
		},
		"FocusLast": func(_ *config.MessageConfig, _ string, _ string) tea.Cmd {
			return func() tea.Msg {
				return NavigationMessage{Action: NavigationActionLast}
			}
		},
		"Enter": func(_ *config.MessageConfig, _ string, _ string) tea.Cmd {
			return func() tea.Msg {
				return NavigationMessage{Action: NavigationActionEnter}
			}
		},
		"Back": func(_ *config.MessageConfig, _ string, _ string) tea.Cmd {
			return func() tea.Msg {
				return NavigationMessage{Action: NavigationActionBack}
			}
		},

		// Selection messages
		"ToggleSelection": func(_ *config.MessageConfig, _ string, _ string) tea.Cmd {
			return func() tea.Msg {
				return SelectionMessage{Action: SelectionActionToggle}
			}
		},
		"ClearSelection": func(_ *config.MessageConfig, _ string, _ string) tea.Cmd {
			return func() tea.Msg {
				return SelectionMessage{Action: SelectionActionClear}
			}
		},
		"SelectAll": func(_ *config.MessageConfig, _ string, _ string) tea.Cmd {
			return func() tea.Msg {
				return SelectionMessage{Action: SelectionActionAll}
			}
		},
		"ToggleSelectionByPath": func(message *config.MessageConfig, _ string, _ string) tea.Cmd {
			return func() tea.Msg {
				if len(message.Args) == 0 {
					return ErrorMessage{Err: fmt.Errorf("ToggleSelectionByPath requires a path argument")}
				}

				return ToggleSelectionByPathMessage{Path: message.Args[0]}
			}
		},

		// Sorting messages
		"SortByName": func(_ *config.MessageConfig, _ string, _ string) tea.Cmd {
			return func() tea.Msg {
				return SortingMessage{SortType: SortTypeName}
			}
		},
		"SortBySize": func(_ *config.MessageConfig, _ string, _ string) tea.Cmd {
			return func() tea.Msg {
				return SortingMessage{SortType: SortTypeSize}
			}
		},
		"SortByDateModified": func(_ *config.MessageConfig, _ string, _ string) tea.Cmd {
			return func() tea.Msg {
				return SortingMessage{SortType: SortTypeDate}
			}
		},
		"SortByExtension": func(_ *config.MessageConfig, _ string, _ string) tea.Cmd {
			return func() tea.Msg {
				return SortingMessage{SortType: SortTypeExtension}
			}
		},
		"SortByDirFirst": func(_ *config.MessageConfig, _ string, _ string) tea.Cmd {
			return func() tea.Msg {
				return SortingMessage{SortType: SortTypeDirFirst}
			}
		},
		"ReverseSort": func(_ *config.MessageConfig, _ string, _ string) tea.Cmd {
			return func() tea.Msg {
				return SortingMessage{SortType: SortTypeReverse}
			}
		},

		// Bash execution
		"BashExec": func(message *config.MessageConfig, _ string, _ string) tea.Cmd {
			return func() tea.Msg {
				return BashExecMessage{Script: message.Args[0]}
			}
		},
		"BashExecSilently": func(message *config.MessageConfig, _ string, _ string) tea.Cmd {
			return func() tea.Msg {
				return BashExecSilentlyMessage{Script: message.Args[0]}
			}
		},

		// Input and logging
		"SetInputBuffer": func(message *config.MessageConfig, _ string, _ string) tea.Cmd {
			return func() tea.Msg {
				value := ""
				if len(message.Args) > 0 {
					value = message.Args[0]
				}

				return SetInputBufferMessage{Value: value}
			}
		},
		"UpdateInputBufferFromKey": func(_ *config.MessageConfig, _ string, _ string) tea.Cmd {
			return func() tea.Msg {
				return UpdateInputBufferFromKeyMessage{}
			}
		},
		"LogSuccess": func(message *config.MessageConfig, _ string, _ string) tea.Cmd {
			return func() tea.Msg {
				return LogMessage{Level: LogLevelSuccess, Message: message.Args[0]}
			}
		},
		"LogError": func(message *config.MessageConfig, _ string, _ string) tea.Cmd {
			return func() tea.Msg {
				return LogMessage{Level: LogLevelError, Message: message.Args[0]}
			}
		},
		"LogInfo": func(message *config.MessageConfig, _ string, _ string) tea.Cmd {
			return func() tea.Msg {
				return LogMessage{Level: LogLevelInfo, Message: message.Args[0]}
			}
		},
		"LogWarning": func(message *config.MessageConfig, _ string, _ string) tea.Cmd {
			return func() tea.Msg {
				return LogMessage{Level: LogLevelWarning, Message: message.Args[0]}
			}
		},

		// UI control
		"ToggleHidden": func(_ *config.MessageConfig, _ string, _ string) tea.Cmd {
			return func() tea.Msg {
				return UIMessage{Action: UIActionToggleHidden}
			}
		},
		"Refresh": func(_ *config.MessageConfig, _ string, _ string) tea.Cmd {
			return func() tea.Msg {
				return UIMessage{Action: UIActionRefresh}
			}
		},
	}
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

// ExecuteMessage executes a single message
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
