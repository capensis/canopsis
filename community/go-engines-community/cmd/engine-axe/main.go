package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/debug"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
)

func main() {
	flagVersion := flag.Bool("version", false, "version infos")

	opts := Options{}

	flag.BoolVar(&opts.ModeDebug, "d", false, "debug")
	flag.BoolVar(&opts.FeaturePrintEventOnError, "printEventOnError", false, "Print event on processing error")
	flag.BoolVar(&opts.FeatureHideResources, "featureHideResources", false, "Enable Hide Resources Management - deprecated")
	flag.BoolVar(&opts.FeatureStatEvents, "featureStatEvents", false, "Send statistic events")
	flag.StringVar(&opts.PublishToQueue, "publishQueue", canopsis.ServiceQueueName, "Publish event to this queue")
	flag.StringVar(&opts.PostProcessorsDirectory, "postProcessorsDirectory", ".", "The path of the directory containing the post-processing plugins.")
	flag.BoolVar(&opts.IgnoreDefaultTomlConfig, "ignoreDefaultTomlConfig", false, "load toml file values into database. - deprecated")
	flag.DurationVar(&opts.PeriodicalWaitTime, "periodicalWaitTime", canopsis.PeriodicalWaitTime, "Duration to wait between two run of periodical process")
	flag.Parse()

	if *flagVersion {
		canopsis.PrintVersionExit()
	}

	logger := log.NewLogger(opts.ModeDebug)

	if opts.FeatureHideResources {
		logger.Info().Msg("featureHideResources option is deprecated")
	}

	if opts.PostProcessorsDirectory != "." {
		logger.Info().Msg("postProcessorsDirectory option is deprecated")
	}

	if opts.IgnoreDefaultTomlConfig {
		logger.Info().Msg("ignoreDefaultTomlConfig option is deprecated")
	}

	if opts.FeatureStatEvents {
		logger.Info().Msg("Statistic Events ENABLED")
	} else {
		logger.Info().Msg("Statistic Events DISABLED")
	}

	trace := debug.Start(logger)

	ctx, cancel := context.WithCancel(context.Background())
	engine := NewEngineAXE(ctx, opts, logger)

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
