package message

import (
	"fmt"

	"github.com/dinhhuy258/fm/pkg/ctx"
	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui/view"
)

func DeleteSelections(ctx ctx.Context, _ ...interface{}) error {
	if len(ctx.State().Selections) == 0 {
		ctx.Gui().Views.Log.SetLog("Select nothing!!!", view.LogLevel(view.WARNING))

		return PopMode(ctx)
	}

	ctx.Gui().Views.Confirm.SetConfirmation(
		"Do you want to delete selected paths?",
	)

	go func() {
		ans := ctx.Gui().Views.Confirm.GetAnswer()

		_ = PopMode(ctx)

		ctx.Gui().Views.Main.SetAsCurrentView()

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
			ctx.Gui().Views.Log.SetLog("Canceled deleting the current file/folder", view.LogLevel(view.WARNING))
		}
	}()

	return nil
}

func DeleteCurrent(ctx ctx.Context, _ ...interface{}) error {
	currentNode := fs.GetFileManager().Dir.VisibleNodes[ctx.State().FocusIdx]

	ctx.Gui().Views.Confirm.SetConfirmation("Do you want to delete " + currentNode.RelativePath + "?")

	go func() {
		ans := ctx.Gui().Views.Confirm.GetAnswer()

		_ = PopMode(ctx)

		ctx.Gui().Views.Main.SetAsCurrentView()

		if ans {
			deletePaths(ctx, []string{currentNode.AbsolutePath})
		} else {
			ctx.Gui().Views.Log.SetLog("Canceled deleting the current file/folder", view.LogLevel(view.WARNING))
		}
	}()

	return nil
}

func deletePaths(ctx ctx.Context, paths []string) {
	ctx.Gui().Views.Progress.StartProgress(len(paths))

	countChan, errChan := fs.GetFileManager().Delete(paths)

	go func() {
		errCount := 0
	loop:
		for {
			select {
			case <-countChan:
				ctx.Gui().Views.Progress.AddCurrent(1)

				if ctx.Gui().Views.Progress.IsFinished() {
					break loop
				}
			case <-errChan:
				errCount++
			}
		}

		if errCount != 0 {
			ctx.Gui().Views.Log.SetLog(
				fmt.Sprintf("Finished to delete %v. Error count: %d", paths, errCount),
				view.LogLevel(view.INFO),
			)
		} else {
			ctx.Gui().Views.Log.SetLog(fmt.Sprintf("Finished to delete file %v", paths), view.LogLevel(view.INFO))
		}

		focusIdx := getFocusIdx(ctx, paths)

		if focusIdx < 0 {
			_ = Refresh(ctx)
		} else {
			_ = Refresh(ctx, fs.GetFileManager().Dir.VisibleNodes[focusIdx].AbsolutePath)
		}
	}()
}

func getFocusIdx(ctx ctx.Context, paths []string) int {
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
