package validation

type AnyFuncValidator[T any] func(value any, ctx *Context[T]) any

type AnyValidator[T any] struct {
	fn       AnyFuncValidator[T]
	dflt     any
	required bool
}

func Any[T any]() *AnyValidator[T] {
	return new(AnyValidator[T])
}

func (v *AnyValidator[T]) Validate(raw any, ctx *Context[T]) any {
	if raw == nil {
		if dflt := v.dflt; dflt != nil {
			return dflt
		}
		if v.required {
			ctx.InvalidField(Required)
		}
		return nil
	}

	if fn := v.fn; fn != nil {
		return fn(raw, ctx)
	}
	return raw
}

func (v *AnyValidator[T]) Required() *AnyValidator[T] {
	v.required = true
	return v
}

// Meant to be used in conjunction with Clone(). Maybe on create, the field is
// required, but on update, it isn't.
func (v *AnyValidator[T]) NotRequired() *AnyValidator[T] {
	v.required = false
	return v
}

func (v *AnyValidator[T]) Default(dflt any) *AnyValidator[T] {
	v.dflt = dflt
	return v
}

func (v *AnyValidator[T]) Func(fn AnyFuncValidator[T]) *AnyValidator[T] {
	v.fn = fn
	return v
}

func (v *AnyValidator[T]) Clone() *AnyValidator[T] {
	return &AnyValidator[T]{
		fn:       v.fn,
		dflt:     v.dflt,
		required: v.required,
	}
}
