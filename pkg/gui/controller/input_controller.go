package controller

import "github.com/dinhhuy258/fm/pkg/gui/view"

type InputController struct {
	view *view.InputView
}

func newInputController() *InputController {
	return &InputController{}
}

func (ic *InputController) SetView(view *view.InputView) {
	ic.view = view

	ic.view.SetOnType(ic.onType)
}

func (ic *InputController) SetInput(msg string, onInput func(string)) {
	ic.view.SetInput(msg, func(ans string) {
		// gui.views.Explorer.SetAsCurrentView()

		onInput(ans)
	})
}

func (ic *InputController) onType(content string, event view.KeyEvent) {
}
