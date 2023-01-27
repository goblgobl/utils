package validation

import (
	"github.com/valyala/fasthttp"
	"src.goblgobl.com/utils/typed"
)

type ArrayRule interface {
	clone() ArrayRule
	Validate(field Field, values []any, object typed.Typed, input typed.Typed, res *Result) []any
}

type ArrayTransformer func(values []any) any

func Array() *ArrayValidator {
	return &ArrayValidator{
		errReq:  Required(),
		errType: InvalidArrayType(),
	}
}

type ArrayValidator struct {
	field       Field
	required    bool
	dflt        []typed.Typed
	rules       []ArrayRule
	validator   InputValidator
	transformer ArrayTransformer
	errReq      Invalid
	errType     Invalid
}

func (v *ArrayValidator) argsToTyped(args *fasthttp.Args, t typed.Typed) {
	panic("ArrayValidator.argstoType not supported")
}

func (v *ArrayValidator) validateObjectField(object typed.Typed, input typed.Typed, res *Result) {
	field := v.field
	fieldName := field.Name

	value, exists := object[fieldName]
	if !exists {
		if v.required {
			res.AddInvalidField(field, v.errReq)
		} else if dflt := v.dflt; dflt != nil {
			object[fieldName] = dflt
		}
		return
	}

	values, ok := value.([]any)
	if !ok {
		res.AddInvalidField(field, v.errType)
	}

	// first we apply validation on the array itself (e.g. min length)
	for _, rule := range v.rules {
		values = rule.Validate(field, values, object, input, res)
	}

	// next we apply validation on every item within the array
	validator := v.validator
	res.BeginArray()
	for i, value := range values {
		res.ArrayIndex(i)
		values[i] = validator.validateArrayValue(value, res)
	}
	res.EndArray()

	if t := v.transformer; t != nil {
		object[fieldName] = t(values)
	}
}

func (v *ArrayValidator) validateArrayValue(value any, res *Result) any {
	panic("nested array validation isn't implemented yet")
}

func (v *ArrayValidator) addField(fieldName string) InputValidator {
	field := v.field.add(fieldName)

	validator := v.validator.
		addField("").
		addField(fieldName)

	rules := make([]ArrayRule, len(v.rules))
	for i, rule := range v.rules {
		rules[i] = rule.clone()
	}

	return &ArrayValidator{
		field:       field,
		required:    v.required,
		dflt:        v.dflt,
		rules:       rules,
		validator:   validator,
		errReq:      v.errReq,
		errType:     v.errType,
		transformer: v.transformer,
	}
}

func (v *ArrayValidator) Required() *ArrayValidator {
	v.required = true
	return v
}

// used when we clone a field that was required, and we want the clone to not be required
func (v *ArrayValidator) NotRequired() *ArrayValidator {
	v.required = false
	return v
}

func (v *ArrayValidator) Default(value []typed.Typed) *ArrayValidator {
	v.dflt = value
	return v
}

func (v *ArrayValidator) Validator(validator InputValidator) *ArrayValidator {
	v.validator = validator
	return v
}

func (v *ArrayValidator) Transformer(transformer ArrayTransformer) *ArrayValidator {
	v.transformer = transformer
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

func (r ArrayMin) Validate(field Field, values []any, object typed.Typed, input typed.Typed, res *Result) []any {
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

func (r ArrayMax) Validate(field Field, values []any, object typed.Typed, input typed.Typed, res *Result) []any {
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

func (r ArrayRange) Validate(field Field, values []any, object typed.Typed, input typed.Typed, res *Result) []any {
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
