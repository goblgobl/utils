package validation

import (
	"testing"

	"src.goblgobl.com/tests/assert"
)

func Test_Pool_Depleted(t *testing.T) {
	p := NewPool[int](2, 1)
	assert.Equal(t, p.Len(), 2)
	assert.Equal(t, p.Depleted(), 0)

	l1 := p.Checkout(1)
	assert.Equal(t, l1.Env, 1)
	assert.Equal(t, p.Len(), 1)
	assert.Equal(t, p.Depleted(), 0)

	l2 := p.Checkout(2)
	assert.Equal(t, l2.Env, 2)
	assert.Equal(t, p.Len(), 0)
	assert.Equal(t, p.Depleted(), 0)

	l3 := p.Checkout(3)
	assert.Equal(t, l3.Env, 3)
	assert.Equal(t, p.Len(), 0)
	assert.Equal(t, p.Depleted(), 1)

	l4 := p.Checkout(4)
	assert.Equal(t, l4.Env, 4)
	assert.Equal(t, p.Len(), 0)
	assert.Equal(t, p.Depleted(), 2)

	assert.NotEqual(t, l1, l2)
	assert.NotEqual(t, l1, l3)
	assert.NotEqual(t, l2, l3)
	assert.NotEqual(t, l2, l4)
	assert.NotEqual(t, l3, l4)
}

func Test_Pool_DynamicCreationWontReleaseToPool(t *testing.T) {
	p := NewPool[any](1, 3)

	l1 := p.Checkout(nil)
	l2 := p.Checkout(nil)
	assert.NotEqual(t, l1, l2)

	l1.Release()
	l2.Release()

	assert.Equal(t, p.Len(), 1)
}
