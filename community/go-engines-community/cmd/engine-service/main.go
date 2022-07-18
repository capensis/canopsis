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
	flagVersion := flag.Bool("version", false, "Show the version information")
	opts := Options{}
	flag.BoolVar(&opts.ModeDebug, "d", false, "debug")
	flag.BoolVar(&opts.FeaturePrintEventOnError, "printEventOnError", false, "Print event on processing error")
	flag.StringVar(&opts.PublishToQueue, "publishQueue", "Engine_action", "Publish event to this queue.")
	flag.DurationVar(&opts.PeriodicalWaitTime, "periodicalWaitTime", canopsis.PeriodicalWaitTime, "Duration to wait between two run of periodical process")
	flag.BoolVar(&opts.AutoRecomputeAll, "autoRecomputeAll", false, "Automatically recompute entity services each minute.")
	flag.BoolVar(&opts.RecomputeAllOnInit, "recomputeAllOnInit", false, "Recompute entity services on init.")

	flag.Parse()

	if *flagVersion {
		canopsis.PrintVersionInfo()
		return
	}

	logger := log.NewLogger(opts.ModeDebug)

	if opts.AutoRecomputeAll {
		logger.Info().Msg("Automatic entity services recomputation ENABLED")
	} else {
		logger.Info().Msg("Automatic entity services recomputation DISABLED")
	}

	if opts.RecomputeAllOnInit {
		logger.Info().Msg("Entity services recomputation on init ENABLED")
	} else {
		logger.Info().Msg("Entity services recomputation on init DISABLED")
	}

	trace := debug.Start(logger)
	// Graceful shutdown.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	engine := NewEngine(ctx, opts, logger)
	err := engine.Run(ctx)
	exitStatus := 0
	if err != nil {
		logger.Err(err).Msg("exit with error")
		exitStatus = 1
	}

	trace.Stop()
	os.Exit(exitStatus)
}
