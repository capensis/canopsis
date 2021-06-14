package cache

import "context"

// Cacheable makes your struct cache-able
type Cacheable interface {
	CacheID() string
}

// Cache unifies cache usage across the project.
type Cache interface {
	// GetLocal returns the object in local cache, if it exists.
	// Returns true if found, false if not found
	Get(ctx context.Context, id string, out interface{}) bool

	// Set both remote and local caches
	Set(ctx context.Context, element Cacheable) error

	// SetRaw lets you choose the id instead of relying on the Cacheable interface
	SetRaw(ctx context.Context, id string, element interface{}) error

	// Drop both remote and local cache entry
	Drop(ctx context.Context, ids ...string) error

	// Flush removes all entries in remote and local cache
	Flush(ctx context.Context) error
}
