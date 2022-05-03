package command

import (
	"fmt"
	"log"

	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui"
	"github.com/dinhhuy258/fm/pkg/gui/view"
)

func PasteSelections(app IApp, params ...interface{}) error {
	operation, _ := params[0].(string)

	selections := app.GetSelections()
	if len(selections) == 0 {
		gui.GetGui().SetLog("Select nothing!!!", view.LogLevel(view.WARNING))

		return nil
	}

	paths := make([]string, 0, len(selections))
	for k := range selections {
		paths = append(paths, k)
	}

	paste(app, paths, fs.GetFileManager().Dir.Path, operation)

	app.ClearSelections()

	return nil
}

func paste(app IApp, paths []string, dest, operation string) {
	gui.GetGui().StartProgress(len(paths))

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
				gui.GetGui().UpdateProgress()

				break loop
			case <-errChan:
				errCount++
			}
		}

		if errCount != 0 {
			gui.GetGui().SetLog(
				fmt.Sprintf("Finished to %s %v. Error count: %d", operation, paths, errCount),
				view.LogLevel(view.INFO),
			)
		} else {
			gui.GetGui().SetLog(fmt.Sprintf("Finished to %s %v", operation, paths), view.LogLevel(view.INFO))
		}

		if err := Refresh(app); err != nil {
			log.Fatalf("failed to refresh %v", err)
		}
	}()
}
