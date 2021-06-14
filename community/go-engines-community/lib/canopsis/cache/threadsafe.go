package cache

import (
	"context"
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
func (tsc *ThreadSafeCache) Get(ctx context.Context, id string, out interface{}) bool {
	tsc.m.RLock()
	defer tsc.m.RUnlock()
	return tsc.c.Get(ctx, id, out)
}

// Set ...
func (tsc *ThreadSafeCache) Set(ctx context.Context, element Cacheable) error {
	tsc.m.Lock()
	defer tsc.m.Unlock()
	return tsc.c.Set(ctx, element)
}

// SetRaw ...
func (tsc *ThreadSafeCache) SetRaw(ctx context.Context, id string, element interface{}) error {
	tsc.m.Lock()
	defer tsc.m.Unlock()
	return tsc.SetRaw(ctx, id, element)
}

// Drop ...
func (tsc *ThreadSafeCache) Drop(ctx context.Context, ids ...string) error {
	tsc.m.Lock()
	defer tsc.m.Unlock()
	return tsc.c.Drop(ctx, ids...)
}

// Flush ...
func (tsc *ThreadSafeCache) Flush(ctx context.Context) error {
	tsc.m.Lock()
	defer tsc.m.Unlock()
	return tsc.c.Flush(ctx)
}
