package controller

import (
	"fmt"

	"github.com/dinhhuy258/fm/pkg/gui/view"
	"github.com/dinhhuy258/fm/pkg/key"
	"github.com/dinhhuy258/fm/pkg/optional"
)

type InputType int8

const (
	InputText InputType = iota
	InputConfirm
)

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

func (ic *InputController) GetInputBuffer() string {
	return ic.view.GetInputBuffer()
}

func (ic *InputController) UpdateInputBufferFromKey(key key.Key) {
	ic.view.UpdateInputBufferFromKey(key)
}

func (ic *InputController) UpdateView() {
	// ic.view.UpdateView(ic.title, ic.prompt, *ic.value.GetOrElse(&defaultInputValue))
}
