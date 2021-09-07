package app

import (
	"math"
)

func focusNext(app *App) error {
	ox, oy := app.Gui.Views.Main.Origin()
	scrollHeight := app.Gui.LinesToScrollDown(app.Gui.Views.Main)

	if scrollHeight > 0 {
		if err := app.Gui.Views.Main.SetOrigin(ox, oy+scrollHeight); err != nil {
			return err
		}

		app.Gui.State.Main.SelectedIdx++
	}

	return nil
}

func focusPrevious(app *App) error {
	ox, oy := app.Gui.Views.Main.Origin()
	scrollHeight := 1
	newOy := int(math.Max(0, float64(oy-scrollHeight)))
	app.Gui.State.Main.SelectedIdx--
	app.Gui.State.Main.SelectedIdx = int(math.Max(float64(app.Gui.State.Main.SelectedIdx), 1))

	return app.Gui.Views.Main.SetOrigin(ox, newOy)
}

func enter(app *App) error {
	currentNode := app.FileManager.Dir.Nodes[app.Gui.State.Main.SelectedIdx-1]

	if currentNode.IsDir {
		app.FileManager.LoadDirectory(currentNode.AbsolutePath)
	}

	return nil
}

func back(app *App) error {
	parent := app.FileManager.Dir.Parent()

	app.FileManager.LoadDirectory(parent)

	return nil
}
