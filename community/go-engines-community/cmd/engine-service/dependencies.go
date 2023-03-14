package main

import (
	"context"
	"errors"
	"runtime/trace"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/techmetrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/depmake"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"github.com/bsm/redislock"
	"github.com/rs/zerolog"
)

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
	mongoClient := m.DepMongoClient(ctx, logger)
	cfg := m.DepConfig(ctx, mongoClient)
	config.SetDbClientRetry(mongoClient, cfg)
	templateConfigProvider := config.NewTemplateConfigProvider(cfg)
	timezoneConfigProvider := config.NewTimezoneConfigProvider(cfg, logger)
	amqpConnection := m.DepAmqpConnection(logger, cfg)
	amqpChannel := m.DepAMQPChannelPub(amqpConnection)
	redisSession := m.DepRedisSession(ctx, redis.CacheService, logger, cfg)
	runInfoRedisSession := m.DepRedisSession(ctx, redis.EngineRunInfo, logger, cfg)
	lockRedisSession := m.DepRedisSession(ctx, redis.EngineLockStorage, logger, cfg)
	periodicalLockClient := redis.NewLockClient(lockRedisSession)
	var serviceLockClient redis.LockClient
	if !options.AutoRecomputeAll {
		serviceLockClient = redis.NewLockClient(redisSession)
	}
	templateExecutor := template.NewExecutor(templateConfigProvider, timezoneConfigProvider)
	entityServicesService := entityservice.NewService(
		amqpChannel,
		canopsis.CheExchangeName,
		canopsis.FIFOQueueName,
		json.NewEncoder(),
		entityservice.NewAdapter(mongoClient),
		entity.NewAdapter(mongoClient),
		alarm.NewAdapter(mongoClient),
		entityservice.NewCountersCache(redisSession, logger),
		entityservice.NewStorage(entityservice.NewAdapter(mongoClient), redisSession, json.NewEncoder(), json.NewDecoder(), logger),
		serviceLockClient,
		redisSession,
		templateExecutor,
		logger,
	)
	runInfoPeriodicalWorker := engine.NewRunInfoPeriodicalWorker(
		options.PeriodicalWaitTime,
		engine.NewRunInfoManager(runInfoRedisSession),
		engine.NewInstanceRunInfo(canopsis.ServiceEngineName, canopsis.ServiceQueueName, options.PublishToQueue),
		amqpChannel,
		logger,
	)

	engineService := engine.New(
		func(ctx context.Context) error {
			ctx, task := trace.NewTask(ctx, "service.Initialize")
			defer task.End()
			runInfoPeriodicalWorker.Work(ctx)

			// Lock periodical, do not release lock to not allow another instance start periodical.
			_, err := periodicalLockClient.Obtain(ctx, redis.ServiceIdleSincePeriodicalLockKey,
				options.PeriodicalWaitTime, &redislock.Options{
					RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(1*time.Second), 1),
				})

			if err == nil {
				logger.Info().Msg("started to recompute idle_since")
				err = entityServicesService.RecomputeIdleSince(ctx)
				if err != nil {
					logger.Error().Err(err).Msg("error while recomputing idle_since")
					return err
				}

				logger.Info().Msg("recomputed idle_since")
			} else {
				const msgIdleSince = "cannot obtain lock to recompute idle_since"
				if !errors.Is(err, redislock.ErrNotObtained) {
					logger.Err(err).Msg(msgIdleSince)
					// fail all following actions only if error isn't ErrNotObtained
					return err
				}
				logger.Warn().Err(err).Msg(msgIdleSince)
			}

			if !options.RecomputeAllOnInit {
				return nil
			}

			// Lock periodical, do not release lock to not allow another instance start periodical.
			_, err = periodicalLockClient.Obtain(ctx, redis.ServicePeriodicalLockKey,
				options.PeriodicalWaitTime, &redislock.Options{
					RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(1*time.Second), 1),
				})
			if err == nil {
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
			} else {
				const msgEntity = "cannot obtain lock to recompute entity services"
				if !errors.Is(err, redislock.ErrNotObtained) {
					logger.Err(err).Msg(msgEntity)
					return err
				}
				logger.Warn().Err(err).Msg(msgEntity)
			}

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

	techMetricsConfigProvider := config.NewTechMetricsConfigProvider(cfg, logger)
	techMetricsSender := techmetrics.NewSender(techMetricsConfigProvider, canopsis.TechMetricsFlushInterval,
		cfg.Global.ReconnectRetries, cfg.Global.GetReconnectTimeout(), logger)

	engineService.AddRoutine(func(ctx context.Context) error {
		techMetricsSender.Run(ctx)
		return nil
	})

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
			TechMetricsSender:        techMetricsSender,
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
	engineService.AddPeriodicalWorker("run info", runInfoPeriodicalWorker)
	if options.AutoRecomputeAll {
		engineService.AddPeriodicalWorker("recompute all", engine.NewLockedPeriodicalWorker(
			periodicalLockClient,
			redis.ServicePeriodicalLockKey,
			&recomputeAllPeriodicalWorker{
				EntityServiceService: entityServicesService,
				PeriodicalInterval:   options.PeriodicalWaitTime,
				Logger:               logger,
			},
			logger,
		))
	}
	engineService.AddPeriodicalWorker("idle since", engine.NewLockedPeriodicalWorker(
		periodicalLockClient,
		redis.ServiceIdleSincePeriodicalLockKey,
		&idleSincePeriodicalWorker{
			EntityServiceService: entityServicesService,
			PeriodicalInterval:   options.PeriodicalWaitTime,
			Logger:               logger,
		},
		logger,
	))
	engineService.AddPeriodicalWorker("config", engine.NewLoadConfigPeriodicalWorker(
		options.PeriodicalWaitTime,
		config.NewAdapter(mongoClient),
		logger,
		techMetricsConfigProvider,
		timezoneConfigProvider,
		templateConfigProvider,
	))

	return engineService
}
