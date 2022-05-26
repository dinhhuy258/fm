package message

import (
	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui/controller"
	"github.com/dinhhuy258/fm/pkg/gui/view"
	"github.com/dinhhuy258/fm/pkg/optional"
)

func DeleteSelections(app IApp, _ ...string) {
	selectionController, _ := app.GetController(controller.Sellection).(*controller.SelectionController)
	logController, _ := app.GetController(controller.Log).(*controller.LogController)

	paths := selectionController.GetSelections()
	if len(paths) == 0 {
		logController.SetLog(view.Warning, "Select nothing!!!")
		logController.UpdateView()

		return
	}

	deletePaths(app, paths)
	// Clear selections after deleting
	selectionController.ClearSelections()
	selectionController.UpdateView()
}

func DeleteCurrent(app IApp, _ ...string) {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)

	entry := explorerController.GetCurrentEntry()

	deletePaths(app, []string{entry.GetPath()})
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
