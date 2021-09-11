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
	if _, err := lv.v.g.SetViewOnTop(lv.v.v.Name()); err != nil {
		return err
	}

	return nil
}
