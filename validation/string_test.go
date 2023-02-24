package validation

import (
	"regexp"
	"strings"
	"testing"

	"src.goblgobl.com/tests/assert"
)

func Test_String_Required(t *testing.T) {
	f1 := String[E]()
	f2 := String[E]().Required()
	o := Object[E]().
		Field("name", f1).Field("name_clone", f1).
		Field("code", f2).Field("code_clone", f2).
		Field("code_not_required", f2.Clone().NotRequired())

	testValidator(t, o).
		FieldsHaveNoErrors("name", "name_clone", "code_not_required").
		Field("code", Required).
		Field("code_clone", Required)

	testValidator(t, o, "code", "1", "code_clone", "1").
		FieldsHaveNoErrors("code", "name", "code_clone", "name_clone", "code_not_required")
}

func Test_String_Nullable(t *testing.T) {
	o := Object[E]().Field("b", String[E]().Nullable())
	data, res := testValidatorData(t, o)
	assert.Nil(t, data["b"])
	res.FieldsHaveNoErrors("b")

	data, res = testValidatorData(t, o, "b", nil)
	assert.Nil(t, data["b"])
	res.FieldsHaveNoErrors("b")

	data, res = testValidatorData(t, o, "b", "hi")
	assert.Equal(t, data["b"], "hi")
	res.FieldsHaveNoErrors("b")
}

func Test_String_Default(t *testing.T) {
	f1 := String[E]().Default("leto")
	o := Object[E]().Field("a", f1)

	data, res := testValidatorData(t, o)
	assert.Equal(t, data.String("a"), "leto")
	res.FieldsHaveNoErrors("a")

	data, res = testValidatorData(t, o, "a", "ghanima")
	assert.Equal(t, data.String("a"), "ghanima")
	res.FieldsHaveNoErrors("a")
}

func Test_String_Type(t *testing.T) {
	o := Object[E]().
		Field("name", String[E]())

	for _, value := range []any{true, []string{}, 3, 5.5} {
		testValidator(t, o, "name", value).Field("name", TypeString)
	}
}

func Test_String_Length(t *testing.T) {
	f1 := String[E]().Max(3)
	f2 := String[E]().Min(2)
	f3 := String[E]().Length(2, 4)
	o := Object[E]().
		Field("f1", f1).Field("f1_clone", f1.Clone()).
		Field("f2", f2).Field("f2_clone", f2.Clone()).
		Field("f3", f3).Field("f3_clone", f3.Clone().Min(2).Max(4))

	testValidator(t, o, "f1", "1234", "f2", "1", "f3", "1", "f1_clone", "1234", "f2_clone", "1", "f3_clone", "1").
		Field("f1", InvalidStringLength(0, 3)).
		Field("f2", InvalidStringLength(2, 0)).
		Field("f3", InvalidStringLength(2, 4)).
		Field("f1_clone", InvalidStringLength(0, 3)).
		Field("f2_clone", InvalidStringLength(2, 0)).
		Field("f3_clone", InvalidStringLength(2, 4))

	testValidator(t, o, "f1", "123", "f2", "12", "f3", "12345", "f1_clone", "123", "f2_clone", "12", "f3_clone", "12345").
		FieldsHaveNoErrors("f1", "f2", "f1_clone", "f2_clone").
		Field("f3", InvalidStringLength(2, 4)).
		Field("f3_clone", InvalidStringLength(2, 4))

	testValidator(t, o, "f1", "1", "f2", "123456677", "f3", "12", "f1_clone", "1", "f2_clone", "123456677", "f3_clone", "12").
		FieldsHaveNoErrors("f1", "f2", "f3", "f1_clone", "f2_clone", "f3_clone")

	testValidator(t, o, "f3", "1234", "f3_clone", "1234").
		FieldsHaveNoErrors("f3", "f3_clone")

	testValidator(t, o, "f3", "123", "f3_clone", "123").
		FieldsHaveNoErrors("f3", "f3_clone")
}

func Test_String_Transform(t *testing.T) {
	f1 := String[E]().Transform(strings.TrimSpace).Length(2, 10)
	o := Object[E]().Field("f1", f1)

	testValidator(t, o, "f1", " 1 ").
		Field("f1", InvalidStringLength(2, 10))

	data, res := testValidatorData(t, o, "f1", " 12 ")
	res.FieldsHaveNoErrors("f1")
	assert.Equal(t, data["f1"].(string), "12")
}

func Test_String_Choice(t *testing.T) {
	f1 := String[E]().Choice("c1", "c2")
	o1 := Object[E]().
		Field("f", f1).Field("f_clone", f1.Clone())

	testValidator(t, o1, "f", "c1", "f_clone", "c2").
		FieldsHaveNoErrors("f", "f_clone")

	testValidator(t, o1, "f", "nope", "f_clone", "C2"). // case sensitive
								Field("f", InvalidStringChoice([]string{"c1", "c2"})).
								Field("f_clone", InvalidStringChoice([]string{"c1", "c2"}))
}

func Test_String_Pattern(t *testing.T) {
	f1 := String[E]().Pattern("\\d.")
	f2 := String[E]().Regexp(regexp.MustCompile("\\d."))
	o1 := Object[E]().
		Field("f", f1).Field("f2", f2).Field("f_clone", f1.Clone())

	testValidator(t, o1, "f", "1d", "f2", "1d", "f_clone", "1d").
		FieldsHaveNoErrors("f", "f_clone", "f2")

	testValidator(t, o1, "f", "1", "f2", "1", "f_clone", "1").
		Field("f", InvalidStringPattern()).
		FieldMessage("f", "is not valid"). // default/generic error
		Field("f2", InvalidStringPattern()).
		FieldMessage("f2", "is not valid"). // default/generic error
		Field("f_clone", InvalidStringPattern()).
		FieldMessage("f_clone", "is not valid") // default/generic error

	// explicit error message
	f3 := String[E]().Pattern("^\\d$", "must be a number")
	o2 := Object[E]().
		Field("f", f3).Field("f_clone", f3.Clone())

	testValidator(t, o2, "f", "1d", "f_clone", "1d").
		Field("f", InvalidStringPattern()).
		FieldMessage("f", "must be a number").
		Field("f_clone", InvalidStringPattern()).
		FieldMessage("f_clone", "must be a number")
}

func Test_String_Func(t *testing.T) {
	f1 := String[E]().Func(func(value string, ctx *Context[E]) any {
		if value == "a" {
			return "a1"
		}
		ctx.InvalidField(InvalidStringPattern())
		return value
	})

	o := Object[E]().Field("f", f1).Field("f_clone", f1.Clone())

	data, res := testValidatorData(t, o, "f", "a", "f_clone", "a")
	assert.Equal(t, data.String("f"), "a1")
	assert.Equal(t, data.String("f_clone"), "a1")
	res.FieldsHaveNoErrors("f", "f_clone")

	data, res = testValidatorData(t, o, "f", "b", "f_clone", "b")
	assert.Equal(t, data.String("f"), "b")
	assert.Equal(t, data.String("f_clone"), "b")
	res.
		Field("f", InvalidStringPattern()).
		Field("f_clone", InvalidStringPattern())
}
