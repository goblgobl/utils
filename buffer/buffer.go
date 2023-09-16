package buffer

import (
	"errors"
	"io"

	"src.goblgobl.com/utils"
)

/*
A wrapper around a []byte with helper methods for writing.
The buffer is also optionally pool-aware and satisfies the io.Reader
and io.Closer interfaces.

While it's general-purpose, a goal is to interact with
fasthttp's Response.SetBodyStream to optimize how data is
written to a response.

The buffer has a minimum and maximum size. The minimum buffer
size is allocated upfront. If we need more space than minimum
but less than maximum, we'll dynamically allocated more memory.
However, when the buffer is reset/released back into the pool,
the dynamically allocated "large" buffer is discard and the
pre-allocated minimal buffer is restored.
*/

var (
	ErrMaxSize = errors.New("code: 3005 - buffer maximum size")
	Empty      = new(Buffer)
)

type Buffer struct {
	release func(*Buffer)

	// Writes might fail due to a full buffer (when we've reached
	// our maximum size). Rather than having each call need to
	// check for err, we just noop every write operation when
	// err != nil and return the error on reads.
	err error

	// fixed-size and pre-allocated data that won't grow
	static []byte

	// active buffer to read/write from. Will either reference
	// data or a dynamically allocated larger space (up to max size)
	data []byte

	// the maximum size we'll allow this buffer to grow
	max int

	// the position within data our last write was at
	pos int

	// the position within out data our last read was at
	read int
}

func New(min uint32, max uint32) *Buffer {
	static := make([]byte, min)
	return &Buffer{
		data:   static,
		static: static,
		max:    int(max),
	}
}

// create a buffer that contains the specified data
// (mostly useful for tests)
func Containing(data []byte, max int) *Buffer {
	return &Buffer{
		data:   data,
		static: data,
		max:    max,
		pos:    len(data),
	}
}

func (b *Buffer) Reset() {
	b.pos = 0
	b.read = 0
	b.err = nil
	b.data = b.static
}

func (b *Buffer) Release() {
	if release := b.release; release != nil {
		b.Reset()
		release(b)
	}
}

// io.Closer
func (b *Buffer) Close() error {
	b.Release()
	return nil
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

func (b Buffer) OKBytes() []byte {
	return b.data[:b.pos]
}

// An advanced function. Some code might want to manipulate a []byte directly
// while leveraging the pre-allocated nature of a pooled buffer. Importantly,
// the returned bytes are not cleared
func (b Buffer) TakeBytes(l int) ([]byte, error) {
	b.Reset()
	b.EnsureCapacity(l)
	return b.data[:l], b.err
}

// optimization for passing the result to the ExecTerminated
// method of sqlite.Conn.
func (b Buffer) SqliteBytes() ([]byte, error) {
	b.WriteByte('\x00')
	return b.data[:b.pos], b.err
}

// ensure that we have enough space for padSize
func (b *Buffer) Pad(padSize int) error {
	if !b.EnsureCapacity(padSize) {
		return b.err
	}
	return nil
}

// Write and ensure enough capacity for len(data) + padSize
// Meant to be used with WriteByteUnsafe.
func (b *Buffer) WritePad(data []byte, padSize int) (int, error) {
	if !b.EnsureCapacity(len(data) + padSize) {
		return 0, b.err
	}

	pos := b.pos
	copy(b.data[pos:], data)

	l := len(data)
	b.pos = pos + l
	return l, nil
}

func (b *Buffer) Write(data []byte) (int, error) {
	return b.WritePad(data, 0)
}

func (b *Buffer) WriteString(data string) {
	b.Write(utils.S2B(data))
}

func (b *Buffer) WriteByte(byte byte) {
	if b.err != nil {
		return
	}
	if !b.EnsureCapacity(1) {
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

func (b *Buffer) Seek(offset int64, whence int) (int64, error) {
	var target int64
	switch whence {
	case io.SeekStart:
		target = offset
	case io.SeekCurrent:
		target = int64(b.pos) + offset
	case io.SeekEnd:
		target = int64(len(b.data)) + offset
	}

	if target < 0 {
		return int64(b.pos), errors.New("Seek before start")
	}
	if target > int64(len(b.data)) {
		return int64(b.pos), errors.New("Seek after end")
	}
	b.pos = int(target)
	return target, nil
}

// io.Reader
func (b *Buffer) Read(p []byte) (int, error) {
	if err := b.err; err != nil {
		return 0, err
	}

	pos := b.pos
	read := b.read
	if pos <= read {
		b.read = 0 // reset this so it can be read again
		if len(p) == 0 {
			return 0, nil
		}
		return 0, io.EOF
	}

	n := copy(p, b.data[read:pos])
	b.read = read + n
	return n, nil
}

func (b *Buffer) EnsureCapacity(l int) bool {
	if b.err != nil {
		return false
	}

	data := b.data
	required := b.pos + l

	max := b.max
	if required > max {
		b.err = ErrMaxSize
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

	//TODO: track how often we expand beyond our static buffer

	newData := make([]byte, newLen)
	copy(newData, data)
	b.data = newData
	return true
}
