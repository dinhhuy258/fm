package ctx

import (
	"github.com/dinhhuy258/fm/pkg/state"
)

type Context interface {
	State() *state.State
	PopMode() error
	PushMode(mode string) error
}
