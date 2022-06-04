package msg

import (
	"github.com/dinhhuy258/fm/pkg/gui/controller"
	"github.com/dinhhuy258/fm/pkg/gui/view"
)

// LogSuccess is a message to log a success message
func LogSuccess(app IApp, params ...string) {
	logController, _ := app.GetController(controller.Log).(*controller.LogController)

	logMessage := params[0]
	logController.SetLog(view.Info, logMessage)
	logController.UpdateView()
}

// LogWarning is a message to log a warning message
func LogWarning(app IApp, params ...string) {
	logController, _ := app.GetController(controller.Log).(*controller.LogController)

	logMessage := params[0]
	logController.SetLog(view.Warning, logMessage)
	logController.UpdateView()
}

// LogError is a message to log an error message
func LogError(app IApp, params ...string) {
	logController, _ := app.GetController(controller.Log).(*controller.LogController)

	logMessage := params[0]
	logController.SetLog(view.Error, logMessage)
	logController.UpdateView()
}
