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

	onYes := func() {
		paths := make([]string, 0, len(ctx.State().Selections))
		for k := range ctx.State().Selections {
			paths = append(paths, k)
		}

		if err := deletePaths(ctx, paths); err != nil {
			log.Fatalf("failed to delete paths %v", err)
		}

		// Clear selections
		for k := range ctx.State().Selections {
			delete(ctx.State().Selections, k)
		}
	}

	onNo := func() {
		if err := PopMode(ctx); err != nil {
			log.Fatalf("failed to pop mode %v", err)
		}

		if err := ctx.Gui().Views.Main.SetAsCurrentView(); err != nil {
			log.Fatalf("failed to set main as the current view %v", err)
		}

		ctx.Gui().Views.Log.SetLog("Canceled deleting the current file/folder", view.LogLevel(view.WARNING))
	}

	return ctx.Gui().Views.Confirm.SetConfirmation(
		"Do you want to delete selected paths?",
		onYes,
		onNo,
	)
}

func DeleteCurrent(ctx ctx.Context, _ ...interface{}) error {
	currentNode := ctx.FileManager().Dir.VisibleNodes[ctx.State().FocusIdx]

	onYes := func() {
		if err := deletePaths(ctx, []string{currentNode.AbsolutePath}); err != nil {
			log.Fatalf("failed to delete paths log %v", err)
		}
	}

	onNo := func() {
		if err := PopMode(ctx); err != nil {
			log.Fatalf("failed to pop mode %v", err)
		}

		if err := ctx.Gui().Views.Main.SetAsCurrentView(); err != nil {
			log.Fatalf("failed to set main as the current view %v", err)
		}

		ctx.Gui().Views.Log.SetLog("Canceled deleting the current file/folder", view.LogLevel(view.WARNING))
	}

	return ctx.Gui().Views.Confirm.SetConfirmation(
		"Do you want to delete "+currentNode.RelativePath+"?",
		onYes,
		onNo,
	)
}

func deletePaths(ctx ctx.Context, paths []string) error {
	if err := PopMode(ctx); err != nil {
		return err
	}

	if err := ctx.Gui().Views.Main.SetAsCurrentView(); err != nil {
		return err
	}

	if err := ctx.Gui().Views.Progress.StartProgress(1); err != nil {
		return err
	}

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

	return nil
}
