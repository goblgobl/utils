package validation

import (
	"testing"

	"src.goblgobl.com/tests/assert"
)

func Test_Context_InvalidField(t *testing.T) {
	ctx := NewContext[int](4)
	assert.True(t, ctx.IsValid())

	ctx.Field = &Field{Flat: "f1"}
	ctx.InvalidField(Required)
	assert.False(t, ctx.IsValid())

	invalid := ctx.Errors()[0].(InvalidField)
	assert.Nil(t, invalid.Data)
	assert.Equal(t, invalid.Field, "f1")
	assert.Equal(t, invalid.Code, 1001)
	assert.Equal(t, invalid.Error, "required")
}

func Test_Context_InvalidWithField(t *testing.T) {
	ctx := NewContext[int](4)
	assert.True(t, ctx.IsValid())

	ctx.Field = &Field{Flat: "f1"}
	ctx.InvalidWithField(Required, &Field{Flat: "nf2"})
	assert.False(t, ctx.IsValid())

	invalid := ctx.Errors()[0].(InvalidField)
	assert.Nil(t, invalid.Data)
	assert.Equal(t, invalid.Field, "nf2")
	assert.Equal(t, invalid.Code, 1001)
	assert.Equal(t, invalid.Error, "required")
}

func Test_Context_Array(t *testing.T) {
	ctx := NewContext[int](4)

	ctx.StartArray()
	ctx.ArrayIndex(1)
	ctx.Field = &Field{Path: []string{"user", "", "name"}}
	ctx.InvalidField(Required)
	invalid := ctx.Errors()[0].(InvalidField)
	assert.Equal(t, invalid.Field, "user.1.name")

	ctx.StartArray()
	ctx.ArrayIndex(2)
	ctx.Field = &Field{Path: []string{"user", "", "tags", ""}}
	ctx.InvalidField(Required)
	invalid = ctx.Errors()[1].(InvalidField)
	assert.Equal(t, invalid.Field, "user.1.tags.2")

	ctx.EndArray()
	ctx.ArrayIndex(3)
	ctx.Field = &Field{Path: []string{"user", "", "name"}}
	ctx.InvalidField(Required)
	invalid = ctx.Errors()[2].(InvalidField)
	assert.Equal(t, invalid.Field, "user.3.name")
}

func Test_Context_Suspend_And_Resume_Array(t *testing.T) {
	ctx := NewContext[int](4)

	ctx.StartArray()
	ctx.ArrayIndex(5)

	d := ctx.SuspendArray()
	ctx.Field = &Field{Flat: "x"}
	ctx.InvalidField(Required)
	invalid := ctx.Errors()[0].(InvalidField)
	assert.Equal(t, invalid.Field, "x")

	ctx.ResumeArray(d)
	ctx.Field = &Field{Path: []string{"user", "", "name"}}
	ctx.InvalidField(Required)
	invalid = ctx.Errors()[1].(InvalidField)
	assert.Equal(t, invalid.Field, "user.5.name")
}

func Test_Context_Release(t *testing.T) {
	p := NewPool[any](1, 3)
	ctx := p.Checkout(nil)

	ctx.InvalidWithField(Required, &Field{Flat: "nf2"})
	assert.False(t, ctx.IsValid())
	assert.Equal(t, len(ctx.Errors()), 1)

	ctx.Release()

	ctx = p.Checkout(nil)
	assert.True(t, ctx.IsValid())
	assert.Equal(t, len(ctx.Errors()), 0)
}
