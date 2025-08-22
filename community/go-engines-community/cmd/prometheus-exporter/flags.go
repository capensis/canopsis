package main

import (
	"flag"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
)

const (
	defaultPort     = 9180
	defaultInterval = time.Second * 10
)

func (f *Flags) ParseArgs() {
	log.BindCmdFlags(&f.Options)
	flag.BoolVar(&f.Version, "version", false, "Show the version information")
	flag.IntVar(&f.Port, "port", defaultPort, "Server port")
	flag.BoolVar(&f.Debug, "d", false, "debug")
	flag.DurationVar(&f.UpdateMetricsInterval, "updateMetricsInterval", defaultInterval, "Duration to wait between two run of update metrics processes")
	flag.Parse()
}

type Flags struct {
	log.Options
	UpdateMetricsInterval time.Duration
	Port                  int
	Version               bool
	Debug                 bool
}
