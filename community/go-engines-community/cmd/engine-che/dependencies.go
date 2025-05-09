package main

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/externaldata"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/che"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/postgres"
	"github.com/rs/zerolog"
)

func NewEngine(ctx context.Context, opts che.Options, logger zerolog.Logger) engine.Engine {
	m := che.DependencyMaker{}
	dbClient := m.DepMongoClient(ctx, logger)
	cfg := m.DepConfig(ctx, dbClient)
	config.SetDbClientRetry(dbClient, cfg)
	pgPoolProvider := postgres.NewPoolProvider(cfg.Global.ReconnectRetries, cfg.Global.GetReconnectTimeout())
	metricsConfigProvider := config.NewMetricsConfigProvider(cfg, logger)
	metricsSender := metrics.NewTimescaleDBSender(pgPoolProvider, metricsConfigProvider, logger)
	e := che.NewEngine(ctx, opts, dbClient, cfg, metricsSender, metrics.NewNullMetaUpdater(),
		externaldata.NewGetterContainer(), config.NewTimezoneConfigProvider(cfg, logger),
		config.NewTemplateConfigProvider(cfg, logger), logger)
	e.AddDeferFunc(func(ctx context.Context) {
		err := dbClient.Disconnect(ctx)
		if err != nil {
			logger.Err(err).Msg("failed to close mongo connection")
		}

		pgPoolProvider.Close()
	})

	return e
}
