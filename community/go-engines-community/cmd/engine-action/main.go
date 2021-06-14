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
	flag.StringVar(&opts.FifoAckExchange, "fifoAckExchange", canopsis.FIFOAckExchangeName, "Publish FIFO Ack event to this exchange.")
	flag.StringVar(&opts.FifoAckQueue, "fifoAckQueue", canopsis.FIFOAckQueueName, "Publish FIFO Ack event to this queue.")
	flag.DurationVar(&opts.PeriodicalWaitTime, "periodicalWaitTime", canopsis.PeriodicalWaitTime, "Duration to wait between two run of periodical process")
	flag.IntVar(&opts.WorkerPoolSize, "workerPoolSize", 10, "Number of workers for scenario executions")
	flag.BoolVar(&opts.WithWebhook, "withWebhook", false, "Handle webhook actions")

	flag.Parse()

	if *flagVersion {
		canopsis.PrintVersionExit()
	}

	logger := log.NewLogger(opts.ModeDebug)

	trace := debug.Start(logger)

	ctx, cancel := context.WithCancel(context.Background())
	engine := NewEngineAction(ctx, opts, logger)

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
