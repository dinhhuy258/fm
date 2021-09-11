package view

import "github.com/dinhhuy258/gocui"

type ConfirmView struct {
	v *View
}

func newConfirmView(g *gocui.Gui, v *gocui.View) *ConfirmView {
	cv := &ConfirmView{
		v: newView(g, v),
	}

	cv.v.v.Title = " Confirm "
	cv.v.v.Editable = true

	return cv
}

func (cv *ConfirmView) SetConfirmation(ask string) error {
	cv.v.SetViewContent([]string{ask})

	_, err := cv.v.g.SetViewOnTop(cv.v.v.Name())
	if err != nil {
		return err
	}

	return cv.v.v.SetCursor(len(ask), 0)
}
