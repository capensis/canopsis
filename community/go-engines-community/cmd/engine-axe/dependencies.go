package main

import (
	"context"
	"fmt"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/serviceweather"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/idlealarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/idlerule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metaalarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation/executor"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/statsng"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/depmake"
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
}

// NewEngineAXE returns the default AXE engine with default connections.
func NewEngineAXE(ctx context.Context, options Options, logger zerolog.Logger) engine.Engine {
	defer depmake.Catch(logger)

	m := DependencyMaker{}
	cfg := m.DepConfig()
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

	dbClient := m.DepMongoClient(cfg)

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
			Executor:                 m.depOperationExecutor(cfg, alarmConfigProvider, logger),
			Decoder:                  json.NewDecoder(),
			Encoder:                  json.NewEncoder(),
			Logger:                   logger,
		},
		amqpChannel,
		logger,
	)

	engineAxe := engine.New(nil, nil, logger)
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
				metaalarm.NewRuleAdapter(dbClient),
				alarmConfigProvider,
				m.depOperationExecutor(cfg, alarmConfigProvider, logger),
				redis.NewLockClient(m.DepRedisSession(ctx, redis.CorrelationLockStorage, logger, cfg)),
				logger,
			),
			StatsService:           m.getDefaultStatsService(logger, cfg),
			TimezoneConfigProvider: timezoneConfigProvider,
			Encoder:                json.NewEncoder(),
			Decoder:                json.NewDecoder(),
			Logger:                 logger,
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
			AlarmAdapter:             alarm.NewAdapter(m.DepMongoClient(cfg)),
			Executor:                 m.depOperationExecutor(cfg, alarmConfigProvider, logger),
			Encoder:                  json.NewEncoder(),
			Decoder:                  json.NewDecoder(),
			Logger:                   logger,
		},
		logger,
	))
	engineAxe.AddConsumer(serviceRpcClient)
	engineAxe.AddConsumer(pbhRpcClient)
	engineAxe.AddPeriodicalWorker(engine.NewRunInfoPeriodicalWorker(
		options.PeriodicalWaitTime,
		engine.NewRunInfoManager(m.DepRedisSession(ctx, redis.EngineRunInfo, logger, cfg)),
		engine.RunInfo{
			Name:         canopsis.AxeExchangeName,
			ConsumeQueue: canopsis.AxeQueueName,
			PublishQueue: options.PublishToQueue,
		},
		logger,
	))
	engineAxe.AddPeriodicalWorker(&periodicalWorker{
		PeriodicalInterval:  options.PeriodicalWaitTime,
		LockerClient:        m.getRedisLockerClient(ctx, logger, cfg),
		ChannelPub:          channelPub,
		AlarmService:        m.getDefaultAlarmService(logger, cfg),
		Encoder:             json.NewEncoder(),
		IdleAlarmService:    m.getIdleAlarmService(ctx, logger, cfg),
		AlarmConfigProvider: alarmConfigProvider,
		Logger:              logger,
	})
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

func (m DependencyMaker) getRedisLockerClient(ctx context.Context, logger zerolog.Logger, cfg config.CanopsisConf) redis.LockClient {
	return redis.NewLockClient(m.DepRedisSession(ctx, redis.AxePeriodicalLockStorage, logger, cfg))
}

func (m DependencyMaker) getDefaultAlarmService(logger zerolog.Logger, cfg config.CanopsisConf) alarm.Service {
	client := m.DepMongoClient(cfg)
	alarmAdapter := alarm.NewAdapter(client)
	return alarm.NewService(
		alarmAdapter,
		logger,
	)
}

func (m DependencyMaker) getDefaultStatsStore(cfg config.CanopsisConf) serviceweather.StatsStore {
	dbClient := m.DepMongoClient(cfg)
	return serviceweather.NewStatsStore(dbClient)
}

func (m DependencyMaker) getDefaultStatsService(logger zerolog.Logger, cfg config.CanopsisConf) statsng.Service {
	amqpSession := m.DepAmqpConnection(logger, cfg)
	pubChannel := m.DepAMQPChannelPub(amqpSession)

	statsService := statsng.NewService(
		pubChannel,
		canopsis.StatsngExchangeName,
		canopsis.StatsngQueueName,
		json.NewEncoder(),
		m.getDefaultStatsStore(cfg),
		logger,
	)

	return statsService
}

func (m DependencyMaker) getIdleAlarmService(ctx context.Context, logger zerolog.Logger, cfg config.CanopsisConf) idlealarm.Service {
	client := m.DepMongoClient(cfg)

	service := idlealarm.NewService(
		idlerule.NewRuleAdapter(client),
		alarm.NewAdapter(client),
		entity.NewAdapter(client),
		redis.NewStore(m.DepRedisSession(ctx, redis.PBehaviorLockStorage, logger, cfg), "pbehaviors", 0),
		pbehavior.NewService(pbehavior.NewModelProvider(client), pbehavior.NewEntityMatcher(client), logger),
		json.NewEncoder(),
		logger,
	)

	return service
}

func (m DependencyMaker) depOperationExecutor(
	cfg config.CanopsisConf,
	configProvider config.AlarmConfigProvider,
	logger zerolog.Logger,
) operation.Executor {
	entityAdapter := entity.NewAdapter(m.DepMongoClient(cfg))
	statsService := m.getDefaultStatsService(logger, cfg)

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
	container.Set(types.EventTypeInstructionJobStarted, executor.NewInstructionExecutor())
	container.Set(types.EventTypeInstructionJobCompleted, executor.NewInstructionExecutor())
	container.Set(types.EventTypeInstructionJobAborted, executor.NewInstructionExecutor())
	container.Set(types.EventTypeInstructionJobFailed, executor.NewInstructionExecutor())
	container.Set(types.EventTypeJunitTestSuiteUpdated, executor.NewJunitExecutor())
	container.Set(types.EventTypeJunitTestCaseUpdated, executor.NewJunitExecutor())

	return executor.NewMongoUpdateExecutor(
		executor.NewCombinedExecutor(container),
		alarm.NewAdapter(m.DepMongoClient(cfg)),
	)
}
