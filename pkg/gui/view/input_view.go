package view

import (
	"github.com/dinhhuy258/fm/pkg/key"
	"github.com/dinhhuy258/gocui"
)

type InputView struct {
	*View
	prompt string
}

func newInputView(g *gocui.Gui, v *gocui.View) *InputView {
	iv := &InputView{
		View: newView(g, v),
	}

	iv.prompt = "> "
	iv.v.Title = " Input "

	return iv
}

func (iv *InputView) SetInputBuffer(input string) {
	iv.v.TextArea.TypeString(iv.prompt + input)
	iv.v.RenderTextArea()

	_, _ = iv.g.SetCurrentView(iv.v.Name())
	iv.SetViewOnTop()
}

func (iv *InputView) GetInputBuffer() string {
	return iv.v.BufferLines()[0][len(iv.prompt):]
}

func (iv *InputView) UpdateInputBufferFromKey(key key.Key) {
	iv.g.Update(func(g *gocui.Gui) error {
		switch k := key.(type) {
		case rune:
			iv.v.TextArea.TypeRune(key.(rune))
		case gocui.Key:
			switch {
			case k == gocui.KeyBackspace || k == gocui.KeyBackspace2:
				iv.v.TextArea.DeleteChar()
			case k == gocui.KeyArrowLeft:
				x, _ := iv.v.Cursor()

				if x > len(iv.prompt) {
					iv.v.TextArea.MoveCursorLeft()
				}
			case k == gocui.KeyArrowRight:
				iv.v.TextArea.MoveCursorRight()
			}
		}

		iv.v.RenderTextArea()

		return nil
	})
}
