package main

import (
	"context"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/serviceweather"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/idlealarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/idlerule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation/executor"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/statsng"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/depmake"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"github.com/rs/zerolog"
)

type Options struct {
	FeatureHideResources     bool
	FeaturePrintEventOnError bool
	FeatureStatEvents        bool
	ModeDebug                bool
	PublishToQueue           string
	PostProcessorsDirectory  string
	IgnoreDefaultTomlConfig  bool
	PeriodicalWaitTime       time.Duration
	WithRemediation          bool
}

// NewEngineAXE returns the default AXE engine with default connections.
func NewEngineAXE(ctx context.Context, options Options, logger zerolog.Logger) engine.Engine {
	defer depmake.Catch(logger)

	m := DependencyMaker{}
	dbClient := m.DepMongoClient(ctx)
	cfg := m.DepConfig(ctx, dbClient)
	config.SetDbClientRetry(dbClient, cfg)
	alarmConfigProvider := config.NewAlarmConfigProvider(cfg, logger)
	timezoneConfigProvider := config.NewTimezoneConfigProvider(cfg, logger)
	amqpConnection := m.DepAmqpConnection(logger, cfg)
	amqpChannel, err := amqpConnection.Channel()
	if err != nil {
		panic(err)
	}

	channelPub, err := amqpConnection.Channel()
	if err != nil {
		panic(fmt.Errorf("dependency error: amqp publish channel: %v", err))
	}

	lockRedisClient := m.DepRedisSession(ctx, redis.EngineLockStorage, logger, cfg)
	corrRedisClient := m.DepRedisSession(ctx, redis.CorrelationLockStorage, logger, cfg)
	pbhRedisClient := m.DepRedisSession(ctx, redis.PBehaviorLockStorage, logger, cfg)
	runInfoRedisClient := m.DepRedisSession(ctx, redis.EngineRunInfo, logger, cfg)

	serviceRpcClient := engine.NewRPCClient(
		canopsis.AxeRPCConsumerName,
		canopsis.ServiceRPCQueueServerName,
		canopsis.AxeServiceRPCClientQueueName,
		cfg.Global.PrefetchCount,
		cfg.Global.PrefetchSize,
		&rpcServiceClientMessageProcessor{
			Logger: logger,
		},
		amqpChannel,
		logger,
	)

	statsService := statsng.NewService(
		m.DepAMQPChannelPub(amqpConnection),
		canopsis.StatsngExchangeName,
		canopsis.StatsngQueueName,
		json.NewEncoder(),
		serviceweather.NewStatsStore(dbClient),
		logger,
	)

	pbhRpcClient := engine.NewRPCClient(
		canopsis.AxeRPCConsumerName,
		canopsis.PBehaviorRPCQueueServerName,
		canopsis.AxePbehaviorRPCClientQueueName,
		cfg.Global.PrefetchCount,
		cfg.Global.PrefetchSize,
		&rpcPBehaviorClientMessageProcessor{
			FeaturePrintEventOnError: options.FeaturePrintEventOnError,
			PublishCh:                channelPub,
			ServiceRpc:               serviceRpcClient,
			Executor:                 m.depOperationExecutor(dbClient, alarmConfigProvider, statsService),
			Decoder:                  json.NewDecoder(),
			Encoder:                  json.NewEncoder(),
			Logger:                   logger,
		},
		amqpChannel,
		logger,
	)

	var remediationRpcClient engine.RPCClient
	if options.WithRemediation {
		remediationRpcClient = engine.NewRPCClient(
			canopsis.AxeRPCConsumerName,
			canopsis.RemediationRPCQueueServerName,
			"",
			cfg.Global.PrefetchCount,
			cfg.Global.PrefetchSize,
			nil,
			amqpChannel,
			logger,
		)
	}

	engineAxe := engine.New(
		nil,
		func(ctx context.Context) {
			err := dbClient.Disconnect(ctx)
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

			err = corrRedisClient.Close()
			if err != nil {
				logger.Error().Err(err).Msg("failed to close redis connection")
			}

			err = pbhRedisClient.Close()
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
	engineAxe.AddConsumer(engine.NewDefaultConsumer(
		canopsis.AxeConsumerName,
		canopsis.AxeQueueName,
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
			FeatureStatEvents:        options.FeatureStatEvents,
			EventProcessor: alarm.NewEventProcessor(
				alarm.NewAdapter(dbClient),
				entity.NewAdapter(dbClient),
				correlation.NewRuleAdapter(dbClient),
				alarmConfigProvider,
				m.depOperationExecutor(dbClient, alarmConfigProvider, statsService),
				redis.NewLockClient(corrRedisClient),
				logger,
			),
			StatsService:           statsService,
			RemediationRpcClient:   remediationRpcClient,
			TimezoneConfigProvider: timezoneConfigProvider,
			Encoder:                json.NewEncoder(),
			Decoder:                json.NewDecoder(),
			Logger:                 logger,
			PbehaviorAdapter:       pbehavior.NewAdapter(dbClient),
			MetricsSender: metrics.NewSender(
				canopsis.MetricsExchangeName,
				canopsis.MetricsQueueName,
				channelPub,
				json.NewEncoder(),
				logger,
			),
		},
		logger,
	))
	engineAxe.AddConsumer(engine.NewRPCServer(
		canopsis.AxeRPCConsumerName,
		canopsis.AxeRPCQueueServerName,
		cfg.Global.PrefetchCount,
		cfg.Global.PrefetchSize,
		amqpConnection,
		&rpcMessageProcessor{
			FeaturePrintEventOnError: options.FeaturePrintEventOnError,
			ServiceRpc:               serviceRpcClient,
			PbhRpc:                   pbhRpcClient,
			AlarmAdapter:             alarm.NewAdapter(dbClient),
			Executor:                 m.depOperationExecutor(dbClient, alarmConfigProvider, statsService),
			Encoder:                  json.NewEncoder(),
			Decoder:                  json.NewDecoder(),
			Logger:                   logger,
			MetricSender: metrics.NewSender(
				canopsis.MetricsExchangeName,
				canopsis.MetricsQueueName,
				channelPub,
				json.NewEncoder(),
				logger,
			),
		},
		logger,
	))
	engineAxe.AddConsumer(serviceRpcClient)
	engineAxe.AddConsumer(pbhRpcClient)
	engineAxe.AddPeriodicalWorker(engine.NewRunInfoPeriodicalWorker(
		options.PeriodicalWaitTime,
		engine.NewRunInfoManager(runInfoRedisClient),
		engine.RunInfo{
			Name:         canopsis.AxeExchangeName,
			ConsumeQueue: canopsis.AxeQueueName,
			PublishQueue: options.PublishToQueue,
		},
		logger,
	))
	engineAxe.AddPeriodicalWorker(engine.NewLockedPeriodicalWorker(
		redis.NewLockClient(lockRedisClient),
		redis.AxePeriodicalLockKey,
		&periodicalWorker{
			PeriodicalInterval: options.PeriodicalWaitTime,
			ChannelPub:         channelPub,
			AlarmService:       alarm.NewService(alarm.NewAdapter(dbClient), logger),
			AlarmAdapter:       alarm.NewAdapter(dbClient),
			Encoder:            json.NewEncoder(),
			IdleAlarmService: idlealarm.NewService(
				idlerule.NewRuleAdapter(dbClient),
				alarm.NewAdapter(dbClient),
				entity.NewAdapter(dbClient),
				redis.NewStore(pbhRedisClient, "pbehaviors", 0),
				pbehavior.NewService(pbehavior.NewModelProvider(dbClient), pbehavior.NewEntityMatcher(dbClient), logger),
				json.NewEncoder(),
				logger,
			),
			AlarmConfigProvider: alarmConfigProvider,
			Logger:              logger,
		},
		logger,
	))
	engineAxe.AddPeriodicalWorker(engine.NewLockedPeriodicalWorker(
		redis.NewLockClient(lockRedisClient),
		redis.AxeResolvedArchiverPeriodicalLockKey,
		&resolvedArchiverWorker{
			PeriodicalInterval:        time.Hour,
			TimezoneConfigProvider:    timezoneConfigProvider,
			DataStorageConfigProvider: config.NewDataStorageConfigProvider(cfg, logger),
			LimitConfigAdapter:        datastorage.NewAdapter(dbClient),
			AlarmAdapter:              alarm.NewAdapter(dbClient),
			Logger:                    logger,
		},
		logger,
	))
	engineAxe.AddPeriodicalWorker(engine.NewLoadConfigPeriodicalWorker(
		options.PeriodicalWaitTime,
		config.NewAdapter(dbClient),
		alarmConfigProvider,
		logger,
	))
	engineAxe.AddPeriodicalWorker(engine.NewLoadConfigPeriodicalWorker(
		options.PeriodicalWaitTime,
		config.NewAdapter(dbClient),
		timezoneConfigProvider,
		logger,
	))

	return engineAxe
}

// DependencyMaker can be created with DependencyMaker{}
type DependencyMaker struct {
	depmake.DependencyMaker
}

func (m DependencyMaker) depOperationExecutor(
	dbClient mongo.DbClient,
	configProvider config.AlarmConfigProvider,
	statsService statsng.Service,
) operation.Executor {
	entityAdapter := entity.NewAdapter(dbClient)
	container := operation.NewExecutorContainer()
	container.Set(types.EventTypeAck, executor.NewAckExecutor(configProvider))
	container.Set(types.EventTypeAckremove, executor.NewAckRemoveExecutor(configProvider))
	container.Set(types.EventTypeActivate, executor.NewActivateExecutor())
	container.Set(types.EventTypeAssocTicket, executor.NewAssocTicketExecutor())
	container.Set(types.EventTypeCancel, executor.NewCancelExecutor(configProvider))
	container.Set(types.EventTypeChangestate, executor.NewChangeStateExecutor(configProvider))
	container.Set(types.EventTypeComment, executor.NewCommentExecutor(configProvider))
	container.Set(types.EventTypeDeclareTicket, executor.NewDeclareTicketExecutor())
	container.Set(types.EventTypeDeclareTicketWebhook, executor.NewDeclareTicketWebhookExecutor(configProvider))
	container.Set(types.EventTypeDone, executor.NewDoneExecutor(configProvider))
	container.Set(types.EventTypeKeepstate, executor.NewChangeStateExecutor(configProvider))
	container.Set(types.EventTypePbhEnter, executor.NewPbhEnterExecutor(configProvider))
	container.Set(types.EventTypePbhLeave, executor.NewPbhLeaveExecutor(configProvider))
	container.Set(types.EventTypePbhLeaveAndEnter, executor.NewPbhLeaveAndEnterExecutor(configProvider))
	container.Set(types.EventTypeResolveDone, executor.NewResolveStatExecutor(executor.NewResolveDoneExecutor(), entityAdapter, statsService))
	container.Set(types.EventTypeResolveCancel, executor.NewResolveStatExecutor(executor.NewResolveCancelExecutor(), entityAdapter, statsService))
	container.Set(types.EventTypeResolveClose, executor.NewResolveStatExecutor(executor.NewResolveCloseExecutor(), entityAdapter, statsService))
	container.Set(types.EventTypeEntityToggled, executor.NewResolveStatExecutor(executor.NewResolveDisabledExecutor(), entityAdapter, statsService))
	container.Set(types.EventTypeSnooze, executor.NewSnoozeExecutor(configProvider))
	container.Set(types.EventTypeUncancel, executor.NewUncancelExecutor(configProvider))
	container.Set(types.EventTypeUnsnooze, executor.NewUnsnoozeExecutor())
	container.Set(types.EventTypeUpdateStatus, executor.NewUpdateStatusExecutor(configProvider))
	container.Set(types.EventTypeInstructionStarted, executor.NewInstructionExecutor())
	container.Set(types.EventTypeInstructionPaused, executor.NewInstructionExecutor())
	container.Set(types.EventTypeInstructionResumed, executor.NewInstructionExecutor())
	container.Set(types.EventTypeInstructionCompleted, executor.NewInstructionExecutor())
	container.Set(types.EventTypeInstructionAborted, executor.NewInstructionExecutor())
	container.Set(types.EventTypeInstructionFailed, executor.NewInstructionExecutor())
	container.Set(types.EventTypeAutoInstructionStarted, executor.NewInstructionExecutor())
	container.Set(types.EventTypeAutoInstructionCompleted, executor.NewInstructionExecutor())
	container.Set(types.EventTypeAutoInstructionFailed, executor.NewInstructionExecutor())
	container.Set(types.EventTypeAutoInstructionAlreadyRunning, executor.NewInstructionExecutor())
	container.Set(types.EventTypeInstructionJobStarted, executor.NewInstructionExecutor())
	container.Set(types.EventTypeInstructionJobCompleted, executor.NewInstructionExecutor())
	container.Set(types.EventTypeInstructionJobAborted, executor.NewInstructionExecutor())
	container.Set(types.EventTypeInstructionJobFailed, executor.NewInstructionExecutor())
	container.Set(types.EventTypeJunitTestSuiteUpdated, executor.NewJunitExecutor())
	container.Set(types.EventTypeJunitTestCaseUpdated, executor.NewJunitExecutor())

	return executor.NewMongoUpdateExecutor(
		executor.NewCombinedExecutor(container),
		alarm.NewAdapter(dbClient),
	)
}
