package message

import (
	"fmt"
	"log"

	"github.com/dinhhuy258/fm/pkg/ctx"
	"github.com/dinhhuy258/fm/pkg/gui/view"
)

func DeleteSelections(ctx *ctx.Context, params ...interface{}) error {
	if len((*ctx).GetState().Selections) == 0 {
		err := (*ctx).GetGui().Views.Log.SetLog("Select nothing!!!", view.LogLevel(view.WARNING))
		if err != nil {
			return err
		}

		return PopMode(ctx)
	}

	onYes := func() {
		paths := make([]string, 0, len((*ctx).GetState().Selections))
		for k := range (*ctx).GetState().Selections {
			paths = append(paths, k)
		}

		if err := deletePaths(ctx, paths); err != nil {
			log.Fatalf("failed to delete paths log %v", err)
		}

		// Clear selections
		for k := range (*ctx).GetState().Selections {
			delete((*ctx).GetState().Selections, k)
		}
	}

	onNo := func() {
		if err := PopMode(ctx); err != nil {
			log.Fatalf("failed to pop mode %v", err)
		}

		if err := (*ctx).GetGui().Views.Main.SetAsCurrentView(); err != nil {
			log.Fatalf("failed to set main as the current view %v", err)
		}

		if err := (*ctx).GetGui().Views.Log.SetLog(
			"Canceled deleting the current file/folder",
			view.LogLevel(view.WARNING)); err != nil {
			log.Fatalf("failed to set log %v", err)
		}
	}

	return (*ctx).GetGui().Views.Confirm.SetConfirmation(
		"Do you want to delete selected paths?",
		onYes,
		onNo,
	)
}

func DeleteCurrent(ctx *ctx.Context, params ...interface{}) error {
	currentNode := (*ctx).GetFileManager().Dir.Nodes[(*ctx).GetState().FocusIdx]

	onYes := func() {
		if err := deletePaths(ctx, []string{currentNode.AbsolutePath}); err != nil {
			log.Fatalf("failed to delete paths log %v", err)
		}
	}

	onNo := func() {
		if err := PopMode(ctx); err != nil {
			log.Fatalf("failed to pop mode %v", err)
		}

		if err := (*ctx).GetGui().Views.Main.SetAsCurrentView(); err != nil {
			log.Fatalf("failed to set main as the current view %v", err)
		}

		if err := (*ctx).GetGui().Views.Log.SetLog("Canceled deleting the current file/folder",
			view.LogLevel(view.WARNING)); err != nil {
			log.Fatalf("failed to set log %v", err)
		}
	}

	return (*ctx).GetGui().Views.Confirm.SetConfirmation(
		"Do you want to delete "+currentNode.RelativePath+"?",
		onYes,
		onNo,
	)
}

func deletePaths(ctx *ctx.Context, paths []string) error {
	if err := PopMode(ctx); err != nil {
		return err
	}

	if err := (*ctx).GetGui().Views.Main.SetAsCurrentView(); err != nil {
		return err
	}

	if err := (*ctx).GetGui().Views.Progress.StartProgress(1); err != nil {
		return err
	}

	(*ctx).GetFileManager().Delete(paths)

	go func() {
		errCount := 0
	loop:
		for {
			select {
			case <-(*ctx).GetFileManager().DeleteCountChan:
				(*ctx).GetGui().Views.Progress.AddCurrent(1)

				break loop
			case <-(*ctx).GetFileManager().DeleteErrChan:
				errCount++
			}
		}

		var err error
		if errCount != 0 {
			err = (*ctx).GetGui().Views.Log.SetLog(
				fmt.Sprintf("Finished to delete %v. Error count: %d", paths, errCount),
				view.LogLevel(view.INFO),
			)
		} else {
			err = (*ctx).GetGui().Views.Log.SetLog(fmt.Sprintf("Finished to delete file %v", paths), view.LogLevel(view.INFO))
		}

		if err != nil {
			log.Fatalf("failed to set log %v", err)
		}

		if err := Refresh(ctx); err != nil {
			log.Fatalf("failed to refresh %v", err)
		}
	}()

	return nil
}
