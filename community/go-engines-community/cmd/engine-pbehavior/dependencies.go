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
	cfg := m.DepConfig()
	timezoneConfigProvider := config.NewTimezoneConfigProvider(cfg, logger)
	amqpConnection := m.DepAmqpConnection(logger, cfg)
	amqpChannel := m.DepAMQPChannelPub(amqpConnection)

	redisClient := m.DepRedisSession(ctx, redis.PBehaviorLockStorage, logger, cfg)
	lockerClient := redis.NewLockClient(redisClient)
	dbClient, err := mongo.NewClient(
		cfg.Global.ReconnectRetries,
		cfg.Global.GetReconnectTimeout(),
	)
	if err != nil {
		panic(err)
	}

	entityMatcher := pbehavior.NewComputedEntityMatcher(dbClient, redisClient,
		json.NewEncoder(), json.NewDecoder())
	pbhStore := pbehavior.NewStore(redisClient, json.NewEncoder(), json.NewDecoder())

	frameDuration := time.Duration(options.FrameDuration) * time.Minute
	eventManager := pbehavior.NewEventManager()
	enginePbehavior := engine.New(
		func(ctx context.Context) error {
			pbhService := pbehavior.NewService(pbehavior.NewModelProvider(dbClient), entityMatcher, pbhStore, lockerClient)

			now := time.Now().In(timezoneConfigProvider.Get().Location)
			newSpan := timespan.New(now, now.Add(frameDuration))

			count, err := pbhService.Compute(ctx, newSpan)
			if err != nil {
				return fmt.Errorf("compute pbehavior's frames failed: %w", err)
			}

			if count >= 0 {
				logger.Info().
					Time("interval_from", newSpan.From()).
					Time("interval_to", newSpan.To()).
					Int("count", count).
					Msg("pbehaviors are recomputed")
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
			PbhService:               pbehavior.NewEntityTypeResolver(pbhStore, entityMatcher),
			TimezoneConfigProvider:   timezoneConfigProvider,
			FeaturePrintEventOnError: options.FeaturePrintEventOnError,
			Encoder:                  json.NewEncoder(),
			Decoder:                  json.NewDecoder(),
			CreatePbehaviorProcessor: createPbehaviorMessageProcessor{
				FeaturePrintEventOnError: options.FeaturePrintEventOnError,
				DbClient:                 dbClient,
				PbhService:               pbehavior.NewService(pbehavior.NewModelProvider(dbClient), entityMatcher, pbhStore, lockerClient),
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
				PbhService:               pbehavior.NewService(pbehavior.NewModelProvider(dbClient), entityMatcher, pbhStore, lockerClient),
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
		PbhService:             pbehavior.NewService(pbehavior.NewModelProvider(dbClient), entityMatcher, pbhStore, lockerClient),
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
