package main

import (
	"context"
	"os"
	"os/signal"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/debug"
	libflag "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/flag"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/che"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
)

func main() {
	// Graceful shutdown.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	opts, deprecatedFlags := che.ParseOptions()
	if opts.Version {
		canopsis.PrintVersionInfo()
		return
	}

	logger := log.NewLogger(ctx, opts.Options)
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

	libflag.LogDeprecatedFlags(logger, deprecatedFlags)

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
