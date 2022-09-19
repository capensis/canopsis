package main

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/axe"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"github.com/rs/zerolog"
)

// NewEngineAXE returns the default AXE engine with default connections.
func NewEngineAXE(ctx context.Context, options axe.Options, logger zerolog.Logger) engine.Engine {
	m := axe.DependencyMaker{}
	dbClient := m.DepMongoClient(ctx, logger)
	cfg := m.DepConfig(ctx, dbClient)
	config.SetDbClientRetry(dbClient, cfg)
	e := axe.NewEngine(ctx, options, dbClient, cfg, metrics.NewNullSender(), logger)
	e.AddDeferFunc(func(ctx context.Context) {
		err := dbClient.Disconnect(ctx)
		if err != nil {
			logger.Err(err).Msg("failed to close mongo connection")
		}
	})
	return e
}
