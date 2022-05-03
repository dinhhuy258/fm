package command

import (
	"github.com/dinhhuy258/fm/pkg/app/context"
)

type IApp interface {
	State() *context.State
	PopMode() error
	PushMode(mode string) error
}

type Command struct {
	Help               string
	Func               func(app IApp, params ...interface{}) error
	Args               []interface{}
}
