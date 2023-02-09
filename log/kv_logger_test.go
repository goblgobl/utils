package log

import (
	"encoding/base64"
	"errors"
	"strconv"
	"strings"
	"testing"
	"time"

	"src.goblgobl.com/tests/assert"
)

func Test_KvLogger_Int(t *testing.T) {
	out := &strings.Builder{}
	l := KvFactory(128)(nil, INFO, true)

	l.Info("i").Int("ms", 0).LogTo(out)
	assertKvLog(t, out, false, map[string]string{"ms": "0"})

	l.Info("i").Int("count", 32).String("x", "b").LogTo(out)
	assertKvLog(t, out, false, map[string]string{"count": "32", "x": "b"})

	l.Warn("i").Int("ms", -99).LogTo(out)
	assertKvLog(t, out, false, map[string]string{"ms": "-99"})
}

func Test_KvLogger_Binary(t *testing.T) {
	out := &strings.Builder{}
	l := KvFactory(128)(nil, INFO, true)

	l.Info("i").Binary("ms", []byte{1, 2, 3}).LogTo(out)
	assertKvLog(t, out, false, map[string]string{"ms": "AQID"})
	// a bit backwards...
	reverse, _ := base64.RawURLEncoding.DecodeString("AQID")
	assert.Bytes(t, []byte{1, 2, 3}, reverse)
}

func Test_KvLogger_Bool(t *testing.T) {
	out := &strings.Builder{}
	l := KvFactory(128)(nil, INFO, true)

	l.Info("i").Bool("active", true).LogTo(out)
	assertKvLog(t, out, false, map[string]string{"active": "Y"})

	l.Info("i").Bool("active", false).String("x", "b").LogTo(out)
	assertKvLog(t, out, false, map[string]string{"active": "N", "x": "b"})
}

func Test_KvLogger_Error(t *testing.T) {
	out := &strings.Builder{}
	l := KvFactory(128)(nil, INFO, true)
	l.Warn("w").Err(errors.New("test_error")).LogTo(out)
	assertKvLog(t, out, false, map[string]string{"_err": "test_error"})
}

func Test_KvLogger_StructuredError_NoData(t *testing.T) {
	out := &strings.Builder{}
	l := KvFactory(128)(nil, INFO, true)
	se := Err(299, errors.New("test_error"))

	l.Warn("w").Err(se).LogTo(out)
	assertKvLog(t, out, false, map[string]string{
		"_code": "299",
		"_err":  "test_error",
	})
}

func Test_KvLogger_StructuredError_Data(t *testing.T) {
	out := &strings.Builder{}
	l := KvFactory(128)(nil, INFO, true)
	se := Err(311, errors.New("test_error2")).String("a", "z").Int("zero", 0)

	l.Warn("w").Err(se).LogTo(out)
	assertKvLog(t, out, false, map[string]string{
		"a":     "z",
		"zero":  "0",
		"_code": "311",
		"_err":  "test_error2",
	})
}

func Test_KvLogger_StructuredError_Nesting_NoData(t *testing.T) {
	out := &strings.Builder{}
	l := KvFactory(128)(nil, INFO, true)
	s1 := Err(311, errors.New("test_error2"))
	s2 := Err(312, s1)

	l.Warn("w").Err(s2).LogTo(out)
	assertKvLog(t, out, false, map[string]string{
		"_code": "312",
		"_err":  `"code: 311 - test_error2"`,
	})
}

func Test_KvLogger_StructuredError_Nesting_Data(t *testing.T) {
	out := &strings.Builder{}
	l := KvFactory(128)(nil, INFO, true)
	s1 := Err(311, errors.New("test_error2")).String("id", "a").Int("x", 9)
	s2 := Err(312, s1).String("other", "b").Int("x", 8)

	l.Warn("w").Err(s2).LogTo(out)
	assertKvLog(t, out, false, map[string]string{
		"id":     "a",
		"other":  "b",
		"x":      "8",
		"_code":  "312",
		"_icode": "311",
		"_err":   `"code: 311 - test_error2"`,
	})

	out.Reset()
	s3 := ErrData(312, s1, map[string]any{"other": "b2", "x": 10})
	l.Warn("w").Err(s3).LogTo(out)
	assertKvLog(t, out, false, map[string]string{
		"id":     "a",
		"other":  "b2",
		"x":      "10",
		"_code":  "312",
		"_icode": "311",
		"_err":   `"code: 311 - test_error2"`,
	})
}

func Test_KvLogger_Timestamp(t *testing.T) {
	out := &strings.Builder{}
	l := KvFactory(128)(nil, INFO, true)

	l.Info("hi").LogTo(out)
	fields := assertKvLog(t, out, false, nil)
	unix, _ := strconv.Atoi(fields["_t"])
	assert.Nowish(t, time.Unix(int64(unix), 0))
}

func Test_KvLogger_UnencodedLenghts(t *testing.T) {
	out := &strings.Builder{}
	// info or warn messages take 26 characters + context length
	l := KvFactory(38)(nil, INFO, true)

	l.Info("ctx1").String("a", "1").LogTo(out)
	assertKvLog(t, out, false, map[string]string{"a": "1"})

	s := string(l.Info("ctx1").String("a", "1").Bytes())
	assert.StringContains(t, s, "_l=info")
	assert.StringContains(t, s, "_c=ctx1 a=1")
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
	// info or warn messages take 26 characters + context length
	l := KvFactory(43)(nil, INFO, true)

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
	l := KvFactory(128)(nil, INFO, true)

	l.Field(NewField().Int("power", 9001).Finalize()).Fixed()
	l.LogTo(out)
	assert.Equal(t, out.String(), "power=9001\n")

	out.Reset()
	l.Reset()

	l.Info("x").String("a", "b").LogTo(out)
	assertKvLog(t, out, true, map[string]string{
		"_l":    "info",
		"_c":    "x",
		"a":     "b",
		"power": "9001",
	})
}

func Test_KvLogger_MultiUse_Common(t *testing.T) {
	out := &strings.Builder{}
	l := KvFactory(128)(nil, INFO, true)

	l.Field(NewField().String("id", "123").Finalize()).MultiUse()
	l.LogTo(out)
	assert.Equal(t, out.String(), "id=123\n")

	out.Reset()
	l.Info("a").LogTo(out)
	assertKvLog(t, out, true, map[string]string{
		"_l": "info",
		"_c": "a",
		"id": "123",
	})

	l.Release()
	l.Info("x").LogTo(out)
	fields := assertKvLog(t, out, true, map[string]string{
		"_l": "info",
		"_c": "x",
	})
	assert.Equal(t, len(fields), 3) // +1 for time
}

func Test_Logger_FixedAndMultiUse(t *testing.T) {
	out := &strings.Builder{}
	l := KvFactory(128)(nil, INFO, true)

	l.Field(NewField().String("f", "one").Finalize()).Fixed()
	l.Field(NewField().Int("m", 2).Finalize()).MultiUse()
	l.LogTo(out)
	assert.Equal(t, out.String(), "f=one m=2\n")

	out.Reset()

	l.Error("e").LogTo(out)
	assertKvLog(t, out, true, map[string]string{
		"_l": "error",
		"_c": "e",
		"f":  "one",
		"m":  "2",
	})

	l.Fatal("f").LogTo(out)
	assertKvLog(t, out, true, map[string]string{
		"_l": "fatal",
		"_c": "f",
		"f":  "one",
		"m":  "2",
	})

	l.Reset()

	l.Fatal("f2").LogTo(out)
	assertKvLog(t, out, true, map[string]string{
		"_l": "fatal",
		"_c": "f2",
		"f":  "one",
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
