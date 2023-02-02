package buffer

import (
	"bytes"
	"io"
	"testing"

	"src.goblgobl.com/tests/assert"
)

func Test_Buffer_Containing(t *testing.T) {
	b := Containing([]byte("hello world"), 20)
	assert.Equal(t, b.Len(), 11)
	b.WriteByte('!')
	assert.Equal(t, testMustString(b), "hello world!")
}

func Test_Buffer_Write_NoGrow(t *testing.T) {
	b := New(10, 20)
	assert.Equal(t, b.Len(), 0)
	assert.Equal(t, b.Max(), 20)

	b.Write([]byte("hello"))
	assert.Nil(t, b.Error())
	assert.Equal(t, b.Len(), 5)
	assert.Equal(t, testMustString(b), "hello")
	assert.Equal(t, &b.data[0], &b.static[0])

	b.Write([]byte("world"))
	assert.Nil(t, b.Error())
	assert.Equal(t, b.Len(), 10)
	assert.Equal(t, testMustString(b), "helloworld")
	assert.Equal(t, &b.data[0], &b.static[0])
}

func Test_Buffer_WriteByte_NoGrow(t *testing.T) {
	b := New(2, 20)
	b.WriteByte('a')
	assert.Equal(t, testMustString(b), "a")
	assert.Equal(t, &b.data[0], &b.static[0])

	b.WriteByte('b')
	assert.Equal(t, testMustString(b), "ab")
	assert.Equal(t, &b.data[0], &b.static[0])
}

func Test_Buffer_Write_GrowFirst(t *testing.T) {
	b := New(2, 20)
	b.Write([]byte("hello"))
	assert.Equal(t, testMustString(b), "hello")
	assert.NotEqual(t, &b.data[0], &b.static[0])
}

func Test_Buffer_Write_Grow(t *testing.T) {
	b := New(8, 20)
	b.Write([]byte("hello"))
	b.WriteString(" world")
	assert.Equal(t, testMustString(b), "hello world")
	assert.NotEqual(t, &b.data[0], &b.static[0])
}

func Test_Buffer_WriteByte_GrowFirst(t *testing.T) {
	b := New(0, 20)
	b.WriteByte('z')
	assert.Equal(t, testMustString(b), "z")
}

func Test_Buffer_WriteByte_Grow(t *testing.T) {
	b := New(1, 20)
	b.WriteByte('y')
	b.WriteByte('z')
	assert.Equal(t, testMustString(b), "yz")
	assert.NotEqual(t, &b.data[0], &b.static[0])
}

func Test_Buffer_TakeBytes(t *testing.T) {
	b := New(20, 40)
	b.WriteByte('z')

	// this resets the buffer
	bytes, err := b.TakeBytes(5)
	assert.Nil(t, err)
	assert.Equal(t, len(bytes), 5)
}

func Test_Buffer_Grow_Doubling(t *testing.T) {
	b := New(4, 25)
	assert.Equal(t, b.Max(), 25)

	// double
	b.Write([]byte("abcde"))
	assert.Equal(t, len(b.data), 8)

	// need more than double, grow to exact
	b.Write([]byte("it's over 9000"))
	assert.Equal(t, len(b.data), 19)

	// double > max, go to max
	b.Write([]byte("yes"))
	assert.Equal(t, len(b.data), 25)
}

func Test_Buffer_Grow_MaxSize(t *testing.T) {
	b := New(4, 8)
	b.Write([]byte("hello world"))
	assert.Equal(t, b.Error().Error(), "code: 3005 - buffer maximum size")

	s, err := b.String()
	assert.Equal(t, s, "")
	assert.Equal(t, err.Error(), "code: 3005 - buffer maximum size")
}

// Totally valid. Min is just statically allocated. But in
// some cases, we might want to limit the actual size
// to some smaller amount
func Test_Buffer_MaxSize_LessThan_Min(t *testing.T) {
	b := New(10, 5)
	b.Write([]byte("over"))

	s, err := b.UnsafeString()
	assert.Nil(t, err)
	assert.Equal(t, s, "over")

	b.Write([]byte(" 9000"))
	s, err = b.String()
	assert.Equal(t, s, "over")
	assert.Equal(t, err.Error(), "code: 3005 - buffer maximum size")
}

func Test_Buffer_Reset_Normal(t *testing.T) {
	b := New(4, 8)
	static := b.static
	b.WriteByte('1')

	b.Reset()
	assert.Nil(t, b.err)
	assert.Equal(t, b.Len(), 0)
	assert.Equal(t, testMustString(b), "")
	assert.Equal(t, &b.data[0], &static[0])
	assert.Equal(t, &b.static[0], &static[0])
}

func Test_Buffer_Reset_Grown(t *testing.T) {
	b := New(4, 8)
	static := b.static
	b.Write([]byte("12345"))

	b.Reset()
	assert.Nil(t, b.err)
	assert.Equal(t, testMustString(b), "")
	assert.Equal(t, testMustString(b), "")
	assert.Equal(t, &b.data[0], &static[0])
	assert.Equal(t, &b.static[0], &static[0])
}

func Test_Buffer_Reset_Error(t *testing.T) {
	b := New(4, 4)
	b.Write([]byte("12345"))
	assert.NotNil(t, b.err)

	b.Reset()
	assert.Nil(t, b.err)
}

func Test_Buffer_WritePad(t *testing.T) {
	b := New(4, 20)

	// without padding, it would have doubled to 8
	b.WritePad([]byte("12345678"), 2)
	assert.Equal(t, cap(b.data), 10)
	assert.Equal(t, testMustString(b), "12345678")

	b.WriteByteUnsafe('9')
	assert.Equal(t, testMustString(b), "123456789")

	b.WriteByteUnsafe('A')
	assert.Equal(t, cap(b.data), 10)
	assert.Equal(t, testMustString(b), "123456789A")
}

func Test_Buffer_Truncate(t *testing.T) {
	b := New(4, 20)

	b.Write([]byte("12345"))
	b.Truncate(3)
	assert.Equal(t, testMustString(b), "12")
	b.Truncate(1)
	assert.Equal(t, testMustString(b), "1")
	b.Truncate(1)
	assert.Equal(t, testMustString(b), "")
}

func Test_Buffer_SqliteBytes(t *testing.T) {
	b := New(5, 20)
	b.Write([]byte("up"))
	sql, err := b.SqliteBytes()
	assert.Nil(t, err)
	assert.Equal(t, string(sql), "up\x00")

	// terminator forces growth
	b = New(5, 20)
	b.Write([]byte("hello"))
	sql, err = b.SqliteBytes()
	assert.Nil(t, err)
	assert.Equal(t, string(sql), "hello\x00")
}

func Test_Buffer_Reader(t *testing.T) {
	b := New(6000, 7000)
	b.Write(bytes.Repeat([]byte("a"), 5000))

	var out bytes.Buffer
	n, err := io.Copy(&out, b)
	assert.Nil(t, err)
	assert.Equal(t, n, 5000)
	assert.Equal(t, out.Len(), 5000)

	b.Reset()
	out.Reset()
	n, err = io.Copy(&out, b)
	assert.Nil(t, err)
	assert.Equal(t, n, 0)
	assert.Equal(t, out.Len(), 0)
}

func testMustString(b *Buffer) string {
	s, err := b.String()
	if err != nil {
		panic(err)
	}
	return s
}
