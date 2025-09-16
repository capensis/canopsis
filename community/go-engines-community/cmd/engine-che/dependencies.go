package main

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/externaldata"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/che"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/postgres"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/rs/zerolog"
)

func NewEngine(ctx context.Context, opts che.Options, logger zerolog.Logger) engine.Engine {
	m := che.DependencyMaker{}

	primaryDbClient := m.DepMongoClient(ctx, mongo.ClientOptions{})

	cfg := m.DepConfig(ctx, primaryDbClient)
	config.SetDbClientRetry(primaryDbClient, cfg)

	secondaryDbClient := m.DepMongoClient(ctx, mongo.ClientOptions{
		RetryCount:      cfg.Global.ReconnectRetries,
		MinRetryTimeout: cfg.Global.GetReconnectTimeout(),
		ReadPreference:  mongo.SecondaryPreferred(),
	})

	// noTimeoutClient should be used by change stream watchers only.
	noTimeoutClient := m.DepMongoClient(ctx, mongo.ClientOptions{
		RetryCount:      cfg.Global.ReconnectRetries,
		MinRetryTimeout: cfg.Global.GetReconnectTimeout(),
		NoClientTimeout: true,
	})

	eventFilterEventCounter := eventfilter.NewEventCounter(primaryDbClient,
		utils.MinDuration(canopsis.DefaultFlushInterval, opts.PeriodicalWaitTime), logger)
	eventFilterFailureService := eventfilter.NewFailureService(primaryDbClient,
		utils.MinDuration(canopsis.DefaultFlushInterval, opts.PeriodicalWaitTime), logger)
	pgPoolProvider := postgres.NewPoolProvider(cfg.Global.ReconnectRetries, cfg.Global.GetReconnectTimeout())
	metricsConfigProvider := config.NewMetricsConfigProvider(cfg, logger)
	metricsSender := metrics.NewTimescaleDBSender(pgPoolProvider, metricsConfigProvider, logger)
	e := che.NewEngine(ctx, opts, primaryDbClient, secondaryDbClient, noTimeoutClient, cfg, metricsSender,
		metrics.NewNullMetaUpdater(), externaldata.NewGetterContainer(), config.NewTimezoneConfigProvider(cfg, logger),
		config.NewTemplateConfigProvider(cfg, logger), eventFilterEventCounter, eventFilterFailureService, logger)
	e.AddDeferFunc(func(ctx context.Context) {
		err := primaryDbClient.Disconnect(ctx)
		if err != nil {
			logger.Err(err).Msg("failed to close primary mongo connection")
		}

		err = secondaryDbClient.Disconnect(ctx)
		if err != nil {
			logger.Err(err).Msg("failed to close secondary mongo connection")
		}

		err = noTimeoutClient.Disconnect(ctx)
		if err != nil {
			logger.Err(err).Msg("failed to close mongo connection without timeout")
		}

		pgPoolProvider.Close()
	})

	return e
}
