package command

import (
	"fmt"
	"github.com/dinhhuy258/fm/pkg/app/context"

	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui"
	"github.com/dinhhuy258/fm/pkg/gui/view"
)

func DeleteSelections(ctx context.Context, _ ...interface{}) error {
	if len(ctx.State().Selections) == 0 {
		gui.GetGui().Views.Log.SetLog("Select nothing!!!", view.LogLevel(view.WARNING))

		return PopMode(ctx)
	}

	gui.GetGui().Views.Confirm.SetConfirmation(
		"Do you want to delete selected paths?",
	)

	go func() {
		ans := gui.GetGui().Views.Confirm.GetAnswer()

		_ = PopMode(ctx)

		gui.GetGui().Views.Main.SetAsCurrentView()

		if ans {
			paths := make([]string, 0, len(ctx.State().Selections))
			for k := range ctx.State().Selections {
				paths = append(paths, k)
			}

			deletePaths(ctx, paths)

			// Clear selections
			for k := range ctx.State().Selections {
				delete(ctx.State().Selections, k)
			}
		} else {
			gui.GetGui().Views.Log.SetLog("Canceled deleting the current file/folder", view.LogLevel(view.WARNING))
		}
	}()

	return nil
}

func DeleteCurrent(ctx context.Context, _ ...interface{}) error {
	currentNode := fs.GetFileManager().Dir.VisibleNodes[ctx.State().FocusIdx]

	gui.GetGui().Views.Confirm.SetConfirmation("Do you want to delete " + currentNode.RelativePath + "?")

	go func() {
		ans := gui.GetGui().Views.Confirm.GetAnswer()

		_ = PopMode(ctx)

		gui.GetGui().Views.Main.SetAsCurrentView()

		if ans {
			deletePaths(ctx, []string{currentNode.AbsolutePath})
		} else {
			gui.GetGui().Views.Log.SetLog("Canceled deleting the current file/folder", view.LogLevel(view.WARNING))
		}
	}()

	return nil
}

func deletePaths(ctx context.Context, paths []string) {
	gui.GetGui().Views.Progress.StartProgress(len(paths))

	countChan, errChan := fs.GetFileManager().Delete(paths)

	go func() {
		errCount := 0
	loop:
		for {
			select {
			case <-countChan:
				gui.GetGui().Views.Progress.AddCurrent(1)

				if gui.GetGui().Views.Progress.IsFinished() {
					break loop
				}
			case <-errChan:
				errCount++
			}
		}

		if errCount != 0 {
			gui.GetGui().Views.Log.SetLog(
				fmt.Sprintf("Finished to delete %v. Error count: %d", paths, errCount),
				view.LogLevel(view.INFO),
			)
		} else {
			gui.GetGui().Views.Log.SetLog(fmt.Sprintf("Finished to delete file %v", paths), view.LogLevel(view.INFO))
		}

		focusIdx := getFocusIdx(ctx, paths)

		if focusIdx < 0 {
			_ = Refresh(ctx)
		} else {
			_ = Refresh(ctx, fs.GetFileManager().Dir.VisibleNodes[focusIdx].AbsolutePath)
		}
	}()
}

func getFocusIdx(ctx context.Context, paths []string) int {
	pathsMap := make(map[string]struct{})
	for _, path := range paths {
		pathsMap[path] = struct{}{}
	}

	visibleNodes := fs.GetFileManager().Dir.VisibleNodes
	visibleNodesSize := len(visibleNodes)
	focusIdx := ctx.State().FocusIdx

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
		focusIdx = ctx.State().FocusIdx

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
