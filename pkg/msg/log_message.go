package msg

import (
	"github.com/dinhhuy258/fm/pkg/gui/key"
	"github.com/dinhhuy258/fm/pkg/gui/view"
)

// ClearLog is a message to log a success message
func ClearLog(app IApp, _ *key.Key, _ MessageContext) {
	app.SetLog(view.LogInfo, "")
}

// LogSuccess is a message to log a success message
func LogSuccess(app IApp, _ *key.Key, ctx MessageContext) {
	logMessage, _ := ctx["arg1"].(string)
	app.SetLog(view.LogInfo, logMessage)
}

// LogWarning is a message to log a warning message
func LogWarning(app IApp, _ *key.Key, ctx MessageContext) {
	logMessage, _ := ctx["arg1"].(string)
	app.SetLog(view.LogWarning, logMessage)
}

// LogError is a message to log an error message
func LogError(app IApp, _ *key.Key, ctx MessageContext) {
	logMessage, _ := ctx["arg1"].(string)
	app.SetLog(view.LogError, logMessage)
}
