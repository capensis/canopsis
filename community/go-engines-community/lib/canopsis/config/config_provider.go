package config

//go:generate mockgen -destination=../../../mocks/lib/canopsis/config/config.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config AlarmConfigProvider,TimezoneConfigProvider,RemediationConfigProvider,UserInterfaceConfigProvider,DataStorageConfigProvider

import (
	"fmt"
	"html/template"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

var weekdays = map[string]time.Weekday{}

func init() {
	for d := time.Sunday; d <= time.Saturday; d++ {
		weekdays[d.String()] = d
	}
}

type Updater interface {
	Update(CanopsisConf) error
}

type AlarmConfigProvider interface {
	Get() AlarmConfig
}

type TimezoneConfigProvider interface {
	Get() TimezoneConfig
}

type RemediationConfigProvider interface {
	Get() RemediationConfig
}

type DataStorageConfigProvider interface {
	Get() DataStorageConfig
}

type UserInterfaceConfigProvider interface {
	Get() UserInterfaceConf
}

type AlarmConfig struct {
	FlappingFreqLimit     int
	FlappingInterval      time.Duration
	StealthyInterval      time.Duration
	BaggotTime            time.Duration
	EnableLastEventDate   bool
	CancelAutosolveDelay  time.Duration
	DisplayNameScheme     *template.Template
	displayNameSchemeText string
	OutputLength          int
	LongOutputLength      int
	// DisableActionSnoozeDelayOnPbh ignores Pbh state to resolve snoozed with Action alarm while is True
	DisableActionSnoozeDelayOnPbh bool
}

type TimezoneConfig struct {
	Location *time.Location
}

type RemediationConfig struct {
	JobExecutorFetchTimeout time.Duration
}

type DataStorageConfig struct {
	TimeToExecute *ScheduledTime
}

type ScheduledTime struct {
	Weekday time.Weekday
	Hour    int
}

func (t ScheduledTime) String() string {
	return fmt.Sprintf("%v,%v", t.Weekday, t.Hour)
}

func NewAlarmConfigProvider(cfg CanopsisConf, logger zerolog.Logger) *BaseAlarmConfigProvider {
	conf := AlarmConfig{}

	if cfg.Alarm.BaggotTime == "" {
		conf.BaggotTime = AlarmBaggotTime
		logger.Error().
			Str("default", conf.BaggotTime.String()).
			Str("invalid", cfg.Alarm.BaggotTime).
			Msg("BaggotTime of alarm config section is not defined, default value is used instead")
	} else {
		var err error
		conf.BaggotTime, err = time.ParseDuration(cfg.Alarm.BaggotTime)
		if err == nil {
			logger.Info().
				Str("value", conf.BaggotTime.String()).
				Msg("BaggotTime of alarm config section is used")
		} else {
			conf.BaggotTime = AlarmBaggotTime
			logger.Err(err).
				Str("default", conf.BaggotTime.String()).
				Str("invalid", cfg.Alarm.BaggotTime).
				Msg("bad value BaggotTime of alarm config section, default value is used instead")
		}
	}

	if cfg.Alarm.CancelAutosolveDelay == "" {
		conf.CancelAutosolveDelay = AlarmCancelAutosolveDelay
		logger.Error().
			Str("default", conf.CancelAutosolveDelay.String()).
			Str("invalid", cfg.Alarm.CancelAutosolveDelay).
			Msg("CancelAutosolveDelay of alarm config section is not defined, default value is used instead")
	} else {
		var err error
		conf.CancelAutosolveDelay, err = time.ParseDuration(cfg.Alarm.CancelAutosolveDelay)
		if err == nil {
			logger.Info().
				Str("value", conf.CancelAutosolveDelay.String()).
				Msg("CancelAutosolveDelay of alarm config section is used")
		} else {
			conf.CancelAutosolveDelay = AlarmCancelAutosolveDelay
			logger.Err(err).
				Str("default", conf.CancelAutosolveDelay.String()).
				Str("invalid", cfg.Alarm.CancelAutosolveDelay).
				Msg("bad value CancelAutosolveDelay of alarm config section, default value is used instead")
		}
	}

	if cfg.Alarm.DisplayNameScheme == "" {
		var err error
		conf.displayNameSchemeText = AlarmDefaultNameScheme
		conf.DisplayNameScheme, err = CreateDisplayNameTpl(conf.displayNameSchemeText)
		if err != nil {
			panic(fmt.Errorf("invalid contant AlarmDefaultNameScheme: %w", err))
		}
		logger.Error().
			Str("default", conf.displayNameSchemeText).
			Str("invalid", cfg.Alarm.DisplayNameScheme).
			Msg("DisplayNameScheme of alarm config section is not defined, default value is used instead")
	} else {
		var err error
		conf.displayNameSchemeText = cfg.Alarm.DisplayNameScheme
		conf.DisplayNameScheme, err = CreateDisplayNameTpl(conf.displayNameSchemeText)
		if err == nil {
			logger.Info().
				Str("value", conf.displayNameSchemeText).
				Msg("DisplayNameScheme of alarm config section is used")
		} else {
			var err error
			conf.displayNameSchemeText = AlarmDefaultNameScheme
			conf.DisplayNameScheme, err = CreateDisplayNameTpl(conf.displayNameSchemeText)
			if err != nil {
				panic(fmt.Errorf("invalid contant AlarmDefaultNameScheme: %w", err))
			}
			logger.Err(err).
				Str("default", conf.displayNameSchemeText).
				Str("invalid", cfg.Alarm.DisplayNameScheme).
				Msg("bad value DisplayNameScheme of alarm config section, default value is used instead")
		}
	}

	if cfg.Alarm.OutputLength <= 0 {
		logger.Warn().Msg("OutputLength of alarm config section is not set or less than 1: the event's output won't be truncated")
	} else {
		conf.OutputLength = cfg.Alarm.OutputLength
		logger.Info().
			Int("value", conf.OutputLength).
			Msg("OutputLength of alarm config section is used")
	}

	if cfg.Alarm.LongOutputLength <= 0 {
		logger.Warn().Msg("LongOutputLength of alarm config section is not set or less than 1: the event's long_output won't be truncated")
	} else {
		conf.LongOutputLength = cfg.Alarm.LongOutputLength
		logger.Info().
			Int("value", conf.LongOutputLength).
			Msg("LongOutputLength of alarm config section is used")
	}

	if cfg.Alarm.FlappingFreqLimit < 0 {
		logger.Error().
			Int("default", conf.FlappingFreqLimit).
			Int("invalid", cfg.Alarm.FlappingFreqLimit).
			Msg("bad value FlappingFreqLimit of alarm config section, default value is used instead")
	} else {
		conf.FlappingFreqLimit = cfg.Alarm.FlappingFreqLimit
		logger.Info().
			Int("value", conf.FlappingFreqLimit).
			Msg("FlappingFreqLimit of alarm config section is used")
	}

	if cfg.Alarm.FlappingInterval < 0 {
		logger.Error().
			Str("default", conf.FlappingInterval.String()).
			Int("invalid", cfg.Alarm.FlappingInterval).
			Msg("bad value FlappingInterval of alarm config section, default value is used instead")
	} else {
		conf.FlappingInterval = time.Second * time.Duration(cfg.Alarm.FlappingInterval)
		logger.Info().
			Str("value", conf.FlappingInterval.String()).
			Msg("FlappingInterval of alarm config section is used")
	}

	if cfg.Alarm.StealthyInterval < 0 {
		logger.Error().
			Str("default", conf.StealthyInterval.String()).
			Int("invalid", cfg.Alarm.StealthyInterval).
			Msg("bad value StealthyInterval of alarm config section, default value is used instead")
	} else {
		conf.StealthyInterval = time.Second * time.Duration(cfg.Alarm.StealthyInterval)
		logger.Info().
			Str("value", conf.StealthyInterval.String()).
			Msg("StealthyInterval of alarm config section is used")
	}

	conf.EnableLastEventDate = cfg.Alarm.EnableLastEventDate
	logger.Info().
		Bool("value", conf.EnableLastEventDate).
		Msg("EnableLastEventDate of alarm config section is used")

	conf.DisableActionSnoozeDelayOnPbh = cfg.Alarm.DisableActionSnoozeDelayOnPbh
	logger.Info().
		Bool("value", conf.DisableActionSnoozeDelayOnPbh).
		Msg("DisableActionSnoozeDelayOnPbh of alarm config section is used")

	return &BaseAlarmConfigProvider{
		conf:   conf,
		logger: logger,
	}
}

type BaseAlarmConfigProvider struct {
	conf   AlarmConfig
	mx     sync.RWMutex
	logger zerolog.Logger
}

func (p *BaseAlarmConfigProvider) Update(cfg CanopsisConf) error {
	if cfg.Alarm.BaggotTime == "" {
		p.logger.Error().
			Str("invalid", cfg.Alarm.BaggotTime).
			Msg("BaggotTime of alarm config section is not defined, previous value is used")
	} else {
		duration, err := time.ParseDuration(cfg.Alarm.BaggotTime)
		if err == nil {
			if p.conf.BaggotTime != duration {
				p.logger.Info().
					Str("previous", p.conf.BaggotTime.String()).
					Str("new", duration.String()).
					Msg("BaggotTime of alarm config section is loaded")

				p.mx.Lock()
				p.conf.BaggotTime = duration
				p.mx.Unlock()
			}
		} else {
			p.logger.Err(err).
				Str("invalid", cfg.Alarm.BaggotTime).
				Msg("bad value BaggotTime of alarm config section, previous value is used instead")
		}
	}

	if cfg.Alarm.CancelAutosolveDelay == "" {
		p.logger.Error().
			Str("invalid", cfg.Alarm.CancelAutosolveDelay).
			Msg("CancelAutosolveDelay of alarm config section is not defined, previous value is used")
	} else {
		duration, err := time.ParseDuration(cfg.Alarm.CancelAutosolveDelay)
		if err == nil {
			if p.conf.CancelAutosolveDelay != duration {
				p.logger.Info().
					Str("previous", p.conf.CancelAutosolveDelay.String()).
					Str("new", duration.String()).
					Msg("CancelAutosolveDelay of alarm config section is loaded")

				p.mx.Lock()
				p.conf.CancelAutosolveDelay = duration
				p.mx.Unlock()
			}
		} else {
			p.logger.Err(err).
				Str("invalid", cfg.Alarm.CancelAutosolveDelay).
				Msg("bad value CancelAutosolveDelay of alarm config section, previous value is used instead")
		}
	}

	if cfg.Alarm.DisplayNameScheme == "" {
		p.logger.Error().
			Str("invalid", cfg.Alarm.DisplayNameScheme).
			Msg("DisplayNameScheme of alarm config section is not defined, previous value is used")
	} else if cfg.Alarm.DisplayNameScheme != p.conf.displayNameSchemeText {
		displayNameScheme, err := CreateDisplayNameTpl(cfg.Alarm.DisplayNameScheme)
		if err == nil {
			p.logger.Info().
				Str("previous", p.conf.displayNameSchemeText).
				Str("new", cfg.Alarm.DisplayNameScheme).
				Msg("DisplayNameScheme of alarm config section is loaded")

			p.mx.Lock()
			p.conf.DisplayNameScheme = displayNameScheme
			p.conf.displayNameSchemeText = cfg.Alarm.DisplayNameScheme
			p.mx.Unlock()
		} else {
			p.logger.Err(err).
				Str("invalid", cfg.Alarm.DisplayNameScheme).
				Msg("bad value DisplayNameScheme of alarm config section, previous value is used instead")
		}
	}

	if cfg.Alarm.OutputLength != p.conf.OutputLength {
		if cfg.Alarm.OutputLength <= 0 {
			p.logger.Warn().
				Int("previous", p.conf.OutputLength).
				Int("new", cfg.Alarm.OutputLength).
				Msg("OutputLength of alarm config section is loaded, value is not set or less than 1: the event's output and long_output won't be truncated")
		} else {
			p.logger.Info().
				Int("previous", p.conf.OutputLength).
				Int("new", cfg.Alarm.OutputLength).
				Msg("OutputLength of alarm config section is loaded")
		}

		p.mx.Lock()
		p.conf.OutputLength = cfg.Alarm.OutputLength
		p.mx.Unlock()
	}

	if cfg.Alarm.FlappingFreqLimit < 0 {
		p.logger.Error().
			Int("invalid", cfg.Alarm.FlappingFreqLimit).
			Msg("FlappingFreqLimit of alarm config section is not defined, previous value is used")
	} else if cfg.Alarm.FlappingFreqLimit != p.conf.FlappingFreqLimit {
		p.logger.Info().
			Int("previous", p.conf.FlappingFreqLimit).
			Int("new", cfg.Alarm.FlappingFreqLimit).
			Msg("FlappingFreqLimit of alarm config section is loaded")

		p.mx.Lock()
		p.conf.FlappingFreqLimit = cfg.Alarm.FlappingFreqLimit
		p.mx.Unlock()
	}

	if cfg.Alarm.FlappingInterval < 0 {
		p.logger.Error().
			Int("invalid", cfg.Alarm.FlappingInterval).
			Msg("FlappingInterval of alarm config section is not defined, previous value is used")
	} else {
		flappingInterval := time.Second * time.Duration(cfg.Alarm.FlappingInterval)
		if flappingInterval != p.conf.FlappingInterval {
			p.logger.Info().
				Str("previous", p.conf.FlappingInterval.String()).
				Str("new", flappingInterval.String()).
				Msg("FlappingInterval of alarm config section is loaded")

			p.mx.Lock()
			p.conf.FlappingInterval = flappingInterval
			p.mx.Unlock()
		}
	}

	if cfg.Alarm.StealthyInterval < 0 {
		p.logger.Error().
			Int("invalid", cfg.Alarm.StealthyInterval).
			Msg("StealthyInterval of alarm config section is not defined, previous value is used")
	} else {
		stealthyInterval := time.Second * time.Duration(cfg.Alarm.StealthyInterval)
		if stealthyInterval != p.conf.StealthyInterval {
			p.logger.Info().
				Str("previous", p.conf.StealthyInterval.String()).
				Str("new", stealthyInterval.String()).
				Msg("StealthyInterval of alarm config section is loaded")

			p.mx.Lock()
			p.conf.StealthyInterval = stealthyInterval
			p.mx.Unlock()
		}
	}

	if cfg.Alarm.EnableLastEventDate != p.conf.EnableLastEventDate {
		p.logger.Info().
			Bool("previous", p.conf.EnableLastEventDate).
			Bool("new", cfg.Alarm.EnableLastEventDate).
			Msg("EnableLastEventDate of alarm config section is loaded")

		p.mx.Lock()
		p.conf.EnableLastEventDate = cfg.Alarm.EnableLastEventDate
		p.mx.Unlock()
	}

	if cfg.Alarm.DisableActionSnoozeDelayOnPbh != p.conf.DisableActionSnoozeDelayOnPbh {
		p.logger.Info().
			Bool("previous", p.conf.DisableActionSnoozeDelayOnPbh).
			Bool("new", cfg.Alarm.DisableActionSnoozeDelayOnPbh).
			Msg("DisableActionSnoozeDelayOnPbh of alarm config section is loaded")

		p.mx.Lock()
		p.conf.DisableActionSnoozeDelayOnPbh = cfg.Alarm.DisableActionSnoozeDelayOnPbh
		p.mx.Unlock()
	}

	return nil
}

func (p *BaseAlarmConfigProvider) Get() AlarmConfig {
	p.mx.RLock()
	defer p.mx.RUnlock()

	return p.conf
}

func NewTimezoneConfigProvider(cfg CanopsisConf, logger zerolog.Logger) *BaseTimezoneConfigProvider {
	var location *time.Location
	defaultLocation := time.UTC
	if cfg.Timezone.Timezone == "" {
		location = defaultLocation
		logger.Error().
			Str("default", location.String()).
			Str("invalid", cfg.Timezone.Timezone).
			Msg("Timezone of timezone config section is not defined, default value is used instead")
	} else {
		var err error
		location, err = time.LoadLocation(cfg.Timezone.Timezone)
		if err == nil {
			logger.Info().
				Str("value", location.String()).
				Msg("Timezone of timezone config section is used")
		} else {
			location = defaultLocation
			logger.Err(err).
				Str("default", location.String()).
				Str("invalid", cfg.Timezone.Timezone).
				Msg("bad value Timezone of timezone config section, default value is used instead")
		}
	}

	return &BaseTimezoneConfigProvider{
		conf: TimezoneConfig{
			Location: location,
		},
		logger: logger,
	}
}

type BaseTimezoneConfigProvider struct {
	conf   TimezoneConfig
	mx     sync.RWMutex
	logger zerolog.Logger
}

func (p *BaseTimezoneConfigProvider) Update(cfg CanopsisConf) error {
	if cfg.Timezone.Timezone == "" {
		p.logger.Error().
			Str("invalid", cfg.Timezone.Timezone).
			Msg("Timezone of timezone config section is not defined, previous value is used")
	} else {
		location, err := time.LoadLocation(cfg.Timezone.Timezone)
		if err == nil {
			if p.conf.Location.String() != location.String() {
				p.logger.Info().
					Str("previous", p.conf.Location.String()).
					Str("new", location.String()).
					Msg("Timezone of timezone config section is loaded")

				p.mx.Lock()
				defer p.mx.Unlock()
				p.conf.Location = location
			}
		} else {
			p.logger.Err(err).
				Str("invalid", cfg.Timezone.Timezone).
				Msg("bad value Timezone of timezone config section, previous value is used instead")
		}
	}

	return nil
}

func (p *BaseTimezoneConfigProvider) Get() TimezoneConfig {
	p.mx.RLock()
	defer p.mx.RUnlock()

	return p.conf
}

func NewRemediationConfigProvider(cfg CanopsisConf, logger zerolog.Logger) *BaseRemediationConfigProvider {
	var jobExecutorFetchTimeout time.Duration
	if cfg.Remediation.JobExecutorFetchTimeoutSeconds <= 0 {
		jobExecutorFetchTimeout = RemediationJobExecutorFetchTimeout
		logger.Error().
			Str("default", jobExecutorFetchTimeout.String()).
			Int64("invalid", cfg.Remediation.JobExecutorFetchTimeoutSeconds).
			Msg("bad value JobExecutorFetchTimeoutSeconds duration of remediation config section, default value is used instead")
	} else {
		jobExecutorFetchTimeout = time.Second * time.Duration(cfg.Remediation.JobExecutorFetchTimeoutSeconds)
		logger.Info().
			Str("value", jobExecutorFetchTimeout.String()).
			Msg("JobExecutorFetchTimeoutSeconds duration of remediation config section is used")
	}

	return &BaseRemediationConfigProvider{
		conf: RemediationConfig{
			JobExecutorFetchTimeout: jobExecutorFetchTimeout,
		},
		logger: logger,
	}
}

type BaseRemediationConfigProvider struct {
	conf   RemediationConfig
	mx     sync.RWMutex
	logger zerolog.Logger
}

func (p *BaseRemediationConfigProvider) Update(cfg CanopsisConf) error {
	if cfg.Remediation.JobExecutorFetchTimeoutSeconds <= 0 {
		p.logger.Error().
			Int64("invalid", cfg.Remediation.JobExecutorFetchTimeoutSeconds).
			Msg("bad value JobExecutorFetchTimeoutSeconds duration of remediation config section, previous value is used")
	} else {
		jobExecutorFetchTimeout := time.Second * time.Duration(cfg.Remediation.JobExecutorFetchTimeoutSeconds)
		if jobExecutorFetchTimeout != p.conf.JobExecutorFetchTimeout {
			p.logger.Info().
				Str("previous", p.conf.JobExecutorFetchTimeout.String()).
				Str("new", jobExecutorFetchTimeout.String()).
				Msg("JobExecutorFetchTimeoutSeconds duration of remediation config section is loaded")

			p.mx.Lock()
			defer p.mx.Unlock()
			p.conf.JobExecutorFetchTimeout = jobExecutorFetchTimeout
		}
	}

	return nil
}

func (p *BaseRemediationConfigProvider) Get() RemediationConfig {
	p.mx.RLock()
	defer p.mx.RUnlock()

	return p.conf
}

type BaseUserInterfaceConfigProvider struct {
	conf   UserInterfaceConf
	mx     sync.RWMutex
	logger zerolog.Logger
}

const DefaultMaxMatchedItems = 10000
const DefaultCheckCountRequestTimeout = 30

func NewUserInterfaceConfigProvider(cfg UserInterfaceConf, logger zerolog.Logger) *BaseUserInterfaceConfigProvider {
	maxMatchedItems := 0
	if cfg.MaxMatchedItems <= 0 {
		maxMatchedItems = DefaultMaxMatchedItems
		logger.Error().
			Int("default", maxMatchedItems).
			Int("invalid", cfg.MaxMatchedItems).
			Msg("bad value MaxMatchedItems of user interface config, default value is used instead")
	} else {
		maxMatchedItems = cfg.MaxMatchedItems
		logger.Info().
			Int("value", maxMatchedItems).
			Msg("MaxMatchedItems of user interface config is used")
	}

	checkCountRequestTimeout := 0
	if cfg.CheckCountRequestTimeout <= 0 {
		checkCountRequestTimeout = DefaultCheckCountRequestTimeout
		logger.Error().
			Int("default", checkCountRequestTimeout).
			Int("invalid", cfg.CheckCountRequestTimeout).
			Msg("bad value CheckCountRequestTimeout of user interface config, default value is used instead")
	} else {
		checkCountRequestTimeout = cfg.CheckCountRequestTimeout
		logger.Info().
			Int("value", checkCountRequestTimeout).
			Msg("CheckCountRequestTimeout of user interface config is used")
	}

	return &BaseUserInterfaceConfigProvider{
		conf: UserInterfaceConf{
			IsAllowChangeSeverityToInfo: cfg.IsAllowChangeSeverityToInfo,
			MaxMatchedItems:             maxMatchedItems,
			CheckCountRequestTimeout:    checkCountRequestTimeout,
		},
		logger: logger,
	}
}

func (p *BaseUserInterfaceConfigProvider) Update(conf UserInterfaceConf) error {
	p.mx.Lock()
	defer p.mx.Unlock()

	if conf.MaxMatchedItems <= 0 {
		p.logger.Error().
			Int("invalid", conf.MaxMatchedItems).
			Msg("bad value MaxMatchedItems of user interface config, previous value is used")
	} else {
		if conf.MaxMatchedItems != p.conf.MaxMatchedItems {
			p.logger.Info().
				Int("previous", p.conf.MaxMatchedItems).
				Int("new", conf.MaxMatchedItems).
				Msg("MaxMatchedItems of user interface config is loaded")

			p.conf.MaxMatchedItems = conf.MaxMatchedItems
		}
	}

	if conf.CheckCountRequestTimeout <= 0 {
		p.logger.Error().
			Int("invalid", conf.CheckCountRequestTimeout).
			Msg("bad value CheckCountRequestTimeout of user interface config, previous value is used")
	} else {
		if conf.CheckCountRequestTimeout != p.conf.CheckCountRequestTimeout {
			p.logger.Info().
				Int("previous", p.conf.CheckCountRequestTimeout).
				Int("new", conf.CheckCountRequestTimeout).
				Msg("CheckCountRequestTimeout of user interface config is loaded")

			p.conf.CheckCountRequestTimeout = conf.CheckCountRequestTimeout
		}
	}

	if conf.IsAllowChangeSeverityToInfo != p.conf.IsAllowChangeSeverityToInfo {
		p.logger.Info().
			Bool("previous", p.conf.IsAllowChangeSeverityToInfo).
			Bool("new", conf.IsAllowChangeSeverityToInfo).
			Msg("IsAllowChangeSeverityToInfo of user interface config is loaded")

		p.conf.IsAllowChangeSeverityToInfo = conf.IsAllowChangeSeverityToInfo
	}

	return nil
}

func (p *BaseUserInterfaceConfigProvider) Get() UserInterfaceConf {
	p.mx.RLock()
	defer p.mx.RUnlock()

	return p.conf
}

func NewDataStorageConfigProvider(cfg CanopsisConf, logger zerolog.Logger) *BaseDataStorageConfigProvider {
	return &BaseDataStorageConfigProvider{
		conf: DataStorageConfig{
			TimeToExecute: parseScheduledTime(cfg.DataStorage.TimeToExecute,
				"TimeToExecute", "data_storage", logger, "data archive and delete are disabled"),
		},
		logger: logger,
	}
}

type BaseDataStorageConfigProvider struct {
	conf     DataStorageConfig
	mx       sync.RWMutex
	logger   zerolog.Logger
	weekdays map[string]time.Weekday
}

func (p *BaseDataStorageConfigProvider) Update(cfg CanopsisConf) error {
	t, ok := parseUpdatedScheduledTime(cfg.DataStorage.TimeToExecute, p.conf.TimeToExecute,
		"TimeToExecute", "data_storage", p.logger)
	if ok {
		p.mx.Lock()
		defer p.mx.Unlock()
		p.conf.TimeToExecute = t
	}

	return nil
}

func (p *BaseDataStorageConfigProvider) Get() DataStorageConfig {
	p.mx.RLock()
	defer p.mx.RUnlock()

	return p.conf
}

func parseScheduledTime(
	v string,
	name, sectionName string,
	logger zerolog.Logger,
	msg string,
) *ScheduledTime {
	if v == "" {
		logger.Info().
			Msgf("missing %s of %s config section, %s", name, sectionName, msg)
		return nil
	}

	t, ok := stringToScheduledTime(v)
	if !ok {
		logger.Error().
			Str("invalid", v).
			Msgf("bad value %s of %s config section, %s", name, sectionName, msg)
		return nil
	}

	logger.Info().
		Str("value", t.String()).
		Msgf("%s of %s config section is used", name, sectionName)

	return &t
}

func parseUpdatedScheduledTime(
	v string,
	oldVal *ScheduledTime,
	name, sectionName string,
	logger zerolog.Logger,
) (*ScheduledTime, bool) {
	if v == "" {
		if oldVal != nil {
			logger.Error().
				Str("invalid", v).
				Msgf("%s of %s config section is not defined, previous value is used", name, sectionName)
		}
		return nil, false
	}
	t, ok := stringToScheduledTime(v)
	if !ok {
		logger.Error().
			Str("invalid", v).
			Msgf("bad value %s of %s config section, previous value is used instead", name, sectionName)
		return nil, false
	}

	if oldVal != nil && oldVal.String() == t.String() {
		return nil, false
	}

	oldValStr := ""
	if oldVal != nil {
		oldValStr = oldVal.String()
	}
	logger.Info().
		Str("previous", oldValStr).
		Str("new", t.String()).
		Msgf("%s of %s config section is loaded", name, sectionName)

	return &t, true
}

func stringToScheduledTime(v string) (ScheduledTime, bool) {
	split := strings.Split(v, ",")
	t := ScheduledTime{}
	if len(split) == 2 {
		if d, ok := weekdays[split[0]]; ok {
			h, err := strconv.Atoi(split[1])
			if err == nil {
				t.Weekday = d
				t.Hour = h
				return t, true
			}
		}
	}

	return t, false
}
