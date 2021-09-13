package app

import (
	"errors"
	"fmt"
	"log"

	"github.com/dinhhuy258/fm/pkg/gui/view"
	"github.com/dinhhuy258/gocui"
)

var ErrInvalidMessageParams = errors.New("invalid message params")

func focusNext(app *App, params ...interface{}) error {
	if app.State.FocusIdx == app.State.NumberOfFiles-1 {
		return nil
	}

	if err := app.Gui.Views.Main.NextCursor(); err != nil {
		return err
	}

	app.State.FocusIdx++

	return app.Gui.Views.Main.RenderDir(app.FileManager.Dir, app.State.Selections, app.State.FocusIdx)
}

func focusPrevious(app *App, params ...interface{}) error {
	if app.State.FocusIdx == 0 {
		return nil
	}

	if err := app.Gui.Views.Main.PreviousCursor(); err != nil {
		return err
	}

	app.State.FocusIdx--

	return app.Gui.Views.Main.RenderDir(app.FileManager.Dir, app.State.Selections, app.State.FocusIdx)
}

func enter(app *App, params ...interface{}) error {
	currentNode := app.FileManager.Dir.Nodes[app.State.FocusIdx]

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
	app.State.History.VisitLast()
	changeDirectory(app, app.State.History.Peek(), false)

	return nil
}

func nextVisitedPath(app *App, params ...interface{}) error {
	app.State.History.VisitNext()
	changeDirectory(app, app.State.History.Peek(), false)

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

		app.State.FocusIdx++
	}

	return nil
}

func toggleSelection(app *App, params ...interface{}) error {
	path := app.FileManager.Dir.Nodes[app.State.FocusIdx].AbsolutePath

	if _, hasPath := app.State.Selections[path]; hasPath {
		delete(app.State.Selections, path)
	} else {
		app.State.Selections[path] = struct{}{}
	}

	app.Gui.Views.Selection.SetTitle(len(app.State.Selections))

	if err := app.Gui.Views.Selection.RenderSelections(app.State.Selections); err != nil {
		return err
	}

	return app.Gui.Views.Main.RenderDir(app.FileManager.Dir, app.State.Selections, app.State.FocusIdx)
}

func clearSelection(app *App, params ...interface{}) error {
	app.State.Selections = make(map[string]struct{})

	app.Gui.Views.Selection.SetTitle(len(app.State.Selections))

	if err := app.Gui.Views.Selection.RenderSelections(app.State.Selections); err != nil {
		return err
	}

	return app.Gui.Views.Main.RenderDir(app.FileManager.Dir, app.State.Selections, app.State.FocusIdx)
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
				view.LogLevel(view.INFO),
			)
		} else {
			err = app.Gui.Views.Log.SetLog(fmt.Sprintf("Finished to delete file %v", paths), view.LogLevel(view.INFO))
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

func deleteSelections(app *App, params ...interface{}) error {
	if len(app.State.Selections) == 0 {
		err := app.Gui.Views.Log.SetLog("Select nothing!!!", view.LogLevel(view.WARNING))
		if err != nil {
			return err
		}

		return popMode(app)
	}

	onYes := func() {
		paths := make([]string, 0, len(app.State.Selections))
		for k := range app.State.Selections {
			paths = append(paths, k)
		}

		if err := deletePaths(app, paths); err != nil {
			log.Fatalf("failed to delete paths log %v", err)
		}

		// Clear selections
		for k := range app.State.Selections {
			delete(app.State.Selections, k)
		}
	}

	onNo := func() {
		if err := popMode(app); err != nil {
			log.Fatalf("failed to pop mode %v", err)
		}

		if err := app.Gui.Views.Main.SetAsCurrentView(); err != nil {
			log.Fatalf("failed to set main as the current view %v", err)
		}

		if err := app.Gui.Views.Log.SetLog(
			"Canceled deleting the current file/folder",
			view.LogLevel(view.WARNING)); err != nil {
			log.Fatalf("failed to set log %v", err)
		}
	}

	if err := app.Gui.Views.Confirm.SetConfirmation(
		"Do you want to delete selected paths?",
		onYes,
		onNo,
	); err != nil {
		return err
	}

	app.onModeChanged()

	return nil
}

func deleteCurrent(app *App, params ...interface{}) error {
	currentNode := app.FileManager.Dir.Nodes[app.State.FocusIdx]

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

		if err := app.Gui.Views.Log.SetLog("Canceled deleting the current file/folder",
			view.LogLevel(view.WARNING)); err != nil {
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
		app.State.History.Push(app.FileManager.Dir.Path)
	}

	app.FileManager.LoadDirectory(path)
}
