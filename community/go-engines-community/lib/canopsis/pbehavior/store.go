package pbehavior

import (
	"context"
	"errors"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	libredis "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/timespan"
	"github.com/redis/go-redis/v9"
)

var ErrNoComputed = errors.New("pbehavior intervals not computed")
var ErrRecomputeNeed = errors.New("provided time is out of computed date, probably need recompute data")

const redisStep = 1000

type Store interface {
	SetSpan(ctx context.Context, span timespan.Span) error
	GetSpan(ctx context.Context) (timespan.Span, error)
	SetComputed(ctx context.Context, computed ComputeResult) error
	GetComputed(ctx context.Context) (ComputeResult, error)
	SetComputedPbehavior(ctx context.Context, pbhID string, computed ComputedPbehavior) error
	DelComputedPbehavior(ctx context.Context, pbhID string) error
	GetComputedByIDs(ctx context.Context, pbehaviorIDs []string) (ComputeResult, error)
}

func NewStore(
	client redis.Cmdable,
	encoder encoding.Encoder,
	decoder encoding.Decoder,
) Store {
	return &store{
		client:               client,
		encoder:              encoder,
		decoder:              decoder,
		spanKey:              libredis.PbehaviorSpanKey,
		typesKey:             libredis.PbehaviorTypesKey,
		defaultActiveTypeKey: libredis.PbehaviorDefaultActiveTypeKey,
		computedKey:          libredis.PbehaviorComputedKey,
	}
}

type store struct {
	client  redis.Cmdable
	encoder encoding.Encoder
	decoder encoding.Decoder

	spanKey, typesKey, defaultActiveTypeKey, computedKey string
}

func (s *store) SetSpan(ctx context.Context, span timespan.Span) error {
	b, err := s.encoder.Encode(span)
	if err != nil {
		return fmt.Errorf("cannot encode span: %w", err)
	}

	err = s.client.Set(ctx, s.spanKey, b, 0).Err()
	if err != nil {
		return fmt.Errorf("cannot set span: %w", err)
	}

	return nil
}

func (s *store) GetSpan(ctx context.Context) (timespan.Span, error) {
	res := s.client.Get(ctx, s.spanKey)
	if err := res.Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return timespan.Span{}, ErrNoComputed
		}

		return timespan.Span{}, fmt.Errorf("cannot get span: %w", err)
	}

	span := timespan.Span{}
	err := s.decoder.Decode([]byte(res.Val()), &span)
	if err != nil {
		return timespan.Span{}, fmt.Errorf("cannot decode span: %w", err)
	}

	return span, nil
}

func (s *store) SetComputed(ctx context.Context, computed ComputeResult) error {
	data := make(map[string]interface{}, len(computed.ComputedPbehaviors)+2)
	var err error
	types := Types{T: computed.TypesByID}
	data[s.typesKey], err = s.encoder.Encode(types)
	if err != nil {
		return fmt.Errorf("cannot encode computed types: %w", err)
	}
	data[s.defaultActiveTypeKey] = computed.DefaultActiveType

	for k, v := range computed.ComputedPbehaviors {
		data[s.computedKey+k], err = s.encoder.Encode(v)
		if err != nil {
			return fmt.Errorf("cannot encode computed pbehavior: %w", err)
		}
	}

	if len(data) > 0 {
		err := s.client.MSet(ctx, data).Err()
		if err != nil {
			return fmt.Errorf("cannot set computed: %w", err)
		}
	}

	var cursor uint64
	for {
		res := s.client.Scan(ctx, cursor, s.computedKey+"*", redisStep)
		if err := res.Err(); err != nil {
			return fmt.Errorf("cannot scan computed: %w", err)
		}

		var keys []string
		keys, cursor = res.Val()
		unprocessedKeys := make([]string, 0)
		for _, key := range keys {
			if _, ok := data[key]; !ok {
				unprocessedKeys = append(unprocessedKeys, key)
			}
		}

		if len(unprocessedKeys) > 0 {
			resGet := s.client.Del(ctx, unprocessedKeys...)
			if err := resGet.Err(); err != nil {
				return fmt.Errorf("cannot del computed: %w", err)
			}
		}

		if cursor == 0 {
			break
		}
	}

	return nil
}

func (s *store) GetComputed(ctx context.Context) (ComputeResult, error) {
	computed, err := s.getTypes(ctx)
	if err != nil {
		return computed, err
	}

	var cursor uint64
	processedKeys := make(map[string]bool)
	pbhs := make(map[string]ComputedPbehavior)

	for {
		res := s.client.Scan(ctx, cursor, s.computedKey+"*", redisStep)
		if err := res.Err(); err != nil {
			return computed, fmt.Errorf("cannot scan computed: %w", err)
		}

		var keys []string
		keys, cursor = res.Val()
		unprocessedKeys := make([]string, 0)
		for _, key := range keys {
			if !processedKeys[key] {
				unprocessedKeys = append(unprocessedKeys, key)
				processedKeys[key] = true
			}
		}

		if len(unprocessedKeys) > 0 {
			resGet := s.client.MGet(ctx, unprocessedKeys...)
			if err := resGet.Err(); err != nil {
				return computed, fmt.Errorf("cannot get computed: %w", err)
			}

			for i, v := range resGet.Val() {
				switch str := v.(type) {
				case string:
					pbh := ComputedPbehavior{}
					err := s.decoder.Decode([]byte(str), &pbh)
					if err != nil {
						return computed, fmt.Errorf("cannot decode computed: %w", err)
					}

					pbhs[unprocessedKeys[i][len(s.computedKey):]] = pbh
				case nil:
					/*do nothing*/
				default:
					return computed, fmt.Errorf("expected string by key %q but got %T %+v", unprocessedKeys[i], v, v)
				}
			}
		}

		if cursor == 0 {
			break
		}
	}

	computed.ComputedPbehaviors = pbhs

	return computed, nil
}

func (s *store) SetComputedPbehavior(ctx context.Context, pbhID string, computed ComputedPbehavior) error {
	b, err := s.encoder.Encode(computed)
	if err != nil {
		return fmt.Errorf("cannot encode computed pbehavior: %w", err)
	}

	err = s.client.Set(ctx, s.computedKey+pbhID, b, 0).Err()
	if err != nil {
		return fmt.Errorf("cannot set computed pbehavior: %w", err)
	}

	return nil
}

func (s *store) DelComputedPbehavior(ctx context.Context, pbhID string) error {
	err := s.client.Del(ctx, s.computedKey+pbhID).Err()
	if err != nil {
		return fmt.Errorf("cannot gel computed pbehavior: %w", err)
	}

	return nil
}

func (s *store) GetComputedByIDs(ctx context.Context, pbhIDs []string) (ComputeResult, error) {
	if len(pbhIDs) == 0 {
		return ComputeResult{}, nil
	}

	computed, err := s.getTypes(ctx)
	if err != nil {
		return computed, err
	}

	keys := make([]string, len(pbhIDs))
	for i, id := range pbhIDs {
		keys[i] = s.computedKey + id
	}

	res := s.client.MGet(ctx, keys...)
	if err := res.Err(); err != nil {
		return computed, fmt.Errorf("cannot get computed: %w", err)
	}

	pbhs := make(map[string]ComputedPbehavior)

	for i, v := range res.Val() {
		switch str := v.(type) {
		case string:
			pbh := ComputedPbehavior{}
			err := s.decoder.Decode([]byte(str), &pbh)
			if err != nil {
				return computed, fmt.Errorf("cannot decode computed: %w", err)
			}
			pbhs[keys[i][len(s.computedKey):]] = pbh
		case nil:
			/*do nothing*/
		default:
			return computed, fmt.Errorf("expected string by key %q but got %T %+v", keys[i], v, v)
		}
	}

	computed.ComputedPbehaviors = pbhs

	return computed, nil
}

func (s *store) getTypes(ctx context.Context) (ComputeResult, error) {
	computed := ComputeResult{}
	keys := []string{s.typesKey, s.defaultActiveTypeKey}
	res := s.client.MGet(ctx, keys...)
	if err := res.Err(); err != nil {
		return computed, fmt.Errorf("cannot get types: %w", err)
	}

	switch str := res.Val()[0].(type) {
	case string:
		types := Types{}
		err := s.decoder.Decode([]byte(str), &types)
		if err != nil {
			return computed, fmt.Errorf("cannot decode types: %w", err)
		}
		computed.TypesByID = types.T
	case nil:
		/*do nothing*/
	default:
		return computed, fmt.Errorf("expected string by key %q but got %T %+v", keys[0], res.Val()[0], res.Val()[0])
	}

	switch str := res.Val()[1].(type) {
	case string:
		computed.DefaultActiveType = str
	case nil:
		/*do nothing*/
	default:
		return computed, fmt.Errorf("expected string by key %q but got %T %+v", keys[1], res.Val()[1], res.Val()[1])
	}

	return computed, nil
}
