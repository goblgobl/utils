package log

import (
	"strings"
	"testing"

	"src.goblgobl.com/tests/assert"
)

func Test_Pool_Level(t *testing.T) {
	assertNoop := func(l Logger) {
		_, ok := l.(Noop)
		assert.True(t, ok)
		l.Release()
	}

	assertKvLogger := func(l Logger) {
		_, ok := l.(*KvLogger)
		assert.True(t, ok)
		l.Release()
	}

	p := NewPool(8, INFO, true, KvFactory(64), nil)
	assertKvLogger(p.Info(""))
	assertKvLogger(p.Warn(""))
	assertKvLogger(p.Error(""))
	assertKvLogger(p.Fatal(""))

	p = NewPool(8, WARN, true, KvFactory(64), nil)
	assertNoop(p.Info(""))
	assertKvLogger(p.Warn(""))
	assertKvLogger(p.Error(""))
	assertKvLogger(p.Fatal(""))

	p = NewPool(8, ERROR, true, KvFactory(64), nil)
	assertNoop(p.Info(""))
	assertNoop(p.Warn(""))
	assertKvLogger(p.Error(""))
	assertKvLogger(p.Fatal(""))

	p = NewPool(8, FATAL, true, KvFactory(64), nil)
	assertNoop(p.Info(""))
	assertNoop(p.Warn(""))
	assertNoop(p.Error(""))
	assertKvLogger(p.Fatal(""))

	p = NewPool(8, NONE, true, KvFactory(64), nil)
	assertNoop(p.Info(""))
	assertNoop(p.Warn(""))
	assertNoop(p.Error(""))
	assertNoop(p.Fatal(""))
}

func Test_Pool_KvLogging(t *testing.T) {
	out := &strings.Builder{}
	p := NewPool(8, INFO, true, KvFactory(128), nil)

	l1 := p.Info("c-info").String("a", "b")
	l1.LogTo(out)
	assertKvLog(t, out, true, map[string]string{
		"a":  "b",
		"_l": "info",
		"_c": "c-info",
	})

	l2 := p.Warn("c-warn").String("a", "b")
	l2.LogTo(out)
	assertKvLog(t, out, true, map[string]string{
		"a":  "b",
		"_l": "warn",
		"_c": "c-warn",
	})

	l3 := p.Error("c-error").String("a", "b")
	l3.LogTo(out)
	assertKvLog(t, out, true, map[string]string{
		"a":  "b",
		"_l": "error",
		"_c": "c-error",
	})

	l4 := p.Fatal("c-fatal").String("a", "b")
	l4.LogTo(out)
	assertKvLog(t, out, true, map[string]string{
		"a":  "b",
		"_l": "fatal",
		"_c": "c-fatal",
	})
}

func Test_Pool_Request(t *testing.T) {
	out := &strings.Builder{}
	p := NewPool(8, FATAL, true, KvFactory(128), nil)

	l1 := p.Request("route1").String("a", "b")
	l1.LogTo(out)
	assertKvLog(t, out, true, map[string]string{
		"a":  "b",
		"_l": "req",
		"_c": "route1",
	})

	// disable request logging
	p = NewPool(8, FATAL, false, KvFactory(128), nil)
	l1 = p.Request("route2").String("a", "b")
	l1.LogTo(out)
	assertKvLog(t, out, true, nil)

}
