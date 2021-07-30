package api

import (
	"context"
	"os"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/contextgraph"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/export"
	apilogger "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/middleware"
	devmiddleware "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/middleware/dev"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	libpbehavior "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	libredis "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	libsecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/proxy"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/session/mongostore"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

const chanBuf = 10
const sessionStoreSessionMaxAge = 24 * time.Hour
const sessionStoreAutoCleanInterval = 10 * time.Second

func Default(
	ctx context.Context,
	addr string,
	configDir string,
	secureSession bool,
	test bool,
	enforcer libsecurity.Enforcer,
	timezoneConfigProvider *config.BaseTimezoneConfigProvider,
	logger zerolog.Logger,
	deferFunc DeferFunc,
) (API, error) {
	// Retrieve config.
	dbClient, err := mongo.NewClient(ctx, 0, 0)
	if err != nil {
		logger.Err(err).Msg("cannot connect to mongodb")
		return nil, err
	}
	configAdapter := config.NewAdapter(dbClient)
	cfg, err := configAdapter.GetConfig(ctx)
	if err != nil {
		logger.Err(err).Msg("cannot load config")
		return nil, err
	}
	if timezoneConfigProvider == nil {
		timezoneConfigProvider = config.NewTimezoneConfigProvider(cfg, logger)
	}
	// Set mongodb setting.
	config.SetDbClientRetry(dbClient, cfg)
	// Connect to rmq.
	amqpConn, err := amqp.NewConnection(logger, -1, cfg.Global.GetReconnectTimeout())
	if err != nil {
		logger.Err(err).Msg("cannot connect to rmq")
		return nil, err
	}
	amqpChannel, err := amqpConn.Channel()
	if err != nil {
		logger.Err(err).Msg("cannot connect to rmq")
		return nil, err
	}
	// Connect to redis.
	pbhRedisSession, err := libredis.NewSession(ctx, libredis.PBehaviorLockStorage, logger,
		cfg.Global.ReconnectRetries, cfg.Global.GetReconnectTimeout())
	if err != nil {
		logger.Err(err).Msg("cannot connect to redis")
		return nil, err
	}
	engineRedisSession, err := libredis.NewSession(ctx, libredis.EngineRunInfo, logger,
		cfg.Global.ReconnectRetries, cfg.Global.GetReconnectTimeout())
	if err != nil {
		logger.Err(err).Msg("cannot connect to redis")
		return nil, err
	}
	securityConfig, err := libsecurity.LoadConfig(configDir)
	if err != nil {
		logger.Err(err).Msg("cannot load security config")
		return nil, err
	}

	sessionStore := mongostore.NewStore(dbClient, []byte(os.Getenv("SESSION_KEY")))
	sessionStore.Options.MaxAge = int(sessionStoreSessionMaxAge.Seconds())
	sessionStore.Options.Secure = secureSession
	security := NewSecurity(securityConfig, dbClient, sessionStore, logger)

	proxyAccessConfig, err := proxy.LoadAccessConfig(configDir)
	if err != nil {
		logger.Err(err).Msg("cannot load access config")
		return nil, err
	}
	// Create pbehavior computer.
	pbhComputeChan := make(chan libpbehavior.ComputeTask, chanBuf)
	pbhStore := libredis.NewStore(pbhRedisSession, libredis.PbehaviorKey, 0)
	pbhService := libpbehavior.NewService(
		libpbehavior.NewModelProvider(dbClient),
		libpbehavior.NewEntityMatcher(dbClient),
		logger,
	)
	// Create entity service event publisher.
	entityPublChan := make(chan entityservice.ChangeEntityMessage, chanBuf)
	entityServiceEventPublisher := entityservice.NewEventPublisher(
		alarm.NewAdapter(dbClient), amqpChannel,
		json.NewEncoder(), canopsis.JsonContentType,
		canopsis.FIFOAckExchangeName, canopsis.FIFOQueueName, logger,
	)

	jobQueue := contextgraph.NewJobQueue()
	importWorker := contextgraph.NewImportWorker(
		cfg,
		dbClient,
		contextgraph.NewRMQPublisher(json.NewEncoder(), amqpChannel),
		contextgraph.NewMongoStatusReporter(dbClient),
		jobQueue,
		logger,
	)

	entityCleanerTaskChan := make(chan entity.CleanTask)
	disabledEntityCleaner := entity.NewDisabledCleaner(
		entity.NewStore(dbClient),
		datastorage.NewAdapter(dbClient),
		logger,
	)

	userInterfaceAdapter := config.NewUserInterfaceAdapter(dbClient)
	userInterfaceConfig, err := userInterfaceAdapter.GetConfig(ctx)
	if err != nil && err != mongodriver.ErrNoDocuments {
		return nil, err
	}
	userInterfaceConfigProvider := config.NewUserInterfaceConfigProvider(userInterfaceConfig, logger)

	// Create csv exporter.
	exportExecutor := export.NewTaskExecutor(dbClient, logger)

	// Create api.
	api := New(
		addr,
		func(ctx context.Context) {
			close(pbhComputeChan)
			close(entityPublChan)
			close(entityCleanerTaskChan)

			err := dbClient.Disconnect(ctx)
			if err != nil {
				logger.Error().Err(err).Msg("failed to close mongo connection")
			}
			err = amqpConn.Close()
			if err != nil {
				logger.Error().Err(err).Msg("failed to close amqp connection")
			}

			err = pbhRedisSession.Close()
			if err != nil {
				logger.Error().Err(err).Msg("failed to close redis connection")
			}

			err = engineRedisSession.Close()
			if err != nil {
				logger.Error().Err(err).Msg("failed to close redis connection")
			}

			if deferFunc != nil {
				deferFunc(ctx)
			}
		},
		logger,
	)
	api.AddRouter(func(router gin.IRouter) {
		corsConfig := cors.DefaultConfig()
		corsConfig.AllowAllOrigins = true
		corsConfig.AllowCredentials = true
		router.Use(cors.New(corsConfig))
		router.Use(middleware.Cache())

		if test {
			router.Use(devmiddleware.ReloadEnforcerPolicy(enforcer))
		}

		RegisterValidators(dbClient)
		RegisterRoutes(
			ctx,
			cfg,
			router,
			security,
			enforcer,
			dbClient,
			timezoneConfigProvider,
			pbhStore,
			pbhService,
			pbhComputeChan,
			entityPublChan,
			entityCleanerTaskChan,
			engine.NewRunInfoManager(engineRedisSession),
			exportExecutor,
			apilogger.NewActionLogger(dbClient, logger),
			amqpChannel,
			jobQueue,
			userInterfaceConfigProvider,
			logger,
		)
	})
	api.AddNoRoute(GetProxy(security, enforcer, proxyAccessConfig))

	api.AddWorker("session clean", func(ctx context.Context) {
		security.GetSessionStore().StartAutoClean(ctx, sessionStoreAutoCleanInterval)
	})
	api.AddWorker("enforce policy load", func(ctx context.Context) {
		enforcer.StartAutoLoadPolicy(ctx)
	})
	api.AddWorker("pbehavior compute", func(ctx context.Context) {
		pbhComputer := libpbehavior.NewCancelableComputer(
			libredis.NewLockClient(pbhRedisSession),
			pbhStore,
			pbhService,
			dbClient,
			amqpChannel,
			libpbehavior.NewEventManager(),
			json.NewEncoder(),
			canopsis.FIFOQueueName,
			logger,
		)
		pbhComputer.Compute(ctx, pbhComputeChan)
	})
	api.AddWorker("entity event publish", func(ctx context.Context) {
		entityServiceEventPublisher.Publish(ctx, entityPublChan)
	})
	api.AddWorker("entity cleaner", func(ctx context.Context) {
		disabledEntityCleaner.RunCleanerProcess(ctx, entityCleanerTaskChan)
	})
	api.AddWorker("import job", func(ctx context.Context) {
		importWorker.Run(ctx)
	})
	api.AddWorker("config reload", updateConfig(timezoneConfigProvider, configAdapter,
		userInterfaceConfigProvider, userInterfaceAdapter, test, logger))
	api.AddWorker("data export", func(ctx context.Context) {
		exportExecutor.Execute(ctx)
	})

	return api, nil
}

func updateConfig(
	timezoneConfigProvider *config.BaseTimezoneConfigProvider,
	configAdapter config.Adapter,
	userInterfaceConfigProvider *config.BaseUserInterfaceConfigProvider,
	userInterfaceAdapter config.UserInterfaceAdapter,
	test bool,
	logger zerolog.Logger,
) func(ctx context.Context) {
	return func(ctx context.Context) {
		timeout := canopsis.PeriodicalWaitTime
		if test {
			timeout = time.Second
		}

		ticker := time.NewTicker(timeout)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				cfg, err := configAdapter.GetConfig(ctx)
				if err != nil {
					logger.Err(err).Msg("fail to load config")
					continue
				}

				err = timezoneConfigProvider.Update(cfg)
				if err != nil {
					logger.Err(err).Msg("fail to load config")
					continue
				}

				userInterfaceConfig, err := userInterfaceAdapter.GetConfig(ctx)
				if err != nil {
					logger.Err(err).Msg("fail to load user interface config")
					continue
				}
				err = userInterfaceConfigProvider.Update(userInterfaceConfig)
				if err != nil {
					logger.Err(err).Msg("fail to load user interface config")
					continue
				}
			case <-ctx.Done():
				return
			}
		}
	}
}
