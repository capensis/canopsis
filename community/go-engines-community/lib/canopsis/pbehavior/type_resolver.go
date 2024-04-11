package pbehavior

import (
	"context"
	"fmt"
	"sort"
	"sync/atomic"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/match"
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
	Resolve(ctx context.Context, t time.Time, entity types.Entity) (ResolveResult, error)
	GetPbehaviors(ctx context.Context, t time.Time, pbehaviorIDs []string) ([]ResolveResult, error)
	GetPbehaviorsCount(ctx context.Context, t time.Time) (int, error)
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
	Type       Type
	ID         string
	Name       string
	ReasonID   string
	ReasonName string
	Color      string
	Created    int64
}

// Resolve checks entity for each pbehavior concurrently. It uses "workerPoolSize" goroutines.
func (r *typeResolver) Resolve(
	ctx context.Context,
	t time.Time,
	entity types.Entity,
) (ResolveResult, error) {
	if !r.Span.In(t) {
		return ResolveResult{}, ErrRecomputeNeed
	}

	pbhRes, err := r.getPbehaviorIntervals(ctx, t, func(id string, computed ComputedPbehavior) bool {
		if len(computed.EntityPattern) == 0 {
			return false
		}

		matched, err := match.MatchEntityPattern(computed.EntityPattern, &entity)
		if err != nil {
			r.logger.Err(err).Str("pbehavior", id).Msg("pbehavior has invalid pattern")
			return false
		}

		return matched
	})
	if err != nil {
		return ResolveResult{}, err
	}

	// Use default active type if no pbehavior is in action for entity.
	res := ResolveResult{}

	if len(pbhRes) > 0 {
		res = pbhRes[0]
	}

	// Empty result represents default active type.
	if res.Type.ID != "" {
		activeType, ok := r.TypesByID[r.DefaultActiveTypeID]
		if !ok {
			return ResolveResult{}, fmt.Errorf("unknown type %v, probably need recompute data", r.DefaultActiveTypeID)
		}
		if res.Type == activeType {
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

	return r.getPbehaviorIntervals(ctx, t, func(id string, _ ComputedPbehavior) bool {
		for i := range pbehaviorIDs {
			if pbehaviorIDs[i] == id {
				return true
			}
		}

		return false
	})
}

func (r *typeResolver) GetPbehaviorsCount(ctx context.Context, t time.Time) (int, error) {
	if !r.Span.In(t) {
		return 0, ErrRecomputeNeed
	}

	g, ctx := errgroup.WithContext(ctx)
	workerCh := make(chan ComputedPbehavior)
	var count int64

	g.Go(func() error {
		defer close(workerCh)

		for _, computed := range r.ComputedPbehaviors {
			select {
			case <-ctx.Done():
				return nil
			case workerCh <- computed:
			}
		}

		return nil
	})

	for i := 0; i < r.workerPoolSize; i++ {
		g.Go(func() error {
			for d := range workerCh {
				for _, computedType := range d.Types {
					if computedType.Span.In(t) {
						atomic.AddInt64(&count, 1)
						break
					}
				}
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return 0, err
	}

	return int(count), nil
}

func (r *typeResolver) getPbehaviorIntervals(
	ctx context.Context,
	t time.Time,
	f func(id string, c ComputedPbehavior) bool,
) ([]ResolveResult, error) {
	resCh := make(chan ResolveResult)
	workerCh := make(chan workerData)
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		defer close(workerCh)

		for id, computed := range r.ComputedPbehaviors {
			if !f(id, computed) {
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
						Type:       resolvedType,
						ID:         d.id,
						Name:       d.computed.Name,
						ReasonID:   d.computed.ReasonID,
						ReasonName: d.computed.ReasonName,
						Color:      d.computed.Color,
						Created:    d.computed.Created,
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
		return res[i].Type.Priority > res[j].Type.Priority ||
			res[i].Type.Priority == res[j].Type.Priority &&
				res[i].Created > res[j].Created ||
			res[i].Type.Priority == res[j].Type.Priority &&
				res[i].Created == res[j].Created &&
				res[i].ID > res[j].ID
	})

	return res, nil
}
