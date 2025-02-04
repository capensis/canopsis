package pbehavior

import (
	"context"
	"errors"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/sync/errgroup"
)

const (
	notInherited = iota
	inherited
)

var ErrCacheNotLoaded = errors.New("cache is not loaded")

// ComputedEntityGetter checks if there are entities which are matched to filters. It saves matched entity ids to local cache.
type ComputedEntityGetter interface {
	Compute(ctx context.Context, notInheritedFilters []bson.M, inheritedFilters []bson.M) error
	GetComputedEntityIDs() ([]string, error)
	GetComputedServiceIDs() ([]string, error)
}

type computeJob struct {
	key int
	t   int
}

type computeJobResult struct {
	ids []string
	t   int
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
	servicesIds  []string
}

func (m *computedEntityGetter) Compute(ctx context.Context, notInheritedFilters []bson.M, inheritedFilters []bson.M) error {
	entityIds := make([]string, 0)
	servicesIds := make([]string, 0)

	if len(notInheritedFilters) == 0 && len(inheritedFilters) == 0 {
		m.entityIds = entityIds
		m.servicesIds = servicesIds

		return nil
	}

	ch := make(chan computeJob)
	resCh := make(chan computeJobResult)
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		defer close(ch)
		for key := range notInheritedFilters {
			select {
			case <-ctx.Done():
				return nil
			case ch <- computeJob{key: key, t: notInherited}:
			}
		}

		for key := range inheritedFilters {
			select {
			case <-ctx.Done():
				return nil
			case ch <- computeJob{key: key, t: inherited}:
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
				case job, ok := <-ch:
					if !ok {
						return nil
					}

					var entityIDs []string
					var err error

					switch job.t {
					case notInherited:
						entityIDs, err = m.findEntityIDs(ctx, notInheritedFilters[job.key], false)
					case inherited:
						entityIDs, err = m.findEntityIDs(ctx, inheritedFilters[job.key], true)
					}

					if err != nil {
						return err
					}

					resCh <- computeJobResult{ids: entityIDs, t: job.t}
				}
			}
		})
	}

	go func() {
		_ = g.Wait()
		close(resCh)
	}()

	for res := range resCh {
		switch res.t {
		case notInherited:
			entityIds = append(entityIds, res.ids...)
		case inherited:
			servicesIds = append(servicesIds, res.ids...)
		}
	}

	err := g.Wait()
	if err != nil {
		return err
	}

	m.entityIds = entityIds
	m.servicesIds = servicesIds

	return nil
}

func (m *computedEntityGetter) GetComputedEntityIDs() ([]string, error) {
	if m.entityIds == nil {
		return nil, ErrCacheNotLoaded
	}

	return m.entityIds, nil
}

func (m *computedEntityGetter) GetComputedServiceIDs() ([]string, error) {
	if m.servicesIds == nil {
		return nil, ErrCacheNotLoaded
	}

	return m.servicesIds, nil
}

func (m *computedEntityGetter) findEntityIDs(ctx context.Context, filter bson.M, onlyServices bool) ([]string, error) {
	mainFilter := bson.M{"$match": filter}
	var additionalFilter bson.M
	if onlyServices {
		additionalFilter = bson.M{"$match": bson.M{"type": types.EntityTypeService, "enabled": true}}
	} else {
		additionalFilter = bson.M{"$match": bson.M{"enabled": true}}
	}

	cursor, err := m.dbCollection.Aggregate(ctx, []bson.M{
		mainFilter,
		additionalFilter,
		{"$project": bson.M{"_id": 1}},
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
