package config

import (
	"time"
)

// Some default values related to configuration
const (
	ConfigKeyName        = "global_config"
	UserInterfaceKeyName = "user_interface"
	VersionKeyName       = "canopsis_version"
	RemediationKeyName   = "remediation"
	HealthCheckName      = "health_check"
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
	PrefetchCount                int `toml:"PrefetchCount"`
	PrefetchSize                 int `toml:"PrefetchSize"`
	ReconnectTimeoutMilliseconds int `toml:"ReconnectTimeoutMilliseconds"`
	ReconnectRetries             int `toml:"ReconnectRetries"`
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
	TokenExpiration    string `toml:"TokenExpiration"`
	TokenSigningMethod string `toml:"TokenSigningMethod"`
	BulkMaxSize        int    `toml:"BulkMaxSize"`
	ExportBulkSize     int    `toml:"ExportBulkSize"`
}

type SectionMetrics struct {
	FlushInterval string `toml:"FlushInterval"`
	SliInterval   string `toml:"SliInterval"`
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
	API         SectionApi         `bson:"api" toml:"api"`
	Metrics     SectionMetrics     `bson:"metrics" toml:"metrics"`
}

// UserInterfaceConf represents a user interface configuration object.
type UserInterfaceConf struct {
	IsAllowChangeSeverityToInfo bool `bson:"allow_change_severity_to_info"`
	// MaxMatchedItems need to warn user when number of items that match patterns is above this value
	MaxMatchedItems          int `bson:"max_matched_items"`
	CheckCountRequestTimeout int `bson:"check_count_request_timeout"`
}
