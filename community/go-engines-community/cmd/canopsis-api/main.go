package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	_ "git.canopsis.net/canopsis/go-engines/docs"
	libapi "git.canopsis.net/canopsis/go-engines/lib/api"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/config"
	liblog "git.canopsis.net/canopsis/go-engines/lib/log"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	libsecurity "git.canopsis.net/canopsis/go-engines/lib/security"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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
	// Init security ACL enforcer.
	dbClient, err = mongo.NewClient(cfg.Global.ReconnectRetries, cfg.Global.GetReconnectTimeout())
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot connect to mongodb")
	}
	enforcer, err := libsecurity.NewEnforcer(flags.ConfigDir, dbClient)
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot create security enforce")
	}

	api, err := libapi.Default(
		fmt.Sprintf(":%d", flags.Port),
		flags.ConfigDir,
		flags.SecureSession,
		flags.Test,
		enforcer,
		logger,
	)
	if err != nil {
		logger.Fatal().Err(err).Msg("fail create api")
	}

	if flags.Debug {
		api.AddRouter(func(router gin.IRouter) {
			router.GET("/swagger/*any", ginswagger.WrapHandler(swaggerfiles.Handler))
		})
	}

	// Graceful shutdown.
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		cancel()
	}()

	err = api.Run(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("fail start api")
	}
}
