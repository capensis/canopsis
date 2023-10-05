package main

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	libengine "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/depmake"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/fifo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/rs/zerolog"
)

func NewEngine(ctx context.Context, options fifo.Options, logger zerolog.Logger) libengine.Engine {
	defer depmake.Catch(logger)

	var m depmake.DependencyMaker
	dbClient := m.DepMongoClient(ctx, logger)
	cfg := m.DepConfig(ctx, dbClient)
	config.SetDbClientRetry(dbClient, cfg)
	eventFilterEventCounter := eventfilter.NewEventCounter(dbClient,
		utils.MinDuration(canopsis.DefaultFlushInterval, options.PeriodicalWaitTime), logger)
	eventFilterFailureService := eventfilter.NewFailureService(dbClient,
		utils.MinDuration(canopsis.DefaultFlushInterval, options.PeriodicalWaitTime), logger)

	return fifo.Default(ctx, options, dbClient, cfg, eventfilter.NewExternalDataGetterContainer(),
		config.NewTimezoneConfigProvider(cfg, logger), config.NewTemplateConfigProvider(cfg, logger), eventFilterEventCounter,
		eventFilterFailureService, logger)
}
