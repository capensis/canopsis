package main

import (
	"context"
	"os"
	"os/signal"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/axe"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/debug"
	libflag "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/flag"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
)

func main() {
	opts, deprecatedFlags := axe.ParseOptions()

	if opts.Version {
		canopsis.PrintVersionInfo()
		return
	}

	logger := log.NewLogger(opts.ModeDebug)
	trace := debug.Start(logger)
	// Graceful shutdown.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	libflag.LogDeprecatedFlags(logger, deprecatedFlags)

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
