package api

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"time"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	alarmapi "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/broadcastmessage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/contextgraph"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/docs"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/export"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/healthcheck"
	apilogger "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/messageratestats"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/middleware"
	apisecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/security"
	apitechmetrics "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/techmetrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/websocket"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/importcontextgraph"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/link"
	linkv1 "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/link/v1"
	linkv2 "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/link/v2"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/link/wrapper"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	libpbehavior "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/statesetting"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/techmetrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/postgres"
	libredis "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	libsecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	securitymodel "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
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
		p.TemplateConfigProvider = config.NewTemplateConfigProvider(cfg, logger)
	}
	// Set mongodb setting.
	config.SetDbClientRetry(dbClient, cfg)
	// Connect to rmq.
	amqpConn, err := libamqp.NewConnection(logger, -1, cfg.Global.GetReconnectTimeout())
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
	lockRedisSession, err := libredis.NewSession(ctx, libredis.EngineLockStorage, logger,
		cfg.Global.ReconnectRetries, cfg.Global.GetReconnectTimeout())
	if err != nil {
		return nil, nil, fmt.Errorf("cannot connect to redis: %w", err)
	}
	runInfoClient, err := libredis.NewSession(ctx, libredis.EngineRunInfo, logger, -1, -1)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot connect to redis: %w", err)
	}
	securityConfig, err := libsecurity.LoadConfig(flags.ConfigDir)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot load security config: %w", err)
	}

	cookieOptions := DefaultCookieOptions()
	sessionStore := mongostore.NewStore(dbClient, []byte(os.Getenv("SESSION_KEY")))
	sessionStore.Options.MaxAge = cookieOptions.MaxAge
	sessionStore.Options.Secure = flags.SecureSession
	if p.ApiConfigProvider == nil {
		p.ApiConfigProvider = config.NewApiConfigProvider(cfg, logger)
	}
	security := NewSecurity(securityConfig, cfg, dbClient, sessionStore, enforcer, p.ApiConfigProvider, config.NewMaintenanceAdapter(dbClient), cookieOptions, logger)

	if flags.EnableSameServiceNames {
		logger.Info().Msg("Non-unique names for services ENABLED")
	}

	dbExportClient, err := mongo.NewClientWithOptions(ctx, 0, 0, 0,
		p.ApiConfigProvider.Get().ExportMongoClientTimeout, logger)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot connect to mongodb: %w", err)
	}

	proxyAccessConfig, err := proxy.LoadAccessConfig(flags.ConfigDir)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot load access config: %w", err)
	}
	// Create pbehavior computer.
	pbhComputeChan := make(chan []string, chanBuf)
	pbhStore := libpbehavior.NewStore(pbhRedisSession, json.NewEncoder(), json.NewDecoder())
	pbhEntityTypeResolver := libpbehavior.NewEntityTypeResolver(pbhStore, logger)
	// Create entity service event publisher.
	entityPublChan := make(chan entityservice.ChangeEntityMessage, chanBuf)
	entityServiceEventPublisher := entityservice.NewEventPublisher(amqpChannel, json.NewEncoder(),
		canopsis.JsonContentType, canopsis.FIFOAckExchangeName, canopsis.FIFOQueueName, logger)

	importWorker := contextgraph.NewImportWorker(
		cfg,
		contextgraph.NewEventPublisher(canopsis.FIFOExchangeName, canopsis.FIFOQueueName, json.NewEncoder(), canopsis.JsonContentType, amqpChannel),
		contextgraph.NewMongoStatusReporter(dbClient),
		importcontextgraph.NewWorker(
			dbClient,
			importcontextgraph.NewEventPublisher(canopsis.FIFOExchangeName, canopsis.FIFOQueueName, json.NewEncoder(), canopsis.JsonContentType, amqpChannel),
			metricsEntityMetaUpdater,
			logger,
		),
		logger,
	)

	entityCleanerTaskChan := make(chan entity.CleanTask)
	disabledEntityCleaner := entity.NewDisabledCleaner(
		lockRedisSession,
		datastorage.NewAdapter(dbClient),
		p.DataStorageConfigProvider,
		metricsEntityMetaUpdater,
		logger,
	)

	userInterfaceAdapter := config.NewUserInterfaceAdapter(dbClient)
	userInterfaceConfig, err := userInterfaceAdapter.GetConfig(ctx)
	if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
		return nil, nil, fmt.Errorf("cannot load user interface config: %w", err)
	}
	if p.UserInterfaceConfigProvider == nil {
		p.UserInterfaceConfigProvider = config.NewUserInterfaceConfigProvider(userInterfaceConfig, logger)
	}

	// Create csv exporter.
	if exportExecutor == nil {
		exportExecutor = export.NewTaskExecutor(dbClient, p.TimezoneConfigProvider, logger)
	}

	tplExecutor := template.NewExecutor(p.TemplateConfigProvider, p.TimezoneConfigProvider)
	alarmStore := alarmapi.NewStore(dbClient, dbExportClient, linkGenerator, p.TimezoneConfigProvider,
		author.NewProvider(dbClient, p.ApiConfigProvider), tplExecutor, json.NewDecoder(), logger)
	websocketStore := websocket.NewStore(dbClient, flags.IntegrationPeriodicalWaitTime)
	websocketHub, err := newWebsocketHub(enforcer, security.GetTokenProviders(), flags.IntegrationPeriodicalWaitTime,
		dbClient, alarmStore, logger)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot create websocket hub: %w", err)
	}

	broadcastMessageChan := make(chan bool)

	techMetricsConfigProvider := config.NewTechMetricsConfigProvider(cfg, logger)
	techMetricsSender := techmetrics.NewSender(techMetricsConfigProvider, canopsis.TechMetricsFlushInterval,
		cfg.Global.ReconnectRetries, cfg.Global.GetReconnectTimeout(), logger)
	techMetricsTaskExecutor := apitechmetrics.NewTaskExecutor(techMetricsConfigProvider, logger)

	healthCheckConfigAdapter := config.NewHealthCheckAdapter(dbClient)
	healthCheckCfg, err := healthCheckConfigAdapter.GetConfig(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot load healthcheck config: %w", err)
	}

	healthCheckConfigProvider := config.NewBaseHealthCheckConfigProvider(healthCheckCfg, logger)
	healthcheckStore := healthcheck.NewStore(
		dbClient,
		engine.NewRunInfoManager(runInfoClient),
		healthCheckConfigAdapter,
		healthCheckConfigProvider,
		logger,
		websocketHub,
	)

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
			err = dbExportClient.Disconnect(ctx)
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

			err = lockRedisSession.Close()
			if err != nil {
				logger.Error().Err(err).Msg("failed to close redis connection")
			}

			err = runInfoClient.Close()
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
			linkv2.NewGenerator(dbClient, tplExecutor, logger),
		}
		if legacyUrl != nil {
			linkGenerators = append(linkGenerators, linkv1.NewGenerator(legacyUrl.String(), dbClient, &http.Client{Timeout: linkFetchTimeout}, json.NewEncoder(), json.NewDecoder()))
		}
		linkGenerator = wrapper.NewGenerator(linkGenerators...)
	}

	stateSettingsUpdatesChan := make(chan statesetting.RuleUpdatedMessage)

	api.AddRouter(func(router *gin.Engine) {
		router.Use(middleware.CacheControl())

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

		RegisterValidators(dbClient)
		RegisterRoutes(
			ctx,
			cfg,
			router,
			security,
			enforcer,
			linkGenerator,
			dbClient,
			dbExportClient,
			pgPoolProvider,
			amqpChannel,
			p.ApiConfigProvider,
			p.TimezoneConfigProvider,
			p.TemplateConfigProvider,
			pbhEntityTypeResolver,
			pbhComputeChan,
			entityPublChan,
			entityCleanerTaskChan,
			exportExecutor,
			techMetricsTaskExecutor,
			apilogger.NewActionLogger(dbClient, logger),
			amqpChannel,
			p.UserInterfaceConfigProvider,
			websocketHub,
			websocketStore,
			broadcastMessageChan,
			metricsEntityMetaUpdater,
			metricsUserMetaUpdater,
			author.NewProvider(dbClient, p.ApiConfigProvider),
			healthcheckStore,
			tplExecutor,
			stateSettingsUpdatesChan,
			flags.EnableSameServiceNames,
			logger,
		)
	})
	if flags.EnableDocs {
		api.AddRouter(func(router *gin.Engine) {
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
			api.AddRouter(func(router *gin.Engine) {
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

	api.AddWorker("tech_metrics", func(ctx context.Context) {
		techMetricsSender.Run(ctx)
	})
	api.AddWorker("session_clean", func(ctx context.Context) {
		security.GetSessionStore().StartAutoClean(ctx, flags.PeriodicalWaitTime)
	})
	api.AddWorker("enforce_policy_load", func(ctx context.Context) {
		enforcer.StartAutoLoadPolicy(ctx, flags.PeriodicalWaitTime)
	})
	api.AddWorker("pbehavior_compute", sendPbhRecomputeEvents(pbhComputeChan, json.NewEncoder(), amqpChannel, logger))

	stateSettingsListener := statesetting.NewListener(dbClient, amqpChannel, json.NewEncoder(), logger)
	api.AddWorker("state_settings_listener", func(ctx context.Context) {
		stateSettingsListener.Listen(ctx, stateSettingsUpdatesChan)
	})
	api.AddWorker("entity_event_publish", func(ctx context.Context) {
		entityServiceEventPublisher.Publish(ctx, entityPublChan)
	})
	api.AddWorker("entity_cleaner", func(ctx context.Context) {
		disabledEntityCleaner.RunCleanerProcess(ctx, entityCleanerTaskChan)
	})
	api.AddWorker("import_job", func(ctx context.Context) {
		importWorker.Run(ctx)
	})
	api.AddWorker("config_reload", updateConfig(p.TimezoneConfigProvider, p.DataStorageConfigProvider, p.ApiConfigProvider,
		p.TemplateConfigProvider, techMetricsConfigProvider, configAdapter, p.UserInterfaceConfigProvider,
		userInterfaceAdapter, flags.PeriodicalWaitTime, logger))
	api.AddWorker("data_export", func(ctx context.Context) {
		exportExecutor.Execute(ctx)
	})
	api.AddWorker("tech_metrics_export", func(ctx context.Context) {
		techMetricsTaskExecutor.Run(ctx)
	})
	tokenStore := token.NewMongoStore(dbClient, logger)
	shareTokenStore := sharetoken.NewMongoStore(dbClient, logger)
	api.AddWorker("auth_token_activity", updateTokenActivity(flags.IntegrationPeriodicalWaitTime, tokenStore, shareTokenStore,
		websocketHub, logger))
	api.AddWorker("auth_token_expiration", removeExpiredTokens(flags.PeriodicalWaitTime, tokenStore, shareTokenStore,
		logger))
	api.AddWorker("websocket", func(ctx context.Context) {
		websocketHub.Start(ctx)
	})
	api.AddWorker("websocket_conns", updateWebsocketConns(flags.IntegrationPeriodicalWaitTime, websocketHub, websocketStore, logger))

	maintenanceAdapter := config.NewMaintenanceAdapter(dbClient)
	broadcastMessageService := broadcastmessage.NewService(broadcastmessage.NewStore(dbClient, maintenanceAdapter), websocketHub, canopsis.PeriodicalWaitTime, logger)
	api.AddWorker("broadcast_message", func(ctx context.Context) {
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
	api.AddWorker("healthcheck", func(ctx context.Context) {
		healthcheckStore.Load(ctx)
	})
	api.AddWorker("messageratestats", messageratestats.GetWebsocketHandler(websocketHub, pgPoolProvider, logger, flags.IntegrationPeriodicalWaitTime))

	return api, docsFile, nil
}

func newWebsocketHub(
	enforcer libsecurity.Enforcer,
	tokenProviders []libsecurity.TokenProvider,
	checkAuthInterval time.Duration,
	dbClient mongo.DbClient,
	alarmStore alarmapi.Store,
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
	websocketHub := websocket.NewHub(websocketUpgrader, websocketAuthorizer, checkAuthInterval, logger)
	if err := websocketHub.RegisterRoom(websocket.RoomBroadcastMessages); err != nil {
		return nil, fmt.Errorf("fail to register websocket room: %w", err)
	}

	if err := websocketHub.RegisterRoom(websocket.RoomLoggedUserCount); err != nil {
		return nil, fmt.Errorf("fail to register websocket room: %w", err)
	}

	if err := websocketHub.RegisterRoom(websocket.RoomHealthcheck, apisecurity.PermHealthcheck, securitymodel.PermissionCan); err != nil {
		return nil, fmt.Errorf("fail to register websocket room: %w", err)
	}

	if err := websocketHub.RegisterRoom(websocket.RoomHealthcheckStatus, apisecurity.PermHealthcheck, securitymodel.PermissionCan); err != nil {
		return nil, fmt.Errorf("fail to register websocket room: %w", err)
	}

	if err := websocketHub.RegisterRoom(websocket.RoomMessageRates, apisecurity.PermMessageRateStatsRead, securitymodel.PermissionCan); err != nil {
		return nil, fmt.Errorf("fail to register websocket room: %w", err)
	}

	if err := websocketHub.RegisterRoom(websocket.RoomIcons); err != nil {
		return nil, fmt.Errorf("fail to register websocket room: %w", err)
	}

	alarmWatcher := alarmapi.NewWatcher(dbClient, websocketHub, alarmStore, json.NewEncoder(), json.NewDecoder(), logger)
	err := websocketHub.RegisterGroup(websocket.RoomAlarmsGroup, websocket.GroupParameters{
		OnJoin:  alarmWatcher.StartWatch,
		OnLeave: alarmWatcher.StopWatch,
	}, apisecurity.PermAlarmRead, securitymodel.PermissionCan)
	if err != nil {
		return nil, fmt.Errorf("fail to register websocket group: %w", err)
	}

	err = websocketHub.RegisterGroup(websocket.RoomAlarmDetailsGroup, websocket.GroupParameters{
		OnJoin:  alarmWatcher.StartWatchDetails,
		OnLeave: alarmWatcher.StopWatch,
	}, apisecurity.PermAlarmRead, securitymodel.PermissionCan)
	if err != nil {
		return nil, fmt.Errorf("fail to register websocket group: %w", err)
	}

	return websocketHub, nil
}
