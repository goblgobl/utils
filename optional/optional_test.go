package optional

import (
	"testing"

	"src.goblgobl.com/tests/assert"
)

func Test_Null(t *testing.T) {
	assert.False(t, NullInt.Exists)
	assert.False(t, Null[int]().Exists)
}

func Test_Int(t *testing.T) {
	v := Int(9001)
	assert.True(t, v.Exists)
	assert.Equal(t, v.Value, 9001)
}

func Test_New(t *testing.T) {
	v := New("over")
	assert.True(t, v.Exists)
	assert.Equal(t, v.Value, "over")
}
