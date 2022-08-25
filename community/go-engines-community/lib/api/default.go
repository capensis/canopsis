package api

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/broadcastmessage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/contextgraph"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/docs"
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
	libentity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/importcontextgraph"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	libpbehavior "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/postgres"
	libredis "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	libsecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/proxy"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/session/mongostore"
	"github.com/gin-gonic/gin"
	gorillawebsocket "github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

const chanBuf = 10
const sessionStoreSessionMaxAge = 24 * time.Hour
const sessionStoreAutoCleanInterval = 10 * time.Second

//go:embed swaggerui/*
var docsUiFile embed.FS

//go:embed docs/*.yaml
var docsFile embed.FS

type ConfigProviders struct {
	TimezoneConfigProvider      *config.BaseTimezoneConfigProvider
	ApiConfigProvider           *config.BaseApiConfigProvider
	UserInterfaceConfigProvider *config.BaseUserInterfaceConfigProvider
}

func Default(
	ctx context.Context,
	flags Flags,
	enforcer libsecurity.Enforcer,
	p *ConfigProviders,
	logger zerolog.Logger,
	metricsEntityMetaUpdater metrics.MetaUpdater,
	metricsUserMetaUpdater metrics.MetaUpdater,
	exportExecutor export.TaskExecutor,
	deferFunc DeferFunc,
	overrideDocs bool,
) (API, fs.ReadFileFS, error) {
	// Retrieve config.
	dbClient, err := mongo.NewClient(ctx, 0, 0, logger)
	if err != nil {
		logger.Err(err).Msg("cannot connect to mongodb")
		return nil, nil, err
	}
	configAdapter := config.NewAdapter(dbClient)
	cfg, err := configAdapter.GetConfig(ctx)
	if err != nil {
		logger.Err(err).Msg("cannot load config")
		return nil, nil, err
	}
	if p.TimezoneConfigProvider == nil {
		p.TimezoneConfigProvider = config.NewTimezoneConfigProvider(cfg, logger)
	}
	// Set mongodb setting.
	config.SetDbClientRetry(dbClient, cfg)
	// Connect to rmq.
	amqpConn, err := amqp.NewConnection(logger, -1, cfg.Global.GetReconnectTimeout())
	if err != nil {
		logger.Err(err).Msg("cannot connect to rmq")
		return nil, nil, err
	}
	amqpChannel, err := amqpConn.Channel()
	if err != nil {
		logger.Err(err).Msg("cannot connect to rmq")
		return nil, nil, err
	}
	// Connect to redis.
	pbhRedisSession, err := libredis.NewSession(ctx, libredis.PBehaviorLockStorage, logger,
		cfg.Global.ReconnectRetries, cfg.Global.GetReconnectTimeout())
	if err != nil {
		logger.Err(err).Msg("cannot connect to redis")
		return nil, nil, err
	}
	engineRedisSession, err := libredis.NewSession(ctx, libredis.EngineRunInfo, logger,
		cfg.Global.ReconnectRetries, cfg.Global.GetReconnectTimeout())
	if err != nil {
		logger.Err(err).Msg("cannot connect to redis")
		return nil, nil, err
	}
	securityConfig, err := libsecurity.LoadConfig(flags.ConfigDir)
	if err != nil {
		logger.Err(err).Msg("cannot load security config")
		return nil, nil, err
	}

	cookieOptions := CookieOptions{
		FileAccessName: "token",
		MaxAge:         int(sessionStoreSessionMaxAge.Seconds()),
		Secure:         flags.SecureSession,
	}
	sessionStore := mongostore.NewStore(dbClient, []byte(os.Getenv("SESSION_KEY")))
	sessionStore.Options.MaxAge = cookieOptions.MaxAge
	sessionStore.Options.Secure = cookieOptions.Secure
	if p.ApiConfigProvider == nil {
		p.ApiConfigProvider = config.NewApiConfigProvider(cfg, logger)
	}
	security := NewSecurity(securityConfig, dbClient, sessionStore, enforcer, p.ApiConfigProvider, cookieOptions, logger)

	if flags.EnableSameServiceNames {
		logger.Info().Msg("Non-unique names for services ENABLED")
	}

	proxyAccessConfig, err := proxy.LoadAccessConfig(flags.ConfigDir)
	if err != nil {
		logger.Err(err).Msg("cannot load access config")
		return nil, nil, err
	}
	// Create pbehavior computer.
	pbhComputeChan := make(chan libpbehavior.ComputeTask, chanBuf)
	pbhStore := libpbehavior.NewStore(pbhRedisSession, json.NewEncoder(), json.NewDecoder())
	pbhService := libpbehavior.NewService(dbClient, libpbehavior.NewTypeComputer(libpbehavior.NewModelProvider(dbClient), json.NewDecoder()),
		pbhStore, libredis.NewLockClient(pbhRedisSession), logger)
	pbhEntityTypeResolver := libpbehavior.NewEntityTypeResolver(pbhStore, libpbehavior.NewEntityMatcher(dbClient), logger)
	// Create entity service event publisher.
	entityPublChan := make(chan entityservice.ChangeEntityMessage, chanBuf)
	entityServiceEventPublisher := entityservice.NewEventPublisher(
		alarm.NewAdapter(dbClient), libentity.NewAdapter(dbClient), amqpChannel,
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
		entity.NewStore(dbClient, p.TimezoneConfigProvider),
		datastorage.NewAdapter(dbClient),
		metricsEntityMetaUpdater,
		logger,
	)

	userInterfaceAdapter := config.NewUserInterfaceAdapter(dbClient)
	userInterfaceConfig, err := userInterfaceAdapter.GetConfig(ctx)
	if err != nil && err != mongodriver.ErrNoDocuments {
		return nil, nil, err
	}
	if p.UserInterfaceConfigProvider == nil {
		p.UserInterfaceConfigProvider = config.NewUserInterfaceConfigProvider(userInterfaceConfig, logger)
	}

	// Create and compute scenario priority intervals.
	scenarioPriorityIntervals := action.NewPriorityIntervals()
	err = scenarioPriorityIntervals.Recalculate(ctx, dbClient.Collection(mongo.ScenarioMongoCollection))
	if err != nil {
		return nil, nil, err
	}

	// Create csv exporter.
	if exportExecutor == nil {
		exportExecutor = export.NewTaskExecutor(dbClient, logger)
	}

	websocketHub := newWebsocketHub(enforcer, security.GetTokenProvider(), logger)

	broadcastMessageChan := make(chan bool)

	metricsConfigProvider := config.NewMetricsConfigProvider(cfg, logger)
	techPostgresPoolProvider := postgres.NewTechMetricsPoolProvider(ctx, metricsConfigProvider, cfg.Global.ReconnectRetries, cfg.Global.GetReconnectTimeout(), logger)
	techMetricsSender := metrics.NewTechMetricsSender(techPostgresPoolProvider, logger)

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
	legacyUrl := GetLegacyURL(logger)
	api.AddRouter(func(router gin.IRouter) {
		router.Use(middleware.Cache())

		router.Use(func(c *gin.Context) {
			start := time.Now()
			c.Next()

			go techMetricsSender.SendApiRequest(context.Background(), start, c.Request.URL.String(), time.Since(start).Microseconds())
		})

		if flags.Test {
			router.Use(devmiddleware.ReloadEnforcerPolicy(enforcer))
		}

		legacyUrlStr := ""
		if legacyUrl != nil {
			legacyUrlStr = legacyUrl.String()
		}
		RegisterValidators(dbClient, flags.EnableSameServiceNames)
		RegisterRoutes(
			ctx,
			cfg,
			router,
			security,
			enforcer,
			legacyUrlStr,
			dbClient,
			p.TimezoneConfigProvider,
			pbhEntityTypeResolver,
			pbhComputeChan,
			entityPublChan,
			entityCleanerTaskChan,
			engine.NewRunInfoManager(engineRedisSession),
			exportExecutor,
			apilogger.NewActionLogger(dbClient, logger),
			amqpChannel,
			jobQueue,
			p.UserInterfaceConfigProvider,
			scenarioPriorityIntervals,
			cfg.File.Upload,
			websocketHub,
			broadcastMessageChan,
			metricsEntityMetaUpdater,
			metricsUserMetaUpdater,
			logger,
		)
	})
	if flags.EnableDocs {
		api.AddRouter(func(router gin.IRouter) {
			router.GET("/swagger/*filepath", func(c *gin.Context) {
				c.FileFromFS(fmt.Sprintf("swaggerui/%s", c.Param("filepath")), http.FS(docsUiFile))
			})
		})
		if !overrideDocs {
			content, err := docsFile.ReadFile("docs/swagger.yaml")
			if err != nil {
				return nil, nil, err
			}
			schemasContent, err := docsFile.ReadFile("docs/schemas_swagger.yaml")
			if err != nil {
				return nil, nil, err
			}
			api.AddRouter(func(router gin.IRouter) {
				router.GET("/swagger.yaml", docs.GetHandler(schemasContent, content))
			})
		}
	}
	if legacyUrl == nil {
		api.AddNoRoute(func(c *gin.Context) {
			c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		})
	} else {
		api.AddNoRoute(GetProxy(legacyUrl, security, enforcer, proxyAccessConfig)...)
	}
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
			json.NewDecoder(),
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
	api.AddWorker("config reload", updateConfig(p.TimezoneConfigProvider, p.ApiConfigProvider, metricsConfigProvider,
		configAdapter, p.UserInterfaceConfigProvider, userInterfaceAdapter, flags.Test, logger))
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

	return api, docsFile, nil
}

func newWebsocketHub(enforcer libsecurity.Enforcer, tokenProvider libsecurity.TokenProvider, logger zerolog.Logger) websocket.Hub {
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
	if err := websocketHub.RegisterRoom(websocket.RoomBroadcastMessages); err != nil {
		logger.Err(err).Msg("Register BroadcastMessages room")
	}
	if err := websocketHub.RegisterRoom(websocket.RoomLoggedUserCount); err != nil {
		logger.Err(err).Msg("Register LoggedUserCount room")
	}
	return websocketHub
}

func updateConfig(
	timezoneConfigProvider *config.BaseTimezoneConfigProvider,
	apiConfigProvider *config.BaseApiConfigProvider,
	metricsConfigProvider *config.BaseMetricsConfigProvider,
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
				metricsConfigProvider.Update(cfg)

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
