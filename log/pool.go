package log

import (
	"sync/atomic"
)

type Level uint8

const (
	INFO Level = iota
	WARN
	ERROR
	FATAL
	NONE
)

type Factory func(p *Pool, level Level, request bool) Logger

type Pool struct {
	field    *Field
	factory  Factory
	list     chan Logger
	depleted uint64
	level    Level
	requests bool
}

func NewPool(count uint16, level Level, requests bool, factory Factory, field *Field) *Pool {
	list := make(chan Logger, count)
	p := &Pool{
		list:     list,
		level:    level,
		field:    field,
		requests: requests,
		factory:  factory,
	}

	for i := uint16(0); i < count; i++ {
		l := factory(p, level, requests)
		if field != nil {
			l.Field(*field).Fixed()
		}
		list <- l
	}
	return p
}

func (p *Pool) Len() int {
	return len(p.list)
}

func (p *Pool) Checkout() Logger {
	select {
	case logger := <-p.list:
		return logger
	default:
		atomic.AddUint64(&p.depleted, 1)
		l := p.factory(nil, p.level, p.requests)
		if field := p.field; field != nil {
			l.Field(*field).Fixed()
		}
		return l
	}
}

func (p *Pool) Info(ctx string) Logger {
	if p.level > INFO {
		return Noop{}
	}
	return p.Checkout().Info(ctx)
}

func (p *Pool) Warn(ctx string) Logger {
	if p.level > WARN {
		return Noop{}
	}
	return p.Checkout().Warn(ctx)
}

func (p *Pool) Error(ctx string) Logger {
	if p.level > ERROR {
		return Noop{}
	}
	return p.Checkout().Error(ctx)
}

func (p *Pool) Fatal(ctx string) Logger {
	if p.level > FATAL {
		return Noop{}
	}
	return p.Checkout().Fatal(ctx)
}

func (p *Pool) Request(route string) Logger {
	if !p.requests {
		return Noop{}
	}
	return p.Checkout().Request(route)
}

func (p *Pool) Depleted() uint64 {
	return atomic.SwapUint64(&p.depleted, 0)
}
