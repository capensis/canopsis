package main

import (
	"context"
	"fmt"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	libcontext "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	libengine "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/depmake"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"github.com/bsm/redislock"
	"github.com/rs/zerolog"
)

type Options struct {
	FeatureEventProcessing bool
	FeatureContextCreation bool
	FeatureContextEnrich   bool
	Purge                  bool
	PrintEventOnError      bool
	ModeDebug              bool
	ConsumeFromQueue       string
	PublishToQueue         string
	EnrichExclude          string
	EnrichInclude          string
	DataSourceDirectory    string
	PeriodicalWaitTime     time.Duration
	FifoAckExchange        string
}

type DependencyMaker struct {
	depmake.DependencyMaker
}

func NewEngineCHE(ctx context.Context, options Options, logger zerolog.Logger) libengine.Engine {
	defer depmake.Catch(logger)

	m := DependencyMaker{}
	mongoClient := m.DepMongoClient(ctx)
	cfg := m.DepConfig(ctx, mongoClient)
	config.SetDbClientRetry(mongoClient, cfg)
	alarmConfigProvider := config.NewAlarmConfigProvider(cfg, logger)
	amqpConnection := m.DepAmqpConnection(logger, cfg)
	amqpChannel := m.DepAMQPChannelPub(amqpConnection)
	entityAdapter := entity.NewAdapter(mongoClient)
	eventFilterAdapter := eventfilter.NewAdapter(mongoClient)
	entityServiceAdapter := entityservice.NewAdapter(mongoClient)
	redisSession := m.DepRedisSession(ctx, redis.EngineLockStorage, logger, cfg)
	runInfoRedisSession := m.DepRedisSession(ctx, redis.EngineRunInfo, logger, cfg)
	serviceRedisSession := m.DepRedisSession(ctx, redis.EntityServiceStorage, logger, cfg)
	periodicalLockClient := redis.NewLockClient(redisSession)

	eventFilterService := eventfilter.NewService(eventFilterAdapter, logger)
	enrichmentCenter := libcontext.NewEnrichmentCenter(
		entityAdapter,
		options.FeatureContextEnrich,
		entityservice.NewManager(
			entityServiceAdapter,
			entityAdapter,
			entityservice.NewStorage(serviceRedisSession, json.NewEncoder(), json.NewDecoder(), logger),
			logger,
		),
		logger,
	)
	enrichFields := libcontext.NewEnrichFields(options.EnrichInclude, options.EnrichExclude)

	engine := libengine.New(
		func(ctx context.Context) error {
			logger.Debug().Msg("Loading event filter data sources")
			err := eventFilterService.LoadDataSourceFactories(
				enrichmentCenter,
				enrichFields,
				options.DataSourceDirectory,
			)
			if err != nil {
				return fmt.Errorf("unable to load data sources: %v", err)
			}

			logger.Debug().Msg("Loading event filter rules")
			err = eventFilterService.LoadRules(ctx)
			if err != nil {
				return fmt.Errorf("unable to load rules: %v", err)
			}

			logger.Debug().Msg("Loading services")
			err = enrichmentCenter.LoadServices(ctx)
			if err != nil {
				logger.Error().Err(err).Msg("unable to load services")
			}

			_, err = periodicalLockClient.Obtain(ctx, redis.ChePeriodicalLockKey,
				options.PeriodicalWaitTime, &redislock.Options{
					RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(1*time.Second), 1),
				})
			if err != nil {
				// Lock is set for options.PeriodicalWaitTime TTL, other instances should skip actions below
				if err == redislock.ErrNotObtained {
					return nil
				}

				logger.Error().Err(err).Msg("cannot obtain lock")
				return err
			}
			// Below are actions locked with ChePeriodicalLockKey for multi-instance configuration

			logger.Debug().Msg("Recompute impacted services for connectors")
			err = enrichmentCenter.UpdateImpactedServices(ctx)
			if err != nil {
				logger.Warn().Err(err).Msg("error while recomputing impacted services for connectors")
			}
			logger.Debug().Msg("Recompute impacted services for connectors finished")
			return nil
		},
		func(ctx context.Context) {
			err := mongoClient.Disconnect(ctx)
			if err != nil {
				logger.Error().Err(err).Msg("failed to close mongo connection")
			}

			err = amqpConnection.Close()
			if err != nil {
				logger.Error().Err(err).Msg("failed to close amqp connection")
			}

			err = redisSession.Close()
			if err != nil {
				logger.Error().Err(err).Msg("failed to close redis connection")
			}

			err = serviceRedisSession.Close()
			if err != nil {
				logger.Error().Err(err).Msg("failed to close redis connection")
			}

			err = runInfoRedisSession.Close()
			if err != nil {
				logger.Error().Err(err).Msg("failed to close redis connection")
			}
		},
		logger,
	)
	engine.AddConsumer(libengine.NewDefaultConsumer(
		canopsis.CheConsumerName,
		options.ConsumeFromQueue,
		cfg.Global.PrefetchCount,
		cfg.Global.PrefetchSize,
		options.Purge,
		"",
		options.PublishToQueue,
		options.FifoAckExchange,
		canopsis.FIFOAckQueueName,
		amqpConnection,
		&messageProcessor{
			FeaturePrintEventOnError: false,
			FeatureEventProcessing:   options.FeatureEventProcessing,
			FeatureContextCreation:   options.FeatureContextCreation,
			AlarmConfigProvider:      alarmConfigProvider,
			EventFilterService:       eventFilterService,
			EnrichmentCenter:         enrichmentCenter,
			EnrichFields:             enrichFields,
			AmqpPublisher:            m.DepAMQPChannelPub(amqpConnection),
			AlarmAdapter:             alarm.NewAdapter(mongoClient),
			Encoder:                  json.NewEncoder(),
			Decoder:                  json.NewDecoder(),
			Logger:                   logger,
		},
		logger,
	))
	engine.AddPeriodicalWorker(&reloadLocalCachePeriodicalWorker{
		EventFilterService: eventFilterService,
		EnrichmentCenter:   enrichmentCenter,
		PeriodicalInterval: options.PeriodicalWaitTime,
		Logger:             logger,
	})
	engine.AddPeriodicalWorker(libengine.NewRunInfoPeriodicalWorker(
		options.PeriodicalWaitTime,
		libengine.NewRunInfoManager(runInfoRedisSession),
		libengine.NewInstanceRunInfo(canopsis.CheEngineName, options.ConsumeFromQueue, options.PublishToQueue),
		amqpChannel,
		logger,
	))
	engine.AddPeriodicalWorker(libengine.NewLoadConfigPeriodicalWorker(
		options.PeriodicalWaitTime,
		config.NewAdapter(mongoClient),
		alarmConfigProvider,
		logger,
	))
	engine.AddPeriodicalWorker(libengine.NewLockedPeriodicalWorker(
		periodicalLockClient,
		redis.ChePeriodicalLockKey,
		&impactedServicesPeriodicalWorker{
			EnrichmentCenter:   enrichmentCenter,
			PeriodicalInterval: options.PeriodicalWaitTime,
			Logger:             logger,
		},
		logger,
	))

	return engine
}
