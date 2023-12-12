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
	getter ComputedEntityGetter,
	resolver TypeResolver,
) ComputedEntityTypeResolver {
	return &computedEntityTypeResolver{
		getter:   getter,
		resolver: resolver,
	}
}

type computedEntityTypeResolver struct {
	getter   ComputedEntityGetter
	resolver TypeResolver
}

func (r *computedEntityTypeResolver) Resolve(
	ctx context.Context,
	entity types.Entity,
	t time.Time,
) (ResolveResult, error) {
	return r.resolver.Resolve(ctx, t, entity)
}

func (r *computedEntityTypeResolver) GetComputedEntityIDs() ([]string, error) {
	return r.getter.GetComputedEntityIDs()
}

func (r *computedEntityTypeResolver) GetPbehaviorsCount(ctx context.Context, t time.Time) (int, error) {
	return r.resolver.GetPbehaviorsCount(ctx, t)
}
