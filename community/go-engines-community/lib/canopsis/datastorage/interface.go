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
	UpdateHistoryRemediation(ctx context.Context, t types.CpsTime) error
	UpdateHistoryAlarm(ctx context.Context, history AlarmHistory) error
	UpdateHistoryEntity(ctx context.Context, history EntityHistory) error
}

type DataStorage struct {
	Config  Config  `bson:"config" json:"config"`
	History History `bson:"history" json:"history"`
}

type Config struct {
	Junit struct {
		DeleteAfter *types.DurationWithEnabled `bson:"delete_after,omitempty" json:"delete_after"`
	} `bson:"junit" json:"junit"`
	Remediation struct {
		AccumulateAfter *types.DurationWithEnabled `bson:"accumulate_after,omitempty" json:"accumulate_after"`
		DeleteAfter     *types.DurationWithEnabled `bson:"delete_after,omitempty" json:"delete_after"`
	} `bson:"remediation" json:"remediation"`
	Alarm struct {
		ArchiveAfter *types.DurationWithEnabled `bson:"archive_after,omitempty" json:"archive_after"`
		DeleteAfter  *types.DurationWithEnabled `bson:"delete_after,omitempty" json:"delete_after"`
	} `bson:"alarm" json:"alarm"`
}

type History struct {
	Junit         *types.CpsTime `bson:"junit" json:"junit" swaggertype:"integer"`
	Remediation   *types.CpsTime `bson:"remediation" json:"remediation" swaggertype:"integer"`
	Alarm         *AlarmHistory  `bson:"alarm" json:"alarm"`
	EntityHistory *EntityHistory `bson:"entity" json:"entity"`
}

type AlarmHistory struct {
	Time           types.CpsTime `bson:"time" json:"time" swaggertype:"integer"`
	AlarmsArchived int64         `bson:"archived_alarms" json:"archived_alarms" swaggertype:"integer"`
	AlarmsDeleted  int64         `bson:"deleted_alarms" json:"deleted_alarms" swaggertype:"integer"`
}

type EntityHistory struct {
	Time             types.CpsTime `bson:"time" json:"time" swaggertype:"integer"`
	EntitiesArchived int64         `bson:"archived_entities" json:"archived_entities" swaggertype:"integer"`
	EntitiesDeleted  int64         `bson:"deleted_entities" json:"deleted_entities" swaggertype:"integer"`
}
