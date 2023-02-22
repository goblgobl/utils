package validation

import (
	"fmt"

	"src.goblgobl.com/utils"
	"src.goblgobl.com/utils/optional"
)

type FloatFuncValidator[T any] func(value float64, ctx *Context[T]) any

type FloatValidator[T any] struct {
	dflt     any
	required bool
	nullable bool
	minValue optional.Float
	maxValue optional.Float
	fn       FloatFuncValidator[T]

	invalidValue *Invalid
}

func Float[T any]() *FloatValidator[T] {
	return new(FloatValidator[T])
}

func (v *FloatValidator[T]) Validate(raw any, ctx *Context[T]) any {
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

	value, ok := raw.(float64)
	if !ok {
		n, ok := raw.(int)
		if !ok {
			ctx.InvalidField(TypeFloat)
			return nil
		}
		value = float64(n)
	}

	if min := v.minValue; min.Exists && value < min.Value {
		ctx.InvalidField(v.invalidValue)
		return value
	}
	if max := v.maxValue; max.Exists && value > max.Value {
		ctx.InvalidField(v.invalidValue)
		return value
	}

	if fn := v.fn; fn != nil {
		return fn(value, ctx)
	}

	return value
}

func (v *FloatValidator[T]) Required() *FloatValidator[T] {
	v.required = true
	return v
}

// Meant to be used in conjunction with Clone(). Maybe on create, the field is
// required, but on update, it isn't.
func (v *FloatValidator[T]) NotRequired() *FloatValidator[T] {
	v.required = false
	return v
}

func (v *FloatValidator[T]) Nullable() *FloatValidator[T] {
	v.nullable = true
	return v
}

func (v *FloatValidator[T]) Default(dflt any) *FloatValidator[T] {
	v.dflt = dflt
	return v
}

func (v *FloatValidator[T]) Min(min float64) *FloatValidator[T] {
	v.minValue = optional.NewFloat(min)
	v.invalidValue = InvalidFloatRange(v.minValue, v.maxValue)
	return v
}

func (v *FloatValidator[T]) Max(max float64) *FloatValidator[T] {
	v.maxValue = optional.NewFloat(max)
	v.invalidValue = InvalidFloatRange(v.minValue, v.maxValue)
	return v
}

func (v *FloatValidator[T]) Range(min float64, max float64) *FloatValidator[T] {
	v.minValue = optional.NewFloat(min)
	v.maxValue = optional.NewFloat(max)
	v.invalidValue = InvalidFloatRange(v.minValue, v.maxValue)
	return v
}

func (v *FloatValidator[T]) Func(fn FloatFuncValidator[T]) *FloatValidator[T] {
	v.fn = fn
	return v
}

func (v *FloatValidator[T]) Clone() *FloatValidator[T] {
	return &FloatValidator[T]{
		fn:       v.fn,
		dflt:     v.dflt,
		required: v.required,
		nullable: v.nullable,
		minValue: v.minValue,
		maxValue: v.maxValue,

		invalidValue: v.invalidValue,
	}
}

func InvalidFloatRange(min optional.Float, max optional.Float) *Invalid {
	hasMin := min.Exists
	hasMax := max.Exists

	if !hasMin && !hasMax {
		return nil
	}

	minValue := min.Value
	maxValue := max.Value

	if hasMin && hasMax {
		return &Invalid{
			Code:  utils.VAL_FLOAT_RANGE,
			Error: fmt.Sprintf("must be between %f and %f", minValue, maxValue),
			Data:  RangeData(minValue, maxValue),
		}
	}

	if hasMin {
		return &Invalid{
			Code:  utils.VAL_FLOAT_MIN,
			Error: fmt.Sprintf("must be greater or equal to %f", minValue),
			Data:  MinData(minValue),
		}
	}

	return &Invalid{
		Code:  utils.VAL_FLOAT_MAX,
		Error: fmt.Sprintf("must be less or equal to %f", maxValue),
		Data:  MaxData(maxValue),
	}
}
