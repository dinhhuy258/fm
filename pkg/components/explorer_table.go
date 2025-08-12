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

// nodeTypes contains all node types with their icons and styles
type nodeTypes struct {
	// Default icons and styles for common file types
	file nodeType
	// Default directory icon and style
	directory nodeType
	// Symlink icons and styles
	fileSymlink nodeType
	// Directory symlink icon and style
	directorySymlink nodeType
	// Maps for extensions and special files
	extensions map[string]nodeType
	// Special files with custom icons
	specials map[string]nodeType
}

// headerStyles contains styles for each column header
type headerStyles struct {
	indexHeader lipgloss.Style
	nameHeader  lipgloss.Style
}

// ExplorerTable represents a table for file explorer
type ExplorerTable struct {
	width  int
	height int
	// Entries (files/directories) to display in the table
	entries []fs.IEntry
	// Selections is a set of selected file paths
	selections set.Set[string]
	// Focused index
	focus       int
	scrollStart int
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
	t.defaultFileStyle = FromStyleConfig(t.explorerConfig.DefaultUI.FileStyle)
	t.defaultDirectoryStyle = FromStyleConfig(t.explorerConfig.DefaultUI.DirectoryStyle)
	t.focusStyle = FromStyleConfig(t.explorerConfig.FocusUI.Style)
	t.selectionStyle = FromStyleConfig(t.explorerConfig.SelectionUI.Style)
	t.focusSelectionStyle = FromStyleConfig(t.explorerConfig.FocusSelectionUI.Style)
	t.headerStyles = headerStyles{
		indexHeader: FromStyleConfig(t.explorerConfig.IndexHeader.Style),
		nameHeader:  FromStyleConfig(t.explorerConfig.NameHeader.Style),
	}
}

// initIcons initializes the icon system from config
func (t *ExplorerTable) initIcons() {
	t.icons = nodeTypes{
		file: nodeType{
			icon:  t.nodeTypesConfig.File.Icon,
			style: FromStyleConfig(t.nodeTypesConfig.File.Style),
		},
		directory: nodeType{
			icon:  t.nodeTypesConfig.Directory.Icon,
			style: FromStyleConfig(t.nodeTypesConfig.Directory.Style),
		},
		fileSymlink: nodeType{
			icon:  t.nodeTypesConfig.FileSymlink.Icon,
			style: FromStyleConfig(t.nodeTypesConfig.FileSymlink.Style),
		},
		directorySymlink: nodeType{
			icon:  t.nodeTypesConfig.DirectorySymlink.Icon,
			style: FromStyleConfig(t.nodeTypesConfig.DirectorySymlink.Style),
		},
		extensions: make(map[string]nodeType),
		specials:   make(map[string]nodeType),
	}

	// Load extension-specific icons
	for ext, ntc := range t.nodeTypesConfig.Extensions {
		t.icons.extensions[strings.ToLower(ext)] = nodeType{
			icon:  ntc.Icon,
			style: FromStyleConfig(ntc.Style),
		}
	}

	// Load special file icons
	for fileName, ntc := range t.nodeTypesConfig.Specials {
		t.icons.specials[strings.ToLower(fileName)] = nodeType{
			icon:  ntc.Icon,
			style: FromStyleConfig(ntc.Style),
		}
	}
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

	t.focus = t.focus + delta
	t.focus = max(0, min(t.focus, len(t.entries)-1))
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

	t.scrollStart = max(0, t.scrollStart)
}

// getVisibleRows calculates how many rows can fit in the current height
func (t *ExplorerTable) getVisibleRows() int {
	// Reserve one row for the header
	return max(t.height-1, 1)
}

// toggleSelection toggles selection for the focused item.
func (t *ExplorerTable) toggleSelection() {
	if t.focus >= len(t.entries) {
		return
	}

	entry := t.entries[t.focus]
	path := entry.GetPath()

	t.ToggleSelectionByPath(path)
}

// View renders the table
func (t *ExplorerTable) View() string {
	if t.width <= 0 || t.height <= 0 {
		return ""
	}

	var sections []string

	// Render header
	sections = append(sections, t.renderHeader())

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
	lastVisibleIndex := min(t.scrollStart+visibleRows, len(t.entries))

	lines := make([]string, 0, lastVisibleIndex-t.scrollStart)

	for i := t.scrollStart; i < lastVisibleIndex; i++ {
		lines = append(lines, t.renderEntry(t.entries[i], i))
	}

	// Pad with empty lines if needed
	for len(lines) < visibleRows {
		lines = append(lines, strings.Repeat(" ", t.width))
	}

	return strings.Join(lines, "\n")
}

// entryDisplayState holds the UI state for rendering an entry
type entryDisplayState struct {
	prefix     string
	suffix     string
	treePrefix string
	style      lipgloss.Style
	isFocused  bool
	isSelected bool
}

// renderEntry renders a single file entry with proper styling
func (t *ExplorerTable) renderEntry(entry fs.IEntry, idx int) string {
	state := t.determineEntryDisplayState(entry, idx)
	entryIcon := t.getEntryIcon(entry, state.isFocused, state.isSelected)
	nameColumn := t.buildEntryDisplayName(entry, entryIcon, state)

	return t.formatEntryRow(idx, nameColumn, state.style)
}

// determineEntryDisplayState calculates the display state for an entry based on focus/selection
func (t *ExplorerTable) determineEntryDisplayState(entry fs.IEntry, idx int) entryDisplayState {
	isFocused := idx == t.focus
	isSelected := t.selections.Contains(entry.GetPath())

	var prefix, suffix string
	var style lipgloss.Style

	switch {
	case isFocused && isSelected:
		prefix = t.explorerConfig.FocusSelectionUI.Prefix
		suffix = t.explorerConfig.FocusSelectionUI.Suffix
		style = t.focusSelectionStyle
	case isFocused:
		prefix = t.explorerConfig.FocusUI.Prefix
		suffix = t.explorerConfig.FocusUI.Suffix
		style = t.focusStyle
	case isSelected:
		prefix = t.explorerConfig.SelectionUI.Prefix
		suffix = t.explorerConfig.SelectionUI.Suffix
		style = t.selectionStyle
	default:
		prefix = t.explorerConfig.DefaultUI.Prefix
		suffix = t.explorerConfig.DefaultUI.Suffix
		if entry.IsDirectory() {
			style = t.defaultDirectoryStyle
		} else {
			style = t.defaultFileStyle
		}
	}

	return entryDisplayState{
		prefix:     prefix,
		suffix:     suffix,
		treePrefix: t.getTreePrefix(idx),
		style:      style,
		isFocused:  isFocused,
		isSelected: isSelected,
	}
}

// getTreePrefix returns the appropriate tree connection prefix based on entry position
func (t *ExplorerTable) getTreePrefix(idx int) string {
	totalEntries := len(t.entries)
	switch idx {
	case totalEntries - 1:
		return t.explorerConfig.LastEntryPrefix
	case 0:
		return t.explorerConfig.FirstEntryPrefix
	default:
		return t.explorerConfig.EntryPrefix
	}
}

// buildEntryDisplayName constructs the complete display name with icon and formatting
func (t *ExplorerTable) buildEntryDisplayName(entry fs.IEntry, entryIcon nodeType, state entryDisplayState) string {
	iconText := entryIcon.icon
	fileName := strings.TrimSpace(entry.GetName())

	// Apply styling to just the icon if needed (but keep it simple)
	var styledIcon string
	if state.isFocused || state.isSelected {
		// For focused/selected items, apply same style to icon as text
		styledIcon = iconText
	} else {
		// For normal items, use icon's default style
		styledIcon = entryIcon.style.Render(iconText)
	}

	return state.treePrefix + state.prefix + styledIcon + " " + fileName + state.suffix
}

// formatEntryRow formats the complete row with index and name columns
func (t *ExplorerTable) formatEntryRow(
	idx int,
	nameColumn string,
	entryStyle lipgloss.Style,
) string {
	columns := []columnConfig{
		{percentage: t.explorerConfig.IndexHeader.Percentage, leftAlign: true},
		{percentage: t.explorerConfig.NameHeader.Percentage, leftAlign: true},
	}

	values := []styledValue{
		{text: strconv.Itoa(idx + 1)},
		{text: nameColumn, style: entryStyle},
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
		return "Invalid row configuration: column count mismatch"
	}
	if len(columns) == 0 {
		return ""
	}
	if t.width <= 0 {
		return ""
	}

	result := ""
	accumulatedColumnWidth := 0
	for i, col := range columns {
		width := int(float32(col.percentage) / 100.0 * float32(t.width))
		// Give remaining width to the last column to avoid rounding errors
		// that could leave empty space or cause overflow
		if i == len(columns)-1 {
			width = t.width - accumulatedColumnWidth
		} else {
			accumulatedColumnWidth += width
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
	default:
		icon = t.icons.file
	}

	switch {
	case isEntrySelected && isEntryFocused:
		icon.style = t.focusSelectionStyle
	case isEntrySelected:
		icon.style = t.selectionStyle
	case isEntryFocused:
		icon.style = t.focusStyle
	}

	return icon
}

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

// SetFocusByIndex sets the focus to a specific index
func (t *ExplorerTable) SetFocusByIndex(index int) bool {
	if len(t.entries) == 0 || index < 0 || index >= len(t.entries) {
		return false
	}

	t.focus = index
	t.ensureVisible()

	return true
}

// ToggleSelectionByPath toggles selection for an entry with the given path
func (t *ExplorerTable) ToggleSelectionByPath(path string) bool {
	if path == "" {
		return false
	}

	if t.selections.Contains(path) {
		t.selections.Remove(path)
	} else {
		t.selections.Add(path)
	}

	return true
}

// FocusPath attempts to focus on an entry with the given path
func (t *ExplorerTable) FocusPath(path string) bool {
	if path == "" || len(t.entries) == 0 {
		return false
	}

	for i, entry := range t.entries {
		if entry.GetPath() == path {
			t.focus = i
			t.ensureVisible()

			return true
		}
	}

	return false
}
