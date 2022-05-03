package view

import (
	"fmt"

	"github.com/dinhhuy258/gocui"
)

const (
	inputViewPrefix = "> "
)

type InputView struct {
	v         *View
	inputChan chan string
}

func newInputView(g *gocui.Gui, v *gocui.View) *InputView {
	iv := &InputView{
		v:         newView(g, v),
		inputChan: make(chan string, 1),
	}

	iv.v.v.Title = " Input "
	iv.v.v.Editable = true
	iv.v.v.Editor = gocui.EditorFunc(iv.inputEditor)

	return iv
}

func (iv *InputView) inputEditor(_ *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case ch != 0 && mod == 0:
		iv.v.v.EditWrite(ch)
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		x, _ := iv.v.v.Cursor()
		if x > len(inputViewPrefix) {
			iv.v.v.EditDelete(true)
		}
	case key == gocui.KeyArrowLeft:
		x, _ := iv.v.v.Cursor()

		if x > len(inputViewPrefix) {
			iv.v.v.MoveCursor(-1, 0, false)
		}
	case key == gocui.KeyArrowRight:
		iv.v.v.MoveCursor(1, 0, false)
	case key == gocui.KeyEnter:
		viewContent := iv.v.v.BufferLines()[0]
		iv.inputChan <- viewContent[len(inputViewPrefix):]
	case key == gocui.KeyEsc:
		iv.inputChan <- ""
	}
}

func (iv *InputView) SetInput(ask string, onInput func(string)) {
	iv.v.SetViewContent([]string{inputViewPrefix})
	iv.v.SetTitle(fmt.Sprintf(" Input [%s] ", ask))
	_ = iv.v.v.SetCursor(len(inputViewPrefix), 0)
	_, _ = iv.v.g.SetCurrentView(iv.v.v.Name())

	iv.v.SetViewOnTop()

	go func() {
		ans := <-iv.inputChan

		onInput(ans)
	}()
}
