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
	"github.com/dinhhuy258/fm/pkg/components"
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
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

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
		m.helpUI.SetSize(msg.Width, msg.Height)

		// Update interactive area size (only 1 line needed)
		m.interactiveArea.SetSize(msg.Width, 1)

		// Set explorer table to use remaining available height
		m.explorerTable.SetSize(msg.Width, availableHeight)

	case tea.KeyMsg:
		// If help UI is visible, let it handle the key first
		if m.helpUI.IsVisible() {
			m.helpUI, cmd = m.helpUI.Update(msg)
			cmds = append(cmds, cmd)
			// If help UI is still visible after update, don't process other keys
			if m.helpUI.IsVisible() {
				return m, tea.Batch(cmds...)
			}
		}

		// Handle key presses using dynamic system
		return m.handleKeyWithDynamicSystem(msg)

	case directoryLoadedMsg:
		// Directory content loaded
		m.currentPath = msg.path
		m.entries = msg.entries

		// Update explorer table
		m.explorerTable.SetEntries(msg.entries)
		m.statusBar.path = m.currentPath
		m.statusBar.total = len(msg.entries)

		// Reset cursor and selection
		m.cursor = 0
		m.selected = make(map[int]struct{})

	case directoryLoadedWithFocusMsg:
		// Directory content loaded with focus requirement
		m.currentPath = msg.path
		m.entries = msg.entries

		// Update explorer table
		m.explorerTable.SetEntries(msg.entries)
		m.statusBar.path = m.currentPath
		m.statusBar.total = len(msg.entries)

		// Reset cursor and selection
		m.cursor = 0
		m.selected = make(map[int]struct{})

		// Focus on the specific path (remove debug logging)
		m.explorerTable.FocusPath(msg.focusPath)

	case errorMsg:
		m.err = msg.err
		// Show error notification
		cmds = append(cmds, m.interactiveArea.ShowError(fmt.Sprintf("Error: %v", msg.err)))

	case ExternalMessage:
		// Parse and execute the pipe message (remove debug logging)
		return m.handlePipeMessage(msg.Content)

	// Handle new message types from MessageExecutor
	case actions.ModeChangedMessage:
		// Actually switch the mode in the mode manager
		err := m.modeManager.SwitchToMode(msg.NewMode)
		if err != nil {
			cmds = append(cmds, m.interactiveArea.ShowError(fmt.Sprintf("Failed to switch mode: %v", err)))
		} else {
			m.statusBar.mode = msg.NewMode
			// Remove debug logging of mode changes

			// Clear dynamic keymap cache when mode changes
			m.dynamicKeyMap.ClearCache()

			// Hide input when switching to default mode
			if msg.NewMode == "default" {
				hideCmd := m.interactiveArea.HideInput()
				if hideCmd != nil {
					cmds = append(cmds, hideCmd)
				}
				m.inputBuffer = "" // Clear input buffer when leaving input modes
			}
		}

	case actions.LogMessage:
		switch msg.Level {
		case "error":
			cmds = append(cmds, m.interactiveArea.ShowError(msg.Message))
		case "warning":
			cmds = append(cmds, m.interactiveArea.ShowWarning(msg.Message))
		case "success":
			cmds = append(cmds, m.interactiveArea.ShowSuccess(msg.Message))
		case "info":
			cmds = append(cmds, m.interactiveArea.ShowInfo(msg.Message))
		default:
			cmds = append(cmds, m.interactiveArea.ShowInfo(msg.Message))
		}

	case actions.ErrorMessage:
		m.err = msg.Err
		cmds = append(cmds, m.interactiveArea.ShowError(fmt.Sprintf("Error: %v", msg.Err)))

	case actions.SetInputBufferMessage:
		m.SetInputBuffer(msg.Value)
		if msg.ShowInput {
			m.interactiveArea.ShowInput(msg.Value)
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
			cmds = append(cmds, m.interactiveArea.ShowInfo(strings.TrimSpace(msg.Output)))
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
	case components.AutoClearMessage:
		// Pass this message to the InteractiveArea
		m.interactiveArea, cmd = m.interactiveArea.Update(msg)
		cmds = append(cmds, cmd)
	case components.InputCompletedMessage:
		// Input was completed - update input buffer and show success notification
		m.inputBuffer = msg.Value
		if msg.Value != "" {
			cmds = append(cmds, m.interactiveArea.ShowSuccess(fmt.Sprintf("Input: %s", msg.Value)))
		}
	default:
		// Always update help UI (it handles its own visibility)
		m.helpUI, cmd = m.helpUI.Update(msg)
		cmds = append(cmds, cmd)

		// Always update interactive area
		m.interactiveArea, cmd = m.interactiveArea.Update(msg)
		cmds = append(cmds, cmd)

		// Update input buffer from interactive area if in input mode
		if m.interactiveArea.IsInputMode() {
			m.inputBuffer = m.interactiveArea.GetInputValue()
		}
	}

	return m, tea.Batch(cmds...)
}

// handleKeyWithDynamicSystem handles key presses using the dynamic mode system
func (m Model) handleKeyWithDynamicSystem(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	updatedModel := m

	// First, try dynamic key mapping
	keyStr, action, matched := m.dynamicKeyMap.MatchesKey(msg)
	if matched {
		// Execute the configured messages
		executedCmds := m.messageExecutor.ExecuteMessages(action.Messages, m.currentPath, m.inputBuffer)
		cmds = append(cmds, executedCmds...)

		// Handle special input buffer updates
		if m.shouldUpdateInputBufferFromKey(action, keyStr) {
			updatedModel.UpdateInputBufferFromKey(keyStr)
			// Update the interactive area if it's in input mode
			if updatedModel.interactiveArea.IsInputMode() {
				updatedModel.interactiveArea.ShowInput(updatedModel.inputBuffer)
			}
		}

		return updatedModel, tea.Batch(cmds...)
	}

	m.interactiveArea.ShowWarning(fmt.Sprintf("No action found for key: %s", m.dynamicKeyMap.keyMsgToString(msg)))

	return m, nil
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
	if m.helpUI.IsVisible() {
		// Overlay the help UI on top
		helpView := m.helpUI.View()

		// Position the help UI in the center
		return lipgloss.Place(
			m.termWidth, m.termHeight,
			lipgloss.Center, lipgloss.Center,
			helpView,
			lipgloss.WithWhitespaceChars(""),
		)
	}

	var sections []string

	// Header
	sections = append(sections, m.renderHeader())

	// Main content area
	sections = append(sections, m.explorerTable.View())

	// Interactive area (handles both input and notifications)
	interactiveView := m.interactiveArea.View()
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
	if len(m.selected) > 0 {
		modeInfo = fmt.Sprintf("Mode: %s | Items: %d | Selected: %d",
			m.GetCurrentMode(), m.statusBar.total, len(m.selected))
	} else {
		modeInfo = fmt.Sprintf("Mode: %s | Items: %d",
			m.GetCurrentMode(), m.statusBar.total)
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
	if !m.interactiveArea.IsInputMode() && m.inputBuffer != "" {
		inputBuffer := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Render(fmt.Sprintf("Buffer: %s | Press ? for help, q to quit", m.inputBuffer))

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
		m.explorerTable.Move(1)
	case "previous":
		m.explorerTable.Move(-1)
	case "first":
		m.explorerTable.MoveFirst()
	case "last":
		m.explorerTable.MoveLast()
	case "enter":
		if entry := m.explorerTable.GetFocusedEntry(); entry != nil {
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

	// Update selection count
	_, selected := m.explorerTable.GetStats()
	m.statusBar.selected = selected

	return m, nil
}

// handleSelectionMessage processes selection actions
func (m Model) handleSelectionMessage(msg actions.SelectionMessage) (tea.Model, tea.Cmd) {
	switch msg.Action {
	case "toggle":
		// Toggle selection of current item using direct method
		m.explorerTable.ToggleSelection()
	case "clear":
		// Clear all selections (this would need to be implemented in the explorer table)
		m.selected = make(map[int]struct{})
		// Selections cleared silently
	case "all":
		// Select all items (this would need to be implemented in the explorer table)
		for i := range m.entries {
			m.selected[i] = struct{}{}
		}
		// Selected all items silently
	}

	// Update selection count
	_, selected := m.explorerTable.GetStats()
	m.statusBar.selected = selected

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
	cmd := m.messageExecutor.ExecuteMessage(message, m.currentPath, m.inputBuffer)
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

	// Validate index bounds
	if index < 0 || index >= len(m.entries) {
		// Index out of bounds - handle silently
		return m, nil
	}

	// Update focus through explorer table (silently)
	m.explorerTable.SetFocusByIndex(index)

	return m, nil
}

// handleToggleSelectionByPathMessage processes toggle selection by path actions
func (m Model) handleToggleSelectionByPathMessage(msg actions.ToggleSelectionByPathMessage) (tea.Model, tea.Cmd) {
	path := msg.Path

	// Find the entry with the given path
	var foundEntry fs.IEntry
	for _, entry := range m.entries {
		if entry.GetPath() == path {
			foundEntry = entry

			break
		}
	}

	if foundEntry == nil {
		// Path not found - handle silently
		return m, nil
	}

	// Toggle selection in explorer table (silently)
	m.explorerTable.ToggleSelectionByPath(path)

	// Update selection count
	_, selected := m.explorerTable.GetStats()
	m.statusBar.selected = selected

	return m, nil
}

// handleWriteSelectionsMessage processes selection writing and bash execution
func (m Model) handleWriteSelectionsMessage(msg actions.WriteSelectionsMessage) (tea.Model, tea.Cmd) {
	// Get current selections and focus index from explorer table
	selectedEntries := m.explorerTable.GetSelectedEntries()
	focusIndex := m.explorerTable.GetFocus()

	// Get the focused entry's path
	focusPath := msg.CurrentPath // Default to current directory
	if focusedEntry := m.explorerTable.GetFocusedEntry(); focusedEntry != nil {
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