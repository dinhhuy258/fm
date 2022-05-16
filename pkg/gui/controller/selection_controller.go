package controller

import (
	set "github.com/deckarep/golang-set/v2"
	"github.com/dinhhuy258/fm/pkg/gui/view"
)

type SelectionController struct {
	*BaseController

	selections set.Set[string]

	view *view.SelectionView
}

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

func (sc *SelectionController) ClearSelections() {
	sc.selections.Clear()
}

func (sc *SelectionController) ToggleSelection(path string) {
	if sc.selections.Contains(path) {
		sc.selections.Remove(path)
	} else {
		sc.selections.Add(path)
	}
}

func (sc *SelectionController) GetSelections() []string {
	return sc.selections.ToSlice()
}

func (sc *SelectionController) UpdateView() {
	sc.view.UpdateView(sc.selections.ToSlice())
}
