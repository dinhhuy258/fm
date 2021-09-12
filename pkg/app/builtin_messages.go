package app

import (
	"errors"
	"fmt"
	"log"

	"github.com/dinhhuy258/gocui"
)

var ErrInvalidMessageParams = errors.New("invalid message params")

func focusNext(app *App, params ...interface{}) error {
	if app.State.Main.FocusIdx == app.State.Main.NumberOfFiles-1 {
		return nil
	}

	if err := app.Gui.Views.Main.NextCursor(); err != nil {
		return err
	}

	app.State.Main.FocusIdx++

	return app.Gui.Views.Main.RenderDir(app.FileManager.Dir, app.State.Selections, app.State.Main.FocusIdx)
}

func focusPrevious(app *App, params ...interface{}) error {
	if app.State.Main.FocusIdx == 0 {
		return nil
	}

	if err := app.Gui.Views.Main.PreviousCursor(); err != nil {
		return err
	}

	app.State.Main.FocusIdx--

	return app.Gui.Views.Main.RenderDir(app.FileManager.Dir, app.State.Selections, app.State.Main.FocusIdx)
}

func enter(app *App, params ...interface{}) error {
	currentNode := app.FileManager.Dir.Nodes[app.State.Main.FocusIdx]

	if currentNode.IsDir {
		changeDirectory(app, currentNode.AbsolutePath, true)
	}

	return nil
}

func back(app *App, params ...interface{}) error {
	parent := app.FileManager.Dir.Parent()

	changeDirectory(app, parent, true)

	return nil
}

func lastVisitedPath(app *App, params ...interface{}) error {
	app.History.VisitLast()
	changeDirectory(app, app.History.Peek(), false)

	return nil
}

func nextVisitedPath(app *App, params ...interface{}) error {
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

func toggleSelection(app *App, params ...interface{}) error {
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

func clearSelection(app *App, params ...interface{}) error {
	app.State.Selections = make(map[string]struct{})

	app.Gui.Views.Selection.SetTitle(len(app.State.Selections))

	if err := app.Gui.Views.Selection.RenderSelections(app.State.Selections); err != nil {
		return err
	}

	return app.Gui.Views.Main.RenderDir(app.FileManager.Dir, app.State.Selections, app.State.Main.FocusIdx)
}

func switchMode(app *App, params ...interface{}) error {
	if len(params) != 1 {
		return ErrInvalidMessageParams
	}

	if err := app.Modes.Push(params[0].(string)); err != nil {
		return err
	}

	app.onModeChanged()

	return nil
}

func popMode(app *App, params ...interface{}) error {
	if err := app.Modes.Pop(); err != nil {
		return err
	}

	app.onModeChanged()

	return nil
}

func refresh(app *App, params ...interface{}) error {
	app.FileManager.Reload()

	return nil
}

func deletePaths(app *App, paths []string) error {
	if err := popMode(app); err != nil {
		return err
	}

	if err := app.Gui.Views.Main.SetAsCurrentView(); err != nil {
		return err
	}

	if err := app.Gui.Views.Progress.StartProgress(1); err != nil {
		return err
	}

	app.FileManager.Delete(paths)

	go func() {
		errCount := 0
	loop:
		for {
			select {
			case <-app.FileManager.DeleteCountChan:
				app.Gui.Views.Progress.AddCurrent(1)

				break loop
			case <-app.FileManager.DeleteErrChan:
				errCount++
			}
		}

		var err error
		if errCount != 0 {
			err = app.Gui.Views.Log.SetLog(
				fmt.Sprintf("Finished to delete %v. Error count: %d", paths, errCount),
			)
		} else {
			err = app.Gui.Views.Log.SetLog(fmt.Sprintf("Finished to delete file %v", paths))
		}

		if err != nil {
			log.Fatalf("failed to set log %v", err)
		}

		if err := refresh(app); err != nil {
			log.Fatalf("failed to refresh %v", err)
		}
	}()

	return nil
}

func deleteCurrent(app *App, params ...interface{}) error {
	currentNode := app.FileManager.Dir.Nodes[app.State.Main.FocusIdx]

	onYes := func() {
		if err := deletePaths(app, []string{currentNode.AbsolutePath}); err != nil {
			log.Fatalf("failed to delete paths log %v", err)
		}
	}

	onNo := func() {
		if err := popMode(app); err != nil {
			log.Fatalf("failed to pop mode %v", err)
		}

		if err := app.Gui.Views.Main.SetAsCurrentView(); err != nil {
			log.Fatalf("failed to set main as the current view %v", err)
		}

		if err := app.Gui.Views.Log.SetLog("Canceled deleting the current file/folder"); err != nil {
			log.Fatalf("failed to set log %v", err)
		}
	}

	if err := app.Gui.Views.Confirm.SetConfirmation(
		"Do you want to delete "+currentNode.RelativePath+"?",
		onYes,
		onNo,
	); err != nil {
		return err
	}

	app.onModeChanged()

	return nil
}

func quit(app *App, params ...interface{}) error {
	return gocui.ErrQuit
}

func changeDirectory(app *App, path string, saveHistory bool) {
	if saveHistory {
		app.History.Push(app.FileManager.Dir.Path)
	}

	app.FileManager.LoadDirectory(path)
}
