package command

import (
	"fmt"
	"log"

	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui"
	"github.com/dinhhuy258/fm/pkg/gui/view"
)

func PasteSelections(app IApp, params ...interface{}) error {
	operation := params[0].(string)
	if len(app.State().Selections) == 0 {
		gui.GetGui().Views.Log.SetLog("Select nothing!!!", view.LogLevel(view.WARNING))

		return nil
	}

	paths := make([]string, 0, len(app.State().Selections))
	for k := range app.State().Selections {
		paths = append(paths, k)
	}

	paste(app, paths, fs.GetFileManager().Dir.Path, operation)

	// Clear selections
	for k := range app.State().Selections {
		delete(app.State().Selections, k)
	}

	return nil
}

func paste(app IApp, paths []string, dest, operation string) {
	gui.GetGui().Views.Progress.StartProgress(len(paths))

	var countChan chan int

	var errChan chan error

	if operation == "copy" {
		countChan, errChan = fs.GetFileManager().Copy(paths, dest)
	} else {
		countChan, errChan = fs.GetFileManager().Move(paths, dest)
	}

	go func() {
		errCount := 0
	loop:
		for {
			select {
			case <-countChan:
				gui.GetGui().Views.Progress.AddCurrent(1)

				break loop
			case <-errChan:
				errCount++
			}
		}

		if errCount != 0 {
			gui.GetGui().Views.Log.SetLog(
				fmt.Sprintf("Finished to %s %v. Error count: %d", operation, paths, errCount),
				view.LogLevel(view.INFO),
			)
		} else {
			gui.GetGui().Views.Log.SetLog(fmt.Sprintf("Finished to %s %v", operation, paths), view.LogLevel(view.INFO))
		}

		if err := Refresh(app); err != nil {
			log.Fatalf("failed to refresh %v", err)
		}
	}()
}
