package log

import (
	"fmt"
	"strconv"
)

/*
A lot of the field we want to log in an entry are either know
at compile time, or used multiple times. For example a 404 response
will want to log an entry with a status=404 field.

For such cases, we can create a Field object and use that
when creating an entry. Depending on the logging format, this might allow
some optimizations to take place.
*/

type Field struct {
	fields map[string]any
	kv     []byte
}

func NewField() *Field {
	return &Field{
		fields: make(map[string]any, 1),
	}
}

func (f *Field) KV() []byte {
	return f.kv
}

func (f *Field) Int(key string, value int) *Field {
	f.fields[key] = value
	return f
}

func (f *Field) String(key string, value string) *Field {
	f.fields[key] = value
	return f
}

// just return Field so that it can be used in chaining
func (f *Field) Finalize() Field {
	pos := uint64(0)
	buffer := make([]byte, 1024)

	for key, value := range f.fields {
		switch v := value.(type) {
		case int:
			pos = writeKeyValue(key, strconv.FormatInt(int64(v), 10), pos, buffer)
		case string:
			pos = writeKeyValue(key, v, pos, buffer)
		default:
			panic(fmt.Sprintf("unsupport field value type: %T (%v)", value, value))
		}
	}

	// We expect fields to be created on startup and be long-lived
	// we should trim out kv data to the exact size to avoid
	// wasting space
	kv := make([]byte, pos)
	copy(kv, buffer)
	f.kv = kv

	return *f
}
