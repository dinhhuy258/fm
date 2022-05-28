package msg

import "errors"

var errMessageNotFound = errors.New("Message name is not found")

type MessageFactory struct {
	messageFunc func(app IApp, params ...string)
}

func newMessageFactory(messageFunc func(app IApp, params ...string)) *MessageFactory {
	return &MessageFactory{
		messageFunc: messageFunc,
	}
}

func (mf *MessageFactory) New(args ...string) *Message {
	return &Message{
		Func: mf.messageFunc,
		Args: args,
	}
}

var messageFactories = map[string]*MessageFactory{
	"BashExecSilently": newMessageFactory(BashExecSilently),
	"BashExec":         newMessageFactory(BashExec),
	"ToggleHidden":     newMessageFactory(ToggleHidden),
	"Refresh":          newMessageFactory(Refresh),
	"Quit":             newMessageFactory(Quit),

	"LogSuccess": newMessageFactory(LogSuccess),
	"LogError":   newMessageFactory(LogError),

	"SwitchMode": newMessageFactory(SwitchMode),
	"PopMode":    newMessageFactory(PopMode),

	"FocusNext":       newMessageFactory(FocusNext),
	"FocusPrevious":   newMessageFactory(FocusPrevious),
	"FocusPath":       newMessageFactory(FocusPath),
	"FocusFirst":      newMessageFactory(FocusFirst),
	"FocusLast":       newMessageFactory(FocusLast),
	"Enter":           newMessageFactory(Enter),
	"Back":            newMessageFactory(Back),
	"ChangeDirectory": newMessageFactory(ChangeDirectory),

	"UpdateInputBufferFromKey": newMessageFactory(UpdateInputBufferFromKey),
	"SetInputBuffer":           newMessageFactory(SetInputBuffer),

	"PasteSelections":  newMessageFactory(PasteSelections),
	"NewFileFromInput": newMessageFactory(NewFileFromInput),
	"DeleteCurrent":    newMessageFactory(DeleteCurrent),
	"DeleteSelections": newMessageFactory(DeleteSelections),
	"SearchFromInput":  newMessageFactory(SearchFromInput),

	"ToggleSelection": newMessageFactory(ToggleSelection),
	"ClearSelection":  newMessageFactory(ClearSelection),

	"MarkSave": newMessageFactory(MarkSave),
	"MarkLoad": newMessageFactory(MarkLoad),
}

func NewMessage(name string, args ...string) (*Message, error) {
	messageFactory, hasMessageFactory := messageFactories[name]
	if !hasMessageFactory {
		return nil, errMessageNotFound
	}

	return messageFactory.New(args...), nil
}
