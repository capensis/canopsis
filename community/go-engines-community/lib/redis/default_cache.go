package redis

import (
	"fmt"
	"time"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/cache"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/encoding/gob"
	mod "github.com/go-redis/redis/v7"
)

type dcache struct {
	backend *mod.Client
	encoder encoding.Encoder
	decoder encoding.Decoder
	timeout time.Duration
}

// Get data from cache if it exists. Returns true if so, else false.
// out param must be a reference to you concrete data.
func (c *dcache) Get(id string, out interface{}) bool {
	r := c.backend.Get(id)

	if r.Err() == mod.Nil {
		return false
	}

	b, _ := r.Bytes()

	err := c.decoder.Decode(b, out)
	if err != nil {
		switch err.(type) {
		case encoding.DecodingError:
			c.Drop(id)
			return false
		default:
			return false
		}
	}

	return true
}

// Flush data cache
func (c *dcache) Flush() error {
	return c.backend.FlushDB().Err()
}

// Drop data matching given ids
func (c *dcache) Drop(ids ...string) error {
	return c.backend.Del(ids...).Err()
}

// Set data in cache
func (c *dcache) Set(data cache.Cacheable) error {
	if data != nil {
		return c.SetRaw(data.CacheID(), data)
	}
	return nil
}

// SetRaw does the Redis magic.
func (c *dcache) SetRaw(id string, data interface{}) error {
	b, err := c.encoder.Encode(data)
	if err != nil {
		switch err.(type) {
		case encoding.DecodingError:
			return encoding.NewDecodingError(fmt.Errorf("data cache set: %v", err))
		default:
			return fmt.Errorf("data cache unexpected: %v", err)
		}
	}

	return c.backend.Set(id, b, c.timeout).Err()
}

// NewDefaultCache returns a default implementation of a redis cache.
// This default cache behaves like this:
// Get() tries to decode any encoded data with GOB from redis
// Set() tries to encode any data with GOB and set to redis, if data != nil
// Other functions just work as expected.
func NewDefaultCache(client *mod.Client, timeout time.Duration) (cache.Cache, error) {
	if client == nil {
		return nil, fmt.Errorf("redis default cache: nil client: unsupported")
	}

	return &dcache{
		backend: client,
		encoder: gob.NewEncoder(),
		decoder: gob.NewDecoder(),
		timeout: timeout,
	}, nil
}
