package validation

import (
	"github.com/valyala/fasthttp"
	"src.goblgobl.com/utils/typed"
)

type ArrayRule interface {
	clone() ArrayRule
	Validate(field Field, value []typed.Typed, object typed.Typed, input typed.Typed, res *Result) []typed.Typed
}

func Array() *ArrayValidator {
	return &ArrayValidator{
		errReq:  Required(),
		errType: InvalidArrayType(),
	}
}

type ArrayValidator struct {
	field     Field
	required  bool
	dflt      []typed.Typed
	rules     []ArrayRule
	validator *ObjectValidator
	errReq    Invalid
	errType   Invalid
}

func (v *ArrayValidator) argsToTyped(args *fasthttp.Args, t typed.Typed) {
	panic("ArrayValidator.argstoType not supported")
}

func (v *ArrayValidator) validate(object typed.Typed, input typed.Typed, res *Result) {
	field := v.field
	fieldName := field.Name

	values, exists := object.ObjectsIf(fieldName)
	if !exists {
		if _, exists := object[fieldName]; !exists {
			if v.required {
				res.AddInvalidField(field, v.errReq)
			} else if dflt := v.dflt; dflt != nil {
				object[fieldName] = dflt
			}
			return
		}
		res.AddInvalidField(field, v.errType)
		return
	}

	// first we apply validation on the array itself (e.g. min length)
	for _, rule := range v.rules {
		values = rule.Validate(field, values, object, input, res)
	}

	// next we apply validation on every item within the array
	res.BeginArray()
	validator := v.validator
	for i, value := range values {
		res.ArrayIndex(i)
		validator.Validate(value, res)
	}
	res.EndArray()
}

func (v *ArrayValidator) addField(fieldName string) InputValidator {
	field := v.field.add(fieldName)

	// add an empty field, so that when we add an error to a result,
	// we can quickly detect which part of the field path needs to
	// be replaced with the index
	validator := v.validator.
		addField("").
		addField(fieldName).(*ObjectValidator)

	rules := make([]ArrayRule, len(v.rules))
	for i, rule := range v.rules {
		rules[i] = rule.clone()
	}

	return &ArrayValidator{
		field:     field,
		required:  v.required,
		dflt:      v.dflt,
		rules:     rules,
		validator: validator,
		errReq:    v.errReq,
		errType:   v.errType,
	}
}

func (v *ArrayValidator) Required() *ArrayValidator {
	v.required = true
	return v
}

func (v *ArrayValidator) Default(value []typed.Typed) *ArrayValidator {
	v.dflt = value
	return v
}

func (v *ArrayValidator) Validator(validator *ObjectValidator) *ArrayValidator {
	v.validator = validator
	return v
}

func (v *ArrayValidator) Min(min int) *ArrayValidator {
	v.rules = append(v.rules, ArrayMin{
		min: min,
		err: InvalidArrayMinLength(min),
	})
	return v
}

func (v *ArrayValidator) Max(max int) *ArrayValidator {
	v.rules = append(v.rules, ArrayMax{
		max: max,
		err: InvalidArrayMaxLength(max),
	})
	return v
}

func (v *ArrayValidator) Range(min int, max int) *ArrayValidator {
	v.rules = append(v.rules, ArrayRange{
		min: min,
		max: max,
		err: InvalidArrayRangeLength(min, max),
	})
	return v
}

type ArrayMin struct {
	min int
	err Invalid
}

func (r ArrayMin) Validate(field Field, values []typed.Typed, object typed.Typed, input typed.Typed, res *Result) []typed.Typed {
	if len(values) < r.min {
		res.AddInvalidField(field, r.err)
	}
	return values
}

func (r ArrayMin) clone() ArrayRule {
	return ArrayMin{
		min: r.min,
		err: r.err,
	}
}

type ArrayMax struct {
	max int
	err Invalid
}

func (r ArrayMax) Validate(field Field, values []typed.Typed, object typed.Typed, input typed.Typed, res *Result) []typed.Typed {
	if len(values) > r.max {
		res.AddInvalidField(field, r.err)
	}
	return values
}

func (r ArrayMax) clone() ArrayRule {
	return ArrayMax{
		max: r.max,
		err: r.err,
	}
}

type ArrayRange struct {
	min int
	max int
	err Invalid
}

func (r ArrayRange) Validate(field Field, values []typed.Typed, object typed.Typed, input typed.Typed, res *Result) []typed.Typed {
	if len(values) < r.min || len(values) > r.max {
		res.AddInvalidField(field, r.err)
	}
	return values
}

func (r ArrayRange) clone() ArrayRule {
	return ArrayRange{
		min: r.min,
		max: r.max,
		err: r.err,
	}
}
