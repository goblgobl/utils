package validation

import (
	"testing"

	"src.goblgobl.com/tests/assert"
	"src.goblgobl.com/utils/optional"
)

func Test_Float_Required(t *testing.T) {
	f2 := Float[E]().Required()
	o := Object[E]().
		Field("name", Float[E]()).
		Field("code", f2).Field("code_not_required", f2.Clone().NotRequired())

	testValidator(t, o).
		FieldsHaveNoErrors("name", "code_not_required").
		Field("code", Required)

	testValidator(t, o, "code", 1.3).
		FieldsHaveNoErrors("code", "name", "code_not_required")

	testValidator(t, o, "code", 1).
		FieldsHaveNoErrors("code", "name", "code_not_required")
}

func Test_Float_Type(t *testing.T) {
	o := Object[E]().Field("a", Float[E]())

	testValidator(t, o, "a", "leto").
		Field("a", TypeFloat)

	data, res := testValidatorData(t, o, "a", -3292.3)
	res.FieldsHaveNoErrors("a")
	assert.Equal(t, data.Float("a"), -3292.3)
}

func Test_Float_Default(t *testing.T) {
	o := Object[E]().
		Field("a", Float[E]().Default(99.1)).
		Field("b", Float[E]().Required())

	data, res := testValidatorData(t, o)
	assert.Equal(t, data.Float("a"), 99.1)
	res.Field("b", Required)
}

func Test_Float_MinMax(t *testing.T) {
	o := Object[E]().
		Field("f1", Float[E]().Min(10.1)).
		Field("f2", Float[E]().Max(10.1))

	testValidator(t, o, "f1", 10, "f2", 10.2).
		Field("f1", InvalidFloatRange(optional.New(10.1), optional.NullFloat)).
		Field("f2", InvalidFloatRange(optional.NullFloat, optional.New(10.1)))

	testValidator(t, o, "f1", 10.1, "f2", 10.1).
		FieldsHaveNoErrors("f1", "f2")

	testValidator(t, o, "f1", 10.2, "f2", 9.9).
		FieldsHaveNoErrors("f1", "f2")
}

func Test_Float_Range(t *testing.T) {
	o := Object[E]().
		Field("f1", Float[E]().Range(10.1, 20.2))

	for _, value := range []float64{10.04, 20.201, 0, 30} {
		testValidator(t, o, "f1", value).Field("f1", InvalidFloatRange(optional.New(10.1), optional.New(20.2)))
	}

	for _, value := range []float64{10.1, 11, 19, 18.82, 20.2} {
		testValidator(t, o, "f1", value).FieldsHaveNoErrors("f1")
	}

	testValidator(t, o, "f1", 20.3).Field("f1", InvalidFloatRange(optional.New(10.1), optional.New(20.2)))
}

func Test_Float_Func(t *testing.T) {
	o := Object[E]().
		Field("f", Float[E]().Func(func(value float64, ctx *Context[E]) any {
			if value == 9001.99 {
				return 9002.22
			}
			ctx.InvalidField(TypeFloat)
			return value
		}))

	data, res := testValidatorData(t, o, "f", 9001.99)
	assert.Equal(t, data.Float("f"), 9002.22)
	res.FieldsHaveNoErrors("f")

	data, res = testValidatorData(t, o, "f", 8000.1)
	assert.Equal(t, data.Float("f"), 8000.1)
	res.Field("f", TypeFloat)
}
