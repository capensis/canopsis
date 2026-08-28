// Package suspendstorage provides a generic redis-backed store for events whose
// processing was parked while their external data is fetched asynchronously through a webhook RPC.
// Entries are keyed by the webhook execution id, written
// with set-if-absent semantics and expire after the configured TTL, so an event
// whose RPC never answers is eventually discarded instead of leaked.
package suspendstorage

import (
	"context"
	"errors"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"github.com/redis/go-redis/v9"
)

const suspendedEventTTL = 24 * time.Hour

// ErrNotFound is returned by Get when no event is parked under the given id.
var ErrNotFound = errors.New("suspended event not found")

// Storage keeps parked events of type T, keyed by id.
type Storage[T any] interface {
	// Add parks v under id. It uses set-if-absent semantics, so an existing
	// entry for id is left untouched rather than overwritten.
	Add(ctx context.Context, id string, v T) error
	// Get returns the event parked under id, or an error when no entry exists.
	Get(ctx context.Context, id string) (T, error)
	// Delete removes the event parked under id. Deleting a missing entry is not an error.
	Delete(ctx context.Context, id string) error
}

func New[T any](
	key string,
	client redis.Cmdable,
	encoder encoding.Encoder,
	decoder encoding.Decoder,
) Storage[T] {
	return &storage[T]{
		key:     key,
		client:  client,
		encoder: encoder,
		decoder: decoder,
	}
}

type storage[T any] struct {
	key     string
	encoder encoding.Encoder
	decoder encoding.Decoder
	client  redis.Cmdable
}

func (s *storage[T]) Add(ctx context.Context, id string, v T) error {
	b, err := s.encoder.Encode(v)
	if err != nil {
		return err
	}

	return s.client.SetNX(ctx, s.getKey(id), b, suspendedEventTTL).Err()
}

func (s *storage[T]) Get(ctx context.Context, id string) (T, error) {
	var v T

	cr := s.client.Get(ctx, s.getKey(id))
	if err := cr.Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return v, ErrNotFound
		}

		return v, err
	}

	if err := s.decoder.Decode([]byte(cr.Val()), &v); err != nil {
		return v, err
	}

	return v, nil
}

func (s *storage[T]) Delete(ctx context.Context, id string) error {
	return s.client.Del(ctx, s.getKey(id)).Err()
}

func (s *storage[T]) getKey(id string) string {
	return s.key + "-" + id
}
