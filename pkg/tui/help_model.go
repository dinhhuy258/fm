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

	// Configuration
	modesConfig *config.ModesConfig

	// Key bindings
	keys HelpUIKeys
}

// NewHelpModel creates a new help model
func NewHelpModel() *HelpModel {
	vp := viewport.New(0, 0)

	return &HelpModel{
		viewport:    vp,
		visible:     false,
		modesConfig: config.AppConfig.Modes,
		keys:        DefaultHelpUIKeys(),
	}
}

// DefaultHelpUIKeys returns default key bindings
func DefaultHelpUIKeys() HelpUIKeys {
	return HelpUIKeys{
		Close: key.NewBinding(
			key.WithKeys("?", "esc", "q"),
			key.WithHelp("?/esc/q", "close help"),
		),
		Up: key.NewBinding(
			key.WithKeys("k", "up"),
			key.WithHelp("↑/k", "scroll up"),
		),
		Down: key.NewBinding(
			key.WithKeys("j", "down"),
			key.WithHelp("↓/j", "scroll down"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("pgup", "ctrl+u"),
			key.WithHelp("pgup", "page up"),
		),
		PageDown: key.NewBinding(
			key.WithKeys("pgdown", "ctrl+d"),
			key.WithHelp("pgdown", "page down"),
		),
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

// Hide hides the help UI
func (m *HelpModel) Hide() {
	m.visible = false
}

// Toggle toggles the help UI visibility
func (m *HelpModel) Toggle() {
	if m.visible {
		m.Hide()
	} else {
		m.Show()
	}
}

// IsVisible returns whether the help UI is visible
func (m *HelpModel) IsVisible() bool {
	return m.visible
}

// GetKeys returns the key bindings
func (m *HelpModel) GetKeys() HelpUIKeys {
	return m.keys
}

// GetViewport returns the viewport for direct manipulation
func (m *HelpModel) GetViewport() *viewport.Model {
	return &m.viewport
}

// UpdateViewport updates the viewport model
func (m *HelpModel) UpdateViewport(vp viewport.Model) {
	m.viewport = vp
}

// GetContent returns the generated content
func (m *HelpModel) GetContent() string {
	return m.content
}

// generateContent generates the help content from configuration
func (m *HelpModel) generateContent() {
	var sections []string

	// Generate help for builtin modes
	if len(m.modesConfig.Builtins) > 0 {
		sections = append(sections, "BUILTIN MODES")

		// Sort builtin modes for consistent display
		var builtinNames []string
		for name := range m.modesConfig.Builtins {
			builtinNames = append(builtinNames, name)
		}
		sort.Strings(builtinNames)

		for _, name := range builtinNames {
			modeConfig := m.modesConfig.Builtins[name]
			modeHelpInfo := m.extractModeHelp(name, &modeConfig.KeyBindings)
			sections = append(sections, m.formatModeHelp(modeHelpInfo))
		}
	}

	// Generate help for custom modes
	if len(m.modesConfig.Customs) > 0 {
		sections = append(sections, "CUSTOM MODES")

		// Sort custom modes for consistent display
		var customNames []string
		for name := range m.modesConfig.Customs {
			customNames = append(customNames, name)
		}
		sort.Strings(customNames)

		for _, name := range customNames {
			modeConfig := m.modesConfig.Customs[name]
			modeHelpInfo := m.extractModeHelp(name, &modeConfig.KeyBindings)
			sections = append(sections, m.formatModeHelp(modeHelpInfo))
		}
	}

	m.content = strings.Join(sections, "\n\n")
	m.viewport.SetContent(m.content)
}

// extractModeHelp extracts help information from a mode's key bindings
func (m *HelpModel) extractModeHelp(modeName string, keyBindings *config.KeyBindingsConfig) modeHelp {
	modeHelpInfo := modeHelp{
		name:    modeName,
		keymaps: []keyMapEntry{},
	}

	// Extract key bindings
	if keyBindings.OnKeys != nil {
		// Sort keys for consistent display
		var keys []string
		for k := range keyBindings.OnKeys {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			actionConfig := keyBindings.OnKeys[k]
			if actionConfig.Help != "" {
				// Get human-readable key display
				keyDisplay := getKeyDisplay(k)

				modeHelpInfo.keymaps = append(modeHelpInfo.keymaps, keyMapEntry{
					key:         keyDisplay,
					description: actionConfig.Help,
				})
			}
		}
	}

	// Extract default action help
	if keyBindings.Default != nil && keyBindings.Default.Help != "" {
		modeHelpInfo.hasDefault = true
		modeHelpInfo.defaultHelp = keyBindings.Default.Help
	}

	// Extract number action help
	if keyBindings.OnNumber != nil && keyBindings.OnNumber.Help != "" {
		modeHelpInfo.hasNumber = true
		modeHelpInfo.numberHelp = keyBindings.OnNumber.Help
	}

	return modeHelpInfo
}

// formatModeHelp formats help information for a single mode
func (m *HelpModel) formatModeHelp(modeHelpInfo modeHelp) string {
	var lines []string

	// Mode title
	lines = append(lines, "Mode: "+modeHelpInfo.name)

	// Key bindings
	if len(modeHelpInfo.keymaps) > 0 {
		for _, km := range modeHelpInfo.keymaps {
			line := "  " + km.key + " " + km.description
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

	// Add empty keymaps note if no keymaps
	if len(modeHelpInfo.keymaps) == 0 && !modeHelpInfo.hasDefault && !modeHelpInfo.hasNumber {
		lines = append(lines, "  (no key bindings)")
	}

	return strings.Join(lines, "\n")
}

// getKeyDisplay converts a key string to a human-readable display format
func getKeyDisplay(keyStr string) string {
	// Simple key display mapping for common keys
	keyDisplayMap := map[string]string{
		"ctrl+c":    "Ctrl+C",
		"ctrl+q":    "Ctrl+Q",
		"ctrl+h":    "Ctrl+H",
		"ctrl+l":    "Ctrl+L",
		"ctrl+u":    "Ctrl+U",
		"ctrl+d":    "Ctrl+D",
		"enter":     "Enter",
		"esc":       "Esc",
		"tab":       "Tab",
		"space":     "Space",
		"backspace": "Backspace",
		"delete":    "Delete",
		"up":        "↑",
		"down":      "↓",
		"left":      "←",
		"right":     "→",
		"home":      "Home",
		"end":       "End",
		"pgup":      "PgUp",
		"pgdown":    "PgDn",
	}

	if display, ok := keyDisplayMap[keyStr]; ok {
		return display
	}

	// For single characters and other keys, return as-is but capitalize
	if len(keyStr) == 1 {
		return strings.ToUpper(keyStr)
	}

	return keyStr
}

// Update handles help model updates and key events
func (m *HelpModel) Update(msg tea.Msg) (*HelpModel, tea.Cmd) {
	if !m.visible {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case msg.String() == "?" || msg.String() == "esc" || msg.String() == "q":
			m.Hide()

			return m, nil
		case msg.String() == "k" || msg.String() == "up":
			m.viewport.ScrollUp(1)
		case msg.String() == "j" || msg.String() == "down":
			m.viewport.ScrollDown(1)
		case msg.String() == "pgup" || msg.String() == "ctrl+u":
			m.viewport.HalfPageUp()
		case msg.String() == "pgdown" || msg.String() == "ctrl+d":
			m.viewport.HalfPageDown()
		}
	}

	return m, nil
}

// View renders the help UI view
func (m *HelpModel) View() string {
	if !m.visible {
		return ""
	}

	width, height := m.GetSize()

	// Title
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FFFF")).
		Bold(true).
		Align(lipgloss.Center).
		Padding(1, 2)
	title := titleStyle.Render("File Manager - Help")

	// Help content with scrollable viewport
	content := m.viewport.View()

	// Instructions at the bottom
	instructionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Align(lipgloss.Center)
	instructions := instructionStyle.Render("Press ? or ESC to close • ↑↓ to scroll • PgUp/PgDn for page navigation")

	// Combine all parts
	helpContent := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		content,
		instructions,
	)

	// Apply border and center on screen
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#00FFFF")).
		Padding(1)

	rendered := borderStyle.Render(helpContent)

	// Position the help UI in the center
	return lipgloss.Place(
		width, height,
		lipgloss.Center, lipgloss.Center,
		rendered,
		lipgloss.WithWhitespaceChars(""),
	)
}
