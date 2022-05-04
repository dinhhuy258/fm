package command

import (
	"fmt"

	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui"
	"github.com/dinhhuy258/fm/pkg/gui/view"
)

func DeleteSelections(app IApp, _ ...interface{}) error {
	appGui := gui.GetGui()

	paths := app.GetSelections()
	if len(paths) == 0 {
		appGui.SetLog("Select nothing!!!", view.LogLevel(view.WARNING))

		return PopMode(app)
	}

	gui.GetGui().SetConfirmation("Do you want to delete selected paths?", func(ans bool) {
		_ = PopMode(app)

		if ans {
			deletePaths(app, paths)
			// Clear selections
			app.ClearSelections()
		} else {
			appGui.SetLog("Canceled deleting the current file/folder", view.LogLevel(view.WARNING))
		}
	})

	return nil
}

func DeleteCurrent(app IApp, _ ...interface{}) error {
	fileExplorer := fs.GetFileExplorer()
	appGui := gui.GetGui()

	entry := fileExplorer.GetEntry(app.GetFocus())

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
	countChan, errChan := fs.GetFileManager().Delete(paths)

	go func() {
		errCount := 0
	loop:
		for {
			select {
			case <-countChan:
				appGui.UpdateProgress()

				if appGui.IsProgressFinished() {
					break loop
				}
			case <-errChan:
				errCount++
				appGui.UpdateProgress()

				if appGui.IsProgressFinished() {
					break loop
				}
			}
		}

		if errCount != 0 {
			appGui.SetLog(
				fmt.Sprintf("Finished to delete %v. Error count: %d", paths, errCount),
				view.LogLevel(view.INFO),
			)
		} else {
			appGui.SetLog(fmt.Sprintf("Finished to delete file %v", paths), view.LogLevel(view.INFO))
		}

		focus := getFocus(app, paths)

		if focus < 0 {
			_ = Refresh(app)
		} else {
			_ = Refresh(app, fileExplorer.GetEntry(focus).GetPath())
		}
	}()
}

func getFocus(app IApp, deletedPaths []string) int {
	fileExplorer := fs.GetFileExplorer()

	paths := make(map[string]struct{})
	for _, deletedPath := range deletedPaths {
		paths[deletedPath] = struct{}{}
	}

	entries := fileExplorer.GetEntries()
	entriesSize := fileExplorer.GetEntriesSize()

	focus := app.GetFocus()
	for {
		if _, hasPath := paths[entries[focus].GetPath()]; !hasPath {
			break
		}

		focus++

		if focus == entriesSize {
			break
		}
	}

	if focus == entriesSize {
		focus = app.GetFocus()

		for {
			if _, hasKey := paths[entries[focus].GetPath()]; !hasKey {
				break
			}

			focus--
			if focus < 0 {
				break
			}
		}
	}

	return focus
}
