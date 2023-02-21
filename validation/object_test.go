package validation

import (
	"testing"

	"src.goblgobl.com/tests/assert"
	"src.goblgobl.com/utils/typed"
)

type E struct{}

func Test_Object_Fields(t *testing.T) {
	o := Object[E]().
		Field("name", String[E]().Length(4, 0)).
		Field("status", String[E]().Required())

	testValidator(t, o, "name", 32).
		Field("name", TypeString).
		Field("status", Required)

	testValidator(t, o, "name", "leto", "status", "nope").
		FieldsHaveNoErrors("name", "status")
}

func Test_Object_Nesting(t *testing.T) {
	inner := Object[E]().
		Required().
		Field("name", String[E]().Length(4, 0)).
		Field("status", String[E]().Required())

	o := Object[E]().Field("user", inner)

	testValidator(t, o).Field("user", Required)
	testValidator(t, o, "user", 3).Field("user", TypeObject)

	testValidator(t, o, "user", map[string]any{"name": "teg"}).
		Field("user.name", InvalidStringLength(4, 0)).Field("user.status", Required)

	testValidator(t, o, "user", map[string]any{"name": "leto", "status": "bull"}).
		FieldsHaveNoErrors("user.name", "user.status")
}

func Test_Object_Deep_Nesting(t *testing.T) {
	inner := Object[E]().
		Required().
		Field("name", String[E]().Length(4, 0)).
		Field("status", String[E]().Required())

	middle := Object[E]().Field("user", inner).Required()

	o := Object[E]().Field("data", middle).Required()

	testValidator(t, o).Field("data", Required)
	testValidator(t, o, "data", 3).Field("data", TypeObject)
	testValidator(t, o, "data", map[string]any{}).Field("data.user", Required)
	testValidator(t, o, "data", map[string]any{"user": 32}).Field("data.user", TypeObject)
	testValidator(t, o, "data", map[string]any{"user": map[string]any{}}).Field("data.user.status", Required)
}

func testValidator[T any](t *testing.T, validator *ObjectValidator[T], args ...any) *assert.V {
	t.Helper()
	_, v := testValidatorData[T](t, validator, args...)
	return v
}

func testValidatorData[T any](t *testing.T, validator *ObjectValidator[T], args ...any) (typed.Typed, *assert.V) {
	t.Helper()

	input := make(typed.Typed, len(args)/2)
	for i := 0; i < len(args); i += 2 {
		input[args[i].(string)] = args[i+1]
	}
	ctx := NewContext[T](10)
	validator.ValidateInput(input, ctx)
	return ctx.Input, assert.Validation(t, ctx)
}
