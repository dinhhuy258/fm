package msg

import (
	"github.com/dinhhuy258/fm/pkg/gui/controller"
	"github.com/dinhhuy258/fm/pkg/gui/view"
)

func LogSuccess(app IApp, params ...string) {
	logController, _ := app.GetController(controller.Log).(*controller.LogController)

	logMessage := params[0]
	logController.SetLog(view.Info, logMessage)
	logController.UpdateView()
}

func LogError(app IApp, params ...string) {
	logController, _ := app.GetController(controller.Log).(*controller.LogController)

	logMessage := params[0]
	logController.SetLog(view.Error, logMessage)
	logController.UpdateView()
}
