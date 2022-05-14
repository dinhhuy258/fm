package view

import (
	"fmt"

	"github.com/dinhhuy258/gocui"
)

type KeyEvent byte

const (
	TYPING  = 0
	CANCEL  = 1
	CONFIRM = 2
)

const (
	inputViewPrefix = "> "
)

type InputView struct {
	v         *View
	onType    func(string, KeyEvent)
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

func (iv *InputView) SetOnType(onType func(string, KeyEvent)) {
	iv.onType = onType
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

	if iv.onType != nil {
		keyEvent := KeyEvent(TYPING)
		if key == gocui.KeyEnter {
			keyEvent = KeyEvent(CONFIRM)
		} else if key == gocui.KeyEsc {
			keyEvent = KeyEvent(CANCEL)
		}

		viewContent := iv.v.v.BufferLines()[0]
		iv.onType(viewContent[len(inputViewPrefix):], keyEvent)
	}
}
