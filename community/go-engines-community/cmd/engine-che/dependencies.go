package main

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/externaldata"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/che"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/depprovider"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/postgres"
	"github.com/rs/zerolog"
)

func NewEngine(ctx context.Context, opts che.Options, dp depprovider.Provider, logger zerolog.Logger) (engine.Engine, error) {
	e, s, err := che.NewEngine(ctx, opts, che.Dependencies{
		NewMetricsSender: func(_ mongo.DbClient, pgPool postgres.PoolProvider, cfg config.CanopsisConf, logger zerolog.Logger) metrics.Sender {
			return metrics.NewTimescaleDBSender(pgPool, config.NewMetricsConfigProvider(cfg, logger), logger)
		},
		NewMetaUpdater: func(mongo.DbClient, postgres.PoolProvider, config.CanopsisConf, zerolog.Logger) metrics.AsyncMetaUpdater {
			return metrics.NewNullAsyncMetaUpdater()
		},
		NewEntityInfosUpdateSender: func(mongo.DbClient, postgres.PoolProvider, config.CanopsisConf, zerolog.Logger) metrics.EntityInfosUpdateSender {
			return metrics.NewNullEntityInfosUpdateSender()
		},
		NewExternalDataGetter: func(mongo.DbClient, postgres.PoolProvider, template.Executor) externaldata.Getter {
			return externaldata.NewNullGetter()
		},
	}, dp, logger)
	if err != nil {
		return nil, err
	}

	s.MessageProcessor.ExternalData = che.NewNullExternalDataCoordinator()

	return e, nil
}
