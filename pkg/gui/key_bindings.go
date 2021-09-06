package gui

import (
	"github.com/jroimartin/gocui"
)

func (gui *Gui) setKeyBindings() error {
	if err := gui.g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit
	}); err != nil {
		return err
	}

	for key := 'a'; key <= 'z'; key++ {
		if err := gui.g.SetKeybinding("", key, gocui.ModNone, gui.wrappedHandler(gui.onKey, key)); err != nil {
			return err
		}
	}

	return nil
}

func (gui *Gui) wrappedHandler(f func(key int32) error, key int32) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		return f(key)
	}
}

func (gui *Gui) onKey(key int32) error {
	return nil
}
