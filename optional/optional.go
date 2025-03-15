package optional

import "encoding/json"

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

func (o Optional[T]) MarshalJSON() ([]byte, error) {
	if !o.set {
		return []byte("null"), nil
	}

	return json.Marshal(o.value)
}

func (o *Optional[T]) UnmarshalJSON(bytes []byte) error {
	err := json.Unmarshal(bytes, &o.value)
	if err != nil {
		return err
	}

	o.set = true
	return nil
}
