package validation

import (
	"src.goblgobl.com/utils/concurrent"
)

type Pool[T any] struct {
	*concurrent.Pool[*Context[T]]
}

func NewPool[T any](count uint16, maxErrors uint16) Pool[T] {
	return Pool[T]{
		Pool: concurrent.NewPool[*Context[T]](uint32(count), pooledContextFactory[T](maxErrors)),
	}
}

func pooledContextFactory[T any](maxErrors uint16) func(func(t *Context[T])) *Context[T] {
	return func(release func(ct *Context[T])) *Context[T] {
		ctx := NewContext[T](maxErrors)
		ctx.release = release
		return ctx
	}
}

func (p Pool[T]) Checkout(env T) *Context[T] {
	context := p.Pool.Checkout()
	context.Env = env
	return context
}
