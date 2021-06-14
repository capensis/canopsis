package datastorage

//go:generate mockgen -destination=../../../mocks/lib/canopsis/datastorage/adapter.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage Adapter

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

const ID = "data_storage"

type Adapter interface {
	Get(ctx context.Context) (DataStorage, error)
	UpdateHistoryJunit(ctx context.Context, t types.CpsTime) error
}

type DataStorage struct {
	Config  Config  `bson:"config" json:"config"`
	History History `bson:"history" json:"history"`
}

type Config struct {
	Junit struct {
		DeleteAfter *types.DurationWithEnabled `bson:"delete_after,omitempty" json:"delete_after"`
	} `bson:"junit" json:"junit"`
}

type History struct {
	Junit *types.CpsTime `bson:"junit" json:"junit" swaggertype:"integer"`
}
