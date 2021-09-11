package view

import "github.com/dinhhuy258/gocui"

type HelpView struct {
	v *View
}

func newHelpView(g *gocui.Gui, v *gocui.View) *HelpView {
	hv := &HelpView{
		v: newView(g, v),
	}

	hv.v.v.Title = " Help "

	return hv
}

func (hv *HelpView) SetViewContent(displayStrings []string) {
	hv.v.SetViewContent(displayStrings)
}
