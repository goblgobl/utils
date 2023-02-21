package validation

import "sync/atomic"

type Pool[T any] struct {
  maxErrors uint16
  depleted  uint64
  list      chan *Context[T]
}

func NewPool[T any](count uint16, maxErrors uint16) *Pool[T] {
  list := make(chan *Context[T], count)
  p := &Pool[T]{list: list, maxErrors: maxErrors}
  for i := uint16(0); i < count; i++ {
    ctx := NewContext[T](maxErrors)
    ctx.pool = p
    list <- ctx
  }
  return p
}

func (p *Pool[T]) Len() int {
  return len(p.list)
}

func (p *Pool[T]) Checkout(env T) *Context[T] {
  select {
  case ctx := <-p.list:
    ctx.Env = env
    return ctx
  default:
    atomic.AddUint64(&p.depleted, 1)
    ctx := NewContext[T](p.maxErrors)
    ctx.Env = env
    return ctx
  }
}

func (p *Pool[T]) Depleted() uint64 {
  return atomic.LoadUint64(&p.depleted)
}
