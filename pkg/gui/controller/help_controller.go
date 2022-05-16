package controller

import "github.com/dinhhuy258/fm/pkg/gui/view"

type HelpController struct {
	*BaseController

	title string
	keys  []string
	msgs  []string

	view *view.HelpView
}

func newHelpController(baseController *BaseController, view *view.HelpView) *HelpController {
	return &HelpController{
		BaseController: baseController,
		view:           view,
	}
}

func (hc *HelpController) SetHelp(title string, keys []string, msgs []string) {
	hc.title = title
	hc.keys = keys
	hc.msgs = msgs
}

func (hc *HelpController) UpdateView() {
	hc.view.UpdateView(hc.title, hc.keys, hc.msgs)
}
