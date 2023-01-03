package log

import (
	"strings"
	"testing"

	"src.goblgobl.com/tests/assert"
)

func Test_Global_Helpers(t *testing.T) {
	out := &strings.Builder{}
	err := Configure(Config{
		PoolSize: 2,
		Format:   "kv",
		Level:    "INFO",
	})
	assert.Nil(t, err)

	Info("i").LogTo(out)
	assertKvLog(t, out, false, map[string]string{"l": "info", "c": "i"})

	Warn("w").LogTo(out)
	assertKvLog(t, out, false, map[string]string{"l": "warn", "c": "w"})

	Error("e").LogTo(out)
	assertKvLog(t, out, false, map[string]string{"l": "error", "c": "e"})

	Fatal("f").LogTo(out)
	assertKvLog(t, out, false, map[string]string{"l": "fatal", "c": "f"})

	Request("r1").LogTo(out)
	assertKvLog(t, out, false, map[string]string{"l": "req", "c": "r1"})

	Checkout().Info("i2").LogTo(out)
	assertKvLog(t, out, false, map[string]string{"l": "info", "c": "i2"})
}
