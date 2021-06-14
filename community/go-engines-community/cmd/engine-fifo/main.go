package main

import (
	"context"
	"flag"
	"os"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/debug"
	"git.canopsis.net/canopsis/go-engines/lib/log"
)

func main() {
	opts := Options{}

	flag.StringVar(&opts.PublishToQueue, "publishQueue", canopsis.CheQueueName, "Publish event to this queue.")
	flag.StringVar(&opts.ConsumeFromQueue, "consumeQueue", canopsis.FIFOQueueName, "Consume events from this queue.")
	flag.BoolVar(&opts.ModeDebug, "d", false, "debug")
	flag.BoolVar(&opts.PrintEventOnError, "printEventOnError", false, "Print event on processing error")
	flag.IntVar(&opts.LockTtl, "lockTtl", 10, "Redis lock ttl time in seconds")
	flag.BoolVar(&opts.EnableMetaAlarmProcessing, "enableMetaAlarmProcessing", true, "Enable meta-alarm processing")

	flagVersion := flag.Bool("version", false, "version infos")

	flag.Parse()

	if *flagVersion {
		canopsis.PrintVersionExit()
	}

	logger := log.NewLogger(opts.ModeDebug)

	trace := debug.Start(logger)

	ctx := context.Background()

	depMaker := DependencyMaker{}
	references := depMaker.GetDefaultReferences(ctx, opts, logger)
	engine := NewEngineFIFO(opts, references)

	exitStatus, err := canopsis.StartEngine(ctx, engine, nil)
	if err != nil {
		logger.Error().Err(err).Msg("")
	}

	trace.Stop()
	os.Exit(exitStatus)
}
