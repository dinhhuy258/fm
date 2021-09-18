package view

import "github.com/dinhhuy258/gocui"

type ConfirmView struct {
	v     *View
	onYes func()
	onNo  func()
}

func newConfirmView(g *gocui.Gui, v *gocui.View) *ConfirmView {
	cv := &ConfirmView{
		v: newView(g, v),
	}

	cv.v.v.Title = " Confirmation "
	cv.v.v.Editable = true
	cv.v.v.Editor = gocui.EditorFunc(cv.confirmEditor)

	return cv
}

func (cv *ConfirmView) confirmEditor(_ *gocui.View, _ gocui.Key, ch rune, _ gocui.Modifier) {
	if ch != 0 {
		key := string(ch)
		if key == "Y" || key == "y" {
			cv.onYes()

			return
		}
	}

	cv.onNo()
}

func (cv *ConfirmView) SetConfirmation(ask string, onYes func(), onNo func()) error {
	ask = "> " + ask + " (y/n) "
	cv.v.SetViewContent([]string{ask})
	cv.onYes = onYes
	cv.onNo = onNo

	_, err := cv.v.g.SetViewOnTop(cv.v.v.Name())
	if err != nil {
		return err
	}

	if _, err := cv.v.g.SetCurrentView(cv.v.v.Name()); err != nil {
		return err
	}

	return cv.v.v.SetCursor(len(ask), 0)
}
