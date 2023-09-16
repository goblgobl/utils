package kv

import (
	"testing"

	"src.goblgobl.com/tests/assert"
)

func Test_Add_And_Put_Get(t *testing.T) {
	kv := New[string, int](3)
	assert.Equal(t, kv.Len, 0)

	assert.Equal(t, true, kv.Add("Leto", 1))
	assertList(t, kv, "Leto", 1)

	assert.Equal(t, true, kv.Add("ghanima", 2))
	assertList(t, kv, "Leto", 1, "ghanima", 2)

	assert.Equal(t, true, kv.Put("Jessica", 3))
	assertList(t, kv, "Leto", 1, "ghanima", 2, "Jessica", 3)

	assert.Equal(t, true, kv.Put("ghanima", 20))
	assertList(t, kv, "Leto", 1, "ghanima", 20, "Jessica", 3)

	assert.Equal(t, true, kv.Put("Leto", 10))
	assertList(t, kv, "Leto", 10, "ghanima", 20, "Jessica", 3)

	assert.Equal(t, true, kv.Put("Jessica", 30))
	assertList(t, kv, "Leto", 10, "ghanima", 20, "Jessica", 30)

	// no more space
	assert.Equal(t, false, kv.Add("Paul", 4))
	assertList(t, kv, "Leto", 10, "ghanima", 20, "Jessica", 30)

	// no more space
	assert.Equal(t, false, kv.Put("Paul", 4))
	assertList(t, kv, "Leto", 10, "ghanima", 20, "Jessica", 30)
}

func Test_Del(t *testing.T) {
	kv := New[string, int](3)

	assert.Equal(t, false, kv.Del("Leto").Exists)

	{
		kv.Add("Leto", 100)
		assert.Equal(t, false, kv.Del("Paul").Exists)
		assertList(t, kv, "Leto", 100)
	}

	{
		d := kv.Del("Leto")
		assert.Equal(t, true, d.Exists)
		assert.Equal(t, 100, d.Value)
		assert.Equal(t, kv.Len, 0)
	}

	{
		kv.Add("A", 100)
		kv.Add("B", 200)
		kv.Add("C", 300)
		assert.Equal(t, 200, kv.Del("B").Value)
		assertList(t, kv, "A", 100, "C", 300)

		assert.Equal(t, 300, kv.Del("C").Value)
		assertList(t, kv, "A", 100)

		assert.Equal(t, 100, kv.Del("A").Value)
		assert.Equal(t, kv.Len, 0)
	}

	{
		kv.Add("A", 100)
		kv.Add("B", 200)
		kv.Add("C", 300)
		assert.Equal(t, 100, kv.Del("A").Value)
		assertList(t, kv, "B", 200, "C", 300)

		assert.Equal(t, 200, kv.Del("B").Value)
		assertList(t, kv, "C", 300)

		assert.Equal(t, 300, kv.Del("C").Value)
		assert.Equal(t, kv.Len, 0)
	}
}

func assertList(t *testing.T, kv KV[string, int], expected ...any) {
	assert.Equal(t, kv.Len, len(expected)/2)

	for i := 0; i < len(expected); i += 2 {
		expectedKey := expected[i].(string)
		expectedValue := expected[i+1].(int)
		assert.Equal(t, kv.Get(expectedKey).Value, expectedValue)
	}
}
