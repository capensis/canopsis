package entityservice

//go:generate go tool go.uber.org/mock/mockgen -destination=../../../mocks/lib/canopsis/entityservice/entityservice.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice Adapter

import (
	"context"

	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
)

// Adapter is an interface that provides methods for database queries regarding
// services and their dependencies.
type Adapter interface {
	GetValid(ctx context.Context) ([]EntityService, error)

	// UpdateBulk bulk update
	UpdateBulk(ctx context.Context, writeModels []mongodriver.WriteModel) error
}

type IdleSinceService interface {
	RecomputeIdleSince(parentCtx context.Context) error
}
