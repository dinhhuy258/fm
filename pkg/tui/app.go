package tui

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/dinhhuy258/fm/pkg/actions"
	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/fm/pkg/pipe"
	"github.com/dinhhuy258/fm/pkg/types"
)

// Model represents the fm application state
type Model struct {
	// Core application state
	currentPath string

	// Display and sorting settings
	showHidden bool
	sortType   types.SortType
	reverse    bool

	// Models for fm components
	explorerModel     *ExplorerModel
	notificationModel *NotificationModel
	inputModel        *InputModel
	helpModel         *HelpModel

	pipe          *pipe.Pipe
	actionHandler *actions.ActionHandler
	modeManager   *ModeManager
	keyManager    *KeyManager

	// Styles for header and footer
	titleStyle    lipgloss.Style
	modeStyle     lipgloss.Style
	helpHintStyle lipgloss.Style
	borderStyle   lipgloss.Style
}

// NewModel creates a new root model
func NewModel(pipe *pipe.Pipe) Model {
	explorerModel := NewExplorerModel()
	notificationModel := NewNotificationModel()
	inputModel := NewInputModel()

	modeManager := NewModeManager()
	helpModel := NewHelpModel(modeManager)
	keyManager := NewKeyManager(modeManager)

	actionHandler := actions.NewActionHandler()

	// Initialize sorting and display settings from config
	showHidden := config.AppConfig.General.ShowHidden
	sortType := types.SortType(config.AppConfig.General.Sorting.SortType)
	reverse := false
	if config.AppConfig.General.Sorting.Reverse != nil {
		reverse = *config.AppConfig.General.Sorting.Reverse
	}

	titleStyle := lipgloss.NewStyle()
	modeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(SecondaryTextColor))
	helpHintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(SecondaryTextColor))
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		PaddingLeft(1).
		PaddingRight(1)

	return Model{
		currentPath:       "",
		showHidden:        showHidden,
		sortType:          sortType,
		reverse:           reverse,
		explorerModel:     explorerModel,
		notificationModel: notificationModel,
		inputModel:        inputModel,
		helpModel:         helpModel,
		pipe:              pipe,
		modeManager:       modeManager,
		keyManager:        keyManager,
		actionHandler:     actionHandler,
		titleStyle:        titleStyle,
		modeStyle:         modeStyle,
		helpHintStyle:     helpHintStyle,
		borderStyle:       borderStyle,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	// Get current working directory and load files
	wd, err := os.Getwd()
	if err != nil {
		return tea.Quit
	}

	return func() tea.Msg {
		return actions.ChangeDirectoryMessage{Path: wd}
	}
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

	content := strings.Join(sections, "\n")

	return m.borderStyle.Render(content)
}

// renderHeader renders the header section
func (m Model) renderHeader() string {
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

	title := lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.titleStyle.Render(ExplorerTitle),
		": ",
		m.titleStyle.Render(m.currentPath),
	)
	mode := m.modeStyle.Render(modeInfo)

	return lipgloss.JoinVertical(lipgloss.Left, title, mode, "")
}

// renderFooter renders the footer section
func (m Model) renderFooter() string {
	return m.helpHintStyle.Render("Press " + HelpToggleKey + " for help")
}
