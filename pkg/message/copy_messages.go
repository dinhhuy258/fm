package message

import (
	"fmt"
	"log"

	"github.com/dinhhuy258/fm/pkg/ctx"
	"github.com/dinhhuy258/fm/pkg/gui/view"
)

func CopySelections(ctx ctx.Context, _ ...interface{}) error {
	if len(ctx.State().Selections) == 0 {
		return ctx.Gui().Views.Log.SetLog("Select nothing!!!", view.LogLevel(view.WARNING))
	}

	paths := make([]string, 0, len(ctx.State().Selections))
	for k := range ctx.State().Selections {
		paths = append(paths, k)
	}

	if err := copyPaths(ctx, paths, ctx.FileManager().Dir.Path); err != nil {
		log.Fatalf("failed to copy %v", err)
	}

	// Clear selections
	for k := range ctx.State().Selections {
		delete(ctx.State().Selections, k)
	}

	return nil
}

func copyPaths(ctx ctx.Context, paths []string, dest string) error {
	if err := ctx.Gui().Views.Progress.StartProgress(len(paths)); err != nil {
		return err
	}

	countChan, errChan := ctx.FileManager().Copy(paths, dest)

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

		var err error
		if errCount != 0 {
			err = ctx.Gui().Views.Log.SetLog(
				fmt.Sprintf("Finished to copy %v. Error count: %d", paths, errCount),
				view.LogLevel(view.INFO),
			)
		} else {
			err = ctx.Gui().Views.Log.SetLog(fmt.Sprintf("Finished to copy %v", paths), view.LogLevel(view.INFO))
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
