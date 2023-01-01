package typed

import (
	"sort"
	"testing"
	"time"

	"src.goblgobl.com/tests/assert"
	"src.goblgobl.com/utils/json"
)

func Test_Must(t *testing.T) {
	typed := Must([]byte(`{"power": 9001}`))
	assert.Equal(t, typed.Int("power"), 9001)

	defer mustTest(t, "expected { character for map value")
	Must([]byte(`"h`))
	t.FailNow()
}

func Test_Nil_Json(t *testing.T) {
	typed, err := Json(nil)
	assert.Nil(t, err)
	assert.Equal(t, len(typed), 0)
}

func Test_Json(t *testing.T) {
	typed, err := Json([]byte(`{"power": 898887678118296}`))
	assert.Equal(t, typed.Int("power"), 898887678118296)
	assert.Nil(t, err)
}

func Test_JsonFile(t *testing.T) {
	typed, err := JsonFile("test.json")
	assert.Equal(t, typed.String("name"), "leto")
	assert.Nil(t, err)

	typed, err = JsonFile("invalid.json")
	assert.Equal(t, err.Error(), "open invalid.json: no such file or directory")
}

func Test_Keys(t *testing.T) {
	typed := New(build("name", "leto", "type", []int{1, 2, 3}, "number", 1))
	keys := typed.Keys()
	sort.Strings(keys)
	assert.List(t, keys, []string{"name", "number", "type"})
}

func Test_Bool(t *testing.T) {
	typed := New(build("log", true, "ace", false, "nope", 99, "s1", "true", "s2", "True", "s3", "TRUE", "s4", 1, "s5", "false", "s6", "False", "s7", "FALSE", "s8", 0))
	assert.Equal(t, typed.Bool("log"), true)
	assert.Equal(t, typed.BoolOr("log", false), true)
	assert.Equal(t, typed.Bool("other"), false)
	assert.Equal(t, typed.BoolOr("other", true), true)

	assert.True(t, typed.BoolMust("log"))
	assert.False(t, typed.BoolMust("ace"))

	// coerce to 'true', 'True', 'TRUE', 1
	assert.True(t, typed.BoolMust("s1"))
	assert.True(t, typed.BoolMust("s2"))
	assert.True(t, typed.BoolMust("s3"))
	assert.True(t, typed.BoolMust("s4"))

	// coerce to 'false', 'False', 'FALSE', 0
	assert.False(t, typed.BoolMust("s5"))
	assert.False(t, typed.BoolMust("s6"))
	assert.False(t, typed.BoolMust("s7"))
	assert.False(t, typed.BoolMust("s8"))

	value, exists := typed.BoolIf("nope")
	assert.False(t, value)
	assert.False(t, exists)

	defer mustTest(t, "expected boolean value for fail")
	typed.BoolMust("fail")
	t.FailNow()
}

func Test_Int(t *testing.T) {
	typed := New(build("port", 84, "string", "30", "i16", int16(1), "i32", int32(2), "i64", int64(3), "f64", float64(4), "number", json.Number("5"), "nope", true, "decimal", 3.1))
	assert.Equal(t, typed.Int("port"), 84)
	assert.Equal(t, typed.IntOr("port", 11), 84)
	value, exists := typed.IntIf("port")
	assert.Equal(t, value, 84)
	assert.True(t, exists)

	assert.Equal(t, typed.Int("string"), 30)
	assert.Equal(t, typed.IntOr("string", 11), 30)
	value, exists = typed.IntIf("string")
	assert.Equal(t, value, 30)
	assert.True(t, exists)

	assert.Equal(t, typed.Int("other"), 0)
	assert.Equal(t, typed.IntOr("other", 33), 33)
	value, exists = typed.IntIf("other")
	assert.Equal(t, value, 0)
	assert.False(t, exists)

	value, exists = typed.IntIf("nope")
	assert.Equal(t, value, 0)
	assert.False(t, exists)

	value, exists = typed.IntIf("decimal")
	assert.Equal(t, value, 0)
	assert.False(t, exists)

	assert.Equal(t, typed.Int("i16"), 1)
	assert.Equal(t, typed.Int("i32"), 2)
	assert.Equal(t, typed.Int("i64"), 3)
	assert.Equal(t, typed.Int("f64"), 4)

	assert.Equal(t, typed.IntMust("port"), 84)

	defer mustTest(t, "expected int value for fail")
	typed.IntMust("fail")
	t.FailNow()
}

func Test_Float(t *testing.T) {
	typed := New(build("pi", 3.14, "string", "30.14", "nope", true, "int", 33))
	assert.Equal(t, typed.Float("int"), 33)
	assert.Equal(t, typed.Float("pi"), 3.14)
	assert.Equal(t, typed.FloatOr("pi", 11.3), 3.14)
	value, exists := typed.FloatIf("pi")
	assert.Equal(t, value, 3.14)
	assert.True(t, exists)

	assert.Equal(t, typed.Float("string"), 30.14)
	assert.Equal(t, typed.FloatOr("string", 11.3), 30.14)
	value, exists = typed.FloatIf("string")
	assert.Equal(t, value, 30.14)
	assert.True(t, exists)

	assert.Equal(t, typed.Float("other"), 0.0)
	assert.Equal(t, typed.FloatOr("other", 11.3), 11.3)
	value, exists = typed.FloatIf("other")
	assert.Equal(t, value, 0.0)
	assert.False(t, exists)

	value, exists = typed.FloatIf("nope")
	assert.Equal(t, value, 0.0)
	assert.False(t, exists)

	assert.Equal(t, typed.FloatMust("pi"), 3.14)

	defer mustTest(t, "expected float value for fail")
	typed.FloatMust("fail")
	t.FailNow()
}

func Test_String(t *testing.T) {
	typed := New(build("host", "localhost", "nope", 1))
	assert.Equal(t, typed.String("host"), "localhost")
	assert.Equal(t, typed.StringOr("host", "openmymind.net"), "localhost")
	value, exists := typed.StringIf("host")
	assert.Equal(t, value, "localhost")
	assert.True(t, exists)

	assert.Equal(t, typed.String("other"), "")
	assert.Equal(t, typed.StringOr("other", "openmymind.net"), "openmymind.net")
	value, exists = typed.StringIf("other")
	assert.Equal(t, value, "")
	assert.False(t, exists)

	value, exists = typed.StringIf("nope")
	assert.Equal(t, value, "")
	assert.False(t, exists)

	assert.Equal(t, typed.StringMust("host"), "localhost")

	defer mustTest(t, "expected string value for fail")
	typed.StringMust("fail")
	t.FailNow()
}

func Test_Bytes(t *testing.T) {
	typed := New(build("host", "localhost", "nope", 1, "yes", []byte{9, 88}))
	assert.Bytes(t, typed.Bytes("host"), []byte("localhost"))
	assert.Bytes(t, typed.Bytes("yes"), []byte{9, 88})
	assert.Bytes(t, typed.BytesOr("host", []byte{1, 9}), []byte("localhost"))
	value, exists := typed.BytesIf("host")
	assert.Bytes(t, value, []byte("localhost"))
	assert.True(t, exists)

	assert.True(t, typed.Bytes("other") == nil)
	assert.Bytes(t, typed.BytesOr("other", []byte{1, 9}), []byte{1, 9})
	value, exists = typed.BytesIf("other")
	assert.True(t, value == nil)
	assert.False(t, exists)

	value, exists = typed.BytesIf("nope")
	assert.True(t, value == nil)
	assert.False(t, exists)

	assert.Bytes(t, typed.BytesMust("host"), []byte("localhost"))

	defer mustTest(t, "expected []byte value for fail")
	typed.BytesMust("fail")
	t.FailNow()
}

func Test_Object(t *testing.T) {
	typed := New(build("server", build("port", 32), "nope", "a"))
	assert.Equal(t, typed.Object("server").Int("port"), 32)
	assert.Equal(t, typed.ObjectOr("server", build("a", "b")).Int("port"), 32)

	assert.Equal(t, len(typed.Object("other")), 0)
	assert.Equal(t, typed.ObjectOr("other", build("x", "y")).String("x"), "y")

	value, exists := typed.ObjectIf("other")
	assert.Equal(t, len(value), 0)
	assert.False(t, exists)

	value, exists = typed.ObjectIf("nope")
	assert.Equal(t, len(value), 0)
	assert.False(t, exists)

	assert.Equal(t, typed.ObjectMust("server").Int("port"), 32)

	defer mustTest(t, "expected map for fail")
	typed.ObjectMust("fail")
	t.FailNow()
}

func Test_ObjectType(t *testing.T) {
	typed := New(build("server", Typed(build("port", 32))))
	assert.Equal(t, typed.Object("server").Int("port"), 32)
}

func Test_Any(t *testing.T) {
	typed := New(build("host", "localhost"))
	assert.Equal(t, typed.Any("host").(string), "localhost")
	assert.Equal(t, typed.AnyOr("host", "openmymind.net").(string), "localhost")
	value, exists := typed.AnyIf("host")
	assert.Equal(t, value.(string), "localhost")
	assert.True(t, exists)

	assert.Nil(t, typed.Any("other"))
	assert.Equal(t, typed.AnyOr("other", "openmymind.net").(string), "openmymind.net")
	value, exists = typed.AnyIf("other")
	assert.Nil(t, value)
	assert.False(t, exists)

	assert.Equal(t, typed.AnyMust("host").(string), "localhost")

	defer mustTest(t, "expected any for fail")
	typed.AnyMust("fail")
	t.FailNow()
}

func Test_Anys(t *testing.T) {
	typed := New(build("input", []any{"a", "b"}))
	a, exists := typed.AnysIf("x")
	assert.False(t, exists)
	assert.Equal(t, len(a), 0)

	a = typed.Anys("x")
	assert.Equal(t, len(a), 0)

	a, exists = typed.AnysIf("input")
	assert.True(t, exists)
	assert.Equal(t, a[0].(string), "a")
	assert.Equal(t, a[1].(string), "b")

	a = typed.Anys("input")
	assert.Equal(t, a[0].(string), "a")
	assert.Equal(t, a[1].(string), "b")

	a = typed.AnysOr("x", []any{"zz"})
	assert.Equal(t, a[0].(string), "zz")

	a = typed.AnysOr("input", []any{"zz"})
	assert.Equal(t, a[0].(string), "a")
	assert.Equal(t, a[1].(string), "b")

	defer mustTest(t, "expected array of anys for fail")
	typed.AnysMust("fail")
	t.FailNow()
}

func Test_Bools(t *testing.T) {
	typed := New(build("boring", []any{true, false}, "fail", []any{true, "goku"}, "bools", []bool{false, false, true}, "nope", 1))
	assert.List(t, typed.Bools("boring"), []bool{true, false})
	assert.Equal(t, len(typed.Bools("other")), 0)
	assert.List(t, typed.BoolsOr("boring", []bool{false, true}), []bool{true, false})
	assert.List(t, typed.BoolsOr("other", []bool{false, true}), []bool{false, true})
	assert.List(t, typed.Bools("bools"), []bool{false, false, true})

	values, exists := typed.BoolsIf("fail")
	assert.Equal(t, len(values), 0)
	assert.False(t, exists)

	values, exists = typed.BoolsIf("nope")
	assert.Equal(t, len(values), 0)
	assert.False(t, exists)
}

func Test_Ints(t *testing.T) {
	typed := New(build("scores", []any{2, 1, "3"}, "fail1", []any{2, "nope"}, "fail2", []any{2, true}, "ints", []int{9, 8}, "empty", []any{}))
	assert.List(t, typed.Ints("scores"), []int{2, 1, 3})
	assert.List(t, typed.Ints("ints"), []int{9, 8})
	assert.List(t, typed.Ints("empty"), []int{})
	assert.Equal(t, len(typed.Ints("other")), 0)
	assert.List(t, typed.IntsOr("scores", []int{3, 4, 5}), []int{2, 1, 3})
	assert.List(t, typed.IntsOr("other", []int{3, 4, 5}), []int{3, 4, 5})

	values, exists := typed.IntsIf("fail1")
	assert.Equal(t, len(values), 0)
	assert.False(t, exists)

	values, exists = typed.IntsIf("fail2")
	assert.Equal(t, len(values), 0)
	assert.False(t, exists)
}

func Test_Ints64(t *testing.T) {
	typed := New(build("scores", []any{2, 1, "3", 2.0, int64(939292992929292)}, "fail1", []any{2, "nope"}, "fail2", []any{2, true}, "fail3", "nope"))
	assert.List(t, typed.Ints64("scores"), []int64{2, 1, 3, 2.0, 939292992929292})
	assert.Equal(t, len(typed.Ints64("other")), 0)
	assert.List(t, typed.Ints64Or("scores", []int64{3, 4, 5}), []int64{2, 1, 3, 2.0, 939292992929292})
	assert.List(t, typed.Ints64Or("other", []int64{3, 4, 5}), []int64{3, 4, 5})

	values, exists := typed.Ints64If("fail1")
	assert.Equal(t, len(values), 0)
	assert.False(t, exists)

	values, exists = typed.Ints64If("fail2")
	assert.Equal(t, len(values), 0)
	assert.False(t, exists)

	values, exists = typed.Ints64If("fail3")
	assert.Equal(t, len(values), 0)
	assert.False(t, exists)
}

func Test_Ints_WithFloats(t *testing.T) {
	typed := New(build("scores", []any{2.1, 7.39}))
	assert.List(t, typed.Ints("scores"), []int{2, 7})
}

func Test_Floats(t *testing.T) {
	typed := New(build("ranks", []any{2.1, 1.2, "3.0"}, "fail1", []any{"a"}, "floats", []float64{9.0}))
	assert.List(t, typed.Floats("floats"), []float64{9.0})
	assert.List(t, typed.Floats("ranks"), []float64{2.1, 1.2, 3.0})
	assert.Equal(t, len(typed.Floats("other")), 0)
	assert.List(t, typed.FloatsOr("ranks", []float64{3.1, 4.2, 5.3}), []float64{2.1, 1.2, 3.0})
	assert.List(t, typed.FloatsOr("other", []float64{3.1, 4.2, 5.3}), []float64{3.1, 4.2, 5.3})

	values, exists := typed.FloatsIf("fail1")
	assert.Equal(t, len(values), 0)
	assert.False(t, exists)
}

func Test_Strings(t *testing.T) {
	typed := New(build("names", []any{"a", "b"}, "strings", []string{"s1", "s2"}, "fail1", []any{1, true}, "fail2", "2"))
	assert.List(t, typed.Strings("names"), []string{"a", "b"})
	assert.List(t, typed.Strings("strings"), []string{"s1", "s2"})
	assert.Equal(t, len(typed.Strings("other")), 0)
	assert.List(t, typed.StringsOr("names", []string{"c", "d"}), []string{"a", "b"})
	assert.List(t, typed.StringsOr("other", []string{"c", "d"}), []string{"c", "d"})

	values, exists := typed.StringsIf("fail1")
	assert.Equal(t, len(values), 0)
	assert.False(t, exists)

	values, exists = typed.StringsIf("fail2")
	assert.Equal(t, len(values), 0)
	assert.False(t, exists)
}

func Test_Objects(t *testing.T) {
	typed := New(build("names", []any{build("first", 1), build("second", 2)}))
	assert.Equal(t, typed.Objects("names")[0].Int("first"), 1)
}

func Test_ObjectsIf(t *testing.T) {
	typed := New(build("names", []any{build("first", 1), build("second", 2)}))
	objects, exists := typed.ObjectsIf("names")
	assert.Equal(t, objects[0].Int("first"), 1)
	assert.True(t, exists)

	objects, exists = typed.ObjectsIf("non_existing")
	assert.Equal(t, len(objects), 0)
	assert.False(t, exists)
}

func Test_ObjectsMust(t *testing.T) {
	typed := New(build("names", []any{build("first", 1), build("second", 2)}))
	objects := typed.ObjectsMust("names")
	assert.Equal(t, objects[0].Int("first"), 1)

	paniced := false
	func() {
		defer func() {
			if recovered := recover(); recovered != nil {
				paniced = true
			}
		}()

		typed.ObjectsMust("non_existing")
	}()

	assert.Equal(t, paniced, true)
}

func Test_ObjectsAsMap(t *testing.T) {
	typed := New(build("names", []map[string]any{build("first", 1), build("second", 2)}))
	assert.Equal(t, typed.Objects("names")[0].Int("first"), 1)
}

func Test_Maps(t *testing.T) {
	typed := New(build("names", []any{build("first", 1), build("second", 2)}))
	assert.Equal(t, typed.Maps("names")[1]["second"].(int), 2)
}

func Test_StringBool(t *testing.T) {
	typed, _ := JsonString(`{"blocked":{"a":true,"b":false}}`)
	m := typed.StringBool("blocked")
	assert.Equal(t, m["a"], true)
	assert.Equal(t, m["b"], false)

	m = typed.StringBool("other")
	assert.Equal(t, len(m), 0)
}

func Test_StringInt(t *testing.T) {
	typed, _ := JsonString(`{"count":{"a":123,"c":"55"}}`)
	m := typed.StringInt("count")
	assert.Equal(t, m["a"], 123)
	assert.Equal(t, m["c"], 55)
	assert.Equal(t, m["xxz"], 0)
	assert.Equal(t, len(typed.StringInt("nope")), 0)

	typed = New(build("count", map[string]any{"a": 99, "b": 8.0, "c": "9"}, "fail1", map[string]any{"a": "nope"}))
	m = typed.StringInt("count")
	assert.Equal(t, m["a"], 99)
	assert.Equal(t, m["b"], 8)
	assert.Equal(t, m["c"], 9)

	assert.Equal(t, len(typed.StringInt("fail")), 0)
}

func Test_StringFloat(t *testing.T) {
	typed, _ := JsonString(`{"rank":{"aa":3.4,"bz":4.2,"cc":"5.5"}}`)
	m := typed.StringFloat("rank")
	assert.Equal(t, m["aa"], 3.4)
	assert.Equal(t, m["bz"], 4.2)
	assert.Equal(t, m["cc"], 5.5)
	assert.Equal(t, m["xx"], 0.0)
	assert.Equal(t, len(typed.StringFloat("nope")), 0)

	typed = New(build("count", map[string]any{"a": 1.1, "b": "2.2"}, "fail1", map[string]any{"a": "nope"}))
	m = typed.StringFloat("count")
	assert.Equal(t, m["a"], 1.1)
	assert.Equal(t, m["b"], 2.2)

	assert.Equal(t, len(typed.StringFloat("fail")), 0)
}

func Test_StringString(t *testing.T) {
	typed, _ := JsonString(`{"atreides":{"leto":"ghanima","paul":"alia"}}`)
	m := typed.StringString("atreides")
	assert.Equal(t, m["leto"], "ghanima")
	assert.Equal(t, m["paul"], "alia")
	assert.Equal(t, m["vladimir"], "")

	m = typed.StringString("other")
	assert.Equal(t, len(m), 0)
}

func Test_StringObject(t *testing.T) {
	typed, _ := JsonString(`{"atreides":{"leto":{"sister": "ghanima"}, "goku": {"power": 9001}}}`)
	m := typed.StringObject("atreides")
	assert.Equal(t, m["leto"].String("sister"), "ghanima")
	assert.Equal(t, m["goku"].Int("power"), 9001)

	m = typed.StringObject("other")
	assert.Equal(t, len(m), 0)
}

func Test_Exists(t *testing.T) {
	typed := New(build("power", 9001))
	assert.True(t, typed.Exists("power"))
	assert.False(t, typed.Exists("spice"))
}

func Test_Map(t *testing.T) {
	typed := New(build("data", map[string]any{"a": 1}, "nope", 11))
	m := typed.Map("data")
	assert.Equal(t, m["a"].(int), 1)

	m = typed.MapOr("data", nil)
	assert.Equal(t, m["a"].(int), 1)

	m = typed.MapOr("other", map[string]any{"a": 2})
	assert.Equal(t, m["a"].(int), 2)

	m, exists := typed.MapIf("data")
	assert.Equal(t, m["a"].(int), 1)
	assert.True(t, exists)

	m, exists = typed.MapIf("other")
	assert.Equal(t, len(m), 0)
	assert.False(t, exists)

	m, exists = typed.MapIf("nope")
	assert.Equal(t, len(m), 0)
	assert.False(t, exists)
}

func Test_Time(t *testing.T) {
	zero := time.Time{}
	now := time.Now().UTC()
	typed := New(build("ts", now, "nope", true))
	assert.Equal(t, typed.Time("ts"), now)
	assert.Equal(t, typed.TimeOr("ts", zero), now)
	assert.Equal(t, typed.TimeOr("other", zero), zero)

	ts, exists := typed.TimeIf("ts")
	assert.Equal(t, ts, now)
	assert.True(t, exists)

	ts, exists = typed.TimeIf("other")
	assert.Equal(t, ts, zero)
	assert.False(t, exists)

	ts, exists = typed.TimeIf("nope")
	assert.Equal(t, ts, zero)
	assert.False(t, exists)

	assert.Equal(t, typed.TimeMust("ts"), now)

	defer mustTest(t, "expected time.Time value for other")
	typed.TimeMust("other")
	t.FailNow()
}

func Test_Time_String(t *testing.T) {
	zero := time.Time{}
	now := time.Now().UTC().Truncate(time.Second)
	typed := New(build("ts", now.Format(time.RFC3339), "nope", true))
	assert.Equal(t, typed.Time("ts"), now)
	assert.Equal(t, typed.TimeOr("ts", zero), now)
	assert.Equal(t, typed.TimeOr("other", zero), zero)

	ts, exists := typed.TimeIf("ts")
	assert.Equal(t, ts, now)
	assert.True(t, exists)

	ts, exists = typed.TimeIf("other")
	assert.Equal(t, ts, zero)
	assert.False(t, exists)

	ts, exists = typed.TimeIf("nope")
	assert.Equal(t, ts, zero)
	assert.False(t, exists)

	assert.Equal(t, typed.TimeMust("ts"), now)

	defer mustTest(t, "expected time.Time value for other")
	typed.TimeMust("other")
	t.FailNow()
}

func Test_Time_Int(t *testing.T) {
	zero := time.Time{}
	now := time.Now().UTC().Truncate(time.Second)
	typed := New(build("ts", now.Unix(), "nope", true))
	assert.Equal(t, typed.Time("ts"), now)
	assert.Equal(t, typed.TimeOr("ts", zero), now)
	assert.Equal(t, typed.TimeOr("other", zero), zero)

	ts, exists := typed.TimeIf("ts")
	assert.Equal(t, ts, now)
	assert.True(t, exists)

	ts, exists = typed.TimeIf("other")
	assert.Equal(t, ts, zero)
	assert.False(t, exists)

	ts, exists = typed.TimeIf("nope")
	assert.Equal(t, ts, zero)
	assert.False(t, exists)

	assert.Equal(t, typed.TimeMust("ts"), now)

	defer mustTest(t, "expected time.Time value for other")
	typed.TimeMust("other")
	t.FailNow()
}

func Test_JsonKey(t *testing.T) {
	typed := Typed(map[string]any{
		"string":  `{"over":9000}`,
		"bytes":   []byte(`{"leto":"atreides"}`),
		"invalid": `"a`,
	})

	typed2, err := typed.Json("unknown")
	assert.Nil(t, err)
	assert.Equal(t, len(typed2), 0)

	typed2, err = typed.Json("string")
	assert.Nil(t, err)
	assert.Equal(t, typed2.Int("over"), 9000)

	typed2, err = typed.Json("bytes")
	assert.Nil(t, err)
	assert.Equal(t, typed2.String("leto"), "atreides")

	typed2, err = typed.Json("invalid")
	assert.Equal(t, err.Error(), `expected { character for map value`)
	assert.Equal(t, len(typed2), 0)

	// MUST
	typed2 = typed.JsonMust("unknown")
	assert.Equal(t, len(typed2), 0)

	typed2 = typed.JsonMust("string")
	assert.Equal(t, typed2.Int("over"), 9000)

	typed2 = typed.JsonMust("bytes")
	assert.Equal(t, typed2.String("leto"), "atreides")

	defer mustTest(t, "expected { character for map value")
	typed2 = typed.JsonMust("invalid")
	assert.Equal(t, len(typed2), 0)

}

func build(values ...any) map[string]any {
	m := make(map[string]any, len(values))
	for i := 0; i < len(values); i += 2 {
		m[values[i].(string)] = values[i+1]
	}
	return m
}

func mustTest(t *testing.T, expected string) {
	t.Helper()
	if err := recover(); err != nil {
		switch e := err.(type) {
		case string:
			assert.Equal(t, e, expected)
		case error:
			assert.Equal(t, e.Error(), expected)
		default:
			panic("unknown recover type")
		}
	}
}
