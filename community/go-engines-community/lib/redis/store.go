package redis

//go:generate mockgen -destination=../../mocks/lib/redis/store.go git.canopsis.net/canopsis/go-engines/lib/redis Store

import (
	"encoding/json"
	"github.com/go-redis/redis/v7"
	"time"
)

// Store interface is used to implement baseStore any variable to redis.
// Basic implementation marshals data to json before stores it to cache.
type Store interface {
	// Save prepares value and saves the result to redis.
	Save(v interface{}) error
	// Restore gets data from redis and stores the result
	// in the value pointed to by v.
	Restore(v interface{}) (bool, error)
}

// baseStore saves data with provided key and expiration.
// baseStore uses json format.
type baseStore struct {
	client     redis.Cmdable
	key        string
	expiration time.Duration
}

// NewStore creates new store.
func NewStore(client redis.Cmdable, key string, expiration time.Duration) Store {
	return &baseStore{
		client:     client,
		key:        key,
		expiration: expiration,
	}
}

func (s *baseStore) Save(v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	res := s.client.Set(s.key, b, s.expiration)
	if err := res.Err(); err != nil {
		return err
	}

	return nil
}

func (s *baseStore) Restore(v interface{}) (bool, error) {
	res := s.client.Get(s.key)
	if err := res.Err(); err != nil {
		if err == redis.Nil {
			return false, nil
		}

		return false, err
	}

	err := json.Unmarshal([]byte(res.Val()), v)
	if err != nil {
		return false, err
	}

	return true, nil
}
