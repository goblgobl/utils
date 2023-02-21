package validation

import (
	"testing"
)

func Test_UUID_Required(t *testing.T) {
	f1 := UUID[E]()
	f2 := UUID[E]().Required()
	o := Object[E]().
		Field("id", f1).Field("id_clone", f1.Clone()).
		Field("parent_id", f2).Field("parent_id_clone", f2.Clone()).Field("parent_id_clone_not_required", f2.Clone().NotRequired())

	testValidator(t, o).
		FieldsHaveNoErrors("id", "id_clone", "parent_id_clone_not_required").
		Field("parent_id", Required).
		Field("parent_id_clone", Required)

	testValidator(t, o, "parent_id", "FFFFFFFF-FFFF-FFFF-FFFF-FFFFFFFFFFFF", "parent_id_clone", "00000000-0000-0000-0000-000000000000").
		FieldsHaveNoErrors("parent_id", "id", "parent_id_clone", "id_clone", "parent_id_clone_not_required")
}

func Test_UUID_Type(t *testing.T) {
	o := Object[E]().Field("id", UUID[E]())

	testValidator(t, o, "id", 3).
		Field("id", TypeUUID)

	testValidator(t, o, "id", "Z0000000-0000-0000-0000-00000000000Z").
		Field("id", TypeUUID)
}
