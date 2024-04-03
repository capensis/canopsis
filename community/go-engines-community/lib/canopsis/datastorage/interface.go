package datastorage

//go:generate mockgen -destination=../../../mocks/lib/canopsis/datastorage/adapter.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage Adapter

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
)

const ID = "data_storage"

const BulkSize = 10000

type Adapter interface {
	Get(ctx context.Context) (DataStorage, error)
	UpdateHistoryJunit(ctx context.Context, t datetime.CpsTime) error
	UpdateHistoryRemediation(ctx context.Context, t datetime.CpsTime) error
	UpdateHistoryAlarm(ctx context.Context, history HistoryWithCount) error
	UpdateHistoryAlarmExternalTag(ctx context.Context, history HistoryWithCount) error
	UpdateHistoryEntityDisabled(ctx context.Context, history HistoryWithCount) error
	UpdateHistoryEntityUnlinked(ctx context.Context, history HistoryWithCount) error
	UpdateHistoryEntityCleaned(ctx context.Context, history HistoryWithCount) error
	UpdateHistoryPbehavior(ctx context.Context, t datetime.CpsTime) error
	UpdateHistoryHealthCheck(ctx context.Context, t datetime.CpsTime) error
	UpdateHistoryWebhook(ctx context.Context, t datetime.CpsTime) error
	UpdateHistoryEventFilterFailure(ctx context.Context, t datetime.CpsTime) error
}

type DataStorage struct {
	Config  Config  `bson:"config" json:"config"`
	History History `bson:"history" json:"history"`
}

type Config struct {
	Junit struct {
		DeleteAfter *datetime.DurationWithEnabled `bson:"delete_after,omitempty" json:"delete_after"`
	} `bson:"junit" json:"junit"`
	Remediation struct {
		DeleteAfter         *datetime.DurationWithEnabled `bson:"delete_after,omitempty" json:"delete_after"`
		DeleteStatsAfter    *datetime.DurationWithEnabled `bson:"delete_stats_after,omitempty" json:"delete_stats_after"`
		DeleteModStatsAfter *datetime.DurationWithEnabled `bson:"delete_mod_stats_after,omitempty" json:"delete_mod_stats_after"`
	} `bson:"remediation" json:"remediation"`
	Alarm struct {
		ArchiveAfter *datetime.DurationWithEnabled `bson:"archive_after,omitempty" json:"archive_after"`
		DeleteAfter  *datetime.DurationWithEnabled `bson:"delete_after,omitempty" json:"delete_after"`
	} `bson:"alarm" json:"alarm"`
	AlarmExternalTag struct {
		DeleteAfter *datetime.DurationWithEnabled `bson:"delete_after,omitempty" json:"delete_after"`
	} `bson:"alarm_external_tag" json:"alarm_external_tag"`
	Pbehavior struct {
		DeleteAfter *datetime.DurationWithEnabled `bson:"delete_after,omitempty" json:"delete_after"`
	} `bson:"pbehavior" json:"pbehavior"`
	HealthCheck struct {
		DeleteAfter *datetime.DurationWithEnabled `bson:"delete_after,omitempty" json:"delete_after"`
	} `bson:"health_check" json:"health_check"`
	Webhook struct {
		LogCredentials bool                          `bson:"log_credentials,omitempty" json:"log_credentials"`
		DeleteAfter    *datetime.DurationWithEnabled `bson:"delete_after,omitempty" json:"delete_after"`
	} `bson:"webhook" json:"webhook"`
	Metrics struct {
		DeleteAfter *datetime.DurationWithEnabled `bson:"delete_after,omitempty" json:"delete_after"`
	} `bson:"metrics" json:"metrics"`
	PerfDataMetrics struct {
		DeleteAfter *datetime.DurationWithEnabled `bson:"delete_after,omitempty" json:"delete_after"`
	} `bson:"perf_data_metrics" json:"perf_data_metrics"`
	EventFilterFailure struct {
		DeleteAfter *datetime.DurationWithEnabled `bson:"delete_after,omitempty" json:"delete_after"`
	} `bson:"event_filter_failure" json:"event_filter_failure"`
}

type History struct {
	Junit              *datetime.CpsTime `bson:"junit" json:"junit" swaggertype:"integer"`
	Remediation        *datetime.CpsTime `bson:"remediation" json:"remediation" swaggertype:"integer"`
	Alarm              *HistoryWithCount `bson:"alarm" json:"alarm"`
	AlarmExternalTag   *HistoryWithCount `bson:"alarm_external_tag" json:"alarm_external_tag"`
	EntityDisabled     *HistoryWithCount `bson:"entity_disabled" json:"entity_disabled"`
	EntityUnlinked     *HistoryWithCount `bson:"entity_unlinked" json:"entity_unlinked"`
	EntityCleaned      *HistoryWithCount `bson:"entity_cleaned" json:"entity_cleaned"`
	Pbehavior          *datetime.CpsTime `bson:"pbehavior" json:"pbehavior" swaggertype:"integer"`
	HealthCheck        *datetime.CpsTime `bson:"health_check" json:"health_check" swaggertype:"integer"`
	Webhook            *datetime.CpsTime `bson:"webhook" json:"webhook" swaggertype:"integer"`
	EventFilterFailure *datetime.CpsTime `bson:"event_filter_failure" json:"event_filter_failure" swaggertype:"integer"`
}

type HistoryWithCount struct {
	Time     datetime.CpsTime `bson:"time" json:"time" swaggertype:"integer"`
	Archived int64            `bson:"archived,omitempty" json:"archived"`
	Deleted  int64            `bson:"deleted,omitempty" json:"deleted"`
}
