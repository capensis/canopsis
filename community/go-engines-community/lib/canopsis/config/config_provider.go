package config

//go:generate mockgen -destination=../../../mocks/lib/canopsis/config/config.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config AlarmConfigProvider,TimezoneConfigProvider,RemediationConfigProvider,UserInterfaceConfigProvider,DataStorageConfigProvider

import (
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"html/template"
	"reflect"
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
	// DisableActionSnoozeDelayOnPbh ignores Pbh state to resolve snoozed with Action alarm while is True
	DisableActionSnoozeDelayOnPbh bool
	TimeToKeepResolvedAlarms      time.Duration
}

type TimezoneConfig struct {
	Location *time.Location
}

type RemediationConfig struct {
	HttpTimeout                    time.Duration
	LaunchJobRetriesAmount         int
	LaunchJobRetriesInterval       time.Duration
	WaitJobCompleteRetriesAmount   int
	WaitJobCompleteRetriesInterval time.Duration
	PauseManualInstructionInterval time.Duration
	ExternalAPI                    map[string]ExternalApiConfig
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
	sectionName := "alarm"
	conf := AlarmConfig{
		FlappingFreqLimit:             parseInt(cfg.Alarm.FlappingFreqLimit, 0, "FlappingFreqLimit", sectionName, logger),
		FlappingInterval:              parseTimeDurationBySeconds(cfg.Alarm.FlappingInterval, 0, "FlappingInterval", sectionName, logger),
		StealthyInterval:              parseTimeDurationBySeconds(cfg.Alarm.StealthyInterval, 0, "StealthyInterval", sectionName, logger),
		BaggotTime:                    parseTimeDurationByStr(cfg.Alarm.BaggotTime, AlarmBaggotTime, "BaggotTime", sectionName, logger),
		EnableLastEventDate:           parseBool(cfg.Alarm.EnableLastEventDate, "EnableLastEventDate", sectionName, logger),
		CancelAutosolveDelay:          parseTimeDurationByStr(cfg.Alarm.CancelAutosolveDelay, AlarmCancelAutosolveDelay, "CancelAutosolveDelay", sectionName, logger),
		DisableActionSnoozeDelayOnPbh: parseBool(cfg.Alarm.DisableActionSnoozeDelayOnPbh, "DisableActionSnoozeDelayOnPbh", sectionName, logger),
		TimeToKeepResolvedAlarms:      parseTimeDurationByStr(cfg.Alarm.TimeToKeepResolvedAlarms, 0, "TimeToKeepResolvedAlarms", sectionName, logger),
	}
	conf.DisplayNameScheme, conf.displayNameSchemeText = parseTemplate(cfg.Alarm.DisplayNameScheme, AlarmDefaultNameScheme, "DisplayNameScheme", sectionName, logger)

	if cfg.Alarm.OutputLength <= 0 {
		logger.Warn().Msgf("OutputLength of %s config section is not set or less than 1: the event's output and long_output won't be truncated", sectionName)
	} else {
		conf.OutputLength = cfg.Alarm.OutputLength
		logger.Info().
			Int("value", conf.OutputLength).
			Msgf("OutputLength of %s config section is used", sectionName)
	}

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
	p.mx.Lock()
	defer p.mx.Unlock()

	sectionName := "alarm"
	d, ok := parseUpdatedTimeDurationByStr(cfg.Alarm.BaggotTime, p.conf.BaggotTime, "BaggotTime", sectionName, p.logger)
	if ok {
		p.conf.BaggotTime = d
	}

	d, ok = parseUpdatedTimeDurationByStr(cfg.Alarm.CancelAutosolveDelay, p.conf.CancelAutosolveDelay, "CancelAutosolveDelay", sectionName, p.logger)
	if ok {
		p.conf.CancelAutosolveDelay = d
	}

	t, s, ok := parseUpdatedTemplate(cfg.Alarm.DisplayNameScheme, p.conf.displayNameSchemeText, "DisplayNameScheme", sectionName, p.logger)
	if ok {
		p.conf.DisplayNameScheme = t
		p.conf.displayNameSchemeText = s
	}

	if cfg.Alarm.OutputLength != p.conf.OutputLength {
		if cfg.Alarm.OutputLength <= 0 {
			p.conf.OutputLength = 0
			p.logger.Warn().
				Int("previous", p.conf.OutputLength).
				Int("new", cfg.Alarm.OutputLength).
				Msg("OutputLength of alarm config section is loaded, value is not set or less than 1: the event's output and long_output won't be truncated")
		} else {
			p.conf.OutputLength = cfg.Alarm.OutputLength
			p.logger.Info().
				Int("previous", p.conf.OutputLength).
				Int("new", cfg.Alarm.OutputLength).
				Msg("OutputLength of alarm config section is loaded")
		}
	}

	i, ok := parseUpdatedInt(cfg.Alarm.FlappingFreqLimit, p.conf.FlappingFreqLimit, "FlappingFreqLimit", sectionName, p.logger)
	if ok {
		p.conf.FlappingFreqLimit = i
	}

	d, ok = parseUpdatedTimeDurationBySeconds(cfg.Alarm.FlappingInterval, p.conf.FlappingInterval, "FlappingInterval", sectionName, p.logger)
	if ok {
		p.conf.FlappingInterval = d
	}

	d, ok = parseUpdatedTimeDurationBySeconds(cfg.Alarm.StealthyInterval, p.conf.StealthyInterval, "StealthyInterval", sectionName, p.logger)
	if ok {
		p.conf.StealthyInterval = d
	}

	d, ok = parseUpdatedTimeDurationByStr(cfg.Alarm.TimeToKeepResolvedAlarms, p.conf.TimeToKeepResolvedAlarms, "TimeToKeepResolvedAlarms", sectionName, p.logger)
	if ok {
		p.mx.Lock()
		p.conf.TimeToKeepResolvedAlarms = d
		p.mx.Unlock()
	}

	b, ok := parseUpdatedBool(cfg.Alarm.EnableLastEventDate, p.conf.EnableLastEventDate, "EnableLastEventDate", sectionName, p.logger)
	if ok {
		p.conf.EnableLastEventDate = b
	}

	b, ok = parseUpdatedBool(cfg.Alarm.DisableActionSnoozeDelayOnPbh, p.conf.DisableActionSnoozeDelayOnPbh, "DisableActionSnoozeDelayOnPbh", sectionName, p.logger)
	if ok {
		p.conf.DisableActionSnoozeDelayOnPbh = b
	}

	return nil
}

func (p *BaseAlarmConfigProvider) Get() AlarmConfig {
	p.mx.RLock()
	defer p.mx.RUnlock()

	return p.conf
}

func NewTimezoneConfigProvider(cfg CanopsisConf, logger zerolog.Logger) *BaseTimezoneConfigProvider {
	return &BaseTimezoneConfigProvider{
		conf: TimezoneConfig{
			Location: parseLocation(cfg.Timezone.Timezone, time.UTC, "Timezone", "timezone", logger),
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
	p.mx.Lock()
	defer p.mx.Unlock()

	l, ok := parseUpdatedLocation(cfg.Timezone.Timezone, p.conf.Location, "Timezone", "timezone", p.logger)
	if ok {
		p.conf.Location = l
	}

	return nil
}

func (p *BaseTimezoneConfigProvider) Get() TimezoneConfig {
	p.mx.RLock()
	defer p.mx.RUnlock()

	return p.conf
}

func NewRemediationConfigProvider(cfg RemediationConf, logger zerolog.Logger) *BaseRemediationConfigProvider {
	sectionName := "Remediation"

	apiKeys := make([]string, len(cfg.ExternalAPI))
	i := 0
	for key := range cfg.ExternalAPI {
		apiKeys[i] = key
		i++
	}
	logger.Info().
		Msgf("%+v is loaded %s of %s config section", apiKeys, "external_api", sectionName)

	return &BaseRemediationConfigProvider{
		conf: RemediationConfig{
			HttpTimeout:                    parseTimeDurationByStr(cfg.HttpTimeout, RemediationHttpTimeout, "http_timeout", sectionName, logger),
			LaunchJobRetriesAmount:         parseInt(cfg.LaunchJobRetriesAmount, RemediationLaunchJobRetriesAmount, "launch_job_retries_amount", sectionName, logger),
			LaunchJobRetriesInterval:       parseTimeDurationByStr(cfg.LaunchJobRetriesInterval, RemediationLaunchJobRetriesInterval, "launch_job_retries_interval", sectionName, logger),
			WaitJobCompleteRetriesAmount:   parseInt(cfg.WaitJobCompleteRetriesAmount, RemediationWaitJobCompleteRetriesAmount, "wait_job_complete_retries_amount", sectionName, logger),
			WaitJobCompleteRetriesInterval: parseTimeDurationByStr(cfg.WaitJobCompleteRetriesInterval, RemediationWaitJobCompleteRetriesInterval, "wait_job_complete_retries_interval", sectionName, logger),
			PauseManualInstructionInterval: parseTimeDurationByStr(cfg.PauseManualInstructionInterval, RemediationPauseManualInstructionInterval, "pause_manual_instruction_interval", sectionName, logger),
			ExternalAPI:                    cfg.ExternalAPI,
		},
		logger: logger,
	}
}

type BaseRemediationConfigProvider struct {
	conf   RemediationConfig
	mx     sync.RWMutex
	logger zerolog.Logger
}

func (p *BaseRemediationConfigProvider) Update(cfg RemediationConf) error {
	p.mx.Lock()
	defer p.mx.Unlock()

	sectionName := "remediation"
	d, ok := parseUpdatedTimeDurationByStr(cfg.HttpTimeout, p.conf.HttpTimeout, "http_timeout", sectionName, p.logger)
	if ok {
		p.conf.HttpTimeout = d
	}
	i, ok := parseUpdatedInt(cfg.LaunchJobRetriesAmount, p.conf.LaunchJobRetriesAmount, "launch_job_retries_amount", sectionName, p.logger)
	if ok {
		p.conf.LaunchJobRetriesAmount = i
	}
	d, ok = parseUpdatedTimeDurationByStr(cfg.LaunchJobRetriesInterval, p.conf.LaunchJobRetriesInterval, "launch_job_retries_interval", sectionName, p.logger)
	if ok {
		p.conf.LaunchJobRetriesInterval = d
	}
	i, ok = parseUpdatedInt(cfg.WaitJobCompleteRetriesAmount, p.conf.WaitJobCompleteRetriesAmount, "wait_job_complete_retries_amount", sectionName, p.logger)
	if ok {
		p.conf.WaitJobCompleteRetriesAmount = i
	}
	d, ok = parseUpdatedTimeDurationByStr(cfg.WaitJobCompleteRetriesInterval, p.conf.WaitJobCompleteRetriesInterval, "wait_job_complete_retries_interval", sectionName, p.logger)
	if ok {
		p.conf.WaitJobCompleteRetriesInterval = d
	}
	d, ok = parseUpdatedTimeDurationByStr(cfg.PauseManualInstructionInterval, p.conf.PauseManualInstructionInterval, "pause_manual_instruction_interval", sectionName, p.logger)
	if ok {
		p.conf.PauseManualInstructionInterval = d
	}

	if !reflect.DeepEqual(cfg.ExternalAPI, p.conf.ExternalAPI) {
		apiKeys := make([]string, len(cfg.ExternalAPI))
		i := 0
		for key := range cfg.ExternalAPI {
			apiKeys[i] = key
			i++
		}
		p.logger.Info().
			Msgf("%+v is updated %s of %s config section", apiKeys, "external_api", sectionName)

		p.conf.ExternalAPI = cfg.ExternalAPI
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
	conf   DataStorageConfig
	mx     sync.RWMutex
	logger zerolog.Logger
}

func (p *BaseDataStorageConfigProvider) Update(cfg CanopsisConf) error {
	p.mx.Lock()
	defer p.mx.Unlock()

	t, ok := parseUpdatedScheduledTime(cfg.DataStorage.TimeToExecute, p.conf.TimeToExecute,
		"TimeToExecute", "data_storage", p.logger)
	if ok {
		p.conf.TimeToExecute = t
	}

	return nil
}

func (p *BaseDataStorageConfigProvider) Get() DataStorageConfig {
	p.mx.RLock()
	defer p.mx.RUnlock()

	return p.conf
}

type BaseHealthCheckConfigProvider struct {
	conf   HealthCheckConf
	mx     sync.RWMutex
	logger zerolog.Logger
}

func NewBaseHealthCheckConfigProvider(cfg HealthCheckConf, logger zerolog.Logger) *BaseHealthCheckConfigProvider {
	return &BaseHealthCheckConfigProvider{
		conf:   HealthCheckConf{EngineOrder: parseEngineOrder(cfg.EngineOrder, defaultConfig, logger)},
		logger: logger,
	}
}

func (p *BaseHealthCheckConfigProvider) Update(cfg HealthCheckConf) error {
	p.mx.Lock()
	defer p.mx.Unlock()

	p.conf.EngineOrder = parseEngineOrder(cfg.EngineOrder, p.conf.EngineOrder, p.logger)

	return nil
}

func (p *BaseHealthCheckConfigProvider) Get() HealthCheckConf {
	p.mx.RLock()
	defer p.mx.RUnlock()

	return p.conf
}

func parseEngineOrder(value []EngineOrder, oldValue []EngineOrder, logger zerolog.Logger) []EngineOrder {
	var possibleEngines = []string{
		canopsis.FIFOEngineName,
		canopsis.CheEngineName,
		canopsis.PBehaviorEngineName,
		canopsis.AxeEngineName,
		canopsis.CorrelationEngineName,
		canopsis.RemediationEngineName,
		canopsis.ServiceEngineName,
		canopsis.DynamicInfosEngineName,
		canopsis.ActionEngineName,
		canopsis.WebhookEngineName,
	}

	for idx, pair := range value {
		fromValid := false
		toValid := false

		for _, engineName := range possibleEngines {
			if pair.From == engineName {
				fromValid = true
			}

			if pair.To == engineName {
				toValid = true
			}

			if fromValid && toValid {
				break
			}
		}

		if !fromValid || !toValid {
			logger.Error().
				Int("index", idx).
				Str("from", pair.From).
				Str("to", pair.To).
				Msgf("from and to values should be one of %v", possibleEngines)

			return oldValue
		}
	}

	return value
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

func parseTimeDurationByStr(
	v string,
	defaultVal time.Duration,
	name, sectionName string,
	logger zerolog.Logger,
) time.Duration {
	if v == "" {
		logger.Error().
			Str("default", defaultVal.String()).
			Str("invalid", v).
			Msgf("%s of %s config section is not defined, default value is used instead", name, sectionName)

		return defaultVal
	}

	d, err := time.ParseDuration(v)
	if err != nil {
		logger.Err(err).
			Str("default", defaultVal.String()).
			Str("invalid", v).
			Msgf("bad value %s of %s config section, default value is used instead", name, sectionName)

		return defaultVal
	}

	logger.Info().
		Str("value", d.String()).
		Msgf("%s of %s config section is used", name, sectionName)

	return d
}

func parseUpdatedTimeDurationByStr(
	v string, oldValal time.Duration,
	name, sectionName string,
	logger zerolog.Logger,
) (time.Duration, bool) {
	if v == "" {
		logger.Error().
			Str("invalid", v).
			Msgf("%s of %s config section is not defined, previous value is used", name, sectionName)
		return 0, false
	}

	d, err := time.ParseDuration(v)
	if err != nil {
		logger.Err(err).
			Str("invalid", v).
			Msgf("bad value %s of %s config section, previous value is used instead", name, sectionName)
		return 0, false
	}

	if d == oldValal {
		return 0, false
	}

	logger.Info().
		Str("previous", oldValal.String()).
		Str("new", d.String()).
		Msgf("%s of %s config section is loaded", name, sectionName)

	return d, true
}

func parseTimeDurationBySeconds(
	v int,
	defaultVal time.Duration,
	name, sectionName string,
	logger zerolog.Logger,
) time.Duration {
	if v < 0 {
		logger.Error().
			Str("default", defaultVal.String()).
			Int("invalid", v).
			Msgf("bad value %s of %s config section, default value is used instead", name, sectionName)

		return defaultVal
	}

	d := time.Duration(v) * time.Second
	logger.Info().
		Str("value", d.String()).
		Msgf("%s of %s config section is used", name, sectionName)

	return d
}

func parseUpdatedTimeDurationBySeconds(
	v int,
	oldValal time.Duration,
	name, sectionName string,
	logger zerolog.Logger,
) (time.Duration, bool) {
	if v < 0 {
		logger.Error().
			Int("invalid", v).
			Msgf("bad value %s of %s config section, previous value is used instead", name, sectionName)
		return 0, false
	}

	d := time.Duration(v) * time.Second
	if d == oldValal {
		return 0, false
	}

	logger.Info().
		Str("previous", oldValal.String()).
		Str("new", d.String()).
		Msgf("%s of %s config section is loaded", name, sectionName)

	return d, true
}

func parseInt(
	v, defaultVal int,
	name, sectionName string,
	logger zerolog.Logger,
) int {
	if v < 0 {
		logger.Error().
			Int("default", defaultVal).
			Int("invalid", v).
			Msgf("bad value %s of %s config section, default value is used instead", name, sectionName)
		return defaultVal
	}

	logger.Info().
		Int("value", v).
		Msgf("%s of %s config section is used", name, sectionName)

	return v
}

func parseUpdatedInt(
	v, oldValal int,
	name, sectionName string,
	logger zerolog.Logger,
	invalidMsg ...string,
) (int, bool) {
	if v < 0 {
		msg := "bad value %s of %s config section, previous value is used instead"
		if len(invalidMsg) == 1 {
			msg = invalidMsg[0]
		}

		logger.Error().
			Int("invalid", v).
			Msgf(msg, name, sectionName)
		return 0, false
	}

	if v == oldValal {
		return 0, false
	}

	logger.Info().
		Int("previous", oldValal).
		Int("new", v).
		Msgf("%s of %s config section is loaded", name, sectionName)

	return v, true
}

func parseTemplate(
	v, defaultVal string,
	name, sectionName string,
	logger zerolog.Logger,
) (*template.Template, string) {
	if v == "" {
		tpl, err := CreateDisplayNameTpl(defaultVal)
		if err != nil {
			panic(fmt.Errorf("invalid contant %s: %w", name, err))
		}
		logger.Error().
			Str("default", defaultVal).
			Str("invalid", v).
			Msgf("%s of %s config section is not defined, default value is used instead", name, sectionName)

		return tpl, defaultVal
	}

	tpl, err := CreateDisplayNameTpl(v)
	if err != nil {
		tpl, parseErr := CreateDisplayNameTpl(defaultVal)
		if parseErr != nil {
			panic(fmt.Errorf("invalid contant %s: %w", name, parseErr))
		}

		logger.Err(err).
			Str("default", defaultVal).
			Str("invalid", v).
			Msgf("bad value %s of %s config section, default value is used instead", name, sectionName)

		return tpl, defaultVal
	}

	logger.Info().
		Str("value", v).
		Msgf("%s of %s config section is used", name, sectionName)

	return tpl, v
}

func parseUpdatedTemplate(
	v, oldVal string,
	name, sectionName string,
	logger zerolog.Logger,
) (*template.Template, string, bool) {
	if v == "" {
		logger.Error().
			Str("invalid", v).
			Msgf("%s of %s config section is not defined, previous value is used", name, sectionName)
		return nil, "", false
	}

	if v == oldVal {
		return nil, "", false
	}

	tpl, err := CreateDisplayNameTpl(v)
	if err != nil {
		logger.Err(err).
			Str("invalid", v).
			Msgf("bad value %s of %s config section, previous value is used instead", name, sectionName)
		return nil, "", false
	}

	logger.Info().
		Str("previous", oldVal).
		Str("new", v).
		Msgf("%s of %s config section is loaded", name, sectionName)

	return tpl, v, true
}

func parseBool(
	v bool,
	name, sectionName string,
	logger zerolog.Logger,
) bool {
	logger.Info().
		Bool("value", v).
		Msgf("%s of %s config section is used", name, sectionName)

	return v
}

func parseUpdatedBool(
	v, oldVal bool,
	name, sectionName string,
	logger zerolog.Logger,
) (bool, bool) {
	if v == oldVal {
		return false, false
	}
	logger.Info().
		Bool("previous", oldVal).
		Bool("new", v).
		Msgf("%s of %s config section is loaded", name, sectionName)

	return v, true
}

func parseLocation(
	v string,
	defaultVal *time.Location,
	name, sectionName string,
	logger zerolog.Logger,
) *time.Location {
	if v == "" {
		logger.Error().
			Str("default", defaultVal.String()).
			Str("invalid", v).
			Msgf("%s of %s config section is not defined, default value is used instead", name, sectionName)
		return defaultVal
	}

	location, err := time.LoadLocation(v)
	if err != nil {
		logger.Err(err).
			Str("default", defaultVal.String()).
			Str("invalid", v).
			Msgf("bad value %s of %s config section, default value is used instead", name, sectionName)
		return defaultVal
	}

	logger.Info().
		Str("value", location.String()).
		Msgf("%s of %s config section is used", name, sectionName)

	return location
}

func parseUpdatedLocation(
	v string,
	oldVal *time.Location,
	name, sectionName string,
	logger zerolog.Logger,
) (*time.Location, bool) {
	if v == "" {
		logger.Error().
			Str("invalid", v).
			Msgf("%s of %s config section is not defined, previous value is used", name, sectionName)
		return nil, false
	}
	location, err := time.LoadLocation(v)
	if err != nil {
		logger.Err(err).
			Str("invalid", v).
			Msgf("bad value %s of %s config section, previous value is used instead", name, sectionName)
		return nil, false
	}

	if oldVal.String() == location.String() {
		return nil, false
	}

	logger.Info().
		Str("previous", oldVal.String()).
		Str("new", location.String()).
		Msgf("%s of %s config section is loaded", name, sectionName)

	return location, true
}
