package optional

type Optional[T any] struct {
	value T
	set   bool
}

func From[T any](value T) Optional[T] {
	return Optional[T]{
		value: value,
		set:   true,
	}
}

func FromPtr[T any](value *T) Optional[T] {
	var optional Optional[T]

	if value != nil {
		optional.value = *value
		optional.set = true
	}

	return optional
}

func (o Optional[T]) Value() (T, bool) {
	if !o.set {
		var zero T
		return zero, false
	}

	return o.value, true
}

func (o Optional[T]) IsSet() bool {
	return o.set
}
