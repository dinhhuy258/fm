package app

import (
	"github.com/dinhhuy258/gocui"
)

func focusNext(app *App) error {
	if app.State.Main.FocusIdx == app.State.Main.NumberOfFiles-1 {
		return nil
	}

	if err := app.Gui.Views.Main.NextCursor(); err != nil {
		return err
	}

	app.State.Main.FocusIdx++

	return app.Gui.Views.Main.RenderDir(app.FileManager.Dir, app.State.Selections, app.State.Main.FocusIdx)
}

func focusPrevious(app *App) error {
	if app.State.Main.FocusIdx == 0 {
		return nil
	}

	if err := app.Gui.Views.Main.PreviousCursor(); err != nil {
		return err
	}

	app.State.Main.FocusIdx--

	return app.Gui.Views.Main.RenderDir(app.FileManager.Dir, app.State.Selections, app.State.Main.FocusIdx)
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

func focus(app *App, path string) error {
	count := 0

	for _, node := range app.FileManager.Dir.Nodes {
		if node.IsDir && node.AbsolutePath == path {
			break
		}

		count++
	}

	if count == len(app.FileManager.Dir.Nodes) {
		return nil
	}

	for i := 0; i < count; i++ {
		if err := app.Gui.Views.Main.NextCursor(); err != nil {
			return err
		}

		app.State.Main.FocusIdx++
	}

	return nil
}

func toggleSelection(app *App) error {
	path := app.FileManager.Dir.Nodes[app.State.Main.FocusIdx].AbsolutePath

	if _, hasPath := app.State.Selections[path]; hasPath {
		delete(app.State.Selections, path)
	} else {
		app.State.Selections[path] = struct{}{}
	}

	app.Gui.Views.Selection.SetTitle(len(app.State.Selections))

	if err := app.Gui.Views.Selection.RenderSelections(app.State.Selections); err != nil {
		return err
	}

	return app.Gui.Views.Main.RenderDir(app.FileManager.Dir, app.State.Selections, app.State.Main.FocusIdx)
}

func clearSelection(app *App) error {
	app.State.Selections = make(map[string]struct{})

	app.Gui.Views.Selection.SetTitle(len(app.State.Selections))

	if err := app.Gui.Views.Selection.RenderSelections(app.State.Selections); err != nil {
		return err
	}

	return app.Gui.Views.Main.RenderDir(app.FileManager.Dir, app.State.Selections, app.State.Main.FocusIdx)
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
