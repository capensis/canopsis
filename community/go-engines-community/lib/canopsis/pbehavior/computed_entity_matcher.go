package pbehavior

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	libredis "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/sync/errgroup"
)

// ComputedEntityMatcher checks if an entity is matched to filter. It precomputes
// filter-entity associations and uses computed data to resolve matched filters.
type ComputedEntityMatcher interface {
	// LoadAll computes filter-entity associations.
	LoadAll(ctx context.Context, filters map[string]string) error
	// Match matches entity to filters by precomputed data.
	Match(ctx context.Context, entityID string) ([]string, error)
	GetComputedEntityIDs(ctx context.Context) ([]string, error)
}

func NewComputedEntityMatcher(
	dbClient mongo.DbClient,
	redisClient redis.Cmdable,
	encoder encoding.Encoder,
	decoder encoding.Decoder,
) ComputedEntityMatcher {
	return &computedEntityMatcher{
		dbCollection: dbClient.Collection(mongo.EntityMongoCollection),
		redisClient:  redisClient,
		encoder:      encoder,
		decoder:      decoder,
		key:          libredis.PbehaviorEntityMatchKey,
	}
}

// entityMatcher parses filter and executes parsed mongo query to check if entity is matched.
type computedEntityMatcher struct {
	dbCollection mongo.DbCollection
	redisClient  redis.Cmdable
	encoder      encoding.Encoder
	decoder      encoding.Decoder
	key          string
}

func (m *computedEntityMatcher) LoadAll(ctx context.Context, filters map[string]string) error {
	ch := make(chan string)
	type workerResult struct {
		key       string
		entityIDs []string
	}
	resCh := make(chan workerResult)

	go func() {
		defer close(ch)
		for key := range filters {
			select {
			case <-ctx.Done():
				return
			case ch <- key:
			}
		}
	}()

	g, gctx := errgroup.WithContext(ctx)

	for i := 0; i < 10; i++ {
		g.Go(func() error {
			for {
				select {
				case <-gctx.Done():
					return nil
				case key, ok := <-ch:
					if !ok {
						return nil
					}

					entityIDs, err := m.findEntityIDs(gctx, filters[key])
					if err != nil {
						return err
					}

					resCh <- workerResult{
						key:       key,
						entityIDs: entityIDs,
					}
				}
			}
		})
	}

	go func() {
		_ = g.Wait()
		close(resCh)
	}()

	keysByEntityID := make(map[string][]string)
	for res := range resCh {
		for _, entityID := range res.entityIDs {
			redisKey := m.key + entityID
			keysByEntityID[redisKey] = append(keysByEntityID[redisKey], res.key)
		}
	}

	err := g.Wait()
	if err != nil {
		return err
	}

	encodedData := make(map[string]interface{}, len(keysByEntityID))
	for k, v := range keysByEntityID {
		encodedData[k], err = m.encoder.Encode(v)
		if err != nil {
			return fmt.Errorf("cannot encode entity ids: %w", err)
		}
	}

	if len(encodedData) > 0 {
		err := m.redisClient.MSet(ctx, encodedData).Err()
		if err != nil {
			return fmt.Errorf("cannot set entity ids: %w", err)
		}
	}

	return m.deleteMissingKeys(ctx, keysByEntityID)
}

func (m *computedEntityMatcher) Match(ctx context.Context, entityID string) ([]string, error) {
	res := m.redisClient.Get(ctx, m.key+entityID)
	if err := res.Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}

		return nil, fmt.Errorf("cannot get keys: %w", err)
	}

	var matchedKeys []string
	err := m.decoder.Decode([]byte(res.Val()), &matchedKeys)
	if err != nil {
		return nil, fmt.Errorf("cannot decode entity ids: %w", err)
	}

	return matchedKeys, nil
}

func (m *computedEntityMatcher) GetComputedEntityIDs(ctx context.Context) ([]string, error) {
	var cursor uint64
	entityIDs := make([]string, 0)
	processedKeys := make(map[string]bool)

	for {
		res := m.redisClient.Scan(ctx, cursor, fmt.Sprintf("%s*", m.key), redisStep)
		if err := res.Err(); err != nil {
			return nil, fmt.Errorf("cannot scan keys: %w", err)
		}

		var keys []string
		keys, cursor = res.Val()
		for _, key := range keys {
			if !processedKeys[key] {
				processedKeys[key] = true
				entityIDs = append(entityIDs, strings.TrimPrefix(key, m.key))
			}
		}

		if cursor == 0 {
			break
		}
	}

	return entityIDs, nil
}

func (m *computedEntityMatcher) findEntityIDs(ctx context.Context, filter string) ([]string, error) {
	match, err := transformFilter(filter)
	if err != nil {
		return nil, err
	}

	cursor, err := m.dbCollection.Aggregate(ctx, []bson.M{
		match,
		{"$project": bson.M{
			"_id": 1,
		}},
	})
	if err != nil {
		return nil, fmt.Errorf("cannot execute filter: %w", err)
	}

	doc := make([]struct {
		ID string `bson:"_id"`
	}, 0)
	err = cursor.All(ctx, &doc)
	if err != nil {
		return nil, fmt.Errorf("cannot decode filter result: %w", err)
	}

	if len(doc) == 0 {
		return nil, nil
	}

	entityIDs := make([]string, len(doc))
	for i, v := range doc {
		entityIDs[i] = v.ID
	}

	return entityIDs, nil
}

func (m *computedEntityMatcher) deleteMissingKeys(ctx context.Context, keysByEntityID map[string][]string) error {
	var cursor uint64

	for {
		res := m.redisClient.Scan(ctx, cursor, fmt.Sprintf("%s*", m.key), redisStep)
		if err := res.Err(); err != nil {
			return fmt.Errorf("cannot scan keys: %w", err)
		}

		var keys []string
		keys, cursor = res.Val()
		unprocessedKeys := make([]string, 0)
		for _, key := range keys {
			if _, ok := keysByEntityID[key]; !ok {
				unprocessedKeys = append(unprocessedKeys, key)
			}
		}

		if len(unprocessedKeys) > 0 {
			resGet := m.redisClient.Del(ctx, unprocessedKeys...)
			if err := resGet.Err(); err != nil {
				return fmt.Errorf("cannot del keys: %w", err)
			}
		}

		if cursor == 0 {
			break
		}
	}

	return nil
}
