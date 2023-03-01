package validation

import (
	"src.goblgobl.com/utils/typed"
)

type ObjectFuncValidator[T any] func(value map[string]any, ctx *Context[T]) any

type ObjectField[T any] struct {
	field     *Field
	validator Validator[T]
}

type ObjectValidator[T any] struct {
	dflt     any
	fn       ObjectFuncValidator[T]
	fields   []ObjectField[T]
	required bool
}

func Object[T any]() *ObjectValidator[T] {
	return &ObjectValidator[T]{
		fields: make([]ObjectField[T], 0, 5),
	}
}

// While Validate(raw, ctx) can be called on any validator, including this
// ObjectValidator, the ValidateInput is a speciality function meant to be called
// at the top level. (ValidateInput is really just a micro-optimization that
// skips some checks that we know we don't need at the top level, like a type
// check on the input)
func (v *ObjectValidator[T]) ValidateInput(input typed.Typed, ctx *Context[T]) bool {
	ctx.Input = input
	v.validate(map[string]any(input), ctx)
	return ctx.IsValid()
}

func (v *ObjectValidator[T]) Validate(raw any, ctx *Context[T]) any {
	if raw == nil {
		if dflt := v.dflt; dflt != nil {
			return dflt
		}
		if v.required {
			ctx.InvalidField(Required)
		}
		return nil
	}

	object, ok := raw.(map[string]any)
	if !ok {
		ctx.InvalidField(TypeObject)
		return nil
	}

	return v.validate(object, ctx)
}

func (v *ObjectValidator[T]) validate(object map[string]any, ctx *Context[T]) any {
	ctx.StartObject(object)
	defer ctx.EndObject()
	for _, vf := range v.fields {
		field := vf.field
		ctx.Field = field
		fieldName := field.Name
		object[fieldName] = vf.validator.Validate(object[fieldName], ctx)
	}

	if fn := v.fn; fn != nil {
		return fn(object, ctx)
	}
	return object
}

// ObjectValidator is the only validator that holds on to a Field (i.e. it's the
// only validator that knows what things are called). You'd think the field is just
// something like "name" or "password", but a field also contains the whole path
// for nested objects, like "user.name" and "user.password". So, when an ObjectValidator
// is nested in another ObjectValidator, we need to build up the field paths throughout
// the graph.
// While an ArrayValidator doesn't hold a Field directly, its validator can be
// an ObjectValidator (which does hold Fields), so we also have to inform the
// ArrayValidator about any nesting so that it can inform its validator (if that
// validator is an ObjectValidator)
func (v *ObjectValidator[T]) Field(name string, validator Validator[T]) *ObjectValidator[T] {
	if _, ok := validator.(*ArrayValidator[T]); ok {
		name += ".#"
	}

	field := BuildField(name)
	validator = nestValidator(field, validator)

	v.fields = append(v.fields, ObjectField[T]{
		field:     field,
		validator: validator,
	})
	return v
}

func (v *ObjectValidator[T]) ForceField(field *Field) *ObjectValidator[T] {
	return nestValidator[T](field, v).(*ObjectValidator[T])
}

func (v *ObjectValidator[T]) Required() *ObjectValidator[T] {
	v.required = true
	return v
}

func (v *ObjectValidator[T]) Default(dflt any) *ObjectValidator[T] {
	v.dflt = dflt
	return v
}

func (v *ObjectValidator[T]) Func(fn ObjectFuncValidator[T]) *ObjectValidator[T] {
	v.fn = fn
	return v
}

// v is a validator that's being nested inside of an ObjectValidator or an
// ArrayValidator. This means that the field names inside of V need to change.
// For example, if v has a field "name" and is now being nested under a "user"
// field, that field is now "user.name".
// We don't mutate the ObjectValidator, but create a clone
func (v *ObjectValidator[T]) nest(field *Field) *ObjectValidator[T] {
	fields := make([]ObjectField[T], len(v.fields))
	for i, vf := range v.fields {
		nested := vf.field.nest(field)
		validator := nestValidator(nested, vf.validator)
		fields[i] = ObjectField[T]{
			validator: validator,
			field:     nested,
		}
	}

	return &ObjectValidator[T]{
		fn:       v.fn,
		dflt:     v.dflt,
		fields:   fields,
		required: v.required,
	}
}

func nestValidator[T any](field *Field, validator Validator[T]) Validator[T] {
	switch inner := validator.(type) {
	case *ObjectValidator[T]:
		return inner.nest(field)
	case *ArrayValidator[T]:
		return inner.nest(field)
	}
	return validator
}
