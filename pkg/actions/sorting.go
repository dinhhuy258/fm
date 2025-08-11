package actions

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/dinhhuy258/fm/pkg/config"
)

// executeSortByName handles sorting by name
func (ah *ActionHandler) executeSortByName(_ *config.MessageConfig) tea.Cmd {
	return func() tea.Msg {
		return SortingMessage{SortType: "name"}
	}
}

// executeSortBySize handles sorting by size
func (ah *ActionHandler) executeSortBySize(_ *config.MessageConfig) tea.Cmd {
	return func() tea.Msg {
		return SortingMessage{SortType: "size"}
	}
}

// executeSortByDateModified handles sorting by date modified
func (ah *ActionHandler) executeSortByDateModified(_ *config.MessageConfig) tea.Cmd {
	return func() tea.Msg {
		return SortingMessage{SortType: "date"}
	}
}

// executeSortByExtension handles sorting by extension
func (ah *ActionHandler) executeSortByExtension(_ *config.MessageConfig) tea.Cmd {
	return func() tea.Msg {
		return SortingMessage{SortType: "extension"}
	}
}

// executeSortByDirFirst handles sorting by directory first
func (ah *ActionHandler) executeSortByDirFirst(_ *config.MessageConfig) tea.Cmd {
	return func() tea.Msg {
		return SortingMessage{SortType: "dir_first"}
	}
}

// executeReverseSort handles reversing sort order
func (ah *ActionHandler) executeReverseSort(_ *config.MessageConfig) tea.Cmd {
	return func() tea.Msg {
		return SortingMessage{SortType: "reverse"}
	}
}
