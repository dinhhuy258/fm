package tui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dinhhuy258/fm/pkg/actions"
	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/types"
)

// PipeMessage represents a message from pipe
type PipeMessage struct {
	Command string
}

// errorMessage contains an error
type errorMessage struct {
	Message string
}

// directoryLoadedMessage indicates that a directory has been loaded
type directoryLoadedMessage struct {
	path    string
	entries []fs.IEntry
}

// handleMessage processes incoming messages and returns updated model with commands
func (m Model) handleMessage(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleWindowSize(msg)
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)
	default:
		return m.handleOtherMessage(msg)
	}
}

// handleWindowSize handles window resize events
func (m Model) handleWindowSize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	headerHeight := 3
	footerHeight := 1
	interactiveHeight := 1
	availableExplorerHeight := msg.Height - headerHeight - footerHeight - interactiveHeight

	m.helpModel.SetSize(msg.Width, msg.Height)
	m.inputModel.SetSize(msg.Width, 1)
	m.notificationModel.SetSize(msg.Width, 1)
	m.explorerModel.SetSize(msg.Width, availableExplorerHeight)

	return m, nil
}

// handleKeyMsg handles keyboard input events
func (m Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.helpModel.IsVisible() {
		m.helpModel.Update(msg)

		return m, nil
	}

	if msg.String() == "?" {
		m.helpModel.Show()

		return m, nil
	}

	return m.handleKeyMap(msg)
}

// handleOtherMessage handles all non-keyboard, non-window-size messages
func (m Model) handleOtherMessage(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case errorMessage:
		m.notificationModel.ShowNotification(NotificationError, msg.Message)

		return m, nil
	case directoryLoadedMessage:
		m.currentPath = msg.path
		m.explorerModel.SetEntries(msg.entries)

		return m, nil
	case PipeMessage:
		return m.handlePipeMessage(msg.Command)
	case actions.ModeChangedMessage:
		err := m.modeManager.SwitchToMode(msg.Mode)
		if err != nil {
			return m, func() tea.Msg {
				return errorMessage{
					Message: fmt.Sprintf("Failed to switch mode: %v", err),
				}
			}
		}

		m.inputModel.Hide()
		m.notificationModel.Show()

		return m, nil
	case AutoClearMessage:
		m.notificationModel.ClearNotification()
	case actions.LogMessage:
		switch msg.Level {
		case actions.LogLevelError:
			return m, m.notificationModel.ShowNotification(NotificationError, msg.Message)
		case actions.LogLevelWarning:
			return m, m.notificationModel.ShowNotification(NotificationWarning, msg.Message)
		case actions.LogLevelSuccess:
			return m, m.notificationModel.ShowNotification(NotificationSuccess, msg.Message)
		default:
			return m, m.notificationModel.ShowNotification(NotificationInfo, msg.Message)
		}
	case actions.SetInputBufferMessage:
		m.notificationModel.Hide()
		m.inputModel.Show(msg.Value)

		return m, nil
	case actions.UpdateInputBufferFromKeyMessage:
		return m, m.inputModel.Update(msg.Key)
	case actions.FocusPathMessage:
		dir := filepath.Dir(msg.Path)
		if dir == m.currentPath {
			m.explorerModel.FocusPath(msg.Path)

			return m, nil
		}

		return m, tea.Sequence(
			m.loadDirectoryCmd(dir),
			func() tea.Msg {
				return actions.FocusPathMessage{Path: msg.Path}
			},
		)
	case actions.NavigationMessage:
		return m.handleNavigationMessage(msg)
	case actions.FocusByIndexMessage:
		return m.handleFocusByIndexMessage(msg)
	case actions.SelectionMessage:
		return m.handleSelectionMessage(msg)
	case actions.ToggleSelectionByPathMessage:
		return m.handleToggleSelectionByPathMessage(msg)
	case actions.UIMessage:
		return m.handleUIMessage(msg)
	case actions.SortingMessage:
		return m.handleSortingMessage(msg)
	case actions.BashExecMessage:
		return m.handleBashExecMessage(msg)
	case actions.BashExecSilentlyMessage:
		return m.handleBashExecSilentlyMessage(msg)
	case actions.ChangeDirectoryMessage:
		return m.handleChangeDirectoryMessage(msg)
	}

	return m, nil
}

// handleKeyMap handles key presses and resolves them to actions
func (m Model) handleKeyMap(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	action := m.keyManager.ResolveKeyAction(msg)
	if action != nil {
		return m, m.actionHandler.ExecuteMessages(action.Messages, msg)
	}

	return m, m.notificationModel.ShowNotification(NotificationWarning,
		fmt.Sprintf("No action found for key: %s", msg.String()),
	)
}

// handlePipeMessage processes messages received from the pipe
func (m Model) handlePipeMessage(command string) (tea.Model, tea.Cmd) {
	// Parse the pipe message - format is usually: CommandName arg1 arg2 ...
	commandName, args := parseCommand(command)
	if commandName == "" {
		return m, nil
	}

	message := &config.MessageConfig{
		Name: commandName,
		Args: args,
	}
	cmd := m.actionHandler.ExecuteMessage(message, tea.KeyMsg{})
	if cmd != nil {
		return m, cmd
	}

	return m, nil
}

// handleNavigationMessage processes navigation actions
func (m Model) handleNavigationMessage(msg actions.NavigationMessage) (tea.Model, tea.Cmd) {
	switch msg.Action {
	case actions.NavigationActionNext:
		m.explorerModel.Move(1)
	case actions.NavigationActionPrevious:
		m.explorerModel.Move(-1)
	case actions.NavigationActionFirst:
		m.explorerModel.MoveFirst()
	case actions.NavigationActionLast:
		m.explorerModel.MoveLast()
	case actions.NavigationActionEnter:
		if entry := m.explorerModel.GetFocusedEntry(); entry != nil {
			if entry.IsDirectory() {
				return m, m.loadDirectoryCmd(entry.GetPath())
			}
		}
	case actions.NavigationActionBack:
		if m.currentPath != "" && m.currentPath != "/" {
			parentPath := filepath.Dir(m.currentPath)
			currentPath := m.currentPath

			return m, tea.Sequence(
				m.loadDirectoryCmd(parentPath),
				func() tea.Msg {
					return actions.FocusPathMessage{Path: currentPath}
				},
			)
		}
	case actions.NavigationActionChangeDirectory:
		if msg.Path != "" {
			return m, m.loadDirectoryCmd(msg.Path)
		}
	}

	return m, nil
}

// handleSelectionMessage processes selection actions
func (m Model) handleSelectionMessage(msg actions.SelectionMessage) (tea.Model, tea.Cmd) {
	switch msg.Action {
	case actions.SelectionActionToggle:
		// Toggle selection of current item using direct method
		m.explorerModel.ToggleSelection()
	case actions.SelectionActionClear:
		// Clear all selections
		m.explorerModel.ClearSelections()
	case actions.SelectionActionAll:
		// Select all items
		m.explorerModel.SelectAll()
	}

	// Selection count is automatically updated via GetStats() in renderHeader()

	return m, nil
}

// handleUIMessage processes UI control actions
func (m Model) handleUIMessage(msg actions.UIMessage) (tea.Model, tea.Cmd) {
	switch msg.Action {
	case actions.UIActionToggleHidden:
		// Toggle hidden file visibility
		m.showHidden = !m.showHidden
		// Toggle hidden files silently

		return m, m.loadDirectoryCmd(m.currentPath) // Reload with new settings
	case actions.UIActionRefresh:
		// Refresh current directory
		return m, m.loadDirectoryCmd(m.currentPath)
	}

	return m, nil
}

// handleSortingMessage processes sorting actions
func (m Model) handleSortingMessage(msg actions.SortingMessage) (tea.Model, tea.Cmd) {
	switch msg.SortType {
	case actions.SortTypeReverse:
		// Toggle reverse sorting
		m.reverse = !m.reverse
	default:
		// Direct assignment since action types now match fs types
		m.sortType = types.SortType(msg.SortType)
	}

	// Reload directory with new sorting
	return m, m.loadDirectoryCmd(m.currentPath)
}

// handleFocusByIndexMessage processes focus by index actions
func (m Model) handleFocusByIndexMessage(msg actions.FocusByIndexMessage) (tea.Model, tea.Cmd) {
	// Parse index from expression
	var index int
	var err error

	indexStr := msg.IndexExpression
	if indexStr == "" {
		// Invalid index expression - handle silently
		return m, nil
	}

	// Handle simple numeric index
	_, err = fmt.Sscanf(indexStr, "%d", &index)
	if err != nil {
		// Invalid index format - handle silently
		return m, nil
	}

	// Update focus through explorer model (silently)
	m.explorerModel.SetFocusByIndex(index)

	return m, nil
}

// handleToggleSelectionByPathMessage processes toggle selection by path actions
func (m Model) handleToggleSelectionByPathMessage(
	msg actions.ToggleSelectionByPathMessage,
) (tea.Model, tea.Cmd) {
	path := msg.Path

	// Toggle selection in explorer model (silently)
	m.explorerModel.ToggleSelectionByPath(path)

	// Selection count is automatically updated via GetStats() in renderHeader()

	return m, nil
}

// handleBashExecMessage executes bash commands
func (m Model) handleBashExecMessage(msg actions.BashExecMessage) (tea.Model, tea.Cmd) {
	return m.handleBashExecution(msg.Script, false)
}

// handleBashExecSilentlyMessage executes bash commands silently
func (m Model) handleBashExecSilentlyMessage(
	msg actions.BashExecSilentlyMessage,
) (tea.Model, tea.Cmd) {
	return m.handleBashExecution(msg.Script, true)
}

// handleChangeDirectoryMessage processes directory change requests with validation
func (m Model) handleChangeDirectoryMessage(
	msg actions.ChangeDirectoryMessage,
) (tea.Model, tea.Cmd) {
	targetPath := msg.Path

	// Handle home directory expansion
	if strings.HasPrefix(targetPath, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return m, func() tea.Msg {
				return errorMessage{
					Message: fmt.Sprintf("Failed to get home directory: %v", err),
				}
			}
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
		return m, func() tea.Msg {
			return errorMessage{
				Message: fmt.Sprintf("Directory does not exist: %s", targetPath),
			}
		}
	}

	// If all validation passes, change to the directory
	return m, m.loadDirectoryCmd(targetPath)
}

// handleBashExecution processes bash execution with environment setup
func (m Model) handleBashExecution(script string, silent bool) (tea.Model, tea.Cmd) {
	selectedEntries := m.explorerModel.GetSelectedEntries()
	focusIndex := m.explorerModel.GetFocus()

	focusPath := m.currentPath
	if focusedEntry := m.explorerModel.GetFocusedEntry(); focusedEntry != nil {
		focusPath = focusedEntry.GetPath()
	}

	selections := make([]string, len(selectedEntries))
	for i, entry := range selectedEntries {
		selections[i] = entry.GetPath()
	}

	if len(selections) > 0 {
		selectionPath := m.pipe.GetSelectionPath()
		if err := writeSelectionsToFile(selectionPath, selections); err != nil {
			return m, func() tea.Msg {
				return errorMessage{
					Message: fmt.Sprintf("Failed to write selections to pipe: %v", err),
				}
			}
		}
	}

	env := os.Environ()
	inputBuffer := m.inputModel.GetValue()
	env = append(env, fmt.Sprintf("FM_FOCUS_PATH=%s", focusPath))
	env = append(env, fmt.Sprintf("FM_PWD=%s", m.currentPath))
	env = append(env, fmt.Sprintf("FM_FOCUS_IDX=%d", focusIndex))
	env = append(env, fmt.Sprintf("FM_INPUT_BUFFER=%s", inputBuffer))
	env = append(env, fmt.Sprintf("FM_PIPE_MSG_IN=%s", m.pipe.GetMessageInPath()))
	env = append(env, fmt.Sprintf("FM_PIPE_SELECTION=%s", m.pipe.GetSelectionPath()))
	env = append(env, fmt.Sprintf("FM_SESSION_PATH=%s", m.pipe.GetSessionPath()))

	cmd := exec.Command("bash", "-c", script)
	cmd.Env = env
	cmd.Dir = m.currentPath

	if silent {
		if err := cmd.Start(); err != nil {
			return m, func() tea.Msg {
				return errorMessage{
					Message: fmt.Sprintf("Failed to start bash command: %v", err),
				}
			}
		}

		go func() { _ = cmd.Wait() }()

		return m, nil
	}

	return m, tea.ExecProcess(&exec.Cmd{
		Path: "/bin/bash",
		Args: []string{"bash", "-c", "clear && " + script},
		Env:  env,
		Dir:  m.currentPath,
	}, func(err error) tea.Msg {
		if err != nil {
			return errorMessage{Message: fmt.Sprintf("Failed to execute bash %v", err)}
		}

		return nil
	})
}

// writeSelectionsToFile writes selected file paths to the selection pipe file
func writeSelectionsToFile(path string, selections []string) error {
	const perm = 0600

	if len(selections) == 0 {
		return os.WriteFile(path, []byte(""), perm)
	}

	content := strings.Join(selections, "\n") + "\n"

	return os.WriteFile(path, []byte(content), perm)
}

// loadDirectoryCmd loads directory contents
func (m Model) loadDirectoryCmd(path string) tea.Cmd {
	return func() tea.Msg {
		entries, err := loadDirectory(path, m.showHidden, m.sortType, m.reverse)
		if err != nil {
			return errorMessage{
				Message: fmt.Sprintf("Failed to load directory %s: %v", path, err),
			}
		}

		return directoryLoadedMessage{path: path, entries: entries}
	}
}

// parseCommand parses a shell command line, handling quoted strings properly
func parseCommand(content string) (string, []string) {
	content = strings.TrimSpace(content)
	if content == "" {
		return "", nil
	}

	var tokens []string
	var current strings.Builder
	inSingleQuote := false
	inDoubleQuote := false

	for _, r := range content {
		switch r {
		case '\'':
			if !inDoubleQuote {
				inSingleQuote = !inSingleQuote
				// Don't include the quotes in the token
			} else {
				current.WriteRune(r)
			}
		case '"':
			if !inSingleQuote {
				inDoubleQuote = !inDoubleQuote
				// Don't include the quotes in the token
			} else {
				current.WriteRune(r)
			}
		case ' ', '\t', '\n':
			if inSingleQuote || inDoubleQuote {
				current.WriteRune(r)
			} else {
				// End of token
				if current.Len() > 0 {
					tokens = append(tokens, current.String())
					current.Reset()
				}
			}
		default:
			current.WriteRune(r)
		}
	}

	// Add the last token if any
	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	if len(tokens) == 0 {
		return "", nil
	}

	return tokens[0], tokens[1:]
}