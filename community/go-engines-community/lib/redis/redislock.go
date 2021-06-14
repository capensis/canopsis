package redis

import (
	"github.com/bsm/redislock"
	"time"
)

type LockClient interface {
	Obtain(key string, ttl time.Duration, opt *redislock.Options) (Lock, error)
}

type Lock interface {
	Key() string
	Token() string
	TTL() (time.Duration, error)
	Refresh(ttl time.Duration, opt *redislock.Options) error
	Release() error
}

type lockClient struct {
	client *redislock.Client
}

func NewLockClient(redisClient redislock.RedisClient) LockClient {
	return &lockClient{
		client: redislock.New(redisClient),
	}
}

func (c *lockClient) Obtain(key string, ttl time.Duration, opt *redislock.Options) (Lock, error) {
	l, err := c.client.Obtain(key, ttl, opt)
	if err != nil {
		return nil, err
	}

	return l, nil
}
