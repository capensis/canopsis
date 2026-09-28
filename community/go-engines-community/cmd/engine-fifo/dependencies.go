package main

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	libengine "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/externaldata"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/depprovider"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/fifo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/postgres"
	"github.com/rs/zerolog"
)

func NewEngine(ctx context.Context, options fifo.Options, logger zerolog.Logger) (libengine.Engine, fifo.Services, error) {
	e, s, err := fifo.Default(ctx, fifo.Dependencies{
		NewMetaUpdater: func(mongo.DbClient, postgres.PoolProvider, config.CanopsisConf, zerolog.Logger) metrics.AsyncMetaUpdater {
			return metrics.NewNullAsyncMetaUpdater()
		},
		NewExternalDataGetter: func(mongo.DbClient, postgres.PoolProvider, template.Executor) externaldata.Getter {
			return externaldata.NewNullGetter()
		},
	}, options, depprovider.NewProvider(), logger)
	if err != nil {
		return e, s, err
	}

	s.MessageProcessor.ExternalData = fifo.NewNullExternalDataCoordinator()

	return e, s, nil
}
