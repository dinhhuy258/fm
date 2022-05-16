package command

import (
	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui/controller"
	"github.com/dinhhuy258/fm/pkg/gui/view"
)

func PasteSelections(app IApp, params ...interface{}) {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)
	logController, _ := app.GetController(controller.Log).(*controller.LogController)
	selectionController, _ := app.GetController(controller.Sellection).(*controller.SelectionController)

	operation, _ := params[0].(string)

	paths := selectionController.GetSelections()

	if len(paths) == 0 {
		logController.SetLog(view.Warning, "Select nothing!!!")
		logController.UpdateView()

		return
	}

	paste(app, paths, explorerController.GetPath(), operation)

	selectionController.ClearSelections()
	selectionController.UpdateView()
}

func paste(app IApp, paths []string, dest, operation string) {
	progressController, _ := app.GetController(controller.Progress).(*controller.ProgressController)
	logController, _ := app.GetController(controller.Log).(*controller.LogController)

	progressController.StartProgress(len(paths))
	progressController.UpdateView()

	onSuccess := func() {
		progressController.UpdateProgress()
		progressController.UpdateView()
	}
	onError := func(error) {
		progressController.UpdateProgress()
		progressController.UpdateView()
	}
	onComplete := func(successCount int, errorCount int) {
		if errorCount != 0 {
			logController.SetLog(
				view.Info,
				"Finished to %s %v. Error count: %d", operation, paths, errorCount,
			)
			logController.UpdateView()
		} else {
			logController.SetLog(view.Info, "Finished to %s %v", operation, paths)
			logController.UpdateView()
		}

		Refresh(app)
	}

	if operation == "copy" {
		fs.Copy(paths, dest, onSuccess, onError, onComplete)
	} else {
		fs.Move(paths, dest, onSuccess, onError, onComplete)
	}
}
