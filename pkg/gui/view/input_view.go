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
	iv.v.SetViewContent([]string{iv.prompt + input})
	_ = iv.v.v.SetCursor(len(iv.prompt)+len(input), 0)

	_, _ = iv.v.g.SetCurrentView(iv.v.v.Name())
	iv.v.SetViewOnTop()
}

func (iv *InputView) GetInputBuffer() string {
	return iv.v.v.BufferLines()[0][len(iv.prompt):]
}

func (iv *InputView) UpdateInputBufferFromKey(key key.Key) {
	switch k := key.(type) {
	case rune:
		// iv.v.v.EditWrite(k)
	case gocui.Key:
		switch {
		case k == gocui.KeyBackspace || k == gocui.KeyBackspace2:
			x, _ := iv.v.v.Cursor()
			if x > len(iv.prompt) {
				// iv.v.v.EditDelete(true)
			}
		case k == gocui.KeyArrowLeft:
			x, _ := iv.v.v.Cursor()

			if x > len(iv.prompt) {
				// iv.v.v.MoveCursor(-1, 0, false)
			}
		case k == gocui.KeyArrowRight:
			// iv.v.v.MoveCursor(1, 0, false)
		}
	}
}
