package api

import (
	"context"
	"os"
	"time"

	"git.canopsis.net/canopsis/go-engines/lib/amqp"
	"git.canopsis.net/canopsis/go-engines/lib/api/export"
	apilogger "git.canopsis.net/canopsis/go-engines/lib/api/logger"
	"git.canopsis.net/canopsis/go-engines/lib/api/middleware"
	devmiddleware "git.canopsis.net/canopsis/go-engines/lib/api/middleware/dev"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/config"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/encoding/json"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/engine"
	libpbehavior "git.canopsis.net/canopsis/go-engines/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	libredis "git.canopsis.net/canopsis/go-engines/lib/redis"
	libsecurity "git.canopsis.net/canopsis/go-engines/lib/security"
	"git.canopsis.net/canopsis/go-engines/lib/security/proxy"
	"git.canopsis.net/canopsis/go-engines/lib/security/session/mongostore"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

const chanBuf = 10
const sessionStoreSessionMaxAge = 24 * time.Hour
const sessionStoreAutoCleanInterval = 10 * time.Second

func Default(
	addr string,
	configDir string,
	secureSession bool,
	test bool,
	enforcer libsecurity.Enforcer,
	logger zerolog.Logger,
) (API, error) {
	// Retrieve config.
	dbClient, err := mongo.NewClient(0, 0)
	if err != nil {
		logger.Err(err).Msg("cannot connect to mongodb")
		return nil, err
	}
	cfg, err := config.NewAdapter(dbClient).GetConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot load config")
	}
	// Connect to mongodb.
	dbClient, err = mongo.NewClient(
		cfg.Global.ReconnectRetries,
		cfg.Global.GetReconnectTimeout(),
	)
	if err != nil {
		logger.Err(err).Msg("cannot connect to mongodb")
		return nil, err
	}
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
	pbhRedisSession, err := libredis.NewSession(libredis.PBehaviorLockStorage, logger,
		cfg.Global.ReconnectRetries, cfg.Global.GetReconnectTimeout())
	if err != nil {
		logger.Err(err).Msg("cannot connect to redis")
		return nil, err
	}
	engineRedisSession, err := libredis.NewSession(libredis.EngineRunInfo, logger,
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
	// Retrieve app timezone.
	location, err := cfg.Timezone.GetLocation()
	if err != nil {
		logger.Err(err).Msg("cannot load timezone")
		return nil, err
	}

	sessionStore := mongostore.NewStore(dbClient, []byte(os.Getenv("SESSION_KEY")))
	sessionStore.Options.MaxAge = int(sessionStoreSessionMaxAge.Seconds())
	sessionStore.Options.Secure = secureSession
	security := NewSecurity(securityConfig, dbClient, sessionStore)

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
	// Create csv exporter.
	exportExecutor := export.NewTaskExecutor(dbClient, logger)
	// Create api.
	api := New(addr, logger)
	api.AddRouter(func(router gin.IRouter) {
		router.Use(middleware.Cache())

		if test {
			router.Use(devmiddleware.ReloadEnforcerPolicy(enforcer))
		}

		RegisterValidators(dbClient)
		RegisterRoutes(
			router,
			security,
			enforcer,
			dbClient,
			location,
			pbhStore,
			pbhService,
			pbhComputeChan,
			engine.NewRunInfoManager(engineRedisSession),
			exportExecutor,
			apilogger.NewActionLogger(dbClient, logger),
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
	api.AddWorker("data export", func(ctx context.Context) {
		exportExecutor.Execute(ctx)
	})

	return api, nil
}
