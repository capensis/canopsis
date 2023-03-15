package config

import (
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

// Some default values related to configuration
const (
	ConfigKeyName        = "global_config"
	UserInterfaceKeyName = "user_interface"
	VersionKeyName       = "canopsis_version"
	RemediationKeyName   = "remediation"
	HealthCheckName      = "health_check"
	AlarmTagColorKeyName = "alarm_tag_color"
)

// SectionAlarm ...
type SectionAlarm struct {
	StealthyInterval     int    `toml:"StealthyInterval"`
	EnableLastEventDate  bool   `toml:"EnableLastEventDate"`
	CancelAutosolveDelay string `toml:"CancelAutosolveDelay"`
	DisplayNameScheme    string `toml:"DisplayNameScheme"`
	OutputLength         int    `toml:"OutputLength"`
	LongOutputLength     int    `toml:"LongOutputLength"`
	// DisableActionSnoozeDelayOnPbh ignores Pbh state to resolve snoozed with Action alarm while is True
	DisableActionSnoozeDelayOnPbh bool `toml:"DisableActionSnoozeDelayOnPbh"`
	// TimeToKeepResolvedAlarms defines how long resolved alarms will be kept in main alarm collection
	TimeToKeepResolvedAlarms string `toml:"TimeToKeepResolvedAlarms"`
	AllowDoubleAck           bool   `toml:"AllowDoubleAck"`
}

// SectionGlobal ...
type SectionGlobal struct {
	PrefetchCount                int   `toml:"PrefetchCount"`
	PrefetchSize                 int   `toml:"PrefetchSize"`
	ReconnectTimeoutMilliseconds int   `toml:"ReconnectTimeoutMilliseconds"`
	ReconnectRetries             int   `toml:"ReconnectRetries"`
	MaxExternalResponseSize      int64 `toml:"MaxExternalResponseSize"`
}

func (s *SectionGlobal) GetReconnectTimeout() time.Duration {
	return time.Duration(s.ReconnectTimeoutMilliseconds) * time.Millisecond
}

type SectionTimezone struct {
	Timezone string `toml:"Timezone"`
}

type SectionImportCtx struct {
	ThdWarnMinPerImport string `toml:"ThdWarnMinPerImport"`
	ThdCritMinPerImport string `toml:"ThdCritMinPerImport"`
	FilePattern         string `toml:"FilePattern"`
}

type SectionFile struct {
	Upload        string `toml:"Upload"`
	UploadMaxSize int64  `toml:"UploadMaxSize"`
	Junit         string `toml:"Junit"`
	JunitApi      string `toml:"JunitApi"`
}

type SectionDataStorage struct {
	TimeToExecute      string `toml:"TimeToExecute"`
	MaxUpdates         int    `toml:"MaxUpdates"`
	MongoClientTimeout string `toml:"MongoClientTimeout"`
}

type SectionApi struct {
	TokenSigningMethod string `toml:"TokenSigningMethod"`
	BulkMaxSize        int    `toml:"BulkMaxSize"`
}

type SectionLogger struct {
	Writer        string        `toml:"Writer"`
	ConsoleWriter ConsoleWriter `toml:"console_writer"`
}

type ConsoleWriter struct {
	Enabled    bool     `toml:"Enabled"`
	NoColor    bool     `toml:"NoColor"`
	TimeFormat string   `toml:"TimeFormat"`
	PartsOrder []string `toml:"PartsOrder"`
}

type SectionMetrics struct {
	FlushInterval             string `toml:"FlushInterval"`
	SliInterval               string `toml:"SliInterval"`
	EnabledManualInstructions bool   `toml:"EnabledManualInstructions"`
	EnabledNotAckedMetrics    bool   `toml:"EnabledNotAckedMetrics"`
}

type SectionTechMetrics struct {
	Enabled          bool   `toml:"Enabled"`
	DumpKeepInterval string `toml:"DumpKeepInterval"`
}

type SectionTemplate struct {
	Vars map[string]any `toml:"vars"`
}

// CanopsisConf represents a generic configuration object.
type CanopsisConf struct {
	ID          string             `bson:"_id,omitempty" toml:"omitempty"`
	Global      SectionGlobal      `bson:"global" toml:"global"`
	Alarm       SectionAlarm       `bson:"alarm" toml:"alarm"`
	Timezone    SectionTimezone    `bson:"timezone" toml:"timezone"`
	ImportCtx   SectionImportCtx   `bson:"import_ctx" toml:"import_ctx"`
	File        SectionFile        `bson:"file" toml:"file"`
	DataStorage SectionDataStorage `bson:"data_storage" toml:"data_storage"`
	Logger      SectionLogger      `bson:"logger" toml:"logger"`
	API         SectionApi         `bson:"api" toml:"api"`
	Metrics     SectionMetrics     `bson:"metrics" toml:"metrics"`
	TechMetrics SectionTechMetrics `bson:"tech_metrics" toml:"tech_metrics"`
	Template    SectionTemplate    `bson:"template" toml:"template"`
}

// UserInterfaceConf represents a user interface configuration object.
type UserInterfaceConf struct {
	IsAllowChangeSeverityToInfo bool `bson:"allow_change_severity_to_info"`
	// MaxMatchedItems need to warn user when number of items that match patterns is above this value
	MaxMatchedItems          int `bson:"max_matched_items"`
	CheckCountRequestTimeout int `bson:"check_count_request_timeout"`
}

type VersionConf struct {
	Edition string `bson:"edition"`
	Stack   string `bson:"stack"`

	Version        string         `bson:"version"`
	VersionUpdated *types.CpsTime `bson:"version_updated,omitempty"`
}
