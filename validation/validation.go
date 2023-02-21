package validation

type Validator[T any] interface {
	Validate(value any, ctx *Context[T]) any
}
