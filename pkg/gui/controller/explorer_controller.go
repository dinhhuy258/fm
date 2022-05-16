package controller

import (
	"strconv"

	set "github.com/deckarep/golang-set/v2"
	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui/view"
	"github.com/dinhhuy258/fm/pkg/optional"
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
	selections set.Set[string],
) *ExplorerController {
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

func (ec *ExplorerController) LoadDirectory(path string, focusPath optional.Optional[string]) {
	if !fs.IsDir(path) {
		// TODO: Showing log here
		return
	}

	ec.path = path

	cfg := config.AppConfig

	entries, err := fs.LoadEntries(path, cfg.ShowHidden)
	if err != nil {
		// TODO: Showing log here
		return
	}

	ec.entries = entries

	title := (" " + path + " (" + strconv.Itoa(len(ec.entries)) + ") ")
	ec.view.SetTitle(title)

	focusPath.IfPresentOrElse(func(focusPath *string) {
		ec.focusPath(*focusPath)
	}, func() {
		ec.FocusFirst()
	})
}

func (ec *ExplorerController) FocusPrevious() {
	if ec.focus <= 0 {
		return
	}

	_ = ec.view.PreviousCursor()
	ec.focus--

	ec.UpdateView()
}

func (ec *ExplorerController) FocusNext() {
	if ec.focus >= len(ec.entries)-1 {
		return
	}

	_ = ec.view.NextCursor()
	ec.focus++

	ec.UpdateView()
}

func (ec *ExplorerController) FocusFirst() {
	_ = ec.view.ResetCursor()

	ec.focus = 0

	ec.UpdateView()
}

func (ec *ExplorerController) FocusPath(path string) {
	if parentPath := fs.Dir(path); ec.path != parentPath {
		ec.LoadDirectory(parentPath, optional.NewOptional(path))
	} else {
		ec.focusPath(path)
	}
}

func (ec *ExplorerController) UpdateView() {
	ec.view.UpdateView(ec.entries, ec.selections, ec.focus)
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

	_ = ec.view.ResetCursor()

	for i := 0; i < focus; i++ {
		_ = ec.view.NextCursor()
	}

	ec.focus = focus

	ec.UpdateView()
}
