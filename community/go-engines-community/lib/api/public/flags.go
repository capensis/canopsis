package api

import (
	"flag"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
)

const (
	defaultPort      = 8083
	defaultConfigDir = "/opt/canopsis/share/config"
)

func BindCmdFlags(o *Options) {
	log.BindCmdFlags(&o.LogOptions)
	flag.BoolVar(&o.Version, "version", false, "Show the version information")
	flag.IntVar(&o.Port, "port", defaultPort, "Server port")
	flag.StringVar(&o.ConfigDir, "c", defaultConfigDir, "Configuration files directory")
	flag.BoolVar(&o.EnableDocs, "docs", false, "Set to enable Swagger docs")
	flag.BoolVar(&o.LogBody, "logBody", false, "Set to enable logging response and request bodies")
	flag.BoolVar(&o.LogBodyOnError, "logBodyOnError", false, "Set to enable logging response and request bodies in case of error")
	flag.DurationVar(&o.PeriodicalWaitTime, "periodicalWaitTime", canopsis.PeriodicalWaitTime, "Interval to wait between two run of periodical process")
}

type Options struct {
	LogOptions log.Options

	Version    bool
	Port       int
	ConfigDir  string
	EnableDocs bool

	LogBody        bool
	LogBodyOnError bool

	PeriodicalWaitTime time.Duration
}
