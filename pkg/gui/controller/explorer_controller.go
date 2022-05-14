package controller

import (
	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui/view"
)

type ExplorerController struct {
	focus      int
	entries    []fs.IEntry
	selections map[string]struct{}

	view *view.ExplorerView
}

func newExplorerController() *ExplorerController {
	return &ExplorerController{
		focus: 0,
	}
}

func (ec *ExplorerController) SetView(view *view.ExplorerView) {
	ec.view = view
}

func (ec *ExplorerController) GetFocus() int {
	return ec.focus
}

func (ec *ExplorerController) GetEntry(idx int) fs.IEntry {
	if idx < 0 || idx >= len(ec.entries) {
		return nil
	}

	return ec.entries[idx]
}

func (ec *ExplorerController) LoadDirectory(path string, focusPath string) {
	fileExplorer := fs.GetFileExplorer()

	// TODO: Find the right way to do this
	fileExplorer.LoadEntries(path, func() {
		ec.entries = fileExplorer.GetEntries()

		if focusPath != "" {
			ec.focusPath(focusPath)
		} else {
			ec.FocusFirst()
		}
	})
}

func (ec *ExplorerController) FocusPrevious() {
	if ec.focus <= 0 {
		return
	}

	ec.view.PreviousCursor()
	ec.focus = ec.focus - 1

	ec.UpdateView()
}

func (ec *ExplorerController) FocusNext() {
	if ec.focus >= len(ec.entries)-1 {
		return
	}

	ec.view.NextCursor()
	ec.focus = ec.focus + 1

	ec.UpdateView()
}

func (ec *ExplorerController) FocusFirst() {
	ec.view.ResetCursor()

	ec.focus = 0

	ec.UpdateView()
}

func (ec *ExplorerController) FocusPath(path string) {
	// TODO: Verify if the path is valid
	fileExplorer := fs.GetFileExplorer()
	parentPath := fs.Dir(path)

	if fileExplorer.GetPath() != parentPath {
		ec.LoadDirectory(parentPath, path)
	} else {
		ec.focusPath(path)
	}
}

func (ec *ExplorerController) UpdateView() {
	// TODO: Rename view function
	ec.view.RenderEntries(ec.entries, ec.selections, ec.focus)
}

func (ec *ExplorerController) focusPath(path string) {
	focus := 0
	// Iterate through the list of entries and find the idx for the current path
	for idx, entry := range ec.entries {
		if entry.GetPath() == path {
			focus = idx

			break
		}
	}

	ec.view.ResetCursor()

	for i := 0; i < focus; i++ {
		ec.view.NextCursor()
	}

	ec.focus = focus

	ec.UpdateView()
}
