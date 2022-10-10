package msg

import (
	"github.com/dinhhuy258/fm/pkg/gui/view"
)

// LogSuccess is a message to log a success message
func LogSuccess(app IApp, params ...string) {
	logMessage := params[0]
	app.SetLog(view.LogInfo, logMessage)
}

// LogWarning is a message to log a warning message
func LogWarning(app IApp, params ...string) {
	logMessage := params[0]
	app.SetLog(view.LogWarning, logMessage)
}

// LogError is a message to log an error message
func LogError(app IApp, params ...string) {
	logMessage := params[0]
	app.SetLog(view.LogError, logMessage)
}
