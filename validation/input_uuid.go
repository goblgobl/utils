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

func (v *UUIDValidator) validateObjectField(object typed.Typed, input typed.Typed, res *Result) {
	field := v.field
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
		return
	}
	v.validateValue(field, value, object, input, res)
}

func (v *UUIDValidator) validateArrayValue(value any, res *Result) {
	field := v.field
	str, ok := value.(string)
	if !ok {
		res.AddInvalidField(field, v.errType)
		return
	}
	v.validateValue(field, str, nil, nil, res)
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

func (v *UUIDValidator) Default(value string) *UUIDValidator {
	v.dflt = value
	return v
}
