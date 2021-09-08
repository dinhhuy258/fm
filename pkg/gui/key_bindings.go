package gui

import (
	"github.com/jroimartin/gocui"
)

func (gui *Gui) setKeyBindings(onKey func(string) error) error {
	if err := gui.g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit
	}); err != nil {
		return err
	}

	for key := 'a'; key <= 'z'; key++ {
		if err := gui.g.SetKeybinding("", key, gocui.ModNone, gui.wrappedHandler(onKey, string(key))); err != nil {
			return err
		}
	}

	return nil
}

func (gui *Gui) wrappedHandler(f func(key string) error, key string) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		return f(key)
	}
}
