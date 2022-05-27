package view

import (
	"github.com/dinhhuy258/fm/pkg/key"
	"github.com/dinhhuy258/gocui"
)

type InputView struct {
	v      *View
	prompt string
}

func newInputView(g *gocui.Gui, v *gocui.View) *InputView {
	iv := &InputView{
		v: newView(g, v),
	}

	iv.prompt = "> "
	iv.v.v.Title = " Input "

	return iv
}

func (iv *InputView) SetInputBuffer(input string) {
	iv.v.v.TextArea.TypeString(iv.prompt + input)
	iv.v.v.RenderTextArea()

	_, _ = iv.v.g.SetCurrentView(iv.v.v.Name())
	iv.v.SetViewOnTop()
}

func (iv *InputView) GetInputBuffer() string {
	return iv.v.v.BufferLines()[0][len(iv.prompt):]
}

func (iv *InputView) UpdateInputBufferFromKey(key key.Key) {
	iv.v.g.Update(func(g *gocui.Gui) error {
		switch k := key.(type) {
		case rune:
			iv.v.v.TextArea.TypeRune(key.(rune))
		case gocui.Key:
			switch {
			case k == gocui.KeyBackspace || k == gocui.KeyBackspace2:
				iv.v.v.TextArea.DeleteChar()
			case k == gocui.KeyArrowLeft:
				x, _ := iv.v.v.Cursor()

				if x > len(iv.prompt) {
					iv.v.v.TextArea.MoveCursorLeft()
				}
			case k == gocui.KeyArrowRight:
				iv.v.v.TextArea.MoveCursorRight()
			}
		}

		iv.v.v.RenderTextArea()

		return nil
	})
}
