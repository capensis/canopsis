package main

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/axe"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"github.com/rs/zerolog"
)

// NewEngineAXE returns the default AXE engine with default connections.
func NewEngineAXE(ctx context.Context, options axe.Options, logger zerolog.Logger) engine.Engine {
	return axe.Default(ctx, options, metrics.NewNullSender(), nil, logger)
}
