package validation

import (
	"strconv"

	"github.com/valyala/fasthttp"
	"src.goblgobl.com/utils"
	"src.goblgobl.com/utils/typed"
)

type IntRule interface {
	clone() IntRule
	Validate(field Field, value int, object typed.Typed, input typed.Typed, res *Result) int
}

func Int() *IntValidator {
	return &IntValidator{
		errReq:  Required(),
		errType: InvalidIntType(),
	}
}

type IntValidator struct {
	field    Field
	dflt     int
	required bool
	rules    []IntRule
	errReq   Invalid
	errType  Invalid
}

func (v *IntValidator) argsToTyped(args *fasthttp.Args, t typed.Typed) {
	fieldName := v.field.Name
	if value := args.Peek(fieldName); value != nil {
		if n, err := strconv.ParseInt(utils.B2S(value), 10, 0); err == nil {
			t[fieldName] = n
		} else {
			t[fieldName] = value
		}
	}
}

func (v *IntValidator) validateObjectField(object typed.Typed, input typed.Typed, res *Result) {
	field := v.field
	fieldName := field.Name
	value, exists := object.IntIf(fieldName)

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

	object[fieldName] = v.validateValue(field, value, object, input, res)
}

func (v *IntValidator) validateArrayValue(value any, res *Result) {
	field := v.field
	n, ok := typed.NumericToInt(value)
	if !ok {
		res.AddInvalidField(field, v.errType)
		return
	}
	v.validateValue(field, n, nil, nil, res)
}

func (v *IntValidator) validateValue(field Field, value int, object typed.Typed, input typed.Typed, res *Result) int {
	for _, rule := range v.rules {
		value = rule.Validate(field, value, object, input, res)
	}
	return value
}

func (v *IntValidator) addField(fieldName string) InputValidator {
	field := v.field.add(fieldName)

	rules := make([]IntRule, len(v.rules))
	for i, rule := range v.rules {
		rules[i] = rule.clone()
	}

	return &IntValidator{
		field:    field,
		dflt:     v.dflt,
		required: v.required,
		rules:    rules,
		errReq:   v.errReq,
		errType:  v.errType,
	}
}

func (v *IntValidator) Required() *IntValidator {
	v.required = true
	return v
}

// used when we clone a field that was required, and we want the clone to not be required
func (v *IntValidator) NotRequired() *IntValidator {
	v.required = false
	return v
}

func (v *IntValidator) Default(value int) *IntValidator {
	v.dflt = value
	return v
}

func (v *IntValidator) Min(min int) *IntValidator {
	v.rules = append(v.rules, IntMin{
		min: min,
		err: InvalidIntMin(min),
	})
	return v
}

func (v *IntValidator) Max(max int) *IntValidator {
	v.rules = append(v.rules, IntMax{
		max: max,
		err: InvalidIntMax(max),
	})
	return v
}

func (v *IntValidator) Range(min int, max int) *IntValidator {
	v.rules = append(v.rules, IntRange{
		min: min,
		max: max,
		err: InvalidIntRange(min, max),
	})
	return v
}

func (v *IntValidator) Func(fn func(field Field, value int, object typed.Typed, input typed.Typed, res *Result) int) *IntValidator {
	v.rules = append(v.rules, IntFunc{fn: fn})
	return v
}

type IntMin struct {
	min int
	err Invalid
}

func (r IntMin) Validate(field Field, value int, object typed.Typed, input typed.Typed, res *Result) int {
	if value < r.min {
		res.AddInvalidField(field, r.err)
	}
	return value
}

func (r IntMin) clone() IntRule {
	return IntMin{
		min: r.min,
		err: r.err,
	}
}

type IntMax struct {
	max int
	err Invalid
}

func (r IntMax) Validate(field Field, value int, object typed.Typed, input typed.Typed, res *Result) int {
	if value > r.max {
		res.AddInvalidField(field, r.err)
	}
	return value
}

func (r IntMax) clone() IntRule {
	return IntMax{
		max: r.max,
		err: r.err,
	}
}

type IntRange struct {
	min int
	max int
	err Invalid
}

func (r IntRange) Validate(field Field, value int, object typed.Typed, input typed.Typed, res *Result) int {
	if value < r.min || value > r.max {
		res.AddInvalidField(field, r.err)
	}
	return value
}

func (r IntRange) clone() IntRule {
	return IntRange{
		min: r.min,
		max: r.max,
		err: r.err,
	}
}

type IntFunc struct {
	fn func(Field, int, typed.Typed, typed.Typed, *Result) int
}

func (v IntFunc) Validate(field Field, value int, object typed.Typed, input typed.Typed, res *Result) int {
	return v.fn(field, value, object, input, res)
}

func (r IntFunc) clone() IntRule {
	return r
}
