package buffer

import (
	"testing"

	"src.goblgobl.com/tests/assert"
)

func Test_Pool_CheckoutMax_and_Release(t *testing.T) {
	p := NewPool(2, 10, 20)

	b := p.CheckoutMax(22)
	b.Write([]byte("abc"))
	assert.Equal(t, b.max, 22)
	assert.Equal(t, len(b.data), 10)
	b.Release()

	s, err := b.String()
	assert.Nil(t, err)
	assert.Equal(t, s, "")
	assert.Equal(t, p.Len(), 2)

	b = p.Checkout()
	b.Write([]byte("abcd"))
	assert.Equal(t, b.max, 20)
	assert.Equal(t, len(b.data), 10)
	b.Release()
}

func Test_Pool_FromConfig(t *testing.T) {
	p := NewPoolFromConfig(Config{Count: 3, Min: 5, Max: 30})
	b := p.Checkout()
	b.Write([]byte("abc"))
	assert.Equal(t, b.max, 30)
	assert.Equal(t, len(b.data), 5)

	b.Release()

	s, err := b.String()
	assert.Nil(t, err)
	assert.Equal(t, s, "")
	assert.Equal(t, p.Len(), 3)
}

func Test_Pool_Depleted(t *testing.T) {
	p := NewPoolFromConfig(Config{Count: 2, Min: 5, Max: 30})

	b1 := p.Checkout()
	assert.Equal(t, p.Depleted(), 0)

	p.Checkout()
	assert.Equal(t, p.Depleted(), 0)

	for i := uint64(0); i < 10; i++ {
		p.Checkout()
		assert.Equal(t, p.Depleted(), i+1)
	}

	b1.Release()
	max := p.Depleted()
	p.Checkout() // grabs b1, so depleted should not increment
	assert.Equal(t, p.Depleted(), max)
}
