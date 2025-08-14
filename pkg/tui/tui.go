package tui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/dinhhuy258/fm/pkg/actions"
	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/fm/pkg/fs"
)

// Init initializes the model
func (m Model) Init() tea.Cmd {
	// Get current working directory and load files
	wd, err := os.Getwd()
	if err != nil {
		return tea.Quit
	}

	return loadDirectoryCmd(wd)
}

// Update handles incoming messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height

		if !m.ready {
			m.ready = true
		}

		// Update component sizes
		headerHeight := 3
		footerHeight := 1      // Further reduced since we removed input from footer
		interactiveHeight := 1 // Always reserve space for interactive area
		availableHeight := msg.Height - headerHeight - footerHeight - interactiveHeight

		// Update help UI size (always update for when it's shown)
		m.helpModel.SetSize(msg.Width, msg.Height)

		// Update interactive area size (only 1 line needed)
		m.notificationModel.SetSize(msg.Width, 1)
		m.inputModel.SetSize(msg.Width, 1)

		// Set explorer table to use remaining available height
		m.explorerModel.SetSize(msg.Width, availableHeight)

	case tea.KeyMsg:
		// If help UI is visible, let it handle the key first
		if m.helpModel.IsVisible() {
			// Handle help UI key events directly since it's now a pure model
			switch {
			case msg.String() == "?" || msg.String() == "esc" || msg.String() == "q":
				m.helpModel.Hide()

				return m, nil
			case msg.String() == "k" || msg.String() == "up":
				viewport := m.helpModel.GetViewport()
				viewport.ScrollUp(1)
				m.helpModel.UpdateViewport(*viewport)
			case msg.String() == "j" || msg.String() == "down":
				viewport := m.helpModel.GetViewport()
				viewport.ScrollDown(1)
				m.helpModel.UpdateViewport(*viewport)
			case msg.String() == "pgup" || msg.String() == "ctrl+u":
				viewport := m.helpModel.GetViewport()
				viewport.HalfPageUp()
				m.helpModel.UpdateViewport(*viewport)
			case msg.String() == "pgdown" || msg.String() == "ctrl+d":
				viewport := m.helpModel.GetViewport()
				viewport.HalfPageDown()
				m.helpModel.UpdateViewport(*viewport)
			}

			// If help UI is still visible after update, don't process other keys
			if m.helpModel.IsVisible() {
				return m, tea.Batch(cmds...)
			}
		}

		// Handle key presses using dynamic system
		// Handle special cases like showing help
		if msg.String() == "?" {
			m.helpModel.Show()

			return m, nil
		}

		return m.handleKeyWithDynamicSystem(msg)

	case directoryLoadedMsg:
		// Directory content loaded
		m.currentPath = msg.path

		// Update explorer model
		m.explorerModel.SetEntries(msg.entries)

	case directoryLoadedWithFocusMsg:
		// Directory content loaded with focus requirement
		m.currentPath = msg.path

		// Update explorer model
		m.explorerModel.SetEntries(msg.entries)

		// Focus on the specific path
		m.explorerModel.FocusPath(msg.focusPath)

	case errorMsg:
		m.err = msg.err
		// Show error notification
		cmds = append(cmds, m.ShowError(fmt.Sprintf("Error: %v", msg.err)))

	case ExternalMessage:
		// Parse and execute the pipe message (remove debug logging)
		return m.handlePipeMessage(msg.Content)

	// Handle new message types from MessageExecutor
	case actions.ModeChangedMessage:
		// Actually switch the mode in the mode manager
		err := m.modeManager.SwitchToMode(msg.NewMode)
		if err != nil {
			cmds = append(cmds, m.ShowError(fmt.Sprintf("Failed to switch mode: %v", err)))
		} else {
			// Hide input when switching to default mode
			if msg.NewMode == "default" {
				m.HideInput()
				m.inputModel.ClearBuffer() // Clear input buffer when leaving input modes
			}
		}

	case actions.LogMessage:
		switch msg.Level {
		case "error":
			cmds = append(cmds, m.ShowError(msg.Message))
		case "warning":
			cmds = append(cmds, m.ShowWarning(msg.Message))
		case "success":
			cmds = append(cmds, m.ShowSuccess(msg.Message))
		case "info":
			cmds = append(cmds, m.ShowInfo(msg.Message))
		default:
			cmds = append(cmds, m.ShowInfo(msg.Message))
		}

	case actions.ErrorMessage:
		m.err = msg.Err
		cmds = append(cmds, m.ShowError(fmt.Sprintf("Error: %v", msg.Err)))

	case actions.SetInputBufferMessage:
		m.SetInputBuffer(msg.Value)
		if msg.ShowInput {
			m.ShowInput(msg.Value)
		}

	case actions.UpdateInputBufferFromKeyMessage:
		// This would be handled in the key press context
	case actions.FocusPathMessage:
		// Load the directory containing the path and focus on it
		dir := filepath.Dir(msg.Path)

		return m, loadDirectoryWithFocusCmd(dir, msg.Path)
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
	case actions.WriteSelectionsMessage:
		return m.handleWriteSelectionsMessage(msg)
	case actions.InteractiveBashMessage:
		return m.handleInteractiveBashMessage(msg)
	case AutoClearMessage:
		// Auto-clear notification in the interactive model
		if m.interactiveMode == InteractiveModeNotification {
			m.ClearNotification()
		}
	case InputCompletedMessage:
		// Input was completed - update input buffer and show success notification
		m.SetInputBuffer(msg.Value)
		if msg.Value != "" {
			cmds = append(cmds, m.ShowSuccess(fmt.Sprintf("Input: %s", msg.Value)))
		}
	default:
		// Handle text input updates for the interactive model
		if m.IsInputMode() {
			// Check for input completion or cancellation
			if keyMsg, ok := msg.(tea.KeyMsg); ok {
				switch keyMsg.String() {
				case "enter":
					// Input completed - return to notification mode and send completion message
					inputValue := m.GetInputValue()
					m.HideInput()

					// Create completion command
					completionCmd := func() tea.Msg {
						return InputCompletedMessage{Value: inputValue}
					}

					cmds = append(cmds, completionCmd)

					return m, tea.Batch(cmds...)
				case "esc":
					// Cancel input - return to notification mode
					m.HideInput()

					return m, tea.Batch(cmds...)
				}
			}

			// Update input model directly
			var cmd tea.Cmd
			m.inputModel, cmd = m.inputModel.Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
	}

	// Update child models
	var childCmd tea.Cmd
	m.notificationModel, childCmd = m.notificationModel.Update(msg)
	if childCmd != nil {
		cmds = append(cmds, childCmd)
	}

	// Only update input model if not handled above
	if !m.IsInputMode() || msg == nil {
		m.inputModel, childCmd = m.inputModel.Update(msg)
		if childCmd != nil {
			cmds = append(cmds, childCmd)
		}
	}

	return m, tea.Batch(cmds...)
}

// handleKeyWithDynamicSystem handles key presses using the dynamic mode system
func (m Model) handleKeyWithDynamicSystem(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	updatedModel := m

	// First, try dynamic key mapping
	keyStr, action, matched := m.keyManager.ResolveKeyAction(msg)
	if matched {
		// Execute the configured messages
		executedCmds := m.messageExecutor.ExecuteMessages(action.Messages, m.currentPath, m.GetInputBuffer())
		cmds = append(cmds, executedCmds...)

		// Handle special input buffer updates
		if m.shouldUpdateInputBufferFromKey(action, keyStr) {
			updatedModel.UpdateInputBufferFromKey(keyStr)
			// Update the input model if it's in input mode
			if updatedModel.IsInputMode() {
				updatedModel.ShowInput(updatedModel.GetInputBuffer())
			}
		}

		return updatedModel, tea.Batch(cmds...)
	}

	warnCmd := m.ShowWarning(
		fmt.Sprintf("No action found for key: %s", m.keyManager.keyMsgToString(msg)),
	)
	cmds = append(cmds, warnCmd)

	return m, tea.Batch(cmds...)
}

// shouldUpdateInputBufferFromKey checks if we should update input buffer from key
func (m Model) shouldUpdateInputBufferFromKey(action *config.ActionConfig, _ string) bool {
	// Check if any message is UpdateInputBufferFromKey
	for _, message := range action.Messages {
		if message.Name == "UpdateInputBufferFromKey" {
			return true
		}
	}

	return false
}

// View renders the UI
func (m Model) View() string {
	if !m.ready {
		return "\n  Loading..."
	}

	// If help UI is visible, render it as an overlay
	if m.helpModel.IsVisible() {
		return m.helpModel.View()
	}

	var sections []string

	// Header
	sections = append(sections, m.renderHeader())

	// Main content area - use the model's View method with cached data
	sections = append(sections, m.explorerModel.View())

	// Interactive area (handles both input and notifications) - use the model's View method with cached styles
	var interactiveView string
	switch m.interactiveMode {
	case InteractiveModeInput, InteractiveModeBuffer:
		if m.inputModel.IsVisible() {
			interactiveView = m.inputModel.View()
		}
	case InteractiveModeNotification:
		if m.notificationModel.IsVisible() {
			interactiveView = m.notificationModel.View()
		}
	}
	if interactiveView != "" {
		sections = append(sections, interactiveView)
	} else {
		sections = append(sections, "") // Empty line to maintain spacing
	}

	// Footer (status bar + input)
	sections = append(sections, m.renderFooter())

	return strings.Join(sections, "\n")
}

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
	if !m.IsInputMode() && m.GetInputBuffer() != "" {
		inputBuffer := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Render(fmt.Sprintf("Buffer: %s | Press ? for help, q to quit", m.GetInputBuffer()))

		return inputBuffer
	}

	// Help hint
	helpHint := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Render("Press ? for help, q to quit")

	return helpHint
}

// Message types for Bubble Tea

type directoryLoadedMsg struct {
	path    string
	entries []fs.IEntry
}

type directoryLoadedWithFocusMsg struct {
	path      string
	entries   []fs.IEntry
	focusPath string
}

type errorMsg struct {
	err error
}

// ExternalMessage represents a message from external sources (pipes, etc.)
type ExternalMessage struct {
	Content string
}

// loadDirectoryCmd loads directory contents
func loadDirectoryCmd(path string) tea.Cmd {
	return func() tea.Msg {
		entries, err := loadDirectory(path)
		if err != nil {
			return errorMsg{err: err}
		}

		return directoryLoadedMsg{path: path, entries: entries}
	}
}

// loadDirectoryWithFocusCmd loads directory contents and focuses on a specific path
func loadDirectoryWithFocusCmd(dirPath, focusPath string) tea.Cmd {
	return func() tea.Msg {
		entries, err := loadDirectory(dirPath)
		if err != nil {
			return errorMsg{err: err}
		}

		return directoryLoadedWithFocusMsg{path: dirPath, entries: entries, focusPath: focusPath}
	}
}

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

// handleNavigationMessage processes navigation actions
func (m Model) handleNavigationMessage(msg actions.NavigationMessage) (tea.Model, tea.Cmd) {
	switch msg.Action {
	case "next":
		m.explorerModel.Move(1)
	case "previous":
		m.explorerModel.Move(-1)
	case "first":
		m.explorerModel.MoveFirst()
	case "last":
		m.explorerModel.MoveLast()
	case "enter":
		if entry := m.explorerModel.GetFocusedEntry(); entry != nil {
			if entry.IsDirectory() {
				return m, loadDirectoryCmd(entry.GetPath())
			}
		}
	case "back":
		if m.currentPath != "" && m.currentPath != "/" {
			parentPath := filepath.Dir(m.currentPath)

			return m, loadDirectoryWithFocusCmd(parentPath, m.currentPath)
		}
	case "change_directory":
		if msg.Path != "" {
			return m, loadDirectoryCmd(msg.Path)
		}
	}

	// Selection count is automatically updated via GetStats() in renderHeader()

	return m, nil
}

// handleSelectionMessage processes selection actions
func (m Model) handleSelectionMessage(msg actions.SelectionMessage) (tea.Model, tea.Cmd) {
	switch msg.Action {
	case "toggle":
		// Toggle selection of current item using direct method
		m.explorerModel.ToggleSelection()
	case "clear":
		// Clear all selections
		m.explorerModel.ClearSelections()
	case "all":
		// Select all items
		m.explorerModel.SelectAll()
	}

	// Selection count is automatically updated via GetStats() in renderHeader()

	return m, nil
}

// handleUIMessage processes UI control actions
func (m Model) handleUIMessage(msg actions.UIMessage) (tea.Model, tea.Cmd) {
	switch msg.Action {
	case "toggle_hidden":
		// Toggle hidden file visibility
		config.AppConfig.General.ShowHidden = !config.AppConfig.General.ShowHidden
		// Toggle hidden files silently

		return m, loadDirectoryCmd(m.currentPath) // Reload with new settings
	case "refresh":
		// Refresh current directory
		return m, loadDirectoryCmd(m.currentPath)
	}

	return m, nil
}

// parseShellCommand parses a shell command line, handling quoted strings properly
func parseShellCommand(content string) (string, []string) {
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

// handlePipeMessage processes messages received from the pipe (from bash scripts)
func (m Model) handlePipeMessage(content string) (tea.Model, tea.Cmd) {
	// Parse the pipe message - format is usually: CommandName arg1 arg2 ...
	commandName, args := parseShellCommand(content)
	if commandName == "" {
		return m, nil
	}

	// Create a synthetic MessageConfig from the pipe command
	message := &config.MessageConfig{
		Name: commandName,
		Args: args,
	}

	// Execute the command using the message executor
	cmd := m.messageExecutor.ExecuteMessage(message, m.currentPath, m.GetInputBuffer())
	if cmd != nil {
		return m, cmd
	}

	return m, nil
}

// handleSortingMessage processes sorting actions
func (m Model) handleSortingMessage(msg actions.SortingMessage) (tea.Model, tea.Cmd) {
	switch msg.SortType {
	case "name":
		// Update config and reload directory
		config.AppConfig.General.Sorting.SortType = "name"
		// Sort by name silently

		return m, loadDirectoryCmd(m.currentPath)
	case "size":
		config.AppConfig.General.Sorting.SortType = "size"
		// Sort by size silently

		return m, loadDirectoryCmd(m.currentPath)
	case "date":
		config.AppConfig.General.Sorting.SortType = "date_modified"
		// Sort by date silently

		return m, loadDirectoryCmd(m.currentPath)
	case "extension":
		config.AppConfig.General.Sorting.SortType = "extension"
		// Sort by extension silently

		return m, loadDirectoryCmd(m.currentPath)
	case "dir_first":
		config.AppConfig.General.Sorting.SortType = "dir_first"
		// Sort by directory first silently

		return m, loadDirectoryCmd(m.currentPath)
	case "reverse":
		// Toggle reverse sorting
		if config.AppConfig.General.Sorting.Reverse != nil {
			*config.AppConfig.General.Sorting.Reverse = !*config.AppConfig.General.Sorting.Reverse
		} else {
			reverse := true
			config.AppConfig.General.Sorting.Reverse = &reverse
		}
		// Sort order changed silently

		return m, loadDirectoryCmd(m.currentPath)
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

// handleWriteSelectionsMessage processes selection writing and bash execution
func (m Model) handleWriteSelectionsMessage(msg actions.WriteSelectionsMessage) (tea.Model, tea.Cmd) {
	// Get current selections and focus index from explorer model
	selectedEntries := m.explorerModel.GetSelectedEntries()
	focusIndex := m.explorerModel.GetFocus()

	// Get the focused entry's path
	focusPath := msg.CurrentPath // Default to current directory
	if focusedEntry := m.explorerModel.GetFocusedEntry(); focusedEntry != nil {
		focusPath = focusedEntry.GetPath()
	}

	// Convert selected entries to paths
	selections := make([]string, len(selectedEntries))
	for i, entry := range selectedEntries {
		selections[i] = entry.GetPath()
	}

	// Execute bash with proper environment
	return m, m.messageExecutor.ExecuteBashWithEnv(
		msg.Script,
		msg.CurrentPath,
		focusPath,
		msg.InputBuffer,
		msg.Silent,
		selections,
		focusIndex,
	)
}

// handleInteractiveBashMessage processes interactive bash execution with TUI suspension
func (m Model) handleInteractiveBashMessage(msg actions.InteractiveBashMessage) (tea.Model, tea.Cmd) {
	return m, tea.ExecProcess(&exec.Cmd{
		Path: "/bin/bash",
		Args: []string{"bash", "-c", "clear && " + msg.Script},
		Env:  msg.Environment,
		Dir:  msg.WorkingDir,
	}, func(err error) tea.Msg {
		if err != nil {
			return actions.ErrorMessage{Err: fmt.Errorf("interactive bash command failed: %w", err)}
		}
		// After interactive command completes, refresh the directory
		return loadDirectoryCmd(m.currentPath)()
	})
}