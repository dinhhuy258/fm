package view

import (
	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/fm/pkg/gui/view/style"
	"github.com/dinhhuy258/gocui"
)

type LogLevel int8

const (
	LogInfo = iota
	LogWarning
	LogError
)

type LogView struct {
	*View

	logInfoTextStyle    style.TextStyle
	logWarningTextStyle style.TextStyle
	logErrorTextStyle   style.TextStyle
}

func newLogView(v *gocui.View) *LogView {
	cfg := config.AppConfig

	lv := &LogView{
		View: newView(v),
	}

	lv.Title = " Logs "
	lv.logInfoTextStyle = style.FromStyleConfig(cfg.General.LogInfoUI.Style)
	lv.logWarningTextStyle = style.FromStyleConfig(cfg.General.LogWarningUI.Style)
	lv.logErrorTextStyle = style.FromStyleConfig(cfg.General.LogErrorUI.Style)

	return lv
}

func (lv *LogView) UpdateView(level LogLevel, log string) {
	cfg := config.AppConfig

	var logStyle style.TextStyle

	switch {
	case level == LogInfo:
		log = cfg.General.LogInfoUI.Prefix + log + cfg.General.LogInfoUI.Suffix
		logStyle = lv.logInfoTextStyle
	case level == LogWarning:
		log = cfg.General.LogWarningUI.Prefix + log + cfg.General.LogWarningUI.Suffix
		logStyle = lv.logWarningTextStyle
	default:
		// Error
		log = cfg.General.LogErrorUI.Prefix + log + cfg.General.LogErrorUI.Suffix
		logStyle = lv.logErrorTextStyle
	}

	lv.SetContent(logStyle.Sprint(log))
}
