package main

import (
	"context"

	libengine "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/depmake"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/fifo"
	"github.com/rs/zerolog"
)

func NewEngine(ctx context.Context, options fifo.Options, logger zerolog.Logger) libengine.Engine {
	defer depmake.Catch(logger)

	engine, _ := fifo.Default(ctx, metrics.NewNullMetaUpdater(), options, logger)

	return engine
}
