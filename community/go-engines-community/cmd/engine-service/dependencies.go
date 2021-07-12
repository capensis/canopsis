package main

import (
	"context"
	"runtime/trace"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/depmake"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"github.com/bsm/redislock"
	"github.com/rs/zerolog"
)

const periodicalLock = "service-periodical-lock"
const periodicalIdleSinceLock = "service-periodical-idle-since-lock"

type Options struct {
	FeaturePrintEventOnError bool
	ModeDebug                bool
	PublishToQueue           string
	PeriodicalWaitTime       time.Duration
	AutoRecomputeAll         bool
	RecomputeAllOnInit       bool
}

// DependencyMaker can be created with DependencyMaker{}
type DependencyMaker struct {
	depmake.DependencyMaker
}

// NewEngine returns the default Service engine with default connections.
func NewEngine(ctx context.Context, options Options, logger zerolog.Logger) engine.Engine {
	m := DependencyMaker{}
	mongoClient := m.DepMongoClient(ctx)
	cfg := m.DepConfig(ctx, mongoClient)
	config.SetDbClientRetry(mongoClient, cfg)
	amqpConnection := m.DepAmqpConnection(logger, cfg)
	amqpChannel, err := amqpConnection.Channel()
	if err != nil {
		panic(err)
	}

	redisSession := m.DepRedisSession(ctx, redis.CacheService, logger, cfg)
	runInfoRedisSession := m.DepRedisSession(ctx, redis.EngineRunInfo, logger, cfg)
	periodicalLockClient := redis.NewLockClient(redisSession)
	var serviceLockClient redis.LockClient
	if !options.AutoRecomputeAll {
		serviceLockClient = redis.NewLockClient(redisSession)
	}
	entityServicesService := entityservice.NewService(
		amqpChannel,
		canopsis.CheExchangeName,
		canopsis.FIFOQueueName,
		json.NewEncoder(),
		entityservice.NewAdapter(mongoClient),
		entity.NewAdapter(mongoClient),
		entityservice.NewCountersCache(redisSession, logger),
		entityservice.NewStorage(redisSession, json.NewEncoder(), json.NewDecoder(), logger),
		serviceLockClient,
		redisSession,
		logger,
	)

	engineService := engine.New(
		func(ctx context.Context) error {
			ctx, task := trace.NewTask(ctx, "service.Initialize")
			defer task.End()

			// Lock periodical, do not release lock to not allow another instance start periodical.
			_, err = periodicalLockClient.Obtain(ctx, periodicalIdleSinceLock,
				options.PeriodicalWaitTime, &redislock.Options{
					RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(1*time.Second), 1),
				})
			if err != nil {
				if err == redislock.ErrNotObtained {
					return nil
				}

				logger.Error().Err(err).Msg("cannot obtain lock")
				return err
			}

			logger.Info().Msg("started to recompute idle_since")
			err = entityServicesService.RecomputeIdleSince(ctx)
			if err != nil {
				logger.Error().Err(err).Msg("error while recomputing idle_since")
				return err
			}

			logger.Info().Msg("recomputed idle_since")

			if !options.RecomputeAllOnInit {
				return nil
			}

			// Lock periodical, do not release lock to not allow another instance start periodical.
			_, err := periodicalLockClient.Obtain(ctx, periodicalLock,
				options.PeriodicalWaitTime, &redislock.Options{
					RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(1*time.Second), 1),
				})
			if err != nil {
				if err == redislock.ErrNotObtained {
					return nil
				}

				logger.Error().Err(err).Msg("cannot obtain lock")
				return err
			}

			logger.Info().Msg("started to recompute entity services")
			err = entityServicesService.ClearCache(ctx)
			if err != nil {
				logger.Error().Err(err).Msg("error while recomputing entity services")
				return err
			}

			err = entityServicesService.ComputeAllServices(ctx)
			if err != nil {
				logger.Error().Err(err).Msg("error while recomputing entity services")
				return err
			}

			logger.Info().Msg("recomputed entity services")

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

			err = runInfoRedisSession.Close()
			if err != nil {
				logger.Error().Err(err).Msg("failed to close redis connection")
			}
		},
		logger,
	)
	engineService.AddConsumer(engine.NewDefaultConsumer(
		canopsis.ServiceConsumerName,
		canopsis.ServiceQueueName,
		cfg.Global.PrefetchCount,
		cfg.Global.PrefetchSize,
		false,
		"",
		options.PublishToQueue,
		canopsis.FIFOAckExchangeName,
		canopsis.FIFOAckQueueName,
		amqpConnection,
		&messageProcessor{
			FeaturePrintEventOnError: options.FeaturePrintEventOnError,
			EntityServiceService:     entityServicesService,
			Encoder:                  json.NewEncoder(),
			Decoder:                  json.NewDecoder(),
			Logger:                   logger,
		},
		logger,
	))
	engineService.AddConsumer(engine.NewRPCServer(
		canopsis.ServiceRPCConsumerName,
		canopsis.ServiceRPCQueueServerName,
		cfg.Global.PrefetchCount,
		cfg.Global.PrefetchSize,
		amqpConnection,
		&rpcServerMessageProcessor{
			FeaturePrintEventOnError: options.FeaturePrintEventOnError,
			EntityServiceService:     entityServicesService,
			Encoder:                  json.NewEncoder(),
			Decoder:                  json.NewDecoder(),
			Logger:                   logger,
		},
		logger,
	))
	engineService.AddPeriodicalWorker(engine.NewRunInfoPeriodicalWorker(
		options.PeriodicalWaitTime,
		engine.NewRunInfoManager(runInfoRedisSession),
		engine.RunInfo{
			Name:         canopsis.ServiceEngineName,
			ConsumeQueue: canopsis.ServiceQueueName,
			PublishQueue: options.PublishToQueue,
		},
		logger,
	))
	if options.AutoRecomputeAll {
		engineService.AddPeriodicalWorker(&periodicalWorker{
			LockClient:           periodicalLockClient,
			EntityServiceService: entityServicesService,
			PeriodicalInterval:   options.PeriodicalWaitTime,
			Logger:               logger,
		})
	}
	engineService.AddPeriodicalWorker(&idleSincePeriodicalWorker{
		LockClient:           periodicalLockClient,
		EntityServiceService: entityServicesService,
		PeriodicalInterval:   options.PeriodicalWaitTime,
		Logger:               logger,
	})

	return engineService
}
