package datastorage

//go:generate mockgen -destination=../../../mocks/lib/canopsis/datastorage/adapter.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage Adapter

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

const ID = "data_storage"

const BulkSize = 10000

type Adapter interface {
	Get(ctx context.Context) (DataStorage, error)
	UpdateHistoryJunit(ctx context.Context, t types.CpsTime) error
	UpdateHistoryRemediation(ctx context.Context, t types.CpsTime) error
	UpdateHistoryAlarm(ctx context.Context, history HistoryWithCount) error
	UpdateHistoryEntity(ctx context.Context, history HistoryWithCount) error
	UpdateHistoryPbehavior(ctx context.Context, t types.CpsTime) error
	UpdateHistoryHealthCheck(ctx context.Context, t types.CpsTime) error
	UpdateHistoryWebhook(ctx context.Context, t types.CpsTime) error
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
		DeleteAfter         *types.DurationWithEnabled `bson:"delete_after,omitempty" json:"delete_after"`
		DeleteStatsAfter    *types.DurationWithEnabled `bson:"delete_stats_after,omitempty" json:"delete_stats_after"`
		DeleteModStatsAfter *types.DurationWithEnabled `bson:"delete_mod_stats_after,omitempty" json:"delete_mod_stats_after"`
	} `bson:"remediation" json:"remediation"`
	Alarm struct {
		ArchiveAfter *types.DurationWithEnabled `bson:"archive_after,omitempty" json:"archive_after"`
		DeleteAfter  *types.DurationWithEnabled `bson:"delete_after,omitempty" json:"delete_after"`
	} `bson:"alarm" json:"alarm"`
	Pbehavior struct {
		DeleteAfter *types.DurationWithEnabled `bson:"delete_after,omitempty" json:"delete_after"`
	} `bson:"pbehavior" json:"pbehavior"`
	HealthCheck struct {
		DeleteAfter *types.DurationWithEnabled `bson:"delete_after,omitempty" json:"delete_after"`
	} `bson:"health_check" json:"health_check"`
	Webhook struct {
		LogCredentials bool                       `bson:"log_credentials,omitempty" json:"log_credentials"`
		DeleteAfter    *types.DurationWithEnabled `bson:"delete_after,omitempty" json:"delete_after"`
	} `bson:"webhook" json:"webhook"`
}

type History struct {
	Junit       *types.CpsTime    `bson:"junit" json:"junit" swaggertype:"integer"`
	Remediation *types.CpsTime    `bson:"remediation" json:"remediation" swaggertype:"integer"`
	Alarm       *HistoryWithCount `bson:"alarm" json:"alarm"`
	Entity      *HistoryWithCount `bson:"entity" json:"entity"`
	Pbehavior   *types.CpsTime    `bson:"pbehavior" json:"pbehavior" swaggertype:"integer"`
	HealthCheck *types.CpsTime    `bson:"health_check" json:"health_check" swaggertype:"integer"`
	Webhook     *types.CpsTime    `bson:"webhook" json:"webhook" swaggertype:"integer"`
}

type HistoryWithCount struct {
	Time     types.CpsTime `bson:"time" json:"time" swaggertype:"integer"`
	Archived int64         `bson:"archived" json:"archived"`
	Deleted  int64         `bson:"deleted" json:"deleted"`
}
