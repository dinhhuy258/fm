package command

import (
	"fmt"

	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui"
	"github.com/dinhhuy258/fm/pkg/gui/view"
)

func PasteSelections(app IApp, params ...interface{}) error {
	fileExplorer := fs.GetFileExplorer()
	appGui := gui.GetGui()

	operation, _ := params[0].(string)

	paths := app.GetSelections()
	if len(paths) == 0 {
		appGui.SetLog("Select nothing!!!", view.LogLevel(view.WARNING))

		return nil
	}

	paste(app, paths, fileExplorer.GetPath(), operation)

	app.ClearSelections()

	return nil
}

func paste(app IApp, paths []string, dest, operation string) {
	appGui := gui.GetGui()

	appGui.StartProgress(len(paths))

	var countChan chan int

	var errChan chan error

	if operation == "copy" {
		countChan, errChan = fs.Copy(paths, dest)
	} else {
		countChan, errChan = fs.Move(paths, dest)
	}

	go func() {
		errCount := 0
	loop:
		for {
			select {
			case <-countChan:
				appGui.UpdateProgress()

				break loop
			case <-errChan:
				errCount++
			}
		}

		if errCount != 0 {
			appGui.SetLog(
				fmt.Sprintf("Finished to %s %v. Error count: %d", operation, paths, errCount),
				view.LogLevel(view.INFO),
			)
		} else {
			appGui.SetLog(fmt.Sprintf("Finished to %s %v", operation, paths), view.LogLevel(view.INFO))
		}

		_ = Refresh(app)
	}()
}
