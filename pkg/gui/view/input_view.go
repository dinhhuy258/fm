package view

import (
	"github.com/dinhhuy258/fm/pkg/key"
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
	return iv.BufferLines()[0][len(iv.prompt):]
}

func (iv *InputView) UpdateInputBufferFromKey(key key.Key) {
	textArea := iv.TextArea

	switch k := key.(type) {
	case rune:
		textArea.TypeRune(key.(rune))
	case gocui.Key:
		switch {
		case key == gocui.KeySpace:
			textArea.TypeRune(' ')
		case k == gocui.KeyBackspace || k == gocui.KeyBackspace2:
			x, _ := iv.Cursor()

			if x > len(iv.prompt) {
				textArea.BackSpaceChar()
			}
		case k == gocui.KeyArrowLeft:
			x, _ := iv.Cursor()

			if x > len(iv.prompt) {
				textArea.MoveCursorLeft()
			}
		case k == gocui.KeyArrowRight:
			textArea.MoveCursorRight()
		case key == gocui.KeyCtrlA || key == gocui.KeyHome:
			textArea.SetCursor2D(len(iv.prompt), 0)
		case key == gocui.KeyCtrlE || key == gocui.KeyEnd:
			textArea.GoToEndOfLine()
		}
	}

	iv.RenderTextArea()
}
