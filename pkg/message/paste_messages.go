package message

import (
	"fmt"
	"log"

	"github.com/dinhhuy258/fm/pkg/ctx"
	"github.com/dinhhuy258/fm/pkg/gui/view"
)

func PasteSelections(ctx ctx.Context, params ...interface{}) error {
	if len(params) != 1 {
		return ErrInvalidMessageParameter
	}

	operation, ok := params[0].(string)
	if !ok || (operation != "cut" && operation != "copy") {
		return ErrInvalidMessageParameter
	}

	if len(ctx.State().Selections) == 0 {
		ctx.Gui().Views.Log.SetLog("Select nothing!!!", view.LogLevel(view.WARNING))

		return nil
	}

	paths := make([]string, 0, len(ctx.State().Selections))
	for k := range ctx.State().Selections {
		paths = append(paths, k)
	}

	paste(ctx, paths, ctx.FileManager().Dir.Path, operation)

	// Clear selections
	for k := range ctx.State().Selections {
		delete(ctx.State().Selections, k)
	}

	return nil
}

func paste(ctx ctx.Context, paths []string, dest, operation string) {
	ctx.Gui().Views.Progress.StartProgress(len(paths))

	var countChan chan int

	var errChan chan error

	if operation == "copy" {
		countChan, errChan = ctx.FileManager().Copy(paths, dest)
	} else {
		countChan, errChan = ctx.FileManager().Move(paths, dest)
	}

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
				fmt.Sprintf("Finished to %s %v. Error count: %d", operation, paths, errCount),
				view.LogLevel(view.INFO),
			)
		} else {
			ctx.Gui().Views.Log.SetLog(fmt.Sprintf("Finished to %s %v", operation, paths), view.LogLevel(view.INFO))
		}

		if err := Refresh(ctx); err != nil {
			log.Fatalf("failed to refresh %v", err)
		}
	}()
}
