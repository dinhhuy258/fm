package view

import (
	"fmt"

	"github.com/dinhhuy258/gocui"
)

type SelectionView struct {
	v *View
}

func newSelectionView(g *gocui.Gui, v *gocui.View) *SelectionView {
	sv := &SelectionView{
		v: newView(g, v),
	}

	sv.setTitle(0)

	return sv
}

func (sv *SelectionView) RenderSelections(selections []string) {
	sv.setTitle(len(selections))
	sv.v.SetViewContent(selections)
}

func (sv *SelectionView) setTitle(selectionsNum int) {
	sv.v.v.Title = fmt.Sprintf(" Selection (%d) ", selectionsNum)
}
