package message

import (
	"fmt"
	"log"

	"github.com/dinhhuy258/fm/pkg/ctx"
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
	currentNode := ctx.FileManager().Dir.VisibleNodes[ctx.State().FocusIdx]

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

	countChan, errChan := ctx.FileManager().Delete(paths)

	go func() {
		errCount := 0
	loop:
		for {
			select {
			case <-countChan:
				ctx.Gui().Views.Progress.AddCurrent(1)

				break loop
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

		if err := Refresh(ctx); err != nil {
			log.Fatalf("failed to refresh %v", err)
		}
	}()
}
