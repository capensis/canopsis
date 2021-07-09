package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	_ "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/docs"
	libapi "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
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
func main() {
	var flags Flags
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
	dbClient, err := mongo.NewClient(0, 0)
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot connect to mongodb")
	}
	cfg, err := config.NewAdapter(dbClient).GetConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot load config")
	}
	// Set mongodb setting.
	dbClient.SetRetryCount(cfg.Global.ReconnectRetries)
	dbClient.SetMinRetryTimeout(cfg.Global.GetReconnectTimeout())
	// Init security ACL enforcer.
	enforcer, err := libsecurity.NewEnforcer(flags.ConfigDir, dbClient)
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot create security enforce")
	}

	api, err := libapi.Default(
		ctx,
		fmt.Sprintf(":%d", flags.Port),
		flags.ConfigDir,
		flags.SecureSession,
		flags.Test,
		enforcer,
		nil,
		logger,
		func(ctx context.Context) error {
			err := dbClient.Disconnect(ctx)
			if err != nil {
				return err
			}

			return nil
		},
	)
	if err != nil {
		logger.Fatal().Err(err).Msg("fail create api")
	}

	if flags.Debug {
		api.AddRouter(func(router gin.IRouter) {
			router.GET("/swagger/*any", ginswagger.WrapHandler(swaggerfiles.Handler))
		})
	}

	err = api.Run(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("fail start api")
	}
}
