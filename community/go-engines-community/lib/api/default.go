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
	apitechmetrics "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/techmetrics"
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
	libcontextgraphV1 "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/importcontextgraph/v1"
	libcontextgraphV2 "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/importcontextgraph/v2"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/link"
	linkv1 "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/link/v1"
	linkv2 "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/link/v2"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/link/wrapper"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	libpbehavior "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/techmetrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/postgres"
	libredis "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	libsecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/proxy"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/session/mongostore"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/sharetoken"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/token"
	"github.com/gin-gonic/gin"
	gorillawebsocket "github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

const chanBuf = 10
const sessionStoreSessionMaxAge = 24 * time.Hour
const linkFetchTimeout = 30 * time.Second

//go:embed swaggerui/*
var docsUiFile embed.FS

//go:embed docs/*.yaml
var docsFile embed.FS

type ConfigProviders struct {
	DataStorageConfigProvider   *config.BaseDataStorageConfigProvider
	TimezoneConfigProvider      *config.BaseTimezoneConfigProvider
	ApiConfigProvider           *config.BaseApiConfigProvider
	TemplateConfigProvider      *config.BaseTemplateConfigProvider
	UserInterfaceConfigProvider *config.BaseUserInterfaceConfigProvider
}

func Default(
	ctx context.Context,
	flags Flags,
	enforcer libsecurity.Enforcer,
	p *ConfigProviders,
	logger zerolog.Logger,
	pgPoolProvider postgres.PoolProvider,
	metricsEntityMetaUpdater metrics.MetaUpdater,
	metricsUserMetaUpdater metrics.MetaUpdater,
	exportExecutor export.TaskExecutor,
	linkGenerator link.Generator,
	deferFunc DeferFunc,
	overrideDocs bool,
) (API, fs.ReadFileFS, error) {
	// Retrieve config.
	dbClient, err := mongo.NewClient(ctx, 0, 0, logger)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot connect to mongodb: %w", err)
	}
	configAdapter := config.NewAdapter(dbClient)
	cfg, err := configAdapter.GetConfig(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot load config: %w", err)
	}
	if p.TimezoneConfigProvider == nil {
		p.TimezoneConfigProvider = config.NewTimezoneConfigProvider(cfg, logger)
	}
	if p.DataStorageConfigProvider == nil {
		p.DataStorageConfigProvider = config.NewDataStorageConfigProvider(cfg, logger)
	}
	if p.TemplateConfigProvider == nil {
		p.TemplateConfigProvider = config.NewTemplateConfigProvider(cfg)
	}
	// Set mongodb setting.
	config.SetDbClientRetry(dbClient, cfg)
	// Connect to rmq.
	amqpConn, err := amqp.NewConnection(logger, -1, cfg.Global.GetReconnectTimeout())
	if err != nil {
		return nil, nil, fmt.Errorf("cannot connect to rmq: %w", err)
	}
	amqpChannel, err := amqpConn.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("cannot connect to rmq: %w", err)
	}
	// Connect to redis.
	pbhRedisSession, err := libredis.NewSession(ctx, libredis.PBehaviorLockStorage, logger,
		cfg.Global.ReconnectRetries, cfg.Global.GetReconnectTimeout())
	if err != nil {
		return nil, nil, fmt.Errorf("cannot connect to redis: %w", err)
	}
	engineRedisSession, err := libredis.NewSession(ctx, libredis.EngineRunInfo, logger,
		cfg.Global.ReconnectRetries, cfg.Global.GetReconnectTimeout())
	if err != nil {
		return nil, nil, fmt.Errorf("cannot connect to redis: %w", err)
	}
	securityConfig, err := libsecurity.LoadConfig(flags.ConfigDir)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot load security config: %w", err)
	}

	cookieOptions := CookieOptions{
		FileAccessName: "token",
		MaxAge:         int(sessionStoreSessionMaxAge.Seconds()),
	}
	sessionStore := mongostore.NewStore(dbClient, []byte(os.Getenv("SESSION_KEY")))
	sessionStore.Options.MaxAge = cookieOptions.MaxAge
	sessionStore.Options.Secure = flags.SecureSession
	if p.ApiConfigProvider == nil {
		p.ApiConfigProvider = config.NewApiConfigProvider(cfg, logger)
	}
	security := NewSecurity(securityConfig, dbClient, sessionStore, enforcer, p.ApiConfigProvider, cookieOptions, logger)

	if flags.EnableSameServiceNames {
		logger.Info().Msg("Non-unique names for services ENABLED")
	}

	proxyAccessConfig, err := proxy.LoadAccessConfig(flags.ConfigDir)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot load access config: %w", err)
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
		libcontextgraphV1.NewWorker(
			dbClient,
			importcontextgraph.NewEventPublisher(canopsis.FIFOExchangeName, canopsis.FIFOQueueName, json.NewEncoder(), canopsis.JsonContentType, amqpChannel),
			metricsEntityMetaUpdater,
		),
		libcontextgraphV2.NewWorker(
			dbClient,
			importcontextgraph.NewEventPublisher(canopsis.FIFOExchangeName, canopsis.FIFOQueueName, json.NewEncoder(), canopsis.JsonContentType, amqpChannel),
			metricsEntityMetaUpdater,
		),
		logger,
	)

	entityCleanerTaskChan := make(chan entity.CleanTask)
	disabledEntityCleaner := entity.NewDisabledCleaner(
		datastorage.NewAdapter(dbClient),
		p.DataStorageConfigProvider,
		metricsEntityMetaUpdater,
		logger,
	)

	userInterfaceAdapter := config.NewUserInterfaceAdapter(dbClient)
	userInterfaceConfig, err := userInterfaceAdapter.GetConfig(ctx)
	if err != nil && err != mongodriver.ErrNoDocuments {
		return nil, nil, fmt.Errorf("cannot load user interface config: %w", err)
	}
	if p.UserInterfaceConfigProvider == nil {
		p.UserInterfaceConfigProvider = config.NewUserInterfaceConfigProvider(userInterfaceConfig, logger)
	}

	// Create and compute scenario priority intervals.
	scenarioPriorityIntervals := action.NewPriorityIntervals()
	err = scenarioPriorityIntervals.Recalculate(ctx, dbClient.Collection(mongo.ScenarioMongoCollection))
	if err != nil {
		return nil, nil, fmt.Errorf("cannot recalculate scenario preority: %w", err)
	}

	// Create csv exporter.
	if exportExecutor == nil {
		exportExecutor = export.NewTaskExecutor(dbClient, p.TimezoneConfigProvider, logger)
	}

	websocketHub, err := newWebsocketHub(enforcer, security.GetTokenProviders(), flags.IntegrationPeriodicalWaitTime, logger)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot create websocket hub: %w", err)
	}

	broadcastMessageChan := make(chan bool)

	techMetricsConfigProvider := config.NewTechMetricsConfigProvider(cfg, logger)
	techMetricsSender := techmetrics.NewSender(techMetricsConfigProvider, canopsis.TechMetricsFlushInterval,
		cfg.Global.ReconnectRetries, cfg.Global.GetReconnectTimeout(), logger)
	techMetricsTaskExecutor := apitechmetrics.NewTaskExecutor(techMetricsConfigProvider, logger)

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
	if linkGenerator == nil {
		linkGenerators := []link.Generator{
			linkv2.NewGenerator(dbClient, template.NewExecutor(p.TemplateConfigProvider, p.TimezoneConfigProvider), logger),
		}
		if legacyUrl != nil {
			linkGenerators = append(linkGenerators, linkv1.NewGenerator(legacyUrl.String(), dbClient, &http.Client{Timeout: linkFetchTimeout}, json.NewEncoder(), json.NewDecoder()))
		}
		linkGenerator = wrapper.NewGenerator(linkGenerators...)
	}

	api.AddRouter(func(router gin.IRouter) {
		router.Use(middleware.Cache())

		router.Use(func(c *gin.Context) {
			start := time.Now()
			c.Next()

			techMetricsSender.SendApiRequest(techmetrics.ApiRequestMetric{
				Timestamp: start,
				Interval:  time.Since(start),
				Method:    c.Request.Method,
				Url:       c.Request.URL.String(),
			})
		})

		RegisterValidators(dbClient, flags.EnableSameServiceNames)
		RegisterRoutes(
			ctx,
			cfg,
			router,
			security,
			enforcer,
			linkGenerator,
			dbClient,
			pgPoolProvider,
			p.TimezoneConfigProvider,
			p.TemplateConfigProvider,
			pbhEntityTypeResolver,
			pbhComputeChan,
			entityPublChan,
			entityCleanerTaskChan,
			engine.NewRunInfoManager(engineRedisSession),
			exportExecutor,
			techMetricsTaskExecutor,
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
				return nil, nil, fmt.Errorf("cannot read swagger: %w", err)
			}
			schemasContent, err := docsFile.ReadFile("docs/schemas_swagger.yaml")
			if err != nil {
				return nil, nil, fmt.Errorf("cannot read swagger: %w", err)
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

	api.AddWorker("tech metrics", func(ctx context.Context) {
		techMetricsSender.Run(ctx)
	})
	api.AddWorker("session clean", func(ctx context.Context) {
		security.GetSessionStore().StartAutoClean(ctx, flags.IntegrationPeriodicalWaitTime)
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
	api.AddWorker("config reload", updateConfig(p.TimezoneConfigProvider, p.DataStorageConfigProvider, p.ApiConfigProvider,
		p.TemplateConfigProvider, techMetricsConfigProvider, configAdapter, p.UserInterfaceConfigProvider,
		userInterfaceAdapter, flags.PeriodicalWaitTime, logger))
	api.AddWorker("data export", func(ctx context.Context) {
		exportExecutor.Execute(ctx)
	})
	api.AddWorker("tech metrics export", func(ctx context.Context) {
		techMetricsTaskExecutor.Run(ctx)
	})
	api.AddWorker("auth token activity", func(ctx context.Context) {
		ticker := time.NewTicker(canopsis.PeriodicalWaitTime)
		defer ticker.Stop()
		tokenStore := token.NewMongoStore(dbClient, logger)
		shareTokenStore := sharetoken.NewMongoStore(dbClient, logger)

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				for _, tokens := range websocketHub.GetUsers() {
					for _, t := range tokens {
						err := tokenStore.Access(ctx, t)
						if err != nil {
							logger.Err(err).Msg("cannot update token access")
						}
						err = shareTokenStore.Access(ctx, t)
						if err != nil {
							logger.Err(err).Msg("cannot update share token access")
						}
					}
				}
			}
		}
	})
	api.AddWorker("auth token expiration", func(ctx context.Context) {
		ticker := time.NewTicker(canopsis.PeriodicalWaitTime)
		defer ticker.Stop()
		tokenStore := token.NewMongoStore(dbClient, logger)
		shareTokenStore := sharetoken.NewMongoStore(dbClient, logger)

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				err := tokenStore.DeleteExpired(ctx)
				if err != nil {
					logger.Err(err).Msg("cannot delete expired tokens")
				}
				err = shareTokenStore.DeleteExpired(ctx)
				if err != nil {
					logger.Err(err).Msg("cannot delete expired share tokens")
				}
			}
		}
	})
	api.AddWorker("websocket", func(ctx context.Context) {
		websocketHub.Start(ctx)
	})
	broadcastMessageService := broadcastmessage.NewService(broadcastmessage.NewStore(dbClient), websocketHub, canopsis.PeriodicalWaitTime, logger)
	api.AddWorker("broadcast message", func(ctx context.Context) {
		broadcastMessageService.Start(ctx, broadcastMessageChan)
	})
	api.AddWorker("links", func(ctx context.Context) {
		ticker := time.NewTicker(flags.PeriodicalWaitTime)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				err := linkGenerator.Load(ctx)
				if err != nil {
					logger.Err(err).Msg("cannot load links")
				}
			}
		}
	})

	return api, docsFile, nil
}

func newWebsocketHub(
	enforcer libsecurity.Enforcer,
	tokenProviders []libsecurity.TokenProvider,
	checkAuthInterval time.Duration,
	logger zerolog.Logger,
) (websocket.Hub, error) {
	websocketUpgrader := websocket.NewUpgrader(gorillawebsocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 2048,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	})
	websocketAuthorizer := websocket.NewAuthorizer(enforcer, tokenProviders)
	websocketHub := websocket.NewHub(websocketUpgrader, websocketAuthorizer,
		checkAuthInterval, logger)
	if err := websocketHub.RegisterRoom(websocket.RoomBroadcastMessages); err != nil {
		return nil, err
	}
	if err := websocketHub.RegisterRoom(websocket.RoomLoggedUserCount); err != nil {
		return nil, err
	}
	return websocketHub, nil
}

func updateConfig(
	timezoneConfigProvider *config.BaseTimezoneConfigProvider,
	dataStorageConfigProvider *config.BaseDataStorageConfigProvider,
	apiConfigProvider *config.BaseApiConfigProvider,
	templateConfigProvider *config.BaseTemplateConfigProvider,
	techMetricsConfigProvider *config.BaseTechMetricsConfigProvider,
	configAdapter config.Adapter,
	userInterfaceConfigProvider *config.BaseUserInterfaceConfigProvider,
	userInterfaceAdapter config.UserInterfaceAdapter,
	interval time.Duration,
	logger zerolog.Logger,
) func(ctx context.Context) {
	return func(ctx context.Context) {
		ticker := time.NewTicker(interval)
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
				techMetricsConfigProvider.Update(cfg)
				dataStorageConfigProvider.Update(cfg)
				templateConfigProvider.Update(cfg)

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
