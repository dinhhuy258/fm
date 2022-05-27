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
	*View
}

func newLogView(g *gocui.Gui, v *gocui.View) *LogView {
	lv := &LogView{
		newView(g, v),
	}

	lv.SetTitle(" Logs ")
	lv.SetViewOnTop()

	return lv
}

func (lv *LogView) UpdateView(level LogLevel, log string) {
	var logStyle style.TextStyle

	switch {
	case level == Info:
		log = config.AppConfig.LogInfoFormat + log
		logStyle = style.FromBasicFg(style.StringToColor(config.AppConfig.LogInfoColor))
	case level == Warning:
		log = config.AppConfig.LogWarningFormat + log
		logStyle = style.FromBasicFg(style.StringToColor(config.AppConfig.LogWarningColor))
	default:
		log = config.AppConfig.LogErrorFormat + log
		logStyle = style.FromBasicFg(style.StringToColor(config.AppConfig.LogErrorColor))
	}

	lv.SetViewContent([]string{logStyle.Sprint(log)})
	lv.SetViewOnTop()
}
