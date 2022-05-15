package command

import (
	"fmt"

	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui/controller"
	"github.com/dinhhuy258/fm/pkg/gui/view"
	"github.com/dinhhuy258/fm/pkg/optional"
)

func DeleteSelections(app IApp, _ ...interface{}) error {
	selectionController, _ := app.GetController(controller.Sellection).(*controller.SelectionController)
	logController, _ := app.GetController(controller.Log).(*controller.LogController)
	inputController, _ := app.GetController(controller.Input).(*controller.InputController)

	paths := selectionController.GetSelections()
	if len(paths) == 0 {
		logController.SetLog(view.Warning, "Select nothing!!!")

		return PopMode(app)
	}

	inputController.SetInput(controller.InputConfirm, "Do you want to delete selected paths?",
		func(ans string) {
			_ = PopMode(app)

			if ans == "y" || ans == "Y" {
				deletePaths(app, paths)
				// Clear selections after deleting
				selectionController.ClearSelections()
			} else {
				logController.SetLog(view.Warning, "Canceled deleting selections files/folders")
			}
		})

	return nil
}

func DeleteCurrent(app IApp, _ ...interface{}) error {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)
	logController, _ := app.GetController(controller.Log).(*controller.LogController)
	inputController, _ := app.GetController(controller.Input).(*controller.InputController)

	entry := explorerController.GetCurrentEntry()

	inputController.SetInput(controller.InputConfirm, fmt.Sprintf("Do you want to delete %s?", entry.GetName()),
		func(ans string) {
			_ = PopMode(app)

			if ans == "y" || ans == "Y" {
				deletePaths(app, []string{entry.GetPath()})
			} else {
				logController.SetLog(view.Warning, "Canceled deleting the current file/folder")
			}
		})

	return nil
}

func deletePaths(app IApp, paths []string) {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)
	logController, _ := app.GetController(controller.Log).(*controller.LogController)
	progressController, _ := app.GetController(controller.Progress).(*controller.ProgressController)

	progressController.StartProgress(len(paths))
	fs.Delete(paths, func() {
		progressController.UpdateProgress()
	}, func(_ error) {
		progressController.UpdateProgress()
	}, func(successCount, errorCount int) {
		if errorCount != 0 {
			logController.SetLog(
				view.Info,
				"Finished to delete %v. Error count: %d", paths, errorCount,
			)
		} else {
			logController.SetLog(view.Info, "Finished to delete file %v", paths)
		}

		focus := getFocus(app, paths)

		if focus < 0 {
			_ = Refresh(app)
		} else {
			entry := explorerController.GetEntry(focus)
			loadDirectory(app, explorerController.GetPath(), optional.NewOptional(entry.GetPath()))
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
