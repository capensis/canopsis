package main

import (
	"flag"
	"time"
)

const (
	defaultPort     = 9180
	defaultInterval = time.Second * 10
)

func (f *Flags) ParseArgs() {
	flag.BoolVar(&f.Version, "version", false, "Show the version information")
	flag.Int64Var(&f.Port, "port", defaultPort, "Server port")
	flag.BoolVar(&f.Debug, "d", false, "debug")
	flag.DurationVar(&f.UpdateMetricsInterval, "updateMetricsInterval", defaultInterval, "Duration to wait between two run of update metrics processes")
	flag.Parse()
}

type Flags struct {
	UpdateMetricsInterval time.Duration
	Port                  int64
	Version               bool
	Debug                 bool
}
