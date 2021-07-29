package pbehavior

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	libtypes "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/timespan"
	"golang.org/x/sync/errgroup"
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
}

// NewTypeResolver creates new type resolver.
func NewTypeResolver(
	matcher EntityMatcher,
	span timespan.Span,
	computedPbehaviors map[string]ComputedPbehavior,
	typesByID map[string]*Type,
	defaultActiveTypeID string,
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

// Resolve checks entity for each pbehavior concurrently. It uses "workerPoolSize" goroutines.
func (r *typeResolver) Resolve(
	ctx context.Context,
	entity *libtypes.Entity,
	t time.Time,
) (ResolveResult, error) {
	// Return error if time is out of timespan.
	if !r.Span.In(t) {
		return ResolveResult{}, ErrRecomputeNeed
	}

	pbhRes, err := r.runWorkers(ctx, t, nil)
	if err != nil {
		return ResolveResult{}, err
	}

	pbhIDs := make([]string, len(pbhRes))
	for i, v := range pbhRes {
		pbhIDs[i] = v.ResolvedPbhID
	}

	matchEntityByPbehavior, err := r.matchEntity(ctx, entity, pbhIDs)
	if err != nil {
		return ResolveResult{}, err
	}

	// Use default active type if no pbehavior is in action for entity.
	res := ResolveResult{}

	for _, v := range pbhRes {
		if matchEntityByPbehavior[v.ResolvedPbhID] {
			res = v
			break
		}
	}

	// Empty result represents default active type.
	if res.ResolvedType != nil {
		activeType, ok := r.TypesByID[r.DefaultActiveTypeID]
		if !ok {
			return ResolveResult{}, fmt.Errorf("unknown type %v, probably need recompute data", r.DefaultActiveTypeID)
		}
		if *res.ResolvedType == *activeType {
			res = ResolveResult{}
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

	pbhRes, err := r.runWorkers(ctx, t, pbehaviorIDs)
	if err != nil {
		return nil, err
	}

	for _, v := range pbhRes {
		res[v.ResolvedPbhID] = true
	}

	return res, nil
}

func (r *typeResolver) GetComputedPbehaviorsCount() int {
	return len(r.ComputedPbehaviors)
}

// matchEntity resolves if entity matches filter of each pbehavior.
func (r *typeResolver) matchEntity(ctx context.Context, entity *libtypes.Entity, pbhIDs []string) (map[string]bool, error) {
	if len(pbhIDs) == 0 {
		return nil, nil
	}

	filters := make(map[string]string, len(pbhIDs))
	for _, id := range pbhIDs {
		filters[id] = r.ComputedPbehaviors[id].Filter
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
	pbehaviorIDs []string,
) ([]ResolveResult, error) {
	resCh := make(chan ResolveResult)
	workerCh := r.createWorkerCh(ctx, pbehaviorIDs)

	g, ctx := errgroup.WithContext(ctx)

	for i := 0; i < r.workerPoolSize; i++ {
		g.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
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
