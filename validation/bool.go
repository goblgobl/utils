package validation

type BoolFuncValidator[T any] func(value bool, ctx *Context[T]) any

type BoolValidator[T any] struct {
	dflt     any
	required bool
	nullable bool
	fn       BoolFuncValidator[T]

	invalidValue *Invalid
}

func Bool[T any]() *BoolValidator[T] {
	return new(BoolValidator[T])
}

func (v *BoolValidator[T]) Validate(raw any, ctx *Context[T]) any {
	if raw == nil {
		if dflt := v.dflt; dflt != nil {
			return dflt
		}
		if v.nullable {
			return nil
		}
		if v.required {
			ctx.InvalidField(Required)
		}
		return nil
	}

	value, ok := raw.(bool)
	if !ok {
		ctx.InvalidField(TypeBool)
		return nil
	}

	if fn := v.fn; fn != nil {
		return fn(value, ctx)
	}

	return value
}

func (v *BoolValidator[T]) Required() *BoolValidator[T] {
	v.required = true
	return v
}

// Meant to be used in conjunction with Clone(). Maybe on create, the field is
// required, but on update, it isn't.
func (v *BoolValidator[T]) NotRequired() *BoolValidator[T] {
	v.required = false
	return v
}

func (v *BoolValidator[T]) Nullable() *BoolValidator[T] {
	v.nullable = true
	return v
}

func (v *BoolValidator[T]) Default(dflt any) *BoolValidator[T] {
	v.dflt = dflt
	return v
}

func (v *BoolValidator[T]) Func(fn BoolFuncValidator[T]) *BoolValidator[T] {
	v.fn = fn
	return v
}

func (v *BoolValidator[T]) Clone() *BoolValidator[T] {
	return &BoolValidator[T]{
		fn:           v.fn,
		dflt:         v.dflt,
		required:     v.required,
		nullable:     v.nullable,
		invalidValue: v.invalidValue,
	}
}
