package controller

import (
	"strconv"

	set "github.com/deckarep/golang-set/v2"
	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui/view"
)

type ExplorerController struct {
	*BaseController

	path       string
	focus      int
	entries    []fs.IEntry
	selections set.Set[string]

	view *view.ExplorerView
}

func newExplorerController(baseController *BaseController,
	view *view.ExplorerView,
	selections set.Set[string]) *ExplorerController {
	return &ExplorerController{
		BaseController: baseController,
		view:           view,

		focus:      0,
		selections: selections,
	}
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

func (ec *ExplorerController) GetCurrentEntry() fs.IEntry {
	return ec.GetEntry(ec.focus)
}

func (ec *ExplorerController) GetEntries() []fs.IEntry {
	return ec.entries
}

func (ec *ExplorerController) GetPath() string {
	return ec.path
}

func (ec *ExplorerController) LoadDirectory(path string, focusPath string) {
	if !fs.IsDir(path) {
		// TODO: Showing log here
		return
	}

	ec.path = path

	//TODO: Goroutine
	cfg := config.AppConfig

	entries, err := fs.LoadEntries(path, cfg.ShowHidden)
	if err != nil {
		//TODO: Showing log here
		return
	}

	ec.entries = entries

	title := (" " + path + " (" + strconv.Itoa(len(ec.entries)) + ") ")
	ec.view.SetTitle(title)

	if focusPath != "" {
		ec.focusPath(focusPath)
	} else {
		ec.FocusFirst()
	}
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
	parentPath := fs.Dir(path)

	if ec.path != parentPath {
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
