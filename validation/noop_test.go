package validation

import "testing"

func Test_Noop_Required(t *testing.T) {
	o := Object[E]().
		Field("f1", Noop[E]())

	testValidator(t, o).FieldsHaveNoErrors("f1")
	testValidator(t, o, "f1", 1).FieldsHaveNoErrors("f1")
	testValidator(t, o, "f1", "fx").FieldsHaveNoErrors("f1")
}
