package types

// Optional represents an optional value.
type Optional[T any] struct {
	value *T
}

// New creates an Optional with the given value.
func New[T any](value T) Optional[T] {
	return Optional[T]{&value}
}

// NewEmpty creates an Optional with no value.
func NewEmpty[T any]() Optional[T] {
	return Optional[T]{}
}

// IsPresent checks if the Optional has a value.
func (o *Optional[T]) IsPresent() bool {
	return o.value != nil
}

// Get returns the value of the Optional.
func (o *Optional[T]) Get() *T {
	return o.value
}

// GetOrElse returns the value of the Optional or the given default value.
func (o *Optional[T]) GetOrElse(defaultValue *T) *T {
	if o.value != nil {
		return o.value
	}

	return defaultValue
}

// IfPresent executes the given function if the Optional has a value.
func (o *Optional[T]) IfPresent(consumer func(*T)) {
	if o.value != nil {
		consumer(o.value)
	}
}

// IfPresentOrElse executes the `consumer` function if the Optional has a value
// otherwise `consumerOrElse` is executed.
func (o *Optional[T]) IfPresentOrElse(consumer func(*T), consumerOrElse func()) {
	if o.value != nil {
		consumer(o.value)
	} else {
		consumerOrElse()
	}
}