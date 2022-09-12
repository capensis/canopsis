package main

import (
	"context"
	"os"
	"os/signal"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/debug"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/che"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
)

func main() {
	opts := che.ParseOptions()

	if opts.Version {
		canopsis.PrintVersionInfo()
		return
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
