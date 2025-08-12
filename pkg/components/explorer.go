package components

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/dinhhuy258/fm/pkg/fs"
)

// Explorer represents the file explorer component
// This is the Bubble Tea equivalent of the original ExplorerController
type Explorer struct {
	list        list.Model
	width       int
	height      int
	entries     []fs.IEntry
	selected    map[int]struct{}
	showHeader  bool
	currentPath string

	// Styles
	titleStyle      lipgloss.Style
	itemStyle       lipgloss.Style
	selectedStyle   lipgloss.Style
	paginationStyle lipgloss.Style
	helpStyle       lipgloss.Style
}

// NewExplorer creates a new explorer component
func NewExplorer() *Explorer {
	// Create list with custom delegate
	delegate := NewExplorerDelegate()
	l := list.New([]list.Item{}, delegate, 0, 0)

	l.Title = "Files"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = lipgloss.NewStyle().
		MarginLeft(2).
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA"))

	return &Explorer{
		list:            l,
		selected:        make(map[int]struct{}),
		showHeader:      true,
		currentPath:     "",
		titleStyle:      lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FAFAFA")),
		itemStyle:       lipgloss.NewStyle().PaddingLeft(2),
		selectedStyle:   lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#FF69B4")),
		paginationStyle: lipgloss.NewStyle().PaddingLeft(2),
		helpStyle:       lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")),
	}
}

// SetSize updates the size of the explorer
func (e *Explorer) SetSize(width, height int) {
	e.width = width
	e.height = height
	e.list.SetSize(width, height)
}

// SetEntries updates the entries displayed in the explorer
func (e *Explorer) SetEntries(entries []fs.IEntry, currentPath string) {
	e.entries = entries
	e.currentPath = currentPath

	// Convert entries to list items
	items := make([]list.Item, len(entries))
	for i, entry := range entries {
		items[i] = ExplorerItem{entry: entry}
	}

	e.list.SetItems(items)
	e.selected = make(map[int]struct{})
}

// ToggleSelection toggles the selection of the current item
func (e *Explorer) ToggleSelection() {
	index := e.list.Index()
	if _, ok := e.selected[index]; ok {
		delete(e.selected, index)
	} else {
		e.selected[index] = struct{}{}
	}
}

// GetSelectedEntry returns the currently focused entry
func (e *Explorer) GetSelectedEntry() fs.IEntry {
	if item := e.list.SelectedItem(); item != nil {
		if explorerItem, ok := item.(ExplorerItem); ok {
			return explorerItem.entry
		}
	}

	return nil
}

// GetSelectedEntries returns all selected entries
func (e *Explorer) GetSelectedEntries() []fs.IEntry {
	var selected []fs.IEntry
	for index := range e.selected {
		if index < len(e.entries) {
			selected = append(selected, e.entries[index])
		}
	}

	return selected
}

// Update handles the explorer's update logic
func (e *Explorer) Update(msg tea.Msg) (*Explorer, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case " ":
			// Toggle selection
			e.ToggleSelection()

			return e, nil
		}
	}

	// Let the list handle other updates
	e.list, cmd = e.list.Update(msg)

	return e, cmd
}

// View renders the explorer
func (e *Explorer) View() string {
	if e.width <= 0 || e.height <= 0 {
		return "Explorer not sized"
	}

	// Render the list
	content := e.list.View()

	// Add selection indicators if any items are selected
	if len(e.selected) > 0 {
		lines := strings.Split(content, "\n")
		for index := range e.selected {
			if index < len(lines) {
				lines[index] = "● " + lines[index]
			}
		}
		content = strings.Join(lines, "\n")
	}

	return content
}

// GetStats returns statistics about the current explorer state
func (e *Explorer) GetStats() (total, selected int) {
	return len(e.entries), len(e.selected)
}

// ExplorerItem represents an item in the explorer list
type ExplorerItem struct {
	entry fs.IEntry
}

// FilterValue returns the value used for filtering
func (i ExplorerItem) FilterValue() string {
	return i.entry.GetName()
}

// Title returns the display title
func (i ExplorerItem) Title() string {
	return i.entry.GetName()
}

// Description returns the display description
func (i ExplorerItem) Description() string {
	return fmt.Sprintf("%s  %s",
		i.entry.GetPermissions(),
		fs.Humanize(i.entry.GetSize()))
}

// ExplorerDelegate handles rendering of explorer items
type ExplorerDelegate struct {
	ShowDesc bool
}

// NewExplorerDelegate creates a new explorer delegate
func NewExplorerDelegate() ExplorerDelegate {
	return ExplorerDelegate{ShowDesc: true}
}

// Height returns the height of each item
func (d ExplorerDelegate) Height() int {
	if d.ShowDesc {
		return 2
	}

	return 1
}

// Spacing returns the spacing between items
func (d ExplorerDelegate) Spacing() int {
	return 0
}

// Update handles updates for the delegate
func (d ExplorerDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

// Render renders an explorer item
func (d ExplorerDelegate) Render(w io.Writer, _ list.Model, index int, item list.Item) {
	str := fmt.Sprintf("%d. %s", index+1, item.FilterValue())

	if explorerItem, ok := item.(ExplorerItem); ok {
		icon := getFileIcon(explorerItem.entry)
		str = fmt.Sprintf("%d. %s %s", index+1, icon, item.FilterValue())

		if d.ShowDesc {
			desc := explorerItem.Description()
			str += "\n   " + desc
		}
	}

	_, _ = fmt.Fprint(w, str)
}

// getFileIcon returns an appropriate icon for the file type
func getFileIcon(entry fs.IEntry) string {
	if entry.IsDirectory() {
		if entry.IsSymlink() {
			return "📂"
		}

		return "📁"
	}

	if entry.IsSymlink() {
		return "🔗"
	}

	// Return icon based on extension
	ext := strings.ToLower(entry.GetExt())
	switch ext {
	case "go":
		return "🐹"
	case "js", "ts":
		return "📜"
	case "py":
		return "🐍"
	case "rs":
		return "🦀"
	case "java":
		return "☕"
	case "cpp", "cc", "cxx", "c":
		return "⚙️"
	case "md", "txt":
		return "📄"
	case "json", "yaml", "yml":
		return "⚙️"
	case "png", "jpg", "jpeg", "gif":
		return "🖼️"
	case "mp3", "wav":
		return "🎵"
	case "mp4", "avi":
		return "🎥"
	case "zip", "tar", "gz":
		return "📦"
	case "pdf":
		return "📕"
	default:
		return "📄"
	}
}
