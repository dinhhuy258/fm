package gui

import (
	"fmt"
	"strings"

	"github.com/jroimartin/gocui"
)

func (gui *Gui) SetViewContent(v *gocui.View, displayStrings []string) {
	gui.g.Update(func(g *gocui.Gui) error {
		v.Clear()
		fmt.Fprint(v, strings.Join(displayStrings, "\n"))

		return nil
	})
}

func (gui *Gui) LinesToScrollDown(view *gocui.View) int {
	_, oy := view.Origin()
	y := oy

	scrollHeight := 1
	viewLinesHeight := strings.Count(view.ViewBuffer(), "\n")
	scrollableLines := viewLinesHeight - y

	if scrollableLines < 0 {
		return 0
	}

	margin := 1

	if scrollableLines-margin < scrollHeight {
		scrollHeight = scrollableLines - margin
	}

	if oy+scrollHeight < 0 {
		return 0
	}

	return scrollHeight
}
