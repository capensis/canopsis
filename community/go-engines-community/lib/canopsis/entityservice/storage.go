package entityservice

import (
	"context"
	"errors"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/oldpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
)

const (
	cacheKey = "services"
	// lastUpdateCacheKey is used to prevent race condition when multiple goroutine or instances update cache.
	lastUpdateCacheKey = "services-last-update"
)

type Storage interface {
	// ReloadAll loads all enabled services from database and saves their data to cache.
	ReloadAll(ctx context.Context) ([]ServiceData, error)
	// Reload loads service by id from database and saves its data to cache.
	// Return parameter isNew contains true if service didn't exist in cache.
	// Return parameter isDisabled contains true if service is disabled. Service data is nil in this case.
	Reload(ctx context.Context, id string) (s *ServiceData, isNew bool, isDisabled bool, isSoftDeleted bool, err error)
	// GetAll returns all services data from cache.
	GetAll(ctx context.Context) ([]ServiceData, error)
	// Get returns service data  by id from cache.
	Get(ctx context.Context, id string) (*ServiceData, error)
}

func NewStorage(
	adapter Adapter,
	client redis.UniversalClient,
	encoder encoding.Encoder,
	decoder encoding.Decoder,
	logger zerolog.Logger,
) Storage {
	return &redisStorage{
		adapter: adapter,
		client:  client,
		encoder: encoder,
		decoder: decoder,
		logger:  logger,
	}
}

type redisStorage struct {
	adapter Adapter
	client  redis.UniversalClient
	encoder encoding.Encoder
	decoder encoding.Decoder
	logger  zerolog.Logger
}

type ServiceData struct {
	ID             string `json:"_id"`
	OutputTemplate string `json:"output_template,omitempty"`

	EntityPattern     pattern.Entity               `json:"entity_pattern,omitempty"`
	OldEntityPatterns oldpattern.EntityPatternList `json:"old_entity_patterns,omitempty"`
}

func (s *redisStorage) ReloadAll(ctx context.Context) ([]ServiceData, error) {
	var data []ServiceData

	txf := func(tx *redis.Tx) error {
		res := tx.HKeys(ctx, cacheKey)
		if err := res.Err(); err != nil {
			return err
		}
		oldKeys := res.Val()

		services, err := s.adapter.GetEnabled(ctx)
		if err != nil {
			return err
		}

		m := make(map[string]interface{}, len(services))
		data = make([]ServiceData, len(services))
		for i, v := range services {
			data[i] = ServiceData{
				ID:                v.ID,
				OutputTemplate:    v.OutputTemplate,
				EntityPattern:     v.EntityPattern,
				OldEntityPatterns: v.OldEntityPatterns,
			}
			str, err := s.encoder.Encode(data[i])
			if err != nil {
				return err
			}
			m[v.ID] = str
		}

		removed := make([]string, 0)
		for _, k := range oldKeys {
			if _, ok := m[k]; !ok {
				removed = append(removed, k)
			}
		}

		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			if len(m) > 0 {
				pipe.HSet(ctx, cacheKey, m)
			}

			if len(removed) > 0 {
				pipe.HDel(ctx, cacheKey, removed...)
			}

			pipe.Set(ctx, lastUpdateCacheKey, time.Now(), 0)
			return nil
		})
		return err
	}

	for i := 0; i < maxRetries; i++ {
		err := s.client.Watch(ctx, txf, lastUpdateCacheKey)

		if err != nil {
			if err == redis.TxFailedErr {
				continue
			}

			return nil, err
		}

		return data, nil
	}

	return nil, errors.New("reached maximum number of retries")
}

func (s *redisStorage) Reload(ctx context.Context, id string) (*ServiceData, bool, bool, bool, error) {
	var data *ServiceData
	var isNew, isDisabled, isSoftDeleted bool

	txf := func(tx *redis.Tx) error {
		data = nil
		isNew = false
		isDisabled = false
		isSoftDeleted = false

		service, err := s.adapter.GetByID(ctx, id)
		if err != nil {
			return err
		}

		var str []byte

		if service != nil && service.SoftDeleted == nil {
			if service.Enabled {
				data = &ServiceData{
					ID:                service.ID,
					OutputTemplate:    service.OutputTemplate,
					EntityPattern:     service.EntityPattern,
					OldEntityPatterns: service.OldEntityPatterns,
				}

				str, err = s.encoder.Encode(data)
				if err != nil {
					return err
				}
			} else {
				isDisabled = true
			}
		}

		isSoftDeleted = service != nil && service.SoftDeleted != nil
		res, err := tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			if data == nil {
				pipe.HDel(ctx, cacheKey, id)
			} else {
				pipe.HSet(ctx, cacheKey, data.ID, str)
			}

			pipe.Set(ctx, lastUpdateCacheKey, time.Now(), 0)
			return nil
		})

		if err != nil {
			return err
		}

		if data != nil {
			if ic, ok := res[0].(*redis.IntCmd); ok {
				isNew = ic.Val() > 0
			}
		}

		return nil
	}

	for i := 0; i < maxRetries; i++ {
		err := s.client.Watch(ctx, txf, lastUpdateCacheKey)

		if err != nil {
			if err == redis.TxFailedErr {
				continue
			}

			return nil, false, false, false, err
		}

		return data, isNew, isDisabled, isSoftDeleted, nil
	}

	return nil, false, false, false, errors.New("reached maximum number of retries")
}

func (s *redisStorage) Get(ctx context.Context, id string) (*ServiceData, error) {
	res := s.client.HGet(ctx, cacheKey, id)
	if err := res.Err(); err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	data := &ServiceData{}
	err := s.decoder.Decode([]byte(res.Val()), data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *redisStorage) GetAll(ctx context.Context) ([]ServiceData, error) {
	res := s.client.HGetAll(ctx, cacheKey)
	if err := res.Err(); err != nil {
		if err == redis.Nil {
			return nil, nil
		}

		return nil, err
	}

	data := make([]ServiceData, len(res.Val()))
	i := 0
	for _, v := range res.Val() {
		err := s.decoder.Decode([]byte(v), &data[i])
		if err != nil {
			return nil, err
		}

		i++
	}

	return data, nil
}
