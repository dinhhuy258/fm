package command

import (
	"fmt"

	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui/controller"
	"github.com/dinhhuy258/fm/pkg/gui/view"
)

func DeleteSelections(app IApp, _ ...interface{}) error {
	appGui := app.GetGui()
	selectionController := appGui.GetControllers().Sellection
	logController := appGui.GetControllers().Log
	inputController := appGui.GetControllers().Input

	paths := selectionController.GetSelections()
	if len(paths) == 0 {
		logController.SetLog(view.LogLevel(view.WARNING), "Select nothing!!!")

		return PopMode(app)
	}

	inputController.SetInput(controller.CONFIRM, "Do you want to delete selected paths?",
		func(ans string) {
			_ = PopMode(app)

			if ans == "y" || ans == "Y" {
				deletePaths(app, paths)
				// Clear selections after deleting
				selectionController.ClearSelections()
			} else {
				logController.SetLog(view.LogLevel(view.WARNING), "Canceled deleting selections files/folders")
			}
		})

	return nil
}

func DeleteCurrent(app IApp, _ ...interface{}) error {
	appGui := app.GetGui()
	explorerController := appGui.GetControllers().Explorer
	logController := appGui.GetControllers().Log
	inputController := appGui.GetControllers().Input

	entry := explorerController.GetCurrentEntry()

	inputController.SetInput(controller.CONFIRM, fmt.Sprintf("Do you want to delete %s?", entry.GetName()),
		func(ans string) {
			_ = PopMode(app)

			if ans == "y" || ans == "Y" {
				deletePaths(app, []string{entry.GetPath()})
			} else {
				logController.SetLog(view.LogLevel(view.WARNING), "Canceled deleting the current file/folder")
			}
		})

	return nil
}

func deletePaths(app IApp, paths []string) {
	appGui := app.GetGui()
	progressController := appGui.GetControllers().Progress
	logController := appGui.GetControllers().Log
	explorerController := appGui.GetControllers().Explorer

	progressController.StartProgress(len(paths))
	fs.Delete(paths, func() {
		progressController.UpdateProgress()
	}, func(_ error) {
		progressController.UpdateProgress()
	}, func(successCount, errorCount int) {
		if errorCount != 0 {
			logController.SetLog(
				view.LogLevel(view.INFO),
				"Finished to delete %v. Error count: %d", paths, errorCount,
			)
		} else {
			logController.SetLog(view.LogLevel(view.INFO), "Finished to delete file %v", paths)
		}

		focus := getFocus(app, paths)

		if focus < 0 {
			_ = Refresh(app)
		} else {
			entry := explorerController.GetEntry(focus)
			LoadDirectory(app, explorerController.GetPath(), entry.GetPath())
		}
	})
}

// getFocus re-calculate the focus after deleting files/folders
func getFocus(app IApp, deletedPaths []string) int {
	appGui := app.GetGui()
	explorerController := appGui.GetControllers().Explorer

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
