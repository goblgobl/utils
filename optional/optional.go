package optional

import (
	"bytes"
	"encoding/json"
)

var (
	NullInt    = Null[int]()
	NullString = Null[string]()
	NullBool   = Null[bool]()
	NullFloat  = Null[float64]()

	jsonNull = []byte("null")
)

type Value[T any] struct {
	Value  T
	Exists bool
}

type Int = Value[int]
type Bool = Value[bool]
type Float = Value[float64]
type String = Value[string]

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

func NewBool(value bool) Value[bool] {
	return New[bool](value)
}

func NewInt(value int) Value[int] {
	return New[int](value)
}

func NewFloat(value float64) Value[float64] {
	return New[float64](value)
}

func NewString(value string) Value[string] {
	return New[string](value)
}
