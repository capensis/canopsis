package entity

import (
	"fmt"
	"time"

	gogob "encoding/gob"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/cache"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/encoding/gob"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/utils"

	"github.com/go-redis/redis/v7"
	"github.com/rs/zerolog"
)

// Cache ...
type Cache struct {
	backend *redis.Client
	encoder encoding.Encoder
	decoder encoding.Decoder
	logger  zerolog.Logger
}

// NewCache creates the local cache and remote cache access for entitys.
func NewCache(redis *redis.Client, logger zerolog.Logger) cache.Cache {
	if redis == nil {
		utils.FailOnError(fmt.Errorf("is nil"), "entity cache redis")
	}

	gogob.Register(types.Entity{})

	c := Cache{
		backend: redis,
		encoder: gob.NewEncoder(),
		decoder: gob.NewDecoder(),
		logger:  logger,
	}

	return &c
}

// Drop entities in local and remote cache. Use Entity.ID for ids.
func (c Cache) Drop(ids ...string) error {
	return c.backend.Del(ids...).Err()
}

// Get entity from cache if it exists. Returns the entity, true or entity, false
func (c *Cache) Get(id string, out interface{}) bool {
	r := c.backend.Get(id)

	if r.Err() == redis.Nil {
		return false
	}

	b, _ := r.Bytes()

	err := c.decoder.Decode(b, &out)
	if err != nil {
		switch err.(type) {
		case encoding.DecodingError:
			c.Drop(id)
			return false
		default:
			return false
		}
	}

	entity, ok := out.(*types.Entity)
	if !ok {
		return false
	}

	entity.EnsureInitialized()

	return true
}

// Flush entity cache
func (c Cache) Flush() error {
	return c.backend.FlushDB().Err()
}

// Set the local and remote cache with this entity. If entity is nil, only
// the local cache is set with this value.
func (c Cache) Set(entity cache.Cacheable) error {
	if entity != nil {
		b, err := c.encoder.Encode(entity)
		if err != nil {
			switch err.(type) {
			case encoding.DecodingError:
				return encoding.NewDecodingError(fmt.Errorf("entity cache set: %v", err))
			default:
				return fmt.Errorf("entity cache unexpected: %v", err)
			}
		}
		return c.backend.Set(entity.CacheID(), b, time.Hour*4).Err()
	}
	return nil
}

// SetRaw is not implemented. Will panic.
func (c *Cache) SetRaw(id string, entity interface{}) error {
	panic("not implemented")
}
