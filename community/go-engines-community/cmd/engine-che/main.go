package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/debug"
	"git.canopsis.net/canopsis/go-engines/lib/log"
)

func main() {
	opts := Options{}

	flag.BoolVar(&opts.FeatureEventProcessing, "processEvent", true, "enable event processing. enabled by default.")
	flag.BoolVar(&opts.FeatureContextCreation, "createContext", true, "enable context graph creation. enabled by default. WARNING: disable the old context-graph engine when using this.")
	flag.BoolVar(&opts.FeatureContextEnrich, "enrichContext", false, "enable context graph enrichment from event. disabled by default. WARNING: disable the old context-graph engine when using this.")
	flag.StringVar(&opts.PublishToQueue, "publishQueue", canopsis.PBehaviorQueueName, "Publish event to this queue.")
	flag.StringVar(&opts.ConsumeFromQueue, "consumeQueue", canopsis.CheQueueName, "Consume events from this queue.")
	flag.StringVar(&opts.EnrichExclude, "enrichExclude", "", "Coma separated list of fields that shall not be part of context enrichment.")
	flag.StringVar(&opts.EnrichInclude, "enrichInclude", "", "Coma separated list of the only fields that will be part of context enrichment. If present, -enrichExclude is ignored.")
	flag.StringVar(&opts.DataSourceDirectory, "dataSourceDirectory", ".", "The path of the directory containing the event filter's data source plugins.")
	flag.BoolVar(&opts.ModeDebug, "d", false, "debug")
	flag.BoolVar(&opts.PrintEventOnError, "printEventOnError", false, "Print event on processing error")
	flag.BoolVar(&opts.Purge, "purge", false, "purge consumer queue(s) before work")
	flag.DurationVar(&opts.PeriodicalWaitTime, "periodicalWaitTime", canopsis.PeriodicalWaitTime, "Duration to wait between two run of periodical process")
	flag.StringVar(&opts.FifoAckExchange, "fifoAckExchange", canopsis.FIFOAckExchangeName, "Publish FIFO Ack event to this exchange.")

	flagVersion := flag.Bool("version", false, "version infos")

	flag.Parse()

	if *flagVersion {
		canopsis.PrintVersionExit()
	}

	logger := log.NewLogger(opts.ModeDebug)

	trace := debug.Start(logger)

	if opts.FeatureEventProcessing {
		logger.Info().Msg("Event processing ENABLED")
	} else {
		logger.Info().Msg("Event processing DISABLED")
	}

	if opts.FeatureContextCreation {
		logger.Info().Msg("Context creation ENABLED")
	} else {
		logger.Info().Msg("Context creation DISABLED")
	}

	if opts.FeatureContextEnrich {
		logger.Info().Msg("Context enrichment ENABLED")
	} else {
		logger.Info().Msg("Context enrichment DISABLED")
	}

	ctx, cancel := context.WithCancel(context.Background())
	engine := NewEngineCHE(ctx, opts, logger)

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
