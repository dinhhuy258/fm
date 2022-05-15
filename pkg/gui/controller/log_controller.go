package controller

import (
	"fmt"

	"github.com/dinhhuy258/fm/pkg/gui/view"
)

type LogController struct {
	*BaseController

	level view.LogLevel
	msg   string

	view *view.LogView
}

func newLogController(baseController *BaseController, view *view.LogView) *LogController {
	return &LogController{
		BaseController: baseController,
		view:           view,
	}
}

func (lc *LogController) SetLog(level view.LogLevel, msgFormat string, args ...interface{}) {
	lc.level = level
	lc.msg = fmt.Sprintf(msgFormat, args...)

	lc.UpdateView()
}

func (lc *LogController) UpdateView() {
	lc.view.UpdateView(lc.level, lc.msg)
}
