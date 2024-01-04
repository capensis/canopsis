package main

//go:generate swag init  -d ../../lib -g ../cmd/canopsis-api-community/$GOFILE -o ../../lib/api/docs --outputTypes yaml --instanceName schemas

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
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"github.com/gin-gonic/gin"
)

// @title Generated schemas
// @description This doc contains auto generated Open API v2 schemas of requests and responses to use in Open Api v3 doc.
// @version 4.0.0
func main() {
	var flags api.Flags
	flags.ParseArgs()

	if flags.Version {
		canopsis.PrintVersionInfo()
		return
	}

	logger := log.NewLogger(flags.Debug)
	// Graceful shutdown.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if flags.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Retrieve config.
	dbClient, err := mongo.NewClient(ctx, 0, 0, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot connect to mongodb")
	}
	cfg, err := config.NewAdapter(dbClient).GetConfig(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot load config")
	}
	// Set mongodb setting.
	config.SetDbClientRetry(dbClient, cfg)
	// Init security ACL enforcer.
	enforcer, err := security.NewEnforcer(flags.ConfigDir, dbClient)
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot create security enforce")
	}

	pgPoolProvider := postgres.NewPoolProvider(cfg.Global.ReconnectRetries, cfg.Global.GetReconnectTimeout())
	providers := &api.ConfigProviders{}
	server, _, err := api.Default(
		ctx,
		flags,
		enforcer,
		providers,
		logger,
		pgPoolProvider,
		metrics.NewNullMetaUpdater(),
		metrics.NewNullMetaUpdater(),
		nil,
		nil,
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
