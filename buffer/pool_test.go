package buffer

import (
	"testing"

	"src.goblgobl.com/tests/assert"
)

func Test_Pool_Checkout_and_Release(t *testing.T) {
	p := NewPool(2, 10, 20)
	b := p.Checkout(22)
	b.Write([]byte("abc"))
	assert.Equal(t, b.max, 22)
	assert.Equal(t, len(b.data), 10)

	b.Release()

	s, err := b.String()
	assert.Nil(t, err)
	assert.Equal(t, s, "")
	assert.Equal(t, p.Len(), 2)
}

func Test_Pool_FromConfig(t *testing.T) {
	p := NewPoolFromConfig(Config{Count: 3, Min: 5, Max: 30})
	b := p.Checkout(30)
	b.Write([]byte("abc"))
	assert.Equal(t, b.max, 30)
	assert.Equal(t, len(b.data), 5)

	b.Release()

	s, err := b.String()
	assert.Nil(t, err)
	assert.Equal(t, s, "")
	assert.Equal(t, p.Len(), 3)
}
