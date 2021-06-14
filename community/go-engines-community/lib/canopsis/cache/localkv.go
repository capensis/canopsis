package cache

import "reflect"

// KVCache is a simple map[string]interface{} so you get a very simple
// local key:value storage.
type KVCache struct {
	cache map[string]interface{}
}

// NewKV is a local Key: Value storage using a stock golang map[string]interface{}.
// Very useful for tests but could be used for any local cache.
func NewKV() Cache {
	c := KVCache{
		cache: make(map[string]interface{}),
	}
	return &c
}

// Store returns the underlying cache map. Useful for tests.
func (c *KVCache) Store() map[string]interface{} {
	return c.cache
}

// Reset allocates a new cache map. Useful for tests.
func (c *KVCache) Reset() {
	c.cache = make(map[string]interface{})
}

// Get ...
func (c *KVCache) Get(id string, out interface{}) bool {
	value, exists := c.cache[id]
	if exists {
		v := reflect.ValueOf(out).Elem()
		if v.CanSet() {
			v.Set(reflect.ValueOf(value))
		} else {
			return false
		}
	}
	return exists
}

// Set ...
func (c *KVCache) Set(element Cacheable) error {
	return c.SetRaw(element.CacheID(), element)
}

// SetRaw ...
func (c *KVCache) SetRaw(id string, element interface{}) error {
	c.cache[id] = element
	return nil
}

// Drop ...
func (c *KVCache) Drop(ids ...string) error {
	for _, id := range ids {
		delete(c.cache, id)
	}
	return nil
}

// Flush ...
func (c *KVCache) Flush() error {
	c.cache = make(map[string]interface{})
	return nil
}
