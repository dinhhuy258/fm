package tui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/dinhhuy258/fm/pkg/actions"
	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/pipe"
	"github.com/dinhhuy258/fm/pkg/type/optional"
)

// Model represents the fm application state
type Model struct {
	// Core application state
	currentPath string

	// Models for fm components
	explorerModel     *ExplorerModel
	notificationModel *NotificationModel
	inputModel        *InputModel
	helpModel         *HelpModel

	pipe          *pipe.Pipe
	modeManager   *ModeManager
	keyManager    *KeyManager
	actionHandler *actions.ActionHandler
}

// Message types for Bubble Tea
type directoryLoadedMsg struct {
	path      string
	entries   []fs.IEntry
	focusPath optional.Optional[string]
}

// PipeMessage represents a message from external sources (pipes, etc.)
type PipeMessage struct {
	Command string
}

// ============================================================================
// CONSTRUCTOR
// ============================================================================

// NewModel creates a new root model
func NewModel(pipe *pipe.Pipe) Model {
	explorerModel := NewExplorerModel()
	notificationModel := NewNotificationModel()
	inputModel := NewInputModel()
	helpModel := NewHelpModel()

	modeManager := NewModeManager()
	keyManager := NewKeyManager(modeManager)

	actionHandler := actions.NewActionHandler()

	return Model{
		currentPath:       "",
		explorerModel:     explorerModel,
		notificationModel: notificationModel,
		inputModel:        inputModel,
		helpModel:         helpModel,
		pipe:              pipe,
		modeManager:       modeManager,
		keyManager:        keyManager,
		actionHandler:     actionHandler,
	}
}

// ============================================================================
// BUBBLETEA LIFECYCLE METHODS
// ============================================================================

// Init initializes the model
func (m Model) Init() tea.Cmd {
	// Get current working directory and load files
	wd, err := os.Getwd()
	if err != nil {
		return tea.Quit
	}

	return loadDirectoryCmd(wd, optional.NewEmpty[string]())
}

// Update handles incoming messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		headerHeight := 3
		footerHeight := 1
		interactiveHeight := 1
		availableExplorerHeight := msg.Height - headerHeight - footerHeight - interactiveHeight

		m.helpModel.SetSize(msg.Width, msg.Height)
		m.inputModel.SetSize(msg.Width, 1)
		m.notificationModel.SetSize(msg.Width, 1)
		m.explorerModel.SetSize(msg.Width, availableExplorerHeight)
	case tea.KeyMsg:
		if m.helpModel.IsVisible() {
			m.helpModel.Update(msg)

			return m, nil
		}
		if msg.String() == "?" {
			m.helpModel.Show()

			return m, nil
		}

		return m.handleKeyMap(msg)
	case AutoClearMessage:
		m.notificationModel.ClearNotification()
	default:
		return m.handleMessage(msg)
	}

	return m, nil
}

// View renders the UI
func (m Model) View() string {
	// If help UI is visible, render it as an overlay
	if m.helpModel.IsVisible() {
		return m.helpModel.View()
	}

	var sections []string

	// Header
	sections = append(sections, m.renderHeader())
	sections = append(sections, m.explorerModel.View())
	if m.inputModel.IsVisible() {
		sections = append(sections, m.inputModel.View())
	} else if m.notificationModel.IsVisible() {
		sections = append(sections, m.notificationModel.View())
	}
	sections = append(sections, m.renderFooter())

	return strings.Join(sections, "\n")
}

// ============================================================================
// MESSAGE HANDLING
// ============================================================================

// handleMessage processes incoming messages and returns updated model with commands
func (m Model) handleMessage(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case directoryLoadedMsg:
		m.currentPath = msg.path
		m.explorerModel.SetEntries(msg.entries)
		msg.focusPath.IfPresent(func(path *string) {
			m.explorerModel.FocusPath(*path)
		})
	case PipeMessage:
		return m.handlePipeMessage(msg.Command)
	// Handle new message types from MessageExecutor
	case actions.ErrorMessage:
		m.ShowNotification(NotificationError, msg.Err.Error())
	case actions.ModeChangedMessage:
		err := m.modeManager.SwitchToMode(msg.Mode)
		if err != nil {
			cmds = append(cmds, m.ShowError(fmt.Sprintf("Failed to switch mode: %v", err)))
		} else {
			// Hide input when switching to default mode
			if msg.Mode == "default" {
				m.HideInput()
				m.inputModel.ClearBuffer() // Clear input buffer when leaving input modes
			}
		}
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
		m.SetInputBuffer(msg.Value)
		m.ShowInput(msg.Value)
	case actions.UpdateInputBufferFromKeyMessage:
		// This would be handled in the key press context
	case actions.FocusPathMessage:
		// Load the directory containing the path and focus on it
		dir := filepath.Dir(msg.Path)

		return m, loadDirectoryCmd(dir, optional.New(msg.Path))
	case actions.BashOutputMessage:
		// Remove debug output logging
		// Only show user-relevant bash output as success notifications
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
		// m.SetInputBuffer(msg.Value)
	}

	return m, tea.Batch(cmds...)
}

// handleKeyMap handles key presses and resolves them to actions
func (m Model) handleKeyMap(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// First, try dynamic key mapping
	action := m.keyManager.ResolveKeyAction(msg)
	if action != nil {
		// Execute the configured messages
		executedCmds := m.actionHandler.ExecuteMessages(action.Messages)
		cmds = append(cmds, executedCmds...)

		// Handle special input buffer updates
		// if m.shouldUpdateInputBufferFromKey(action) {
		if m.inputModel.IsVisible() {
			m.UpdateInputBufferFromKey(msg.String())
			// Update the input model if it's in input mode
			m.ShowInput(m.GetInputBuffer())
		}

		return m, tea.Batch(cmds...)
	}

	cmds = append(cmds, m.ShowWarning(
		fmt.Sprintf("No action found for key: %s", msg.String()),
	))

	return m, tea.Batch(cmds...)
}

// handlePipeMessage processes messages received from the pipe (from bash scripts)
func (m Model) handlePipeMessage(command string) (tea.Model, tea.Cmd) {
	// Parse the pipe message - format is usually: CommandName arg1 arg2 ...
	commandName, args := parseCommand(command)
	if commandName == "" {
		return m, nil
	}

	// Create a synthetic MessageConfig from the pipe command
	message := &config.MessageConfig{
		Name: commandName,
		Args: args,
	}

	// Execute the command using the message executor
	cmd := m.actionHandler.ExecuteMessage(message)
	if cmd != nil {
		return m, cmd
	}

	return m, nil
}

// ============================================================================
// SPECIFIC MESSAGE HANDLERS
// ============================================================================

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

	// Selection count is automatically updated via GetStats() in renderHeader()

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
	env = append(env, fmt.Sprintf("FM_FOCUS_PATH=%s", focusPath))
	env = append(env, fmt.Sprintf("FM_PWD=%s", m.currentPath))
	env = append(env, fmt.Sprintf("FM_FOCUS_IDX=%d", focusIndex))
	env = append(env, fmt.Sprintf("FM_INPUT_BUFFER=%s", m.GetInputBuffer()))

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

// ============================================================================
// MODEL METHODS
// ============================================================================

// GetCurrentMode returns the current mode
func (m Model) GetCurrentMode() string {
	return m.modeManager.GetCurrentMode()
}

// SwitchMode switches to a new mode
func (m *Model) SwitchMode(modeName string) error {
	err := m.modeManager.SwitchToMode(modeName)
	if err != nil {
		return err
	}

	// Handle input state when switching modes
	if modeName == "default" {
		// Hide input and show notifications when switching to default mode
		m.inputModel.Hide()
		m.notificationModel.Show()
		m.inputModel.ClearBuffer()
	}

	return nil
}

// SetInputBuffer sets the input buffer value
func (m *Model) SetInputBuffer(value string) {
	m.inputModel.SetBuffer(value)
}

// GetInputBuffer returns the current input buffer value
func (m Model) GetInputBuffer() string {
	return m.inputModel.GetBuffer()
}

// UpdateInputBufferFromKey updates the input buffer with the last key press
func (m *Model) UpdateInputBufferFromKey(keyStr string) {
	m.inputModel.AppendToBuffer(keyStr)
}

// ShowInput switches to input mode and displays the text input field
func (m *Model) ShowInput(initialValue string) {
	m.notificationModel.Hide()
	m.inputModel.Show(initialValue)
}

// HideInput switches back to notification mode
func (m *Model) HideInput() {
	m.inputModel.Hide()
	m.notificationModel.Show()
}

// GetInputValue returns the current input value
func (m *Model) GetInputValue() string {
	return m.inputModel.GetValue()
}

// ShowNotification displays a notification
func (m *Model) ShowNotification(notificationType NotificationType, message string) tea.Cmd {
	cmd := m.notificationModel.ShowNotification(notificationType, message)

	if m.inputModel.IsVisible() {
		// In input mode, notification is stored but not displayed yet
		return cmd
	}

	// In notification mode, ensure notification is visible
	m.notificationModel.Show()

	return cmd
}

// ShowSuccess displays a success notification (auto-clears in 5 seconds)
func (m *Model) ShowSuccess(message string) tea.Cmd {
	return m.ShowNotification(NotificationSuccess, message)
}

// ShowInfo displays an info notification
func (m *Model) ShowInfo(message string) tea.Cmd {
	return m.ShowNotification(NotificationInfo, message)
}

// ShowWarning displays a warning notification
func (m *Model) ShowWarning(message string) tea.Cmd {
	return m.ShowNotification(NotificationWarning, message)
}

// ShowError displays an error notification
func (m *Model) ShowError(message string) tea.Cmd {
	return m.ShowNotification(NotificationError, message)
}

// GetActiveNotification returns the current active notification
func (m *Model) GetActiveNotification() *Notification {
	return m.notificationModel.GetActiveNotification()
}

// ClearNotification clears the active notification
func (m *Model) ClearNotification() {
	m.notificationModel.ClearNotification()
}

// GetTextInput returns the text input model for direct manipulation
func (m *Model) GetTextInput() *textinput.Model {
	return m.inputModel.GetTextInput()
}

// UpdateTextInput updates the text input model
func (m *Model) UpdateTextInput(textInput textinput.Model) {
	m.inputModel.UpdateTextInput(textInput)
}

// ============================================================================
// RENDERING FUNCTIONS
// ============================================================================

// renderHeader renders the header section
func (m Model) renderHeader() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		Render("File Manager")

	path := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		Render(m.currentPath)

	// Combine mode and items information
	var modeInfo string
	totalCount, selectedCount := m.explorerModel.GetStats()
	if selectedCount > 0 {
		modeInfo = fmt.Sprintf("Mode: %s | Items: %d | Selected: %d",
			m.GetCurrentMode(), totalCount, selectedCount)
	} else {
		modeInfo = fmt.Sprintf("Mode: %s | Items: %d",
			m.GetCurrentMode(), totalCount)
	}

	mode := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Render(modeInfo)

	line1 := lipgloss.JoinHorizontal(lipgloss.Left, title, " ", path)
	line2 := mode

	header := lipgloss.JoinVertical(lipgloss.Left, line1, line2, "")

	return header
}

// renderFooter renders the footer section
func (m Model) renderFooter() string {
	// Show input buffer status if not in input mode but buffer has content
	// if !m.IsInputMode() && m.GetInputBuffer() != "" {
	// 	inputBuffer := lipgloss.NewStyle().
	// 		Foreground(lipgloss.Color("#626262")).
	// 		Render(fmt.Sprintf("Buffer: %s | Press ? for help, q to quit", m.GetInputBuffer()))

	// 	return inputBuffer
	// }

	// Help hint
	helpHint := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Render("Press ? for help, q to quit")

	return helpHint
}

// ============================================================================
// COMMAND FUNCTIONS
// ============================================================================

// loadDirectoryCmd loads directory contents
func loadDirectoryCmd(path string, focusPath optional.Optional[string]) tea.Cmd {
	return func() tea.Msg {
		entries, err := loadDirectory(path)
		if err != nil {
			return actions.ErrorMessage{Err: err}
		}

		return directoryLoadedMsg{path: path, entries: entries, focusPath: focusPath}
	}
}

// ============================================================================
// UTILITY FUNCTIONS
// ============================================================================

// loadDirectory loads and returns directory entries
func loadDirectory(path string) ([]fs.IEntry, error) {
	// Get configuration values
	cfg := config.AppConfig.General
	showHidden := cfg.ShowHidden
	sortType := cfg.Sorting.SortType
	reverse := false
	if cfg.Sorting.Reverse != nil {
		reverse = *cfg.Sorting.Reverse
	}

	// Use the existing fs.LoadEntries function with config values
	entries, err := fs.LoadEntries(path, showHidden, sortType, reverse, false, false)
	if err != nil {
		return nil, err
	}

	return entries, nil
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