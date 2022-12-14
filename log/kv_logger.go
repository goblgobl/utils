package log

/*
Loggers generally come from a pool in one of three ways:
1 - Call pool.Info/Warn/Error/Fatal and you'll get a logger from it.
    In these cases, the pool will check the log level and
    return a Noop logger if the level of the logged message < the configured log level

2 - Call to pool.Checkout(). This is used with the MultiUse feature, where
    we'll want to possibly write multiple message all including the same fields
    (like a requestId). In such cases, the log level cannot be checked by the pool
    (because the returned logger can be used multiple times to log different message levels).
    Therefore the logger itself must check the level.

In case #2, we only check the log level once per call (yay!). In case #2
we check the log level twice: once in the pool and once again in the logger.
We need to do the logger check to cover case #2. The first check _could_ be skipped
without any change in behavior, but it would put higher contention and usage on our
pool. That seems worth avoiding since it's reasonable to think most messages will be
INFO and our configured level will be > INFO.

3 - The 3rd way a logger come from the pool is when the pool is depleted and
    we create a new logger just-in-time. This type of logger is only different
    in that, on release, it isn't put back into the pool.
*/

import (
	"io"
	"strconv"
	"strings"
	"time"
)

type KvLogger struct {
	// the position in buffer to write to next
	pos uint64

	// reference back into our pool
	pool *Pool

	// buffer that we write our message to
	buffer []byte

	// A logger can have a fixed piece of data which is
	// always included (e.g pid=$PROJECT_ID for a project-owned
	// logger). Once our fixed data is set, pos will never be
	// less than fixedLen.
	fixedLen uint64

	// A logger can also have temporary repeated data
	// (e.g. rid=$REQUEST_ID for an env-owned logger).
	// After logging a message, pos == multiUseLen. Only
	// on reset/release will pos == fixedLen
	multiUseLen uint64

	// The log level that we're logging.
	level Level

	// Whether or not we're logging request messages
	requests bool
}

func NewKvLogger(maxSize uint32, pool *Pool, level Level, requests bool) *KvLogger {
	return &KvLogger{
		pool:     pool,
		level:    level,
		requests: requests,
		buffer:   make([]byte, maxSize),
	}
}

func KvFactory(maxSize uint32) Factory {
	return func(pool *Pool, level Level, requests bool) Logger {
		return NewKvLogger(maxSize, pool, level, requests)
	}
}

// Get the bytes from the logger. This is only valid before Log is called (after
// log is called, you'll get an empty slice). Only really useful for testing.
func (l *KvLogger) Bytes() []byte {
	return l.buffer[:l.pos]
}

// Logger will _always_ include this data. Meant to be used with the Field builder.
// Even once released to the pool and re-checked out, this data will still be in the logger.
// For checkout-specific data, see MultiUse().
func (l *KvLogger) Fixed() {
	l.fixedLen = l.pos
}

// Similar to Fixed, but exists only while checked out
func (l *KvLogger) MultiUse() Logger {
	l.multiUseLen = l.pos
	return l
}

// Add a field (key=value) where value is a string
func (l *KvLogger) String(key string, value string) Logger {
	l.writeKeyValue(key, value, false)
	return l
}

// Add a field (key=value) where value is an int
func (l *KvLogger) Int(key string, value int) Logger {
	return l.Int64(key, int64(value))
}

// Add a field (key=value) where value is an int
func (l *KvLogger) Int64(key string, value int64) Logger {
	l.writeKeyValue(key, strconv.FormatInt(value, 10), true)
	return l
}

// Add a field (key=value) where value is a boolean
func (l *KvLogger) Bool(key string, value bool) Logger {
	if value {
		l.writeKeyValue(key, "1", true)
	} else {
		l.writeKeyValue(key, "0", true)
	}
	return l
}

// Add a field (key=value) where value is an error
func (l *KvLogger) Err(err error) Logger {
	se, ok := err.(*StructuredError)
	if !ok {
		return l.String("err", err.Error())
	}

	l.Int("code", se.Code).String("err", se.Err.Error())
	for key, value := range se.Data {
		switch v := value.(type) {
		case string:
			l.String(key, v)
		case int:
			l.Int(key, v)
		}
	}
	return l
}

// Write the log to our globally configured writer
func (l *KvLogger) Log() {
	l.LogTo(Out)
}

func (l *KvLogger) LogTo(out io.Writer) {
	pos := l.pos
	buffer := l.buffer

	// no length check, if we did everything right, there should
	// always be at least 1 space in our buffer
	buffer[pos] = '\n'
	out.Write(buffer[:pos+1])
	l.conditionalRelease()
}

func (l *KvLogger) Reset() {
	l.pos = l.fixedLen
}

func (l *KvLogger) Release() {
	l.pos = l.fixedLen // Reset()
	if pool := l.pool; pool != nil {
		pool.list <- l
	}
}

// Normally, logger is automatically released when Log or LogTo is called
// unless we've enabled multiUse.
func (l *KvLogger) conditionalRelease() {
	if l.multiUseLen == 0 {
		l.Release()
	}
}

// Log an info-level message.
func (l *KvLogger) Info(ctx string) Logger {
	if l.level > INFO {
		l.conditionalRelease()
		return Noop{}
	}
	return l.start(ctx, []byte("l=info t="))
}

// Log an warn-level message.
func (l *KvLogger) Warn(ctx string) Logger {
	if l.level > WARN {
		l.conditionalRelease()
		return Noop{}
	}
	return l.start(ctx, []byte("l=warn t="))
}

// Log an error-level message.
func (l *KvLogger) Error(ctx string) Logger {
	if l.level > ERROR {
		l.conditionalRelease()
		return Noop{}
	}
	return l.start(ctx, []byte("l=error t="))
}

// Log an fatal-level message.
func (l *KvLogger) Fatal(ctx string) Logger {
	if l.level > FATAL {
		l.conditionalRelease()
		return Noop{}
	}
	return l.start(ctx, []byte("l=fatal t="))
}

// Log a request message.
func (l *KvLogger) Request(route string) Logger {
	if !l.requests {
		l.conditionalRelease()
		return Noop{}
	}
	return l.start(route, []byte("l=req t="))
}

func (l *KvLogger) Field(field Field) Logger {
	pos := l.pos
	buffer := l.buffer
	bl := uint64(len(buffer))

	// might already have data
	if pos != 0 && pos < bl {
		buffer[pos] = ' '
		pos += 1
	}

	if pos < uint64(len(buffer)) {
		data := field.kv
		copy(buffer[pos:], data)
		l.pos = pos + uint64(len(data))
	}
	return l
}

// "starts" a new log message. Every message always contains a timestamp (t) a
// context (c) and a level (l).
func (l *KvLogger) start(ctx string, meta []byte) Logger {
	pos := l.pos
	buffer := l.buffer

	bl := uint64(len(buffer))

	// pos > 0 when MultiUse is enabled
	if pos > 0 && pos < bl {
		buffer[pos] = ' '
		pos = pos + 1
	}

	copy(buffer[pos:], meta)
	pos += uint64(len(meta))

	t := strconv.FormatInt(time.Now().Unix(), 10)
	copy(buffer[pos:], t)
	pos += uint64(len(t))

	// we always expect the ctx to be safe and to outlive this log
	copy(buffer[pos:], []byte(" c="))
	pos += 3

	copy(buffer[pos:], ctx)
	pos += uint64(len(ctx))

	l.pos = pos
	return l
}

// When safe, we're being told that value 100% does not need
// to be escaped (e.g. we know the value is an int), so we don't need to check/encode it.
func (l *KvLogger) writeKeyValue(key string, value string, safe bool) {
	l.pos = writeKeyValue(key, value, safe, l.pos, l.buffer)
}

// We expect key to always be safe to write as-is.
// We only encode newline and quotes. If either is present, the value is quote encoded.
func writeKeyValue(key string, value string, safe bool, pos uint64, buffer []byte) uint64 {
	bl := uint64(len(buffer))

	// Need at least enough room for:
	// space sperator + equal separator + trailing newline
	// + our key + our value
	if bl-pos < uint64(len(key)+len(value))+3 {
		return pos
	}

	if pos > 0 {
		buffer[pos] = ' '
		pos += 1
	}

	copy(buffer[pos:], key)
	pos += uint64(len(key))

	buffer[pos] = '='
	pos += 1

	if safe || !strings.ContainsAny(value, " =\"\n") {
		copy(buffer[pos:], value)
		return pos + uint64(len(value))
	}

	buffer[pos] = '"'
	pos += 1

	// -2 because we need enough space for our quote and final newline
	var i int
	for ; i < len(value) && pos < bl-5; i++ {
		c := value[i]
		switch c {
		case '\n':
			buffer[pos] = '\\'
			buffer[pos+1] = 'n'
			pos += 2
		case '"':
			buffer[pos] = '\\'
			buffer[pos+1] = '"'
			pos += 2
		default:
			buffer[pos] = c
			pos += 1
		}
	}

	if pos == bl-5 && i < len(value) {
		copy(buffer[pos:], "...")
		pos += 3
	}

	buffer[pos] = '"'
	return pos + 1
}
