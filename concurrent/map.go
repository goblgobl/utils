package concurrent

/*
A concurrent-safe string => generic map
*/

import (
	"sync"

	"golang.org/x/sync/singleflight"
)

type Loader[V any] func(id string) (V, error)

func NewMap[V any](loader Loader[V], cleaner func(v V)) Map[V] {
	shards := make([]*shard[V], 64)
	for i := 0; i < len(shards); i++ {
		shards[i] = &shard[V]{
			loader:  loader,
			cleaner: cleaner,
			lookup:  make(map[string]V),
			sf:      new(singleflight.Group),
		}
	}

	return Map[V]{shards}
}

type Map[V any] struct {
	shards []*shard[V]
}

func (m Map[V]) Get(id string) (V, error) {
	return m.shard(id).get(id)
}

func (m Map[V]) Put(id string, value V) {
	m.shard(id).put(id, value)
}

func (m Map[V]) shard(id string) *shard[V] {
	var h uint32
	for i := 0; i < len(id); i++ {
		h ^= uint32(id[i])
		h *= 16777619
	}
	return m.shards[h&63]
}

type shard[V any] struct {
	sf      *singleflight.Group
	cleaner func(v V)
	loader  Loader[V]
	lookup  map[string]V
	sync.RWMutex
}

func (s *shard[V]) get(id string) (V, error) {
	s.RLock()
	value, exists := s.lookup[id]
	s.RUnlock()

	if exists {
		return value, nil
	}

	ivalue, err, _ := s.sf.Do(id, func() (interface{}, error) {
		value, err := s.loader(id)
		if err != nil {
			var dflt V
			return dflt, err
		}
		s.Lock()
		s.lookup[id] = value
		s.Unlock()
		return value, nil
	})

	if err != nil {
		var dflt V
		return dflt, err
	}

	return ivalue.(V), nil
}

func (s *shard[V]) put(id string, value V) {
	var existing V
	var existed bool
	cleaner := s.cleaner

	s.Lock()
	lookup := s.lookup
	if cleaner != nil {
		existing, existed = lookup[id]
	}
	lookup[id] = value
	s.Unlock()

	if existed {
		cleaner(existing)
	}
}
