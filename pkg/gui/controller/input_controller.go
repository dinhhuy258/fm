package controller

import (
	"fmt"

	"github.com/dinhhuy258/fm/pkg/gui/view"
)

type InputType int8

const (
	InputText InputType = iota
	InputConfirm
)

type InputController struct {
	*BaseController

	inputType InputType
	onConfirm func(string)

	view *view.InputView
}

func newInputController(baseController *BaseController, view *view.InputView) *InputController {
	inputController := &InputController{
		BaseController: baseController,
		view:           view,
	}

	view.SetOnType(inputController.onType)

	return inputController
}

func (ic *InputController) SetInput(inputType InputType, msg string, onConfirm func(string)) {
	ic.inputType = inputType
	ic.onConfirm = onConfirm

	title := ""
	prompt := ""

	if inputType == InputConfirm {
		title = " Confirmation "
		prompt = "> " + msg + " (y/n) "
	} else {
		title = fmt.Sprintf(" Input [%s] ", msg)
		prompt = "> "
	}

	ic.view.SetInput(title, prompt)
}

func (ic *InputController) onType(content string, event view.InputEvent) {
	if ic.inputType == InputConfirm {
		ic.mediator.notify(InputDone, content)
		ic.onConfirm(content)

		return
	}

	if event == view.Confirm || event == view.Cancel {
		ic.mediator.notify(InputDone, content)

		ic.onConfirm(content)
	}
}
