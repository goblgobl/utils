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

	p := NewPool(1, INFO, KvFactory(64), nil)
	assertKvLogger(p.Info(""))
	assertKvLogger(p.Warn(""))
	assertKvLogger(p.Error(""))
	assertKvLogger(p.Fatal(""))

	p = NewPool(1, WARN, KvFactory(64), nil)
	assertNoop(p.Info(""))
	assertKvLogger(p.Warn(""))
	assertKvLogger(p.Error(""))
	assertKvLogger(p.Fatal(""))

	p = NewPool(1, ERROR, KvFactory(64), nil)
	assertNoop(p.Info(""))
	assertNoop(p.Warn(""))
	assertKvLogger(p.Error(""))
	assertKvLogger(p.Fatal(""))

	p = NewPool(1, FATAL, KvFactory(64), nil)
	assertNoop(p.Info(""))
	assertNoop(p.Warn(""))
	assertNoop(p.Error(""))
	assertKvLogger(p.Fatal(""))

	p = NewPool(1, NONE, KvFactory(64), nil)
	assertNoop(p.Info(""))
	assertNoop(p.Warn(""))
	assertNoop(p.Error(""))
	assertNoop(p.Fatal(""))
}

func Test_Pool_Checkout(t *testing.T) {
	out := &strings.Builder{}
	p := NewPool(1, WARN, KvFactory(64), nil)

	l1 := p.Checkout()
	l1.Info("should not log").LogTo(out)
	assert.Equal(t, out.String(), "")

	out.Reset()
	l2 := p.Checkout()
	l2.Warn("should log").LogTo(out)
	assert.True(t, strings.Contains(out.String(), "should log"))

	assert.Equal(t, l1.(*KvLogger), l2.(*KvLogger))
}

func Test_Pool_Depleted(t *testing.T) {
	p := NewPool(2, INFO, KvFactory(64), nil)
	assert.Equal(t, p.Len(), 2)
	assert.Equal(t, p.Depleted(), 0)

	l1 := p.Checkout().(*KvLogger)
	assert.Equal(t, p.Len(), 1)
	assert.Equal(t, p.Depleted(), 0)

	l2 := p.Checkout().(*KvLogger)
	assert.Equal(t, p.Len(), 0)
	assert.Equal(t, p.Depleted(), 0)

	l3 := p.Checkout().(*KvLogger)
	assert.Equal(t, p.Len(), 0)
	assert.Equal(t, p.Depleted(), 1)
	assert.Equal(t, p.Depleted(), 0) // calling Delpeted resets it

	assert.NotEqual(t, l1, l2)
	assert.NotEqual(t, l1, l3)
	assert.NotEqual(t, l2, l3)
}

func Test_Pool_DynamicCreationWontReleaseToPool(t *testing.T) {
	p := NewPool(1, INFO, KvFactory(64), nil)

	l1 := p.Checkout().(*KvLogger)
	l2 := p.Checkout().(*KvLogger)
	assert.NotEqual(t, l1, l2)

	l1.Release()
	l2.Release()

	assert.Equal(t, p.Len(), 1)
}

func Test_Pool_KvLogging(t *testing.T) {
	out := &strings.Builder{}
	p := NewPool(1, INFO, KvFactory(128), nil)

	l1 := p.Info("c-info").String("a", "b")
	l1.LogTo(out)
	assertKvLog(t, out, true, map[string]string{
		"a": "b",
		"l": "info",
		"c": "c-info",
	})

	l2 := p.Warn("c-warn").String("a", "b")
	l2.LogTo(out)
	assertKvLog(t, out, true, map[string]string{
		"a": "b",
		"l": "warn",
		"c": "c-warn",
	})

	l3 := p.Error("c-error").String("a", "b")
	l3.LogTo(out)
	assertKvLog(t, out, true, map[string]string{
		"a": "b",
		"l": "error",
		"c": "c-error",
	})

	l4 := p.Fatal("c-fatal").String("a", "b")
	l4.LogTo(out)
	assertKvLog(t, out, true, map[string]string{
		"a": "b",
		"l": "fatal",
		"c": "c-fatal",
	})
}
