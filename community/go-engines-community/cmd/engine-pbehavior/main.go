package main

import (
	"context"
	"flag"
	"os"
	"os/signal"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/debug"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
)

func main() {
	opts := Options{}

	flag.StringVar(&opts.PublishToQueue, "publishQueue", canopsis.AxeQueueName, "Publish event to this queue.")
	flag.StringVar(&opts.ConsumeFromQueue, "consumeQueue", canopsis.PBehaviorQueueName, "Consume events from this queue.")
	flag.BoolVar(&opts.ModeDebug, "d", false, "debug")
	flag.IntVar(&opts.FrameDuration, "frameDuration", 120, "The engine computes all pbehaviors for a further interval which duration controls this parameter. The default value is 120 minutes. This could be reduced when pre-compute takes too much system resources.")
	flag.BoolVar(&opts.FeaturePrintEventOnError, "printEventOnError", false, "Print event on processing error")
	flag.DurationVar(&opts.PeriodicalWaitTime, "periodicalWaitTime", canopsis.PeriodicalWaitTime, "Duration to wait between two run of periodical process")
	flag.StringVar(&opts.FifoAckExchange, "fifoAckExchange", canopsis.FIFOAckExchangeName, "Publish FIFO Ack event to this exchange.")
	flag.BoolVar(&opts.FeatureResolveDeadlineDisabled, "disableResolveDeadline", false, "Resolve entites may last unlimited period")

	flagVersion := flag.Bool("version", false, "version infos")

	flag.Parse()

	if *flagVersion {
		canopsis.PrintVersionExit()
	}

	logger := log.NewLogger(opts.ModeDebug)
	trace := debug.Start(logger)

	// Graceful shutdown.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	engine := NewEnginePBehavior(ctx, opts, logger)
	err := engine.Run(ctx)
	exitStatus := 0
	if err != nil {
		logger.Err(err).Msg("exit with error")
		exitStatus = 1
	}

	trace.Stop()
	os.Exit(exitStatus)
}
