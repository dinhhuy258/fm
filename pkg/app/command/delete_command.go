package command

import (
	"fmt"

	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui"
	"github.com/dinhhuy258/fm/pkg/gui/view"
)

func DeleteSelections(app IApp, _ ...interface{}) error {
	if len(app.State().Selections) == 0 {
		gui.GetGui().SetLog("Select nothing!!!", view.LogLevel(view.WARNING))

		return PopMode(app)
	}

	gui.GetGui().Views.Confirm.SetConfirmation(
		"Do you want to delete selected paths?",
	)

	go func() {
		ans := gui.GetGui().Views.Confirm.GetAnswer()

		_ = PopMode(app)

		gui.GetGui().Views.Main.SetAsCurrentView()

		if ans {
			paths := make([]string, 0, len(app.State().Selections))
			for k := range app.State().Selections {
				paths = append(paths, k)
			}

			deletePaths(app, paths)

			// Clear selections
			for k := range app.State().Selections {
				delete(app.State().Selections, k)
			}
		} else {
			gui.GetGui().SetLog("Canceled deleting the current file/folder", view.LogLevel(view.WARNING))
		}
	}()

	return nil
}

func DeleteCurrent(app IApp, _ ...interface{}) error {
	currentNode := fs.GetFileManager().Dir.VisibleNodes[app.State().FocusIdx]

	gui.GetGui().Views.Confirm.SetConfirmation("Do you want to delete " + currentNode.RelativePath + "?")

	go func() {
		ans := gui.GetGui().Views.Confirm.GetAnswer()

		_ = PopMode(app)

		gui.GetGui().Views.Main.SetAsCurrentView()

		if ans {
			deletePaths(app, []string{currentNode.AbsolutePath})
		} else {
			gui.GetGui().SetLog("Canceled deleting the current file/folder", view.LogLevel(view.WARNING))
		}
	}()

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

		focusIdx := getFocusIdx(app, paths)

		if focusIdx < 0 {
			_ = Refresh(app)
		} else {
			_ = Refresh(app, fs.GetFileManager().Dir.VisibleNodes[focusIdx].AbsolutePath)
		}
	}()
}

func getFocusIdx(app IApp, paths []string) int {
	pathsMap := make(map[string]struct{})
	for _, path := range paths {
		pathsMap[path] = struct{}{}
	}

	visibleNodes := fs.GetFileManager().Dir.VisibleNodes
	visibleNodesSize := len(visibleNodes)
	focusIdx := app.State().FocusIdx

	for {
		if _, hasKey := pathsMap[visibleNodes[focusIdx].AbsolutePath]; !hasKey {
			break
		}

		focusIdx++

		if focusIdx == visibleNodesSize {
			break
		}
	}

	if focusIdx == visibleNodesSize {
		focusIdx = app.State().FocusIdx

		for {
			if _, hasKey := pathsMap[visibleNodes[focusIdx].AbsolutePath]; !hasKey {
				break
			}

			focusIdx--
			if focusIdx < 0 {
				break
			}
		}
	}

	return focusIdx
}
