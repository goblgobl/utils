package validation

import (
	"testing"

	"src.goblgobl.com/tests/assert"
)

func Test_Any_Required(t *testing.T) {
	f2 := Any[E]().Required()
	o := Object[E]().
		Field("name", Any[E]()).
		Field("code", f2).Field("code_not_required", f2.Clone().NotRequired())

	testValidator(t, o).
		FieldsHaveNoErrors("name", "code_not_required").
		Field("code", Required)

	testValidator(t, o, "code", 1).
		FieldsHaveNoErrors("code", "name", "code_not_required")
}

func Test_Any_Default(t *testing.T) {
	o := Object[E]().Field("name", Any[E]().Default(32))

	data, _ := testValidatorData(t, o)
	assert.Equal(t, data.Int("name"), 32)
}

func Test_Any_Func(t *testing.T) {
	o := Object[E]().Field("name", Any[E]().Func(func(value any, ctx *Context[E]) any {
		assert.Equal(t, value.(string), "one-one")
		return 11
	}))

	data, _ := testValidatorData(t, o, "name", "one-one")
	assert.Equal(t, data.Int("name"), 11)
}
