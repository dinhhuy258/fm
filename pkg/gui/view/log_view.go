package view

import (
	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/fm/pkg/style"
	"github.com/dinhhuy258/gocui"
)

type LogLevel byte

const (
	INFO    = 0
	WARNING = 1
	ERROR   = 2
)

type LogView struct {
	v *View
}

func newLogView(g *gocui.Gui, v *gocui.View) *LogView {
	lv := &LogView{
		v: newView(g, v),
	}

	lv.v.v.Title = " Logs "
	lv.v.SetViewOnTop()

	return lv
}

func (lv *LogView) SetLog(log string, level LogLevel) {
	var logStyle style.TextStyle

	switch {
	case level == INFO:
		log = config.AppConfig.LogInfoFormat + log
		logStyle = config.AppConfig.LogInfoStyle
	case level == WARNING:
		log = config.AppConfig.LogWarningFormat + log
		logStyle = config.AppConfig.LogWarningStyle
	default:
		log = config.AppConfig.LogErrorFormat + log
		logStyle = config.AppConfig.LogErrorStyle
	}

	lv.v.SetViewContent([]string{logStyle.Sprint(log)})
	lv.v.SetViewOnTop()
}
