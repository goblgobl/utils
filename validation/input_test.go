package validation

import (
	"encoding/hex"
	"testing"

	"github.com/valyala/fasthttp"
	"src.goblgobl.com/tests/assert"
	"src.goblgobl.com/utils/ascii"
	"src.goblgobl.com/utils/typed"
)

func Test_String_Required(t *testing.T) {
	f1 := String()
	f2 := String().Required()
	o := Object().
		Field("name", f1).Field("name_clone", f1).
		Field("code", f2).Field("code_clone", f2).Field("code_not_required", f2.NotRequired())

	_, res := testInput(o)
	assert.Validation(t, res).
		FieldsHaveNoErrors("name", "name_clone", "code_not_required").
		Field("code", Required()).
		Field("code_clone", Required())

	_, res = testInput(o, "code", "1", "code_clone", "1")
	assert.Validation(t, res).
		FieldsHaveNoErrors("code", "name", "code_clone", "name_clone", "code_not_required")
}

func Test_String_Default(t *testing.T) {
	f1 := String().Default("leto")
	f2 := String().Required().Default("leto")
	o := Object().
		Field("a", f1).Field("a_clone", f1).
		Field("b", f2).Field("b_clone", f2)

	// default doesn't really make sense with required, required
	// takes precedence
	data, res := testInput(o)
	assert.Equal(t, data.String("a"), "leto")
	assert.Equal(t, data.String("a_clone"), "leto")
	assert.Validation(t, res).
		Field("b", Required()).
		Field("b_clone", Required())
}

func Test_String_Type(t *testing.T) {
	o := Object().
		Field("name", String())

	_, res := testInput(o, "name", 3)
	assert.Validation(t, res).
		Field("name", InvalidStringType())
}

func Test_String_Length(t *testing.T) {
	f1 := String().Length(0, 3)
	f2 := String().Length(2, 0)
	f3 := String().Length(2, 4)
	o := Object().
		Field("f1", f1).Field("f1_clone", f1).
		Field("f2", f2).Field("f2_clone", f2).
		Field("f3", f3).Field("f3_clone", f3)

	_, res := testInput(o, "f1", "1234", "f2", "1", "f3", "1", "f1_clone", "1234", "f2_clone", "1", "f3_clone", "1")
	assert.Validation(t, res).
		Field("f1", InvalidStringLength(0, 3)).
		Field("f2", InvalidStringLength(2, 0)).
		Field("f3", InvalidStringLength(2, 4)).
		Field("f1_clone", InvalidStringLength(0, 3)).
		Field("f2_clone", InvalidStringLength(2, 0)).
		Field("f3_clone", InvalidStringLength(2, 4))

	_, res = testInput(o, "f1", "123", "f2", "12", "f3", "12345", "f1_clone", "123", "f2_clone", "12", "f3_clone", "12345")
	assert.Validation(t, res).
		FieldsHaveNoErrors("f1", "f2", "f1_clone", "f2_clone").
		Field("f3", InvalidStringLength(2, 4)).
		Field("f3_clone", InvalidStringLength(2, 4))

	_, res = testInput(o, "f1", "1", "f2", "123456677", "f3", "12", "f1_clone", "1", "f2_clone", "123456677", "f3_clone", "12")
	assert.Validation(t, res).
		FieldsHaveNoErrors("f1", "f2", "f3", "f1_clone", "f2_clone", "f3_clone")

	_, res = testInput(o, "f3", "1234", "f3_clone", "1234")
	assert.Validation(t, res).
		FieldsHaveNoErrors("f3", "f3_clone")

	_, res = testInput(o, "f3", "123", "f3_clone", "123")
	assert.Validation(t, res).
		FieldsHaveNoErrors("f3", "f3_clone")
}

func Test_String_Choice(t *testing.T) {
	f1 := String().Choice("c1", "c2")
	o1 := Object().
		Field("f", f1).Field("f_clone", f1)

	_, res := testInput(o1, "f", "c1", "f_clone", "c2")
	assert.Validation(t, res).
		FieldsHaveNoErrors("f", "f_clone")

	_, res = testInput(o1, "f", "nope", "f_clone", "C2") // case sensitive
	assert.Validation(t, res).
		Field("f", InvalidStringChoice([]string{"c1", "c2"})).
		Field("f_clone", InvalidStringChoice([]string{"c1", "c2"}))
}

func Test_String_Pattern(t *testing.T) {
	f1 := String().Pattern("\\d.")
	o1 := Object().
		Field("f", f1).Field("f_clone", f1)

	_, res := testInput(o1, "f", "1d", "f_clone", "1d")
	assert.Validation(t, res).
		FieldsHaveNoErrors("f", "f_clone")

	_, res = testInput(o1, "f", "1", "f_clone", "1")
	assert.Validation(t, res).
		Field("f", InvalidStringPattern()).
		FieldMessage("f", "is not valid"). // default/generic error
		Field("f_clone", InvalidStringPattern()).
		FieldMessage("f_clone", "is not valid") // default/generic error

	// explicit error message
	f2 := String().Pattern("^\\d$", "must be a number")
	o2 := Object().
		Field("f", f2).Field("f_clone", f2)

	_, res = testInput(o2, "f", "1d", "f_clone", "1d")
	assert.Validation(t, res).
		Field("f", InvalidStringPattern()).
		FieldMessage("f", "must be a number").
		Field("f_clone", InvalidStringPattern()).
		FieldMessage("f_clone", "must be a number")
}

func Test_String_Func(t *testing.T) {
	f1 := String().Func(func(field Field, value string, object typed.Typed, input typed.Typed, res *Result) string {
		if value == "a" {
			return "a1"
		}
		res.AddInvalidField(field, InvalidStringPattern())
		return value
	})

	o := Object().Field("f", f1).Field("f_clone", f1)

	data, res := testInput(o, "f", "a", "f_clone", "a")
	assert.Equal(t, data.String("f"), "a1")
	assert.Equal(t, data.String("f_clone"), "a1")
	assert.Validation(t, res).
		FieldsHaveNoErrors("f", "f_clone")

	data, res = testInput(o, "f", "b", "f_clone", "b")
	assert.Equal(t, data.String("f"), "b")
	assert.Equal(t, data.String("f_clone"), "b")
	assert.Validation(t, res).
		Field("f", InvalidStringPattern()).
		Field("f_clone", InvalidStringPattern())
}

func Test_String_Converter(t *testing.T) {
	f1 := String().Convert(func(field Field, value string, object typed.Typed, input typed.Typed, res *Result) any {
		b, err := hex.DecodeString(value)
		if err == nil {
			return b
		}
		res.AddInvalidField(field, InvalidStringPattern())
		return nil
	})

	o := Object().Field("f", f1).Field("f_clone", f1)

	data, res := testInput(o, "f", "FFFe", "f_clone", "FFFe")
	assert.Bytes(t, data.Bytes("f"), []byte{255, 254})
	assert.Bytes(t, data.Bytes("f_clone"), []byte{255, 254})
	assert.Validation(t, res).
		FieldsHaveNoErrors("f", "f_clone")

	data, res = testInput(o, "f", "z", "f_clone", "z")
	assert.True(t, data.Bytes("f") == nil)
	assert.True(t, data.Bytes("f_clone") == nil)
	assert.Validation(t, res).
		Field("f", InvalidStringPattern()).
		Field("f_clone", InvalidStringPattern())
}

func Test_String_Args(t *testing.T) {
	o := Object().
		Field("name", String().Required().Length(4, 4))

	_, res := testArgs(o, "name", "leto")
	assert.Validation(t, res).FieldsHaveNoErrors("name")
}

func Test_String_Transformer(t *testing.T) {
	o := Object().
		Field("name", String().Transformer(ascii.Lowercase))

	data, res := testInput(o, "name", "LeTO_9001 !!")
	assert.Validation(t, res).FieldsHaveNoErrors("name")
	assert.Equal(t, data.String("name"), "leto_9001 !!")
}

func Test_String_Multiple_Transformer(t *testing.T) {
	o := Object().
		Field("name", String().Transformer(func(input string) string {
			return input + "ZZ"
		}).Transformer(ascii.Lowercase))

	data, res := testInput(o, "name", "AB")
	assert.Validation(t, res).FieldsHaveNoErrors("name")
	assert.Equal(t, data.String("name"), "abzz")
}

func Test_Int_Required(t *testing.T) {
	f2 := Int().Required()
	o := Object().
		Field("name", Int()).
		Field("code", f2).Field("code_not_required", f2.NotRequired())

	_, res := testInput(o)
	assert.Validation(t, res).
		FieldsHaveNoErrors("name", "code_not_required").
		Field("code", Required())

	_, res = testInput(o, "code", 1)
	assert.Validation(t, res).
		FieldsHaveNoErrors("code", "name", "code_not_required")
}

func Test_Int_Type(t *testing.T) {
	o := Object().
		Field("a", Int())

	_, res := testInput(o, "a", "leto")
	assert.Validation(t, res).
		Field("a", InvalidIntType())

	data, res := testInput(o, "a", "-3292")
	assert.Validation(t, res).
		FieldsHaveNoErrors("a")
	assert.Equal(t, data.Int("a"), -3292)
}

func Test_Int_Default(t *testing.T) {
	o := Object().
		Field("a", Int().Default(99)).
		Field("b", Int().Required().Default(88))

	// default doesn't really make sense with required, required
	// takes precedence
	data, res := testInput(o)
	assert.Equal(t, data.Int("a"), 99)
	assert.Validation(t, res).
		Field("b", Required())
}

func Test_Int_MinMax(t *testing.T) {
	o := Object().
		Field("f1", Int().Min(10)).
		Field("f2", Int().Max(10))

	_, res := testInput(o, "f1", 9, "f2", 11)
	assert.Validation(t, res).
		Field("f1", InvalidIntMin(10)).
		Field("f2", InvalidIntMax(10))

	_, res = testInput(o, "f1", 10, "f2", 10)
	assert.Validation(t, res).
		FieldsHaveNoErrors("f1", "f2")

	_, res = testInput(o, "f1", 11, "f2", 9)
	assert.Validation(t, res).
		FieldsHaveNoErrors("f1", "f2")
}

func Test_Int_Range(t *testing.T) {
	o := Object().
		Field("f1", Int().Range(10, 20))

	for _, value := range []int{9, 21, 0, 30} {
		_, res := testInput(o, "f1", value)
		assert.Validation(t, res).
			Field("f1", InvalidIntRange(10, 20))
	}

	for _, value := range []int{10, 11, 19, 20} {
		_, res := testInput(o, "f1", value)
		assert.Validation(t, res).
			FieldsHaveNoErrors("f1")
	}

	_, res := testInput(o, "f1", 21)
	assert.Validation(t, res).
		Field("f1", InvalidIntRange(10, 20))
}

func Test_Int_Func(t *testing.T) {
	o := Object().
		Field("f", Int().Func(func(field Field, value int, object typed.Typed, input typed.Typed, res *Result) int {
			if value == 9001 {
				return 9002
			}
			res.AddInvalidField(field, InvalidIntMax(444))
			return value
		}))

	data, res := testInput(o, "f", 9001)
	assert.Equal(t, data.Int("f"), 9002)
	assert.Validation(t, res).
		FieldsHaveNoErrors("f")

	data, res = testInput(o, "f", 8000)
	assert.Equal(t, data.Int("f"), 8000)
	assert.Validation(t, res).
		Field("f", InvalidIntMax(444))
}

func Test_Int_Args(t *testing.T) {
	o := Object().Field("id", Int().Required().Range(4, 4))
	input, res := testArgs(o, "id", "4")
	assert.Validation(t, res).FieldsHaveNoErrors("id")
	assert.Equal(t, input.Int("id"), 4)

	input, res = testArgs(o, "id", "nope")
	assert.Validation(t, res).Field("id", InvalidIntType())
	assert.Equal(t, input.IntOr("id", -1), -1)
}

func Test_Float_Required(t *testing.T) {
	f2 := Float().Required()
	o := Object().
		Field("name", Float()).
		Field("code", f2).Field("code_not_required", f2.NotRequired())

	_, res := testInput(o)
	assert.Validation(t, res).
		FieldsHaveNoErrors("name", "code_not_required").
		Field("code", Required())

	_, res = testInput(o, "code", 1.2)
	assert.Validation(t, res).
		FieldsHaveNoErrors("code", "name", "code_not_required")

	// accepts ints
	_, res = testInput(o, "code", 2)
	assert.Validation(t, res).
		FieldsHaveNoErrors("code", "name", "code_not_required")
}

func Test_Float_Type(t *testing.T) {
	o := Object().
		Field("a", Float())

	_, res := testInput(o, "a", "leto")
	assert.Validation(t, res).
		Field("a", InvalidFloatType())

	data, res := testInput(o, "a", "-3292")
	assert.Validation(t, res).
		FieldsHaveNoErrors("a")
	assert.Equal(t, data.Float("a"), -3292)
}

func Test_Float_Default(t *testing.T) {
	o := Object().
		Field("a", Float().Default(99)).
		Field("b", Float().Required().Default(88))

	// default doesn't really make sense with required, required
	// takes precedence
	data, res := testInput(o)
	assert.Equal(t, data.Float("a"), 99)
	assert.Validation(t, res).
		Field("b", Required())
}

func Test_Float_MinMax(t *testing.T) {
	o := Object().
		Field("f1", Float().Min(10.3)).
		Field("f2", Float().Max(10.3))

	_, res := testInput(o, "f1", 10.2, "f2", 10.4)
	assert.Validation(t, res).
		Field("f1", InvalidFloatMin(10.3)).
		Field("f2", InvalidFloatMax(10.3))

	_, res = testInput(o, "f1", 10.3, "f2", 10.3)
	assert.Validation(t, res).
		FieldsHaveNoErrors("f1", "f2")

	_, res = testInput(o, "f1", 10.4, "f2", 10.2)
	assert.Validation(t, res).
		FieldsHaveNoErrors("f1", "f2")
}

func Test_Float_Range(t *testing.T) {
	o := Object().
		Field("f1", Float().Range(10.1, 20.2))

	for _, value := range []float64{9, 10.0, 20.3, 21, 0, 30} {
		_, res := testInput(o, "f1", value)
		assert.Validation(t, res).
			Field("f1", InvalidFloatRange(10.1, 20.2))
	}

	for _, value := range []float64{10.1, 10.2, 11.4, 19.2, 20, 20.2} {
		_, res := testInput(o, "f1", value)
		assert.Validation(t, res).
			FieldsHaveNoErrors("f1")
	}

	_, res := testInput(o, "f1", 20.3)
	assert.Validation(t, res).
		Field("f1", InvalidFloatRange(10.1, 20.2))
}

func Test_Float_Func(t *testing.T) {
	o := Object().
		Field("f", Float().Func(func(field Field, value float64, object typed.Typed, input typed.Typed, res *Result) float64 {
			if value == 9001.0 {
				return 9002.1
			}
			res.AddInvalidField(field, InvalidFloatMax(32.2))
			return value
		}))

	data, res := testInput(o, "f", 9001.0)
	assert.Equal(t, data.Float("f"), 9002.1)
	assert.Validation(t, res).
		FieldsHaveNoErrors("f")

	data, res = testInput(o, "f", 8001.2)
	assert.Equal(t, data.Float("f"), 8001.2)
	assert.Validation(t, res).
		Field("f", InvalidFloatMax(32.2))
}

func Test_Float_Args(t *testing.T) {
	o := Object().Field("id", Float().Required().Range(4, 4))
	input, res := testArgs(o, "id", "4")
	assert.Validation(t, res).FieldsHaveNoErrors("id")
	assert.Equal(t, input.Float("id"), 4)

	input, res = testArgs(o, "id", "nope")
	assert.Validation(t, res).Field("id", InvalidFloatType())
	assert.Equal(t, input.FloatOr("id", -1), -1)
}

func Test_Bool_Required(t *testing.T) {
	f2 := Bool().Required()
	o := Object().
		Field("required", Bool()).
		Field("agree", f2).Field("agree_not_required", f2.NotRequired())

	_, res := testInput(o)
	assert.Validation(t, res).
		FieldsHaveNoErrors("required", "agree_not_required").
		Field("agree", Required())

	_, res = testInput(o, "agree", true)
	assert.Validation(t, res).
		FieldsHaveNoErrors("required", "agree", "agree_not_required")
}

func Test_Bool_Type(t *testing.T) {
	o := Object().
		Field("a", Bool())

	_, res := testInput(o, "a", "leto")
	assert.Validation(t, res).
		Field("a", InvalidBoolType())

	data, res := testInput(o, "a", "true")
	assert.Validation(t, res).
		FieldsHaveNoErrors("a")
	assert.Equal(t, data.Bool("a"), true)
}

func Test_Bool_Default(t *testing.T) {
	o := Object().
		Field("a", Bool().Default(true)).
		Field("b", Bool().Required().Default(true))

	// default doesn't really make sense with required, required
	// takes precedence
	data, res := testInput(o)
	assert.Equal(t, data.Bool("a"), true)
	assert.Validation(t, res).
		Field("b", Required())
}

func Test_Bool_Func(t *testing.T) {
	o := Object().
		Field("f", Bool().Func(func(field Field, value bool, object typed.Typed, input typed.Typed, res *Result) bool {
			if value == false {
				return true
			}
			res.AddInvalidField(field, InvalidBoolType())
			return value
		}))

	data, res := testInput(o, "f", false)
	assert.Equal(t, data.Bool("f"), true)
	assert.Validation(t, res).
		FieldsHaveNoErrors("f")

	data, res = testInput(o, "f", true)
	assert.Equal(t, data.Bool("f"), true)
	assert.Validation(t, res).
		Field("f", InvalidBoolType())
}

func Test_Bool_Args(t *testing.T) {
	o := Object().Field("agree", Bool().Required())
	for _, value := range []string{"true", "TRUE", "True"} {
		input, res := testArgs(o, "agree", value)
		assert.Validation(t, res).FieldsHaveNoErrors("agree")
		assert.True(t, input.Bool("agree"))
	}

	for _, value := range []string{"false", "FALSE", "False"} {
		input, res := testArgs(o, "agree", value)
		assert.Validation(t, res).FieldsHaveNoErrors("agree")
		assert.False(t, input.Bool("agree"))
	}

	input, res := testArgs(o, "agree", "other")
	assert.Validation(t, res).Field("agree", InvalidBoolType())
	_, isBool := input.BoolIf("agree")
	assert.False(t, isBool)
}

func Test_UUID_Required(t *testing.T) {
	f1 := UUID()
	f2 := UUID().Required()
	o := Object().
		Field("id", f1).Field("id_clone", f1).
		Field("parent_id", f2).Field("parent_id_clone", f2).Field("parent_id_clone_not_required", f2.NotRequired())

	_, res := testInput(o)
	assert.Validation(t, res).
		FieldsHaveNoErrors("id", "id_clone", "parent_id_clone_not_required").
		Field("parent_id", Required()).
		Field("parent_id_clone", Required())

	_, res = testInput(o, "parent_id", "FFFFFFFF-FFFF-FFFF-FFFF-FFFFFFFFFFFF", "parent_id_clone", "00000000-0000-0000-0000-000000000000")
	assert.Validation(t, res).
		FieldsHaveNoErrors("parent_id", "id", "parent_id_clone", "id_clone", "parent_id_clone_not_required")
}

func Test_UUID_Type(t *testing.T) {
	o := Object().Field("id", UUID())

	_, res := testInput(o, "id", 3)
	assert.Validation(t, res).
		Field("id", InvalidUUIDType())

	_, res = testInput(o, "id", "Z0000000-0000-0000-0000-00000000000Z")
	assert.Validation(t, res).
		Field("id", InvalidUUIDType())
}

func Test_Nested_Object_Fields(t *testing.T) {
	child := Object().
		Field("age", Int().Required()).
		Field("name", String().Required())

	o1 := Object().Field("user", child)
	_, res := testInput(o1, "user", 3)
	assert.Validation(t, res).
		Field("user", InvalidObjectType())

	o2 := Object().Field("user", child)
	_, res = testInput(o2, "user", map[string]any{})
	assert.Validation(t, res).
		Field("user.age", Required()).
		Field("user.name", Required())

	o3 := Object().Field("entry", o2)
	_, res = testInput(o3, "entry", map[string]any{"user": map[string]any{}})
	assert.Validation(t, res).
		Field("entry.user.age", Required()).
		Field("entry.user.name", Required())

	_, res = testInput(o3, "entry", typed.Typed{"user": typed.Typed{"age": 3000, "name": "Leto"}})
	assert.Validation(t, res).FieldsHaveNoErrors("entry.user.age", "entry.user.name")
}

func Test_Nested_Object_Required(t *testing.T) {
	o1 := Object().Field("user", Object().Required())

	_, res := testInput(o1)
	assert.Validation(t, res).Field("user", Required())

	_, res = testInput(Object().Field("user", o1.NotRequired()))
	assert.Validation(t, res).FieldsHaveNoErrors("user")
}

func Test_Nested_Object_Default(t *testing.T) {
	o1 := Object().Field("user", Object().Default(typed.Typed{"id": 3}))

	data, res := testInput(o1, "user", map[string]any{"name": "leto"})
	assert.Validation(t, res).FieldsHaveNoErrors("user")
	assert.Equal(t, data.Object("user").String("name"), "leto")

	data, res = testInput(o1)
	assert.Validation(t, res).FieldsHaveNoErrors("user")
	assert.Equal(t, data.Object("user").Int("id"), 3)
}

func Test_Object_Func(t *testing.T) {
	userField := BuildField("user.name")
	o1 := Object().Field("user", Object().Func(func(field Field, value typed.Typed, input typed.Typed, res *Result) any {
		res.AddInvalidField(userField, InvalidStringPattern())
		return value
	}))

	_, res := testInput(o1, "user", map[string]any{"name": "leto"})
	assert.Validation(t, res).Field("user.name", InvalidStringPattern())
}

func Test_Array_Objects(t *testing.T) {
	child := Object().Field("name", String().Required())
	o1 := Object().
		Field("users", Array().Required().Validator(child))

	_, res := testInput(o1)
	assert.Validation(t, res).Field("users", Required())

	_, res = testInput(o1, "users", 1)
	assert.Validation(t, res).Field("users", InvalidArrayType())

	_, res = testInput(o1, "users", []any{map[string]any{}, map[string]any{}})
	assert.Validation(t, res).
		Field("users.0.name", Required()).
		Field("users.1.name", Required())

	_, res = testInput(o1, "users", []any{map[string]any{"name": "leto"}})
	assert.Validation(t, res).FieldsHaveNoErrors("users.0.name")

	_, res = testInput(o1, "users", []any{
		map[string]any{"name": "leto"},
		map[string]any{"name": 3},
	})
	assert.Validation(t, res).
		Field("users.1.name", InvalidStringType()).
		FieldsHaveNoErrors("users.0.name")
}

func Test_Array_MinAndMax(t *testing.T) {
	createItem := func() typed.Typed {
		return typed.Typed{"name": "n"}
	}

	child := Object().Field("name", String())
	o1 := Object().Field("users", Array().Min(2).Max(3).Required().Validator(child))

	_, res := testInput(o1, "users", []any{createItem()})
	assert.Validation(t, res).Field("users", InvalidArrayMinLength(2))

	// 4 items, too many
	_, res = testInput(o1, "users", []any{
		createItem(), createItem(), createItem(), createItem(),
	})
	assert.Validation(t, res).Field("users", InvalidArrayMaxLength(3))

	// 2 items, good
	_, res = testInput(o1, "users", []any{
		createItem(), createItem(),
	})
	assert.Validation(t, res).FieldsHaveNoErrors("users")

	// 3 items, good
	_, res = testInput(o1, "users", []any{
		createItem(), createItem(), createItem(),
	})
	assert.Validation(t, res).FieldsHaveNoErrors("users")
}

func Test_Array_Range(t *testing.T) {
	createItem := func() typed.Typed {
		return typed.Typed{"name": "n"}
	}

	child := Object().Field("name", String())
	o1 := Object().Field("users", Array().Range(2, 3).Required().Validator(child))

	_, res := testInput(o1, "users", []any{createItem()})
	assert.Validation(t, res).Field("users", InvalidArrayRangeLength(2, 3))

	// 4 items, too many
	_, res = testInput(o1, "users", []any{
		createItem(), createItem(), createItem(), createItem(),
	})
	assert.Validation(t, res).Field("users", InvalidArrayRangeLength(2, 3))

	// 2 items, good
	_, res = testInput(o1, "users", []any{
		createItem(), createItem(),
	})
	assert.Validation(t, res).FieldsHaveNoErrors("users")

	// 3 items, good
	_, res = testInput(o1, "users", []any{
		createItem(), createItem(), createItem(),
	})
	assert.Validation(t, res).FieldsHaveNoErrors("users")
}

func Test_Array_Strings(t *testing.T) {
	child := String().Length(2, 4)
	o1 := Object().
		Field("users", Array().Required().Validator(child))

	_, res := testInput(o1, "users", []any{"leto", 123, "2"})
	assert.Validation(t, res).
		Field("users.1", InvalidStringType()).
		Field("users.2", InvalidStringLength(2, 4)).
		FieldsHaveNoErrors("users.0")
}

func Test_Array_Bools(t *testing.T) {
	child := Bool()
	o1 := Object().
		Field("users", Array().Required().Validator(child))

	_, res := testInput(o1, "users", []any{true, 123, false})
	assert.Validation(t, res).
		Field("users.1", InvalidBoolType()).
		FieldsHaveNoErrors("users.0", "users.2")
}

func Test_Array_Ints(t *testing.T) {
	child := Int().Min(10)
	o1 := Object().
		Field("users", Array().Required().Validator(child))

	_, res := testInput(o1, "users", []any{10, 9, false, 11.0})
	assert.Validation(t, res).
		Field("users.1", InvalidIntMin(10)).
		Field("users.2", InvalidIntType()).
		FieldsHaveNoErrors("users.0", "users.3")
}

func Test_Array_Floats(t *testing.T) {
	child := Float().Min(10.1)
	o1 := Object().
		Field("users", Array().Required().Validator(child))

	_, res := testInput(o1, "users", []any{10.1, 10.0, false})
	assert.Validation(t, res).
		Field("users.1", InvalidFloatMin(10.1)).
		Field("users.2", InvalidFloatType()).
		FieldsHaveNoErrors("users.0")
}

func Test_Array_Transformer(t *testing.T) {
	o1 := Object().Field("ids", Array().Required().Validator(Int()).Transformer(func(values []any) any {
		ids := make([]int, len(values))
		for i, id := range values {
			ids[i] = id.(int)
		}
		return ids
	}))

	data, res := testInput(o1, "ids", []any{1, 2, 3})
	assert.Validation(t, res).FieldsHaveNoErrors("ids")
	ids := data["ids"].([]int)
	assert.Equal(t, len(ids), 3)
}

func Test_Any_Required(t *testing.T) {
	f2 := Any().Required()
	o := Object().
		Field("name", Any()).
		Field("code", f2).Field("code_not_required", f2.NotRequired())

	_, res := testInput(o)
	assert.Validation(t, res).
		FieldsHaveNoErrors("name", "code_not_required").
		Field("code", Required())

	_, res = testInput(o, "code", 1)
	assert.Validation(t, res).
		FieldsHaveNoErrors("code", "name", "code_not_required")
}

func Test_Any_Default(t *testing.T) {
	o := Object().Field("name", Any().Default(32))

	data, _ := testInput(o)
	assert.Equal(t, data.Int("name"), 32)
}

func Test_Any_Func(t *testing.T) {
	o := Object().Field("name", Any().Func(func(field Field, value any, object typed.Typed, input typed.Typed, res *Result) any {
		assert.Equal(t, value.(string), "one-one")
		return 11
	}))

	data, _ := testInput(o, "name", "one-one")
	assert.Equal(t, data.Int("name"), 11)
}

func testInput(o *ObjectValidator, args ...any) (typed.Typed, *Result) {
	m := make(typed.Typed, len(args)/2)
	for i := 0; i < len(args); i += 2 {
		m[args[i].(string)] = args[i+1]
	}

	res := NewResult(10)
	o.Validate(m, res)
	return m, res
}

func testArgs(o *ObjectValidator, args ...string) (typed.Typed, *Result) {
	m := new(fasthttp.Args)
	for i := 0; i < len(args); i += 2 {
		m.Add(args[i], args[i+1])
	}

	res := NewResult(5)
	input, _ := o.ValidateArgs(m, res)
	return input, res
}
