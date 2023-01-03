package log

import (
	"errors"
	"strconv"
	"strings"
	"testing"
	"time"

	"src.goblgobl.com/tests/assert"
)

func Test_KvLogger_Int(t *testing.T) {
	out := &strings.Builder{}
	l := KvFactory(128)(nil, INFO)

	l.Info("i").Int("ms", 0).LogTo(out)
	assertKvLog(t, out, false, map[string]string{"ms": "0"})

	l.Info("i").Int("count", 32).String("x", "b").LogTo(out)
	assertKvLog(t, out, false, map[string]string{"count": "32", "x": "b"})

	l.Warn("i").Int("ms", -99).LogTo(out)
	assertKvLog(t, out, false, map[string]string{"ms": "-99"})
}

func Test_KvLogger_Error(t *testing.T) {
	out := &strings.Builder{}
	l := KvFactory(128)(nil, INFO)
	l.Warn("w").Err(errors.New("test_error")).LogTo(out)
	assertKvLog(t, out, false, map[string]string{"err": "test_error"})
}

func Test_KvLogger_StructuredError_NoData(t *testing.T) {
	out := &strings.Builder{}
	l := KvFactory(128)(nil, INFO)
	se := Err(299, errors.New("test_error"))

	l.Warn("w").Err(se).LogTo(out)
	assertKvLog(t, out, false, map[string]string{
		"code": "299",
		"err":  "test_error",
	})
}

func Test_KvLogger_StructuredError_Data(t *testing.T) {
	out := &strings.Builder{}
	l := KvFactory(128)(nil, INFO)
	se := Err(311, errors.New("test_error2")).String("a", "z").Int("zero", 0)

	l.Warn("w").Err(se).LogTo(out)
	assertKvLog(t, out, false, map[string]string{
		"a":    "z",
		"zero": "0",
		"code": "311",
		"err":  "test_error2",
	})
}

func Test_KvLogger_Timestamp(t *testing.T) {
	out := &strings.Builder{}
	l := KvFactory(128)(nil, INFO)

	l.Info("hi").LogTo(out)
	fields := assertKvLog(t, out, false, nil)
	unix, _ := strconv.Atoi(fields["t"])
	assert.Nowish(t, time.Unix(int64(unix), 0))
}

func Test_KvLogger_UnencodedLenghts(t *testing.T) {
	out := &strings.Builder{}
	// info or warn messages take 23 characters + context length
	l := KvFactory(35)(nil, INFO)

	l.Info("ctx1").String("a", "1").LogTo(out)
	assertKvLog(t, out, false, map[string]string{"a": "1"})

	s := string(l.Info("ctx1").String("a", "1").Bytes())
	assert.StringContains(t, s, "l=info")
	assert.StringContains(t, s, "c=ctx1 a=1")
	l.Reset()

	l.Info("ctx1").String("a", "12").LogTo(out)
	assertKvLog(t, out, false, map[string]string{"a": "12"})

	l.Info("ctx1").String("a", "123").LogTo(out)
	assertKvLog(t, out, false, map[string]string{"a": "123"})

	l.Info("ctx1").String("a", "1234").LogTo(out)
	assertKvLog(t, out, false, map[string]string{"a": "1234"})

	l.Info("ctx1").String("a", "12345").LogTo(out)
	assertKvLog(t, out, false, map[string]string{"a": "12345"})

	l.Info("ctx2").String("ab", "1").LogTo(out)
	assertKvLog(t, out, false, map[string]string{"ab": "1"})

	l.Info("ctx2").String("ab", "12").LogTo(out)
	assertKvLog(t, out, false, map[string]string{"ab": "12"})

	l.Info("ctx2").String("ab", "123").LogTo(out)
	assertKvLog(t, out, false, map[string]string{"ab": "123"})

	l.Info("ctx2").String("ab", "1234").LogTo(out)
	assertKvLog(t, out, false, map[string]string{"ab": "1234"})

	l.Info("ctx1").String("a", "123456").LogTo(out)
	assertNoField(t, out, "a")

	l.Info("ctx1").String("ab", "12345").LogTo(out)
	assertNoField(t, out, "ab")
}

func Test_KvLogger_EncodedLenghts(t *testing.T) {
	out := &strings.Builder{}
	// info or warn messages take 23 characters + context length
	l := KvFactory(40)(nil, INFO)

	l.Info("ctx1").String("a", "\"").LogTo(out)
	assertKvLog(t, out, false, map[string]string{"a": `"\""`})

	l.Info("ctx1").String("a", "1\"").LogTo(out)
	assertKvLog(t, out, false, map[string]string{"a": `"1\""`})

	l.Info("ctx1").String("a", "1\"b").LogTo(out)
	assertKvLog(t, out, false, map[string]string{"a": `"1\"b"`})

	l.Info("ctx1").String("a", "1\"bc").LogTo(out)
	assertKvLog(t, out, false, map[string]string{"a": `"1\"bc"`})

	l.Info("ctx1").String("a", "1\"bcd").LogTo(out)
	assertKvLog(t, out, false, map[string]string{"a": `"1\"bc..."`})

	l.Info("ctx1").String("a", "1\"bcde").LogTo(out)
	assertKvLog(t, out, false, map[string]string{"a": `"1\"bc..."`})

	l.Info("ctx1").String("ab", "\"").LogTo(out)
	assertKvLog(t, out, false, map[string]string{"ab": `"\""`})

	l.Info("ctx1").String("ab", "1\"").LogTo(out)
	assertKvLog(t, out, false, map[string]string{"ab": `"1\""`})

	l.Info("ctx1").String("ab", "1\"b").LogTo(out)
	assertKvLog(t, out, false, map[string]string{"ab": `"1\"b"`})

	l.Info("ctx1").String("ab", "1\"bc").LogTo(out)
	assertKvLog(t, out, false, map[string]string{"ab": `"1\"b..."`})

	l.Info("ctx1").String("ab", "1\"bcd").LogTo(out)
	assertKvLog(t, out, false, map[string]string{"ab": `"1\"b..."`})
}

func Test_KvLogger_Fixed(t *testing.T) {
	out := &strings.Builder{}
	l := KvFactory(128)(nil, INFO)

	l.Field(NewField().Int("power", 9001).Finalize()).Fixed()
	l.LogTo(out)
	assert.Equal(t, out.String(), "power=9001\n")

	out.Reset()
	l.Reset()

	l.Info("x").String("a", "b").LogTo(out)
	assertKvLog(t, out, true, map[string]string{
		"l":     "info",
		"c":     "x",
		"a":     "b",
		"power": "9001",
	})
}

func Test_KvLogger_MultiUse_Common(t *testing.T) {
	out := &strings.Builder{}
	l := KvFactory(128)(nil, INFO)

	l.Field(NewField().String("id", "123").Finalize()).MultiUse()
	l.LogTo(out)
	assert.Equal(t, out.String(), "id=123\n")

	out.Reset()
	l.Info("a").LogTo(out)
	assertKvLog(t, out, true, map[string]string{
		"l":  "info",
		"c":  "a",
		"id": "123",
	})

	l.Release()
	l.Info("x").LogTo(out)
	fields := assertKvLog(t, out, true, map[string]string{
		"l": "info",
		"c": "x",
	})
	assert.Equal(t, len(fields), 3) // +1 for time
}

func Test_Logger_FixedAndMultiUse(t *testing.T) {
	out := &strings.Builder{}
	l := KvFactory(128)(nil, INFO)

	l.Field(NewField().String("f", "one").Finalize()).Fixed()
	l.Field(NewField().Int("m", 2).Finalize()).MultiUse()
	l.LogTo(out)
	assert.Equal(t, out.String(), "f=one m=2\n")

	out.Reset()

	l.Error("e").LogTo(out)
	assertKvLog(t, out, true, map[string]string{
		"l": "error",
		"c": "e",
		"f": "one",
		"m": "2",
	})

	l.Fatal("f").LogTo(out)
	assertKvLog(t, out, true, map[string]string{
		"l": "fatal",
		"c": "f",
		"f": "one",
		"m": "2",
	})

	l.Reset()

	l.Fatal("f2").LogTo(out)
	assertKvLog(t, out, true, map[string]string{
		"l": "fatal",
		"c": "f2",
		"f": "one",
	})
}

func assertKvLog(t *testing.T, out *strings.Builder, strict bool, expected map[string]string) map[string]string {
	t.Helper()
	lookup := KvParse(out.String())

	if lookup == nil {
		assert.Nil(t, expected)
		return nil
	}

	for expectedKey, expectedValue := range expected {
		assert.Equal(t, lookup[expectedKey], expectedValue)
	}

	if strict {
		// -1 to remove the timestamp
		assert.Equal(t, len(lookup)-1, len(expected))
	}

	out.Reset()
	return lookup
}

func assertNoField(t *testing.T, out *strings.Builder, field string) {
	t.Helper()
	fields := assertKvLog(t, out, false, nil)
	_, exists := fields[field]
	assert.False(t, exists)
}
