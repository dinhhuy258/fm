package controller

import "github.com/dinhhuy258/fm/pkg/gui/view"

type HelpController struct {
	keys []string
	msgs []string

	view *view.HelpView
}

func newHelpController() *HelpController {
	return &HelpController{}
}

func (hc *HelpController) SetView(view *view.HelpView) {
	hc.view = view
}

func (hc *HelpController) SetHelp(keys []string, msgs []string) {
	hc.keys = keys
	hc.msgs = msgs

	hc.UpdateView()
}

func (hc *HelpController) UpdateView() {
	hc.view.SetHelp(hc.keys, hc.msgs)
}
