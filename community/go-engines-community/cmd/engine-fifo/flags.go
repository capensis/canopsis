package main

import (
	"flag"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"time"
)

func (o *Options) ParseArgs() {
	flag.StringVar(&o.PublishToQueue, "publishQueue", canopsis.CheQueueName, "Publish event to this queue.")
	flag.StringVar(&o.ConsumeFromQueue, "consumeQueue", canopsis.FIFOQueueName, "Consume events from this queue.")
	flag.BoolVar(&o.ModeDebug, "d", false, "debug")
	flag.BoolVar(&o.PrintEventOnError, "printEventOnError", false, "Print event on processing error")
	flag.IntVar(&o.LockTtl, "lockTtl", 10, "Redis lock ttl time in seconds")
	flag.DurationVar(&o.EventsStatsFlushInterval, "eventsStatsFlushInterval", 60*time.Second, "Interval between saving statistics from redis to mongo")
	flag.Bool("enableMetaAlarmProcessing", true, "Enable meta-alarm processing - deprecated")

	flagVersion := flag.Bool("version", false, "version infos")

	flag.Parse()

	if *flagVersion {
		canopsis.PrintVersionExit()
	}
}
