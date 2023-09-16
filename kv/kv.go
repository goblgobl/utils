package kv

import "src.goblgobl.com/utils/optional"

// Poolable key-value optimized for small sets

type KV[K comparable, V any] struct {
	release func(*KV[K, V])
	Len     int
	Keys    []K
	Values  []V
}

func New[K comparable, V any](max uint16) KV[K, V] {
	return KV[K, V]{
		Len:    0,
		Keys:   make([]K, max),
		Values: make([]V, max),
	}
}

// Adds the key=>value. If key already exists, a duplicate is added
func (kv *KV[K, V]) Add(key K, value V) bool {
	l := kv.Len
	if l == len(kv.Keys) {
		return false
	}
	kv.Keys[l] = key
	kv.Values[l] = value
	kv.Len = l + 1
	return true
}

// Adds the key=>value. If key already exists, a duplicate is added
func (kv *KV[K, V]) Put(key K, value V) bool {
	for i, k := range kv.Keys {
		if k == key {
			kv.Values[i] = value
			return true
		}
	}
	return kv.Add(key, value)
}

func (kv *KV[K, V]) Get(key K) optional.Value[V] {
	for i, k := range kv.Keys {
		if k == key {
			return optional.New(kv.Values[i])
		}
	}
	return optional.Null[V]()
}

func (kv *KV[K, V]) Del(key K) optional.Value[V] {
	for i, k := range kv.Keys {
		if k == key {
			l := kv.Len - 1
			value := kv.Values[i]
			if l > 0 {
				// move the last item into this slot
				kv.Keys[i] = kv.Keys[l]
				kv.Values[i] = kv.Values[l]
			}
			kv.Len = l
			return optional.New(value)
		}
	}
	return optional.Null[V]()
}

func (kv *KV[K, V]) Release() {
	if release := kv.release; release != nil {
		kv.Len = 0
		release(kv)
	}
}
