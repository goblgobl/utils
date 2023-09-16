package buffer

import (
	"src.goblgobl.com/utils/concurrent"
)

type Config struct {
	Count uint16 `json:"count"`
	Min   uint32 `json:"min"`
	Max   uint32 `json:"max"`
}

type Pool struct {
	*concurrent.Pool[*Buffer]
	maxSize uint32
}

func NewPoolFromConfig(config Config) Pool {
	return NewPool(config.Count, config.Min, config.Max)
}

func NewPool(count uint16, minSize uint32, maxSize uint32) Pool {
	return Pool{
		maxSize: maxSize,
		Pool:    concurrent.NewPool[*Buffer](uint32(count), pooledBufferFactory(minSize, maxSize)),
	}
}

func pooledBufferFactory(minSize uint32, maxSize uint32) func(func(b *Buffer)) *Buffer {
	return func(release func(b *Buffer)) *Buffer {
		buffer := New(minSize, maxSize)
		buffer.release = release
		return buffer
	}
}

func (p Pool) Checkout() *Buffer {
	return p.CheckoutMax(p.maxSize)
}

func (p *Pool) CheckoutMax(maxSize uint32) *Buffer {
	buffer := p.Pool.Checkout()
	buffer.max = int(maxSize)
	return buffer
}
