package validation

import (
	"github.com/valyala/fasthttp"
	"src.goblgobl.com/utils/typed"
)

func Object() *ObjectValidator {
	return &ObjectValidator{
		errReq:  Required(),
		errType: InvalidObjectType(),
	}
}

type ObjectValidator struct {
	field      Field
	required   bool
	errReq     Invalid
	errType    Invalid
	dflt       typed.Typed
	validators []InputValidator
}

func (v *ObjectValidator) Default(value typed.Typed) *ObjectValidator {
	v.dflt = value
	return v
}

func (v *ObjectValidator) Required() *ObjectValidator {
	v.required = true
	return v
}

// used when we clone a field that was required, and we want the clone to not be required
func (v *ObjectValidator) NotRequired() *ObjectValidator {
	v.required = false
	return v
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
	field := v.field
	fieldName := field.Name

	value, exists := object.ObjectIf(fieldName)
	if !exists {
		if _, exists := object[fieldName]; exists {
			res.AddInvalidField(field, v.errType)
			return
		}
		if v.required {
			res.AddInvalidField(field, v.errReq)
		} else if dflt := v.dflt; dflt != nil {
			object[fieldName] = dflt
		}
		return
	}

	for _, validator := range v.validators {
		validator.validateObjectField(value, input, res)
	}
}

func (v *ObjectValidator) validateArrayValue(value any, res *Result) {
	field := v.field
	t, ok := value.(map[string]any)
	if !ok {
		res.AddInvalidField(field, v.errType)
		return
	}
	v.Validate(typed.Typed(t), res)
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
		required:   v.required,
		dflt:       v.dflt,
		errReq:     v.errReq,
		errType:    v.errType,
		validators: validators,
	}
}
