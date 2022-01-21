package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/broadcastmessage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/contextgraph"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/export"
	apilogger "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/middleware"
	devmiddleware "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/middleware/dev"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/websocket"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/action"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/importcontextgraph"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	libpbehavior "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	libredis "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	libsecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/session/mongostore"
	"github.com/gin-gonic/gin"
	gorillawebsocket "github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

const chanBuf = 10
const sessionStoreSessionMaxAge = 24 * time.Hour
const sessionStoreAutoCleanInterval = 10 * time.Second

func Default(
	ctx context.Context,
	flags Flags,
	enforcer libsecurity.Enforcer,
	timezoneConfigProvider *config.BaseTimezoneConfigProvider,
	apiConfigProvider *config.BaseApiConfigProvider,
	logger zerolog.Logger,
	metricsEntityMetaUpdater metrics.MetaUpdater,
	metricsUserMetaUpdater metrics.MetaUpdater,
	exportExecutor export.TaskExecutor,
	deferFunc DeferFunc,
) (API, error) {
	// Retrieve config.
	dbClient, err := mongo.NewClient(ctx, 0, 0, logger)
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
	securityConfig, err := libsecurity.LoadConfig(flags.ConfigDir)
	if err != nil {
		logger.Err(err).Msg("cannot load security config")
		return nil, err
	}

	cookieOptions := CookieOptions{
		FileAccessName: "token",
		MaxAge:         int(sessionStoreSessionMaxAge.Seconds()),
		Secure:         flags.SecureSession,
	}
	sessionStore := mongostore.NewStore(dbClient, []byte(os.Getenv("SESSION_KEY")))
	sessionStore.Options.MaxAge = cookieOptions.MaxAge
	sessionStore.Options.Secure = cookieOptions.Secure
	if apiConfigProvider == nil {
		apiConfigProvider = config.NewApiConfigProvider(cfg, logger)
	}
	security := NewSecurity(securityConfig, dbClient, sessionStore, enforcer, apiConfigProvider, cookieOptions, logger)
	// Create pbehavior computer.
	pbhComputeChan := make(chan libpbehavior.ComputeTask, chanBuf)
	pbhEntityMatcher := libpbehavior.NewComputedEntityMatcher(dbClient, pbhRedisSession, json.NewEncoder(), json.NewDecoder())
	pbhStore := libpbehavior.NewStore(pbhRedisSession, json.NewEncoder(), json.NewDecoder())
	pbhService := libpbehavior.NewService(libpbehavior.NewModelProvider(dbClient), pbhEntityMatcher, pbhStore, libredis.NewLockClient(pbhRedisSession))
	pbhEntityTypeResolver := libpbehavior.NewEntityTypeResolver(pbhStore, pbhEntityMatcher)
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
		contextgraph.NewEventPublisher(canopsis.FIFOExchangeName, canopsis.FIFOQueueName, json.NewEncoder(), canopsis.JsonContentType, amqpChannel),
		contextgraph.NewMongoStatusReporter(dbClient),
		jobQueue,
		importcontextgraph.NewWorker(
			dbClient,
			importcontextgraph.NewEventPublisher(canopsis.FIFOExchangeName, canopsis.FIFOQueueName, json.NewEncoder(), canopsis.JsonContentType, amqpChannel),
			metricsEntityMetaUpdater,
		),
		logger,
	)

	entityCleanerTaskChan := make(chan entity.CleanTask)
	disabledEntityCleaner := entity.NewDisabledCleaner(
		entity.NewStore(dbClient),
		datastorage.NewAdapter(dbClient),
		metricsEntityMetaUpdater,
		logger,
	)

	userInterfaceAdapter := config.NewUserInterfaceAdapter(dbClient)
	userInterfaceConfig, err := userInterfaceAdapter.GetConfig(ctx)
	if err != nil && err != mongodriver.ErrNoDocuments {
		return nil, err
	}
	userInterfaceConfigProvider := config.NewUserInterfaceConfigProvider(userInterfaceConfig, logger)

	// Create and compute scenario priority intervals.
	scenarioPriorityIntervals := action.NewPriorityIntervals()
	err = scenarioPriorityIntervals.Recalculate(ctx, dbClient.Collection(mongo.ScenarioMongoCollection))
	if err != nil {
		return nil, err
	}

	// Create csv exporter.
	if exportExecutor == nil {
		exportExecutor = export.NewTaskExecutor(dbClient, logger)
	}

	websocketHub := newWebsocketHub(enforcer, security.GetTokenProvider(), logger)

	broadcastMessageChan := make(chan bool)

	// Create api.
	api := New(
		fmt.Sprintf(":%d", flags.Port),
		func(ctx context.Context) {
			close(pbhComputeChan)
			close(entityPublChan)
			close(entityCleanerTaskChan)
			close(broadcastMessageChan)

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
		router.Use(middleware.Cache())

		if flags.Test {
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
			pbhEntityTypeResolver,
			pbhComputeChan,
			entityPublChan,
			entityCleanerTaskChan,
			engine.NewRunInfoManager(engineRedisSession),
			exportExecutor,
			apilogger.NewActionLogger(dbClient, logger),
			amqpChannel,
			jobQueue,
			userInterfaceConfigProvider,
			scenarioPriorityIntervals,
			cfg.File.Upload,
			websocketHub,
			broadcastMessageChan,
			metricsEntityMetaUpdater,
			metricsUserMetaUpdater,
			logger,
		)
	})
	api.AddNoRoute(func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
	})
	api.AddNoMethod(func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusMethodNotAllowed, common.MethodNotAllowedResponse)
	})
	api.SetWebsocketHub(websocketHub)

	api.AddWorker("session clean", func(ctx context.Context) {
		security.GetSessionStore().StartAutoClean(ctx, sessionStoreAutoCleanInterval)
	})
	api.AddWorker("enforce policy load", func(ctx context.Context) {
		enforcer.StartAutoLoadPolicy(ctx)
	})
	api.AddWorker("pbehavior compute", func(ctx context.Context) {
		pbhComputer := libpbehavior.NewCancelableComputer(
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
	api.AddWorker("config reload", updateConfig(timezoneConfigProvider, apiConfigProvider,
		configAdapter, userInterfaceConfigProvider, userInterfaceAdapter, flags.Test, logger))
	api.AddWorker("data export", func(ctx context.Context) {
		exportExecutor.Execute(ctx)
	})
	api.AddWorker("auth token", func(ctx context.Context) {
		security.GetTokenStore().DeleteExpired(ctx, canopsis.PeriodicalWaitTime)
	})
	api.AddWorker("websocket", func(ctx context.Context) {
		websocketHub.Start(ctx)
	})
	broadcastMessageService := broadcastmessage.NewService(broadcastmessage.NewStore(dbClient), websocketHub, canopsis.PeriodicalWaitTime, logger)
	api.AddWorker("broadcast message", func(ctx context.Context) {
		broadcastMessageService.Start(ctx, broadcastMessageChan)
	})

	return api, nil
}

func newWebsocketHub(enforcer libsecurity.Enforcer, tokenProvider libsecurity.TokenProvider, logger zerolog.Logger, roomPerms ...map[string][]string) websocket.Hub {
	websocketUpgrader := websocket.NewUpgrader(gorillawebsocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 2048,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	})
	websocketAuthorizer := websocket.NewAuthorizer(enforcer, tokenProvider)
	websocketHub := websocket.NewHub(websocketUpgrader, websocketAuthorizer,
		canopsis.PeriodicalWaitTime, logger)
	websocketHub.RegisterRoom(websocket.RoomBroadcastMessages)
	websocketHub.RegisterRoom(websocket.RoomLoggedUserCount)
	return websocketHub
}

func updateConfig(
	timezoneConfigProvider *config.BaseTimezoneConfigProvider,
	apiConfigProvider *config.BaseApiConfigProvider,
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

				timezoneConfigProvider.Update(cfg)
				apiConfigProvider.Update(cfg)

				userInterfaceConfig, err := userInterfaceAdapter.GetConfig(ctx)
				if err != nil {
					logger.Err(err).Msg("fail to load user interface config")
					continue
				}
				userInterfaceConfigProvider.Update(userInterfaceConfig)
			case <-ctx.Done():
				return
			}
		}
	}
}
