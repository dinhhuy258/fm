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

func newLogView(v *gocui.View) *LogView {
	lv := &LogView{
		newView(v),
	}

	lv.Title = " Logs "

	return lv
}

func (lv *LogView) UpdateView(level LogLevel, log string) {
	var logStyle style.TextStyle

	switch {
	case level == Info:
		log = config.AppConfig.LogInfoFormat + log
		logStyle = style.FromBasicFg(style.ColorMap[config.AppConfig.LogInfoColor].Foreground)
	case level == Warning:
		log = config.AppConfig.LogWarningFormat + log
		logStyle = style.FromBasicFg(style.ColorMap[config.AppConfig.LogWarningColor].Foreground)
	default:
		log = config.AppConfig.LogErrorFormat + log
		logStyle = style.FromBasicFg(style.ColorMap[config.AppConfig.LogErrorColor].Foreground)
	}

	lv.SetContent(logStyle.Sprint(log))
}
