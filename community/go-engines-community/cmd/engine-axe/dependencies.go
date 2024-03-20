package main

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/axe"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/axe/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/postgres"
	"github.com/rs/zerolog"
)

// NewEngineAXE returns the default AXE engine with default connections.
func NewEngineAXE(ctx context.Context, options axe.Options, logger zerolog.Logger) engine.Engine {
	m := axe.DependencyMaker{}
	dbClient := m.DepMongoClient(ctx, logger)
	cfg := m.DepConfig(ctx, dbClient)
	config.SetDbClientRetry(dbClient, cfg)
	amqpConnection := m.DepAmqpConnection(logger, cfg)
	pgPoolProvider := postgres.NewPoolProvider(cfg.Global.ReconnectRetries, cfg.Global.GetReconnectTimeout())
	metricsConfigProvider := config.NewMetricsConfigProvider(cfg, logger)
	metricsSender := metrics.NewTimescaleDBSender(pgPoolProvider, metricsConfigProvider, logger)
	e := axe.NewEngine(ctx, options, dbClient, amqpConnection, cfg, metricsSender, event.NewNullAutoInstructionMatcher(), nil, nil, nil, logger)
	e.AddDeferFunc(func(ctx context.Context) {
		err := dbClient.Disconnect(ctx)
		if err != nil {
			logger.Err(err).Msg("failed to close mongo connection")
		}

		err = amqpConnection.Close()
		if err != nil {
			logger.Error().Err(err).Msg("failed to close amqp connection")
		}

		pgPoolProvider.Close()
	})

	return e
}
