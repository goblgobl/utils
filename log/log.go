package log

import (
	"io"
	"os"
)

var (
	Out io.Writer = os.Stderr

	globalPool *Pool
)

func init() {
	// Under normal cirumstances, we expect the application to setup the pool
	// object early in the startup sequence, but we'll configure a default one
	// incase it's needed before the app has the chane to configure it (e.g. if the
	// app fails to read the configuration file)
	globalPool = NewPool(16, INFO, true, KvFactory(2048), nil)
}

func Checkout() Logger {
	return globalPool.Checkout()
}

func Info(ctx string) Logger {
	return globalPool.Info(ctx)
}

func Warn(ctx string) Logger {
	return globalPool.Warn(ctx)
}

func Error(ctx string) Logger {
	return globalPool.Error(ctx)
}

func Fatal(ctx string) Logger {
	return globalPool.Fatal(ctx)
}

func Request(route string) Logger {
	return globalPool.Request(route)
}

func Detach() Logger {
	return globalPool.Detach()
}

type Logger interface {
	// Actually log the data to the configured output
	// If the logger was not configured for MultiUse, this
	// must release the logger back to the pool (if any).
	Log()

	// Log to the specific writer
	LogTo(io.Writer)

	// Resets the logger, without releasing it back to the pool.
	// Reset must not erased any fixed data
	Reset()

	// Releases the logger back to the pool, if one is configured.
	Release()

	// Gets the log data
	Bytes() []byte

	// Set the level and context for a new log entry
	Info(ctx string) Logger
	Warn(ctx string) Logger
	Error(ctx string) Logger
	Fatal(ctx string) Logger
	Request(route string) Logger

	// Log an error
	Err(err error) Logger

	// Any data already in the logger will be including in _every_
	// message ever generated from this logger. This data survives both
	// a reset and a release.
	// As an example of where this could be used is a project owned loggers
	// in which case these loggers could have the project_id always included
	// in any log entry.
	Fixed()

	// Any data already in the logger will be including in _every_
	// message generated until the logger is released. This data survives a
	// a reset but not a release.
	// This is meant to be used for a request-owned logger where every entry
	// logged with the logger has request_id.
	MultiUse() Logger

	// Add an int value to the current entry
	Int(key string, value int) Logger

	// Add an int64 value to the current entry
	Int64(key string, value int64) Logger

	// Add a string value to the current entry
	String(key string, value string) Logger

	// Add a binary value to the current entry
	Binary(key string, value []byte) Logger

	// Add a boolean value to the current entry
	Bool(key string, value bool) Logger

	// Log a field
	Field(field Field) Logger
}
