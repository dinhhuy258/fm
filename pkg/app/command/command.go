package command

import (
	"github.com/dinhhuy258/fm/pkg/app/context"
)

type Command struct {
	Help string
	Func func(context context.Context, params ...interface{}) error
	Args []interface{}
}
