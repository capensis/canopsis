package main

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/depmake"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/timespan"
	"github.com/bsm/redislock"
	"github.com/rs/zerolog"
	"time"
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

func (m DependencyMaker) getRedisLockerClient(ctx context.Context, logger zerolog.Logger, cfg config.CanopsisConf) redis.LockClient {
	return redis.NewLockClient(m.DepRedisSession(ctx, redis.PBehaviorLockStorage, logger, cfg))
}

func (m DependencyMaker) getRedisStore(ctx context.Context, logger zerolog.Logger, cfg config.CanopsisConf) redis.Store {
	store := redis.NewStore(m.DepRedisSession(ctx, redis.PBehaviorLockStorage, logger, cfg), "pbehaviors", 0)

	return store
}

func NewEnginePBehavior(ctx context.Context, options Options, logger zerolog.Logger) engine.Engine {
	m := DependencyMaker{}
	cfg := m.DepConfig()
	timezoneConfigProvider := config.NewTimezoneConfigProvider(cfg, logger)
	amqpConnection := m.DepAmqpConnection(logger, cfg)
	amqpChannel, err := amqpConnection.Channel()
	if err != nil {
		panic(err)
	}

	lockerClient := m.getRedisLockerClient(ctx, logger, cfg)
	store := m.getRedisStore(ctx, logger, cfg)

	dbClient, err := mongo.NewClient(
		cfg.Global.ReconnectRetries,
		cfg.Global.GetReconnectTimeout(),
	)
	if err != nil {
		panic(err)
	}

	frameDuration := time.Duration(options.FrameDuration) * time.Minute
	eventManager := pbehavior.NewEventManager()
	enginePbehavior := engine.New(
		func(ctx context.Context) error {
			logger.Debug().Msg("Initialize process")

			computeLock, err := lockerClient.Obtain(ctx, redis.RecomputeLockKey, redis.RecomputeLockDuration, &redislock.Options{
				RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(1*time.Second), 1),
			})

			defer func() {
				if computeLock != nil {
					err := computeLock.Release(ctx)
					if err != nil && err != redislock.ErrLockNotHeld {
						logger.Warn().Msg("Initialize: failed to manually release compute-lock, the lock will be released by ttl")
					}
				}
			}()

			if err != nil {
				logger.Err(err).Msg("Initialize: obtain redlock failed! The engine will be stopped")

				return err
			}

			pbhService := pbehavior.NewService(pbehavior.NewModelProvider(dbClient), pbehavior.NewEntityMatcher(dbClient), logger)
			ok, err := store.Restore(ctx, pbhService)
			if err != nil {
				logger.Err(err).Msg("Initialize: get pbehavior's frames from redis failed! The engine will be stopped")

				return err
			}

			now := time.Now().In(timezoneConfigProvider.Get().Location)
			span := pbhService.GetSpan()

			if !ok || span.To().Before(now.Add(frameDuration/2)) {
				err = pbhService.Compute(ctx, timespan.New(now, now.Add(frameDuration)))
				if err != nil {
					logger.Err(err).Msg("Initialize: compute pbehavior's frames failed! The engine will be stopped")

					return err
				}

				err = store.Save(ctx, pbhService)
				if err != nil {
					logger.Err(err).Msg("Initialize: save pbehavior's frames to redis failed! The engine will be stopped")

					return err
				}
			}

			err = computeLock.Release(ctx)
			if err != nil {
				if err == redislock.ErrLockNotHeld {
					logger.Err(err).Msg("Initialize: the pbehavior's frames computing took more time than redlock ttl, the data might be inconsistent, engine will be stopped")

					return err
				} else {
					logger.Warn().Msg("Initialize: failed to manually release compute-lock, the lock will be released by ttl")
				}
			}

			return nil
		},
		nil,
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
			Store:                    store,
			PbhService:               pbehavior.NewService(pbehavior.NewModelProvider(dbClient), pbehavior.NewEntityMatcher(dbClient), logger),
			TimezoneConfigProvider:   timezoneConfigProvider,
			FeaturePrintEventOnError: options.FeaturePrintEventOnError,
			Encoder:                  json.NewEncoder(),
			Decoder:                  json.NewDecoder(),
			CreatePbehaviroProcessor: createPbehaviorMessageProcessor{
				FeaturePrintEventOnError: options.FeaturePrintEventOnError,
				DbClient:                 dbClient,
				LockerClient:             lockerClient,
				Store:                    store,
				PbhService:               pbehavior.NewService(pbehavior.NewModelProvider(dbClient), pbehavior.NewEntityMatcher(dbClient), logger),
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
				LockerClient:             lockerClient,
				Store:                    store,
				PbhService:               pbehavior.NewService(pbehavior.NewModelProvider(dbClient), pbehavior.NewEntityMatcher(dbClient), logger),
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
		engine.NewRunInfoManager(m.DepRedisSession(ctx, redis.EngineRunInfo, logger, cfg)),
		engine.RunInfo{
			Name:         canopsis.PBehaviorEngineName,
			ConsumeQueue: canopsis.PBehaviorQueueName,
			PublishQueue: options.PublishToQueue,
		},
		logger,
	))
	enginePbehavior.AddPeriodicalWorker(&periodicalWorker{
		ChannelPub:             amqpChannel,
		PeriodicalInterval:     options.PeriodicalWaitTime,
		LockerClient:           lockerClient,
		Store:                  store,
		PbhService:             pbehavior.NewService(pbehavior.NewModelProvider(dbClient), pbehavior.NewEntityMatcher(dbClient), logger),
		DbClient:               dbClient,
		EventManager:           eventManager,
		FrameDuration:          frameDuration,
		Encoder:                json.NewEncoder(),
		Logger:                 logger,
		TimezoneConfigProvider: timezoneConfigProvider,
	})
	enginePbehavior.AddPeriodicalWorker(engine.NewLoadConfigPeriodicalWorker(
		options.PeriodicalWaitTime,
		config.NewAdapter(dbClient),
		timezoneConfigProvider,
		logger,
	))

	return enginePbehavior
}