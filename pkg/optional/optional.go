package optional

type Optional[T any] struct {
	value *T
}

func NewOptional[T any](value T) Optional[T] {
	return Optional[T]{&value}
}

func NewEmptyOptional[T any]() Optional[T] {
	return Optional[T]{}
}

func (o *Optional[T]) IsPresent() bool {
	return o.value != nil
}

func (o *Optional[T]) Get() *T {
	return o.value
}

func (o *Optional[T]) GetOrElse(defaultValue *T) *T {
	if o.value != nil {
		return o.value
	}

	return defaultValue
}

func (o *Optional[T]) IfPresent(consumer func(*T)) {
	if o.value != nil {
		consumer(o.value)
	}
}

func (o *Optional[T]) IfPresentOrElse(consumer func(*T), consumerOrElse func()) {
	if o.value != nil {
		consumer(o.value)
	} else {
		consumerOrElse()
	}
}
