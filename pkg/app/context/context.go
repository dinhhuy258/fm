package context

type Context interface {
	State() *State
	PopMode() error
	PushMode(mode string) error
}
