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
	opts := Options{}

	flag.StringVar(&opts.PublishToQueue, "publishQueue", canopsis.AxeQueueName, "Publish event to this queue.")
	flag.StringVar(&opts.ConsumeFromQueue, "consumeQueue", canopsis.PBehaviorQueueName, "Consume events from this queue.")
	flag.BoolVar(&opts.ModeDebug, "d", false, "debug")
	flag.IntVar(&opts.FrameDuration, "frameDuration", 120, "Pbehavior frame duration in minutes")
	flag.BoolVar(&opts.FeaturePrintEventOnError, "printEventOnError", false, "Print event on processing error")
	flag.DurationVar(&opts.PeriodicalWaitTime, "periodicalWaitTime", 60*time.Second, "Duration to wait between two run of periodical process")

	flagVersion := flag.Bool("version", false, "version infos")

	flag.Parse()

	if *flagVersion {
		canopsis.PrintVersionExit()
	}

	logger := log.NewLogger(opts.ModeDebug)
	trace := debug.Start(logger)

	ctx, cancel := context.WithCancel(context.Background())
	engine := NewEnginePBehavior(opts, logger)

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
