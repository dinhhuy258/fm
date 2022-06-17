package controller

import (
	set "github.com/deckarep/golang-set/v2"
	"github.com/dinhhuy258/fm/pkg/gui/view"
)

// SelectionController is a controller for selection view.
type SelectionController struct {
	*BaseController

	selections set.Set[string]

	view *view.SelectionView
}

// newSelectionController creates a new selection controller.
func newSelectionController(baseController *BaseController,
	view *view.SelectionView,
	selections set.Set[string],
) *SelectionController {
	return &SelectionController{
		BaseController: baseController,
		view:           view,
		selections:     selections,
	}
}

// ClearSelections clears all selections.
func (sc *SelectionController) ClearSelections() {
	sc.selections.Clear()
}

// ToggleSelection toggles the selection of the given item.
func (sc *SelectionController) ToggleSelection(path string) {
	if sc.selections.Contains(path) {
		sc.selections.Remove(path)
	} else {
		sc.selections.Add(path)
	}
}

// SelectPath add the given path to selection list.
func (sc *SelectionController) SelectPath(path string) {
	sc.selections.Add(path)
}

// GetSelections returns the current selections.
func (sc *SelectionController) GetSelections() []string {
	return sc.selections.ToSlice()
}

// UpdateView updates the view.
func (sc *SelectionController) UpdateView() {
	sc.view.UpdateView(sc.selections.ToSlice())
}
