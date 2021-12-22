package main

import (
	"context"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	libengine "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/neweventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/ratelimit"
	libscheduler "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/scheduler"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/statistics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/depmake"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"github.com/rs/zerolog"
	"io/ioutil"
	"path/filepath"
	"plugin"
	"strings"
	"time"
)

type Options struct {
	PrintEventOnError         bool
	ModeDebug                 bool
	ConsumeFromQueue          string
	PublishToQueue            string
	LockTtl                   int
	EnableMetaAlarmProcessing bool
	EventsStatsFlushInterval  time.Duration
	PeriodicalWaitTime        time.Duration
	DataSourceDirectory       string
}

// DependencyMaker can be created with DependencyMaker{}
type DependencyMaker struct {
	depmake.DependencyMaker
}

func NewEngine(ctx context.Context, options Options, logger zerolog.Logger) libengine.Engine {
	defer depmake.Catch(logger)

	m := DependencyMaker{}
	mongoClient := m.DepMongoClient(ctx)
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
		options.EnableMetaAlarmProcessing,
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

	logger.Debug().Msg("Loading event filter data sources")
	factories, err := LoadDataSourceFactories(options.DataSourceDirectory)
	if err != nil {
		panic(fmt.Errorf("unable to load data sources: %w", err))
	}

	ruleAdapter := neweventfilter.NewRuleAdapter(mongoClient)
	ruleApplicatorContainer := neweventfilter.NewRuleApplicatorContainer()
	ruleApplicatorContainer.Set(neweventfilter.RuleTypeChangeEntity, neweventfilter.NewChangeEntityApplicator(factories))
	eventFilterService := neweventfilter.NewRuleService(ruleAdapter, ruleApplicatorContainer, logger)

	engine := libengine.New(
		func(ctx context.Context) error {
			scheduler.Start(ctx)

			err := eventFilterService.LoadRules(ctx)
			if err != nil {
				return err
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
			EventFilterService:       eventFilterService,
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
	engine.AddPeriodicalWorker(&periodicalWorker{
		RuleService:        eventFilterService,
		PeriodicalInterval: options.PeriodicalWaitTime,
		Logger:             logger,
	})
	engine.AddPeriodicalWorker(libengine.NewRunInfoPeriodicalWorker(
		canopsis.PeriodicalWaitTime,
		libengine.NewRunInfoManager(runInfoRedisClient),
		libengine.NewInstanceRunInfo(canopsis.FIFOEngineName, options.ConsumeFromQueue, options.PublishToQueue),
		amqpChannel,
		logger,
	))
	engine.AddPeriodicalWorker(libengine.NewLockedPeriodicalWorker(
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
	engine.AddPeriodicalWorker(libengine.NewLoadConfigPeriodicalWorker(
		canopsis.PeriodicalWaitTime,
		config.NewAdapter(mongoClient),
		timezoneConfigProvider,
		logger,
	))

	return engine
}

//TODO: copy from eventfilter package, all mongo plugin feature should be refactored
func LoadDataSourceFactories(dataSourceDirectory string) (map[string]eventfilter.DataSourceFactory, error) {
	factories := make(map[string]eventfilter.DataSourceFactory)

	files, err := ioutil.ReadDir(dataSourceDirectory)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), canopsis.PluginExtension) {
			sourceName := strings.TrimSuffix(file.Name(), canopsis.PluginExtension)
			fileName := filepath.Join(dataSourceDirectory, file.Name())

			plug, err := plugin.Open(fileName)
			if err != nil {
				return nil, fmt.Errorf("unable to open plugin: %w", err)
			}

			factorySymbol, err := plug.Lookup("DataSourceFactory")
			if err != nil {
				return nil, fmt.Errorf("unable to load plugin: %w", err)
			}

			factory, isFactory := factorySymbol.(eventfilter.DataSourceFactory)
			if !isFactory {
				return nil, fmt.Errorf("the plugin does not define a valid data source")
			}

			factories[sourceName] = factory
		}
	}

	return factories, nil
}
