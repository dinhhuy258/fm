package components

import (
	"strconv"
	"strings"
	"unicode/utf8"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rivo/uniseg"

	set "github.com/deckarep/golang-set/v2"
	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/fm/pkg/fs"
)

// nodeType represents an icon with its style
type nodeType struct {
	icon  string
	style lipgloss.Style
}

// nodeTypes matches the GoCUI version structure
type nodeTypes struct {
	file             nodeType
	directory        nodeType
	fileSymlink      nodeType
	directorySymlink nodeType
	extensions       map[string]nodeType
	specials         map[string]nodeType
}

// headerStyles contains styles for each column header
type headerStyles struct {
	indexHeader lipgloss.Style
	nameHeader  lipgloss.Style
}

// ExplorerTable provides exact feature parity with GoCUI ExplorerView
// This is a standalone component that doesn't use bubbles/list
type ExplorerTable struct {
	width       int
	height      int
	entries     []fs.IEntry
	selections  set.Set[string]
	focus       int
	scrollStart int
	showHeader  bool
	currentPath string

	// Configuration
	explorerConfig  *config.ExplorerTableConfig
	nodeTypesConfig *config.NodeTypesConfig

	// Computed styles from config
	defaultFileStyle      lipgloss.Style
	defaultDirectoryStyle lipgloss.Style
	focusStyle            lipgloss.Style
	selectionStyle        lipgloss.Style
	focusSelectionStyle   lipgloss.Style

	// Header styles
	headerStyles headerStyles

	// Icons mapping
	icons nodeTypes
}

// NewExplorerTable creates a new explorer table
func NewExplorerTable() *ExplorerTable {
	explorerConfig := config.AppConfig.General.ExplorerTable
	nodeTypesConfig := config.AppConfig.NodeTypes

	table := &ExplorerTable{
		selections:      set.NewSet[string](),
		showHeader:      true,
		explorerConfig:  explorerConfig,
		nodeTypesConfig: nodeTypesConfig,
		focus:           0,
		scrollStart:     0,
	}

	// Initialize styles and icons
	table.initStyles()
	table.initIcons()

	return table
}

// initStyles initializes lipgloss styles from config
func (t *ExplorerTable) initStyles() {
	// Default styles with good visibility
	t.defaultFileStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
	t.defaultDirectoryStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF")).Bold(true)

	// Override with config if available
	if t.explorerConfig.DefaultUI.FileStyle != nil {
		t.defaultFileStyle = convertConfigToLipgloss(t.explorerConfig.DefaultUI.FileStyle)
	}
	if t.explorerConfig.DefaultUI.DirectoryStyle != nil {
		t.defaultDirectoryStyle = convertConfigToLipgloss(t.explorerConfig.DefaultUI.DirectoryStyle)
	}

	// Enhanced focus styling with background highlight and visible text
	t.focusStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("#444444")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true)

	// Enhanced selection styling with green background
	t.selectionStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("#004400")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true)

	// Enhanced focus+selection styling
	t.focusSelectionStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("#006600")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true)

	// Header styles - using default if not specified
	defaultHeaderStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00FFFF"))
	t.headerStyles = headerStyles{
		indexHeader: getStyleOrDefault(t.explorerConfig.IndexHeader.Style, defaultHeaderStyle),
		nameHeader:  getStyleOrDefault(t.explorerConfig.NameHeader.Style, defaultHeaderStyle),
	}
}

// initIcons initializes the icon system from config
func (t *ExplorerTable) initIcons() {
	t.icons = nodeTypes{
		file: nodeType{
			icon:  t.nodeTypesConfig.File.Icon,
			style: convertConfigToLipgloss(t.nodeTypesConfig.File.Style),
		},
		directory: nodeType{
			icon:  t.nodeTypesConfig.Directory.Icon,
			style: convertConfigToLipgloss(t.nodeTypesConfig.Directory.Style),
		},
		fileSymlink: nodeType{
			icon:  t.nodeTypesConfig.FileSymlink.Icon,
			style: convertConfigToLipgloss(t.nodeTypesConfig.FileSymlink.Style),
		},
		directorySymlink: nodeType{
			icon:  t.nodeTypesConfig.DirectorySymlink.Icon,
			style: convertConfigToLipgloss(t.nodeTypesConfig.DirectorySymlink.Style),
		},
		extensions: make(map[string]nodeType),
		specials:   make(map[string]nodeType),
	}

	// Load extension-specific icons
	for ext, ntc := range t.nodeTypesConfig.Extensions {
		t.icons.extensions[strings.ToLower(ext)] = nodeType{
			icon:  ntc.Icon,
			style: convertConfigToLipgloss(ntc.Style),
		}
	}

	// Load special file icons
	for fileName, ntc := range t.nodeTypesConfig.Specials {
		t.icons.specials[strings.ToLower(fileName)] = nodeType{
			icon:  ntc.Icon,
			style: convertConfigToLipgloss(ntc.Style),
		}
	}
}

// convertConfigToLipgloss converts config.StyleConfig to lipgloss.Style
func convertConfigToLipgloss(styleConfig *config.StyleConfig) lipgloss.Style {
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
		case "reverse":
			// For reverse, we swap foreground and background
			if styleConfig.Fg != "" && styleConfig.Bg != "" {
				fgColor := styleConfig.Fg
				bgColor := styleConfig.Bg
				if mappedFg, ok := colorMap[fgColor]; ok {
					fgColor = mappedFg
				}
				if mappedBg, ok := colorMap[bgColor]; ok {
					bgColor = mappedBg
				}
				style = style.Foreground(lipgloss.Color(bgColor))
				style = style.Background(lipgloss.Color(fgColor))
			}
		}
	}

	return style
}

// getStyleOrDefault returns a converted style or default if config is nil
func getStyleOrDefault(styleConfig *config.StyleConfig, defaultStyle lipgloss.Style) lipgloss.Style {
	if styleConfig != nil {
		return convertConfigToLipgloss(styleConfig)
	}

	return defaultStyle
}

// SetSize updates the table size
func (t *ExplorerTable) SetSize(width, height int) {
	t.width = width
	t.height = height
}

// SetEntries updates the entries and resets focus/selection
func (t *ExplorerTable) SetEntries(entries []fs.IEntry, currentPath string) {
	t.entries = entries
	t.currentPath = currentPath
	t.focus = 0
	t.scrollStart = 0
	t.selections = set.NewSet[string]()
}

// Update handles input messages
func (t *ExplorerTable) Update(msg tea.Msg) (*ExplorerTable, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			t.moveCursor(-1)
		case tea.KeyDown:
			t.moveCursor(1)
		case tea.KeyPgUp:
			t.moveCursor(-t.getVisibleRows())
		case tea.KeyPgDown:
			t.moveCursor(t.getVisibleRows())
		case tea.KeyHome:
			t.focus = 0
			t.scrollStart = 0
		case tea.KeyEnd:
			if len(t.entries) > 0 {
				t.focus = len(t.entries) - 1
				t.ensureVisible()
			}
		case tea.KeySpace:
			t.toggleSelection()
		}
	}

	return t, nil
}

// moveCursor moves the cursor by delta positions
func (t *ExplorerTable) moveCursor(delta int) {
	if len(t.entries) == 0 {
		return
	}

	newFocus := t.focus + delta
	if newFocus < 0 {
		newFocus = 0
	} else if newFocus >= len(t.entries) {
		newFocus = len(t.entries) - 1
	}

	t.focus = newFocus
	t.ensureVisible()
}

// ensureVisible ensures the focused item is visible by adjusting scroll
func (t *ExplorerTable) ensureVisible() {
	visibleRows := t.getVisibleRows()

	if t.focus < t.scrollStart {
		t.scrollStart = t.focus
	} else if t.focus >= t.scrollStart+visibleRows {
		t.scrollStart = t.focus - visibleRows + 1
	}

	if t.scrollStart < 0 {
		t.scrollStart = 0
	}
}

// getVisibleRows calculates how many rows can fit in the current height
func (t *ExplorerTable) getVisibleRows() int {
	available := t.height
	if t.showHeader {
		available--
	}
	if available < 1 {
		available = 1
	}

	return available
}

// toggleSelection toggles selection for the focused item
func (t *ExplorerTable) toggleSelection() {
	if t.focus < len(t.entries) {
		entry := t.entries[t.focus]
		path := entry.GetPath()

		if t.selections.Contains(path) {
			t.selections.Remove(path)
		} else {
			t.selections.Add(path)
		}
	}
}

// View renders the table
func (t *ExplorerTable) View() string {
	if t.width <= 0 || t.height <= 0 {
		return ""
	}

	var sections []string

	// Render header if enabled
	if t.showHeader {
		sections = append(sections, t.renderHeader())
	}

	// Render visible entries
	sections = append(sections, t.renderEntries())

	return strings.Join(sections, "\n")
}

// renderHeader renders the column headers
func (t *ExplorerTable) renderHeader() string {
	columns := []columnConfig{
		{percentage: t.explorerConfig.IndexHeader.Percentage, leftAlign: true},
		{percentage: t.explorerConfig.NameHeader.Percentage, leftAlign: true},
	}

	values := []styledValue{
		{text: t.explorerConfig.IndexHeader.Name, style: t.headerStyles.indexHeader},
		{text: t.explorerConfig.NameHeader.Name, style: t.headerStyles.nameHeader},
	}

	return t.formatRow(columns, values)
}

// renderEntries renders the visible file entries
func (t *ExplorerTable) renderEntries() string {
	if len(t.entries) == 0 {
		return ""
	}

	visibleRows := t.getVisibleRows()
	endIndex := min(t.scrollStart+visibleRows, len(t.entries))

	lines := make([]string, 0, endIndex-t.scrollStart)

	for i := t.scrollStart; i < endIndex; i++ {
		lines = append(lines, t.renderEntry(t.entries[i], i))
	}

	// Pad with empty lines if needed
	for len(lines) < visibleRows {
		lines = append(lines, strings.Repeat(" ", t.width))
	}

	return strings.Join(lines, "\n")
}

// renderEntry renders a single file entry with proper styling
func (t *ExplorerTable) renderEntry(entry fs.IEntry, idx int) string {
	isEntryFocused := idx == t.focus
	isEntrySelected := t.selections.Contains(entry.GetPath())

	// Get the appropriate icon with state-based styling
	entryIcon := t.getEntryIcon(entry, isEntryFocused, isEntrySelected)

	// Determine prefix, suffix, and style based on state
	var prefix, suffix, entryTreePrefix string
	var entryStyle lipgloss.Style

	switch {
	case isEntryFocused && isEntrySelected:
		prefix = t.explorerConfig.FocusSelectionUI.Prefix
		suffix = t.explorerConfig.FocusSelectionUI.Suffix
		entryStyle = t.focusSelectionStyle
	case isEntryFocused:
		prefix = t.explorerConfig.FocusUI.Prefix
		suffix = t.explorerConfig.FocusUI.Suffix
		entryStyle = t.focusStyle
	case isEntrySelected:
		prefix = t.explorerConfig.SelectionUI.Prefix
		suffix = t.explorerConfig.SelectionUI.Suffix
		entryStyle = t.selectionStyle
	default:
		prefix = t.explorerConfig.DefaultUI.Prefix
		suffix = t.explorerConfig.DefaultUI.Suffix
		if entry.IsDirectory() {
			entryStyle = t.defaultDirectoryStyle
		} else {
			entryStyle = t.defaultFileStyle
		}
	}

	// Tree prefix based on position
	entriesSize := len(t.entries)
	switch idx {
	case entriesSize - 1:
		entryTreePrefix = t.explorerConfig.LastEntryPrefix
	case 0:
		entryTreePrefix = t.explorerConfig.FirstEntryPrefix
	default:
		entryTreePrefix = t.explorerConfig.EntryPrefix
	}

	// Get icon and apply state-specific styling if needed
	iconText := entryIcon.icon
	if iconText == "" {
		if entry.IsDirectory() {
			iconText = "📁"
		} else {
			iconText = "📄"
		}
	}

	// Build the complete name string without applying styling yet
	fileName := strings.TrimSpace(entry.GetName())

	// Apply styling to just the icon if needed (but keep it simple)
	var styledIcon string
	if isEntryFocused || isEntrySelected {
		// For focused/selected items, apply same style to icon as text
		styledIcon = iconText
	} else {
		// For normal items, use icon's default style
		styledIcon = entryIcon.style.Render(iconText)
	}

	// Build the complete name column - don't apply styling here
	nameColumn := entryTreePrefix + prefix + styledIcon + " " + fileName + suffix

	// Column configuration - only two columns now
	columns := []columnConfig{
		{percentage: t.explorerConfig.IndexHeader.Percentage, leftAlign: true},
		{percentage: t.explorerConfig.NameHeader.Percentage, leftAlign: true},
	}

	// Column values - apply styling to the entire nameColumn to ensure consistent background
	indexStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")) // Dimmed index
	values := []styledValue{
		{text: strconv.Itoa(idx + 1), style: indexStyle},
		{text: nameColumn, style: entryStyle}, // Apply entry style to complete name column
	}

	return t.formatRow(columns, values)
}

// columnConfig represents column configuration
type columnConfig struct {
	percentage int
	leftAlign  bool
}

// styledValue represents a value with its style
type styledValue struct {
	text  string
	style lipgloss.Style
}

// formatRow formats a row with proper column alignment and styling
func (t *ExplorerTable) formatRow(columns []columnConfig, values []styledValue) string {
	if len(columns) != len(values) {
		return "Invalid row configuration"
	}

	result := ""
	consumedWidth := 0
	for i, col := range columns {
		width := int(float32(col.percentage) / 100.0 * float32(t.width))
		if i == len(columns)-1 {
			width = t.width - consumedWidth
		} else {
			consumedWidth += width
		}
		result += t.formatColumn(values[i], width, col.leftAlign)
	}

	// Ensure the row doesn't exceed terminal width
	if uniseg.StringWidth(result) > t.width {
		runes := []rune(result)
		if len(runes) > t.width {
			result = string(runes[:t.width])
		}
	}

	return result
}

// formatColumn formats a single column with proper alignment
func (t *ExplorerTable) formatColumn(value styledValue, width int, leftAlign bool) string {
	if width <= 0 {
		return ""
	}

	text := value.text
	displayWidth := uniseg.StringWidth(text)

	// Truncate if too long
	if displayWidth > width {
		truncated := ""
		currentWidth := 0
		for _, r := range text {
			runeWidth := utf8.RuneLen(r)
			if runeWidth < 0 {
				runeWidth = 1 // fallback for invalid runes
			}
			if currentWidth+runeWidth > width {
				break
			}
			truncated += string(r)
			currentWidth += runeWidth
		}
		text = truncated
		displayWidth = currentWidth
	}

	// Calculate padding
	padding := max(width-displayWidth, 0)

	// Apply styling to the text content
	styledText := value.style.Render(text)

	// Add padding
	if leftAlign {
		return styledText + strings.Repeat(" ", padding)
	} else {
		return strings.Repeat(" ", padding) + styledText
	}
}

// getEntryIcon returns the appropriate icon for an entry with state-based styling
func (t *ExplorerTable) getEntryIcon(entry fs.IEntry, isEntryFocused, isEntrySelected bool) nodeType {
	var icon nodeType

	// Find the appropriate icon based on file type
	extensionIcon, hasExtIcon := t.icons.extensions[strings.ToLower(entry.GetExt())]
	specialIcon, hasSpecialIcon := t.icons.specials[strings.ToLower(entry.GetName())]

	switch {
	case entry.IsSymlink() && entry.IsDirectory():
		icon = t.icons.directorySymlink
	case entry.IsSymlink():
		icon = t.icons.fileSymlink
	case hasSpecialIcon:
		icon = specialIcon
	case hasExtIcon:
		icon = extensionIcon
	case entry.IsDirectory():
		icon = t.icons.directory
		// Use a visible folder icon if none is configured
		if icon.icon == "" {
			icon.icon = "📁"
		}
	default:
		icon = t.icons.file
		// Use a visible file icon if none is configured
		if icon.icon == "" {
			icon.icon = "📄"
		}
	}

	// Apply state-based styling that inherits from base icon style
	baseStyle := icon.style
	switch {
	case isEntrySelected && isEntryFocused:
		icon.style = baseStyle.
			Background(lipgloss.Color("#006600")).
			Bold(true)
	case isEntrySelected:
		icon.style = baseStyle.
			Background(lipgloss.Color("#004400")).
			Bold(true)
	case isEntryFocused:
		icon.style = baseStyle.
			Background(lipgloss.Color("#444444")).
			Bold(true)
		// Default case - use base icon style
	}

	return icon
}

// Public interface methods

// GetFocusedEntry returns the currently focused entry
func (t *ExplorerTable) GetFocusedEntry() fs.IEntry {
	if t.focus < len(t.entries) {
		return t.entries[t.focus]
	}

	return nil
}

// GetSelectedEntries returns all selected entries
func (t *ExplorerTable) GetSelectedEntries() []fs.IEntry {
	var selected []fs.IEntry
	for _, entry := range t.entries {
		if t.selections.Contains(entry.GetPath()) {
			selected = append(selected, entry)
		}
	}

	return selected
}

// GetStats returns total and selected entry counts
func (t *ExplorerTable) GetStats() (total, selected int) {
	return len(t.entries), t.selections.Cardinality()
}

// GetFocus returns the current focus index
func (t *ExplorerTable) GetFocus() int {
	return t.focus
}

// IsSelected returns true if the entry at the given index is selected
func (t *ExplorerTable) IsSelected(index int) bool {
	if index < len(t.entries) {
		return t.selections.Contains(t.entries[index].GetPath())
	}

	return false
}

// SetShowHeader sets whether to show the column header
func (t *ExplorerTable) SetShowHeader(show bool) {
	t.showHeader = show
}

// SetFocusByIndex sets the focus to a specific index
func (t *ExplorerTable) SetFocusByIndex(index int) bool {
	if index < 0 || index >= len(t.entries) {
		return false
	}

	t.focus = index
	t.ensureVisible()

	return true
}

// ToggleSelectionByPath toggles selection for an entry with the given path
func (t *ExplorerTable) ToggleSelectionByPath(path string) bool {
	for _, entry := range t.entries {
		if entry.GetPath() == path {
			if t.selections.Contains(path) {
				t.selections.Remove(path)
			} else {
				t.selections.Add(path)
			}

			return true
		}
	}

	return false
}

// FocusPath attempts to focus on an entry with the given path
func (t *ExplorerTable) FocusPath(path string) bool {
	for i, entry := range t.entries {
		if entry.GetPath() == path {
			t.focus = i
			t.ensureVisible()

			return true
		}
	}

	return false
}
