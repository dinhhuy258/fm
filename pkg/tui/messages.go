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
	"github.com/dinhhuy258/fm/pkg/type/optional"
)

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
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case directoryLoadedMsg:
		m.currentPath = msg.path
		m.explorerModel.SetEntries(msg.entries)
		msg.focusPath.IfPresent(func(path *string) {
			m.explorerModel.FocusPath(*path)
		})

		return m, nil
	case PipeMessage:
		return m.handlePipeMessage(msg.Command)
	case actions.ErrorMessage:
		m.notificationModel.ShowNotification(NotificationError, msg.Err.Error())

		return m, nil
	case actions.ModeChangedMessage:
		err := m.modeManager.SwitchToMode(msg.Mode)
		if err != nil {
			cmds = append(cmds, m.ShowError(fmt.Sprintf("Failed to switch mode: %v", err)))
		}
		m.inputModel.Hide()
		m.notificationModel.Show()

		return m, tea.Batch(cmds...)
	case actions.LogMessage:
		switch msg.Level {
		case actions.LogLevelError:
			cmds = append(cmds, m.ShowError(msg.Message))
		case actions.LogLevelWarning:
			cmds = append(cmds, m.ShowWarning(msg.Message))
		case actions.LogLevelSuccess:
			cmds = append(cmds, m.ShowSuccess(msg.Message))
		case actions.LogLevelInfo:
			cmds = append(cmds, m.ShowInfo(msg.Message))
		}
	case actions.SetInputBufferMessage:
		m.notificationModel.Hide()
		m.inputModel.Show(msg.Value)
	case actions.UpdateInputBufferFromKeyMessage:
		return m, m.inputModel.Update(msg.Key)
	case actions.FocusPathMessage:
		dir := filepath.Dir(msg.Path)

		return m, loadDirectoryCmd(dir, optional.New(msg.Path))
	case actions.BashOutputMessage:
		if !msg.Silent && strings.TrimSpace(msg.Output) != "" {
			cmds = append(cmds, m.ShowInfo(strings.TrimSpace(msg.Output)))
		}
	case actions.NavigationMessage:
		return m.handleNavigationMessage(msg)
	case actions.SelectionMessage:
		return m.handleSelectionMessage(msg)
	case actions.UIMessage:
		return m.handleUIMessage(msg)
	case actions.SortingMessage:
		return m.handleSortingMessage(msg)
	case actions.FocusByIndexMessage:
		return m.handleFocusByIndexMessage(msg)
	case actions.ToggleSelectionByPathMessage:
		return m.handleToggleSelectionByPathMessage(msg)
	case actions.BashExecMessage:
		return m.handleBashExecMessage(msg)
	case actions.BashExecSilentlyMessage:
		return m.handleBashExecSilentlyMessage(msg)
	case actions.ChangeDirectoryMessage:
		return m.handleChangeDirectoryMessage(msg)
	case AutoClearMessage:
		m.ClearNotification()
	case InputCompletedMessage:
		// Input was completed - update input buffer and show success notification
	}

	return m, tea.Batch(cmds...)
}

// handleKeyMap handles key presses and resolves them to actions
func (m Model) handleKeyMap(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	action := m.keyManager.ResolveKeyAction(msg)
	if action != nil {
		return m, m.actionHandler.ExecuteMessages(action.Messages, msg)
	}

	return m, m.ShowWarning(
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
				return m, loadDirectoryCmd(entry.GetPath(), optional.NewEmpty[string]())
			}
		}
	case actions.NavigationActionBack:
		if m.currentPath != "" && m.currentPath != "/" {
			parentPath := filepath.Dir(m.currentPath)

			return m, loadDirectoryCmd(parentPath, optional.New(m.currentPath))
		}
	case actions.NavigationActionChangeDirectory:
		if msg.Path != "" {
			return m, loadDirectoryCmd(msg.Path, optional.NewEmpty[string]())
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
		config.AppConfig.General.ShowHidden = !config.AppConfig.General.ShowHidden
		// Toggle hidden files silently

		return m, loadDirectoryCmd(m.currentPath, optional.Optional[string]{}) // Reload with new settings
	case actions.UIActionRefresh:
		// Refresh current directory
		return m, loadDirectoryCmd(m.currentPath, optional.Optional[string]{})
	}

	return m, nil
}

// handleSortingMessage processes sorting actions
func (m Model) handleSortingMessage(msg actions.SortingMessage) (tea.Model, tea.Cmd) {
	switch msg.SortType {
	case actions.SortTypeName:
		// Update config and reload directory
		config.AppConfig.General.Sorting.SortType = "name"
		// Sort by name silently

		return m, loadDirectoryCmd(m.currentPath, optional.Optional[string]{})
	case actions.SortTypeSize:
		config.AppConfig.General.Sorting.SortType = "size"
		// Sort by size silently

		return m, loadDirectoryCmd(m.currentPath, optional.Optional[string]{})
	case actions.SortTypeDate:
		config.AppConfig.General.Sorting.SortType = "date_modified"
		// Sort by date silently

		return m, loadDirectoryCmd(m.currentPath, optional.Optional[string]{})
	case actions.SortTypeExtension:
		config.AppConfig.General.Sorting.SortType = "extension"
		// Sort by extension silently

		return m, loadDirectoryCmd(m.currentPath, optional.Optional[string]{})
	case actions.SortTypeDirFirst:
		config.AppConfig.General.Sorting.SortType = "dir_first"
		// Sort by directory first silently

		return m, loadDirectoryCmd(m.currentPath, optional.Optional[string]{})
	case actions.SortTypeReverse:
		// Toggle reverse sorting
		if config.AppConfig.General.Sorting.Reverse != nil {
			*config.AppConfig.General.Sorting.Reverse = !*config.AppConfig.General.Sorting.Reverse
		} else {
			reverse := true
			config.AppConfig.General.Sorting.Reverse = &reverse
		}
		// Sort order changed silently

		return m, loadDirectoryCmd(m.currentPath, optional.Optional[string]{})
	}

	return m, nil
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
func (m Model) handleToggleSelectionByPathMessage(msg actions.ToggleSelectionByPathMessage) (tea.Model, tea.Cmd) {
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
func (m Model) handleBashExecSilentlyMessage(msg actions.BashExecSilentlyMessage) (tea.Model, tea.Cmd) {
	return m.handleBashExecution(msg.Script, true)
}

// handleChangeDirectoryMessage processes directory change requests with validation
func (m Model) handleChangeDirectoryMessage(msg actions.ChangeDirectoryMessage) (tea.Model, tea.Cmd) {
	targetPath := msg.Path

	// Handle home directory expansion
	if strings.HasPrefix(targetPath, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return m, func() tea.Msg {
				return actions.ErrorMessage{Err: fmt.Errorf("failed to get home directory: %w", err)}
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
			return actions.ErrorMessage{Err: fmt.Errorf("directory does not exist: %s", targetPath)}
		}
	}

	// If all validation passes, change to the directory
	return m, loadDirectoryCmd(targetPath, optional.NewEmpty[string]())
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

	env := os.Environ()
	inputBuffer := m.inputModel.GetTextInput().Value()
	env = append(env, fmt.Sprintf("FM_FOCUS_PATH=%s", focusPath))
	env = append(env, fmt.Sprintf("FM_PWD=%s", m.currentPath))
	env = append(env, fmt.Sprintf("FM_FOCUS_IDX=%d", focusIndex))
	env = append(env, fmt.Sprintf("FM_INPUT_BUFFER=%s", inputBuffer))

	if m.pipe != nil {
		if len(selections) > 0 {
			selectionPath := m.pipe.GetSelectionPath()
			if err := writeSelectionsToFile(selectionPath, selections); err != nil {
				return m, func() tea.Msg {
					return actions.ErrorMessage{Err: fmt.Errorf("failed to write selections to pipe: %w", err)}
				}
			}
		}

		env = append(env, fmt.Sprintf("FM_PIPE_MSG_IN=%s", m.pipe.GetMessageInPath()))
		env = append(env, fmt.Sprintf("FM_PIPE_SELECTION=%s", m.pipe.GetSelectionPath()))
		env = append(env, fmt.Sprintf("FM_SESSION_PATH=%s", m.pipe.GetSessionPath()))
	}

	cmd := exec.Command("bash", "-c", script)
	cmd.Env = env
	cmd.Dir = m.currentPath

	if silent {
		if err := cmd.Start(); err != nil {
			return m, func() tea.Msg {
				return actions.ErrorMessage{Err: fmt.Errorf("failed to start bash command: %w", err)}
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
			return actions.ErrorMessage{Err: fmt.Errorf("failed to execute bash %v", err)}
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