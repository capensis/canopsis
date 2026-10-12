package token

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
)

type Store interface {
	Save(ctx context.Context, token Token) error
	Delete(ctx context.Context, id string) (bool, error)
	DeleteBy(ctx context.Context, user, provider string) error
	DeleteByUserIDs(ctx context.Context, ids []string) error
}

type Token struct {
	ID       string           `bson:"_id"`
	User     string           `bson:"user"`
	Provider string           `bson:"provider,omitempty"`
	Created  datetime.CpsTime `bson:"created"`
	Accessed datetime.CpsTime `bson:"accessed"`

	Expired             *datetime.CpsTime          `bson:"expired,omitempty"`
	ExpiredByInactivity *datetime.CpsTime          `bson:"expired_by_inactivity,omitempty"`
	MaxInactiveInterval *datetime.DurationWithUnit `bson:"max_inactive_interval,omitempty"`
}
