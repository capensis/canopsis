package pbehavior

import (
	"context"
	"fmt"
	"sort"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/timespan"
	"golang.org/x/sync/errgroup"
)

// TypeResolver figures out in which state provided entity at the moment is.
type TypeResolver interface {
	// Resolve returns current type for entity if there is corresponding periodical behavior.
	// Otherwise it returns default active type.
	Resolve(context.Context, time.Time, []string) (ResolveResult, error)
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
}

// NewTypeResolver creates new type resolver.
func NewTypeResolver(
	span timespan.Span,
	computedPbehaviors map[string]ComputedPbehavior,
	typesByID map[string]Type,
	defaultActiveTypeID string,
) TypeResolver {
	return &typeResolver{
		Span:                span,
		ComputedPbehaviors:  computedPbehaviors,
		DefaultActiveTypeID: defaultActiveTypeID,
		TypesByID:           typesByID,
		workerPoolSize:      DefaultPoolSize,
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
	t time.Time,
	pbhIDs []string,
) (ResolveResult, error) {
	if !r.Span.In(t) {
		return ResolveResult{}, ErrRecomputeNeed
	}

	pbhRes, err := r.GetPbehaviors(ctx, t, pbhIDs)
	if err != nil {
		return ResolveResult{}, err
	}

	// Use default active type if no pbehavior is in action for entity.
	res := ResolveResult{}

	if len(pbhRes) > 0 {
		res = pbhRes[0]
	}

	// Empty result represents default active type.
	if res.ResolvedType != nil {
		activeType, ok := r.TypesByID[r.DefaultActiveTypeID]
		if !ok {
			return ResolveResult{}, fmt.Errorf("unknown type %v, probably need recompute data", r.DefaultActiveTypeID)
		}
		if *res.ResolvedType == activeType {
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
	if !r.Span.In(t) {
		return nil, ErrRecomputeNeed
	}

	resCh := make(chan ResolveResult)
	workerCh := make(chan workerData)
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		defer close(workerCh)

		for id, computed := range r.ComputedPbehaviors {
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

			select {
			case <-ctx.Done():
				return nil
			case workerCh <- workerData{
				id:       id,
				computed: computed,
			}:
			}
		}

		return nil
	})

	for i := 0; i < r.workerPoolSize; i++ {
		g.Go(func() error {
			for d := range workerCh {
				for _, computedType := range d.computed.Types {
					if !computedType.Span.In(t) {
						continue
					}

					resolvedType, ok := r.TypesByID[computedType.ID]
					if !ok {
						return fmt.Errorf("unknown type %v, probably need recompute data", computedType.ID)
					}

					resCh <- ResolveResult{
						ResolvedType:      &resolvedType,
						ResolvedPbhID:     d.id,
						ResolvedPbhName:   d.computed.Name,
						ResolvedPbhReason: d.computed.Reason,
						ResolvedCreated:   d.computed.Created,
					}
					break
				}
			}

			return nil
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
