package command

import (
	"fmt"

	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui"
	"github.com/dinhhuy258/fm/pkg/gui/view"
)

func DeleteSelections(app IApp, _ ...interface{}) error {
	if len(app.GetSelections()) == 0 {
		gui.GetGui().SetLog("Select nothing!!!", view.LogLevel(view.WARNING))

		return PopMode(app)
	}

	gui.GetGui().SetConfirmation("Do you want to delete selected paths?", func(ans bool) {
		_ = PopMode(app)

		if ans {
			paths := app.GetSelections()
			deletePaths(app, paths)
			// Clear selections
			app.ClearSelections()
		} else {
			gui.GetGui().SetLog("Canceled deleting the current file/folder", view.LogLevel(view.WARNING))
		}
	})

	return nil
}

func DeleteCurrent(app IApp, _ ...interface{}) error {
	currentNode := fs.GetFileManager().GetNodeAtIdx(app.GetFocus())

	gui.GetGui().SetConfirmation("Do you want to delete "+currentNode.RelativePath+"?", func(ans bool) {
		_ = PopMode(app)

		if ans {
			deletePaths(app, []string{currentNode.AbsolutePath})
		} else {
			gui.GetGui().SetLog("Canceled deleting the current file/folder", view.LogLevel(view.WARNING))
		}
	})

	return nil
}

func deletePaths(app IApp, paths []string) {
	gui.GetGui().StartProgress(len(paths))

	countChan, errChan := fs.GetFileManager().Delete(paths)

	go func() {
		errCount := 0
	loop:
		for {
			select {
			case <-countChan:
				gui.GetGui().UpdateProgress()

				if gui.GetGui().IsProgressFinished() {
					break loop
				}
			case <-errChan:
				errCount++
				gui.GetGui().UpdateProgress()

				if gui.GetGui().IsProgressFinished() {
					break loop
				}
			}
		}

		if errCount != 0 {
			gui.GetGui().SetLog(
				fmt.Sprintf("Finished to delete %v. Error count: %d", paths, errCount),
				view.LogLevel(view.INFO),
			)
		} else {
			gui.GetGui().SetLog(fmt.Sprintf("Finished to delete file %v", paths), view.LogLevel(view.INFO))
		}

		focus := getFocusIdx(app, paths)

		if focus < 0 {
			_ = Refresh(app)
		} else {
			_ = Refresh(app, fs.GetFileManager().GetNodeAtIdx(focus).AbsolutePath)
		}
	}()
}

func getFocusIdx(app IApp, paths []string) int {
	pathsMap := make(map[string]struct{})
	for _, path := range paths {
		pathsMap[path] = struct{}{}
	}

	visibleNodes := fs.GetFileManager().GetVisibleNodes()
	visibleNodesSize := len(visibleNodes)
	focus := app.GetFocus()

	for {
		if _, hasKey := pathsMap[visibleNodes[focus].AbsolutePath]; !hasKey {
			break
		}

		focus++

		if focus == visibleNodesSize {
			break
		}
	}

	if focus == visibleNodesSize {
		focus = app.GetFocus()

		for {
			if _, hasKey := pathsMap[visibleNodes[focus].AbsolutePath]; !hasKey {
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
