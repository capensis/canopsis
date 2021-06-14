package pbehavior

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	libtypes "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/timespan"
	"github.com/rs/zerolog"
)

var ErrRecomputeNeed = errors.New("provided time is out of computed date, probably need recompute data")

// TypeResolver figures out in which state provided entity at the moment is.
type TypeResolver interface {
	// Resolve returns current type for entity if there is corresponding periodical behavior.
	// Otherwise it returns default active type.
	Resolve(context.Context, *libtypes.Entity, time.Time) (ResolveResult, error)
	GetSpan() timespan.Span
}

// typeResolver resolves entity state by computed data.
type typeResolver struct {
	// matcher checks if entity matches pbehavior filter.
	matcher EntityMatcher
	// workerPoolSize restricts amount of goroutine which can be used during type resolving.
	workerPoolSize int
	// compute data only for this Span.
	Span timespan.Span `json:"span"`
	// ComputedPbehaviors contains computed data for each pbehavior.
	ComputedPbehaviors map[string]ComputedPbehavior `json:"computed"`
	// DefaultActiveTypeID uses if there aren't any behaviors.
	DefaultActiveTypeID string `json:"default_active_type"`
	// TypesByID contains all types.
	TypesByID map[string]*Type `json:"types"`
	logger    zerolog.Logger
}

// NewTypeResolver creates new type resolver.
func NewTypeResolver(
	matcher EntityMatcher,
	span timespan.Span,
	computedPbehaviors map[string]ComputedPbehavior,
	typesByID map[string]*Type,
	defaultActiveTypeID string,
	logger zerolog.Logger,
	workerPoolSize ...int,
) *typeResolver {
	poolSize := DefaultPoolSize

	if len(workerPoolSize) == 1 {
		poolSize = workerPoolSize[0]
	} else if len(workerPoolSize) > 1 {
		panic("too much arguments")
	}

	return &typeResolver{
		matcher:             matcher,
		Span:                span,
		ComputedPbehaviors:  computedPbehaviors,
		DefaultActiveTypeID: defaultActiveTypeID,
		TypesByID:           typesByID,
		logger:              logger,
		workerPoolSize:      poolSize,
	}
}

// ResolveResult represents current state of entity.
type ResolveResult struct {
	ResolvedType      *Type
	ResolvedPbhID     string
	ResolvedPbhName   string
	ResolvedPbhReason string
	ResolvedCreated   int64
}

// resolveResult uses for concurrency in Resolve.
type resolveResult struct {
	ResolveResult
	Err error
}

// Resolve checks entity for each pbehavior concurrently. It uses "workerPoolSize" goroutines.
func (r *typeResolver) Resolve(
	parentCtx context.Context,
	entity *libtypes.Entity,
	t time.Time,
) (ResolveResult, error) {
	// Return error if time is out of timespan.
	if !r.Span.In(t) {
		return ResolveResult{}, ErrRecomputeNeed
	}

	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	matchEntityByPbehavior, err := r.matchEntity(ctx, entity)
	if err != nil {
		return ResolveResult{}, err
	}

	resCh := r.runWorkers(ctx, t, matchEntityByPbehavior, nil)
	// Use default active type if no pbehavior is in action for entity.
	res := ResolveResult{}

	// Get results from result channel and select result with max type priority.
	for v := range resCh {
		if v.Err == nil {
			if res.ResolvedType == nil ||
				res.ResolvedType.Priority < v.ResolvedType.Priority ||
				res.ResolvedType.Priority == v.ResolvedType.Priority &&
					res.ResolvedCreated < v.ResolvedCreated ||
				res.ResolvedType.Priority == v.ResolvedType.Priority &&
					res.ResolvedCreated == v.ResolvedCreated &&
					res.ResolvedPbhID < v.ResolvedPbhID {
				res = v.ResolveResult
			}
		} else {
			err = v.Err
		}
	}

	if err != nil {
		return ResolveResult{}, err
	}

	// Empty result represents default active type.
	if res.ResolvedType != nil {
		activeType, ok := r.TypesByID[r.DefaultActiveTypeID]
		if !ok {
			return ResolveResult{}, fmt.Errorf("unknown type %v, probably need recompute data", r.DefaultActiveTypeID)
		}
		if *res.ResolvedType == *activeType {
			res = ResolveResult{}
			r.logger.Debug().Msgf("resolve default active type for entity %v", entity.ID)
		} else {
			r.logger.Debug().Msgf("resolve type %v for entity %v", res.ResolvedType.ID, entity.ID)
		}
	}

	return res, nil
}

func (r *typeResolver) GetSpan() timespan.Span {
	return r.Span
}

func (r *typeResolver) UpdateData(pbehaviorID string, c *ComputedPbehavior) {
	if _, ok := r.ComputedPbehaviors[pbehaviorID]; ok && c == nil {
		delete(r.ComputedPbehaviors, pbehaviorID)
	} else if c != nil {
		r.ComputedPbehaviors[pbehaviorID] = *c
	}
}

// GetPbehaviorStatus returns true if pbehavior is in action for provided time.
func (r *typeResolver) GetPbehaviorStatus(
	ctx context.Context,
	pbehaviorIDs []string,
	t time.Time,
) (map[string]bool, error) {
	// Return error if time is out of timespan.
	if !r.Span.In(t) {
		return nil, ErrRecomputeNeed
	}

	res := make(map[string]bool, len(pbehaviorIDs))
	for _, id := range pbehaviorIDs {
		res[id] = false
	}

	resCh := r.runWorkers(ctx, t, nil, pbehaviorIDs)
	var err error
	for v := range resCh {
		if v.Err == nil {
			res[v.ResolvedPbhID] = true
		} else {
			err = v.Err
		}
	}

	if err != nil {
		return nil, err
	}

	return res, nil
}

// matchEntity resolves if entity matches filter of each pbehavior.
func (r *typeResolver) matchEntity(ctx context.Context, entity *libtypes.Entity) (map[string]bool, error) {
	filters := make(map[string]string, len(r.ComputedPbehaviors))
	for id, computed := range r.ComputedPbehaviors {
		filters[id] = computed.Filter
	}

	matchEntityByPbehavior, err := r.matcher.MatchAll(ctx, entity.ID, filters)
	if err != nil {
		return nil, err
	}

	return matchEntityByPbehavior, nil
}

type workerData struct {
	id       string
	computed ComputedPbehavior
}

// runWorkers creates res chan and runs workers to fill it.
func (r *typeResolver) runWorkers(
	ctx context.Context,
	t time.Time,
	matchEntityByPbehavior map[string]bool,
	pbehaviorIDs []string,
) <-chan resolveResult {
	resCh := make(chan resolveResult)
	workerCh := r.createWorkerCh(ctx, pbehaviorIDs)

	go func() {
		wg := sync.WaitGroup{}
		defer close(resCh)

		for i := 0; i < r.workerPoolSize; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				workerCtx, cancel := context.WithCancel(ctx)
				defer cancel()
				for {
					select {
					case <-workerCtx.Done():
						resCh <- resolveResult{Err: workerCtx.Err()}
						return
					case d, ok := <-workerCh:
						if !ok {
							return
						}

						if matchEntityByPbehavior != nil {
							if m, ok := matchEntityByPbehavior[d.id]; !ok || !m {
								continue
							}
						}

						for _, computedType := range d.computed.Types {
							if !computedType.Span.In(t) {
								continue
							}

							resolvedType, ok := r.TypesByID[computedType.ID]
							if !ok {
								resCh <- resolveResult{Err: fmt.Errorf("unknown type %v, probably need recompute data", computedType.ID)}
								return
							}

							resCh <- resolveResult{
								ResolveResult: ResolveResult{
									ResolvedType:      resolvedType,
									ResolvedPbhID:     d.id,
									ResolvedPbhName:   d.computed.Name,
									ResolvedPbhReason: d.computed.Reason,
									ResolvedCreated:   d.computed.Created,
								},
							}
							break
						}
					}
				}
			}()
		}

		wg.Wait()
	}()

	return resCh
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
				if pbehaviorIDs != nil {
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
