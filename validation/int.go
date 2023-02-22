package validation

import (
	"fmt"

	"src.goblgobl.com/utils"
	"src.goblgobl.com/utils/optional"
	"src.goblgobl.com/utils/typed"
)

type IntFuncValidator[T any] func(value int, ctx *Context[T]) any

type IntValidator[T any] struct {
	fn           IntFuncValidator[T]
	invalidValue *Invalid
	dflt         any
	minValue     optional.Int
	maxValue     optional.Int
	required     bool
	nullable     bool
}

func Int[T any]() *IntValidator[T] {
	return new(IntValidator[T])
}

func (v *IntValidator[T]) Validate(raw any, ctx *Context[T]) any {
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

	value, ok := typed.NumericToInt(raw)
	if !ok {
		ctx.InvalidField(TypeInt)
		return nil
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

func (v *IntValidator[T]) Required() *IntValidator[T] {
	v.required = true
	return v
}

// Meant to be used in conjunction with Clone(). Maybe on create, the field is
// required, but on update, it isn't.
func (v *IntValidator[T]) NotRequired() *IntValidator[T] {
	v.required = false
	return v
}

func (v *IntValidator[T]) Nullable() *IntValidator[T] {
	v.nullable = true
	return v
}

func (v *IntValidator[T]) Default(dflt any) *IntValidator[T] {
	v.dflt = dflt
	return v
}

func (v *IntValidator[T]) Min(min int) *IntValidator[T] {
	v.minValue = optional.NewInt(min)
	v.invalidValue = InvalidIntRange(v.minValue, v.maxValue)
	return v
}

func (v *IntValidator[T]) Max(max int) *IntValidator[T] {
	v.maxValue = optional.NewInt(max)
	v.invalidValue = InvalidIntRange(v.minValue, v.maxValue)
	return v
}

func (v *IntValidator[T]) Range(min int, max int) *IntValidator[T] {
	v.minValue = optional.NewInt(min)
	v.maxValue = optional.NewInt(max)
	v.invalidValue = InvalidIntRange(v.minValue, v.maxValue)
	return v
}

func (v *IntValidator[T]) Func(fn IntFuncValidator[T]) *IntValidator[T] {
	v.fn = fn
	return v
}

func (v *IntValidator[T]) Clone() *IntValidator[T] {
	return &IntValidator[T]{
		fn:       v.fn,
		dflt:     v.dflt,
		required: v.required,
		nullable: v.nullable,
		minValue: v.minValue,
		maxValue: v.maxValue,

		invalidValue: v.invalidValue,
	}
}

func InvalidIntRange(min optional.Int, max optional.Int) *Invalid {
	hasMin := min.Exists
	hasMax := max.Exists

	if !hasMin && !hasMax {
		return nil
	}

	minValue := min.Value
	maxValue := max.Value

	if hasMin && hasMax {
		return &Invalid{
			Code:  utils.VAL_INT_RANGE,
			Error: fmt.Sprintf("must be between %d and %d", minValue, maxValue),
			Data:  RangeData(minValue, maxValue),
		}
	}

	if hasMin {
		return &Invalid{
			Code:  utils.VAL_INT_MIN,
			Error: fmt.Sprintf("must be greater or equal to %d", minValue),
			Data:  MinData(minValue),
		}
	}

	return &Invalid{
		Code:  utils.VAL_INT_MAX,
		Error: fmt.Sprintf("must be less or equal to %d", maxValue),
		Data:  MaxData(maxValue),
	}
}
