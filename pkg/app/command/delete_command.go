package command

import (
	"fmt"

	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui"
	"github.com/dinhhuy258/fm/pkg/gui/view"
)

func DeleteSelections(app IApp, _ ...interface{}) error {
	appGui := gui.GetGui()
	selectionController := appGui.GetControllers().Sellection

	paths := selectionController.GetSelections()
	if len(paths) == 0 {
		appGui.SetLog("Select nothing!!!", view.LogLevel(view.WARNING))

		return PopMode(app)
	}

	gui.GetGui().SetConfirmation("Do you want to delete selected paths?", func(ans bool) {
		_ = PopMode(app)

		if ans {
			deletePaths(app, paths)
			// Clear selections after deleting
			selectionController.ClearSelections()
		} else {
			appGui.SetLog("Canceled deleting selections files/folders", view.LogLevel(view.WARNING))
		}
	})

	return nil
}

func DeleteCurrent(app IApp, _ ...interface{}) error {
	appGui := gui.GetGui()
	explorerController := appGui.GetControllers().Explorer

	entry := explorerController.GetCurrentEntry()

	appGui.SetConfirmation(fmt.Sprintf("Do you want to delete %s?", entry.GetName()), func(ans bool) {
		_ = PopMode(app)

		if ans {
			deletePaths(app, []string{entry.GetPath()})
		} else {
			appGui.SetLog("Canceled deleting the current file/folder", view.LogLevel(view.WARNING))
		}
	})

	return nil
}

func deletePaths(app IApp, paths []string) {
	fileExplorer := fs.GetFileExplorer()
	appGui := gui.GetGui()

	appGui.StartProgress(len(paths))
	fs.Delete(paths, func() {
		appGui.UpdateProgress()
	}, func(_ error) {
		appGui.UpdateProgress()
	}, func(successCount, errorCount int) {
		if errorCount != 0 {
			appGui.SetLog(
				fmt.Sprintf("Finished to delete %v. Error count: %d", paths, errorCount),
				view.LogLevel(view.INFO),
			)
		} else {
			appGui.SetLog(fmt.Sprintf("Finished to delete file %v", paths), view.LogLevel(view.INFO))
		}

		focus := getFocus(app, paths)

		if focus < 0 {
			_ = Refresh(app)
		} else {
			entry := fileExplorer.GetEntry(focus)
			LoadDirectory(app, fileExplorer.GetPath(), entry.GetPath())
		}
	})
}

// getFocus re-calculate the focus after deleting files/folders
func getFocus(app IApp, deletedPaths []string) int {
	appGui := gui.GetGui()
	explorerController := appGui.GetControllers().Explorer
	fileExplorer := fs.GetFileExplorer()

	// Put deleted paths to hash map
	paths := make(map[string]struct{})
	for _, deletedPath := range deletedPaths {
		paths[deletedPath] = struct{}{}
	}

	entries := fileExplorer.GetEntries()
	entriesSize := fileExplorer.GetEntriesSize()

	focus := explorerController.GetFocus()
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
