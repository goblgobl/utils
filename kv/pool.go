package kv

import (
	"src.goblgobl.com/utils/concurrent"
)

type Config struct {
	Count uint16 `json:"count"`
	Max   uint16 `json:"max"`
}

type Pool[K comparable, V any] struct {
	max uint16
	*concurrent.Pool[*KV[K, V]]
}

func NewPoolFromConfig[K comparable, V any](config Config) Pool[K, V] {
	return NewPool[K, V](config.Count, config.Max)
}

func NewPool[K comparable, V any](count uint16, max uint16) Pool[K, V] {
	return Pool[K, V]{
		max:  max,
		Pool: concurrent.NewPool[*KV[K, V]](uint32(count), pooledKVFactory[K, V](max)),
	}
}

func pooledKVFactory[K comparable, V any](max uint16) func(func(kv *KV[K, V])) *KV[K, V] {
	return func(release func(kv *KV[K, V])) *KV[K, V] {
		kv := New[K, V](max)
		kv.release = release
		return &kv
	}
}

func (p Pool[K, V]) Checkout() *KV[K, V] {
	return p.Pool.Checkout()
}
