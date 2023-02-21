package validation

import "src.goblgobl.com/utils/uuid"

type UUIDValidator[T any] struct {
	dflt     any
	required bool
}

func UUID[T any]() *UUIDValidator[T] {
	return new(UUIDValidator[T])
}

func (v *UUIDValidator[T]) Validate(raw any, ctx *Context[T]) any {
	if raw == nil {
		if dflt := v.dflt; dflt != nil {
			return dflt
		}
		if v.required {
			ctx.InvalidField(Required)
		}
		return nil
	}

	value, ok := raw.(string)
	if !ok {
		ctx.InvalidField(TypeUUID)
		return nil
	}

	if !uuid.IsValid(value) {
		ctx.InvalidField(TypeUUID)
	}

	return value
}

func (v *UUIDValidator[T]) Required() *UUIDValidator[T] {
	v.required = true
	return v
}

// Meant to be used in conjunction with Clone(). Maybe on create, the field is
// required, but on update, it isn't.
func (v *UUIDValidator[T]) NotRequired() *UUIDValidator[T] {
	v.required = false
	return v
}

func (v *UUIDValidator[T]) Default(dflt any) *UUIDValidator[T] {
	v.dflt = dflt
	return v
}

func (v *UUIDValidator[T]) Clone() *UUIDValidator[T] {
	return &UUIDValidator[T]{
		dflt:     v.dflt,
		required: v.required,
	}
}
