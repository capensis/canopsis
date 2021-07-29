package main

import (
	"context"
	"fmt"
	"time"

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
			computeLock, err := lockerClient.Obtain(ctx, redis.RecomputeLockKey, redis.RecomputeLockDuration, &redislock.Options{
				RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(1*time.Second), 1),
			})

			defer func() {
				if computeLock != nil {
					err := computeLock.Release(ctx)
					if err != nil && err != redislock.ErrLockNotHeld {
						logger.Warn().Msg("failed to manually release compute-lock, the lock will be released by ttl")
					}
				}
			}()

			if err != nil {
				return fmt.Errorf("obtain redlock failed: %w", err)
			}

			pbhService := pbehavior.NewService(pbehavior.NewModelProvider(dbClient), pbehavior.NewEntityMatcher(dbClient))
			ok, err := store.Restore(ctx, pbhService)
			if err != nil {
				return fmt.Errorf("get pbehavior's frames from redis failed: %w", err)
			}

			now := time.Now().In(timezoneConfigProvider.Get().Location)
			span := pbhService.GetSpan()

			if !ok || span.To().Before(now.Add(frameDuration/2)) {
				err = pbhService.Compute(ctx, timespan.New(now, now.Add(frameDuration)))
				if err != nil {
					return fmt.Errorf("compute pbehavior's frames failed: %w", err)
				}

				err = store.Save(ctx, pbhService)
				if err != nil {
					return fmt.Errorf("save pbehavior's frames to redis failed: %w", err)
				}

				newSpan := pbhService.GetSpan()
				logger.Info().
					Time("interval from", newSpan.From()).
					Time("interval to", newSpan.To()).
					Int("count", pbhService.GetComputedPbehaviorsCount()).
					Msg("pbehaviors are recomputed")
			}

			err = computeLock.Release(ctx)
			if err != nil {
				if err == redislock.ErrLockNotHeld {
					return fmt.Errorf("the pbehavior's frames computing took more time than redlock ttl, the data might be inconsistent: %w", err)
				}

				logger.Warn().Msg("failed to manually release compute-lock, the lock will be released by ttl")
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
			PbhService:               pbehavior.NewService(pbehavior.NewModelProvider(dbClient), pbehavior.NewEntityMatcher(dbClient)),
			TimezoneConfigProvider:   timezoneConfigProvider,
			FeaturePrintEventOnError: options.FeaturePrintEventOnError,
			Encoder:                  json.NewEncoder(),
			Decoder:                  json.NewDecoder(),
			CreatePbehaviroProcessor: createPbehaviorMessageProcessor{
				FeaturePrintEventOnError: options.FeaturePrintEventOnError,
				DbClient:                 dbClient,
				LockerClient:             lockerClient,
				Store:                    store,
				PbhService:               pbehavior.NewService(pbehavior.NewModelProvider(dbClient), pbehavior.NewEntityMatcher(dbClient)),
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
				PbhService:               pbehavior.NewService(pbehavior.NewModelProvider(dbClient), pbehavior.NewEntityMatcher(dbClient)),
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
		PbhService:             pbehavior.NewService(pbehavior.NewModelProvider(dbClient), pbehavior.NewEntityMatcher(dbClient)),
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
