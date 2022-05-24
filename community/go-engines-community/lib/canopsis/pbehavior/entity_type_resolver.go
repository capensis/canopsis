package pbehavior

import (
	"context"
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
	cachedMatcher ComputedEntityMatcher,
) EntityTypeResolver {
	return &entityTypeResolver{
		store:         store,
		matcher:       matcher,
		cachedMatcher: cachedMatcher,
	}
}

type entityTypeResolver struct {
	matcher       EntityMatcher
	cachedMatcher ComputedEntityMatcher

	store Store
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

	pbhIDs, err := r.cachedMatcher.Match(ctx, entity.ID)
	if err != nil {
		return ResolveResult{}, err
	}
	computed := ComputeResult{}
	// todo
	//if len(pbhIDs) == 0 {
	//	if entity.Created.Time.After(t.Add(-2 * canopsis.PeriodicalWaitTime)) {
	//		computed, err = r.store.GetComputed(ctx)
	//		if err != nil {
	//			return ResolveResult{}, err
	//		}
	//
	//		filters := make(map[string]string, len(computed.computedPbehaviors))
	//		for id, pbehavior := range computed.computedPbehaviors {
	//			filters[id] = pbehavior.Filter
	//		}
	//
	//		pbhIDs, err = r.matcher.MatchAll(ctx, entity.ID, filters)
	//		if err != nil {
	//			return ResolveResult{}, err
	//		}
	//	}
	//}
	//if len(pbhIDs) == 0 {
	//	return ResolveResult{}, nil
	//}

	if computed.defaultActiveType == "" {
		computed, err = r.store.GetComputedByIDs(ctx, pbhIDs)
		if err != nil {
			return ResolveResult{}, err
		}
	}

	resolver := NewTypeResolver(span, computed.computedPbehaviors, computed.typesByID, computed.defaultActiveType)

	return resolver.Resolve(ctx, t, pbhIDs)
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

	resolver := NewTypeResolver(span, data.computedPbehaviors, data.typesByID, data.defaultActiveType)
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
