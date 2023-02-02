package validation

import (
	"regexp"
	"strings"

	"github.com/valyala/fasthttp"
	"src.goblgobl.com/utils/typed"
)

type StringRule interface {
	clone() StringRule
	Validate(field Field, value string, object typed.Typed, input typed.Typed, res *Result) string
}

type StringTransformer func(value string) string
type StringConverter func(field Field, value string, object typed.Typed, input typed.Typed, res *Result) any
type StringFuncValidator func(field Field, value string, object typed.Typed, input typed.Typed, res *Result) string

func String() *StringValidator {
	return &StringValidator{
		errReq:  Required(),
		errType: InvalidStringType(),
	}
}

type StringValidator struct {
	field        Field
	dflt         string
	required     bool
	rules        []StringRule
	converter    StringConverter
	transformers []StringTransformer
	errReq       Invalid
	errType      Invalid
}

func (v *StringValidator) argsToTyped(args *fasthttp.Args, t typed.Typed) {
	fieldName := v.field.Name
	if value := args.Peek(fieldName); value != nil {
		t[fieldName] = string(value)
	}
}

// This is exposed in case some caller wants to execute the validator directly
// This most likely happens when the object is being manually validated with the
// use of an object validator (i.e. Object().Func(...))
func (v *StringValidator) ValidateObjectField(field Field, object typed.Typed, input typed.Typed, res *Result) any {
	fieldName := field.Name

	value, exists := object.StringIf(fieldName)
	if !exists {
		if _, exists = object[fieldName]; !exists && v.required {
			res.AddInvalidField(field, v.errReq)
		} else if exists {
			res.AddInvalidField(field, v.errType)
		}
		if dflt := v.dflt; dflt != "" {
			object[fieldName] = dflt
			return dflt
		}
		return value
	}

	validated := v.validateValue(field, value, object, input, res)
	object[fieldName] = validated
	return validated
}

// this is called internally when we're validating an object and the nested fields
func (v *StringValidator) validateObjectField(object typed.Typed, input typed.Typed, res *Result) {
	v.ValidateObjectField(v.field, object, input, res)
}

func (v *StringValidator) validateArrayValue(value any, res *Result) any {
	field := v.field
	str, ok := value.(string)
	if !ok {
		res.AddInvalidField(field, v.errType)
		return ""
	}
	return v.validateValue(field, str, nil, nil, res)
}

func (v *StringValidator) validateValue(field Field, value string, object typed.Typed, input typed.Typed, res *Result) any {
	for _, rule := range v.rules {
		value = rule.Validate(field, value, object, input, res)
	}

	for _, transformer := range v.transformers {
		value = transformer(value)
	}

	if converter := v.converter; converter != nil {
		return converter(field, value, object, input, res)
	}
	return value
}

func (v *StringValidator) addField(fieldName string) InputValidator {
	field := v.field.add(fieldName)

	rules := make([]StringRule, len(v.rules))
	for i, rule := range v.rules {
		rules[i] = rule.clone()
	}

	return &StringValidator{
		field:        field,
		dflt:         v.dflt,
		required:     v.required,
		converter:    v.converter,
		transformers: v.transformers,
		rules:        rules,
		errReq:       v.errReq,
		errType:      v.errType,
	}
}

func (v *StringValidator) Required() *StringValidator {
	v.required = true
	return v
}

// used when we clone a field that was required, and we want the clone to not be required
func (v *StringValidator) NotRequired() *StringValidator {
	v.required = false
	return v
}

func (v *StringValidator) Default(value string) *StringValidator {
	v.dflt = value
	return v
}

func (v *StringValidator) Choice(valid ...string) *StringValidator {
	v.rules = append(v.rules, StringChoice{
		valid: valid,
		err:   InvalidStringChoice(valid),
	})
	return v
}

func (v *StringValidator) Length(min int, max int) *StringValidator {
	v.rules = append(v.rules, StringLen{
		min: min,
		max: max,
		err: InvalidStringLength(min, max),
	})
	return v
}

func (v *StringValidator) TrimSpace() *StringValidator {
	v.rules = append(v.rules, StringTrimSpace{})
	return v
}

func (v *StringValidator) Pattern(pattern string, errorMessage ...string) *StringValidator {
	v.rules = append(v.rules, StringPattern{
		pattern: regexp.MustCompile(pattern),
		err:     InvalidStringPattern(errorMessage...),
	})
	return v
}

func (v *StringValidator) Func(fn StringFuncValidator) *StringValidator {
	v.rules = append(v.rules, StringFunc{fn})
	return v
}

func (v *StringValidator) Convert(fn StringConverter) *StringValidator {
	v.converter = fn
	return v
}

func (v *StringValidator) Transformer(transformer StringTransformer) *StringValidator {
	v.transformers = append(v.transformers, transformer)
	return v
}

type StringLen struct {
	min int
	max int
	err Invalid
}

func (r StringLen) Validate(field Field, value string, object typed.Typed, input typed.Typed, res *Result) string {
	if min := r.min; min > 0 && len(value) < min {
		res.AddInvalidField(field, r.err)
	}
	if max := r.max; max > 0 && len(value) > max {
		res.AddInvalidField(field, r.err)
	}
	return value
}

func (r StringLen) clone() StringRule {
	return StringLen{
		min: r.min,
		max: r.max,
		err: r.err,
	}
}

type StringPattern struct {
	pattern *regexp.Regexp
	err     Invalid
}

func (r StringPattern) Validate(field Field, value string, object typed.Typed, input typed.Typed, res *Result) string {
	if !r.pattern.MatchString(value) {
		res.AddInvalidField(field, r.err)
	}
	return value
}

func (r StringPattern) clone() StringRule {
	return StringPattern{
		pattern: r.pattern,
		err:     r.err,
	}
}

type StringChoice struct {
	valid []string
	err   Invalid
}

func (r StringChoice) Validate(field Field, value string, object typed.Typed, input typed.Typed, res *Result) string {
	for _, valid := range r.valid {
		if value == valid {
			return value
		}
	}
	res.AddInvalidField(field, r.err)
	return value
}

func (r StringChoice) clone() StringRule {
	return StringChoice{
		valid: r.valid,
		err:   r.err,
	}
}

type StringFunc struct {
	fn StringFuncValidator
}

func (r StringFunc) Validate(field Field, value string, object typed.Typed, input typed.Typed, res *Result) string {
	return r.fn(field, value, object, input, res)
}

func (r StringFunc) clone() StringRule {
	return r
}

type StringTrimSpace struct {
}

func (r StringTrimSpace) Validate(field Field, value string, object typed.Typed, input typed.Typed, res *Result) string {
	return strings.TrimSpace(value)
}

func (r StringTrimSpace) clone() StringRule {
	return r
}
