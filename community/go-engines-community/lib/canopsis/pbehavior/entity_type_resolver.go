package pbehavior

import (
	"context"
	"time"
)

type EntityTypeResolver interface {
	Resolve(
		ctx context.Context,
		entityID string,
		t time.Time,
	) (ResolveResult, error)
	ResolveAll(
		ctx context.Context,
		entityIDs []string,
		t time.Time,
	) (map[string]ResolveResult, error)
	GetPbehaviors(ctx context.Context, pbhIDs []string, t time.Time) (map[string]ResolveResult, error)
}

func NewEntityTypeResolver(
	store Store,
	matcher ComputedEntityMatcher,
) EntityTypeResolver {
	return &entityTypeResolver{
		store:   store,
		matcher: matcher,
	}
}

type entityTypeResolver struct {
	matcher ComputedEntityMatcher
	store   Store
}

func (r *entityTypeResolver) Resolve(
	ctx context.Context,
	entityID string,
	t time.Time,
) (ResolveResult, error) {
	span, err := r.store.GetSpan(ctx)
	if err != nil {
		return ResolveResult{}, err
	}

	if !span.In(t) {
		return ResolveResult{}, ErrRecomputeNeed
	}

	pbhIDs, err := r.matcher.Match(ctx, entityID)
	if err != nil || len(pbhIDs) == 0 {
		return ResolveResult{}, err
	}

	computed, err := r.store.GetComputedByIDs(ctx, pbhIDs)
	if err != nil {
		return ResolveResult{}, err
	}

	resolver := NewTypeResolver(span, computed.computedPbehaviors, computed.typesByID, computed.defaultActiveType)

	return resolver.Resolve(ctx, t, pbhIDs)
}

func (r *entityTypeResolver) ResolveAll(
	ctx context.Context,
	entityIDs []string,
	t time.Time,
) (map[string]ResolveResult, error) {
	span, err := r.store.GetSpan(ctx)
	if err != nil {
		return nil, err
	}

	if !span.In(t) {
		return nil, ErrRecomputeNeed
	}

	pbhIDsByEntityID, err := r.matcher.MatchAll(ctx, entityIDs)
	if err != nil || len(pbhIDsByEntityID) == 0 {
		return nil, err
	}

	pbhIDs := make([]string, 0)
	for _, v := range pbhIDsByEntityID {
		pbhIDs = append(pbhIDs, v...)
	}

	computed, err := r.store.GetComputedByIDs(ctx, pbhIDs)
	if err != nil {
		return nil, err
	}

	resolver := NewTypeResolver(span, computed.computedPbehaviors, computed.typesByID, computed.defaultActiveType)
	resolveResult := make(map[string]ResolveResult, len(entityIDs))
	for _, entityID := range entityIDs {
		result, err := resolver.Resolve(ctx, t, pbhIDsByEntityID[entityID])
		if err != nil {
			return nil, err
		}

		resolveResult[entityID] = result
	}

	return resolveResult, nil
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
