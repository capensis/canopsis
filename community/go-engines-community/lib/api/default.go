package api

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"time"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	alarmapi "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/broadcastmessage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/contextgraph"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/docs"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/export"
	apiexternaldata "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/externaldatatable"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/healthcheck"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	apilogger "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/messageratestats"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/middleware"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/notification"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/patternfields"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pbehavior"
	apisecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/security"
	apitechmetrics "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/techmetrics"
	libcommtemplate "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/websocket"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/workers"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/wsconn"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	libentity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/externaldata"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/importcontextgraph"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/link"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	libpbehavior "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/statesetting"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/techmetrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/usernotification"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/postgres"
	libredis "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	libsecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	securitymodel "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/session/mongostore"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/sharetoken"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/token"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/gin-gonic/gin"
	gorillawebsocket "github.com/gorilla/websocket"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
)

const (
	chanBuf = 10

	websocketReadBufferSize  = 1024
	websocketWriteBufferSize = 2048

	jobKeyExport          = "export"
	jobKeyImport          = "import"
	jobKeyExtDataImport   = "extdataimport"
	jobKeyPbhPatterns     = "pbhpatterns"
	jobKeyPatternOptimize = "patternoptimize"
)

//go:embed swaggerui/*
var docsUiFile embed.FS

//go:embed docs/*.yaml
var docsFile embed.FS

type Services struct {
	ExportTaskExecutor export.TaskExecutor
	LinkGenerator      link.Generator
	Enforcer           libsecurity.Enforcer
	DocsFile           fs.ReadFileFS

	WebsocketHub                websocket.Hub
	WebsocketRoomRegistry       websocket.RoomRegistry
	WebsocketAuthorize          func(perms ...string) websocket.Authorize
	DataStorageConfigProvider   *config.BaseDataStorageConfigProvider
	TimezoneConfigProvider      *config.BaseTimezoneConfigProvider
	ApiConfigProvider           *config.BaseApiConfigProvider
	TemplateConfigProvider      *config.BaseTemplateConfigProvider
	UserInterfaceConfigProvider *config.BaseUserInterfaceConfigProvider
	ExternalDataContainer       *externaldata.GetterContainer
	NotificationStore           usernotification.Store
	ErrorResponder              httperror.Responder
}

func Default(
	ctx context.Context,
	flags Flags,
	logger zerolog.Logger,
	pgPoolProvider postgres.PoolProvider,
	metricsEntityMetaUpdater metrics.MetaUpdater,
	metricsUserMetaUpdater metrics.MetaUpdater,
	tplTestTypePermMapping map[int][]any,
	deferFunc DeferFunc,
	overrideDocs bool,
) (API, Services, error) {
	services := Services{
		DocsFile: docsFile,
	}

	// Retrieve config.
	primaryDbClient, err := mongo.NewClient(ctx)
	if err != nil {
		return nil, services, fmt.Errorf("cannot create primary mongodb client: %w", err)
	}

	configAdapter := config.NewAdapter(primaryDbClient)
	cfg, err := configAdapter.GetConfig(ctx)
	if err != nil {
		return nil, services, fmt.Errorf("cannot load config: %w", err)
	}

	// Set mongodb setting.
	config.SetDbClientRetry(primaryDbClient, cfg)

	secondaryDbClient, err := mongo.NewClient(ctx, mongo.ClientOptions{
		RetryCount:      cfg.Global.ReconnectRetries,
		MinRetryTimeout: cfg.Global.GetReconnectTimeout(),
		ReadPreference:  mongo.SecondaryPreferred(),
	})
	if err != nil {
		return nil, services, fmt.Errorf("cannot create secondary mongodb client: %w", err)
	}

	// noTimeoutClient should be used by change stream watchers only.
	noTimeoutClient, err := mongo.NewClient(ctx, mongo.ClientOptions{
		RetryCount:      cfg.Global.ReconnectRetries,
		MinRetryTimeout: cfg.Global.GetReconnectTimeout(),
		ReadPreference:  mongo.SecondaryPreferred(),
		NoClientTimeout: true,
	})
	if err != nil {
		return nil, services, fmt.Errorf("cannot create mongodb client without timeout: %w", err)
	}

	services.Enforcer, err = libsecurity.NewEnforcer(flags.ConfigDir, primaryDbClient)
	if err != nil {
		return nil, services, fmt.Errorf("cannot create security enforcer: %w", err)
	}

	services.TimezoneConfigProvider = config.NewTimezoneConfigProvider(cfg, logger)
	services.DataStorageConfigProvider = config.NewDataStorageConfigProvider(cfg, logger)
	services.TemplateConfigProvider = config.NewTemplateConfigProvider(cfg, logger)
	// Connect to rmq.
	amqpConn, err := libamqp.NewConnection(logger, -1, cfg.Global.GetReconnectTimeout())
	if err != nil {
		return nil, services, fmt.Errorf("cannot connect to rmq: %w", err)
	}
	amqpPublisher, err := amqpConn.Channel()
	if err != nil {
		return nil, services, fmt.Errorf("cannot connect to rmq: %w", err)
	}
	// Connect to redis.
	pbhRedisSession, err := libredis.NewSession(ctx, libredis.PBehaviorLockStorage, logger,
		cfg.Global.ReconnectRetries, cfg.Global.GetReconnectTimeout())
	if err != nil {
		return nil, services, fmt.Errorf("cannot connect to redis: %w", err)
	}
	lockRedisSession, err := libredis.NewSession(ctx, libredis.EngineLockStorage, logger,
		cfg.Global.ReconnectRetries, cfg.Global.GetReconnectTimeout())
	if err != nil {
		return nil, services, fmt.Errorf("cannot connect to redis: %w", err)
	}
	runInfoClient, err := libredis.NewSession(ctx, libredis.EngineRunInfo, logger, -1, -1)
	if err != nil {
		return nil, services, fmt.Errorf("cannot connect to redis: %w", err)
	}
	securityConfig, err := libsecurity.LoadConfig(flags.ConfigDir)
	if err != nil {
		return nil, services, fmt.Errorf("cannot load security config: %w", err)
	}

	cookieOptions := DefaultCookieOptions()
	sessionStore := mongostore.NewStore(primaryDbClient, GetSessionKeyVar(logger))
	sessionStore.Options.MaxAge = cookieOptions.MaxAge
	sessionStore.Options.Secure = flags.SecureSession
	services.ApiConfigProvider = config.NewApiConfigProvider(cfg, logger)
	tplExecutor := template.NewExecutor(services.TemplateConfigProvider, services.TimezoneConfigProvider)
	uniTrans, err := RegisterValidators(securityConfig, tplExecutor)
	if err != nil {
		return nil, services, fmt.Errorf("cannot register request validators: %w", err)
	}

	services.ErrorResponder = httperror.NewResponder(validation.NewErrorTranslator(uniTrans, logger), logger)
	security := NewSecurity(securityConfig, cfg, primaryDbClient, sessionStore, services.Enforcer,
		services.ApiConfigProvider, config.NewMaintenanceAdapter(primaryDbClient), cookieOptions,
		services.ErrorResponder, logger)

	if flags.EnableSameServiceNames {
		logger.Info().Msg("Non-unique names for services ENABLED")
	}

	dbExportClient, err := mongo.NewClient(ctx, mongo.ClientOptions{
		ClientTimeout:  services.ApiConfigProvider.Get().ExportMongoClientTimeout,
		ReadPreference: mongo.SecondaryPreferred(),
	})
	if err != nil {
		return nil, services, fmt.Errorf("cannot connect to mongodb: %w", err)
	}

	// Create pbehavior computer.
	pbhComputeChan := make(chan rpc.PbehaviorRecomputeEvent, chanBuf)
	pbhStore := libpbehavior.NewStore(pbhRedisSession, json.NewEncoder(), json.NewDecoder())
	pbhEntityTypeResolver := libpbehavior.NewEntityTypeResolver(pbhStore, logger)
	// Create entity service event publisher.
	entityPublChan := make(chan entityservice.ChangeEntityMessage, chanBuf)
	entityServiceEventPublisher := entityservice.NewEventPublisher(amqpPublisher, json.NewEncoder(),
		canopsis.JsonContentType, canopsis.DefaultExchangeName, canopsis.FIFOQueueName, canopsis.ApiConnector, logger)

	entityCleanerTaskChan := make(chan libentity.CleanTask)
	disabledEntityCleaner := libentity.NewCleaner(
		libredis.NewLockClient(lockRedisSession),
		datastorage.NewAdapter(primaryDbClient),
		services.DataStorageConfigProvider,
		metricsEntityMetaUpdater,
		logger,
	)

	userInterfaceAdapter := config.NewUserInterfaceAdapter(primaryDbClient)
	userInterfaceConfig, err := userInterfaceAdapter.GetConfig(ctx)
	if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
		return nil, services, fmt.Errorf("cannot load user interface config: %w", err)
	}

	services.UserInterfaceConfigProvider = config.NewUserInterfaceConfigProvider(userInterfaceConfig, logger)
	workersRunner := workers.NewRunner(amqpConn, logger)
	// Create csv exporter.
	services.ExportTaskExecutor = export.NewTaskExecutor(primaryDbClient, workers.NewJobPublisher(jobKeyExport, amqpPublisher),
		services.TimezoneConfigProvider, filepath.Join(cfg.File.Dir, canopsis.SubDirExport), logger)
	workersRunner.AddJobExecutor(jobKeyExport, func(ctx context.Context, id string) error {
		return services.ExportTaskExecutor.ExecuteTask(ctx, id)
	})
	importWorker := contextgraph.NewImportWorker(
		cfg,
		contextgraph.NewEventPublisher(canopsis.DefaultExchangeName, canopsis.FIFOQueueName, json.NewEncoder(), canopsis.JsonContentType, amqpPublisher),
		contextgraph.NewMongoStatusReporter(primaryDbClient),
		importcontextgraph.NewWorker(
			primaryDbClient,
			importcontextgraph.NewEventPublisher(canopsis.DefaultExchangeName, canopsis.FIFOQueueName, json.NewEncoder(), canopsis.JsonContentType, amqpPublisher),
			metricsEntityMetaUpdater,
			canopsis.ApiConnector,
			patternfields.NewTransformer(primaryDbClient),
			logger,
		),
		workers.NewJobPublisher(jobKeyImport, amqpPublisher),
		logger,
	)
	workersRunner.AddJobExecutor(jobKeyImport, func(ctx context.Context, _ string) error {
		return importWorker.ProcessFirstJob(ctx)
	})
	websocketStore := wsconn.NewStore(primaryDbClient, flags.IntegrationPeriodicalWaitTime)
	websocketUpgrader := websocket.NewUpgrader(&gorillawebsocket.Upgrader{
		ReadBufferSize:  websocketReadBufferSize,
		WriteBufferSize: websocketWriteBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	})
	services.WebsocketRoomRegistry = websocket.NewRoomRegistry()
	wsRoomAuthenticate := func(ctx context.Context, token string) (string, error) {
		tokenProviders := security.GetTokenProviders()
		for _, provider := range tokenProviders {
			user, err := provider.Auth(ctx, token)
			if err != nil {
				return "", err
			}

			if user != nil {
				return user.ID, nil
			}
		}

		return "", errors.New("invalid token")
	}
	services.WebsocketHub = websocket.NewHub(websocketUpgrader, services.WebsocketRoomRegistry, wsRoomAuthenticate, services.ApiConfigProvider, flags.IntegrationPeriodicalWaitTime, logger)
	services.ExternalDataContainer = externaldata.NewGetterContainer()
	services.LinkGenerator = link.NewGenerator(primaryDbClient, tplExecutor, services.ExternalDataContainer, logger)
	authorProvider := author.NewProvider(services.ApiConfigProvider)
	alarmStore := alarmapi.NewStore(secondaryDbClient, dbExportClient, services.LinkGenerator, patternfields.NewTransformer(primaryDbClient),
		services.TimezoneConfigProvider, authorProvider, tplExecutor, json.NewDecoder(), logger)
	alarmWatcher := alarmapi.NewWatcher(noTimeoutClient, services.WebsocketHub, alarmStore, json.NewEncoder(), json.NewDecoder(), logger)

	messageRateWatcher := messageratestats.NewWatcher(services.WebsocketHub, messageratestats.NewStore(pgPoolProvider),
		json.NewEncoder(), json.NewDecoder(), flags.IntegrationPeriodicalWaitTime, logger)

	services.WebsocketAuthorize = func(perms ...string) websocket.Authorize {
		return func(ctx context.Context, userID string) (bool, error) {
			vals := make([]any, len(perms)+1)
			vals[0] = userID
			for i, v := range perms {
				vals[i+1] = v
			}

			return services.Enforcer.Enforce(vals...)
		}
	}
	err = registerWebsocketRooms(services.WebsocketRoomRegistry, services.WebsocketAuthorize, alarmWatcher, messageRateWatcher)
	if err != nil {
		return nil, services, fmt.Errorf("cannot register websocket rooms: %w", err)
	}

	broadcastMessageChan := make(chan bool, chanBuf)

	techMetricsConfigProvider := config.NewTechMetricsConfigProvider(cfg, logger)
	techMetricsSender := techmetrics.NewSender(canopsis.ApiName+"/"+utils.NewID(), techMetricsConfigProvider, canopsis.TechMetricsFlushInterval,
		cfg.Global.ReconnectRetries, cfg.Global.GetReconnectTimeout(), logger)
	techMetricsTaskExecutor := apitechmetrics.NewTaskExecutor(apitechmetrics.NewStore(primaryDbClient), filepath.Join(cfg.File.Dir, canopsis.SubDirExport), logger)

	healthCheckConfigAdapter := config.NewHealthCheckAdapter(primaryDbClient)
	healthCheckCfg, err := healthCheckConfigAdapter.GetConfig(ctx)
	if err != nil {
		return nil, services, fmt.Errorf("cannot load healthcheck config: %w", err)
	}

	healthCheckConfigProvider := config.NewBaseHealthCheckConfigProvider(healthCheckCfg, logger)
	healthcheckStore := healthcheck.NewStore(
		primaryDbClient,
		engine.NewRunInfoManager(runInfoClient),
		healthCheckConfigAdapter,
		healthCheckConfigProvider,
		logger,
		services.WebsocketHub,
	)

	exdataImportWorker := apiexternaldata.NewImportWorker(primaryDbClient, pgPoolProvider,
		filepath.Join(cfg.File.Dir, canopsis.SubDirExDataImport), workers.NewJobPublisher(jobKeyExtDataImport, amqpPublisher),
		logger)
	workersRunner.AddJobExecutor(jobKeyExtDataImport, func(ctx context.Context, id string) error {
		return exdataImportWorker.ProcessJob(ctx, id)
	})

	stateSettingsUpdatesChan := make(chan statesetting.RuleUpdatedMessage)

	patternStore := pattern.NewStore(primaryDbClient, secondaryDbClient, pbhComputeChan, entityPublChan, stateSettingsUpdatesChan, authorProvider, patternfields.NewTransformer(primaryDbClient), logger)
	patternOptimizeWorker := pattern.NewOptimizeWorker(patternStore, primaryDbClient, workers.NewJobPublisher(jobKeyPatternOptimize, amqpPublisher), patternfields.NewTransformer(primaryDbClient), logger)
	workersRunner.AddJobExecutor(jobKeyPatternOptimize, func(ctx context.Context, id string) error {
		return patternOptimizeWorker.ProcessJob(ctx, id)
	})
	apiPbhStore := pbehavior.NewStore(primaryDbClient, secondaryDbClient, lockRedisSession, pbhEntityTypeResolver,
		libpbehavior.NewTypeComputer(libpbehavior.NewModelProvider(primaryDbClient, authorProvider), json.NewDecoder()),
		services.TimezoneConfigProvider, authorProvider, patternfields.NewTransformer(primaryDbClient),
		services.WebsocketHub, services.UserInterfaceConfigProvider)
	workersRunner.AddJobExecutor(jobKeyPbhPatterns, func(ctx context.Context, _ string) error {
		return apiPbhStore.ExecPatternsAndUpdate(ctx)
	})

	services.NotificationStore = usernotification.NewStore(primaryDbClient, amqpPublisher, json.NewEncoder(),
		canopsis.ApiNotificationExchangeName, "", canopsis.JsonContentType)
	notifQueueListener := notification.NewQueueListener(primaryDbClient, amqpConn, services.WebsocketHub,
		notification.NewStore(primaryDbClient, authorProvider), json.NewDecoder(), services.ApiConfigProvider, logger)

	if tplTestTypePermMapping == nil {
		tplTestTypePermMapping = make(map[int][]any)
	}

	tplTestTypePermMapping[libcommtemplate.TypeTestEventFilterRule] = []any{
		apisecurity.ObjEventFilterRule,
		securitymodel.PermissionCreate,
	}
	tplTestTypePermMapping[libcommtemplate.TypeTestLinkRule] = []any{
		apisecurity.ObjLinkRule,
		securitymodel.PermissionCreate,
	}
	tplTestTypePermMapping[libcommtemplate.TypeTestActionScenario] = []any{
		apisecurity.ObjAction,
		securitymodel.PermissionCreate,
	}
	tplTestTypePermMapping[libcommtemplate.TypeTestWidget] = []any{
		apisecurity.ObjView,
		securitymodel.PermissionCreate,
	}

	// Create api.
	api := New(
		fmt.Sprintf(":%d", flags.Port),
		func(ctx context.Context) {
			close(pbhComputeChan)
			close(entityPublChan)
			close(entityCleanerTaskChan)
			close(broadcastMessageChan)

			err := primaryDbClient.Disconnect(context.WithoutCancel(ctx))
			if err != nil {
				logger.Error().Err(err).Msg("failed to close primary mongo connection")
			}

			err = secondaryDbClient.Disconnect(context.WithoutCancel(ctx))
			if err != nil {
				logger.Error().Err(err).Msg("failed to close secondary mongo connection")
			}

			err = noTimeoutClient.Disconnect(context.WithoutCancel(ctx))
			if err != nil {
				logger.Error().Err(err).Msg("failed to close mongo connection without timeout")
			}

			err = dbExportClient.Disconnect(context.WithoutCancel(ctx))
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

	api.AddRouter(func(router *gin.Engine) error {
		router.Use(middleware.Logger(services.ErrorResponder, logger, flags.LogBody, flags.LogBodyOnError))
		router.Use(middleware.Recovery(logger))
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

		return RegisterRoutes(
			ctx,
			cfg,
			router,
			security,
			services.Enforcer,
			services.ErrorResponder,
			services.LinkGenerator,
			primaryDbClient,
			secondaryDbClient,
			dbExportClient,
			pgPoolProvider,
			amqpPublisher,
			lockRedisSession,
			services.ApiConfigProvider,
			services.TimezoneConfigProvider,
			services.TemplateConfigProvider,
			pbhEntityTypeResolver,
			pbhComputeChan,
			entityPublChan,
			entityCleanerTaskChan,
			services.ExportTaskExecutor,
			techMetricsTaskExecutor,
			services.UserInterfaceConfigProvider,
			services.WebsocketHub,
			websocketStore,
			broadcastMessageChan,
			metricsEntityMetaUpdater,
			metricsUserMetaUpdater,
			authorProvider,
			healthcheckStore,
			tplExecutor,
			stateSettingsUpdatesChan,
			flags.EnableSameServiceNames,
			event.NewGenerator(canopsis.ApiConnector, canopsis.ApiConnector),
			securityConfig,
			exdataImportWorker,
			patternOptimizeWorker,
			services.NotificationStore,
			services.ExternalDataContainer,
			tplTestTypePermMapping,
			logger,
		)
	})
	if flags.EnableDocs {
		api.AddRouter(func(router *gin.Engine) error {
			router.GET("/swagger/*filepath", func(c *gin.Context) {
				c.FileFromFS("swaggerui/"+c.Param("filepath"), http.FS(docsUiFile))
			})

			return nil
		})
		if !overrideDocs {
			content, err := docsFile.ReadFile("docs/swagger.yaml")
			if err != nil {
				return nil, services, fmt.Errorf("cannot read swagger: %w", err)
			}
			schemasContent, err := docsFile.ReadFile("docs/schemas_swagger.yaml")
			if err != nil {
				return nil, services, fmt.Errorf("cannot read swagger: %w", err)
			}
			api.AddRouter(func(router *gin.Engine) error {
				router.GET("/swagger.yaml", docs.GetHandler(
					services.ErrorResponder,
					[][]byte{schemasContent},
					[][]byte{content},
				))

				return nil
			})
		}
	}

	api.AddNoRoute(func(c *gin.Context) {
		services.ErrorResponder.Respond(c, httperror.ErrNotFound)
	})
	api.AddNoMethod(func(c *gin.Context) {
		services.ErrorResponder.Respond(c, httperror.ErrMethodNotAllowed)
	})

	actionLogger := apilogger.NewActionLogger(noTimeoutClient, libredis.NewLockClient(lockRedisSession), pgPoolProvider, logger, cfg.Global.ReconnectRetries, cfg.Global.GetReconnectTimeout())
	api.AddWorker("action_log", func(ctx context.Context) {
		err := actionLogger.Watch(ctx)
		if err != nil && !errors.Is(err, context.Canceled) {
			panic(FatalWorkerError{err: err})
		}
	})

	api.AddWorker("amqp_workers", func(ctx context.Context) {
		err = workersRunner.Run(ctx)
		if err != nil {
			var amqpErr *amqp.Error
			if errors.As(err, &amqpErr) && amqpErr.Code == amqp.NotFound {
				panic(NewFatalWorkerError(err))
			}

			panic(err)
		}
	})

	api.AddWorker("tech_metrics", func(ctx context.Context) {
		techMetricsSender.Run(ctx)
	})
	api.AddWorker("session_clean", func(ctx context.Context) {
		security.GetSessionStore().StartAutoClean(ctx, flags.PeriodicalWaitTime)
	})
	api.AddWorker("enforce_policy_load", func(ctx context.Context) {
		services.Enforcer.StartAutoLoadPolicy(ctx, flags.PeriodicalWaitTime)
	})
	api.AddWorker("pbehavior_compute", sendPbhRecomputeEvents(pbhComputeChan, json.NewEncoder(), amqpPublisher, logger))

	stateSettingsListener := statesetting.NewListener(primaryDbClient, amqpPublisher, canopsis.ApiConnector,
		flags.IntegrationPeriodicalWaitTime, flags.StateSettingRecomputeDelay, json.NewEncoder(), logger)
	api.AddWorker("state_settings_listener", func(ctx context.Context) {
		stateSettingsListener.Listen(ctx, stateSettingsUpdatesChan)
	})
	api.AddWorker("state_settings_worker", func(ctx context.Context) {
		stateSettingsListener.Work(ctx)
	})
	api.AddWorker("entity_event_publish", func(ctx context.Context) {
		entityServiceEventPublisher.Publish(ctx, entityPublChan)
	})
	api.AddWorker("entity_cleaner", func(ctx context.Context) {
		disabledEntityCleaner.RunCleanerProcess(ctx, entityCleanerTaskChan)
	})
	api.AddWorker("data_import_abandoned", func(ctx context.Context) {
		importWorker.ProcessAbandonedJob(ctx)
	})
	api.AddWorker("data_import_delete_old", func(ctx context.Context) {
		importWorker.DeleteOldJobs(ctx)
	})
	api.AddWorker("config_reload", updateConfig(services.TimezoneConfigProvider, services.DataStorageConfigProvider,
		services.ApiConfigProvider, services.TemplateConfigProvider, techMetricsConfigProvider, configAdapter,
		services.UserInterfaceConfigProvider, userInterfaceAdapter, flags.PeriodicalWaitTime, logger))
	api.AddWorker("data_export_abandoned", func(ctx context.Context) {
		services.ExportTaskExecutor.ProcessAbandonedTasks(ctx)
	})
	api.AddWorker("data_export_delete_old", func(ctx context.Context) {
		services.ExportTaskExecutor.DeleteOldTasks(ctx)
	})
	api.AddWorker("tech_metrics_export", func(ctx context.Context) {
		techMetricsTaskExecutor.Run(ctx)
	})
	tokenStore := token.NewMongoStore(primaryDbClient, logger)
	shareTokenStore := sharetoken.NewMongoStore(primaryDbClient, logger)
	api.AddWorker("auth_token_activity", updateTokenActivity(flags.IntegrationPeriodicalWaitTime, tokenStore, shareTokenStore,
		services.WebsocketHub, logger))
	api.AddWorker("auth_token_expiration", removeExpiredTokens(flags.PeriodicalWaitTime, tokenStore, shareTokenStore,
		logger))
	api.AddWorker("websocket", func(ctx context.Context) {
		services.WebsocketHub.Run(ctx)
	})
	api.AddWorker("websocket_conns", updateWebsocketConns(flags.IntegrationPeriodicalWaitTime, services.WebsocketHub, websocketStore, logger))

	maintenanceAdapter := config.NewMaintenanceAdapter(primaryDbClient)
	broadcastMessageService := broadcastmessage.NewService(
		broadcastmessage.NewStore(primaryDbClient, maintenanceAdapter, authorProvider), services.WebsocketHub,
		flags.BroadcastMessagePeriodicalWaitTime, logger)
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
				err := services.LinkGenerator.Load(ctx)
				if err != nil {
					logger.Err(err).Msg("cannot load links")
				}
			}
		}
	})
	api.AddWorker("data_exdata_import_abandoned", func(ctx context.Context) {
		exdataImportWorker.ProcessAbandonedJobs(ctx)
	})
	api.AddWorker("data_extdata_import_delete_old", func(ctx context.Context) {
		exdataImportWorker.DeleteOldJobs(ctx)
	})
	api.AddWorker("pattern_optimize_abandoned", func(ctx context.Context) {
		patternOptimizeWorker.ProcessAbandonedJobs(ctx)
	})
	api.AddWorker("healthcheck", func(ctx context.Context) {
		healthcheckStore.Load(ctx)
	})
	api.AddWorker("notification_queue_listen", func(ctx context.Context) {
		err := notifQueueListener.Listen(ctx)
		if err != nil {
			var amqpErr *amqp.Error
			if errors.As(err, &amqpErr) && amqpErr.Code == amqp.NotFound {
				panic(NewFatalWorkerError(err))
			}

			panic(err)
		}
	})

	return api, services, nil
}

func registerWebsocketRooms(
	roomRegistry websocket.RoomRegistry,
	authorize func(perms ...string) websocket.Authorize,
	alarmWatcher alarmapi.Watcher,
	messageRateWatcher messageratestats.Watcher,
) error {
	if err := roomRegistry.Register(websocket.RoomBroadcastMessages, websocket.RoomHandlers{}); err != nil {
		return fmt.Errorf("fail to register websocket room: %w", err)
	}

	if err := roomRegistry.Register(websocket.RoomLoggedUserCount, websocket.RoomHandlers{}); err != nil {
		return fmt.Errorf("fail to register websocket room: %w", err)
	}

	if err := roomRegistry.Register(websocket.RoomNotifications, websocket.RoomHandlers{}); err != nil {
		return fmt.Errorf("fail to register websocket room: %w", err)
	}

	if err := roomRegistry.Register(websocket.RoomHealthcheck, websocket.RoomHandlers{Authorize: authorize(apisecurity.PermHealthcheck, securitymodel.PermissionCan)}); err != nil {
		return fmt.Errorf("fail to register websocket room: %w", err)
	}

	if err := roomRegistry.Register(websocket.RoomHealthcheckStatus, websocket.RoomHandlers{Authorize: authorize(apisecurity.PermHealthcheck, securitymodel.PermissionCan)}); err != nil {
		return fmt.Errorf("fail to register websocket room: %w", err)
	}

	if err := roomRegistry.Register(websocket.RoomIcons, websocket.RoomHandlers{}); err != nil {
		return fmt.Errorf("fail to register websocket room: %w", err)
	}

	if err := roomRegistry.Register(websocket.RoomPbhPatterns, websocket.RoomHandlers{Authorize: authorize(apisecurity.PermPbhPatterns, securitymodel.PermissionCan)}); err != nil {
		return fmt.Errorf("fail to register websocket room: %w", err)
	}

	err := roomRegistry.RegisterGroup(websocket.RoomAlarmsGroup, websocket.RoomHandlers{
		OnJoin:    alarmWatcher.StartWatch,
		OnLeave:   alarmWatcher.StopWatch,
		Authorize: authorize(apisecurity.PermAlarmRead, securitymodel.PermissionCan),
	})
	if err != nil {
		return fmt.Errorf("fail to register websocket group: %w", err)
	}

	err = roomRegistry.RegisterGroup(websocket.RoomAlarmDetailsGroup, websocket.RoomHandlers{
		OnJoin:    alarmWatcher.StartWatchDetails,
		OnLeave:   alarmWatcher.StopWatch,
		Authorize: authorize(apisecurity.PermAlarmRead, securitymodel.PermissionCan),
	})
	if err != nil {
		return fmt.Errorf("fail to register websocket group: %w", err)
	}

	err = roomRegistry.RegisterGroup(websocket.RoomMessageRates, websocket.RoomHandlers{
		OnJoin:    messageRateWatcher.StartWatch,
		OnLeave:   messageRateWatcher.StopWatch,
		Authorize: authorize(apisecurity.PermMessageRateStatsRead, securitymodel.PermissionCan),
	})
	if err != nil {
		return fmt.Errorf("fail to register websocket group: %w", err)
	}
	// TODO: add handlers for anomaly monitored connector state

	return nil
}

func GetSessionKeyVar(logger zerolog.Logger) []byte {
	sessionKey := os.Getenv("SESSION_KEY")
	if sessionKey == "" {
		logger.Warn().Msg("SESSION_KEY is not set, using default value")
		sessionKey = "canopsis"
	}
	return []byte(sessionKey)
}
