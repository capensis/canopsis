package main

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/debug"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/fifo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
	"os"
	"os/signal"
)

func main() {
	opts := fifo.ParseOptions()
	logger := log.NewLogger(opts.ModeDebug)
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
