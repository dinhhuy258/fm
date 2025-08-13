package tui

import (
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
	set "github.com/deckarep/golang-set/v2"
	"github.com/rivo/uniseg"

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

// ExplorerViewData holds the computed styles and icons for rendering
type ExplorerViewData struct {
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

// ExplorerModel represents the pure state for the file explorer table
type ExplorerModel struct {
	// Display dimensions
	width  int
	height int

	// File system state
	entries []fs.IEntry

	// Navigation state
	focus       int
	scrollStart int

	// Selection state
	selections set.Set[string]

	// Cached view data for performance
	viewData *ExplorerViewData
}

// NewExplorerModel creates a new explorer model
func NewExplorerModel() *ExplorerModel {
	model := &ExplorerModel{
		selections:  set.NewSet[string](),
		focus:       0,
		scrollStart: 0,
		entries:     make([]fs.IEntry, 0),
	}

	// Initialize view data on creation
	model.initViewData()

	return model
}

// SetSize updates the model dimensions
func (m *ExplorerModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// GetSize returns the current dimensions
func (m *ExplorerModel) GetSize() (int, int) {
	return m.width, m.height
}

// SetEntries updates the entries and resets focus/selection state
func (m *ExplorerModel) SetEntries(entries []fs.IEntry) {
	m.entries = entries
	m.focus = 0
	m.scrollStart = 0
	m.selections = set.NewSet[string]()
}

// GetEntries returns the current entries
func (m *ExplorerModel) GetEntries() []fs.IEntry {
	return m.entries
}

// Move moves the cursor by delta positions
func (m *ExplorerModel) Move(delta int) {
	if len(m.entries) == 0 {
		return
	}

	m.focus = m.focus + delta
	m.focus = max(0, min(m.focus, len(m.entries)-1))
	m.ensureVisible()
}

// MoveFirst moves focus to the first item
func (m *ExplorerModel) MoveFirst() {
	if len(m.entries) == 0 {
		return
	}
	m.focus = 0
	m.scrollStart = 0
}

// MoveLast moves focus to the last item
func (m *ExplorerModel) MoveLast() {
	if len(m.entries) == 0 {
		return
	}
	m.focus = len(m.entries) - 1
	m.ensureVisible()
}

// GetFocus returns the current focus index
func (m *ExplorerModel) GetFocus() int {
	return m.focus
}

// SetFocusByIndex sets the focus to a specific index
func (m *ExplorerModel) SetFocusByIndex(index int) bool {
	if len(m.entries) == 0 || index < 0 || index >= len(m.entries) {
		return false
	}

	m.focus = index
	m.ensureVisible()

	return true
}

// GetFocusedEntry returns the currently focused entry
func (m *ExplorerModel) GetFocusedEntry() fs.IEntry {
	if m.focus < len(m.entries) {
		return m.entries[m.focus]
	}

	return nil
}

// ToggleSelection toggles selection for the focused item
func (m *ExplorerModel) ToggleSelection() {
	if m.focus >= len(m.entries) {
		return
	}

	entry := m.entries[m.focus]
	path := entry.GetPath()
	m.ToggleSelectionByPath(path)
}

// ToggleSelectionByPath toggles selection for an entry with the given path
func (m *ExplorerModel) ToggleSelectionByPath(path string) bool {
	if path == "" {
		return false
	}

	if m.selections.Contains(path) {
		m.selections.Remove(path)
	} else {
		m.selections.Add(path)
	}

	return true
}

// ClearSelections clears all selections
func (m *ExplorerModel) ClearSelections() {
	m.selections = set.NewSet[string]()
}

// SelectAll selects all entries
func (m *ExplorerModel) SelectAll() {
	for _, entry := range m.entries {
		m.selections.Add(entry.GetPath())
	}
}

// GetSelectedEntries returns all selected entries
func (m *ExplorerModel) GetSelectedEntries() []fs.IEntry {
	var selected []fs.IEntry
	for _, entry := range m.entries {
		if m.selections.Contains(entry.GetPath()) {
			selected = append(selected, entry)
		}
	}

	return selected
}

// GetStats returns total and selected entry counts
func (m *ExplorerModel) GetStats() (total, selected int) {
	return len(m.entries), m.selections.Cardinality()
}

// IsSelected returns whether an entry path is selected
func (m *ExplorerModel) IsSelected(path string) bool {
	return m.selections.Contains(path)
}

// FocusPath attempts to focus on an entry with the given path
func (m *ExplorerModel) FocusPath(path string) bool {
	if path == "" || len(m.entries) == 0 {
		return false
	}

	for i, entry := range m.entries {
		if entry.GetPath() == path {
			m.focus = i
			m.ensureVisible()

			return true
		}
	}

	return false
}

// GetScrollStart returns the current scroll position
func (m *ExplorerModel) GetScrollStart() int {
	return m.scrollStart
}

// GetVisibleRows calculates how many rows can fit in the current height
func (m *ExplorerModel) GetVisibleRows() int {
	// Reserve one row for the header
	return max(m.height-1, 1)
}

// ensureVisible ensures the focused item is visible by adjusting scroll
func (m *ExplorerModel) ensureVisible() {
	visibleRows := m.GetVisibleRows()

	if m.focus < m.scrollStart {
		m.scrollStart = m.focus
	} else if m.focus >= m.scrollStart+visibleRows {
		m.scrollStart = m.focus - visibleRows + 1
	}

	m.scrollStart = max(0, m.scrollStart)
}

// initViewData initializes the cached view data with styles and icons
func (m *ExplorerModel) initViewData() {
	m.viewData = &ExplorerViewData{}
	m.viewData.initStyles()
	m.viewData.initIcons()
}

// InvalidateViewData clears cached view data (call when config changes)
func (m *ExplorerModel) InvalidateViewData() {
	m.viewData = nil
	m.initViewData()
}

// GetViewData returns the cached view data, initializing if needed
func (m *ExplorerModel) GetViewData() *ExplorerViewData {
	if m.viewData == nil {
		m.initViewData()
	}

	return m.viewData
}

// initStyles initializes lipgloss styles from config
func (d *ExplorerViewData) initStyles() {
	explorerConfig := config.AppConfig.General.ExplorerTable
	d.defaultFileStyle = fromStyleConfig(explorerConfig.DefaultUI.FileStyle)
	d.defaultDirectoryStyle = fromStyleConfig(explorerConfig.DefaultUI.DirectoryStyle)
	d.focusStyle = fromStyleConfig(explorerConfig.FocusUI.Style)
	d.selectionStyle = fromStyleConfig(explorerConfig.SelectionUI.Style)
	d.focusSelectionStyle = fromStyleConfig(explorerConfig.FocusSelectionUI.Style)
	d.headerStyles = headerStyles{
		indexHeader: fromStyleConfig(explorerConfig.IndexHeader.Style),
		nameHeader:  fromStyleConfig(explorerConfig.NameHeader.Style),
	}
}

// initIcons initializes the icon system from config
func (d *ExplorerViewData) initIcons() {
	nodeTypesConfig := config.AppConfig.NodeTypes
	d.icons = nodeTypes{
		file: nodeType{
			icon:  nodeTypesConfig.File.Icon,
			style: fromStyleConfig(nodeTypesConfig.File.Style),
		},
		directory: nodeType{
			icon:  nodeTypesConfig.Directory.Icon,
			style: fromStyleConfig(nodeTypesConfig.Directory.Style),
		},
		fileSymlink: nodeType{
			icon:  nodeTypesConfig.FileSymlink.Icon,
			style: fromStyleConfig(nodeTypesConfig.FileSymlink.Style),
		},
		directorySymlink: nodeType{
			icon:  nodeTypesConfig.DirectorySymlink.Icon,
			style: fromStyleConfig(nodeTypesConfig.DirectorySymlink.Style),
		},
		extensions: make(map[string]nodeType),
		specials:   make(map[string]nodeType),
	}

	// Load extension-specific icons
	for ext, ntc := range nodeTypesConfig.Extensions {
		d.icons.extensions[strings.ToLower(ext)] = nodeType{
			icon:  ntc.Icon,
			style: fromStyleConfig(ntc.Style),
		}
	}

	// Load special file icons
	for fileName, ntc := range nodeTypesConfig.Specials {
		d.icons.specials[strings.ToLower(fileName)] = nodeType{
			icon:  ntc.Icon,
			style: fromStyleConfig(ntc.Style),
		}
	}
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

// entryDisplayState holds the UI state for rendering an entry
type entryDisplayState struct {
	prefix     string
	suffix     string
	treePrefix string
	style      lipgloss.Style
	isFocused  bool
	isSelected bool
}

// View renders the explorer table view using cached view data
func (m *ExplorerModel) View() string {
	width, height := m.GetSize()
	if width <= 0 || height <= 0 {
		return ""
	}

	// Use cached view data
	viewData := m.GetViewData()

	var sections []string

	// Render header
	sections = append(sections, m.renderHeader(viewData, width))

	// Render visible entries
	sections = append(sections, m.renderEntries(viewData, width))

	return strings.Join(sections, "\n")
}

// renderHeader renders the column headers
func (m *ExplorerModel) renderHeader(viewData *ExplorerViewData, width int) string {
	explorerConfig := config.AppConfig.General.ExplorerTable
	columns := []columnConfig{
		{percentage: explorerConfig.IndexHeader.Percentage, leftAlign: true},
		{percentage: explorerConfig.NameHeader.Percentage, leftAlign: true},
	}

	values := []styledValue{
		{text: explorerConfig.IndexHeader.Name, style: viewData.headerStyles.indexHeader},
		{text: explorerConfig.NameHeader.Name, style: viewData.headerStyles.nameHeader},
	}

	return m.formatRow(columns, values, width)
}

// renderEntries renders the visible file entries
func (m *ExplorerModel) renderEntries(viewData *ExplorerViewData, width int) string {
	if len(m.entries) == 0 {
		return ""
	}

	visibleRows := m.GetVisibleRows()
	lastVisibleIndex := min(m.scrollStart+visibleRows, len(m.entries))

	lines := make([]string, 0, lastVisibleIndex-m.scrollStart)

	for i := m.scrollStart; i < lastVisibleIndex; i++ {
		lines = append(lines, m.renderEntry(m.entries[i], i, viewData, width))
	}

	// Pad with empty lines if needed
	for len(lines) < visibleRows {
		lines = append(lines, strings.Repeat(" ", width))
	}

	return strings.Join(lines, "\n")
}

// renderEntry renders a single file entry with proper styling
func (m *ExplorerModel) renderEntry(entry fs.IEntry, idx int, viewData *ExplorerViewData, width int) string {
	state := m.determineEntryDisplayState(entry, idx, viewData)
	entryIcon := m.getEntryIcon(entry, state.isFocused, state.isSelected, viewData)
	nameColumn := m.buildEntryDisplayName(entry, entryIcon, state)

	return m.formatEntryRow(idx, nameColumn, state.style, width)
}

// determineEntryDisplayState calculates the display state for an entry based on focus/selection
func (m *ExplorerModel) determineEntryDisplayState(entry fs.IEntry, idx int, viewData *ExplorerViewData) entryDisplayState {
	explorerConfig := config.AppConfig.General.ExplorerTable
	isFocused := idx == m.focus
	isSelected := m.IsSelected(entry.GetPath())

	var prefix, suffix string
	var style lipgloss.Style

	switch {
	case isFocused && isSelected:
		prefix = explorerConfig.FocusSelectionUI.Prefix
		suffix = explorerConfig.FocusSelectionUI.Suffix
		style = viewData.focusSelectionStyle
	case isFocused:
		prefix = explorerConfig.FocusUI.Prefix
		suffix = explorerConfig.FocusUI.Suffix
		style = viewData.focusStyle
	case isSelected:
		prefix = explorerConfig.SelectionUI.Prefix
		suffix = explorerConfig.SelectionUI.Suffix
		style = viewData.selectionStyle
	default:
		prefix = explorerConfig.DefaultUI.Prefix
		suffix = explorerConfig.DefaultUI.Suffix
		if entry.IsDirectory() {
			style = viewData.defaultDirectoryStyle
		} else {
			style = viewData.defaultFileStyle
		}
	}

	return entryDisplayState{
		prefix:     prefix,
		suffix:     suffix,
		treePrefix: m.getTreePrefix(idx),
		style:      style,
		isFocused:  isFocused,
		isSelected: isSelected,
	}
}

// getTreePrefix returns the appropriate tree connection prefix based on entry position
func (m *ExplorerModel) getTreePrefix(idx int) string {
	explorerConfig := config.AppConfig.General.ExplorerTable
	switch idx {
	case len(m.entries) - 1:
		return explorerConfig.LastEntryPrefix
	case 0:
		return explorerConfig.FirstEntryPrefix
	default:
		return explorerConfig.EntryPrefix
	}
}

// buildEntryDisplayName constructs the complete display name with icon and formatting
func (m *ExplorerModel) buildEntryDisplayName(entry fs.IEntry, entryIcon nodeType, state entryDisplayState) string {
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
func (m *ExplorerModel) formatEntryRow(idx int, nameColumn string, entryStyle lipgloss.Style, width int) string {
	explorerConfig := config.AppConfig.General.ExplorerTable
	columns := []columnConfig{
		{percentage: explorerConfig.IndexHeader.Percentage, leftAlign: true},
		{percentage: explorerConfig.NameHeader.Percentage, leftAlign: true},
	}

	values := []styledValue{
		{text: strconv.Itoa(idx + 1)},
		{text: nameColumn, style: entryStyle},
	}

	return m.formatRow(columns, values, width)
}

// getEntryIcon returns the appropriate icon for an entry with state-based styling
func (m *ExplorerModel) getEntryIcon(entry fs.IEntry, isEntryFocused, isEntrySelected bool, viewData *ExplorerViewData) nodeType {
	var icon nodeType

	// Find the appropriate icon based on file type
	extensionIcon, hasExtIcon := viewData.icons.extensions[strings.ToLower(entry.GetExt())]
	specialIcon, hasSpecialIcon := viewData.icons.specials[strings.ToLower(entry.GetName())]

	switch {
	case entry.IsSymlink() && entry.IsDirectory():
		icon = viewData.icons.directorySymlink
	case entry.IsSymlink():
		icon = viewData.icons.fileSymlink
	case hasSpecialIcon:
		icon = specialIcon
	case hasExtIcon:
		icon = extensionIcon
	case entry.IsDirectory():
		icon = viewData.icons.directory
	default:
		icon = viewData.icons.file
	}

	switch {
	case isEntrySelected && isEntryFocused:
		icon.style = viewData.focusSelectionStyle
	case isEntrySelected:
		icon.style = viewData.selectionStyle
	case isEntryFocused:
		icon.style = viewData.focusStyle
	}

	return icon
}

// formatRow formats a row with proper column alignment and styling
func (m *ExplorerModel) formatRow(columns []columnConfig, values []styledValue, width int) string {
	if len(columns) != len(values) {
		return "Invalid row configuration: column count mismatch"
	}
	if len(columns) == 0 {
		return ""
	}
	if width <= 0 {
		return ""
	}

	result := ""
	accumulatedColumnWidth := 0
	for i, col := range columns {
		columnWidth := int(float32(col.percentage) / 100.0 * float32(width))
		// Give remaining width to the last column to avoid rounding errors
		// that could leave empty space or cause overflow
		if i == len(columns)-1 {
			columnWidth = width - accumulatedColumnWidth
		} else {
			accumulatedColumnWidth += columnWidth
		}
		result += m.formatColumn(values[i], columnWidth, col.leftAlign)
	}

	// Ensure the row doesn't exceed terminal width
	if uniseg.StringWidth(result) > width {
		runes := []rune(result)
		if len(runes) > width {
			result = string(runes[:width])
		}
	}

	return result
}

// formatColumn formats a single column with proper alignment
func (m *ExplorerModel) formatColumn(value styledValue, width int, leftAlign bool) string {
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
