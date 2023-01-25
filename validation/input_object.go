package validation

import (
	"github.com/valyala/fasthttp"
	"src.goblgobl.com/utils/typed"
)

func Object() *ObjectValidator {
	return &ObjectValidator{}
}

type ObjectValidator struct {
	field      Field
	validators []InputValidator
}

func (o *ObjectValidator) Field(fieldName string, validator InputValidator) *ObjectValidator {
	o.validators = append(o.validators, validator.addField(fieldName))
	return o
}

// object validation called on the root
func (o *ObjectValidator) Validate(input typed.Typed, res *Result) bool {
	len := res.Len()
	for _, validator := range o.validators {
		validator.validateObjectField(input, input, res)
	}
	return res.Len() == len
}

func (o *ObjectValidator) ValidateArgs(args *fasthttp.Args, res *Result) (typed.Typed, bool) {
	validators := o.validators
	input := make(typed.Typed, len(validators))
	for _, validator := range validators {
		validator.argsToTyped(args, input)
	}
	return input, o.Validate(input, res)
}

// called when the object is nested, unlike the public Validate which is
// the main entry point into validation.
func (v *ObjectValidator) validateObjectField(object typed.Typed, input typed.Typed, res *Result) {
	object = object.Object(v.field.Name)
	for _, validator := range v.validators {
		validator.validateObjectField(object, input, res)
	}
}

func (v *ObjectValidator) validateArrayValue(value any, res *Result) {
	v.Validate(value.(typed.Typed), res)
}

func (v *ObjectValidator) argsToTyped(args *fasthttp.Args, t typed.Typed) {
	panic("ObjectValidator.argstoType not supported")
}

func (v *ObjectValidator) addField(fieldName string) InputValidator {
	validators := make([]InputValidator, len(v.validators))
	for i, validator := range v.validators {
		validators[i] = validator.addField(fieldName)
	}
	field := v.field.add(fieldName)
	return &ObjectValidator{
		field:      field,
		validators: validators,
	}
}
