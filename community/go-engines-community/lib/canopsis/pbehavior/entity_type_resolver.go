package pbehavior

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/rs/zerolog"
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
	logger zerolog.Logger,
) EntityTypeResolver {
	return &entityTypeResolver{
		store:  store,
		logger: logger,
	}
}

type entityTypeResolver struct {
	store  Store
	logger zerolog.Logger
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

	resolver := NewTypeResolver(
		span,
		computed.ComputedPbehaviors,
		computed.TypesByID,
		computed.DefaultActiveType,
		r.logger,
	)

	return resolver.Resolve(ctx, t, entity)
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

	resolver := NewTypeResolver(
		span,
		data.ComputedPbehaviors,
		data.TypesByID,
		data.DefaultActiveType,
		r.logger,
	)
	pbhs, err := resolver.GetPbehaviors(ctx, t, pbhIDs)
	if err != nil {
		return nil, err
	}

	res := make(map[string]ResolveResult, len(pbhs))
	for _, pbh := range pbhs {
		res[pbh.ID] = pbh
	}

	return res, nil
}
