package log

import (
	"strings"
	"testing"

	"src.goblgobl.com/tests/assert"
)

func Test_Global_Helpers(t *testing.T) {
	out := &strings.Builder{}
	err := Configure(Config{
		PoolSize: 8,
		Format:   "kv",
		Level:    "INFO",
	})
	assert.Nil(t, err)

	Info("i").LogTo(out)
	assertKvLog(t, out, false, map[string]string{"_l": "info", "_c": "i"})

	Warn("w").LogTo(out)
	assertKvLog(t, out, false, map[string]string{"_l": "warn", "_c": "w"})

	Error("e").LogTo(out)
	assertKvLog(t, out, false, map[string]string{"_l": "error", "_c": "e"})

	Fatal("f").LogTo(out)
	assertKvLog(t, out, false, map[string]string{"_l": "fatal", "_c": "f"})

	Request("r1").LogTo(out)
	assertKvLog(t, out, false, map[string]string{"_l": "req", "_c": "r1"})

	Checkout().Info("i2").LogTo(out)
	assertKvLog(t, out, false, map[string]string{"_l": "info", "_c": "i2"})
}

func Test_Detach(t *testing.T) {
	out := &strings.Builder{}
	err := Configure(Config{
		PoolSize: 8,
		Format:   "kv",
		Level:    "WARN",
	})
	assert.Nil(t, err)

	l1 := Detach()
	l1.Info("i").LogTo(out)
	assertKvLog(t, out, true, nil)

	l1.Warn("w").String("x", "a").LogTo(out)
	assertKvLog(t, out, false, map[string]string{
		"_l": "warn",
		"x":  "a",
	})
}
