package entityservice

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/pattern"
	"github.com/go-redis/redis/v8"
)

const cacheKey = "services"

type Storage interface {
	SaveAll(ctx context.Context, data []ServiceData) error
	Save(ctx context.Context, data ServiceData) error
	Load(ctx context.Context) ([]ServiceData, error)
	Get(ctx context.Context, id string) (*ServiceData, error)
	Delete(ctx context.Context, id string) error
}

func NewStorage(
	client redis.Cmdable,
	encoder encoding.Encoder,
	decoder encoding.Decoder,
) Storage {
	return &redisStorage{
		client:  client,
		encoder: encoder,
		decoder: decoder,
	}
}

type redisStorage struct {
	client  redis.Cmdable
	encoder encoding.Encoder
	decoder encoding.Decoder
}

type ServiceData struct {
	ID             string                    `json:"_id"`
	OutputTemplate string                    `json:"output_template,omitempty"`
	EntityPatterns pattern.EntityPatternList `json:"entity_patterns"`
	Impacts        []string                  `json:"impacts"`
}

func (s *redisStorage) SaveAll(ctx context.Context, data []ServiceData) error {
	m := make(map[string]interface{}, len(data))
	for _, v := range data {
		str, err := s.encoder.Encode(v)
		if err != nil {
			return err
		}
		m[v.ID] = str
	}
	// Save services.
	if len(m) > 0 {
		err := s.client.HSet(ctx, cacheKey, m).Err()
		if err != nil {
			return err
		}
	}
	// Remove deleted services.
	res := s.client.HKeys(ctx, cacheKey)
	if err := res.Err(); err != nil {
		return err
	}

	removed := make([]string, 0)
	for _, k := range res.Val() {
		if _, ok := m[k]; !ok {
			removed = append(removed, k)
		}
	}

	if len(removed) > 0 {
		err := s.client.HDel(ctx, cacheKey, removed...).Err()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *redisStorage) Save(ctx context.Context, data ServiceData) error {
	str, err := s.encoder.Encode(data)
	if err != nil {
		return err
	}

	err = s.client.HSet(ctx, cacheKey, data.ID, str).Err()
	if err != nil {
		return err
	}

	return nil
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

func (s *redisStorage) Delete(ctx context.Context, id string) error {
	return s.client.HDel(ctx, cacheKey, id).Err()
}

func (s *redisStorage) Load(ctx context.Context) ([]ServiceData, error) {
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
