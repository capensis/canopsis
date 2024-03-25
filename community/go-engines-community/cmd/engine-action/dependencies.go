package main

import (
	"context"

	libaction "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/action"
	libconfig "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	libengine "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"github.com/rs/zerolog"
)

func NewEngineAction(
	ctx context.Context,
	options libaction.Options,
	logger zerolog.Logger,
) libengine.Engine {
	m := libaction.DependencyMaker{}
	dbClient := m.DepMongoClient(ctx, logger)
	cfg := m.DepConfig(ctx, dbClient)
	libconfig.SetDbClientRetry(dbClient, cfg)
	amqpConnection := m.DepAmqpConnection(logger, cfg)
	engine := libaction.NewEngineAction(ctx, options, cfg, dbClient, amqpConnection, nil, nil, logger)
	engine.AddDeferFunc(func(ctx context.Context) {
		err := dbClient.Disconnect(ctx)
		if err != nil {
			logger.Err(err).Msg("failed to close mongo connection")
		}

		err = amqpConnection.Close()
		if err != nil {
			logger.Error().Err(err).Msg("failed to close amqp connection")
		}
	})

	return engine
}
