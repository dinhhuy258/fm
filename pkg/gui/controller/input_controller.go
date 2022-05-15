package controller

import (
	"fmt"

	"github.com/dinhhuy258/fm/pkg/gui/view"
)

type InputType int8

const (
	Input InputType = iota
	Confirm
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
	inputPrefix := ""

	if inputType == Confirm {
		title = " Confirmation "
		inputPrefix = "> " + msg + " (y/n) "
	} else {
		title = fmt.Sprintf(" Input [%s] ", msg)
		inputPrefix = "> "
	}

	ic.view.SetInput(title, inputPrefix)
}

func (ic *InputController) onType(content string, event view.InputEvent) {
	if ic.inputType == Confirm {
		ic.mediator.notify(InputDone, content)
		ic.onConfirm(content)

		return
	}

	if event == view.Confirm || event == view.Cancel {
		ic.mediator.notify(InputDone, content)

		ic.onConfirm(content)
	}
}
