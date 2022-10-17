package pbehavior

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

// ComputedEntityTypeResolver uses data in memory to resolve type for an entity.
type ComputedEntityTypeResolver interface {
	Resolve(
		ctx context.Context,
		entity types.Entity,
		t time.Time,
	) (ResolveResult, error)
	GetComputedEntityIDs() ([]string, error)
	GetPbehaviorsCount(ctx context.Context, t time.Time) (int, error)
}

func NewComputedEntityTypeResolver(
	matcher ComputedEntityMatcher,
	resolver TypeResolver,
) ComputedEntityTypeResolver {
	return &computedEntityTypeResolver{
		matcher:  matcher,
		resolver: resolver,
	}
}

type computedEntityTypeResolver struct {
	matcher  ComputedEntityMatcher
	resolver TypeResolver
}

func (r *computedEntityTypeResolver) Resolve(
	ctx context.Context,
	entity types.Entity,
	t time.Time,
) (ResolveResult, error) {
	pbhIDs, err := r.matcher.Match(entity.ID)
	if err != nil {
		return ResolveResult{}, err
	}

	return r.resolver.Resolve(ctx, t, entity, pbhIDs)
}

func (r *computedEntityTypeResolver) GetComputedEntityIDs() ([]string, error) {
	return r.matcher.GetComputedEntityIDs()
}

func (r *computedEntityTypeResolver) GetPbehaviorsCount(ctx context.Context, t time.Time) (int, error) {
	return r.resolver.GetPbehaviorsCount(ctx, t)
}
