package validation

import (
	"testing"

	"src.goblgobl.com/tests/assert"
	"src.goblgobl.com/utils/optional"
	"src.goblgobl.com/utils/typed"
)

func Test_Array_Simple(t *testing.T) {
	a := Array[E]().Validator(String[E]().Length(4, 0)).Required()
	o := Object[E]().Field("names", a)

	testValidator(t, o).Field("names", Required)
	testValidator(t, o, "names", 3).Field("names", TypeArray)
	testValidator(t, o, "names", []any{"hello", "teg", 3}).
		Field("names.2", TypeString).
		Field("names.1", InvalidStringLength(4, 0))
}

func Test_Array_Nesting(t *testing.T) {
	innerA := Array[E]().Validator(String[E]().Length(4, 0)).Required()
	innerO := Object[E]().Field("names", innerA)
	outerA := Array[E]().Validator(innerO).Required()
	outerO := Object[E]().Field("entries", outerA)

	testValidator(t, outerO).Field("entries", Required)
	testValidator(t, outerO, "entries", 3).Field("entries", TypeArray)
	testValidator(t, outerO, "entries", []any{3}).Field("entries.0", TypeObject)
	testValidator(t, outerO, "entries", []any{
		map[string]any{"names": 4},
	}).Field("entries.0.names", TypeArray)
}

func Test_Array_MinAndMax(t *testing.T) {
	createItem := func() typed.Typed {
		return typed.Typed{"name": "n"}
	}

	child := Object[E]().Field("name", String[E]())
	o1 := Object[E]().Field("users", Array[E]().Min(2).Max(3).Required().Validator(child))

	testValidator(t, o1, "users", []any{createItem()}).
		Field("users", InvalidArrayLen(optional.New(2), optional.New(3)))

	// 4 items, too many
	testValidator(t, o1, "users", []any{
		createItem(), createItem(), createItem(), createItem(),
	}).Field("users", InvalidArrayLen(optional.New(2), optional.New(3)))

	// 2 items, good
	testValidator(t, o1, "users", []any{
		createItem(), createItem(),
	}).FieldsHaveNoErrors("users")

	// 3 items, good
	testValidator(t, o1, "users", []any{
		createItem(), createItem(), createItem(),
	}).FieldsHaveNoErrors("users")
}

func Test_ConvertToType_Valid(t *testing.T) {
	o1 := Object[E]().Field(
		"values", Array[E]().ConvertToType().Validator(String[E]()),
	)
	data, _ := testValidatorData(t, o1, "values", []any{"a", "bb"})
	assert.List(t, data["values"].([]string), []string{"a", "bb"})

	o2 := Object[E]().Field(
		"values", Array[E]().ConvertToType().Validator(Int[E]()),
	)
	data, _ = testValidatorData(t, o2, "values", []any{33, 6668})
	assert.List(t, data["values"].([]int), []int{33, 6668})

	o3 := Object[E]().Field(
		"values", Array[E]().ConvertToType().Validator(Bool[E]()),
	)
	data, _ = testValidatorData(t, o3, "values", []any{true, true, false})
	assert.List(t, data["values"].([]bool), []bool{true, true, false})

	o4 := Object[E]().Field(
		"values", Array[E]().ConvertToType().Validator(Float[E]()),
	)
	data, _ = testValidatorData(t, o4, "values", []any{1.2, 21.3})
	assert.List(t, data["values"].([]float64), []float64{1.2, 21.3})
}
