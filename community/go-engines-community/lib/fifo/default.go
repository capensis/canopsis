package fifo

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	apisecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/axe"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	libengine "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	libentity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/externaldata"
	libflag "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/flag"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/healthcheck"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	libprometheus "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics/prometheus"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	libscheduler "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/scheduler"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/techmetrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/usernotification"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/che"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/depprovider"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/postgres"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/bsm/redislock"
	goredis "github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type Options struct {
	log.Options
	Version            bool
	PrintEventOnError  bool
	LockTtl            int
	PeriodicalWaitTime time.Duration
	Workers            int
	RpcWorkers         int
	DataStorageCleanUp bool

	PrometheusExporterPort   int
	EnablePrometheusExporter bool
}

type Services struct {
	Cfg                         config.CanopsisConf
	DbClient                    mongo.DbClient
	PgPoolProvider              postgres.PoolProvider
	QueueRedisClient            goredis.Cmdable
	AmqpPublisher               amqp.Publisher
	AmqpConsumeChPool           amqp.ChannelPool
	DataStoragePeriodicalWorker datastorage.PeriodicalWorker
	EventFilterService          eventfilter.Service
	TechMetricsSender           techmetrics.Sender
	Scheduler                   libscheduler.Scheduler
	MessageProcessor            *messageProcessor
}

// Dependencies injects the engine-flavour-specific collaborators that need the
// shared infrastructure Default builds.
type Dependencies struct {
	NewMetaUpdater        func(mongo.DbClient, postgres.PoolProvider, config.CanopsisConf, zerolog.Logger) metrics.AsyncMetaUpdater
	NewExternalDataGetter func(mongo.DbClient, postgres.PoolProvider, template.Executor) externaldata.Getter
}

func ParseOptions() (Options, []string) {
	var opts Options

	log.BindCmdFlags(&opts.Options)
	flag.BoolVar(&opts.PrintEventOnError, "printEventOnError", false, "Print event on processing error")
	flag.IntVar(&opts.LockTtl, "lockTtl", 10, "Redis lock ttl time in seconds")
	flag.DurationVar(&opts.PeriodicalWaitTime, "periodicalWaitTime", canopsis.PeriodicalWaitTime, "Duration to wait between two run of periodical process")
	flag.BoolVar(&opts.Version, "version", false, "Show the version information")
	flag.IntVar(&opts.Workers, "workers", canopsis.DefaultEventWorkers, "Amount of workers to process fifo_ack events flow")
	flag.IntVar(&opts.RpcWorkers, "rpcWorkers", canopsis.DefaultRpcWorkers, "Amount of workers to process rpc event flow.")
	flag.BoolVar(&opts.DataStorageCleanUp, "cleanUp", false, "Immediately execute all data storage archive and delete and exit after.")
	flag.BoolVar(&opts.EnablePrometheusExporter, "enablePrometheusExporter", false, "Enable prometheus exporter")
	flag.IntVar(&opts.PrometheusExporterPort, "prometheusExporterPort", libprometheus.DefaultExporterPort, "Prometheus exporter port")

	flag.Duration("externalDataApiTimeout", 30*time.Second, "Deprecated: External API HTTP Request Timeout.")

	flag.Parse()

	return opts, libflag.FindDeprecatedFlags("externalDataApiTimeout")
}

func Default(
	ctx context.Context,
	deps Dependencies,
	options Options,
	dp depprovider.Provider,
	logger zerolog.Logger,
) (e libengine.Engine, s Services, err error) {
	dbClient, err := dp.MongoClient(ctx, mongo.ClientOptions{})
	if err != nil {
		return e, s, err
	}

	cfg, err := dp.Config(ctx, dbClient)
	if err != nil {
		return e, s, err
	}

	config.SetDbClientRetry(dbClient, cfg)

	// noTimeoutClient should be used by change stream watchers only.
	noTimeoutClient, err := dp.MongoClient(ctx, mongo.ClientOptions{
		RetryCount:      cfg.Global.ReconnectRetries,
		MinRetryTimeout: cfg.Global.GetReconnectTimeout(),
		NoClientTimeout: true,
	})
	if err != nil {
		return e, s, err
	}

	pgPoolProvider := postgres.NewPoolProvider(cfg.Global.ReconnectRetries, cfg.Global.GetReconnectTimeout())
	metricsConfigProvider := config.NewMetricsConfigProvider(cfg, logger)
	metricsSender := metrics.NewTimescaleDBSender(pgPoolProvider, metricsConfigProvider, logger)
	timezoneConfigProvider := config.NewTimezoneConfigProvider(cfg, logger)
	templateConfigProvider := config.NewTemplateConfigProvider(cfg, logger)
	dataStorageConfigProvider := config.NewDataStorageConfigProvider(cfg, logger)
	amqpPubConn, err := dp.AMQPConnection(logger, cfg)
	if err != nil {
		return e, s, err
	}

	amqpConsumeConn, err := dp.AMQPConnection(logger, cfg)
	if err != nil {
		return e, s, err
	}

	amqpPubChPool := dp.AMQPPubChannelPool(amqpPubConn)
	amqpPublisher := amqp.NewPooledPublisher(amqpPubChPool)
	amqpConsumeChPool := dp.AMQPConsumeChannelPool(amqpConsumeConn)
	lockRedisClient, err := dp.RedisClient(ctx, redis.LockStorage, logger, cfg)
	if err != nil {
		return e, s, err
	}

	engineLockRedisClient, err := dp.RedisClient(ctx, redis.EngineLockStorage, logger, cfg)
	if err != nil {
		return e, s, err
	}

	queueRedisClient, err := dp.RedisClient(ctx, redis.QueueStorage, logger, cfg)
	if err != nil {
		return e, s, err
	}

	runInfoRedisClient, err := dp.RedisClient(ctx, redis.EngineRunInfo, logger, cfg)
	if err != nil {
		return e, s, err
	}

	s.Scheduler = libscheduler.NewSchedulerService(
		lockRedisClient,
		queueRedisClient,
		amqpPublisher,
		canopsis.CheQueuePrefix,
		logger,
		options.LockTtl,
		json.NewDecoder(),
		json.NewEncoder(),
	)
	eventFilterEventCounter := eventfilter.NewEventCounter(dbClient,
		utils.MinDuration(canopsis.DefaultFlushInterval, options.PeriodicalWaitTime), logger)
	notifStore := usernotification.NewStore(dbClient, amqpPublisher, json.NewEncoder(),
		canopsis.ApiNotificationExchangeName, "", canopsis.JsonContentType)
	eventFilterFailureService := eventfilter.NewFailureService(dbClient, notifStore,
		utils.MinDuration(canopsis.DefaultFlushInterval, options.PeriodicalWaitTime), apisecurity.ObjEventFilterRule, logger)
	templateExecutor := template.NewExecutor(templateConfigProvider, timezoneConfigProvider)
	externalDataGetter := deps.NewExternalDataGetter(dbClient, pgPoolProvider, templateExecutor)
	ruleAdapter := eventfilter.NewRuleAdapter(dbClient)
	ruleApplicatorContainer := eventfilter.NewRuleApplicatorContainer()
	ruleApplicatorContainer.Set(eventfilter.RuleTypeChangeEntity, eventfilter.NewChangeEntityApplicator(eventFilterFailureService, templateExecutor))
	s.EventFilterService = eventfilter.NewRuleService(ruleAdapter, ruleApplicatorContainer, externalDataGetter, eventFilterEventCounter, eventFilterFailureService, templateExecutor, logger)
	techMetricsConfigProvider := config.NewTechMetricsConfigProvider(cfg, logger)
	s.TechMetricsSender = techmetrics.NewSender(canopsis.FIFOEngineName+"/"+utils.NewID(), techMetricsConfigProvider, canopsis.TechMetricsFlushInterval,
		cfg.Global.ReconnectRetries, cfg.Global.GetReconnectTimeout(), logger)

	healthCheckCfg, err := config.NewHealthCheckAdapter(dbClient).GetConfig(ctx)
	if err != nil {
		return e, s, fmt.Errorf("cannot load healthcheck config: %w", err)
	}

	runInfoPeriodicalWorker := libengine.NewRunInfoPeriodicalWorker(
		healthCheckCfg.ParseUpdateInterval(logger),
		libengine.NewRunInfoManager(runInfoRedisClient),
		libengine.NewInstanceRunInfo(canopsis.FIFOEngineName, canopsis.FIFOQueueName, canopsis.CheQueuePrefix, []string{canopsis.FIFOQueueName}),
		amqpConsumeChPool,
		logger,
	)

	queueMetricsPeriodicalWorker := libengine.NewQueueMetricsPeriodicalWorker(
		options.PeriodicalWaitTime,
		amqpConsumeChPool,
		s.TechMetricsSender,
		[]string{canopsis.FIFOQueueName},
		techmetrics.FIFOQueue,
		logger,
	)

	prometheusMetrics := libprometheus.NewFifoMetrics()

	mainMessageProcessor := NewMessageProcessor(
		s.EventFilterService,
		s.Scheduler,
		metricsSender,
		json.NewEncoder(),
		json.NewDecoder(),
		logger,
		s.TechMetricsSender,
		prometheusMetrics,
		options.PrintEventOnError,
	)

	s.Cfg = cfg
	s.DbClient = dbClient
	s.PgPoolProvider = pgPoolProvider
	s.QueueRedisClient = queueRedisClient
	s.AmqpPublisher = amqpPublisher
	s.AmqpConsumeChPool = amqpConsumeChPool
	s.MessageProcessor = mainMessageProcessor

	lockDuration := max(options.PeriodicalWaitTime, lockMinDuration) + lockBackoff*time.Duration(lockRetries)

	engine := libengine.New(
		func(ctx context.Context) (err error) {
			err = mainMessageProcessor.ObtainExclusive(ctx, redis.NewLockClient(engineLockRedisClient), lockDuration)
			if err != nil {
				if !errors.Is(err, redislock.ErrNotObtained) {
					logger.Err(err).Msg("cannot obtain lock for engine initialization, exiting")
				}
				return err
			}
			runInfoPeriodicalWorker.Work(ctx)
			queueMetricsPeriodicalWorker.Work(ctx)

			s.Scheduler.Start(ctx)

			return nil
		},
		func(ctx context.Context) {
			s.Scheduler.Stop(context.WithoutCancel(ctx))

			err := dbClient.Disconnect(context.WithoutCancel(ctx))
			if err != nil {
				logger.Error().Err(err).Msg("failed to close mongo connection")
			}

			err = noTimeoutClient.Disconnect(context.WithoutCancel(ctx))
			if err != nil {
				logger.Error().Err(err).Msg("failed to close mongo connection without timeout")
			}

			err = amqpPubChPool.Close()
			if err != nil {
				logger.Err(err).Msg("failed to close amqp publish channels")
			}

			err = amqpConsumeChPool.Close()
			if err != nil {
				logger.Err(err).Msg("failed to close amqp consumer channels")
			}

			err = amqpPubConn.Close()
			if err != nil {
				logger.Error().Err(err).Msg("failed to close amqp connection")
			}

			err = amqpConsumeConn.Close()
			if err != nil {
				logger.Error().Err(err).Msg("failed to close amqp connection")
			}

			if err = mainMessageProcessor.ReleaseExclusive(ctx); err != nil && !errors.Is(err, redislock.ErrLockNotHeld) {
				logger.Warn().Err(err).Msg("failed to release exclusive processor lock")
			}

			err = lockRedisClient.Close()
			if err != nil {
				logger.Error().Err(err).Msg("failed to close redis connection")
			}

			err = engineLockRedisClient.Close()
			if err != nil {
				logger.Error().Err(err).Msg("failed to close redis connection")
			}

			err = queueRedisClient.Close()
			if err != nil {
				logger.Error().Err(err).Msg("failed to close redis connection")
			}

			err = runInfoRedisClient.Close()
			if err != nil {
				logger.Error().Err(err).Msg("failed to close redis connection")
			}

			pgPoolProvider.Close()
		},
		logger,
	)

	engine.AddRoutine(func(ctx context.Context) error {
		s.TechMetricsSender.Run(ctx)
		return nil
	})

	engine.AddRoutine(func(ctx context.Context) error {
		return mainMessageProcessor.RefreshExclusiveProcessor(ctx, options.PeriodicalWaitTime, lockDuration)
	})

	if options.EnablePrometheusExporter {
		engine.AddRoutine(func(ctx context.Context) error {
			return libprometheus.RunPrometheusExporter(ctx, options.PrometheusExporterPort, logger, prometheusMetrics)
		})
	}

	engine.AddConsumer(libengine.NewConcurrentConsumer(
		canopsis.FIFOConsumerName,
		canopsis.FIFOQueueName,
		cfg.Global.PrefetchCount,
		cfg.Global.PrefetchSize,
		false,
		"",
		"",
		"",
		"",
		1, // TODO: 1 worker for now, to think about making fifo concurrent
		false,
		amqpPublisher,
		amqpConsumeChPool,
		mainMessageProcessor,
		logger,
	))
	engine.AddConsumer(libengine.NewConcurrentConsumer(
		canopsis.FIFOAckConsumerName,
		canopsis.FIFOAckQueueName,
		cfg.Global.PrefetchCount,
		cfg.Global.PrefetchSize,
		false,
		"",
		"",
		"",
		"",
		options.Workers,
		false,
		amqpPublisher,
		amqpConsumeChPool,
		&ackMessageProcessor{
			FeaturePrintEventOnError: options.PrintEventOnError,

			Scheduler:         s.Scheduler,
			TechMetricsSender: s.TechMetricsSender,
			Decoder:           json.NewDecoder(),
			Logger:            logger,
		},
		logger,
	))

	engine.AddPeriodicalWorker("run_info", runInfoPeriodicalWorker)
	engine.AddPeriodicalWorker("queue_metrics", queueMetricsPeriodicalWorker)

	s.DataStoragePeriodicalWorker = datastorage.NewPeriodicalWorker(
		func(ctx context.Context, clientTimeout time.Duration) (mongo.DbClient, error) {
			return mongo.NewClient(ctx, mongo.ClientOptions{
				ClientTimeout: clientTimeout,
			})
		},
		time.Hour,
		timezoneConfigProvider,
		dataStorageConfigProvider,
		logger,
	)

	metricsEntityMetaUpdater := deps.NewMetaUpdater(dbClient, pgPoolProvider, cfg, logger)
	disabledEntityCleaner := libentity.NewCleaner(
		redis.NewLockClient(engineLockRedisClient),
		datastorage.NewAdapter(dbClient),
		dataStorageConfigProvider,
		metricsEntityMetaUpdater,
		logger,
	)
	s.DataStoragePeriodicalWorker.OnSchedule(true)
	s.DataStoragePeriodicalWorker.AddCleaner("entity", disabledEntityCleaner)
	s.DataStoragePeriodicalWorker.AddCleaner("alarm", alarm.NewCleaner(logger))
	s.DataStoragePeriodicalWorker.AddCleaner("alarm_external_tag", axe.NewExternalTagCleaner(logger))
	s.DataStoragePeriodicalWorker.AddCleaner("pbehavior", pbehavior.NewCleaner(logger))
	s.DataStoragePeriodicalWorker.AddCleaner("event_filter_failure", che.NewEventFailureCleaner(amqpPublisher, json.NewEncoder(),
		canopsis.ApiNotificationExchangeName, "", canopsis.JsonContentType, logger))
	engine.AddPeriodicalWorker("datastorage", libengine.NewLockedPeriodicalWorker(
		redis.NewLockClient(engineLockRedisClient),
		redis.FifoDataStorageLockKey,
		s.DataStoragePeriodicalWorker,
		logger,
	))
	engine.AddPeriodicalWorker("config", libengine.NewLoadConfigPeriodicalWorker(
		options.PeriodicalWaitTime,
		config.NewAdapter(dbClient),
		logger,
		timezoneConfigProvider,
		techMetricsConfigProvider,
		dataStorageConfigProvider,
		templateConfigProvider,
		metricsConfigProvider,
	))
	engine.AddRoutine(func(ctx context.Context) error {
		w := eventfilter.NewRulesChangesWatcher(noTimeoutClient, s.EventFilterService)

		logger.Debug().Msg("Loading event filter rules")

		for {
			select {
			case <-ctx.Done():
				return nil
			default:
				err := w.Watch(ctx, []string{eventfilter.RuleTypeChangeEntity})
				if err != nil {
					logger.Error().Err(err).Msg("failed to watch eventfilter collection")
				}
			}
		}
	})
	engine.AddRoutine(func(ctx context.Context) error {
		eventFilterEventCounter.Run(ctx)

		return nil
	})
	engine.AddRoutine(func(ctx context.Context) error {
		eventFilterFailureService.Run(ctx)

		return nil
	})
	engine.AddRoutine(func(ctx context.Context) error {
		metricsSender.Run(ctx)

		return nil
	})
	engine.AddRoutine(func(ctx context.Context) error {
		metricsEntityMetaUpdater.Run(ctx)

		return nil
	})

	healthcheck.Start(ctx, healthcheck.NewChecker(
		"fifo",
		mainMessageProcessor,
		json.NewEncoder(),
		false,
		false,
	), logger)

	return engine, s, nil
}
