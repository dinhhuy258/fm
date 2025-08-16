package tui

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/dinhhuy258/fm/pkg/actions"
	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/pipe"
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

	return loadDirectoryCmd(wd)
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