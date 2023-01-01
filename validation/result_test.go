package validation

import (
	"testing"

	"src.goblgobl.com/tests/assert"
)

func Test_Result_FieldError_NoData(t *testing.T) {
	r := NewResult(4)
	assert.True(t, r.IsValid())

	f := Field{}
	r.AddInvalidField(f.add("field1"), Required())
	assert.False(t, r.IsValid())

	invalid := r.Errors()[0].(InvalidField)
	assert.Nil(t, invalid.Data)
	assert.Equal(t, invalid.Field, "field1")
	assert.Equal(t, invalid.Code, 1001)
	assert.Equal(t, invalid.Error, "required")
}

func Test_Result_FieldError_Data(t *testing.T) {
	r := NewResult(4)

	f := Field{}
	r.AddInvalidField(f.add("field2"), InvalidIntMin(33))
	invalid := r.Errors()[0].(InvalidField)
	assert.Equal(t, invalid.Data.(DataMin).Min.(int), 33)
	assert.Equal(t, invalid.Code, 1006)
	assert.Equal(t, invalid.Field, "field2")
	assert.Equal(t, invalid.Error, "must be greater or equal to 33")
}

func Test_Result_NoField_NoData(t *testing.T) {
	r := NewResult(4)
	assert.True(t, r.IsValid())

	r.AddInvalid(Required())
	assert.False(t, r.IsValid())

	invalid := r.Errors()[0].(Invalid)
	assert.Nil(t, invalid.Data)
	assert.Equal(t, invalid.Code, 1001)
	assert.Equal(t, invalid.Error, "required")
}

func Test_Result_NoField_Data(t *testing.T) {
	r := NewResult(4)
	r.AddInvalid(InvalidIntMin(33))
	invalid := r.Errors()[0].(Invalid)
	assert.Equal(t, invalid.Data.(DataMin).Min.(int), 33)
	assert.Equal(t, invalid.Code, 1006)
	assert.Equal(t, invalid.Error, "must be greater or equal to 33")
}
