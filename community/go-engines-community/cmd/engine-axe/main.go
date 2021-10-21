package main

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/axe"
	"os"
	"os/signal"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/debug"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
)

func main() {
	opts := axe.ParseOptions()
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
	// Graceful shutdown.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	engine := NewEngineAXE(ctx, opts, logger)
	err := engine.Run(ctx)
	exitStatus := 0
	if err != nil {
		logger.Err(err).Msg("exit with error")
		exitStatus = 1
	}

	trace.Stop()
	os.Exit(exitStatus)
}
