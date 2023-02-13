package buffer

import (
	"sync"
	"sync/atomic"
)

type Pool struct {
	minSize  uint32
	maxSize  uint32
	depleted uint64
	list     chan *Buffer
}

func NewPoolFromConfig(config Config) *Pool {
	return NewPool(config.Count, config.Min, config.Max)
}

func NewPool(count uint16, minSize uint32, maxSize uint32) *Pool {
	list := make(chan *Buffer, count)
	p := &Pool{list: list, minSize: minSize, maxSize: maxSize}
	for i := uint16(0); i < count; i++ {
		buffer := New(minSize, maxSize)
		buffer.pool = p
		list <- buffer
	}
	return p
}

func (p *Pool) Len() int {
	return len(p.list)
}

func (p *Pool) Checkout() *Buffer {
	return p.CheckoutMax(p.maxSize)
}

func (p *Pool) CheckoutMax(maxSize uint32) *Buffer {
	select {
	case buffer := <-p.list:
		buffer.max = int(maxSize)
		return buffer
	default:
		atomic.AddUint64(&p.depleted, 1)
		return New(p.minSize, maxSize)
	}
}

func (p *Pool) Depleted() uint64 {
	return atomic.LoadUint64(&p.depleted)
}

/*
A lot of our object pools are encapsulated inside of project
environments. This has a lot of benefits. It simplifies the
code, minimizes contention, and further isolates projects.

But for buffers, which are used both to generate SQL and
read results, having these per-env would be memory-
inefficient. Our buffers need to be relatively large (for
responses), so sharing a large pool across projects is likely
to result in much better usage (so we need far less of them).

Also, these objects have a lifecycle that's different than
an env. Specifically, it has to interact with the http framework
in a way that minimize the amount of copying we need to do.
*/

var buffers = sync.Pool{
	New: func() any {
		// max size doesn't really matter, we're going to
		// reset it to a per-project value on checkout.
		return New(65536, 65536)
	},
}

func Checkout(maxSize int) *Buffer {
	b := buffers.Get().(*Buffer)
	b.max = maxSize
	return b
}

func Release(b *Buffer) {

}
