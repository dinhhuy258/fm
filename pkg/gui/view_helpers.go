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
