package main

//go:generate swag init  -d ../../lib -g ../cmd/canopsis-api-community/main.go -o ../../docs

import (
	"context"
	"os"
	"os/signal"

	_ "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/docs"
	libapi "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	liblog "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	libsecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
)

// @title Canopsis API
// @version 4.0
// @description This is a Canopsis server.

// @BasePath /api/v4
// @query.collection.format multi

// @securityDefinitions.basic BasicAuth

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name x-canopsis-authkey

// @securityDefinitions.apikey JWTAuth
// @in header
// @name Authorization
func main() {
	var flags libapi.Flags
	flags.ParseArgs()
	logger := liblog.NewLogger(flags.Debug)
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
	enforcer, err := libsecurity.NewEnforcer(flags.ConfigDir, dbClient)
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot create security enforce")
	}

	api, err := libapi.Default(
		ctx,
		flags,
		enforcer,
		nil, nil,
		logger,
		metrics.NewNullMetaUpdater(),
		metrics.NewNullMetaUpdater(),
		nil,
		func(ctx context.Context) {
			err := dbClient.Disconnect(ctx)
			if err != nil {
				logger.Error().Err(err).Msg("failed to close mongo connection")
			}
		},
	)
	if err != nil {
		logger.Fatal().Err(err).Msg("fail create api")
	}

	if flags.EnableDocs {
		api.AddRouter(func(router gin.IRouter) {
			router.GET("/swagger/*any", ginswagger.WrapHandler(swaggerfiles.Handler))
		})
	}

	err = api.Run(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("fail start api")
	}
}
