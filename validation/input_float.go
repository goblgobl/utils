package validation

import (
	"strconv"

	"github.com/valyala/fasthttp"
	"src.goblgobl.com/utils"
	"src.goblgobl.com/utils/typed"
)

type FloatRule interface {
	clone() FloatRule
	Validate(field Field, value float64, object typed.Typed, input typed.Typed, res *Result) float64
}

func Float() *FloatValidator {
	return &FloatValidator{
		errReq:  Required(),
		errType: InvalidFloatType(),
	}
}

type FloatValidator struct {
	field    Field
	dflt     float64
	required bool
	rules    []FloatRule
	errReq   Invalid
	errType  Invalid
}

func (v *FloatValidator) argsToTyped(args *fasthttp.Args, t typed.Typed) {
	fieldName := v.field.Name
	if value := args.Peek(fieldName); value != nil {
		if n, err := strconv.ParseFloat(utils.B2S(value), 64); err == nil {
			t[fieldName] = n
		} else {
			t[fieldName] = value
		}
	}
}

func (v *FloatValidator) validate(object typed.Typed, input typed.Typed, res *Result) {
	field := v.field
	fieldName := field.Name
	value, exists := object.FloatIf(fieldName)

	if !exists {
		if _, exists := object[fieldName]; !exists {
			if v.required {
				res.AddInvalidField(field, v.errReq)
			} else if dflt := v.dflt; dflt != 0 {
				object[fieldName] = dflt
			}
			return
		}
		res.AddInvalidField(field, v.errType)
		return
	}

	for _, rule := range v.rules {
		value = rule.Validate(field, value, object, input, res)
	}
	object[fieldName] = value
}

func (v *FloatValidator) addField(fieldName string) InputValidator {
	field := v.field.add(fieldName)

	rules := make([]FloatRule, len(v.rules))
	for i, rule := range v.rules {
		rules[i] = rule.clone()
	}

	return &FloatValidator{
		field:    field,
		dflt:     v.dflt,
		required: v.required,
		rules:    rules,
		errReq:   v.errReq,
		errType:  v.errType,
	}
}

func (v *FloatValidator) Required() *FloatValidator {
	v.required = true
	return v
}

func (v *FloatValidator) Default(value float64) *FloatValidator {
	v.dflt = value
	return v
}

func (v *FloatValidator) Min(min float64) *FloatValidator {
	v.rules = append(v.rules, FloatMin{
		min: min,
		err: InvalidFloatMin(min),
	})
	return v
}

func (v *FloatValidator) Max(max float64) *FloatValidator {
	v.rules = append(v.rules, FloatMax{
		max: max,
		err: InvalidFloatMax(max),
	})
	return v
}

func (v *FloatValidator) Range(min float64, max float64) *FloatValidator {
	v.rules = append(v.rules, FloatRange{
		min: min,
		max: max,
		err: InvalidFloatRange(min, max),
	})
	return v
}

func (v *FloatValidator) Func(fn func(field Field, value float64, object typed.Typed, input typed.Typed, res *Result) float64) *FloatValidator {
	v.rules = append(v.rules, FloatFunc{fn: fn})
	return v
}

type FloatMin struct {
	min float64
	err Invalid
}

func (r FloatMin) Validate(field Field, value float64, object typed.Typed, input typed.Typed, res *Result) float64 {
	if value < r.min {
		res.AddInvalidField(field, r.err)
	}
	return value
}

func (r FloatMin) clone() FloatRule {
	return FloatMin{
		min: r.min,
		err: r.err,
	}
}

type FloatMax struct {
	max float64
	err Invalid
}

func (r FloatMax) Validate(field Field, value float64, object typed.Typed, input typed.Typed, res *Result) float64 {
	if value > r.max {
		res.AddInvalidField(field, r.err)
	}
	return value
}

func (r FloatMax) clone() FloatRule {
	return FloatMax{
		max: r.max,
		err: r.err,
	}
}

type FloatRange struct {
	min float64
	max float64
	err Invalid
}

func (r FloatRange) Validate(field Field, value float64, object typed.Typed, input typed.Typed, res *Result) float64 {
	if value < r.min || value > r.max {
		res.AddInvalidField(field, r.err)
	}
	return value
}

func (r FloatRange) clone() FloatRule {
	return FloatRange{
		min: r.min,
		max: r.max,
		err: r.err,
	}
}

type FloatFunc struct {
	fn func(Field, float64, typed.Typed, typed.Typed, *Result) float64
}

func (v FloatFunc) Validate(field Field, value float64, object typed.Typed, input typed.Typed, res *Result) float64 {
	return v.fn(field, value, object, input, res)
}

func (r FloatFunc) clone() FloatRule {
	return r
}
