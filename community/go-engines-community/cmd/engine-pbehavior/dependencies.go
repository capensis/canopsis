package main

import (
	"context"
	"fmt"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/depmake"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/timespan"
	"github.com/rs/zerolog"
)

type Options struct {
	FeaturePrintEventOnError bool
	ModeDebug                bool
	ConsumeFromQueue         string
	PublishToQueue           string
	FrameDuration            int
	PeriodicalWaitTime       time.Duration
	FifoAckExchange          string
}

type DependencyMaker struct {
	depmake.DependencyMaker
}

func NewEnginePBehavior(ctx context.Context, options Options, logger zerolog.Logger) engine.Engine {
	m := DependencyMaker{}
	dbClient := m.DepMongoClient(ctx)
	cfg := m.DepConfig(ctx, dbClient)
	config.SetDbClientRetry(dbClient, cfg)
	timezoneConfigProvider := config.NewTimezoneConfigProvider(cfg, logger)
	dataStorageConfigProvider := config.NewDataStorageConfigProvider(cfg, logger)
	amqpConnection := m.DepAmqpConnection(logger, cfg)
	amqpChannel := m.DepAMQPChannelPub(amqpConnection)
	pbhRedisSession := m.DepRedisSession(ctx, redis.PBehaviorLockStorage, logger, cfg)
	runInfoRedisSession := m.DepRedisSession(ctx, redis.EngineRunInfo, logger, cfg)
	lockRedisSession := m.DepRedisSession(ctx, redis.EngineLockStorage, logger, cfg)
	pbhLockerClient := redis.NewLockClient(pbhRedisSession)

	entityMatcher := pbehavior.NewComputedEntityMatcher(dbClient, pbhRedisSession,
		json.NewEncoder(), json.NewDecoder())
	pbhStore := pbehavior.NewStore(pbhRedisSession, json.NewEncoder(), json.NewDecoder())

	frameDuration := time.Duration(options.FrameDuration) * time.Minute
	eventManager := pbehavior.NewEventManager()
	enginePbehavior := engine.New(
		func(ctx context.Context) error {
			pbhService := pbehavior.NewService(pbehavior.NewModelProvider(dbClient), entityMatcher, pbhStore, pbhLockerClient)

			now := time.Now().In(timezoneConfigProvider.Get().Location)
			newSpan := timespan.New(now, now.Add(frameDuration))

			count, err := pbhService.Compute(ctx, newSpan)
			if err != nil {
				return fmt.Errorf("compute pbehavior's frames failed: %w", err)
			}

			if count >= 0 {
				logger.Info().
					Time("interval from", newSpan.From()).
					Time("interval to", newSpan.To()).
					Int("count", count).
					Msg("pbehaviors are recomputed")
			}

			return nil
		},
		func(ctx context.Context) {
			err := dbClient.Disconnect(ctx)
			if err != nil {
				logger.Error().Err(err).Msg("failed to close mongo connection")
			}

			err = amqpConnection.Close()
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

			err = runInfoRedisSession.Close()
			if err != nil {
				logger.Error().Err(err).Msg("failed to close redis connection")
			}
		},
		logger,
	)
	enginePbehavior.AddConsumer(engine.NewDefaultConsumer(
		canopsis.PBehaviorConsumerName,
		canopsis.PBehaviorQueueName,
		cfg.Global.PrefetchCount,
		cfg.Global.PrefetchSize,
		false,
		"",
		options.PublishToQueue,
		options.FifoAckExchange,
		canopsis.FIFOAckQueueName,
		amqpConnection,
		&messageProcessor{
			PbhService:               pbehavior.NewEntityTypeResolver(pbhStore, entityMatcher),
			TimezoneConfigProvider:   timezoneConfigProvider,
			FeaturePrintEventOnError: options.FeaturePrintEventOnError,
			Encoder:                  json.NewEncoder(),
			Decoder:                  json.NewDecoder(),
			CreatePbehaviorProcessor: createPbehaviorMessageProcessor{
				FeaturePrintEventOnError: options.FeaturePrintEventOnError,
				DbClient:                 dbClient,
				PbhService:               pbehavior.NewService(pbehavior.NewModelProvider(dbClient), entityMatcher, pbhStore, pbhLockerClient),
				EventManager:             pbehavior.NewEventManager(),
				AlarmAdapter:             alarm.NewAdapter(dbClient),
				TimezoneConfigProvider:   timezoneConfigProvider,
				Logger:                   logger,
			},
			ChannelPub: amqpChannel,
			Logger:     logger,
		},
		logger,
	))
	enginePbehavior.AddConsumer(engine.NewRPCServer(
		canopsis.PBehaviorRPCConsumerName,
		canopsis.PBehaviorRPCQueueServerName,
		cfg.Global.PrefetchCount,
		cfg.Global.PrefetchSize,
		amqpConnection,
		&rpcServerMessageProcessor{
			FeaturePrintEventOnError: options.FeaturePrintEventOnError,
			Processor: createPbehaviorMessageProcessor{
				FeaturePrintEventOnError: options.FeaturePrintEventOnError,
				DbClient:                 dbClient,
				PbhService:               pbehavior.NewService(pbehavior.NewModelProvider(dbClient), entityMatcher, pbhStore, pbhLockerClient),
				EventManager:             pbehavior.NewEventManager(),
				AlarmAdapter:             alarm.NewAdapter(dbClient),
				TimezoneConfigProvider:   timezoneConfigProvider,
				Logger:                   logger,
			},
			Decoder: json.NewDecoder(),
			Encoder: json.NewEncoder(),
			Logger:  logger,
		},
		logger,
	))
	enginePbehavior.AddPeriodicalWorker(engine.NewRunInfoPeriodicalWorker(
		options.PeriodicalWaitTime,
		engine.NewRunInfoManager(runInfoRedisSession),
		engine.NewInstanceRunInfo(canopsis.PBehaviorEngineName, canopsis.PBehaviorQueueName, options.PublishToQueue),
		amqpChannel,
		logger,
	))
	enginePbehavior.AddPeriodicalWorker(engine.NewLockedPeriodicalWorker(
		redis.NewLockClient(lockRedisSession),
		redis.PbehaviorPeriodicalLockKey,
		&periodicalWorker{
			ChannelPub:             amqpChannel,
			PeriodicalInterval:     options.PeriodicalWaitTime,
			PbhService:             pbehavior.NewService(pbehavior.NewModelProvider(dbClient), entityMatcher, pbhStore, pbhLockerClient),
			DbClient:               dbClient,
			EventManager:           eventManager,
			FrameDuration:          frameDuration,
			Encoder:                json.NewEncoder(),
			Logger:                 logger,
			TimezoneConfigProvider: timezoneConfigProvider,
		},
		logger,
	))
	enginePbehavior.AddPeriodicalWorker(engine.NewLockedPeriodicalWorker(
		redis.NewLockClient(lockRedisSession),
		redis.PbehaviorCleanPeriodicalLockKey,
		&cleanPeriodicalWorker{
			PeriodicalInterval:        time.Hour,
			TimezoneConfigProvider:    timezoneConfigProvider,
			DataStorageConfigProvider: dataStorageConfigProvider,
			LimitConfigAdapter:        datastorage.NewAdapter(dbClient),
			Logger:                    logger,
		},
		logger,
	))
	enginePbehavior.AddPeriodicalWorker(engine.NewLoadConfigPeriodicalWorker(
		options.PeriodicalWaitTime,
		config.NewAdapter(dbClient),
		timezoneConfigProvider,
		logger,
	))
	enginePbehavior.AddPeriodicalWorker(engine.NewLoadConfigPeriodicalWorker(
		options.PeriodicalWaitTime,
		config.NewAdapter(dbClient),
		dataStorageConfigProvider,
		logger,
	))

	return enginePbehavior
}
