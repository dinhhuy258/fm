package controller

import (
	"github.com/dinhhuy258/fm/pkg/gui/view"
	"github.com/dinhhuy258/fm/pkg/key"
)

type InputController struct {
	*BaseController

	view *view.InputView
}

func newInputController(baseController *BaseController, view *view.InputView) *InputController {
	inputController := &InputController{
		BaseController: baseController,
		view:           view,
	}

	return inputController
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
