package entityservice

import (
	"context"
	"math"
	"sync"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	maxWorkersCount = 10
)

// NewManager creates new manager.
func NewManager(
	adapter Adapter,
	entityAdapter entity.Adapter,
	storage Storage,
	logger zerolog.Logger,
) Manager {
	return &manager{
		storage:       storage,
		adapter:       adapter,
		entityAdapter: entityAdapter,
		logger:        logger,
	}
}

type manager struct {
	storage       Storage
	adapter       Adapter
	entityAdapter entity.Adapter
	logger        zerolog.Logger
}

func (m *manager) LoadServices(ctx context.Context) error {
	data, err := m.storage.ReloadAll(ctx)
	if err != nil {
		return err
	}

	ids := make([]string, len(data))
	for i := range data {
		ids[i] = data[i].ID
	}

	m.logger.Debug().Strs("services", ids).Msg("load services")

	return nil
}

func (m *manager) UpdateServices(ctx context.Context, entities []types.Entity) (map[string][]string, map[string][]string, error) {
	if len(entities) == 0 {
		return nil, nil, nil
	}

	services, err := m.storage.GetAll(ctx)
	if err != nil {
		return nil, nil, err
	}

	workerCh := make(chan ServiceData, len(services))
	go func() {
		defer close(workerCh)
		for i := range services {
			workerCh <- services[i]
		}
	}()

	workerCount := int(math.Min(float64(maxWorkersCount), float64(len(services))))
	outCh, err := m.runWorkers(ctx, workerCount, workerCh, entities)
	if err != nil {
		return nil, nil, err
	}

	return m.combineResult(outCh)
}

func (m *manager) UpdateService(ctx context.Context, serviceID string) (bool, []string, error) {
	data, isNew, isDisabled, err := m.storage.Reload(ctx, serviceID)
	if err != nil {
		return false, nil, err
	}

	// Change context graph only for completely removed service.
	if data == nil && !isDisabled {
		removedFrom, err := m.removeService(ctx, serviceID)
		if err != nil {
			return false, nil, err
		}

		return true, removedFrom, nil
	}

	if isDisabled {
		return true, nil, nil
	}

	query, negativeQuery, err := getServiceQueries(*data)
	if err != nil {
		m.logger.Err(err).Str("service", data.ID).Msgf("service has invalid pattern")
	}
	// Do not ignore empty negativeQuery to remove service from context graph.
	removedIDs, err := m.entityAdapter.RemoveImpactByQuery(ctx, negativeQuery, data.ID)
	if err != nil {
		return false, nil, err
	}

	var addedIDs []string
	// Ignore empty query to not add service to context graph.
	if query != nil {
		addedIDs, err = m.entityAdapter.AddImpactByQuery(ctx, query, data.ID)
		if err != nil {
			return false, nil, err
		}
	}

	if len(removedIDs) > 0 {
		_, err := m.adapter.RemoveDepends(ctx, data.ID, removedIDs)
		if err != nil {
			return false, nil, err
		}
	}

	if len(addedIDs) > 0 {
		_, err := m.adapter.AddDepends(ctx, data.ID, addedIDs)
		if err != nil {
			return false, nil, err
		}
	}

	return isNew || len(removedIDs) > 0 || len(addedIDs) > 0, nil, nil
}

func (m *manager) ReloadService(ctx context.Context, serviceID string) error {
	_, _, _, err := m.storage.Reload(ctx, serviceID)
	return err
}

func (m *manager) HasEntityServiceByComponentInfos(ctx context.Context) (bool, error) {
	services, err := m.storage.GetAll(ctx)
	if err != nil {
		return false, err
	}

	for _, s := range services {
		if len(s.EntityPattern) > 0 {
			if s.EntityPattern.HasComponentInfosField() {
				return true, nil
			}
		} else if s.OldEntityPatterns.IsSet() && s.OldEntityPatterns.IsValid() {
			for _, p := range s.OldEntityPatterns.Patterns {
				if len(p.ComponentInfos) > 0 {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

func (m *manager) combineResult(ch <-chan serviceUpdate) (map[string][]string, map[string][]string, error) {
	addedTo := make(map[string][]string)
	removedFrom := make(map[string][]string)
	for v := range ch {
		if v.IsAdded {
			for _, id := range v.Entities {
				addedTo[id] = append(addedTo[id], v.ID)
			}
		} else {
			for _, id := range v.Entities {
				removedFrom[id] = append(removedFrom[id], v.ID)
			}
		}
	}

	return addedTo, removedFrom, nil
}

type serviceUpdate struct {
	ID       string
	Entities []string
	IsAdded  bool
}

// runWorkers updates entity services concurrently.
func (m *manager) runWorkers(
	ctx context.Context,
	workerCount int,
	inCh <-chan ServiceData,
	entities []types.Entity,
) (<-chan serviceUpdate, error) {
	errCh := make(chan error, workerCount)
	outCh := make(chan serviceUpdate)

	go func() {
		defer close(outCh)
		wg := sync.WaitGroup{}

		for i := 0; i < workerCount; i++ {
			wg.Add(1)

			go func() {
				defer wg.Done()

				for {
					select {
					case <-ctx.Done():
						return
					case data, ok := <-inCh:
						if !ok {
							return
						}

						r, err := m.processService(ctx, data, entities)
						if err != nil {
							errCh <- err
							return
						}

						for _, v := range r {
							outCh <- v
						}
					}
				}
			}()
		}

		wg.Wait()
	}()

	select {
	case err := <-errCh:
		return nil, err
	default:
	}

	return outCh, nil
}

// processService removes entities from entity service dependencies if they don't match
// pattern anymore and adds entities to entity service dependencies if they matches pattern.
func (m *manager) processService(ctx context.Context, data ServiceData, entities []types.Entity) ([]serviceUpdate, error) {
	added := make([]string, 0)
	removed := make([]string, 0)

	for _, e := range entities {
		found := false
		enabled := e.Enabled

		for _, impact := range e.Impacts {
			if impact == data.ID {
				found = true
				break
			}
		}

		match := false
		if len(data.EntityPattern) > 0 {
			var err error
			match, _, err = data.EntityPattern.Match(e)
			if err != nil {
				m.logger.Err(err).Str("service", data.ID).Msgf("service has invalid pattern")
			}
		} else if data.OldEntityPatterns.IsSet() {
			if data.OldEntityPatterns.IsValid() {
				match = data.OldEntityPatterns.Matches(&e)
			} else {
				m.logger.Err(pattern.ErrInvalidOldEntityPattern).Str("service", data.ID).Msgf("service has invalid pattern")
			}
		}

		if match {
			if !found && enabled {
				added = append(added, e.ID)
			}

			if found && !enabled {
				removed = append(removed, e.ID)
			}
		} else if found {
			removed = append(removed, e.ID)
		}
	}

	res := make([]serviceUpdate, 0)
	if len(added) > 0 {
		ok, err := m.adapter.AddDepends(ctx, data.ID, added)
		if err != nil {
			return nil, err
		}

		if ok {
			err := m.entityAdapter.AddImpacts(ctx, added, []string{data.ID})
			if err != nil {
				return nil, err
			}

			res = append(res, serviceUpdate{
				ID:       data.ID,
				Entities: added,
				IsAdded:  true,
			})
		}
	}

	if len(removed) > 0 {
		ok, err := m.adapter.RemoveDepends(ctx, data.ID, removed)
		if err != nil {
			return nil, err
		}

		if ok {
			err := m.entityAdapter.RemoveImpacts(ctx, removed, []string{data.ID})
			if err != nil {
				return nil, err
			}

			res = append(res, serviceUpdate{
				ID:       data.ID,
				Entities: removed,
				IsAdded:  false,
			})
		}
	}

	return res, nil
}

// removeService removes service from cache and context graph.
func (m *manager) removeService(ctx context.Context, serviceID string) ([]string, error) {
	_, err := m.entityAdapter.RemoveImpactByQuery(ctx, bson.M{"impact": serviceID}, serviceID)
	if err != nil {
		return nil, err
	}

	ids, err := m.adapter.RemoveDependByQuery(ctx, bson.M{"depends": serviceID}, serviceID)
	if err != nil {
		return nil, err
	}

	removedFromIDs := make([]string, 0)
	data, err := m.storage.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		for _, s := range data {
			if s.ID == id {
				removedFromIDs = append(removedFromIDs, id)
				break
			}
		}
	}

	return removedFromIDs, nil
}

func getServiceQueries(data ServiceData) (interface{}, interface{}, error) {
	var query, negativeQuery interface{}
	var err error

	if len(data.EntityPattern) > 0 {
		query, err = data.EntityPattern.ToMongoQuery("")
		if err != nil {
			return nil, nil, err
		}

		negativeQuery, err = data.EntityPattern.ToNegativeMongoQuery("")
		if err != nil {
			return nil, nil, err
		}
	} else if data.OldEntityPatterns.IsSet() {
		if !data.OldEntityPatterns.IsValid() {
			return nil, nil, pattern.ErrInvalidOldEntityPattern
		}
		query = data.OldEntityPatterns.AsMongoDriverQuery()
		negativeQuery = data.OldEntityPatterns.AsNegativeMongoDriverQuery()
	}

	return query, negativeQuery, nil
}
