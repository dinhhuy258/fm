package msg

import (
	"errors"
	"strconv"

	"github.com/dinhhuy258/fm/pkg/gui/key"
)

type MessageContext map[string]interface{}

var errMessageNotFound = errors.New("Message name is not found")

// MessageFactory is a factory for creating messages.
type MessageFactory struct {
	messageFunc func(app IApp, key *key.Key, ctx MessageContext)
}

// newMessageFactory creates a new message factory.
func newMessageFactory(messageFunc func(app IApp, key *key.Key, ctx MessageContext)) *MessageFactory {
	return &MessageFactory{
		messageFunc: messageFunc,
	}
}

// New creates a new message.
func (mf *MessageFactory) New(args ...string) *Message {
	ctx := make(MessageContext)
	for idx, arg := range args {
		ctx["arg"+strconv.Itoa(idx+1)] = arg
	}

	return &Message{
		Func: mf.messageFunc,
		Ctx:  ctx,
	}
}

// messageFactories is a map of message factories.
var messageFactories = map[string]*MessageFactory{
	"BashExecSilently": newMessageFactory(BashExecSilently),
	"BashExec":         newMessageFactory(BashExec),
	"ToggleHidden":     newMessageFactory(ToggleHidden),
	"Refresh":          newMessageFactory(Refresh),
	"Quit":             newMessageFactory(Quit),
	"Null":             newMessageFactory(Null),

	"ClearLog":   newMessageFactory(ClearLog),
	"LogSuccess": newMessageFactory(LogSuccess),
	"LogWarning": newMessageFactory(LogWarning),
	"LogError":   newMessageFactory(LogError),

	"SwitchMode": newMessageFactory(SwitchMode),

	"FocusByIndex":    newMessageFactory(FocusByIndex),
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

	"ToggleSelection":       newMessageFactory(ToggleSelection),
	"ClearSelection":        newMessageFactory(ClearSelection),
	"ToggleSelectionByPath": newMessageFactory(ToggleSelectionByPath),
	"SelectAll":             newMessageFactory(SelectAll),

	"SortByDirFirst":     newMessageFactory(SortByDirFirst),
	"SortByDateModified": newMessageFactory(SortByDateModified),
	"SortByName":         newMessageFactory(SortByName),
	"SortBySize":         newMessageFactory(SortBySize),
	"SortByExtension":    newMessageFactory(SortByExtension),
	"ReverseSort":        newMessageFactory(ReverseSort),
}

// NewMessage creates a new message.
func NewMessage(name string, args ...string) (*Message, error) {
	messageFactory, hasMessage := messageFactories[name]
	if !hasMessage {
		return nil, errMessageNotFound
	}

	return messageFactory.New(args...), nil
}
