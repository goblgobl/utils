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
