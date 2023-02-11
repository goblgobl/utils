package validation

import (
	"github.com/valyala/fasthttp"
	"src.goblgobl.com/utils/typed"
	"src.goblgobl.com/utils/uuid"
)

func UUID() *UUIDValidator {
	return &UUIDValidator{
		errReq:  Required(),
		errType: InvalidUUIDType(),
	}
}

type UUIDValidator struct {
	field    Field
	dflt     string
	required bool
	errReq   Invalid
	errType  Invalid
}

func (v *UUIDValidator) argsToTyped(args *fasthttp.Args, t typed.Typed) {
	fieldName := v.field.Name
	if value := args.Peek(fieldName); value != nil {
		t[fieldName] = string(value)
	}
}

// This is exposed in case some caller wants to execute the validator directly
// This most likely happens when the object is being manually validated with the
// use of an object validator (i.e. Object().Func(...))
func (v *UUIDValidator) ValidateObjectField(field Field, object typed.Typed, input typed.Typed, res *Result) string {
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
		}
		return value
	}
	// The UUID validator never transforms the value in any way, so we don't have
	// to write it back into object
	v.validateValue(field, value, object, input, res)
	return value
}

func (v *UUIDValidator) validateObjectField(object typed.Typed, input typed.Typed, res *Result) {
	v.ValidateObjectField(v.field, object, input, res)
}

func (v *UUIDValidator) validateArrayValue(value any, input typed.Typed, res *Result) any {
	field := v.field
	str, ok := value.(string)
	if !ok {
		res.AddInvalidField(field, v.errType)
		return str
	}
	v.validateValue(field, str, nil, nil, res)
	return str
}

func (v *UUIDValidator) validateValue(field Field, value string, object typed.Typed, input typed.Typed, res *Result) {
	if !uuid.IsValid(value) {
		res.AddInvalidField(field, v.errType)
	}
}

func (v *UUIDValidator) addField(fieldName string) InputValidator {
	field := v.field.add(fieldName)
	return &UUIDValidator{
		field:    field,
		dflt:     v.dflt,
		required: v.required,
		errReq:   v.errReq,
		errType:  v.errType,
	}
}

func (v *UUIDValidator) Required() *UUIDValidator {
	v.required = true
	return v
}

// used when we clone a field that was required, and we want the clone to not be required
func (v *UUIDValidator) NotRequired() *UUIDValidator {
	v.required = false
	return v
}

func (v *UUIDValidator) Default(value string) *UUIDValidator {
	v.dflt = value
	return v
}
