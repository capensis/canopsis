package datastorage

//go:generate go tool go.uber.org/mock/mockgen -destination=../../../mocks/lib/canopsis/datastorage/adapter.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage Adapter,Cleaner

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
)

const ID = "data_storage"

const BulkSize = 10000

type Adapter interface {
	Get(ctx context.Context) (DataStorage, error)
	UpdateHistoryEntityDisabled(ctx context.Context, history HistoryWithCount) error
	UpdateHistoryEntityUnlinked(ctx context.Context, history HistoryWithCount) error
	UpdateHistoryEntityCleaned(ctx context.Context, history HistoryWithCount) error
}

type DataStorage struct {
	Config  Config                      `bson:"config" json:"config"`
	History map[string]HistoryWithCount `bson:"history" json:"history"`
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
	Entity struct {
		ArchiveAfter *datetime.DurationWithEnabled `bson:"archive_after,omitempty" json:"archive_after"`
	} `bson:"entity" json:"entity"`
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
	EventRecords struct {
		DeleteAfter *datetime.DurationWithEnabled `bson:"delete_after,omitempty" json:"delete_after"`
	} `bson:"event_records" json:"event_records"`
	EntityInfosLog struct {
		DeleteAfter *datetime.DurationWithEnabled `bson:"delete_after,omitempty" json:"delete_after"`
	} `bson:"entity_infos_log" json:"entity_infos_log"`
	LLMChat struct {
		DeleteAfter *datetime.DurationWithEnabled `bson:"delete_after,omitempty" json:"delete_after"`
	} `bson:"llm_chat" json:"llm_chat"`
}

type HistoryWithCount struct {
	Time     datetime.CpsTime `bson:"time" json:"time" swaggertype:"integer"`
	Archived int64            `bson:"archived,omitempty" json:"archived"`
	Deleted  int64            `bson:"deleted,omitempty" json:"deleted"`
}
