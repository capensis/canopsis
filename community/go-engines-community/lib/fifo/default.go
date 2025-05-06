package fifo

import (
	"context"
	"flag"
	"fmt"
	"time"

	apisecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/axe"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	libengine "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/externaldata"
	libflag "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/flag"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/healthcheck"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	libscheduler "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/scheduler"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/techmetrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/usernotification"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/che"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/depmake"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/postgres"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/rs/zerolog"
)

type Options struct {
	Version                bool
	PrintEventOnError      bool
	ModeDebug              bool
	LockTtl                int
	PeriodicalWaitTime     time.Duration
	ExternalDataApiTimeout time.Duration
	Workers                int
}

type Services struct {
	DbClient                    mongo.DbClient
	PgPoolProvider              postgres.PoolProvider
	Cfg                         config.CanopsisConf
	ExternalDataContainer       *externaldata.GetterContainer
	TimezoneConfigProvider      config.TimezoneConfigProvider
	TemplateConfigProvider      config.TemplateConfigProvider
	EventFilterFailureService   eventfilter.FailureService
	DataStoragePeriodicalWorker datastorage.PeriodicalWorker
}

func ParseOptions() (Options, []string) {
	var opts Options

	flag.BoolVar(&opts.ModeDebug, "d", false, "debug")
	flag.BoolVar(&opts.PrintEventOnError, "printEventOnError", false, "Print event on processing error")
	flag.IntVar(&opts.LockTtl, "lockTtl", 10, "Redis lock ttl time in seconds")
	flag.DurationVar(&opts.PeriodicalWaitTime, "periodicalWaitTime", canopsis.PeriodicalWaitTime, "Duration to wait between two run of periodical process")
	flag.DurationVar(&opts.ExternalDataApiTimeout, "externalDataApiTimeout", 30*time.Second, "External API HTTP Request Timeout.")
	flag.BoolVar(&opts.Version, "version", false, "Show the version information")
	flag.IntVar(&opts.Workers, "workers", canopsis.DefaultEventWorkers, "Amount of workers to process fifo_ack events flow")

	flag.Duration("eventsStatsFlushInterval", 60*time.Second, "Deprecated: interval between saving statistics from redis to mongo")
	flag.String("publishQueue", "", "Deprecated: publish event to this queue.")
	flag.String("consumeQueue", "", "Deprecated: consume events from this queue.")

	flag.Parse()

	return opts, libflag.FindDeprecatedFlags("eventsStatsFlushInterval", "consumeQueue", "publishQueue")
}

func Default(
	ctx context.Context,
	options Options,
	logger zerolog.Logger,
) (libengine.Engine, Services) {
	var m depmake.DependencyMaker
	s := Services{}

	s.DbClient = m.DepMongoClient(ctx, logger)
	s.Cfg = m.DepConfig(ctx, s.DbClient)
	config.SetDbClientRetry(s.DbClient, s.Cfg)
	s.PgPoolProvider = postgres.NewPoolProvider(s.Cfg.Global.ReconnectRetries, s.Cfg.Global.GetReconnectTimeout())
	metricsConfigProvider := config.NewMetricsConfigProvider(s.Cfg, logger)
	metricsSender := metrics.NewTimescaleDBSender(s.PgPoolProvider, metricsConfigProvider, logger)
	s.ExternalDataContainer = externaldata.NewGetterContainer()
	timezoneConfigProvider := config.NewTimezoneConfigProvider(s.Cfg, logger)
	s.TimezoneConfigProvider = timezoneConfigProvider
	templateConfigProvider := config.NewTemplateConfigProvider(s.Cfg, logger)
	s.TemplateConfigProvider = templateConfigProvider
	dataStorageConfigProvider := config.NewDataStorageConfigProvider(s.Cfg, logger)
	amqpConnection := m.DepAmqpConnection(logger, s.Cfg)
	amqpChannel := m.DepAMQPChannelPub(amqpConnection)
	lockRedisClient := m.DepRedisSession(ctx, redis.LockStorage, logger, s.Cfg)
	engineLockRedisClient := m.DepRedisSession(ctx, redis.EngineLockStorage, logger, s.Cfg)
	queueRedisClient := m.DepRedisSession(ctx, redis.QueueStorage, logger, s.Cfg)
	runInfoRedisClient := m.DepRedisSession(ctx, redis.EngineRunInfo, logger, s.Cfg)
	scheduler := libscheduler.NewSchedulerService(
		lockRedisClient,
		queueRedisClient,
		m.DepAMQPChannelPub(m.DepAmqpConnection(logger, s.Cfg)),
		canopsis.CheQueuePrefix,
		logger,
		options.LockTtl,
		json.NewDecoder(),
		json.NewEncoder(),
	)
	eventFilterEventCounter := eventfilter.NewEventCounter(s.DbClient,
		utils.MinDuration(canopsis.DefaultFlushInterval, options.PeriodicalWaitTime), logger)
	notifStore := usernotification.NewStore(s.DbClient, amqpChannel, json.NewEncoder(),
		canopsis.ApiNotificationExchangeName, "", canopsis.JsonContentType)
	s.EventFilterFailureService = eventfilter.NewFailureService(s.DbClient, notifStore,
		utils.MinDuration(canopsis.DefaultFlushInterval, options.PeriodicalWaitTime), apisecurity.ObjEventFilter, logger)
	templateExecutor := template.NewExecutor(templateConfigProvider, timezoneConfigProvider)
	ruleAdapter := eventfilter.NewRuleAdapter(s.DbClient)
	ruleApplicatorContainer := eventfilter.NewRuleApplicatorContainer()
	ruleApplicatorContainer.Set(eventfilter.RuleTypeChangeEntity, eventfilter.NewChangeEntityApplicator(s.ExternalDataContainer, s.EventFilterFailureService, templateExecutor))
	eventfilterService := eventfilter.NewRuleService(ruleAdapter, ruleApplicatorContainer, eventFilterEventCounter, s.EventFilterFailureService, templateExecutor, logger)
	techMetricsConfigProvider := config.NewTechMetricsConfigProvider(s.Cfg, logger)
	techMetricsSender := techmetrics.NewSender(canopsis.FIFOEngineName+"/"+utils.NewID(), techMetricsConfigProvider, canopsis.TechMetricsFlushInterval,
		s.Cfg.Global.ReconnectRetries, s.Cfg.Global.GetReconnectTimeout(), logger)

	healthCheckCfg, err := config.NewHealthCheckAdapter(s.DbClient).GetConfig(ctx)
	if err != nil {
		panic(fmt.Errorf("cannot load healthcheck config: %w", err))
	}

	runInfoPeriodicalWorker := libengine.NewRunInfoPeriodicalWorker(
		healthCheckCfg.ParseUpdateInterval(logger),
		libengine.NewRunInfoManager(runInfoRedisClient),
		libengine.NewInstanceRunInfo(canopsis.FIFOEngineName, canopsis.FIFOQueueName, canopsis.CheQueuePrefix, []string{canopsis.FIFOQueueName}),
		amqpChannel,
		logger,
	)

	queueMetricsPeriodicalWorker := libengine.NewQueueMetricsPeriodicalWorker(
		options.PeriodicalWaitTime,
		amqpChannel,
		techMetricsSender,
		[]string{canopsis.FIFOQueueName},
		techmetrics.FIFOQueue,
		logger,
	)

	engine := libengine.New(
		func(ctx context.Context) error {
			runInfoPeriodicalWorker.Work(ctx)
			queueMetricsPeriodicalWorker.Work(ctx)

			scheduler.Start(ctx)

			if !s.DbClient.IsDistributed() {
				err := eventfilterService.LoadRules(ctx, []string{eventfilter.RuleTypeChangeEntity})
				if err != nil {
					return err
				}
			}

			return nil
		},
		func(ctx context.Context) {
			scheduler.Stop(context.WithoutCancel(ctx))
			err := s.DbClient.Disconnect(context.WithoutCancel(ctx))
			if err != nil {
				logger.Error().Err(err).Msg("failed to close mongo connection")
			}

			err = amqpConnection.Close()
			if err != nil {
				logger.Error().Err(err).Msg("failed to close amqp connection")
			}

			err = lockRedisClient.Close()
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

			s.PgPoolProvider.Close()
		},
		logger,
	)

	engine.AddRoutine(func(ctx context.Context) error {
		techMetricsSender.Run(ctx)
		return nil
	})

	mainMessageProcessor := &messageProcessor{
		FeaturePrintEventOnError: options.PrintEventOnError,

		EventFilterService: eventfilterService,
		TechMetricsSender:  techMetricsSender,
		Scheduler:          scheduler,
		MetricsSender:      metricsSender,
		Decoder:            json.NewDecoder(),
		Logger:             logger,
	}

	engine.AddConsumer(libengine.NewConcurrentConsumer(
		canopsis.FIFOConsumerName,
		canopsis.FIFOQueueName,
		s.Cfg.Global.PrefetchCount,
		s.Cfg.Global.PrefetchSize,
		false,
		"",
		"",
		"",
		"",
		1, // TODO: 1 worker for now, to think about making fifo concurrent
		false,
		amqpConnection,
		mainMessageProcessor,
		logger,
	))
	engine.AddConsumer(libengine.NewConcurrentConsumer(
		canopsis.FIFOAckConsumerName,
		canopsis.FIFOAckQueueName,
		s.Cfg.Global.PrefetchCount,
		s.Cfg.Global.PrefetchSize,
		false,
		"",
		"",
		"",
		"",
		options.Workers,
		false,
		amqpConnection,
		&ackMessageProcessor{
			FeaturePrintEventOnError: options.PrintEventOnError,

			Scheduler:         scheduler,
			TechMetricsSender: techMetricsSender,
			Decoder:           json.NewDecoder(),
			Logger:            logger,
		},
		logger,
	))

	engine.AddPeriodicalWorker("run_info", runInfoPeriodicalWorker)
	engine.AddPeriodicalWorker("queue_metrics", queueMetricsPeriodicalWorker)

	s.DataStoragePeriodicalWorker = datastorage.NewPeriodicalWorker(
		func(ctx context.Context, clientTimeout time.Duration) (mongo.DbClient, error) {
			return mongo.NewClientWithOptions(ctx, 0, 0, mongo.DefaultServerSelectionTimeout,
				clientTimeout, logger)
		},
		time.Hour,
		timezoneConfigProvider,
		dataStorageConfigProvider,
		logger,
	)
	s.DataStoragePeriodicalWorker.AddCleaner("alarm", alarm.NewCleaner(logger))
	s.DataStoragePeriodicalWorker.AddCleaner("alarm_external_tag", axe.NewExternalTagCleaner(logger))
	s.DataStoragePeriodicalWorker.AddCleaner("pbehavior", pbehavior.NewCleaner(logger))
	s.DataStoragePeriodicalWorker.AddCleaner("event_filter_failure", che.NewEventFailureCleaner(amqpChannel, json.NewEncoder(),
		canopsis.ApiNotificationExchangeName, "", canopsis.JsonContentType, logger))
	engine.AddPeriodicalWorker("datastorage", libengine.NewLockedPeriodicalWorker(
		redis.NewLockClient(engineLockRedisClient),
		redis.FifoDataStorageLockKey,
		s.DataStoragePeriodicalWorker,
		logger,
	))
	engine.AddPeriodicalWorker("config", libengine.NewLoadConfigPeriodicalWorker(
		options.PeriodicalWaitTime,
		config.NewAdapter(s.DbClient),
		logger,
		timezoneConfigProvider,
		techMetricsConfigProvider,
		dataStorageConfigProvider,
		templateConfigProvider,
		metricsConfigProvider,
	))
	if s.DbClient.IsDistributed() {
		engine.AddRoutine(func(ctx context.Context) error {
			w := eventfilter.NewRulesChangesWatcher(s.DbClient, eventfilterService)

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
	} else {
		engine.AddPeriodicalWorker("local_cache", &periodicalWorker{
			RuleService:        eventfilterService,
			PeriodicalInterval: options.PeriodicalWaitTime,
			Logger:             logger,
		})
	}
	engine.AddRoutine(func(ctx context.Context) error {
		eventFilterEventCounter.Run(ctx)

		return nil
	})
	engine.AddRoutine(func(ctx context.Context) error {
		s.EventFilterFailureService.Run(ctx)

		return nil
	})
	engine.AddRoutine(func(ctx context.Context) error {
		metricsSender.Run(ctx)

		return nil
	})

	healthcheck.Start(ctx, healthcheck.NewChecker(
		"fifo",
		mainMessageProcessor,
		json.NewEncoder(),
		false,
		false,
	), logger)

	return engine, s
}
