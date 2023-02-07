package ascii

import (
	"testing"

	"src.goblgobl.com/tests/assert"
)

func Test_HasPrefixIgnoreCase(t *testing.T) {
	assert.False(t, HasPrefixIgnoreCase("", "a"))
	assert.False(t, HasPrefixIgnoreCase("", "ab"))
	assert.False(t, HasPrefixIgnoreCase("", " "))
	assert.False(t, HasPrefixIgnoreCase("h", "a"))
	assert.False(t, HasPrefixIgnoreCase("h", " h"))
	assert.False(t, HasPrefixIgnoreCase("h", "h "))
	assert.False(t, HasPrefixIgnoreCase("he", "ha"))
	assert.False(t, HasPrefixIgnoreCase("he", "eh"))
	assert.False(t, HasPrefixIgnoreCase("he", " he"))
	assert.False(t, HasPrefixIgnoreCase("he", " he "))
	assert.True(t, HasPrefixIgnoreCase("", ""))
	assert.True(t, HasPrefixIgnoreCase("a", "a"))
	assert.True(t, HasPrefixIgnoreCase("a", "A"))
	assert.True(t, HasPrefixIgnoreCase("A", "a"))
	assert.True(t, HasPrefixIgnoreCase("A", "A"))
	assert.True(t, HasPrefixIgnoreCase("Abc", "A"))
	assert.True(t, HasPrefixIgnoreCase("A ", "A"))
	assert.True(t, HasPrefixIgnoreCase("abc", "abc"))
	assert.True(t, HasPrefixIgnoreCase("abc213", "Abc"))
	assert.True(t, HasPrefixIgnoreCase("abc", "ABc"))
	assert.True(t, HasPrefixIgnoreCase("abc", "ABC"))
	assert.True(t, HasPrefixIgnoreCase("Abc", "abc"))
	assert.True(t, HasPrefixIgnoreCase("ABc", "abc"))
	assert.True(t, HasPrefixIgnoreCase("ABC", "abc"))
	assert.True(t, HasPrefixIgnoreCase("ABC", "ABC"))
	assert.True(t, HasPrefixIgnoreCase("aBc", "AbC"))
	assert.True(t, HasPrefixIgnoreCase("aBc", "AbC"))
	assert.True(t, HasPrefixIgnoreCase("aBc_123", "abc_123"))
	assert.True(t, HasPrefixIgnoreCase("abc_12345", "ABC_123"))
}
