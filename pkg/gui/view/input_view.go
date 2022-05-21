package view

import (
	"github.com/dinhhuy258/gocui"
)

type InputEvent int8

const (
	Typing InputEvent = iota
	Cancel
	Confirm
)

type InputView struct {
	v      *View
	prompt string
	input  string
	onType func(string, InputEvent)
}

func newInputView(g *gocui.Gui, v *gocui.View) *InputView {
	iv := &InputView{
		v: newView(g, v),
	}

	iv.prompt = "> "
	iv.v.v.Title = " Input "
	// iv.v.v.Editable = false
	// iv.v.v.Editor = gocui.EditorFunc(iv.inputEditor)

	return iv
}

func (iv *InputView) UpdateView(title string, prompt string, value string) {
	iv.prompt = prompt
	iv.input = ""

	iv.v.SetViewContent([]string{prompt + value})
	iv.v.SetTitle(title)
	_ = iv.v.v.SetCursor(len(prompt)+len(value), 0)
	_, _ = iv.v.g.SetCurrentView(iv.v.v.Name())

	iv.v.SetViewOnTop()
}

func (iv *InputView) SetOnType(onType func(string, InputEvent)) {
	iv.onType = onType
}

func (iv *InputView) SetInputBuffer(input string) {
	iv.input = input
	iv.v.SetViewContent([]string{iv.prompt + iv.input})
	_ = iv.v.v.SetCursor(len(iv.prompt)+len(iv.input), 0)

	iv.v.SetViewOnTop()
}

func (iv *InputView) InputEditor(key string) {
	iv.input += key
	iv.v.SetViewContent([]string{iv.prompt + iv.input})
	_ = iv.v.v.SetCursor(len(iv.prompt)+len(iv.input), 0)
	// switch {
	// case ch != 0 && mod == 0:
	// 	iv.v.v.EditWrite(ch)
	// case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
	// 	x, _ := iv.v.v.Cursor()
	// 	if x > len(iv.prompt) {
	// 		iv.v.v.EditDelete(true)
	// 	}
	// case key == gocui.KeyArrowLeft:
	// 	x, _ := iv.v.v.Cursor()
	//
	// 	if x > len(iv.prompt) {
	// 		iv.v.v.MoveCursor(-1, 0, false)
	// 	}
	// case key == gocui.KeyArrowRight:
	// 	iv.v.v.MoveCursor(1, 0, false)
	// }
	//
	// if iv.onType != nil {
	// 	keyEvent := Typing
	// 	inputValue := iv.v.v.BufferLines()[0][len(iv.prompt):]
	//
	// 	if key == gocui.KeyEnter {
	// 		keyEvent = Confirm
	// 	} else if key == gocui.KeyEsc {
	// 		keyEvent = Cancel
	// 		inputValue = ""
	// 	}
	//
	// 	iv.onType(inputValue, keyEvent)
	// }
}
