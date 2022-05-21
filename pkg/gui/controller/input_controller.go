package controller

import (
	"fmt"

	"github.com/dinhhuy258/fm/pkg/gui/view"
	"github.com/dinhhuy258/fm/pkg/optional"
)

type InputType int8

const (
	InputText InputType = iota
	InputConfirm
)

var defaultInputValue = ""

type InputController struct {
	*BaseController

	title     string
	prompt    string
	value     optional.Optional[string]
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

func (ic *InputController) SetInput(inputType InputType, msg string,
	value optional.Optional[string], onConfirm func(string),
) {
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

	ic.title = title
	ic.prompt = prompt
	ic.value = value
}

func (ic *InputController) SetInputBuffer(input string) {
	ic.view.SetInputBuffer(input)
}

func (ic *InputController) UpdateInputBufferFromKey(key string) {
	ic.view.InputEditor(key)
}

func (ic *InputController) UpdateView() {
	ic.view.UpdateView(ic.title, ic.prompt, *ic.value.GetOrElse(&defaultInputValue))
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
