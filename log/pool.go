package log

import (
	"src.goblgobl.com/utils/concurrent"
)

type Level uint8

const (
	INFO Level = iota
	WARN
	ERROR
	FATAL
	NONE
)

type Factory func(release func(Logger), level Level, request bool) Logger

type Pool struct {
	*concurrent.Pool[Logger]
	field    *Field
	level    Level
	requests bool
}

func NewPool(count uint16, level Level, requests bool, factory Factory, field *Field) *Pool {
	return &Pool{
		level:    level,
		field:    field,
		requests: requests,
		Pool:     concurrent.NewPool[Logger](uint32(count), pooledLoggerFactory(factory, level, requests, field)),
	}
}

func pooledLoggerFactory(factory Factory, level Level, requests bool, field *Field) func(func(l Logger)) Logger {
	return func(release func(ct Logger)) Logger {
		logger := factory(release, level, requests)
		if field != nil {
			logger.Field(*field).Fixed()
		}
		return logger
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
