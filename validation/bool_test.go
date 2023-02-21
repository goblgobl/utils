package validation

import (
	"testing"

	"src.goblgobl.com/tests/assert"
)

func Test_Bool_Required(t *testing.T) {
	f2 := Bool[E]().Required()
	o := Object[E]().
		Field("required", Bool[E]()).
		Field("agree", f2).Field("agree_not_required", f2.Clone().NotRequired())

	testValidator(t, o).
		FieldsHaveNoErrors("required", "agree_not_required").
		Field("agree", Required)

	testValidator(t, o, "agree", true).
		FieldsHaveNoErrors("required", "agree", "agree_not_required")
}

func Test_Bool_Type(t *testing.T) {
	o := Object[E]().Field("a", Bool[E]())

	testValidator(t, o, "a", "leto").Field("a", TypeBool)

	data, res := testValidatorData(t, o, "a", true)
	res.FieldsHaveNoErrors("a")
	assert.Equal(t, data.Bool("a"), true)

	data, res = testValidatorData(t, o, "a", false)
	res.FieldsHaveNoErrors("a")
	assert.Equal(t, data.Bool("a"), false)
}

func Test_Bool_Default(t *testing.T) {
	o := Object[E]().
		Field("a", Bool[E]().Default(true)).
		Field("b", Bool[E]().Required())

	data, res := testValidatorData(t, o)
	assert.Equal(t, data.Bool("a"), true)
	res.Field("b", Required)
}

func Test_Bool_Func(t *testing.T) {
	o := Object[E]().
		Field("f", Bool[E]().Func(func(value bool, ctx *Context[E]) any {
			if value == false {
				return true
			}
			ctx.InvalidField(TypeBool)
			return value
		}))

	data, res := testValidatorData(t, o, "f", false)
	assert.Equal(t, data.Bool("f"), true)
	res.FieldsHaveNoErrors("f")

	data, res = testValidatorData(t, o, "f", true)
	assert.Equal(t, data.Bool("f"), true)
	res.Field("f", TypeBool)
}
