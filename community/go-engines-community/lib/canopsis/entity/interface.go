package entity

//go:generate mockgen -destination=../../../mocks/lib/canopsis/entity/adapter.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity Adapter

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

// Adapter ...
type Adapter interface {
	UpdateIdleFields(ctx context.Context, id string, idleSince *datetime.CpsTime, lastIdleRuleApply string) error

	GetAllWithLastUpdateDateBefore(ctx context.Context, time datetime.CpsTime, exclude []string) (mongo.Cursor, error)

	GetWithIdleSince(ctx context.Context) (mongo.Cursor, error)

	Bulk(ctx context.Context, models []mongodriver.WriteModel) error

	FindToCheckPbehaviorInfo(ctx context.Context, idsWithPbehaviors []string, exceptIds []string) (mongo.Cursor, error)
}
