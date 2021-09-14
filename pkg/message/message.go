package message

import "github.com/dinhhuy258/fm/pkg/ctx"

type Message struct {
	Func func(context *ctx.Context, params ...interface{}) error
	Args []interface{}
}
