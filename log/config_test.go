package log

import (
	"testing"

	"src.goblgobl.com/tests/assert"
)

func Test_Configure_InvalidLevel(t *testing.T) {
	err := Configure(Config{Level: "invalid"})
	assert.Equal(t, err.Error(), "code: 3001 - log.level is invalid. Should be one of: INFO, WARN, ERROR, FATAL or NONE")
}

func Test_Configure_InvalidFormat(t *testing.T) {
	err := Configure(Config{Format: "unknown"})
	assert.Equal(t, err.Error(), "code: 3002 - log.format is invalid. Should be one of: kv")
}

func Test_Configure_Defaults(t *testing.T) {
	err := Configure(Config{})
	assert.Nil(t, err)
	// pool is divided into 8 buckets.
	// 100/8 == 12.5
	// int(12.5) == 12
	// 12 * 8 == 96
	assert.Equal(t, globalPool.Len(), 96)

	l := globalPool.Checkout().(*KvLogger)
	defer l.Release()
	assert.Equal(t, l.buffer.Max(), 4096)
}

func Test_Configure_Custom(t *testing.T) {
	err := Configure(Config{
		PoolSize: 24,
		Format:   "kv",
		Level:    "error",
		KV:       KvConfig{MaxSize: 100},
	})

	assert.Nil(t, err)
	assert.Equal(t, globalPool.level, ERROR)
	assert.Equal(t, globalPool.Len(), 24)

	l := globalPool.Checkout().(*KvLogger)
	defer l.Release()
	assert.Equal(t, l.buffer.Max(), 100)

	levels := map[string]Level{
		"infO":  INFO,
		"WARN":  WARN,
		"ErrOR": ERROR,
		"faTAL": FATAL,
		"none":  NONE,
	}

	for level, typed := range levels {
		err := Configure(Config{Level: level})
		assert.Nil(t, err)
		assert.Equal(t, globalPool.level, typed)
	}
}
