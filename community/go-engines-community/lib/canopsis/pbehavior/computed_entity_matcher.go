package pbehavior

import (
	"context"
	"errors"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/sync/errgroup"
)

var ErrCacheNotLoaded = errors.New("cache is not loaded")

// ComputedEntityMatcher checks if an entity is matched to filter. It precomputes
// filter-entity associations and uses computed data to resolve matched filters.
type ComputedEntityMatcher interface {
	// LoadAll computes filter-entity associations to memory.
	LoadAll(ctx context.Context, filters map[string]interface{}) error
	// Match matches entity to filters by precomputed data in memory.
	Match(entityID string) ([]string, error)
	GetComputedEntityIDs() ([]string, error)
}

func NewComputedEntityMatcher(dbClient mongo.DbClient) ComputedEntityMatcher {
	return &computedEntityMatcher{
		dbCollection: dbClient.Collection(mongo.EntityMongoCollection),
	}
}

// entityMatcher executes mongo query to check if entity is matched.
type computedEntityMatcher struct {
	dbCollection   mongo.DbCollection
	keysByEntityID map[string][]string
}

func (m *computedEntityMatcher) LoadAll(ctx context.Context, filters map[string]interface{}) error {
	keysByEntityID := make(map[string][]string, len(filters))
	if len(filters) == 0 {
		m.keysByEntityID = keysByEntityID
		return nil
	}

	ch := make(chan string)
	type workerResult struct {
		key       string
		entityIDs []string
	}
	resCh := make(chan workerResult)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		defer close(ch)
		for key := range filters {
			select {
			case <-ctx.Done():
				return nil
			case ch <- key:
			}
		}

		return nil
	})

	for i := 0; i < DefaultPoolSize; i++ {
		g.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					return nil
				case key, ok := <-ch:
					if !ok {
						return nil
					}

					entityIDs, err := m.findEntityIDs(ctx, filters[key])
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

	for res := range resCh {
		for _, entityID := range res.entityIDs {
			keysByEntityID[entityID] = append(keysByEntityID[entityID], res.key)
		}
	}

	err := g.Wait()
	if err != nil {
		return err
	}

	m.keysByEntityID = keysByEntityID
	return nil
}

func (m *computedEntityMatcher) Match(entityID string) ([]string, error) {
	if m.keysByEntityID == nil {
		return nil, ErrCacheNotLoaded
	}

	return m.keysByEntityID[entityID], nil
}

func (m *computedEntityMatcher) GetComputedEntityIDs() ([]string, error) {
	if m.keysByEntityID == nil {
		return nil, ErrCacheNotLoaded
	}

	entityIDs := make([]string, len(m.keysByEntityID))
	i := 0
	for entityID := range m.keysByEntityID {
		entityIDs[i] = entityID
		i++
	}

	return entityIDs, nil
}

func (m *computedEntityMatcher) findEntityIDs(ctx context.Context, filter interface{}) ([]string, error) {
	cursor, err := m.dbCollection.Aggregate(ctx, []bson.M{
		{"$match": filter},
		{"$match": bson.M{"enabled": true}},
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
