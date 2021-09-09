package app

import "github.com/dinhhuy258/gocui"

func focusNext(app *App) error {
	if app.State.Main.FocusIdx == app.State.Main.NumberOfFiles-1 {
		return nil
	}

	v := app.Gui.Views.Main

	cx, cy := v.Cursor()
	if err := v.SetCursor(cx, cy+1); err != nil {
		ox, oy := v.Origin()
		if err := v.SetOrigin(ox, oy+1); err != nil {
			return err
		}
	}

	app.State.Main.FocusIdx++
	app.Gui.RenderDir(app.FileManager.Dir, app.State.Main.FocusIdx)

	return nil
}

func focusPrevious(app *App) error {
	if app.State.Main.FocusIdx == 0 {
		return nil
	}

	v := app.Gui.Views.Main

	cx, cy := v.Cursor()

	if err := v.SetCursor(cx, cy-1); err != nil {
		ox, oy := v.Origin()
		if err := v.SetOrigin(ox, oy-1); err != nil {
			return err
		}
	}

	app.State.Main.FocusIdx--
	app.Gui.RenderDir(app.FileManager.Dir, app.State.Main.FocusIdx)

	return nil
}

func enter(app *App) error {
	currentNode := app.FileManager.Dir.Nodes[app.State.Main.FocusIdx]

	if currentNode.IsDir {
		changeDirectory(app, currentNode.AbsolutePath, true)
	}

	return nil
}

func back(app *App) error {
	parent := app.FileManager.Dir.Parent()

	changeDirectory(app, parent, true)

	return nil
}

func lastVisitedPath(app *App) error {
	app.History.VisitLast()
	changeDirectory(app, app.History.Peek(), false)

	return nil
}

func nextVisitedPath(app *App) error {
	app.History.VisitNext()
	changeDirectory(app, app.History.Peek(), false)

	return nil
}

func changeDirectory(app *App, path string, saveHistory bool) {
	if saveHistory {
		app.History.Push(app.FileManager.Dir.Path)
	}

	app.FileManager.LoadDirectory(path)
}

func quit(app *App) error {
	return gocui.ErrQuit
}
