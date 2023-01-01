package optional

var (
	NullInt = Null[int]()
)

type Value[T any] struct {
	Value  T
	Exists bool
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
