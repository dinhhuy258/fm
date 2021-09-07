package app

import "github.com/dinhhuy258/fm/pkg/gui"

func focusNext(gui *gui.Gui) {
	gui.SetViewContent(gui.Views.Main, []string{"NEXT"})
}

func focusPrevious(gui *gui.Gui) {
	gui.SetViewContent(gui.Views.Main, []string{"PREVIOUS"})
}
