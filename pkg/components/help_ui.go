package components

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/dinhhuy258/fm/pkg/config"
)

// HelpUI provides a help interface showing all keymaps for all modes
type HelpUI struct {
	width    int
	height   int
	viewport viewport.Model
	visible  bool

	// Configuration
	modesConfig *config.ModesConfig

	// Styling
	titleStyle       lipgloss.Style
	headerStyle      lipgloss.Style
	keyStyle         lipgloss.Style
	descriptionStyle lipgloss.Style
	borderStyle      lipgloss.Style

	// Key bindings
	keys HelpUIKeys
}

// HelpUIKeys defines key bindings for the help UI
type HelpUIKeys struct {
	Close    key.Binding
	Up       key.Binding
	Down     key.Binding
	PageUp   key.Binding
	PageDown key.Binding
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

// NewHelpUI creates a new help UI
func NewHelpUI() *HelpUI {
	vp := viewport.New(0, 0)

	return &HelpUI{
		viewport:    vp,
		visible:     false,
		modesConfig: config.AppConfig.Modes,
		keys:        DefaultHelpUIKeys(),

		// Initialize styles
		titleStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FFFF")).
			Bold(true).
			Align(lipgloss.Center).
			Padding(1, 2),

		headerStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true).
			Underline(true).
			Margin(1, 0, 0, 0),

		keyStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFF00")).
			Bold(true).
			Width(15),

		descriptionStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")),

		borderStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#00FFFF")).
			Padding(1),
	}
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

// SetSize updates the help UI size
func (h *HelpUI) SetSize(width, height int) {
	h.width = width
	h.height = height

	// Account for border and padding
	innerWidth := width - 4   // border + padding
	innerHeight := height - 6 // border + padding + title

	if innerWidth < 0 {
		innerWidth = 0
	}
	if innerHeight < 0 {
		innerHeight = 0
	}

	h.viewport.Width = innerWidth
	h.viewport.Height = innerHeight
}

// Show displays the help UI
func (h *HelpUI) Show() {
	h.visible = true
	h.generateContent()
}

// Hide hides the help UI
func (h *HelpUI) Hide() {
	h.visible = false
}

// Toggle toggles the help UI visibility
func (h *HelpUI) Toggle() {
	if h.visible {
		h.Hide()
	} else {
		h.Show()
	}
}

// IsVisible returns whether the help UI is visible
func (h *HelpUI) IsVisible() bool {
	return h.visible
}

// Update handles input messages
func (h *HelpUI) Update(msg tea.Msg) (*HelpUI, tea.Cmd) {
	if !h.visible {
		return h, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, h.keys.Close):
			h.Hide()

			return h, nil
		case key.Matches(msg, h.keys.Up):
			h.viewport.ScrollUp(1)
		case key.Matches(msg, h.keys.Down):
			h.viewport.ScrollDown(1)
		case key.Matches(msg, h.keys.PageUp):
			h.viewport.HalfPageUp()
		case key.Matches(msg, h.keys.PageDown):
			h.viewport.HalfPageDown()
		}
	}

	var cmd tea.Cmd
	h.viewport, cmd = h.viewport.Update(msg)

	return h, cmd
}

// generateContent generates the help content from configuration
func (h *HelpUI) generateContent() {
	var sections []string

	// Generate help for builtin modes
	if len(h.modesConfig.Builtins) > 0 {
		sections = append(sections, h.headerStyle.Render("BUILTIN MODES"))

		// Sort builtin modes for consistent display
		var builtinNames []string
		for name := range h.modesConfig.Builtins {
			builtinNames = append(builtinNames, name)
		}
		sort.Strings(builtinNames)

		for _, name := range builtinNames {
			modeConfig := h.modesConfig.Builtins[name]
			modeHelpInfo := h.extractModeHelp(name, &modeConfig.KeyBindings)
			sections = append(sections, h.renderModeHelp(modeHelpInfo))
		}
	}

	// Generate help for custom modes
	if len(h.modesConfig.Customs) > 0 {
		sections = append(sections, h.headerStyle.Render("CUSTOM MODES"))

		// Sort custom modes for consistent display
		var customNames []string
		for name := range h.modesConfig.Customs {
			customNames = append(customNames, name)
		}
		sort.Strings(customNames)

		for _, name := range customNames {
			modeConfig := h.modesConfig.Customs[name]
			modeHelpInfo := h.extractModeHelp(name, &modeConfig.KeyBindings)
			sections = append(sections, h.renderModeHelp(modeHelpInfo))
		}
	}

	content := strings.Join(sections, "\n\n")
	h.viewport.SetContent(content)
}

// extractModeHelp extracts help information from a mode's key bindings
func (h *HelpUI) extractModeHelp(modeName string, keyBindings *config.KeyBindingsConfig) modeHelp {
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

// renderModeHelp renders help information for a single mode
func (h *HelpUI) renderModeHelp(modeHelpInfo modeHelp) string {
	var lines []string

	// Mode title
	modeTitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FFFF")).
		Bold(true).
		Render(fmt.Sprintf("Mode: %s", modeHelpInfo.name))
	lines = append(lines, modeTitle)

	// Key bindings
	if len(modeHelpInfo.keymaps) > 0 {
		for _, km := range modeHelpInfo.keymaps {
			keyPart := h.keyStyle.Render(km.key)
			descPart := h.descriptionStyle.Render(km.description)
			line := keyPart + " " + descPart
			lines = append(lines, "  "+line)
		}
	}

	// Default action
	if modeHelpInfo.hasDefault {
		keyPart := h.keyStyle.Render("<other>")
		descPart := h.descriptionStyle.Render(modeHelpInfo.defaultHelp)
		line := keyPart + " " + descPart
		lines = append(lines, "  "+line)
	}

	// Number action
	if modeHelpInfo.hasNumber {
		keyPart := h.keyStyle.Render("<number>")
		descPart := h.descriptionStyle.Render(modeHelpInfo.numberHelp)
		line := keyPart + " " + descPart
		lines = append(lines, "  "+line)
	}

	// Add empty keymaps note if no keymaps
	if len(modeHelpInfo.keymaps) == 0 && !modeHelpInfo.hasDefault && !modeHelpInfo.hasNumber {
		emptyNote := h.descriptionStyle.
			Foreground(lipgloss.Color("#888888")).
			Render("(no key bindings)")
		lines = append(lines, "  "+emptyNote)
	}

	return strings.Join(lines, "\n")
}

// View renders the help UI
func (h *HelpUI) View() string {
	if !h.visible {
		return ""
	}

	// Title
	title := h.titleStyle.Render("File Manager - Help")

	// Help content with scrollable viewport
	content := h.viewport.View()

	// Instructions at the bottom
	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Align(lipgloss.Center).
		Render("Press ? or ESC to close • ↑↓ to scroll • PgUp/PgDn for page navigation")

	// Combine all parts
	helpContent := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		content,
		instructions,
	)

	// Apply border and center on screen
	return h.borderStyle.Render(helpContent)
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
