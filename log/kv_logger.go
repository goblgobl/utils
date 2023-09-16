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
	"encoding/base64"
	"io"
	"strconv"
	"time"

	"src.goblgobl.com/utils"
	"src.goblgobl.com/utils/buffer"
)

var binaryEncoder = base64.RawURLEncoding

type KvLogger struct {
	release func(Logger)

	// buffer that we write our message to
	buffer *buffer.Buffer

	// The log level that we're logging.
	level Level

	// Whether or not we're logging request messages
	requests bool

	// A logger can have a fixed piece of data which is
	// always included (e.g pid=$PROJECT_ID for a project-owned
	// logger). Once our fixed data is set, pos will never be
	// less than fixedLen.
	fixedLen int64

	// A logger can also have temporary repeated data
	// (e.g. rid=$REQUEST_ID for an env-owned logger).
	// After logging a message, pos == multiUseLen. Only
	// on reset/release will pos == fixedLen
	multiUseLen uint64
}

func NewKvLogger(maxSize uint32, release func(Logger), level Level, requests bool) *KvLogger {
	return &KvLogger{
		level:    level,
		release:  release,
		requests: requests,
		buffer:   buffer.New(4096, maxSize),
	}
}

func KvFactory(maxSize uint32) Factory {
	return func(release func(Logger), level Level, requests bool) Logger {
		return NewKvLogger(maxSize, release, level, requests)
	}
}

// Get the bytes from the logger. This is only valid before Log is called (after
// log is called, you'll get an empty slice). Only really useful for testing.
func (l *KvLogger) Bytes() []byte {
	return l.buffer.OKBytes()
}

// Logger will _always_ include this data. Meant to be used with the Field builder.
// Even once released to the pool and re-checked out, this data will still be in the logger.
// For checkout-specific data, see MultiUse().
func (l *KvLogger) Fixed() {
	l.fixedLen = int64(l.buffer.Len())
}

// Similar to Fixed, but exists only while checked out
func (l *KvLogger) MultiUse() Logger {
	l.multiUseLen = uint64(l.buffer.Len())
	return l
}

// Add a field (key=value) where value is a string
func (l *KvLogger) String(key string, value string) Logger {
	l.writeKeyValue(key, value)
	return l
}

func (l *KvLogger) Binary(key string, value []byte) Logger {
	len := binaryEncoder.EncodedLen(len(value))
	if l.writeKeyForValueLen(key, len) {
		enc := base64.NewEncoder(binaryEncoder, l.buffer)
		enc.Write(value)
		enc.Close()
	}
	return l
}

// Add a field (key=value) where value is an int
func (l *KvLogger) Int(key string, value int) Logger {
	return l.Int64(key, int64(value))
}

// Add a field (key=value) where value is an int
func (l *KvLogger) Int64(key string, value int64) Logger {
	s := strconv.FormatInt(value, 10)
	if l.writeKeyForValue(key, s) {
		l.buffer.WriteString(s)
	}
	return l
}

// Add a field (key=value) where value is a boolean
func (l *KvLogger) Bool(key string, value bool) Logger {
	if l.writeKeyForValueLen(key, 1) {
		if value {
			l.buffer.WriteByte('Y')
		} else {
			l.buffer.WriteByte('N')
		}
	}
	return l
}

// Add a field (key=value) where value is an error
func (l *KvLogger) Err(err error) Logger {
	se, ok := err.(*StructuredError)
	if !ok {
		return l.String("_err", err.Error())
	}

	l.Int("_code", se.Code).String("_err", se.Err.Error())
	for key, value := range se.Data {
		switch v := value.(type) {
		case string:
			l.String(key, v)
		case int:
			l.Int(key, v)
		case []byte:
			l.Binary(key, v)
		}
	}
	return l
}

// Write the log to our globally configured writer
func (l *KvLogger) Log() {
	l.LogTo(Out)
}

func (l *KvLogger) LogTo(out io.Writer) {
	buffer := l.buffer

	buffer.WriteByte('\n')
	out.Write(buffer.OKBytes())
	l.conditionalRelease()
}

func (l *KvLogger) Reset() {
	l.buffer.Seek(l.fixedLen, io.SeekStart)
}

func (l *KvLogger) Release() {
	l.buffer.Reset()
	l.buffer.Seek(l.fixedLen, io.SeekStart)
	if release := l.release; release != nil {
		release(l)
	}
}

// Normally, logger is automatically released when Log or LogTo is called
// unless we've enabled multiUse.
func (l *KvLogger) conditionalRelease() {
	if l.multiUseLen == 0 {
		l.Release()
	} else {
		// remove out trailing newline
		l.buffer.Truncate(1)
	}
}

// Log an info-level message.
func (l *KvLogger) Info(ctx string) Logger {
	if l.level > INFO {
		l.conditionalRelease()
		return Noop{}
	}
	return l.start(ctx, []byte("_l=info _t="))
}

// Log an warn-level message.
func (l *KvLogger) Warn(ctx string) Logger {
	if l.level > WARN {
		l.conditionalRelease()
		return Noop{}
	}
	return l.start(ctx, []byte("_l=warn _t="))
}

// Log an error-level message.
func (l *KvLogger) Error(ctx string) Logger {
	if l.level > ERROR {
		l.conditionalRelease()
		return Noop{}
	}
	return l.start(ctx, []byte("_l=error _t="))
}

// Log an fatal-level message.
func (l *KvLogger) Fatal(ctx string) Logger {
	if l.level > FATAL {
		l.conditionalRelease()
		return Noop{}
	}
	return l.start(ctx, []byte("_l=fatal _t="))
}

// Log a request message.
func (l *KvLogger) Request(route string) Logger {
	if !l.requests {
		l.conditionalRelease()
		return Noop{}
	}
	return l.start(route, []byte("_l=req _t="))
}

func (l *KvLogger) Field(field Field) Logger {
	buffer := l.buffer

	// might already have data
	if buffer.Len() != 0 {
		buffer.WriteByte(' ')
	}
	buffer.Write(field.kv)
	return l
}

// "starts" a new log message. Every message always contains a timestamp (t) a
// context (c) and a level (l).
func (l *KvLogger) start(ctx string, meta []byte) Logger {
	buffer := l.buffer

	// len > 0 when MultiUse is enabled
	if buffer.Len() > 0 {
		buffer.WriteByte(' ')
	}

	// includes the log level + the _t= key (but not the timestamp itself)
	buffer.Write(meta)
	buffer.WriteString(strconv.FormatInt(time.Now().Unix(), 10))
	buffer.WriteString(" _c=")
	buffer.WriteString(ctx)
	return l
}

// Writes "$key=" and returns the position where the value can be written.
func (l *KvLogger) writeKeyForValue(key string, value string) bool {
	return l.writeKeyForValueLen(key, len(value))
}

// We expect key to always be safe to write as-is.
func (l *KvLogger) writeKeyForValueLen(key string, valueLen int) bool {
	return writeKeyForValueLen(key, valueLen, l.buffer)
}

// We only encode newline and quotes. If either is present, the value is quote encoded.
func (l *KvLogger) writeKeyValue(key string, value string) {
	writeKeyValue(key, value, l.buffer)
}

func writeKeyForValueLen(key string, valueLen int, buffer *buffer.Buffer) bool {
	// + 3 for space & equal & final newline
	if buffer.EnsureCapacity(len(key)+valueLen+3) == false {
		return false
	}

	if buffer.Len() > 0 {
		buffer.WriteByte(' ')
	}

	buffer.WriteString(key)
	buffer.WriteByte('=')
	return true
}

func writeKeyValue(key string, value string, buffer *buffer.Buffer) {
	escapeCount := escapeCount(value)

	spaceRequiredForValue := len(value)
	if escapeCount > 0 {
		// +2 because we need to wrap the entire thing in double quotes
		spaceRequiredForValue += escapeCount + 2
	}

	if writeKeyForValueLen(key, spaceRequiredForValue, buffer) == false {
		return
	}

	if escapeCount == 0 {
		buffer.WriteString(value)
		return
	}

	buffer.WriteByteUnsafe('"')

	for _, c := range utils.S2B(value) {
		switch c {
		case '\n':
			buffer.WriteByteUnsafe('\\')
			buffer.WriteByteUnsafe('n')
		case '"':
			buffer.WriteByteUnsafe('\\')
			buffer.WriteByteUnsafe('"')
		default:
			buffer.WriteByteUnsafe(c)
		}
	}

	buffer.WriteByteUnsafe('"')
}

func escapeCount(input string) int {
	count := 0
	for i := 0; i < len(input); i++ {
		c := input[i]
		if c == '=' || c == '"' || c == '\n' || c == ' ' {
			count += 1
		}
	}
	return count
}
