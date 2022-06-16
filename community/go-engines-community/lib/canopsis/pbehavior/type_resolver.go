package pbehavior

import (
	"context"
	"fmt"
	"sort"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/timespan"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

// TypeResolver figures out in which state provided entity at the moment is.
type TypeResolver interface {
	// Resolve returns current type for entity if there is corresponding periodical behavior.
	// Otherwise it returns default active type.
	// An entity is matched to a pbehavior by an entity pattern or by cachedMatchedPbehaviorIds for old pbehaviors' queries.
	Resolve(ctx context.Context, t time.Time, entity types.Entity, cachedMatchedPbehaviorIds []string) (ResolveResult, error)
	GetPbehaviors(ctx context.Context, t time.Time, pbehaviorIDs []string) ([]ResolveResult, error)
}

// typeResolver resolves entity state by computed data.
type typeResolver struct {
	// workerPoolSize restricts amount of goroutine which can be used during type resolving.
	workerPoolSize int
	// compute data only for this Span.
	Span timespan.Span
	// ComputedPbehaviors contains computed data for each pbehavior.
	ComputedPbehaviors map[string]ComputedPbehavior
	// DefaultActiveTypeID uses if there aren't any behaviors.
	DefaultActiveTypeID string
	// TypesByID contains all types.
	TypesByID map[string]Type

	logger zerolog.Logger
}

// NewTypeResolver creates new type resolver.
func NewTypeResolver(
	span timespan.Span,
	computedPbehaviors map[string]ComputedPbehavior,
	typesByID map[string]Type,
	defaultActiveTypeID string,
	logger zerolog.Logger,
) TypeResolver {
	return &typeResolver{
		Span:                span,
		ComputedPbehaviors:  computedPbehaviors,
		DefaultActiveTypeID: defaultActiveTypeID,
		TypesByID:           typesByID,
		workerPoolSize:      DefaultPoolSize,
		logger:              logger,
	}
}

// ResolveResult represents current state of entity.
type ResolveResult struct {
	ResolvedType      Type
	ResolvedPbhID     string
	ResolvedPbhName   string
	ResolvedPbhReason string
	ResolvedCreated   int64
}

// Resolve checks entity for each pbehavior concurrently. It uses "workerPoolSize" goroutines.
func (r *typeResolver) Resolve(
	ctx context.Context,
	t time.Time,
	entity types.Entity,
	cachedMatchedPbehaviorIds []string,
) (ResolveResult, error) {
	// Return error if time is out of timespan.
	if !r.Span.In(t) {
		return ResolveResult{}, ErrRecomputeNeed
	}

	workerCh := r.createWorkerChByEntity(ctx, entity, cachedMatchedPbehaviorIds)
	pbhRes, err := r.getPbehaviorIntervals(ctx, t, workerCh)
	if err != nil {
		return ResolveResult{}, err
	}

	// Use default active type if no pbehavior is in action for entity.
	res := ResolveResult{}

	if len(pbhRes) > 0 {
		res = pbhRes[0]
	}

	// Empty result represents default active type.
	if res.ResolvedType.ID != "" {
		activeType, ok := r.TypesByID[r.DefaultActiveTypeID]
		if !ok {
			return ResolveResult{}, fmt.Errorf("unknown type %v, probably need recompute data", r.DefaultActiveTypeID)
		}
		if res.ResolvedType == activeType {
			res = ResolveResult{}
		}
	}

	return res, nil
}

type workerData struct {
	id       string
	computed ComputedPbehavior
}

func (r *typeResolver) GetPbehaviors(
	ctx context.Context,
	t time.Time,
	pbehaviorIDs []string,
) ([]ResolveResult, error) {
	// Return error if time is out of timespan.
	if !r.Span.In(t) {
		return nil, ErrRecomputeNeed
	}

	workerCh := r.createWorkerCh(ctx, pbehaviorIDs)

	return r.getPbehaviorIntervals(ctx, t, workerCh)
}

func (r *typeResolver) getPbehaviorIntervals(
	ctx context.Context,
	t time.Time,
	workerCh <-chan workerData,
) ([]ResolveResult, error) {
	resCh := make(chan ResolveResult)

	g, ctx := errgroup.WithContext(ctx)

	for i := 0; i < r.workerPoolSize; i++ {
		g.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					return nil
				case d, ok := <-workerCh:
					if !ok {
						return nil
					}

					for _, computedType := range d.computed.Types {
						if !computedType.Span.In(t) {
							continue
						}

						resolvedType, ok := r.TypesByID[computedType.ID]
						if !ok {
							return fmt.Errorf("unknown type %v, probably need recompute data", computedType.ID)
						}

						resCh <- ResolveResult{
							ResolvedType:      resolvedType,
							ResolvedPbhID:     d.id,
							ResolvedPbhName:   d.computed.Name,
							ResolvedPbhReason: d.computed.Reason,
							ResolvedCreated:   d.computed.Created,
						}
						break
					}
				}
			}
		})
	}

	go func() {
		_ = g.Wait() // check error in wrapper func
		close(resCh)
	}()

	res := make([]ResolveResult, 0)
	for v := range resCh {
		res = append(res, v)
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].ResolvedType.Priority > res[j].ResolvedType.Priority ||
			res[i].ResolvedType.Priority == res[j].ResolvedType.Priority &&
				res[i].ResolvedCreated > res[j].ResolvedCreated ||
			res[i].ResolvedType.Priority == res[j].ResolvedType.Priority &&
				res[i].ResolvedCreated == res[j].ResolvedCreated &&
				res[i].ResolvedPbhID > res[j].ResolvedPbhID
	})

	return res, nil
}

// createWorkerCh creates worker ch and fills it.
func (r *typeResolver) createWorkerCh(ctx context.Context, pbehaviorIDs []string) <-chan workerData {
	workerCh := make(chan workerData)

	go func() {
		defer close(workerCh)

		for id, computed := range r.ComputedPbehaviors {
			select {
			case <-ctx.Done():
				return
			default:
				found := false
				for i := range pbehaviorIDs {
					if pbehaviorIDs[i] == id {
						found = true
						break
					}
				}

				if !found {
					continue
				}

				workerCh <- workerData{
					id:       id,
					computed: computed,
				}
			}
		}
	}()

	return workerCh
}

func (r *typeResolver) createWorkerChByEntity(ctx context.Context, entity types.Entity, cachedMatchedPbehaviorIds []string) <-chan workerData {
	workerCh := make(chan workerData)

	cachedPbehaviorIds := make(map[string]bool, len(cachedMatchedPbehaviorIds))
	for _, id := range cachedMatchedPbehaviorIds {
		cachedPbehaviorIds[id] = true
	}

	go func() {
		defer close(workerCh)

		for id, computed := range r.ComputedPbehaviors {
			select {
			case <-ctx.Done():
				return
			default:
				if len(computed.Patten) > 0 {
					ok, _, err := computed.Patten.Match(entity)
					if err != nil {
						r.logger.Err(err).Str("pbehavior", id).Msg("pbehavior has invalid pattern")
						continue
					}

					if ok {
						workerCh <- workerData{
							id:       id,
							computed: computed,
						}
					}

					continue
				}

				if cachedPbehaviorIds[id] {
					workerCh <- workerData{
						id:       id,
						computed: computed,
					}
				}
			}
		}
	}()

	return workerCh
}
