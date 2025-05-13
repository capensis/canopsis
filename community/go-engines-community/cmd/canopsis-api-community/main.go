package main

//go:generate go tool github.com/swaggo/swag/cmd/swag init -d ../../lib -g ../cmd/canopsis-api-community/$GOFILE -o ../../lib/api/docs --outputTypes yaml --instanceName schemas

import (
	"context"
	"os"
	"os/signal"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/postgres"
	"github.com/gin-gonic/gin"
)

// @title Generated schemas
// @description This doc contains auto generated Open API v2 schemas of requests and responses to use in Open Api v3 doc.
// @version 4.0.0
func main() {
	// Graceful shutdown.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	var flags api.Flags
	flags.ParseArgs()

	if flags.Version {
		canopsis.PrintVersionInfo()
		return
	}

	if flags.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	logger := log.NewLogger(ctx, flags.Debug)

	// Retrieve config.
	dbClient, err := mongo.NewClient(ctx, mongo.ClientOptions{})
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot connect to mongodb")
	}
	cfg, err := config.NewAdapter(dbClient).GetConfig(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot load config")
	}
	// Set mongodb setting.
	config.SetDbClientRetry(dbClient, cfg)

	pgPoolProvider := postgres.NewPoolProvider(cfg.Global.ReconnectRetries, cfg.Global.GetReconnectTimeout())
	server, _, err := api.Default(
		ctx,
		flags,
		logger,
		pgPoolProvider,
		metrics.NewNullMetaUpdater(),
		metrics.NewNullMetaUpdater(),
		func(ctx context.Context) {
			err := dbClient.Disconnect(ctx)
			if err != nil {
				logger.Error().Err(err).Msg("failed to close mongo connection")
			}

			pgPoolProvider.Close()
		},
		false,
	)
	if err != nil {
		logger.Fatal().Err(err).Msg("fail create api")
	}

	err = server.Run(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("fail start api")
	}
}
