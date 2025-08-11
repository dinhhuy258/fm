package components

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/dinhhuy258/fm/pkg/config"
)

// LogLevel represents different types of log messages
type LogLevel int8

const (
	LogInfo LogLevel = iota
	LogWarning
	LogError
)

// LogEntry represents a single log message with level and timestamp
type LogEntry struct {
	Level     LogLevel
	Message   string
	Timestamp time.Time
}

// LogView provides a scrollable log view with different message types
type LogView struct {
	width  int
	height int

	// Log entries storage
	entries    []LogEntry
	maxEntries int // Maximum number of entries to keep

	// Scrolling state
	scrollOffset int
	autoScroll   bool // Whether to auto-scroll to bottom on new entries

	// Styling based on config
	infoStyle    lipgloss.Style
	warningStyle lipgloss.Style
	errorStyle   lipgloss.Style
	titleStyle   lipgloss.Style

	// Configuration
	config *config.Config
}

// NewLogView creates a new log view
func NewLogView() *LogView {
	cfg := config.AppConfig

	logView := &LogView{
		entries:      make([]LogEntry, 0),
		maxEntries:   1000, // Keep last 1000 entries
		autoScroll:   true,
		scrollOffset: 0,
		config:       cfg,
	}

	logView.initStyles()

	return logView
}

// initStyles initializes lipgloss styles from config
func (lv *LogView) initStyles() {
	// Default styles with good visibility
	lv.infoStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(false)

	lv.warningStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFF00")).
		Bold(true)

	lv.errorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF0000")).
		Bold(true)

	lv.titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00FFFF")).
		Background(lipgloss.Color("#444444")).
		Padding(0, 1)

	// Override with config if available
	if lv.config != nil && lv.config.General.LogInfoUI.Style != nil {
		lv.infoStyle = convertConfigStyleToLipgloss(lv.config.General.LogInfoUI.Style)
	}
	if lv.config != nil && lv.config.General.LogWarningUI.Style != nil {
		lv.warningStyle = convertConfigStyleToLipgloss(lv.config.General.LogWarningUI.Style)
	}
	if lv.config != nil && lv.config.General.LogErrorUI.Style != nil {
		lv.errorStyle = convertConfigStyleToLipgloss(lv.config.General.LogErrorUI.Style)
	}
}

// convertConfigStyleToLipgloss converts config.StyleConfig to lipgloss.Style
func convertConfigStyleToLipgloss(styleConfig *config.StyleConfig) lipgloss.Style {
	style := lipgloss.NewStyle()

	if styleConfig == nil {
		return style
	}

	// Map common color names to better alternatives
	colorMap := map[string]string{
		"white":   "#FFFFFF",
		"cyan":    "#00FFFF",
		"green":   "#00FF00",
		"yellow":  "#FFFF00",
		"red":     "#FF0000",
		"blue":    "#0000FF",
		"magenta": "#FF00FF",
		"black":   "#000000",
	}

	// Set foreground color
	if styleConfig.Fg != "" {
		color := styleConfig.Fg
		if mappedColor, ok := colorMap[color]; ok {
			color = mappedColor
		}
		style = style.Foreground(lipgloss.Color(color))
	}

	// Set background color
	if styleConfig.Bg != "" {
		color := styleConfig.Bg
		if mappedColor, ok := colorMap[color]; ok {
			color = mappedColor
		}
		style = style.Background(lipgloss.Color(color))
	}

	// Apply decorations
	for _, decoration := range styleConfig.Decorations {
		switch decoration {
		case "bold":
			style = style.Bold(true)
		case "italic":
			style = style.Italic(true)
		case "underline", "underscore":
			style = style.Underline(true)
		}
	}

	return style
}

// SetSize updates the log view size
func (lv *LogView) SetSize(width, height int) {
	lv.width = width
	lv.height = height
}

// AddLog adds a new log entry
func (lv *LogView) AddLog(level LogLevel, message string) {
	entry := LogEntry{
		Level:     level,
		Message:   message,
		Timestamp: time.Now(),
	}

	lv.entries = append(lv.entries, entry)

	// Trim entries if we exceed max
	if len(lv.entries) > lv.maxEntries {
		lv.entries = lv.entries[1:]
	}

	// Auto-scroll to bottom if enabled
	if lv.autoScroll {
		lv.scrollToBottom()
	}
}

// AddInfo adds an info log entry
func (lv *LogView) AddInfo(message string) {
	lv.AddLog(LogInfo, message)
}

// AddWarning adds a warning log entry
func (lv *LogView) AddWarning(message string) {
	lv.AddLog(LogWarning, message)
}

// AddError adds an error log entry
func (lv *LogView) AddError(message string) {
	lv.AddLog(LogError, message)
}

// Clear clears all log entries
func (lv *LogView) Clear() {
	lv.entries = make([]LogEntry, 0)
	lv.scrollOffset = 0
}

// ScrollUp scrolls the log view up
func (lv *LogView) ScrollUp() {
	if lv.scrollOffset > 0 {
		lv.scrollOffset--
		lv.autoScroll = false
	}
}

// ScrollDown scrolls the log view down
func (lv *LogView) ScrollDown() {
	maxScroll := lv.getMaxScrollOffset()
	if lv.scrollOffset < maxScroll {
		lv.scrollOffset++
	}

	// Re-enable auto-scroll if we're at the bottom
	if lv.scrollOffset >= maxScroll {
		lv.autoScroll = true
	}
}

// ScrollToTop scrolls to the top of the log
func (lv *LogView) ScrollToTop() {
	lv.scrollOffset = 0
	lv.autoScroll = false
}

// scrollToBottom scrolls to the bottom of the log
func (lv *LogView) scrollToBottom() {
	lv.scrollOffset = lv.getMaxScrollOffset()
	lv.autoScroll = true
}

// getMaxScrollOffset returns the maximum scroll offset
func (lv *LogView) getMaxScrollOffset() int {
	if lv.height <= 2 { // Account for title and borders
		return 0
	}

	visibleLines := lv.height - 2 // Account for title
	totalLines := len(lv.entries)

	maxScroll := totalLines - visibleLines
	if maxScroll < 0 {
		maxScroll = 0
	}

	return maxScroll
}

// Update handles input messages
func (lv *LogView) Update(msg tea.Msg) (*LogView, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			lv.ScrollUp()
		case "down", "j":
			lv.ScrollDown()
		case "home", "g":
			lv.ScrollToTop()
		case "end", "G":
			lv.scrollToBottom()
		case "pgup":
			// Page up
			for i := 0; i < lv.height-2 && lv.scrollOffset > 0; i++ {
				lv.scrollOffset--
			}
			lv.autoScroll = false
		case "pgdown":
			// Page down
			maxScroll := lv.getMaxScrollOffset()
			for i := 0; i < lv.height-2 && lv.scrollOffset < maxScroll; i++ {
				lv.scrollOffset++
			}
			// Re-enable auto-scroll if we're at the bottom
			if lv.scrollOffset >= maxScroll {
				lv.autoScroll = true
			}
		}
	}

	return lv, nil
}

// View renders the log view
func (lv *LogView) View() string {
	if lv.width <= 0 || lv.height <= 0 {
		return ""
	}

	var sections []string

	// Title bar
	title := lv.titleStyle.Render(" Logs ")
	sections = append(sections, title)

	// If no entries, show empty message
	if len(lv.entries) == 0 {
		emptyMsg := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Render("No log messages")
		sections = append(sections, emptyMsg)

		// Pad with empty lines
		for len(sections) < lv.height {
			sections = append(sections, "")
		}

		return strings.Join(sections, "\n")
	}

	// Render visible log entries
	sections = append(sections, lv.renderEntries())

	return strings.Join(sections, "\n")
}

// renderEntries renders the visible log entries
func (lv *LogView) renderEntries() string {
	if len(lv.entries) == 0 {
		return ""
	}

	visibleLines := lv.height - 1 // Account for title
	if visibleLines <= 0 {
		return ""
	}

	startIndex := lv.scrollOffset
	endIndex := startIndex + visibleLines
	if endIndex > len(lv.entries) {
		endIndex = len(lv.entries)
	}

	lines := make([]string, 0, endIndex-startIndex)

	for i := startIndex; i < endIndex; i++ {
		lines = append(lines, lv.renderEntry(lv.entries[i]))
	}

	// Pad with empty lines if needed
	for len(lines) < visibleLines {
		lines = append(lines, "")
	}

	return strings.Join(lines, "\n")
}

// renderEntry renders a single log entry
func (lv *LogView) renderEntry(entry LogEntry) string {
	var style lipgloss.Style
	var prefix, suffix string

	// Get style and prefix/suffix based on log level
	switch entry.Level {
	case LogInfo:
		style = lv.infoStyle
		if lv.config != nil {
			prefix = lv.config.General.LogInfoUI.Prefix
			suffix = lv.config.General.LogInfoUI.Suffix
		}
	case LogWarning:
		style = lv.warningStyle
		if lv.config != nil {
			prefix = lv.config.General.LogWarningUI.Prefix
			suffix = lv.config.General.LogWarningUI.Suffix
		}
	case LogError:
		style = lv.errorStyle
		if lv.config != nil {
			prefix = lv.config.General.LogErrorUI.Prefix
			suffix = lv.config.General.LogErrorUI.Suffix
		}
	}

	// Format message with prefix and suffix
	message := prefix + entry.Message + suffix

	// Apply styling and truncate if needed
	if len(message) > lv.width {
		message = message[:lv.width-3] + "..."
	}

	return style.Render(message)
}

// GetEntryCount returns the total number of log entries
func (lv *LogView) GetEntryCount() int {
	return len(lv.entries)
}

// IsAutoScrollEnabled returns whether auto-scroll is enabled
func (lv *LogView) IsAutoScrollEnabled() bool {
	return lv.autoScroll
}

// SetAutoScroll enables or disables auto-scroll
func (lv *LogView) SetAutoScroll(enabled bool) {
	lv.autoScroll = enabled
	if enabled {
		lv.scrollToBottom()
	}
}
