package main

import (
	"context"
	"flag"
	"os"
	"sync"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/debug"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
)

func main() {
	opts := Options{}

	flag.StringVar(&opts.PublishToQueue, "publishQueue", canopsis.CheQueueName, "Publish event to this queue.")
	flag.StringVar(&opts.ConsumeFromQueue, "consumeQueue", canopsis.FIFOQueueName, "Consume events from this queue.")
	flag.BoolVar(&opts.ModeDebug, "d", false, "debug")
	flag.BoolVar(&opts.PrintEventOnError, "printEventOnError", false, "Print event on processing error")
	flag.IntVar(&opts.LockTtl, "lockTtl", 10, "Redis lock ttl time in seconds")
	flag.BoolVar(&opts.EnableMetaAlarmProcessing, "enableMetaAlarmProcessing", true, "Enable meta-alarm processing")
	flag.DurationVar(&opts.EventsStatsFlushInterval, "eventsStatsFlushInterval", 60*time.Second, "Interval between saving statistics from redis to mongo")

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
	defer close(references.StatsCh)
	engine := NewEngineFIFO(opts, references)

	wg := sync.WaitGroup{}

	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	go func() {
		defer wg.Done()
		references.StatsListener.Listen(ctx, references.StatsCh)
	}()

	exitStatus, err := canopsis.StartEngine(ctx, engine, nil)
	if err != nil {
		logger.Error().Err(err).Msg("")
	}

	cancel()
	wg.Wait()

	trace.Stop()
	os.Exit(exitStatus)
}
