package fifo

import (
	"context"
	"flag"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	libengine "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/ratelimit"
	libscheduler "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/scheduler"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/statistics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/depmake"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"github.com/rs/zerolog"
)

type Options struct {
	Version                  bool
	PrintEventOnError        bool
	ModeDebug                bool
	ConsumeFromQueue         string
	PublishToQueue           string
	LockTtl                  int
	EventsStatsFlushInterval time.Duration
	PeriodicalWaitTime       time.Duration
	DataSourceDirectory      string
	ExternalDataApiTimeout   time.Duration
}

func ParseOptions() Options {
	var opts Options

	flag.StringVar(&opts.PublishToQueue, "publishQueue", canopsis.CheQueueName, "Publish event to this queue.")
	flag.StringVar(&opts.ConsumeFromQueue, "consumeQueue", canopsis.FIFOQueueName, "Consume events from this queue.")
	flag.BoolVar(&opts.ModeDebug, "d", false, "debug")
	flag.BoolVar(&opts.PrintEventOnError, "printEventOnError", false, "Print event on processing error")
	flag.IntVar(&opts.LockTtl, "lockTtl", 10, "Redis lock ttl time in seconds")
	flag.DurationVar(&opts.EventsStatsFlushInterval, "eventsStatsFlushInterval", 60*time.Second, "Interval between saving statistics from redis to mongo")
	flag.DurationVar(&opts.PeriodicalWaitTime, "periodicalWaitTime", canopsis.PeriodicalWaitTime, "Duration to wait between two run of periodical process")
	flag.StringVar(&opts.DataSourceDirectory, "dataSourceDirectory", ".", "The path of the directory containing the event filter's data source plugins.")
	flag.DurationVar(&opts.ExternalDataApiTimeout, "externalDataApiTimeout", 30*time.Second, "External API HTTP Request Timeout.")
	flag.Bool("enableMetaAlarmProcessing", true, "Enable meta-alarm processing - deprecated")
	flag.BoolVar(&opts.Version, "version", false, "Show the version information")

	flag.Parse()

	return opts
}

func Default(ctx context.Context, options Options, mongoClient mongo.DbClient, ExternalDataContainer *eventfilter.ExternalDataContainer, logger zerolog.Logger) libengine.Engine {
	var m depmake.DependencyMaker

	cfg := m.DepConfig(ctx, mongoClient)
	config.SetDbClientRetry(mongoClient, cfg)
	timezoneConfigProvider := config.NewTimezoneConfigProvider(cfg, logger)
	amqpConnection := m.DepAmqpConnection(logger, cfg)
	amqpChannel := m.DepAMQPChannelPub(amqpConnection)
	lockRedisClient := m.DepRedisSession(ctx, redis.LockStorage, logger, cfg)
	engineLockRedisClient := m.DepRedisSession(ctx, redis.EngineLockStorage, logger, cfg)
	queueRedisClient := m.DepRedisSession(ctx, redis.QueueStorage, logger, cfg)
	statsRedisClient := m.DepRedisSession(ctx, redis.FIFOMessageStatisticsStorage, logger, cfg)
	runInfoRedisClient := m.DepRedisSession(ctx, redis.EngineRunInfo, logger, cfg)
	scheduler := libscheduler.NewSchedulerService(
		lockRedisClient,
		queueRedisClient,
		m.DepAMQPChannelPub(m.DepAmqpConnection(logger, cfg)),
		options.PublishToQueue,
		logger,
		options.LockTtl,
		json.NewDecoder(),
		json.NewEncoder(),
	)
	statsCh := make(chan statistics.Message)
	statsSender := ratelimit.NewStatsSender(statsCh, logger)
	statsListener := statistics.NewStatsListener(
		mongoClient,
		statsRedisClient,
		options.EventsStatsFlushInterval,
		map[string]int64{
			mongo.MessageRateStatsMinuteCollectionName: 1,  // 1 minute
			mongo.MessageRateStatsHourCollectionName:   60, // 60 minutes
		},
		logger,
	)

	ruleAdapter := eventfilter.NewRuleAdapter(mongoClient)
	ruleApplicatorContainer := eventfilter.NewRuleApplicatorContainer()
	ruleApplicatorContainer.Set(eventfilter.RuleTypeChangeEntity, eventfilter.NewChangeEntityApplicator(ExternalDataContainer))
	eventfilterService := eventfilter.NewRuleService(ruleAdapter, ruleApplicatorContainer, config.NewTimezoneConfigProvider(cfg, logger), logger)

	runInfoPeriodicalWorker := libengine.NewRunInfoPeriodicalWorker(
		canopsis.PeriodicalWaitTime,
		libengine.NewRunInfoManager(runInfoRedisClient),
		libengine.NewInstanceRunInfo(canopsis.FIFOEngineName, options.ConsumeFromQueue, options.PublishToQueue),
		amqpChannel,
		logger,
	)

	engine := libengine.New(
		func(ctx context.Context) error {
			runInfoPeriodicalWorker.Work(ctx)
			scheduler.Start(ctx)

			if !mongoClient.IsDistributed() {
				err := eventfilterService.LoadRules(ctx, []string{eventfilter.RuleTypeChangeEntity})
				if err != nil {
					return err
				}
			}

			go statsListener.Listen(ctx, statsCh)

			return nil
		},
		func(ctx context.Context) {
			close(statsCh)

			scheduler.Stop(ctx)
			err := mongoClient.Disconnect(ctx)
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

			err = statsRedisClient.Close()
			if err != nil {
				logger.Error().Err(err).Msg("failed to close redis connection")
			}

			err = runInfoRedisClient.Close()
			if err != nil {
				logger.Error().Err(err).Msg("failed to close redis connection")
			}
		},
		logger,
	)

	engine.AddConsumer(libengine.NewDefaultConsumer(
		canopsis.FIFOConsumerName,
		options.ConsumeFromQueue,
		cfg.Global.PrefetchCount,
		cfg.Global.PrefetchSize,
		false,
		"",
		"",
		"",
		"",
		amqpConnection,
		&messageProcessor{
			FeaturePrintEventOnError: options.PrintEventOnError,
			EventFilterService:       eventfilterService,
			Scheduler:                scheduler,
			StatsSender:              statsSender,
			Decoder:                  json.NewDecoder(),
			Logger:                   logger,
		},
		logger,
	))
	engine.AddConsumer(libengine.NewDefaultConsumer(
		canopsis.FIFOAckConsumerName,
		canopsis.FIFOAckQueueName,
		cfg.Global.PrefetchCount,
		cfg.Global.PrefetchSize,
		false,
		"",
		"",
		"",
		"",
		amqpConnection,
		&ackMessageProcessor{
			FeaturePrintEventOnError: options.PrintEventOnError,
			Scheduler:                scheduler,
			Decoder:                  json.NewDecoder(),
			Logger:                   logger,
		},
		logger,
	))
	engine.AddPeriodicalWorker("run info", runInfoPeriodicalWorker)
	engine.AddPeriodicalWorker("outdated rates", libengine.NewLockedPeriodicalWorker(
		redis.NewLockClient(engineLockRedisClient),
		redis.FifoDeleteOutdatedRatesLockKey,
		&deleteOutdatedRatesWorker{
			PeriodicalInterval:        time.Hour,
			TimezoneConfigProvider:    timezoneConfigProvider,
			DataStorageConfigProvider: config.NewDataStorageConfigProvider(cfg, logger),
			LimitConfigAdapter:        datastorage.NewAdapter(mongoClient),
			RateLimitAdapter:          ratelimit.NewAdapter(mongoClient),
			Logger:                    logger,
		},
		logger,
	))
	engine.AddPeriodicalWorker("config", libengine.NewLoadConfigPeriodicalWorker(
		canopsis.PeriodicalWaitTime,
		config.NewAdapter(mongoClient),
		timezoneConfigProvider,
		logger,
	))
	if mongoClient.IsDistributed() {
		engine.AddRoutine(func(ctx context.Context) error {
			w := eventfilter.NewRulesChangesWatcher(mongoClient, eventfilterService, logger)

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
		engine.AddPeriodicalWorker("local cache", &periodicalWorker{
			RuleService:        eventfilterService,
			PeriodicalInterval: options.PeriodicalWaitTime,
			Logger:             logger,
		})
	}

	return engine
}
