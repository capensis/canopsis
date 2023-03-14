package main

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/che"
	"github.com/rs/zerolog"
)

func NewEngine(ctx context.Context, opts che.Options, logger zerolog.Logger) engine.Engine {
	m := che.DependencyMaker{}
	dbClient := m.DepMongoClient(ctx, logger)
	cfg := m.DepConfig(ctx, dbClient)
	config.SetDbClientRetry(dbClient, cfg)

	e := che.NewEngine(ctx, opts, dbClient, cfg, metrics.NewNullMetaUpdater(), eventfilter.NewExternalDataGetterContainer(),
		config.NewTimezoneConfigProvider(cfg, logger), config.NewTemplateConfigProvider(cfg), logger)
	e.AddDeferFunc(func(ctx context.Context) {
		err := dbClient.Disconnect(ctx)
		if err != nil {
			logger.Err(err).Msg("failed to close mongo connection")
		}
	})
	return e
}
