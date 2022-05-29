package controller

import (
	"os"
	"strconv"

	set "github.com/deckarep/golang-set/v2"
	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui/view"
	"github.com/dinhhuy258/fm/pkg/optional"
)

// ExplorerController is the controller for the explorer view
type ExplorerController struct {
	*BaseController

	path       string
	focus      int
	entries    []fs.IEntry
	selections set.Set[string]

	view       *view.ExplorerView
	headerView *view.ExplorerHeaderView
}

// newExplorerController creates a new explorer controller
func newExplorerController(baseController *BaseController,
	view *view.ExplorerView,
	headerView *view.ExplorerHeaderView,
	selections set.Set[string],
) *ExplorerController {
	return &ExplorerController{
		BaseController: baseController,
		view:           view,
		headerView:     headerView,

		focus:      0,
		selections: selections,
	}
}

// GetFocus returns the index of the focused entry
func (ec *ExplorerController) GetFocus() int {
	return ec.focus
}

// GetCurrentEntry returns the current entry
func (ec *ExplorerController) GetCurrentEntry() fs.IEntry {
	if ec.focus < 0 || ec.focus >= len(ec.entries) {
		return nil
	}

	return ec.entries[ec.focus]
}

// GetEntries returns the list of entries
func (ec *ExplorerController) GetEntries() []fs.IEntry {
	return ec.entries
}

// GetPath returns the current path
func (ec *ExplorerController) GetPath() string {
	return ec.path
}

// LoadDirectory loads the directory at the given path
func (ec *ExplorerController) LoadDirectory(path string, focusPath optional.Optional[string]) {
	if !fs.IsDir(path) {
		ec.mediator.notify(ShowErrorLog, optional.New(path+" is not a directory"))

		return
	}

	ec.path = path

	cfg := config.AppConfig

	entries, err := fs.LoadEntries(path, cfg.ShowHidden)
	if err != nil {
		ec.mediator.notify(ShowErrorLog, optional.New("Failed to load directory: "+path))

		return
	}

	// Change the working directory according to the new path
	// We can be sure that the path is a directory and its existence is checked before => no need to
	// check the error here
	_ = os.Chdir(ec.path)

	ec.entries = entries
	ec.headerView.Title = " " + path + " (" + strconv.Itoa(len(ec.entries)) + ") "

	focusPath.IfPresentOrElse(func(focusPath *string) {
		ec.focusPath(*focusPath)
	}, func() {
		ec.FocusFirst()
	})
}

// Focus focuses cusor to the current entry
func (ec *ExplorerController) Focus() {
	ec.view.FocusPoint(0, ec.focus)
}

// FocusPrevious focuses the previous entry
func (ec *ExplorerController) FocusPrevious() {
	if ec.focus <= 0 {
		return
	}

	ec.focus--
	ec.Focus()
}

// FocusNext focuses the next entry
func (ec *ExplorerController) FocusNext() {
	if ec.focus >= len(ec.entries)-1 {
		return
	}

	ec.focus++
	ec.Focus()
}

// FocusFirst focuses the first entry
func (ec *ExplorerController) FocusFirst() {
	ec.focus = 0
	ec.Focus()
}

// FocusLast focuses the last entry
func (ec *ExplorerController) FocusLast() {
	ec.focus = len(ec.entries) - 1
	ec.Focus()
}

// FocusPath focuses the entry with the given path
func (ec *ExplorerController) FocusPath(path string) {
	if parentPath := fs.Dir(path); ec.path != parentPath {
		ec.LoadDirectory(parentPath, optional.New(path))
	} else {
		ec.focusPath(path)
	}
}

// UpdateView updates the view
func (ec *ExplorerController) UpdateView() {
	ec.view.UpdateView(ec.entries, ec.selections, ec.focus)
}

// focusPath focuses the entry with the given path
func (ec *ExplorerController) focusPath(path string) {
	focus := 0
	// Iterate through the list of entries and find the idx for the current path
	for idx, entry := range ec.entries {
		if entry.GetPath() == path {
			focus = idx

			break
		}
	}

	ec.focus = focus
	ec.Focus()
}
