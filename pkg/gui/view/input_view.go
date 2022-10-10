package view

import (
	"unicode"

	"github.com/dinhhuy258/fm/pkg/gui/key"
	"github.com/dinhhuy258/gocui"
)

const inputPrompt = "> "

type InputView struct {
	*View
	prompt string
}

func newInputView(v *gocui.View) *InputView {
	iv := &InputView{
		View:   newView(v),
		prompt: inputPrompt,
	}

	iv.Title = " Input "

	return iv
}

func (iv *InputView) SetInputBuffer(input string) {
	textArea := iv.TextArea

	textArea.Clear()
	textArea.TypeString(iv.prompt + input)
	iv.RenderTextArea()
}

func (iv *InputView) GetInputBuffer() string {
	inputContent := iv.TextArea.GetContent()
	if len(inputContent) <= len(iv.prompt) {
		return ""
	}

	return inputContent[len(iv.prompt):]
}

func (iv *InputView) UpdateInputBufferFromKey(k key.Key) {
	textArea := iv.TextArea

	key := k.Key
	ch := k.Ch
	mod := k.Mod

	switch {
	case key == gocui.KeySpace:
		textArea.TypeRune(' ')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		x, _ := iv.Cursor()

		if x > len(iv.prompt) {
			textArea.BackSpaceChar()
		}
	case key == gocui.KeyArrowLeft:
		x, _ := iv.Cursor()

		if x > len(iv.prompt) {
			textArea.MoveCursorLeft()
		}
	case key == gocui.KeyArrowRight:
		textArea.MoveCursorRight()
	case key == gocui.KeyCtrlA || key == gocui.KeyHome:
		textArea.SetCursor2D(len(iv.prompt), 0)
	case key == gocui.KeyCtrlE || key == gocui.KeyEnd:
		textArea.GoToEndOfLine()
		// TODO: see if we need all three of these conditions: maybe the final one is sufficient
	case ch != 0 && mod == gocui.ModNone && unicode.IsPrint(ch):
		textArea.TypeRune(ch)
	}

	iv.RenderTextArea()
}
