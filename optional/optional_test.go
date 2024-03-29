package optional

import (
	"testing"

	"src.goblgobl.com/utils/json"

	"src.goblgobl.com/tests/assert"
)

func Test_Null(t *testing.T) {
	assert.False(t, NullInt.Exists)
	assert.False(t, Null[int]().Exists)
}

func Test_Int(t *testing.T) {
	v := NewInt(9001)
	assert.True(t, v.Exists)
	assert.Equal(t, v.Value, 9001)
}

func Test_Bool(t *testing.T) {
	v := NewBool(true)
	assert.True(t, v.Exists)
	assert.Equal(t, v.Value, true)
}

func Test_String(t *testing.T) {
	v := NewString("teg atreides")
	assert.True(t, v.Exists)
	assert.Equal(t, v.Value, "teg atreides")
}

func Test_Float(t *testing.T) {
	v := NewFloat(9000.2)
	assert.True(t, v.Exists)
	assert.Equal(t, v.Value, 9000.2)
}

func Test_New(t *testing.T) {
	v := New("over")
	assert.True(t, v.Exists)
	assert.Equal(t, v.Value, "over")
}

func Test_Int_Json(t *testing.T) {
	power := NewInt(9001)
	data, err := json.Marshal(power)
	assert.Nil(t, err)
	assert.Equal(t, string(data), `9001`)

	power = NullInt
	data, err = json.Marshal(power)
	assert.Nil(t, err)
	assert.Equal(t, string(data), `null`)

	var p1 Int
	assert.Nil(t, json.Unmarshal([]byte(`9002`), &p1))
	assert.True(t, p1.Exists)
	assert.Equal(t, p1.Value, 9002)

	var p2 Value[int]
	assert.Nil(t, json.Unmarshal([]byte(`null`), &p2))
	assert.False(t, p2.Exists)
}

func Test_String_Json(t *testing.T) {
	power := NewString("over")
	data, err := json.Marshal(power)
	assert.Nil(t, err)
	assert.Equal(t, string(data), `"over"`)

	power = NullString
	data, err = json.Marshal(power)
	assert.Nil(t, err)
	assert.Equal(t, string(data), `null`)

	var p1 String
	assert.Nil(t, json.Unmarshal([]byte(`"ninet"`), &p1))
	assert.True(t, p1.Exists)
	assert.Equal(t, p1.Value, "ninet")

	var p2 Value[string]
	assert.Nil(t, json.Unmarshal([]byte(`null`), &p2))
	assert.False(t, p2.Exists)
}

func Test_FromAny(t *testing.T) {
	assert.False(t, FromAny[int](nil).Exists)
	assert.False(t, FromAny[int]("nope").Exists)
	assert.True(t, FromAny[int](32).Exists)
	assert.Equal(t, FromAny[int](32).Value, 32)
}
