package pbehavior

import (
	"context"
	"github.com/rs/zerolog"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type EntityTypeResolver interface {
	Resolve(
		ctx context.Context,
		entity types.Entity,
		t time.Time,
	) (ResolveResult, error)
	GetPbehaviors(ctx context.Context, pbhIDs []string, t time.Time) (map[string]ResolveResult, error)
}

func NewEntityTypeResolver(
	store Store,
	matcher EntityMatcher,
	logger zerolog.Logger,
) EntityTypeResolver {
	return &entityTypeResolver{
		store:   store,
		matcher: matcher,
		logger:  logger,
	}
}

type entityTypeResolver struct {
	matcher EntityMatcher
	store   Store
	logger  zerolog.Logger
}

func (r *entityTypeResolver) Resolve(
	ctx context.Context,
	entity types.Entity,
	t time.Time,
) (ResolveResult, error) {
	span, err := r.store.GetSpan(ctx)
	if err != nil {
		return ResolveResult{}, err
	}

	if !span.In(t) {
		return ResolveResult{}, ErrRecomputeNeed
	}

	computed, err := r.store.GetComputed(ctx)
	if err != nil {
		return ResolveResult{}, err
	}

	filters := make(map[string]interface{}, len(computed.computedPbehaviors))
	for id, pbehavior := range computed.computedPbehaviors {
		if len(pbehavior.OldMongoQuery) > 0 {
			filters[id] = pbehavior.OldMongoQuery
		}
	}

	pbhIDs, err := r.matcher.MatchAll(ctx, entity.ID, filters)
	if err != nil {
		return ResolveResult{}, err
	}

	resolver := NewTypeResolver(span, computed.computedPbehaviors, computed.typesByID, computed.defaultActiveType, r.logger)

	return resolver.Resolve(ctx, t, entity, pbhIDs)
}

func (r *entityTypeResolver) GetPbehaviors(ctx context.Context, pbhIDs []string, t time.Time) (map[string]ResolveResult, error) {
	span, err := r.store.GetSpan(ctx)
	if err != nil {
		return nil, err
	}

	if !span.In(t) {
		return nil, ErrRecomputeNeed
	}

	data, err := r.store.GetComputedByIDs(ctx, pbhIDs)
	if err != nil {
		return nil, err
	}

	resolver := NewTypeResolver(span, data.computedPbehaviors, data.typesByID, data.defaultActiveType, r.logger)
	pbhs, err := resolver.GetPbehaviors(ctx, t, pbhIDs)
	if err != nil {
		return nil, err
	}

	res := make(map[string]ResolveResult, len(pbhs))
	for _, pbh := range pbhs {
		res[pbh.ResolvedPbhID] = pbh
	}

	return res, nil
}
