package gui

import (
	"fmt"
	"strings"

	"github.com/dinhhuy258/gocui"
)

func (gui *Gui) SetViewContent(v *gocui.View, displayStrings []string) {
	gui.g.Update(func(g *gocui.Gui) error {
		v.Clear()
		fmt.Fprint(v, strings.Join(displayStrings, "\n"))

		return nil
	})
}

func (gui *Gui) NextCursor(v *gocui.View) error {
	cx, cy := v.Cursor()
	if err := v.SetCursor(cx, cy+1); err != nil {
		ox, oy := v.Origin()
		if err := v.SetOrigin(ox, oy+1); err != nil {
			return err
		}
	}

	return nil
}

func (gui *Gui) PreviousCursor(v *gocui.View) error {
	cx, cy := v.Cursor()
	if err := v.SetCursor(cx, cy-1); err != nil {
		ox, oy := v.Origin()
		if err := v.SetOrigin(ox, oy-1); err != nil {
			return err
		}
	}

	return nil
}
