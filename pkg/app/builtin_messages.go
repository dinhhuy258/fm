package app

import "github.com/dinhhuy258/fm/pkg/gui"

func focusNext(gui *gui.Gui) error {
	idx := gui.State.Main.SelectedIdx + 1

	if idx > gui.State.Main.NumberOfFiles {
		idx = 1
	}

	if err := gui.Views.Main.SetCursor(0, idx); err != nil {
		return err
	}

	gui.State.Main.SelectedIdx = idx

	return nil
}

func focusPrevious(gui *gui.Gui) error {
	idx := gui.State.Main.SelectedIdx - 1

	if idx < 0 {
		idx = gui.State.Main.NumberOfFiles
	}

	if err := gui.Views.Main.SetCursor(0, idx); err != nil {
		return err
	}

	gui.State.Main.SelectedIdx = idx

	return nil
}
