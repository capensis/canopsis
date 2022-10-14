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
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/techmetrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/depmake"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/timespan"
	"github.com/rs/zerolog"
)

type Options struct {
	FeaturePrintEventOnError bool
	ModeDebug                bool
	FrameDuration            int
	PeriodicalWaitTime       time.Duration
}

type DependencyMaker struct {
	depmake.DependencyMaker
}

func NewEnginePBehavior(ctx context.Context, options Options, logger zerolog.Logger) engine.Engine {
	m := DependencyMaker{}
	dbClient := m.DepMongoClient(ctx, logger)
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
	pbhStore := pbehavior.NewStore(pbhRedisSession, json.NewEncoder(), json.NewDecoder())
	pbhTypeComputer := pbehavior.NewTypeComputer(pbehavior.NewModelProvider(dbClient), json.NewDecoder())
	frameDuration := time.Duration(options.FrameDuration) * time.Minute
	eventManager := pbehavior.NewEventManager()
	runInfoPeriodicalWorker := engine.NewRunInfoPeriodicalWorker(
		options.PeriodicalWaitTime,
		engine.NewRunInfoManager(runInfoRedisSession),
		engine.NewInstanceRunInfo(canopsis.PBehaviorEngineName, "", "", []string{canopsis.PBehaviorRPCQueueServerName}),
		amqpChannel,
		logger,
	)

	techMetricsConfigProvider := config.NewTechMetricsConfigProvider(cfg, logger)
	techMetricsSender := techmetrics.NewSender(techMetricsConfigProvider, canopsis.TechMetricsFlushInterval,
		cfg.Global.ReconnectRetries, cfg.Global.GetReconnectTimeout(), logger)

	enginePbehavior := engine.New(
		func(ctx context.Context) error {
			runInfoPeriodicalWorker.Work(ctx)
			pbhService := pbehavior.NewService(dbClient, pbhTypeComputer, pbhStore, pbhLockerClient, logger)

			now := time.Now().In(timezoneConfigProvider.Get().Location)
			newSpan := timespan.New(now, now.Add(frameDuration))

			_, count, err := pbhService.Compute(ctx, newSpan)
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
	enginePbehavior.AddRoutine(func(ctx context.Context) error {
		techMetricsSender.Run(ctx)
		return nil
	})
	enginePbehavior.AddConsumer(engine.NewRPCServer(
		canopsis.PBehaviorRPCConsumerName,
		canopsis.PBehaviorRPCQueueServerName,
		cfg.Global.PrefetchCount,
		cfg.Global.PrefetchSize,
		amqpConnection,
		&rpcServerMessageProcessor{
			FeaturePrintEventOnError: options.FeaturePrintEventOnError,
			DbClient:                 dbClient,
			PbhService:               pbehavior.NewService(dbClient, pbhTypeComputer, pbhStore, pbhLockerClient, logger),
			EventManager:             pbehavior.NewEventManager(),
			TimezoneConfigProvider:   timezoneConfigProvider,
			Decoder:                  json.NewDecoder(),
			Encoder:                  json.NewEncoder(),
			Logger:                   logger,
		},
		logger,
	))
	enginePbehavior.AddPeriodicalWorker("run info", runInfoPeriodicalWorker)
	enginePbehavior.AddPeriodicalWorker("alarms", engine.NewLockedPeriodicalWorker(
		redis.NewLockClient(lockRedisSession),
		redis.PbehaviorPeriodicalLockKey,
		&periodicalWorker{
			TechMetricsSender:      techMetricsSender,
			ChannelPub:             amqpChannel,
			PeriodicalInterval:     options.PeriodicalWaitTime,
			PbhService:             pbehavior.NewService(dbClient, pbhTypeComputer, pbhStore, pbhLockerClient, logger),
			AlarmAdapter:           alarm.NewAdapter(dbClient),
			EntityAdapter:          entity.NewAdapter(dbClient),
			EventManager:           eventManager,
			FrameDuration:          frameDuration,
			Encoder:                json.NewEncoder(),
			Logger:                 logger,
			TimezoneConfigProvider: timezoneConfigProvider,
		},
		logger,
	))
	enginePbehavior.AddPeriodicalWorker("cleaner", engine.NewLockedPeriodicalWorker(
		redis.NewLockClient(lockRedisSession),
		redis.PbehaviorCleanPeriodicalLockKey,
		&cleanPeriodicalWorker{
			PeriodicalInterval:        time.Hour,
			TimezoneConfigProvider:    timezoneConfigProvider,
			DataStorageConfigProvider: dataStorageConfigProvider,
			LimitConfigAdapter:        datastorage.NewAdapter(dbClient),
			PbehaviorCleaner:          pbehavior.NewCleaner(dbClient, logger),
			Logger:                    logger,
		},
		logger,
	))
	enginePbehavior.AddPeriodicalWorker("config", engine.NewLoadConfigPeriodicalWorker(
		options.PeriodicalWaitTime,
		config.NewAdapter(dbClient),
		logger,
		timezoneConfigProvider,
		dataStorageConfigProvider,
		techMetricsConfigProvider,
	))

	return enginePbehavior
}
