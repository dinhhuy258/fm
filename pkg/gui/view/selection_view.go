package view

import (
	"fmt"

	"github.com/dinhhuy258/gocui"
)

type SelectionView struct {
	*View
}

func newSelectionView(g *gocui.Gui, v *gocui.View) *SelectionView {
	sv := &SelectionView{
		newView(g, v),
	}

	sv.setTitle(0)

	return sv
}

func (sv *SelectionView) UpdateView(selections []string) {
	sv.setTitle(len(selections))
	sv.SetViewContent(selections)
}

func (sv *SelectionView) setTitle(selectionsNum int) {
	sv.Title = fmt.Sprintf(" Selection (%d) ", selectionsNum)
}
