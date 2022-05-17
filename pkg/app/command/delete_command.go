package command

import (
	"fmt"

	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui/controller"
	"github.com/dinhhuy258/fm/pkg/gui/view"
	"github.com/dinhhuy258/fm/pkg/optional"
)

func DeleteSelections(app IApp, _ ...interface{}) {
	selectionController, _ := app.GetController(controller.Sellection).(*controller.SelectionController)
	logController, _ := app.GetController(controller.Log).(*controller.LogController)
	inputController, _ := app.GetController(controller.Input).(*controller.InputController)

	paths := selectionController.GetSelections()
	if len(paths) == 0 {
		logController.SetLog(view.Warning, "Select nothing!!!")
		logController.UpdateView()

		return
	}

	inputController.SetInput(controller.InputConfirm, "Do you want to delete selected paths?",
		optional.NewEmpty[string](),
		func(ans string) {
			if ans == "y" || ans == "Y" {
				deletePaths(app, paths)
				// Clear selections after deleting
				selectionController.ClearSelections()
				selectionController.UpdateView()
			} else {
				logController.SetLog(view.Warning, "Canceled deleting selections files/folders")
				logController.UpdateView()
			}
		})
	inputController.UpdateView()
}

func DeleteCurrent(app IApp, _ ...interface{}) {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)
	logController, _ := app.GetController(controller.Log).(*controller.LogController)
	inputController, _ := app.GetController(controller.Input).(*controller.InputController)

	entry := explorerController.GetCurrentEntry()

	inputController.SetInput(controller.InputConfirm, fmt.Sprintf("Do you want to delete %s?", entry.GetName()),
		optional.NewEmpty[string](),
		func(ans string) {
			if ans == "y" || ans == "Y" {
				deletePaths(app, []string{entry.GetPath()})
			} else {
				logController.SetLog(view.Warning, "Canceled deleting the current file/folder")
				logController.UpdateView()
			}
		})
	inputController.UpdateView()
}

func deletePaths(app IApp, paths []string) {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)
	logController, _ := app.GetController(controller.Log).(*controller.LogController)
	progressController, _ := app.GetController(controller.Progress).(*controller.ProgressController)

	progressController.StartProgress(len(paths))
	progressController.UpdateView()
	fs.Delete(paths, func() {
		progressController.UpdateProgress()
		progressController.UpdateView()
	}, func(_ error) {
		progressController.UpdateProgress()
		progressController.UpdateView()
	}, func(successCount, errorCount int) {
		if errorCount != 0 {
			logController.SetLog(
				view.Info,
				"Finished to delete %v. Error count: %d", paths, errorCount,
			)
			logController.UpdateView()
		} else {
			logController.SetLog(view.Info, "Finished to delete file %v", paths)
			logController.UpdateView()
		}

		focus := getFocus(app, paths)

		if focus < 0 {
			Refresh(app)
		} else {
			entry := explorerController.GetEntry(focus)
			loadDirectory(app, explorerController.GetPath(), optional.New(entry.GetPath()))
		}
	})
}

// getFocus re-calculate the focus after deleting files/folders
func getFocus(app IApp, deletedPaths []string) int {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)

	// Put deleted paths to hash map
	paths := make(map[string]struct{})
	for _, deletedPath := range deletedPaths {
		paths[deletedPath] = struct{}{}
	}

	focus := explorerController.GetFocus()
	entries := explorerController.GetEntries()
	entriesSize := len(entries)

	// Move the focus until it focus to non-deleted paths
	for focus < entriesSize {
		if _, deleted := paths[entries[focus].GetPath()]; !deleted {
			break
		}

		focus++
	}

	if focus == entriesSize {
		// In case, all paths below the current focus was deleted, let's try with paths above
		focus = explorerController.GetFocus()

		for focus < 0 {
			if _, deleted := paths[entries[focus].GetPath()]; !deleted {
				break
			}

			focus--
		}
	}

	return focus
}
