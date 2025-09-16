package main

import (
	"flag"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics/prometheus"
)

const (
	defaultInterval = time.Second * 10
)

func (f *Flags) ParseArgs() {
	flag.BoolVar(&f.Version, "version", false, "Show the version information")
	flag.IntVar(&f.Port, "port", prometheus.DefaultExporterPort, "Prometheus exporter port")
	flag.BoolVar(&f.Debug, "d", false, "debug")
	flag.DurationVar(&f.UpdateMetricsInterval, "updateMetricsInterval", defaultInterval, "Duration to wait between two run of update metrics processes")
	flag.Parse()
}

type Flags struct {
	UpdateMetricsInterval time.Duration
	Port                  int
	Version               bool
	Debug                 bool
}
