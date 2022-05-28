package view

import (
	"fmt"
	"strings"

	"github.com/dinhhuy258/gocui"
)

type SelectionView struct {
	*View
}

func newSelectionView(v *gocui.View) *SelectionView {
	sv := &SelectionView{
		newView(v),
	}

	sv.setTitle(0)

	return sv
}

func (sv *SelectionView) UpdateView(selections []string) {
	sv.setTitle(len(selections))
	sv.SetContent(strings.Join(selections, "\n"))
}

func (sv *SelectionView) setTitle(selectionsNum int) {
	sv.Title = fmt.Sprintf(" Selection (%d) ", selectionsNum)
}
