package cache

import (
	"sync"
)

// ThreadSafeCache wrap Cache with locks
type ThreadSafeCache struct {
	c Cache
	m sync.RWMutex
}

// NewThreadSafeCache creates the wrapped cache
func NewThreadSafeCache(c Cache) Cache {
	return &ThreadSafeCache{
		c: c,
	}
}

// Get ...
func (tsc *ThreadSafeCache) Get(id string, out interface{}) bool {
	tsc.m.RLock()
	defer tsc.m.RUnlock()
	return tsc.c.Get(id, out)
}

// Set ...
func (tsc *ThreadSafeCache) Set(element Cacheable) error {
	tsc.m.Lock()
	defer tsc.m.Unlock()
	return tsc.c.Set(element)
}

// SetRaw ...
func (tsc *ThreadSafeCache) SetRaw(id string, element interface{}) error {
	tsc.m.Lock()
	defer tsc.m.Unlock()
	return tsc.SetRaw(id, element)
}

// Drop ...
func (tsc *ThreadSafeCache) Drop(ids ...string) error {
	tsc.m.Lock()
	defer tsc.m.Unlock()
	return tsc.c.Drop(ids...)
}

// Flush ...
func (tsc *ThreadSafeCache) Flush() error {
	tsc.m.Lock()
	defer tsc.m.Unlock()
	return tsc.c.Flush()
}
