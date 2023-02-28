package concurrent

/*
sync.Pool is impressive. But, sync.Pool is for "temporary objects" and:

	"Any item stored in the Pool may be removed automatically at any time
	without notification"

The details aren't part of the contract, but sync.Pool currently removes some
(1/2?) of all pooled items from the pool during a GC. The result is that using
sync.Pool can cause unpredictable performance, allocations and, as a result,
latency. It's the type of performance issue that might not show up in a micro-
benchark when there's little memory pressure and thus little GC.

This pool structure uses a buffered channel. The buffered channel guarantees
that a minimum number of items are kept (assuming consumers release items back
to the pool). It's pretty typical stuff. To spice it up, we have 2 modifications:

1 - Multiple channels are actually used, using a `[8]chan T`. An atomic counter
is used on checkout to pick a channel (e.g. `channels[atomic.Add(1) & 7]`).
Items are always returned to the same channel they came from.

2 - When the channels are empty, we fallback to a sync.Pool. This allows the
pool to efficiently deal with spikes or being under-provisioned. Thus setting a
conservative minimum is an acceptable solution. This also leverages the
temporary nature of sync.Pool: as load decreases, additionally created items will
be GC'd (up to the configured minimum).

In a microbenchmark that focuses on the pool's performance, this performs much
closer to a single channel than sync.Pool. Plus in an extreme case where items
are checked out from the pool for inconsistent periods of time, it's possible
for 1 channel to become empty while others are full, which would result in falling
back to sync.Pool and more allocations.
*/

import (
	"sync"
	"sync/atomic"
)

type Pool[T any] struct {
	fallback *sync.Pool
	buckets  [8]chan T
	createT  func(release func(T)) T
	index    atomic.Uint64
	depleted atomic.Uint64
}

func NewPool[T any](min uint32, createT func(release func(T)) T) *Pool[T] {
	if min < 8 {
		min = 8
	}

	minPerBucket := min / 8
	p := &Pool[T]{createT: createT}

	for i := uint32(0); i < 8; i++ {
		releaser := p.makeReleaser(i)
		bucket := make(chan T, minPerBucket)
		for j := uint32(0); j < minPerBucket; j++ {
			bucket <- createT(releaser)
		}
		p.buckets[i] = bucket
	}

	p.fallback = &sync.Pool{
		New: func() any {
			return createT(p.releaseToFallback)
		},
	}

	return p
}

func (p *Pool[T]) Checkout() T {
	index := p.index.Add(1) & 7
	select {
	case t := <-p.buckets[index]:
		return t
	default:
		p.depleted.Add(1)
		return p.fallback.Get().(T)
	}
}

func (p *Pool[T]) makeReleaser(bucketIndex uint32) func(t T) {
	return func(t T) {
		p.buckets[bucketIndex] <- t
	}
}

func (p *Pool[T]) releaseToFallback(t T) {
	p.fallback.Put(t)
}

// how often the pool was empty and we had to create a buffer on the fly
func (p *Pool[T]) Depleted() uint64 {
	return p.depleted.Load()
}

func (p *Pool[T]) Len() int {
	buckets := p.buckets
	return len(buckets[0]) +
		len(buckets[1]) +
		len(buckets[2]) +
		len(buckets[3]) +
		len(buckets[4]) +
		len(buckets[5]) +
		len(buckets[6]) +
		len(buckets[7])
}
