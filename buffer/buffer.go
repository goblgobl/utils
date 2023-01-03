package buffer

import (
	"errors"
	"fmt"

	"src.goblgobl.com/utils"
)

/*
A wrapper around a []byte with helper methods for writing.
The buffer is also optionally pool-aware and satisfies io.Reader
and io.Closer interfaces.

While it's general-purpose, a goal is to interact with
fasthttp's Response.SetBodyStream to optimize how data is
written to a response.

The buffer has a minimum and maximum size. The minimum buffer
size is allocated upfront. If we need more space that minimum
but less than maximum, we'll dynamically allocated more memory.
However, when the buffer is reset/released back into the pool,
the dynamically allocated "large" buffer is discard and the
pre-allocated minimal buffer is restored.

It doesn't really matter to the implementation, but the maximum
size is something we expect to change on each usage, as many max
sizes are project-specific and, most likely, pools of buffers
are going to be cross-project.
*/

var (
	ErrMaxSize = errors.New("buffer maximum size")
)

type Buffer struct {
	// Writes might fail due to a full buffer (when we've reached
	// our maximum size). Rather than having each call need to
	// check for err, we just noop every write operation when
	// err != nil and return the error on reads.
	err error

	// the maximum size we'll allow this buffer to grow
	max int

	// the position within data our last write was at
	pos int

	// fixed-size and pre-allocated data that won't grow
	static []byte

	// active buffer to read/write from. Will either reference
	// data or a dynamically allocated larger space (up to max size)
	data []byte

	// could be nil (if created outside of the pool, or if the pool
	// was empty and created it on the fly)
	pool *Pool
}

func New(min uint32, max uint32) *Buffer {
	static := make([]byte, min)
	return &Buffer{
		data:   static,
		static: static,
		max:    int(max),
	}
}

func (b *Buffer) Reset() {
	b.pos = 0
	b.err = nil
	b.data = b.static
}

func (b *Buffer) Release() {
	if pool := b.pool; pool != nil {
		b.Reset()
		pool.list <- b
	}
}

func (b *Buffer) Len() int {
	return b.pos
}

func (b *Buffer) Max() int {
	return b.max
}

func (b *Buffer) Error() error {
	return b.err
}

func (b *Buffer) String() (string, error) {
	bytes, err := b.Bytes()
	return string(bytes), err
}

func (b *Buffer) UnsafeString() (string, error) {
	bytes, err := b.Bytes()
	return utils.B2S(bytes), err
}

func (b *Buffer) MustString() string {
	str, err := b.String()
	if err != nil {
		panic(err)
	}
	return str
}

func (b Buffer) Bytes() ([]byte, error) {
	return b.data[:b.pos], b.err
}

// optimization for passing the result to the ExecTerminated
// method of sqlite.Conn.
func (b Buffer) SqliteBytes() ([]byte, error) {
	b.WriteByte('\x00')
	return b.data[:b.pos], b.err
}

// Write and ensure enough capacity for len(data) + padSize
// Meant to be used with WriteByteUnsafe.
func (b *Buffer) WritePad(data []byte, padSize int) {
	if !b.ensureCapacity(len(data) + padSize) {
		return
	}

	pos := b.pos
	copy(b.data[pos:], data)
	b.pos = pos + len(data)
}

func (b *Buffer) Write(data []byte) {
	b.WritePad(data, 0)
}

func (b *Buffer) WriteUnsafe(data string) {
	b.Write(utils.S2B(data))
}

func (b *Buffer) WriteString(data string) {
	b.Write([]byte(data))
}

func (b *Buffer) WriteByte(byte byte) {
	if b.err != nil {
		return
	}
	if !b.ensureCapacity(1) {
		return
	}

	pos := b.pos
	b.data[pos] = byte
	b.pos = pos + 1
}

// Our caller knows that there's enough space in data
// (probably because it used WritePad)
func (b *Buffer) WriteByteUnsafe(byte byte) {
	pos := b.pos
	b.data[pos] = byte
	b.pos = pos + 1
}

func (b *Buffer) Truncate(n int) {
	b.pos -= n
}

func (b *Buffer) ensureCapacity(l int) bool {
	if b.err != nil {
		return false
	}

	data := b.data
	required := b.pos + l

	max := b.max
	if required > max {
		b.err = fmt.Errorf("buffer.ensureCapacity(%d, %d) - %w", required, max, ErrMaxSize)
		return false
	}

	// Whatever data is (static or dynamic), we have
	// enough space as-is. happiness.
	if required <= len(data) {
		return true
	}

	newLen := len(data) * 2
	if newLen < required {
		newLen = required
	} else if newLen > max {
		newLen = max
	}

	newData := make([]byte, newLen)
	copy(newData, data)
	b.data = newData
	return true
}
