package tui

import (
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/dinhhuy258/fm/pkg/config"
)

// HelpUIKeys defines key bindings for the help UI
type HelpUIKeys struct {
	Close    key.Binding
	Up       key.Binding
	Down     key.Binding
	PageUp   key.Binding
	PageDown key.Binding
}

// keyMapEntry represents a key mapping with its description
type keyMapEntry struct {
	key         string
	description string
}

// modeHelp represents help information for a mode
type modeHelp struct {
	name        string
	keymaps     []keyMapEntry
	hasDefault  bool
	defaultHelp string
	hasNumber   bool
	numberHelp  string
}

// HelpModel represents the pure state for the help interface
type HelpModel struct {
	// Display dimensions
	width  int
	height int

	// Visibility state
	visible bool

	// Content and viewport state
	viewport viewport.Model
	content  string

	// Configuration and mode management
	modesConfig *config.ModesConfig
	modeManager *ModeManager

	// Styles
	titleStyle       lipgloss.Style
	instructionStyle lipgloss.Style
	borderStyle      lipgloss.Style
}

// NewHelpModel creates a new help model
func NewHelpModel(modeManager *ModeManager) *HelpModel {
	vp := viewport.New(0, 0)

	// Initialize styles once
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Align(lipgloss.Center)

	instructionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(SecondaryTextColor)).
		Align(lipgloss.Center)

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		PaddingLeft(1).
		PaddingRight(1)

	return &HelpModel{
		viewport:         vp,
		visible:          false,
		modesConfig:      config.AppConfig.Modes,
		modeManager:      modeManager,
		titleStyle:       titleStyle,
		instructionStyle: instructionStyle,
		borderStyle:      borderStyle,
	}
}

// SetSize updates the help UI size
func (m *HelpModel) SetSize(width, height int) {
	m.width = width
	m.height = height

	// Account for border and padding
	innerWidth := width - 4   // border + padding
	innerHeight := height - 6 // border + padding + title

	if innerWidth < 0 {
		innerWidth = 0
	}
	if innerHeight < 0 {
		innerHeight = 0
	}

	m.viewport.Width = innerWidth
	m.viewport.Height = innerHeight
}

// GetSize returns the current dimensions
func (m *HelpModel) GetSize() (int, int) {
	return m.width, m.height
}

// Show displays the help UI
func (m *HelpModel) Show() {
	m.visible = true
	m.generateContent()
}

// IsVisible returns whether the help UI is visible
func (m *HelpModel) IsVisible() bool {
	return m.visible
}

// generateContent generates the help content for the current mode only
func (m *HelpModel) generateContent() {
	currentMode := m.modeManager.GetCurrentMode()
	modeConfig := m.modeManager.GetModeConfig(currentMode)

	if modeConfig == nil {
		m.content = "No configuration found for mode: " + currentMode
		m.viewport.SetContent(m.content)

		return
	}

	// Generate help for the current mode only
	modeHelpInfo := m.extractModeHelp(currentMode, &modeConfig.KeyBindings)

	// Create a simple header and format the help
	var sections []string
	sections = append(sections, m.formatModeHelp(modeHelpInfo))

	m.content = strings.Join(sections, "\n\n")
	m.viewport.SetContent(m.content)
}

// extractModeHelp extracts help information from a mode's key bindings
func (m *HelpModel) extractModeHelp(
	modeName string,
	keyBindings *config.KeyBindingsConfig,
) modeHelp {
	modeHelpInfo := modeHelp{
		name:    modeName,
		keymaps: []keyMapEntry{},
	}

	if keyBindings.OnKeys != nil {
		var keys []string
		for k := range keyBindings.OnKeys {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			actionConfig := keyBindings.OnKeys[k]
			if actionConfig.Help != "" {
				modeHelpInfo.keymaps = append(modeHelpInfo.keymaps, keyMapEntry{
					key:         k,
					description: actionConfig.Help,
				})
			}
		}
	}

	if keyBindings.Default != nil && keyBindings.Default.Help != "" {
		modeHelpInfo.hasDefault = true
		modeHelpInfo.defaultHelp = keyBindings.Default.Help
	}

	if keyBindings.OnNumber != nil && keyBindings.OnNumber.Help != "" {
		modeHelpInfo.hasNumber = true
		modeHelpInfo.numberHelp = keyBindings.OnNumber.Help
	}

	return modeHelpInfo
}

// formatModeHelp formats help information for a single mode
func (m *HelpModel) formatModeHelp(modeHelpInfo modeHelp) string {
	var lines []string

	// Key bindings
	if len(modeHelpInfo.keymaps) > 0 {
		for _, km := range modeHelpInfo.keymaps {
			line := km.key + " " + km.description
			lines = append(lines, line)
		}
	}

	// Default action
	if modeHelpInfo.hasDefault {
		line := "  " + "<other>" + " " + modeHelpInfo.defaultHelp
		lines = append(lines, line)
	}

	// Number action
	if modeHelpInfo.hasNumber {
		line := "  " + "<number>" + " " + modeHelpInfo.numberHelp
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

// Update handles help model updates and key events
func (m *HelpModel) Update(msg tea.Msg) {
	if !m.visible {
		return
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case msg.Type == tea.KeyEsc:
			m.visible = false
		case msg.String() == "k" || msg.String() == "up":
			m.viewport.ScrollUp(1)
		case msg.String() == "j" || msg.String() == "down":
			m.viewport.ScrollDown(1)
		case msg.String() == "ctrl+u":
			m.viewport.HalfPageUp()
		case msg.String() == "ctrl+d":
			m.viewport.HalfPageDown()
		}
	}
}

// View renders the help UI view
func (m *HelpModel) View() string {
	if !m.visible {
		return ""
	}

	width, height := m.GetSize()

	// Title
	title := m.titleStyle.Render("Help")

	// Help content with scrollable viewport
	content := m.viewport.View()

	// Instructions at the bottom
	instructions := m.instructionStyle.Render(
		"Press esc to close • ↑↓ to scroll",
	)

	// Combine all parts
	helpContent := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		content,
		instructions,
	)

	// Apply border and center on screen
	rendered := m.borderStyle.Render(helpContent)

	// Position the help UI in the center
	return lipgloss.Place(
		width, height,
		lipgloss.Center, lipgloss.Center,
		rendered,
		lipgloss.WithWhitespaceChars(""),
	)
}
