package api

import (
	"flag"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
)

const (
	defaultPort      = 8082
	defaultConfigDir = "/opt/canopsis/share/config"
)

func (f *Flags) ParseArgs() {
	flag.BoolVar(&f.Version, "version", false, "Show the version information")
	flag.Int64Var(&f.Port, "port", defaultPort, "Server port")
	flag.StringVar(&f.ConfigDir, "c", defaultConfigDir, "Configuration files directory")
	flag.BoolVar(&f.Debug, "d", false, "debug")
	flag.BoolVar(&f.SecureSession, "secure", false, "Secure session")
	flag.BoolVar(&f.EnableDocs, "docs", false, "Set to enable Swagger docs")
	flag.DurationVar(&f.PeriodicalWaitTime, "periodicalWaitTime", canopsis.PeriodicalWaitTime, "Duration to wait between two run of periodical process")
	flag.DurationVar(&f.IntegrationPeriodicalWaitTime, "integrationPeriodicalWaitTime", 5*time.Second, "Duration to periodically check results of engines' tasks")
	flag.DurationVar(&f.EntityCategoryMetaPeriodicalWaitTime, "entityCategoryMetaPeriodicalWaitTime", 1*time.Minute, "Duration to wait between two run of periodical process to update entity category meta")
	flag.BoolVar(&f.EnableSameServiceNames, "enableSameServiceNames", false, "Enable same service names, services have unique names by default")
	flag.Parse()
}

type Flags struct {
	Version       bool
	Port          int64
	ConfigDir     string
	Debug         bool
	SecureSession bool
	EnableDocs    bool

	PeriodicalWaitTime                   time.Duration
	IntegrationPeriodicalWaitTime        time.Duration
	EntityCategoryMetaPeriodicalWaitTime time.Duration

	// EnableSameServiceNames affects entityservice Create/Update payload validation
	EnableSameServiceNames bool
}
