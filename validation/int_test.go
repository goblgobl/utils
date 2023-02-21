package validation

import (
	"testing"

	"src.goblgobl.com/tests/assert"
	"src.goblgobl.com/utils/optional"
)

func Test_Int_Required(t *testing.T) {
	f2 := Int[E]().Required()
	o := Object[E]().
		Field("name", Int[E]()).
		Field("code", f2).Field("code_not_required", f2.Clone().NotRequired())

	testValidator(t, o).
		FieldsHaveNoErrors("name", "code_not_required").
		Field("code", Required)

	testValidator(t, o, "code", 1).
		FieldsHaveNoErrors("code", "name", "code_not_required")
}

func Test_Int_Type(t *testing.T) {
	o := Object[E]().Field("a", Int[E]())

	testValidator(t, o, "a", "leto").
		Field("a", TypeInt)

	data, res := testValidatorData(t, o, "a", "-3292")
	res.FieldsHaveNoErrors("a")
	assert.Equal(t, data.Int("a"), -3292)
}

func Test_Int_Default(t *testing.T) {
	o := Object[E]().
		Field("a", Int[E]().Default(99)).
		Field("b", Int[E]().Required())

	data, res := testValidatorData(t, o)
	assert.Equal(t, data.Int("a"), 99)
	res.Field("b", Required)
}

func Test_Int_MinMax(t *testing.T) {
	o := Object[E]().
		Field("f1", Int[E]().Min(10)).
		Field("f2", Int[E]().Max(10))

	testValidator(t, o, "f1", 9, "f2", 11).
		Field("f1", InvalidIntRange(optional.New(10), optional.NullInt)).
		Field("f2", InvalidIntRange(optional.NullInt, optional.New(10)))

	testValidator(t, o, "f1", 10, "f2", 10).
		FieldsHaveNoErrors("f1", "f2")

	testValidator(t, o, "f1", 11, "f2", 9).
		FieldsHaveNoErrors("f1", "f2")
}

func Test_Int_Range(t *testing.T) {
	o := Object[E]().
		Field("f1", Int[E]().Range(10, 20))

	for _, value := range []int{9, 21, 0, 30} {
		testValidator(t, o, "f1", value).Field("f1", InvalidIntRange(optional.New(10), optional.New(20)))
	}

	for _, value := range []int{10, 11, 19, 20} {
		testValidator(t, o, "f1", value).FieldsHaveNoErrors("f1")
	}

	testValidator(t, o, "f1", 21).Field("f1", InvalidIntRange(optional.New(10), optional.New(20)))
}

func Test_Int_Func(t *testing.T) {
	o := Object[E]().
		Field("f", Int[E]().Func(func(value int, ctx *Context[E]) any {
			if value == 9001 {
				return 9002
			}
			ctx.InvalidField(TypeInt)
			return value
		}))

	data, res := testValidatorData(t, o, "f", 9001)
	assert.Equal(t, data.Int("f"), 9002)
	res.FieldsHaveNoErrors("f")

	data, res = testValidatorData(t, o, "f", 8000)
	assert.Equal(t, data.Int("f"), 8000)
	res.Field("f", TypeInt)
}
