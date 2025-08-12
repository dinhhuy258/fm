package tui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/dinhhuy258/fm/pkg/fs"
)

// FileItem represents a file entry for the list component
type FileItem struct {
	entry fs.IEntry
}

// NewFileItem creates a new file item
func NewFileItem(entry fs.IEntry) FileItem {
	return FileItem{entry: entry}
}

// FilterValue returns the value used for filtering
func (f FileItem) FilterValue() string {
	return f.entry.GetName()
}

// Title returns the display title
func (f FileItem) Title() string {
	return f.entry.GetName()
}

// Description returns the display description
func (f FileItem) Description() string {
	return f.entry.GetName()
}

// GetEntry returns the underlying file entry
func (f FileItem) GetEntry() fs.IEntry {
	return f.entry
}

// FileItemDelegate handles the rendering of file items
type FileItemDelegate struct {
	styles FileItemStyles
}

// FileItemStyles defines the styles for file items
type FileItemStyles struct {
	NormalTitle   lipgloss.Style
	NormalDesc    lipgloss.Style
	SelectedTitle lipgloss.Style
	SelectedDesc  lipgloss.Style
	DimmedTitle   lipgloss.Style
	DimmedDesc    lipgloss.Style
	FilterMatch   lipgloss.Style
	DirectoryIcon lipgloss.Style
	FileIcon      lipgloss.Style
	SymlinkIcon   lipgloss.Style
}

// NewFileItemDelegate creates a new file item delegate
func NewFileItemDelegate() FileItemDelegate {
	return FileItemDelegate{
		styles: DefaultFileItemStyles(),
	}
}

// DefaultFileItemStyles returns the default styles for file items
func DefaultFileItemStyles() FileItemStyles {
	return FileItemStyles{
		NormalTitle: lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"}).
			Padding(0, 0, 0, 2),

		NormalDesc: lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"}),

		SelectedTitle: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"}).
			Foreground(lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: "#EE6FF8"}).
			Padding(0, 0, 0, 1),

		SelectedDesc: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"}).
			Foreground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"}).
			Padding(0, 0, 0, 1),

		DimmedTitle: lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"}).
			Padding(0, 0, 0, 2),

		DimmedDesc: lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#C2B8C2", Dark: "#4D4D4D"}),

		FilterMatch: lipgloss.NewStyle().
			Underline(true),

		DirectoryIcon: lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#3C82F6", Dark: "#60A5FA"}),

		FileIcon: lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#6B7280", Dark: "#9CA3AF"}),

		SymlinkIcon: lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#F59E0B", Dark: "#FCD34D"}),
	}
}

// Height returns the height of the item
func (d FileItemDelegate) Height() int {
	return 1
}

// Spacing returns the spacing between items
func (d FileItemDelegate) Spacing() int {
	return 0
}

// Update handles updates for the delegate
func (d FileItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

// Render renders the file item
func (d FileItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	if fileItem, ok := listItem.(FileItem); ok {
		var (
			title, desc  string
			matchedRunes []int
			s            = &d.styles
		)

		if m.Width() <= 0 {
			// short-circuit
			return
		}

		// Prepare content
		title = fileItem.Title()
		desc = fileItem.Description()

		if m.FilterState() == list.Filtering {
			// Get the matched character indices
			matchedRunes = m.MatchesForItem(index)
		}

		// Get the file icon based on type
		icon := d.getFileIcon(fileItem.entry)
		title = icon + " " + title

		// Determine if this item is selected
		isSelected := index == m.Index()

		// Style the title and description
		if isSelected && m.FilterState() != list.Filtering {
			title = s.SelectedTitle.Render(title)
			desc = s.SelectedDesc.Render(desc)
		} else if isSelected && m.FilterState() == list.Filtering {
			title = s.DimmedTitle.Render(title)
			desc = s.DimmedDesc.Render(desc)
		} else {
			title = s.NormalTitle.Render(title)
			desc = s.NormalDesc.Render(desc)
		}

		// Handle filter matches
		if m.FilterState() == list.Filtering && matchedRunes != nil {
			unmatched := s.DimmedTitle.Inline(true)
			matched := unmatched.Inherit(s.FilterMatch)

			title = lipgloss.StyleRunes(title, matchedRunes, matched, unmatched)
		}

		_, _ = fmt.Fprintf(w, "%s", title)
		if desc != "" {
			_, _ = fmt.Fprintf(w, "\n%s", desc)
		}
	}
}

// getFileIcon returns the appropriate icon for the file type
func (d FileItemDelegate) getFileIcon(entry fs.IEntry) string {
	var icon string
	var style lipgloss.Style

	switch {
	case entry.IsSymlink() && entry.IsDirectory():
		icon = "📂" // Directory symlink
		style = d.styles.SymlinkIcon
	case entry.IsSymlink():
		icon = "🔗" // File symlink
		style = d.styles.SymlinkIcon
	case entry.IsDirectory():
		icon = "📁" // Directory
		style = d.styles.DirectoryIcon
	default:
		// Get icon based on extension
		ext := strings.ToLower(entry.GetExt())
		switch ext {
		case ".go":
			icon = "🐹"
		case ".js", ".ts":
			icon = "📜"
		case ".py":
			icon = "🐍"
		case ".rs":
			icon = "🦀"
		case ".java":
			icon = "☕"
		case ".cpp", ".cc", ".cxx", ".c":
			icon = "⚙️"
		case ".md", ".txt":
			icon = "📄"
		case ".json", ".yaml", ".yml", ".toml":
			icon = "⚙️"
		case ".png", ".jpg", ".jpeg", ".gif", ".svg":
			icon = "🖼️"
		case ".mp3", ".wav", ".flac":
			icon = "🎵"
		case ".mp4", ".avi", ".mkv":
			icon = "🎥"
		case ".zip", ".tar", ".gz", ".7z":
			icon = "📦"
		case ".pdf":
			icon = "📕"
		default:
			icon = "📄"
		}
		style = d.styles.FileIcon
	}

	return style.Render(icon)
}