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

	sv.SetTitle(0)

	return sv
}

func (sv *SelectionView) SetTitle(selectionsNum int) {
	sv.v.v.Title = fmt.Sprintf(" Selection (%d) ", selectionsNum)
}

func (sv *SelectionView) RenderSelections(selections map[string]struct{}) {
	s := make([]string, 0, len(selections))
	for k := range selections {
		s = append(s, k)
	}

	sv.v.SetViewContent(s)
}
