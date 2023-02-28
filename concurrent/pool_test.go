package concurrent

import (
	"sync/atomic"
	"testing"
	"time"

	"src.goblgobl.com/tests/assert"
)

func Test_Pool(t *testing.T) {
	p := NewPool(8, poolItemFactory())
	seen := make(map[int]*PoolItem, 8)
	for i := 0; i < 8; i++ {
		t := p.Checkout()
		seen[t.id] = t
		t.Release()
	}
	assert.Equal(t, len(seen), 8)

	for i := 0; i < 8; i++ {
		t := p.Checkout()
		seen[t.id] = t
		t.Release()
	}
	assert.Equal(t, len(seen), 8)
}

func Test_Pool_Depleted(t *testing.T) {
	p := NewPool(8, poolItemFactory())
	seen := make(map[int]*PoolItem, 8)
	for i := 0; i < 8; i++ {
		t := p.Checkout()
		seen[t.id] = t
	}
	assert.Equal(t, p.Depleted(), 0)

	p.Checkout().Release()
	assert.Equal(t, p.Depleted(), 1)

	p.Checkout().Release()
	assert.Equal(t, p.Depleted(), 2)

	seen[4].Release()
	p.Checkout().Release()
	assert.Equal(t, p.Depleted(), 2)
}

// not really sure this is testing anything meaningful
func Test_Pool_Concurrency(t *testing.T) {
	p := NewPool(100, poolItemFactory())
	for i := 0; i < 100; i++ {
		go func() {
			t := p.Checkout()
			time.Sleep(time.Millisecond * 10)
			t.Release()
		}()
	}
}

type PoolItem struct {
	release func(*PoolItem)
	id      int
}

func (t *PoolItem) Release() {
	t.release(t)
}

func poolItemFactory() func(func(t *PoolItem)) *PoolItem {
	var testId atomic.Uint64
	return func(release func(t *PoolItem)) *PoolItem {
		return &PoolItem{
			release: release,
			id:      int(testId.Add(1)),
		}
	}
}
