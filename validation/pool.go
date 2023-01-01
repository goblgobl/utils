package validation

import "sync/atomic"

type Pool struct {
	maxErrors uint16
	depleted  uint64
	list      chan *Result
}

func NewPool(count uint16, maxErrors uint16) *Pool {
	list := make(chan *Result, count)
	p := &Pool{list: list, maxErrors: maxErrors}
	for i := uint16(0); i < count; i++ {
		result := NewResult(maxErrors)
		result.pool = p
		list <- result
	}
	return p
}

func (p *Pool) Len() int {
	return len(p.list)
}

func (p *Pool) Checkout() *Result {
	select {
	case result := <-p.list:
		return result
	default:
		atomic.AddUint64(&p.depleted, 1)
		return NewResult(p.maxErrors)
	}
}

func (p *Pool) Depleted() uint64 {
	return atomic.SwapUint64(&p.depleted, 0)
}
