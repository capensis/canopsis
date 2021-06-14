package main

import (
	"context"
	"runtime/trace"
	"time"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/encoding/json"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/engine"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/entity"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/watcher"
	"git.canopsis.net/canopsis/go-engines/lib/depmake"
	"git.canopsis.net/canopsis/go-engines/lib/mongo/bulk"
	"git.canopsis.net/canopsis/go-engines/lib/redis"
	"github.com/rs/zerolog"
)

type Options struct {
	FeaturePrintEventOnError bool
	ModeDebug                bool
	PublishToQueue           string
	PeriodicalWaitTime       time.Duration
	AutoRecomputeWatchers    bool
}

// DependencyMaker can be created with DependencyMaker{}
type DependencyMaker struct {
	depmake.DependencyMaker
}

// NewEngineWatcher returns the default Watcher engine with default connections.
func NewEngineWatcher(options Options, logger zerolog.Logger) engine.Engine {
	m := DependencyMaker{}
	cfg := m.DepConfig()
	amqpConnection := m.DepAmqpConnection(logger, cfg)
	amqpChannel, err := amqpConnection.Channel()
	if err != nil {
		panic(err)
	}
	mongoSession := m.DepMongoSession()

	watcherCollection := watcher.DefaultCollection(mongoSession)
	watcherBulk := watcherCollection.NewBulk(bulk.BulkSizeMax)
	alarmAdapter := alarm.NewAdapter(m.DepMongoClient(cfg))

	watcherAdapter := watcher.NewAdapter(
		watcherCollection,
		watcherBulk,
		entity.EntityCollectionName,
		alarm.AlarmCollectionName,
		logger)

	redisSession := m.DepRedisSession(redis.CacheWatcher, logger, cfg)
	countersCache := watcher.NewCountersCache(redisSession, logger)

	watcherService := watcher.NewService(
		redisSession,
		amqpChannel,
		canopsis.CheExchangeName,
		canopsis.FIFOQueueName,
		json.NewEncoder(),
		watcherAdapter,
		alarmAdapter,
		countersCache,
		logger,
	)

	engineWatcher := engine.New(
		func(ctx context.Context) error {
			ctx, task := trace.NewTask(context.Background(), "watcher.Initialize")
			defer task.End()

			logger.Info().Msg("started to recompute watchers")
			err := watcherService.FlushDB()
			if err != nil {
				logger.Error().Err(err).Msg("error while recomputing watchers")
				return err
			}

			err = watcherService.ComputeAllWatchers(ctx)
			if err != nil {
				logger.Error().Err(err).Msg("error while recomputing watchers")
				return err
			}

			logger.Info().Msg("recomputed watchers")

			return nil
		},
		nil,
		logger,
	)
	engineWatcher.AddConsumer(engine.NewDefaultConsumer(
		canopsis.WatcherConsumerName,
		canopsis.WatcherQueueName,
		cfg.Global.PrefetchCount,
		cfg.Global.PrefetchSize,
		options.PublishToQueue,
		canopsis.FIFOAckQueueName,
		amqpConnection,
		&messageProcessor{
			FeaturePrintEventOnError: options.FeaturePrintEventOnError,
			WatcherService:           watcherService,
			Encoder:                  json.NewEncoder(),
			Decoder:                  json.NewDecoder(),
			Logger:                   logger,
		},
		logger,
	))
	engineWatcher.AddConsumer(engine.NewRPCServer(
		canopsis.WatcherRPCConsumerName,
		canopsis.WatcherRPCQueueServerName,
		cfg.Global.PrefetchCount,
		cfg.Global.PrefetchSize,
		amqpConnection,
		&rpcServerMessageProcessor{
			FeaturePrintEventOnError: options.FeaturePrintEventOnError,
			WatcherService:           watcherService,
			Encoder:                  json.NewEncoder(),
			Decoder:                  json.NewDecoder(),
			Logger:                   logger,
		},
		logger,
	))
	engineWatcher.AddPeriodicalWorker(engine.NewRunInfoPeriodicalWorker(
		options.PeriodicalWaitTime,
		engine.NewRunInfoManager(m.DepRedisSession(redis.EngineRunInfo, logger, cfg)),
		engine.RunInfo{
			Name:         canopsis.WatcherEngineName,
			ConsumeQueue: canopsis.WatcherQueueName,
			PublishQueue: options.PublishToQueue,
		},
		logger,
	))
	if options.AutoRecomputeWatchers {
		engineWatcher.AddPeriodicalWorker(&periodicalWorker{
			WatcherService:     watcherService,
			PeriodicalInterval: options.PeriodicalWaitTime,
			Logger:             logger,
		})
	}

	return engineWatcher
}
