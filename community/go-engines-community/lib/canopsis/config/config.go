package config

import (
	"time"
)

// Some default values related to configuration
const (
	ConfigKeyName        = "global_config"
	UserInterfaceKeyName = "user_interface"
	VersionKeyName       = "canopsis_version"
)

// SectionAlarm ...
type SectionAlarm struct {
	FlappingFreqLimit    int    `toml:"FlappingFreqLimit"`
	FlappingInterval     int    `toml:"FlappingInterval"`
	StealthyInterval     int    `toml:"StealthyInterval"`
	BaggotTime           string `toml:"BaggotTime"`
	EnableLastEventDate  bool   `toml:"EnableLastEventDate"`
	CancelAutosolveDelay string `toml:"CancelAutosolveDelay"`
	DisplayNameScheme    string `toml:"DisplayNameScheme"`
	OutputLength         int    `toml:"OutputLength"`
	// DisableActionSnoozeDelayOnPbh ignores Pbh state to resolve snoozed with Action alarm while is True
	DisableActionSnoozeDelayOnPbh bool `toml:"DisableActionSnoozeDelayOnPbh"`
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

type SectionRemediation struct {
	JobExecutorFetchTimeoutSeconds int64 `toml:"JobExecutorFetchTimeoutSeconds"`
}

type SectionImportCtx struct {
	ThdWarnMinPerImport string `toml:"ThdWarnMinPerImport"`
	ThdCritMinPerImport string `toml:"ThdCritMinPerImport"`
	FilePattern         string `toml:"FilePattern"`
}

type SectionFile struct {
	Remediation string `toml:"Remediation"`
	Junit       string `toml:"Junit"`
	JunitApi    string `toml:"JunitApi"`
}

type SectionDataStorage struct {
	TimeToExecute string `toml:"TimeToExecute"`
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

// CanopsisConf represents a generic configuration object.
type CanopsisConf struct {
	ID          string             `bson:"_id,omitempty" toml:"omitempty"`
	Global      SectionGlobal      `bson:"global" toml:"global"`
	Alarm       SectionAlarm       `bson:"alarm" toml:"alarm"`
	Timezone    SectionTimezone    `bson:"timezone" toml:"timezone"`
	Remediation SectionRemediation `bson:"remediation" toml:"remediation"`
	ImportCtx   SectionImportCtx   `bson:"import_ctx" toml:"import_ctx"`
	File        SectionFile        `bson:"file" toml:"file"`
	DataStorage SectionDataStorage `bson:"data_storage" toml:"data_storage"`
	Logger      SectionLogger      `bson:"logger" toml:"logger"`
}

// UserInterfaceConf represents a user interface configuration object.
type UserInterfaceConf struct {
	IsAllowChangeSeverityToInfo bool `bson:"allow_change_severity_to_info"`
	// MaxMatchedItems need to warn user when number of items that match patterns is above this value
	MaxMatchedItems          int `bson:"max_matched_items"`
	CheckCountRequestTimeout int `bson:"check_count_request_timeout"`
}
