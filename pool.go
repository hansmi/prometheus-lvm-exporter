package main

import (
	"sync"
)

var defaultStringSlicePool = newStringSlicePool()

type stringSlicePool struct {
	pool sync.Pool
}

func newStringSlicePool() *stringSlicePool {
	return &stringSlicePool{
		pool: sync.Pool{
			New: func() any {
				return []string(nil)
			},
		},
	}
}

// get returns an empty string slice. The caller has ownership of the slice
// until the slice is put back into the pool.
func (p *stringSlicePool) get() []string {
	return p.pool.Get().([]string)[:0]
}

func (p *stringSlicePool) put(s []string) {
	// All elements must be accessible.
	s = s[:cap(s)]

	// Remove references to values.
	clear(s)

	p.pool.Put(s)
}
