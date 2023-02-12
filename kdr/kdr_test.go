package kdr

import (
	"testing"

	"src.goblgobl.com/tests/assert"
)

func Test_Keep(t *testing.T) {
	var change Value[string]
	change = Keep[string]()
	assert.Equal(t, change.Action, KDR_ACTION_KEEP)
	assert.True(t, change.IsKeep())
	assert.False(t, change.IsDelete())
	assert.False(t, change.IsReplace())
}

func Test_Delete(t *testing.T) {
	var change Value[string]
	change = Delete[string]()
	assert.Equal(t, change.Action, KDR_ACTION_DELETE)
	assert.False(t, change.IsKeep())
	assert.True(t, change.IsDelete())
	assert.False(t, change.IsReplace())
}

func Test_Replace(t *testing.T) {
	change := Replace("new value")
	assert.Equal(t, change.Action, KDR_ACTION_REPLACE)
	assert.False(t, change.IsKeep())
	assert.False(t, change.IsDelete())
	assert.True(t, change.IsReplace())
	assert.Equal(t, change.Replacement, "new value")
}
