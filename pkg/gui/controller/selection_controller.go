package controller

import (
	set "github.com/deckarep/golang-set/v2"
	"github.com/dinhhuy258/fm/pkg/gui/view"
)

type SelectionController struct {
	selections set.Set[string]

	view *view.SelectionView
}

func newSelectionController(selections set.Set[string]) *SelectionController {
	return &SelectionController{
		selections: selections,
	}
}

func (sc *SelectionController) SetView(view *view.SelectionView) {
	sc.view = view
}

func (sc *SelectionController) ClearSelections() {
	sc.selections.Clear()

	sc.UpdateView()
}

func (sc *SelectionController) ToggleSelection(path string) {
	if sc.selections.Contains(path) {
		sc.selections.Remove(path)
	} else {
		sc.selections.Add(path)
	}

	sc.UpdateView()
}

func (sc *SelectionController) UpdateView() {
	sc.view.RenderSelections(sc.selections.ToSlice())
}
