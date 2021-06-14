package config

import (
	"time"
)

// Some default values related to configuration
const (
	ConfigKeyName = "global_config"
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

func (s *SectionTimezone) GetLocation() (*time.Location, error) {
	location := time.UTC
	if s.Timezone != "" {
		var err error
		location, err = time.LoadLocation(s.Timezone)
		if err != nil {
			return nil, err
		}
	}

	return location, nil
}

type SectionRemediation struct {
	JobExecutorFetchTimeoutSeconds int64 `toml:"JobExecutorFetchTimeoutSeconds"`
}

// CanopsisConf represents a generic configuration object.
type CanopsisConf struct {
	ID          string             `bson:"_id,omitempty" toml:"omitempty"`
	Global      SectionGlobal      `bson:"global" toml:"global"`
	Alarm       SectionAlarm       `bson:"alarm" toml:"alarm"`
	Timezone    SectionTimezone    `bson:"timezone" toml:"timezone"`
	Remediation SectionRemediation `bson:"remediation" toml:"remediation"`
}
