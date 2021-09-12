package view

import "github.com/dinhhuy258/gocui"

type LogView struct {
	v *View
}

func newLogView(g *gocui.Gui, v *gocui.View) *LogView {
	lv := &LogView{
		v: newView(g, v),
	}

	lv.v.v.Title = " Logs "

	return lv
}

func (lv *LogView) SetViewOnTop() error {
	return lv.v.SetViewOnTop()
}

func (lv *LogView) SetLog(log string) error {
	lv.v.SetViewContent([]string{log})

	return lv.SetViewOnTop()
}
