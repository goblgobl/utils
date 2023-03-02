package validation

import (
	"fmt"

	"src.goblgobl.com/utils"
	"src.goblgobl.com/utils/optional"
)

type ArrayFuncValidator[T any] func(value []any, ctx *Context[T]) any

type ArrayValidator[T any] struct {
	invalidLength *Invalid
	validator     Validator[T]
	fn            ArrayFuncValidator[T]
	dflt          any
	minLength     optional.Int
	maxLength     optional.Int
	required      bool
	convertToType bool
}

func Array[T any]() *ArrayValidator[T] {
	return &ArrayValidator[T]{}
}

func (v *ArrayValidator[T]) Validate(raw any, ctx *Context[T]) any {
	if raw == nil {
		if dflt := v.dflt; dflt != nil {
			return dflt
		}
		if v.required {
			ctx.InvalidField(Required)
		}
		return nil
	}

	values, ok := raw.([]any)
	if !ok {
		ctx.InvalidField(TypeArray)
		return nil
	}

	if min := v.minLength; min.Exists && len(values) < min.Value {
		ctx.InvalidField(v.invalidLength)
		return values
	}

	if max := v.maxLength; max.Exists && len(values) > max.Value {
		ctx.InvalidField(v.invalidLength)
		return values
	}

	errorCount := ctx.ErrorCount()
	validator := v.validator
	ctx.StartArray()
	for i, value := range values {
		ctx.ArrayIndex(i)
		values[i] = validator.Validate(value, ctx)
	}
	ctx.EndArray()

	if fn := v.fn; fn != nil {
		return fn(values, ctx)
	}

	// Let's only do this if the above validator didn't add any errors
	if v.convertToType && ctx.ErrorCount() == errorCount {
		switch validator.(type) {
		case *StringValidator[T]:
			t := make([]string, len(values))
			for i, value := range values {
				t[i] = value.(string)
			}
			return t
		case *IntValidator[T]:
			t := make([]int, len(values))
			for i, value := range values {
				t[i] = value.(int)
			}
			return t
		case *FloatValidator[T]:
			t := make([]float64, len(values))
			for i, value := range values {
				t[i] = value.(float64)
			}
			return t
		case *BoolValidator[T]:
			t := make([]bool, len(values))
			for i, value := range values {
				t[i] = value.(bool)
			}
			return t
		case *UUIDValidator[T]:
			t := make([]string, len(values))
			for i, value := range values {
				t[i] = value.(string)
			}
			return t
		}
	}

	return values
}

func (v *ArrayValidator[T]) Required() *ArrayValidator[T] {
	v.required = true
	return v
}

func (v *ArrayValidator[T]) Default(dflt any) *ArrayValidator[T] {
	v.dflt = dflt
	return v
}

func (v *ArrayValidator[T]) ConvertToType() *ArrayValidator[T] {
	v.convertToType = true
	return v
}

func (v *ArrayValidator[T]) Min(min int) *ArrayValidator[T] {
	v.minLength = optional.NewInt(min)
	v.invalidLength = InvalidArrayLen(v.minLength, v.maxLength)
	return v
}

func (v *ArrayValidator[T]) Max(max int) *ArrayValidator[T] {
	v.maxLength = optional.NewInt(max)
	v.invalidLength = InvalidArrayLen(v.minLength, v.maxLength)
	return v
}

func (v *ArrayValidator[T]) Range(min int, max int) *ArrayValidator[T] {
	v.minLength = optional.NewInt(min)
	v.maxLength = optional.NewInt(max)
	v.invalidLength = InvalidArrayLen(v.minLength, v.maxLength)
	return v
}

func (v *ArrayValidator[T]) Func(fn ArrayFuncValidator[T]) *ArrayValidator[T] {
	v.fn = fn
	return v
}

func (v *ArrayValidator[T]) Validator(validator Validator[T]) *ArrayValidator[T] {
	if ov, ok := validator.(*ObjectValidator[T]); ok {
		validator = ov.nest(BuildField("#"))
	}
	v.validator = validator
	return v
}

// See the ObjectValidator[T] Field method
func (v *ArrayValidator[T]) nest(field *Field) *ArrayValidator[T] {
	inner, ok := v.validator.(*ObjectValidator[T])
	if !ok {
		return v
	}

	return &ArrayValidator[T]{
		fn:            v.fn,
		dflt:          v.dflt,
		required:      v.required,
		minLength:     v.minLength,
		maxLength:     v.maxLength,
		invalidLength: v.invalidLength,
		convertToType: v.convertToType,
		// When we add a validation error, empty strings will be replaced
		// by the current array index. It's handeld in context.go
		validator: inner.nest(field),
	}
}

func InvalidArrayLen(min optional.Int, max optional.Int) *Invalid {
	hasMin := min.Exists
	hasMax := max.Exists

	if !hasMin && !hasMax {
		return nil
	}

	minValue := min.Value
	maxValue := max.Value

	if hasMin && hasMax {
		return &Invalid{
			Code:  utils.VAL_ARRAY_RANGE_LENGTH,
			Error: fmt.Sprintf("must have between %d and %d values", minValue, maxValue),
			Data:  RangeData(minValue, maxValue),
		}
	}

	if hasMin {
		return &Invalid{
			Code:  utils.VAL_ARRAY_MIN_LENGTH,
			Error: fmt.Sprintf("must have at least %d values", minValue),
			Data:  MinData(minValue),
		}
	}

	return &Invalid{
		Code:  utils.VAL_ARRAY_MAX_LENGTH,
		Error: fmt.Sprintf("must have no more than %d values", maxValue),
		Data:  MaxData(maxValue),
	}
}
