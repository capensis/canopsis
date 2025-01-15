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

// ComputedEntityGetter checks if there are entities which are matched to filters. It saves matched entity ids to local cache.
type ComputedEntityGetter interface {
	Compute(ctx context.Context, filters []bson.M) error
	GetComputedEntityIDs() ([]string, error)
}

func NewComputedEntityGetter(dbClient mongo.DbClient) ComputedEntityGetter {
	return &computedEntityGetter{
		dbCollection: dbClient.Collection(mongo.EntityMongoCollection),
	}
}

// computedEntityGetter executes mongo query to check if entity is matched.
type computedEntityGetter struct {
	dbCollection mongo.DbCollection
	entityIds    []string
}

func (m *computedEntityGetter) Compute(ctx context.Context, filters []bson.M) error {
	entityIds := make([]string, 0)
	if len(filters) == 0 {
		m.entityIds = entityIds
		return nil
	}

	ch := make(chan int)
	resCh := make(chan []string)
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

					resCh <- entityIDs
				}
			}
		})
	}

	go func() {
		_ = g.Wait()
		close(resCh)
	}()

	for res := range resCh {
		entityIds = append(entityIds, res...)
	}

	err := g.Wait()
	if err != nil {
		return err
	}

	m.entityIds = entityIds
	return nil
}

func (m *computedEntityGetter) GetComputedEntityIDs() ([]string, error) {
	if m.entityIds == nil {
		return nil, ErrCacheNotLoaded
	}

	return m.entityIds, nil
}

func (m *computedEntityGetter) findEntityIDs(ctx context.Context, filter bson.M) ([]string, error) {
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
