package msg

import (
	"github.com/dinhhuy258/fm/pkg/gui/key"
	"github.com/dinhhuy258/fm/pkg/gui/view"
)

// LogSuccess is a message to log a success message
func LogSuccess(app IApp, _ key.Key, ctx MessageContext) {
	logMessage := ctx["arg1"].(string)
	app.SetLog(view.LogInfo, logMessage)
}

// LogWarning is a message to log a warning message
func LogWarning(app IApp, _ key.Key, ctx MessageContext) {
	logMessage := ctx["arg1"].(string)
	app.SetLog(view.LogWarning, logMessage)
}

// LogError is a message to log an error message
func LogError(app IApp, key key.Key, ctx MessageContext) {
	logMessage := ctx["arg1"].(string)
	app.SetLog(view.LogError, logMessage)
}
