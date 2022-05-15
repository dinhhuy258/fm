package controller

import (
	"fmt"

	"github.com/dinhhuy258/fm/pkg/gui/view"
)

type InputType int8

const (
	INPUT InputType = iota
	CONFIRM
)

type InputController struct {
	*BaseController

	inputType InputType
	onConfirm func(string)

	view *view.InputView
}

func newInputController(baseController *BaseController, view *view.InputView) *InputController {
	return &InputController{
		BaseController: baseController,
		view:           view,
	}
}

func (ic *InputController) SetInput(inputType InputType, msg string, onConfirm func(string)) {
	ic.inputType = inputType
	ic.onConfirm = onConfirm

	title := ""
	inputPrefix := ""
	if inputType == CONFIRM {
		title = " Confirmation "
		inputPrefix = "> " + msg + " (y/n) "
	} else {
		title = fmt.Sprintf(" Input [%s] ", msg)
		inputPrefix = "> "
	}

	ic.view.SetInput(title, inputPrefix)
}

func (ic *InputController) onType(content string, event view.KeyEvent) {
	if ic.inputType == CONFIRM {
		ic.mediator.notify(INPUT_DONE, content)
		ic.onConfirm(content)

		return
	}

	if event == view.CONFIRM || event == view.CANCEL {
		ic.mediator.notify(INPUT_DONE, content)

		ic.onConfirm(content)
	}
}
