package log

import (
	"strings"

	"src.goblgobl.com/utils"
)

type Config struct {
	Requests *bool    `json:"requests"`
	Level    string   `json:"level"`
	Format   string   `json:"format"`
	PoolSize uint16   `json:"pool_size"`
	KV       KvConfig `json:"kv"`
}

type KvConfig struct {
	MaxSize uint32 `json:"max_size"`
}

func Configure(config Config) error {
	var level Level
	levelName := strings.ToUpper(config.Level)

	switch levelName {
	case "INFO":
		level = INFO
	case "", "WARN":
		level = WARN
		levelName = "WARN" // reset this incase it was empty/default
	case "ERROR":
		level = ERROR
	case "FATAL":
		level = FATAL
	case "NONE":
		level = NONE
	default:
		return Errf(utils.ERR_INVALID_LOG_LEVEL, "log.level is invalid. Should be one of: INFO, WARN, ERROR, FATAL or NONE")
	}

	var factory Factory
	formatName := strings.ToUpper(config.Format)
	switch formatName {
	case "", "KV":
		maxSize := config.KV.MaxSize
		if maxSize == 0 {
			maxSize = 131072 // 128KB
		}
		factory = KvFactory(maxSize)
		formatName = "KV" // reset this incase it was empty
	default:
		return Errf(utils.ERR_INVALID_LOG_FORMAT, "log.format is invalid. Should be one of: kv")
	}

	poolSize := config.PoolSize
	if poolSize == 0 {
		poolSize = 100
	}

	requests := true
	if r := config.Requests; r != nil && *r == false {
		requests = false
	}

	globalPool = NewPool(poolSize, level, requests, factory, nil)
	Info("log_config").
		String("level", levelName).
		String("format", formatName).
		Int("pool_size", int(poolSize)).
		Bool("requests", requests).
		Log()
	return nil
}
