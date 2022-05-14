package controller

import "github.com/dinhhuy258/fm/pkg/gui/view"

type InputController struct {
	*BaseController
	view *view.InputView
}

func newInputController(baseController *BaseController) *InputController {
	return &InputController{
		BaseController: baseController,
	}
}

func (ic *InputController) SetView(view *view.InputView) {
	ic.view = view

	ic.view.SetOnType(ic.onType)
}

func (ic *InputController) SetInput(msg string, onInput func(string)) {
	ic.view.SetInput(msg, func(ans string) {
		ic.mediator.notify(INPUT_DONE, nil)

		onInput(ans)
	})
}

func (ic *InputController) onType(content string, event view.KeyEvent) {
}
