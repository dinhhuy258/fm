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
	onType func(string, InputEvent)
}

func newInputView(g *gocui.Gui, v *gocui.View) *InputView {
	iv := &InputView{
		v: newView(g, v),
	}

	iv.v.v.Title = " Input "
	iv.v.v.Editable = true
	iv.v.v.Editor = gocui.EditorFunc(iv.inputEditor)

	return iv
}

func (iv *InputView) UpdateView(title string, prompt string, value string) {
	iv.prompt = prompt

	iv.v.SetViewContent([]string{prompt + value})
	iv.v.SetTitle(title)
	_ = iv.v.v.SetCursor(len(prompt) + len(value), 0)
	_, _ = iv.v.g.SetCurrentView(iv.v.v.Name())

	iv.v.SetViewOnTop()
}

func (iv *InputView) SetOnType(onType func(string, InputEvent)) {
	iv.onType = onType
}

func (iv *InputView) inputEditor(_ *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case ch != 0 && mod == 0:
		iv.v.v.EditWrite(ch)
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		x, _ := iv.v.v.Cursor()
		if x > len(iv.prompt) {
			iv.v.v.EditDelete(true)
		}
	case key == gocui.KeyArrowLeft:
		x, _ := iv.v.v.Cursor()

		if x > len(iv.prompt) {
			iv.v.v.MoveCursor(-1, 0, false)
		}
	case key == gocui.KeyArrowRight:
		iv.v.v.MoveCursor(1, 0, false)
	}

	if iv.onType != nil {
		keyEvent := Typing
		if key == gocui.KeyEnter {
			keyEvent = Confirm
		} else if key == gocui.KeyEsc {
			keyEvent = Cancel
		}

		viewContent := iv.v.v.BufferLines()[0]
		iv.onType(viewContent[len(iv.prompt):], keyEvent)
	}
}
