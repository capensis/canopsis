package api

import (
	"flag"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
)

const (
	defaultPort      = 8082
	defaultConfigDir = "/opt/canopsis/share/config"
)

func (f *Flags) ParseArgs() {
	log.BindCmdFlags(&f.Options)
	flag.BoolVar(&f.Version, "version", false, "Show the version information")
	flag.Int64Var(&f.Port, "port", defaultPort, "Server port")
	flag.StringVar(&f.ConfigDir, "c", defaultConfigDir, "Configuration files directory")
	flag.BoolVar(&f.SecureSession, "secure", false, "Secure session")
	flag.BoolVar(&f.EnableDocs, "docs", false, "Set to enable Swagger docs")
	flag.DurationVar(&f.PeriodicalWaitTime, "periodicalWaitTime", canopsis.PeriodicalWaitTime, "Interval to wait between two run of periodical process")
	flag.DurationVar(&f.IntegrationPeriodicalWaitTime, "integrationPeriodicalWaitTime", 5*time.Second, "Interval to periodically check results of engines' tasks")
	flag.DurationVar(&f.EntityCategoryMetaPeriodicalWaitTime, "entityCategoryMetaPeriodicalWaitTime", time.Minute, "Interval to periodically update entity category meta")
	flag.DurationVar(&f.InstructionRateNotificationPeriodicalWaitTime, "instructionRateNotificationPeriodicalWaitTime", time.Hour, "Interval to periodically check instructions and create rate notifications")
	flag.DurationVar(&f.BroadcastMessagePeriodicalWaitTime, "broadcastMessagePeriodicalWaitTime", time.Minute, "Interval to periodically check broadcast messages")
	flag.DurationVar(&f.StateSettingRecomputeDelay, "stateSettingRecomputeDelay", time.Second, "Minimum duration to wait before send recompute event for services and components")
	flag.BoolVar(&f.EnableSameServiceNames, "enableSameServiceNames", false, "Enable same service names, services have unique names by default")
	flag.DurationVar(&f.ExternalDataAPITimeout, "externalDataAPITimeout", 30*time.Second, "External API HTTP Request Timeout.")
	flag.BoolVar(&f.LogBody, "logBody", false, "Set to enable logging response and request bodies")
	flag.BoolVar(&f.LogBodyOnError, "logBodyOnError", false, "Set to enable logging response and request bodies in case of error")
	flag.Parse()
}

type Flags struct {
	log.Options

	Version       bool
	Port          int64
	ConfigDir     string
	SecureSession bool
	EnableDocs    bool

	PeriodicalWaitTime                            time.Duration
	IntegrationPeriodicalWaitTime                 time.Duration
	EntityCategoryMetaPeriodicalWaitTime          time.Duration
	InstructionRateNotificationPeriodicalWaitTime time.Duration
	BroadcastMessagePeriodicalWaitTime            time.Duration

	StateSettingRecomputeDelay time.Duration

	// EnableSameServiceNames affects entityservice Create/Update payload validation
	EnableSameServiceNames bool

	ExternalDataAPITimeout time.Duration

	LogBody        bool
	LogBodyOnError bool
}
