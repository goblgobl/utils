package optional

import (
	"bytes"
	"encoding/json"
)

var (
	NullInt    = Null[int]()
	NullString = Null[string]()
	jsonNull   = []byte("null")
)

type Value[T any] struct {
	Value  T
	Exists bool
}

func (v Value[T]) MarshalJSON() ([]byte, error) {
	if !v.Exists {
		return jsonNull, nil
	}
	return json.Marshal(v.Value)
}

func (v *Value[T]) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, jsonNull) {
		return nil
	}

	var value T
	err := json.Unmarshal(data, &value)
	v.Value = value
	v.Exists = true
	return err
}

func New[T any](value T) Value[T] {
	return Value[T]{Value: value, Exists: true}
}

func Null[T any]() Value[T] {
	return Value[T]{}
}

func Int(value int) Value[int] {
	return New[int](value)
}

func String(value string) Value[string] {
	return New[string](value)
}
