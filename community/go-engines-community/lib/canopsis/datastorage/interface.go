package datastorage

//go:generate mockgen -destination=../../../mocks/lib/canopsis/datastorage/adapter.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage Adapter

import (
	"context"

	libtime "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/time"
)

const ID = "data_storage"

const BulkSize = 10000

type Adapter interface {
	Get(ctx context.Context) (DataStorage, error)
	UpdateHistoryJunit(ctx context.Context, t libtime.CpsTime) error
	UpdateHistoryRemediation(ctx context.Context, t libtime.CpsTime) error
	UpdateHistoryAlarm(ctx context.Context, history HistoryWithCount) error
	UpdateHistoryEntityDisabled(ctx context.Context, history HistoryWithCount) error
	UpdateHistoryEntityUnlinked(ctx context.Context, history HistoryWithCount) error
	UpdateHistoryEntityCleaned(ctx context.Context, history HistoryWithCount) error
	UpdateHistoryPbehavior(ctx context.Context, t libtime.CpsTime) error
	UpdateHistoryHealthCheck(ctx context.Context, t libtime.CpsTime) error
	UpdateHistoryWebhook(ctx context.Context, t libtime.CpsTime) error
	UpdateHistoryEventFilterFailure(ctx context.Context, t libtime.CpsTime) error
}

type DataStorage struct {
	Config  Config  `bson:"config" json:"config"`
	History History `bson:"history" json:"history"`
}

type Config struct {
	Junit struct {
		DeleteAfter *libtime.DurationWithEnabled `bson:"delete_after,omitempty" json:"delete_after"`
	} `bson:"junit" json:"junit"`
	Remediation struct {
		DeleteAfter         *libtime.DurationWithEnabled `bson:"delete_after,omitempty" json:"delete_after"`
		DeleteStatsAfter    *libtime.DurationWithEnabled `bson:"delete_stats_after,omitempty" json:"delete_stats_after"`
		DeleteModStatsAfter *libtime.DurationWithEnabled `bson:"delete_mod_stats_after,omitempty" json:"delete_mod_stats_after"`
	} `bson:"remediation" json:"remediation"`
	Alarm struct {
		ArchiveAfter *libtime.DurationWithEnabled `bson:"archive_after,omitempty" json:"archive_after"`
		DeleteAfter  *libtime.DurationWithEnabled `bson:"delete_after,omitempty" json:"delete_after"`
	} `bson:"alarm" json:"alarm"`
	Pbehavior struct {
		DeleteAfter *libtime.DurationWithEnabled `bson:"delete_after,omitempty" json:"delete_after"`
	} `bson:"pbehavior" json:"pbehavior"`
	HealthCheck struct {
		DeleteAfter *libtime.DurationWithEnabled `bson:"delete_after,omitempty" json:"delete_after"`
	} `bson:"health_check" json:"health_check"`
	Webhook struct {
		LogCredentials bool                         `bson:"log_credentials,omitempty" json:"log_credentials"`
		DeleteAfter    *libtime.DurationWithEnabled `bson:"delete_after,omitempty" json:"delete_after"`
	} `bson:"webhook" json:"webhook"`
	Metrics struct {
		DeleteAfter *libtime.DurationWithEnabled `bson:"delete_after,omitempty" json:"delete_after"`
	} `bson:"metrics" json:"metrics"`
	PerfDataMetrics struct {
		DeleteAfter *libtime.DurationWithEnabled `bson:"delete_after,omitempty" json:"delete_after"`
	} `bson:"perf_data_metrics" json:"perf_data_metrics"`
	EventFilterFailure struct {
		DeleteAfter *libtime.DurationWithEnabled `bson:"delete_after,omitempty" json:"delete_after"`
	} `bson:"event_filter_failure" json:"event_filter_failure"`
}

type History struct {
	Junit              *libtime.CpsTime  `bson:"junit" json:"junit" swaggertype:"integer"`
	Remediation        *libtime.CpsTime  `bson:"remediation" json:"remediation" swaggertype:"integer"`
	Alarm              *HistoryWithCount `bson:"alarm" json:"alarm"`
	EntityDisabled     *HistoryWithCount `bson:"entity_disabled" json:"entity_disabled"`
	EntityUnlinked     *HistoryWithCount `bson:"entity_unlinked" json:"entity_unlinked"`
	EntityCleaned      *HistoryWithCount `bson:"entity_cleaned" json:"entity_cleaned"`
	Pbehavior          *libtime.CpsTime  `bson:"pbehavior" json:"pbehavior" swaggertype:"integer"`
	HealthCheck        *libtime.CpsTime  `bson:"health_check" json:"health_check" swaggertype:"integer"`
	Webhook            *libtime.CpsTime  `bson:"webhook" json:"webhook" swaggertype:"integer"`
	EventFilterFailure *libtime.CpsTime  `bson:"event_filter_failure" json:"event_filter_failure" swaggertype:"integer"`
}

type HistoryWithCount struct {
	Time     libtime.CpsTime `bson:"time" json:"time" swaggertype:"integer"`
	Archived int64           `bson:"archived,omitempty" json:"archived"`
	Deleted  int64           `bson:"deleted,omitempty" json:"deleted"`
}
