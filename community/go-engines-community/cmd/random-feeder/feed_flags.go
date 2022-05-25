package main

import (
	"flag"
	cps "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
)

func (f *Flags) ParseArgs() error {
	flagEventsPerSec := flag.Int("opersec", int(500), "events per second")
	flagNewResourcesPerSec := flag.Int("npersec", int(5), "new resources per second")
	flagDumpFile := flag.String("dumpfile", "", "file with dump")
	flagFeederTime := flag.Int("time", int(60), "feeder duration")
	flagAlarms := flag.Int("alarms", 20, "percent of alarms")
	flagExchange := flag.String("exchange", cps.CheExchangeName, "exchange name to publish events to")

	flag.Parse()

	f.EventsPerSec = *flagEventsPerSec
	f.NewResourcesPerSec = *flagNewResourcesPerSec
	f.DumpFile = *flagDumpFile
	f.Alarms = *flagAlarms
	f.FeederTime = *flagFeederTime
	f.ExchangeName = *flagExchange

	return nil
}

type Flags struct {
	EventsPerSec       int
	NewResourcesPerSec int
	Alarms             int
	DumpFile           string
	FeederTime         int
	ExchangeName       string
}
