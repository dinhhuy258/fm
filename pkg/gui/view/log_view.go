package view

import (
	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/fm/pkg/gui/view/style"
	"github.com/dinhhuy258/gocui"
)

type LogLevel int8

const (
	Info = iota
	Warning
	Error
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

func (lv *LogView) SetViewOnTop() {
	lv.v.SetViewOnTop()
}

func (lv *LogView) SetLog(level LogLevel, log string) {
	var logStyle style.TextStyle

	switch {
	case level == Info:
		log = config.AppConfig.LogInfoFormat + log
		logStyle = style.FromBasicFg(config.AppConfig.LogInfoColor)
	case level == Warning:
		log = config.AppConfig.LogWarningFormat + log
		logStyle = style.FromBasicFg(config.AppConfig.LogWarningColor)
	default:
		log = config.AppConfig.LogErrorFormat + log
		logStyle = style.FromBasicFg(config.AppConfig.LogErrorColor)
	}

	lv.v.SetViewContent([]string{logStyle.Sprint(log)})
	lv.v.SetViewOnTop()
}
