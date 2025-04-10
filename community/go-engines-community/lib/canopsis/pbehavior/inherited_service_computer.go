package pbehavior

import (
	"context"
	"errors"
	"fmt"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"github.com/bsm/redislock"
	"go.mongodb.org/mongo-driver/bson"
)

const maxRedisLockRetries = 10

type ServiceEventsData struct {
	ServiceEvents []types.Event
	ServicesMap   map[string]types.Entity
}

type InheritedServicePbhResolver interface {
	ComputeAndResolveInheritedServicePbh(
		ctx context.Context,
		resolver ComputedEntityTypeResolver,
	) (InheritedServicesPbhResolveResult, ServiceEventsData, error)
	GetResolvedInheritedServicePbh(ctx context.Context) (InheritedServicesPbhResolveResult, error)
}

type InheritedServicesPbhResolveResult struct {
	// IDs should contain all map keys
	IDs []string `json:"ids"`
	// PersonalPbh map contains max prioritized PERSONAL pbh for an entity service,
	// which should be used in service's personal pbh resolve result, e.g. service's pbh has higher priority than inherited one.
	// Events should be sent by PersonalPbh map.
	PersonalPbh map[string]ResolveResult `json:"-"`
	// InheritedPbh map contains max prioritized INHERITED pbh for an entity service,
	// which should be used in inherited pbh propagation from a parent service to a child service or from a service to a dependent resource.
	// A separate map is needed in case where personal pbh may have higher priority, but we still need to propagate inherited pbh to dependent entities.
	// Events shouldn't be sent by InheritedPbh map, it's only for calculation.
	InheritedPbh map[string]ResolveResult `json:"results"`
}

func (r *InheritedServicesPbhResolveResult) ResolveForEntity(
	entityResult ResolveResult,
	entityServices []string,
) ResolveResult {
	for _, serviceID := range entityServices {
		serviceResult, ok := r.InheritedPbh[serviceID]
		if ok && serviceResult.Type.Priority >= entityResult.Type.Priority {
			entityResult = serviceResult
		}
	}

	return entityResult
}

func (r *InheritedServicesPbhResolveResult) ResolveParentServicePbh(
	ctx context.Context,
	resolver ComputedEntityTypeResolver,
	parent types.Entity,
) (ResolveResult, error) {
	var err error

	parentServicePbh, ok := r.InheritedPbh[parent.ID]
	if !ok {
		parentServicePbh, err = resolver.Resolve(ctx, parent, time.Now())
		if err != nil {
			return ResolveResult{}, fmt.Errorf("failed to resolve pbehavior for a parent service %q: %w", parent.ID, err)
		}

		// if it wasn't resolved before and resolve shows that is has not inherited pbh now, then we can skip.
		// if it has any parent with inherited pbhs, then current parent's children will be resolved later.
		// if it's a root parent service, and it has higher prioritized pbh than inherited one, then its inherited pbh won't be propagated!
		// parentServicePbh should always be inherited for future calculations.
		if !parentServicePbh.Inherited {
			return parentServicePbh, nil
		}

		r.PersonalPbh[parent.ID] = parentServicePbh
		r.InheritedPbh[parent.ID] = parentServicePbh
	}

	return parentServicePbh, nil
}

func (r *InheritedServicesPbhResolveResult) ResolveChildServicePbh(
	ctx context.Context,
	resolver ComputedEntityTypeResolver,
	child types.Entity,
	parentInheritedPbh ResolveResult,
) error {
	var err error

	childResolvedPbh, ok := r.PersonalPbh[child.ID]
	if !ok {
		childResolvedPbh, err = resolver.Resolve(ctx, child, time.Now())
		if err != nil {
			return fmt.Errorf("failed to resolve pbh for service %q: %w", child.ID, err)
		}
	}

	if parentInheritedPbh.Type.Priority >= childResolvedPbh.Type.Priority {
		// If parent's priority is higher, then it overlaps both personal and inherited pbh.
		// No need to check a child's inherited result, because if it's not the same as personal pbh,
		// then it's 100% has lower priority.
		r.PersonalPbh[child.ID] = parentInheritedPbh
		r.InheritedPbh[child.ID] = parentInheritedPbh
	} else {
		// if parent's priority is lower, then we save childPbh as personal
		r.PersonalPbh[child.ID] = childResolvedPbh

		// but we need to check child's inherited pbh.
		childServiceInheritedPbh, ok := r.InheritedPbh[child.ID]
		if !ok && childResolvedPbh.Inherited {
			childServiceInheritedPbh = childResolvedPbh
		}

		// check priorities to set inherited pbh with more priority.
		if parentInheritedPbh.Type.Priority >= childServiceInheritedPbh.Type.Priority {
			r.InheritedPbh[child.ID] = parentInheritedPbh
		} else {
			r.InheritedPbh[child.ID] = childServiceInheritedPbh
		}
	}

	return nil
}

type inheritedServicePbhResolver struct {
	entityCollection mongo.DbCollection
	eventManager     EventManager
	store            Store
	lockClient       redis.LockClient
}

func NewInheritedServicePbhResolver(
	dbClient mongo.DbClient,
	eventManager EventManager,
	store Store,
	lockClient redis.LockClient,
) InheritedServicePbhResolver {
	return &inheritedServicePbhResolver{
		entityCollection: dbClient.Collection(mongo.EntityMongoCollection),
		eventManager:     eventManager,
		store:            store,
		lockClient:       lockClient,
	}
}

func (s *inheritedServicePbhResolver) ComputeAndResolveInheritedServicePbh(
	ctx context.Context,
	resolver ComputedEntityTypeResolver,
) (_ InheritedServicesPbhResolveResult, _ ServiceEventsData, resErr error) {
	lock, err := s.lockClient.Obtain(ctx, redis.RecomputeInheritedLockKey, redis.RecomputeLockDuration, &redislock.Options{
		RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(time.Second), maxRedisLockRetries),
	})
	if err != nil {
		return InheritedServicesPbhResolveResult{}, ServiceEventsData{}, fmt.Errorf("failed to obtain lock: %w", err)
	}

	defer func() {
		err = lock.Release(context.WithoutCancel(ctx))
		if err != nil && !errors.Is(err, redislock.ErrLockNotHeld) && resErr == nil {
			resErr = fmt.Errorf("failed to release lock: %w", err)
		}
	}()

	parentServiceIDs, err := resolver.GetComputedServicesWithInheritedPbhIDs()
	if err != nil {
		return InheritedServicesPbhResolveResult{}, ServiceEventsData{}, fmt.Errorf("failed to get service ids which have inherited pbehavior: %w", err)
	}

	if len(parentServiceIDs) == 0 {
		_, err = s.store.GetInheritedServicesPbhResolveResult(ctx)
		if err != nil && !errors.Is(err, ErrNoComputed) {
			return InheritedServicesPbhResolveResult{}, ServiceEventsData{}, err
		} else {
			err = s.store.SetInheritedServicesPbhResolveResult(ctx, InheritedServicesPbhResolveResult{})
			if err != nil {
				return InheritedServicesPbhResolveResult{}, ServiceEventsData{}, fmt.Errorf("failed to save default inherited service pbehavior resolve result: %w", err)
			}
		}

		return InheritedServicesPbhResolveResult{}, ServiceEventsData{}, nil
	}

	services := make(map[string]types.Entity)
	resolveResult := InheritedServicesPbhResolveResult{
		IDs:          make([]string, 0),
		InheritedPbh: make(map[string]ResolveResult),
		PersonalPbh:  make(map[string]ResolveResult),
	}

	parentCursor, err := s.entityCollection.Aggregate(ctx, []bson.M{
		{
			"$match": bson.M{
				"_id":     bson.M{"$in": parentServiceIDs},
				"type":    types.EntityTypeService,
				"enabled": true,
			},
		},
		{
			"$graphLookup": bson.M{
				"from":             mongo.EntityMongoCollection,
				"startWith":        "$_id",
				"connectFromField": "_id",
				"connectToField":   "services",
				"as":               "children",
				"maxDepth":         mongo.DefaultGraphLookupMaxDepth,
				"restrictSearchWithMatch": bson.M{
					"type":    types.EntityTypeService,
					"enabled": true,
				},
			},
		},
		{
			"$addFields": bson.M{
				"children": bson.M{
					"$map": bson.M{
						"input": "$children",
						"as":    "entity",
						"in":    "$$entity._id",
					},
				},
			},
		},
	})
	if err != nil {
		return InheritedServicesPbhResolveResult{}, ServiceEventsData{}, fmt.Errorf("failed to query parent services: %w", err)
	}

	defer parentCursor.Close(ctx)

	for parentCursor.Next(ctx) {
		var parent struct {
			types.Entity     `bson:",inline"`
			ChildrenServices []string `bson:"children"`
		}

		err = parentCursor.Decode(&parent)
		if err != nil {
			return InheritedServicesPbhResolveResult{}, ServiceEventsData{}, fmt.Errorf("failed to save inherited service pbehavior resolve result: %w", err)
		}

		if _, ok := services[parent.ID]; !ok {
			services[parent.ID] = parent.Entity
		}

		// check if a parent inherited pbh was already resolved, if it was a child for some other parent.
		parentServicePbh, err := resolveResult.ResolveParentServicePbh(ctx, resolver, parent.Entity)
		if err != nil {
			return InheritedServicesPbhResolveResult{}, ServiceEventsData{}, fmt.Errorf("failed to resolve pbehavior for a parent service %q: %w", parent.ID, err)
		}

		if !parentServicePbh.Inherited {
			continue
		}

		err = s.processChildren(ctx, resolver, parent.ChildrenServices, services, &resolveResult, parentServicePbh)
		if err != nil {
			return InheritedServicesPbhResolveResult{}, ServiceEventsData{}, fmt.Errorf("failed to process dependent services for a parent service %q: %w", parent.ID, err)
		}
	}

	err = s.store.SetInheritedServicesPbhResolveResult(ctx, resolveResult)
	if err != nil {
		return InheritedServicesPbhResolveResult{}, ServiceEventsData{}, fmt.Errorf("failed to save inherited service pbehavior resolve result: %w", err)
	}

	serviceEvents := make([]types.Event, 0, len(services))

	for k, v := range resolveResult.PersonalPbh {
		resolveResult.IDs = append(resolveResult.IDs, k)

		event, err := s.eventManager.GetEvent(v, services[k], datetime.CpsTime{Time: time.Now()})
		if err != nil {
			return InheritedServicesPbhResolveResult{}, ServiceEventsData{}, fmt.Errorf("failed to get event for service %q: %w", services[k].ID, err)
		}

		if event.EventType != "" {
			serviceEvents = append(serviceEvents, event)
		}
	}

	return resolveResult, ServiceEventsData{ServiceEvents: serviceEvents, ServicesMap: services}, nil
}

func (s *inheritedServicePbhResolver) processChildren(
	ctx context.Context,
	resolver ComputedEntityTypeResolver,
	childrenIDs []string,
	services map[string]types.Entity,
	resolveResult *InheritedServicesPbhResolveResult,
	parentServiceInheritedPbh ResolveResult,
) error {
	childrenCursor, err := s.entityCollection.Find(ctx, bson.M{
		"_id":     bson.M{"$in": childrenIDs},
		"enabled": true,
	})
	if err != nil {
		return fmt.Errorf("failed to find dependent context graph services: %w", err)
	}

	defer childrenCursor.Close(ctx)

	for childrenCursor.Next(ctx) {
		var child types.Entity

		err = childrenCursor.Decode(&child)
		if err != nil {
			return fmt.Errorf("failed to decode child service: %w", err)
		}

		if _, ok := services[child.ID]; !ok {
			services[child.ID] = child
		}

		err = resolveResult.ResolveChildServicePbh(ctx, resolver, child, parentServiceInheritedPbh)
		if err != nil {
			return fmt.Errorf("failed calculate child pbh: %w", err)
		}
	}

	return nil
}

func (s *inheritedServicePbhResolver) GetResolvedInheritedServicePbh(ctx context.Context) (InheritedServicesPbhResolveResult, error) {
	return s.store.GetInheritedServicesPbhResolveResult(ctx)
}
