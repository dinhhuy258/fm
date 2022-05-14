package controller

import (
	"fmt"

	"github.com/dinhhuy258/fm/pkg/gui/view"
)

type LogController struct {
	view *view.LogView
}

func newLogController() *LogController {
	return &LogController{}
}

func (lc *LogController) SetView(view *view.LogView) {
	lc.view = view
}

func (lc *LogController) SetLog(level view.LogLevel, msgFormat string, args ...interface{}) {
	lc.view.SetLog(level, fmt.Sprintf(msgFormat, args...))
}

//TODO: Remove
func (lc *LogController) SetViewOnTop() {
	lc.view.SetViewOnTop()
}
