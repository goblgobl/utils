package validation

type NoopValidator[T any] struct {
	invalidValue *Invalid
	fn           BoolFuncValidator[T]
	dflt         any
	required     bool
	nullable     bool
}

func Noop[T any]() *NoopValidator[T] {
	return new(NoopValidator[T])
}

func (v *NoopValidator[T]) Validate(value any, ctx *Context[T]) any {
	return value
}
