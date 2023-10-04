package controller

import "github.com/dinhhuy258/fm/pkg/gui/view"

// HelpController is a controller for help view.
type HelpController struct {
	*BaseController

	title string
	keys  []string
	msgs  []string

	view *view.HelpView
}

// newHelpController creates a new help controller.
func newHelpController(baseController *BaseController, view *view.HelpView) *HelpController {
	return &HelpController{
		BaseController: baseController,
		view:           view,
	}
}

// SetHelp sets the help information.
func (hc *HelpController) SetHelp(title string, keys []string, msgs []string) {
	hc.title = title
	hc.keys = keys
	hc.msgs = msgs
}

// UpdateView updates the view.
func (hc *HelpController) UpdateView() {
	hc.view.UpdateView(hc.title, hc.keys, hc.msgs)
}
