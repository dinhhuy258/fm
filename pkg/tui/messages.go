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
	// Account for border: 4 characters width (2 border + 2 padding), 4 characters height
	borderPadding := 4
	availableWidth := max(msg.Width-borderPadding, 1)
	availableHeight := max(msg.Height-borderPadding, 1)

	headerHeight := 3
	footerHeight := 1
	interactiveHeight := 1
	availableExplorerHeight := availableHeight - headerHeight - footerHeight - interactiveHeight

	m.helpModel.SetSize(msg.Width, msg.Height)
	m.inputModel.SetSize(availableWidth, 1)
	m.notificationModel.SetSize(availableWidth, 1)
	m.explorerModel.SetSize(availableWidth, availableExplorerHeight)

	return m, nil
}

// handleKeyMsg handles keyboard input events
func (m Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.helpModel.IsVisible() {
		m.helpModel.Update(msg)

		return m, nil
	}

	if msg.String() == HelpToggleKey {
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
		m.modeManager.SwitchToMode(msg.Mode)
		// Notification is always shown by default
		m.notificationModel.Show()
		m.inputModel.Hide()

		return m, nil
	case AutoClearMessage:
		m.notificationModel.ClearNotification()

		return m, nil
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

		if err := m.loadDirectory(dir); err != nil {
			return m, func() tea.Msg {
				return errorMessage{
					Message: err.Error(),
				}
			}
		}
		m.explorerModel.FocusPath(msg.Path)

		return m, nil
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
				if err := m.loadDirectory(entry.GetPath()); err != nil {
					return m, func() tea.Msg {
						return errorMessage{
							Message: err.Error(),
						}
					}
				}

				return m, nil
			}
		}

		return m, nil
	case actions.NavigationActionBack:
		parentPath := filepath.Dir(m.currentPath)
		lastPath := m.currentPath

		if err := m.loadDirectory(parentPath); err != nil {
			return m, func() tea.Msg {
				return errorMessage{
					Message: err.Error(),
				}
			}
		}
		m.explorerModel.FocusPath(lastPath)

		return m, nil
	case actions.NavigationActionChangeDirectory:
		if err := m.loadDirectory(msg.Path); err != nil {
			return m, func() tea.Msg {
				return errorMessage{
					Message: err.Error(),
				}
			}
		}

		return m, nil
	}

	return m, nil
}

// handleSelectionMessage processes selection actions
func (m Model) handleSelectionMessage(msg actions.SelectionMessage) (tea.Model, tea.Cmd) {
	switch msg.Action {
	case actions.SelectionActionToggle:
		m.explorerModel.ToggleSelection()
	case actions.SelectionActionClear:
		m.explorerModel.ClearSelections()
	case actions.SelectionActionAll:
		m.explorerModel.SelectAll()
	}

	return m, nil
}

// handleUIMessage processes UI control actions
func (m Model) handleUIMessage(msg actions.UIMessage) (tea.Model, tea.Cmd) {
	switch msg.Action {
	case actions.UIActionToggleHidden:
		m.showHidden = !m.showHidden
		if err := m.loadDirectory(m.currentPath); err != nil {
			return m, func() tea.Msg {
				return errorMessage{
					Message: err.Error(),
				}
			}
		}

		return m, nil
	case actions.UIActionRefresh:
		if err := m.loadDirectory(m.currentPath); err != nil {
			return m, func() tea.Msg {
				return errorMessage{
					Message: err.Error(),
				}
			}
		}

		return m, nil
	}

	return m, nil
}

// handleSortingMessage processes sorting actions
func (m Model) handleSortingMessage(msg actions.SortingMessage) (tea.Model, tea.Cmd) {
	switch msg.SortType {
	case actions.SortTypeReverse:
		m.reverse = !m.reverse
	default:
		m.sortType = types.SortType(msg.SortType)
	}

	// Reload directory with new sorting
	if err := m.loadDirectory(m.currentPath); err != nil {
		return m, func() tea.Msg {
			return errorMessage{
				Message: err.Error(),
			}
		}
	}

	return m, nil
}

// handleFocusByIndexMessage processes focus by index actions
func (m Model) handleFocusByIndexMessage(msg actions.FocusByIndexMessage) (tea.Model, tea.Cmd) {
	m.explorerModel.SetFocusByIndex(msg.Index)

	return m, nil
}

// handleToggleSelectionByPathMessage processes toggle selection by path actions
func (m Model) handleToggleSelectionByPathMessage(
	msg actions.ToggleSelectionByPathMessage,
) (tea.Model, tea.Cmd) {
	m.explorerModel.ToggleSelectionByPath(msg.Path)

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
	if err := m.loadDirectory(msg.Path); err != nil {
		return m, func() tea.Msg {
			return errorMessage{
				Message: err.Error(),
			}
		}
	}

	return m, nil
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

// loadDirectory loads directory contents synchronously and updates model state
func (m *Model) loadDirectory(path string) error {
	// Use the existing fs.LoadEntries function with provided values
	entries, err := fs.LoadEntries(path, m.showHidden, m.sortType.String(), m.reverse, false, false)
	if err != nil {
		return fmt.Errorf("failed to load directory %s: %w", path, err)
	}

	m.currentPath = path
	m.explorerModel.SetEntries(entries)

	return nil
}

// parseCommand parses a shell command line, properly handling:
// - Single quotes: preserve all characters literally (no variable expansion)
// - Double quotes: preserve spaces but allow variable expansion
// - Unquoted spaces: act as token separators
// Returns the command name and its arguments as separate values
func parseCommand(content string) (string, []string) {
	const (
		singleQuote = '\''
		doubleQuote = '"'
		space       = ' '
		tab         = '\t'
		newline     = '\n'
	)

	content = strings.TrimSpace(content)
	if content == "" {
		return "", nil
	}

	var tokens []string
	var tokenBuilder strings.Builder
	insideSingleQuotes := false
	insideDoubleQuotes := false

	for _, char := range content {
		switch char {
		case singleQuote:
			if !insideDoubleQuotes {
				insideSingleQuotes = !insideSingleQuotes
				// Don't include the quotes in the token
			} else {
				tokenBuilder.WriteRune(char)
			}
		case doubleQuote:
			if !insideSingleQuotes {
				insideDoubleQuotes = !insideDoubleQuotes
				// Don't include the quotes in the token
			} else {
				tokenBuilder.WriteRune(char)
			}
		case space, tab, newline:
			insideAnyQuotes := insideSingleQuotes || insideDoubleQuotes
			if insideAnyQuotes {
				// Preserve whitespace inside quotes
				tokenBuilder.WriteRune(char)
			} else {
				// End of token - whitespace acts as separator
				if tokenBuilder.Len() > 0 {
					tokens = append(tokens, tokenBuilder.String())
					tokenBuilder.Reset()
				}
			}
		default:
			tokenBuilder.WriteRune(char)
		}
	}

	// Add the last token if any remains
	if tokenBuilder.Len() > 0 {
		tokens = append(tokens, tokenBuilder.String())
	}

	if len(tokens) == 0 {
		return "", nil
	}

	// Return command name and arguments separately
	return tokens[0], tokens[1:]
}
