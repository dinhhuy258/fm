package message

import (
	"github.com/dinhhuy258/fm/pkg/app/context"
)

type Message struct {
	Func func(context context.Context, params ...interface{}) error
	Args []interface{}
}
