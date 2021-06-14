package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/debug"
	"git.canopsis.net/canopsis/go-engines/lib/log"
)

func main() {
	flagVersion := flag.Bool("version", false, "version infos")

	opts := Options{}

	flag.BoolVar(&opts.ModeDebug, "d", false, "debug")
	flag.BoolVar(&opts.FeaturePrintEventOnError, "printEventOnError", false, "Print event on processing error")
	flag.StringVar(&opts.PublishToQueue, "publishQueue", "Engine_action", "Publish event to this queue.")
	flag.DurationVar(&opts.PeriodicalWaitTime, "periodicalWaitTime", 60 * time.Second, "Duration to wait between two run of periodical process")
	flag.BoolVar(&opts.AutoRecomputeWatchers, "autoRecomputeWatchers", false, "Automatically recompute watchers each minute.")

	flag.Parse()

	if *flagVersion {
		canopsis.PrintVersionExit()
	}

	logger := log.NewLogger(opts.ModeDebug)

	if opts.AutoRecomputeWatchers {
		logger.Info().Msg("Automatic watcher recomputation ENABLED")
	} else {
		logger.Info().Msg("Automatic watcher recomputation DISABLED")
	}

	trace := debug.Start(logger)

	ctx, cancel := context.WithCancel(context.Background())
	engine := NewEngineWatcher(opts, logger)

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		logger.Info().Msg("engine is stopping")
		cancel()
	}()

	err := engine.Run(ctx)
	exitStatus := 0
	if err != nil {
		logger.Err(err).Msg("exit with error")
		exitStatus = 1
	}

	trace.Stop()
	os.Exit(exitStatus)
}
