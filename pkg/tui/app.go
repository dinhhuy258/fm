package tui

import (
	"fmt"
	"os"
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
	actionHandler *actions.ActionHandler
	modeManager   *ModeManager
	keyManager    *KeyManager
}

// directoryLoadedMsg indicates that a directory has been loaded
type directoryLoadedMsg struct {
	path      string
	entries   []fs.IEntry
	focusPath optional.Optional[string]
}

// PipeMessage represents a message from pipe
type PipeMessage struct {
	Command string
}

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
	return m.handleMessage(msg)
}

// View renders the UI
func (m Model) View() string {
	// If help UI is visible, render it as an overlay
	if m.helpModel.IsVisible() {
		return m.helpModel.View()
	}

	var sections []string

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
	currentMode := m.modeManager.GetCurrentMode()
	if selectedCount > 0 {
		modeInfo = fmt.Sprintf("Mode: %s | Items: %d | Selected: %d",
			currentMode, totalCount, selectedCount)
	} else {
		modeInfo = fmt.Sprintf("Mode: %s | Items: %d",
			currentMode, totalCount)
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
	// Help hint
	helpHint := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Render("Press ? for help, q to quit")

	return helpHint
}

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