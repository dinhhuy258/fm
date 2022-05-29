package controller

import (
	"fmt"

	"github.com/dinhhuy258/fm/pkg/gui/view"
	"github.com/dinhhuy258/fm/pkg/optional"
)

// LogController is a controller for log view
type LogController struct {
	*BaseController

	level view.LogLevel
	msg   string

	view *view.LogView
}

// newLogController creates a new log controller
func newLogController(baseController *BaseController, view *view.LogView) *LogController {
	return &LogController{
		BaseController: baseController,
		view:           view,
	}
}

// SetLog sets the log level and message
func (lc *LogController) SetLog(level view.LogLevel, msgFormat string, args ...interface{}) {
	lc.level = level
	lc.msg = fmt.Sprintf(msgFormat, args...)
}

// SetVisible sets the visibility of the log view
func (lc *LogController) SetVisible(visible bool) {
	lc.view.Visible = visible

	if visible {
		lc.mediator.notify(CursorDisabled, optional.NewEmpty[string]())
	} else {
		lc.mediator.notify(CursorEnabled, optional.NewEmpty[string]())
	}
}

// UpdateView updates the log view
func (lc *LogController) UpdateView() {
	lc.SetVisible(true)
	lc.view.UpdateView(lc.level, lc.msg)
}
