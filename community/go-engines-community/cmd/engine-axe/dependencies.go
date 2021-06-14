package main

import (
	"fmt"
	"time"

	"git.canopsis.net/canopsis/go-engines/lib/api/watcherweather"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/config"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/encoding/json"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/engine"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/entity"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/idlealarm"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/idlerule"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/metaalarm"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/operation"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/operation/executor"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/statsng"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/depmake"
	"git.canopsis.net/canopsis/go-engines/lib/redis"
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
func NewEngineAXE(options Options, logger zerolog.Logger) engine.Engine {
	defer depmake.Catch(logger)

	m := DependencyMaker{}
	cfg := m.DepConfig()
	amqpConnection := m.DepAmqpConnection(logger, cfg)
	amqpChannel, err := amqpConnection.Channel()
	if err != nil {
		panic(err)
	}

	alarmBaggotTime, err := time.ParseDuration(cfg.Alarm.BaggotTime)
	if err != nil {
		panic(fmt.Errorf("error parsing baggot duration: %v", err))
	}

	alarmCancelAutosolveDelay, err := time.ParseDuration(cfg.Alarm.CancelAutosolveDelay)
	if err != nil || alarmCancelAutosolveDelay <= 0 {
		logger.Error().Err(err).Msg("error parsing CancelAutosolveDelay duration of alarm config section")
		alarmCancelAutosolveDelay = canopsis.CancelAutosolveDelay * time.Second
	}

	channelPub, err := amqpConnection.Channel()
	if err != nil {
		panic(fmt.Errorf("dependency error: amqp publish channel: %v", err))
	}

	dbClient := m.DepMongoClient(cfg)

	watcherRpcClient := engine.NewRPCClient(
		canopsis.AxeRPCConsumerName,
		canopsis.WatcherRPCQueueServerName,
		canopsis.AxeWatcherRPCClientQueueName,
		cfg.Global.PrefetchCount,
		cfg.Global.PrefetchSize,
		&rpcWatcherClientMessageProcessor{
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
			WatcherRpc:               watcherRpcClient,
			Executor:                 m.depOperationExecutor(cfg, logger),
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
		options.PublishToQueue,
		canopsis.FIFOAckQueueName,
		amqpConnection,
		&messageProcessor{
			FeaturePrintEventOnError: options.FeaturePrintEventOnError,
			FeatureStatEvents:        options.FeatureStatEvents,
			EventProcessor: alarm.NewEventProcessor(
				alarm.NewAdapter(dbClient),
				metaalarm.NewRuleAdapter(metaalarm.DefaultRulesCollection(m.DepMongoSession())),
				cfg,
				m.depOperationExecutor(cfg, logger),
				redis.NewLockClient(m.DepRedisSession(redis.CorrelationLockStorage, logger, cfg)),
				logger,
			),
			StatsService: m.getDefaultStatsService(logger, cfg),
			Encoder:      json.NewEncoder(),
			Decoder:      json.NewDecoder(),
			Logger:       logger,
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
			WatcherRpc:               watcherRpcClient,
			PbhRpc:                   pbhRpcClient,
			AlarmAdapter:             alarm.NewAdapter(m.DepMongoClient(cfg)),
			Executor:                 m.depOperationExecutor(cfg, logger),
			Encoder:                  json.NewEncoder(),
			Decoder:                  json.NewDecoder(),
			Logger:                   logger,
		},
		logger,
	))
	engineAxe.AddConsumer(watcherRpcClient)
	engineAxe.AddConsumer(pbhRpcClient)
	engineAxe.AddPeriodicalWorker(engine.NewRunInfoPeriodicalWorker(
		options.PeriodicalWaitTime,
		engine.NewRunInfoManager(m.DepRedisSession(redis.EngineRunInfo, logger, cfg)),
		engine.RunInfo{
			Name:         canopsis.AxeExchangeName,
			ConsumeQueue: canopsis.AxeQueueName,
			PublishQueue: options.PublishToQueue,
		},
		logger,
	))
	engineAxe.AddPeriodicalWorker(&periodicalWorker{
		PeriodicalInterval:        options.PeriodicalWaitTime,
		LockerClient:              m.getRedisLockerClient(logger, cfg),
		ChannelPub:                channelPub,
		AlarmService:              m.getDefaultAlarmService(logger, cfg),
		AlarmBaggotTime:           alarmBaggotTime,
		AlarmCancelAutosolveDelay: alarmCancelAutosolveDelay,
		Encoder:                   json.NewEncoder(),
		IdleAlarmService:          m.getIdleAlarmService(logger, cfg),
		// ignore Pbh state to resolve snoozed with Action alarm while is True
		DisableActionSnoozeDelayOnPbh: cfg.Alarm.DisableActionSnoozeDelayOnPbh,
		Logger:                        logger,
	})

	return engineAxe
}

// DependencyMaker can be created with DependencyMaker{}
type DependencyMaker struct {
	depmake.DependencyMaker
}

func (m DependencyMaker) getRedisLockerClient(logger zerolog.Logger, cfg config.CanopsisConf) redis.LockClient {
	return redis.NewLockClient(m.DepRedisSession(redis.AxePeriodicalLockStorage, logger, cfg))
}

func (m DependencyMaker) getDefaultAlarmService(logger zerolog.Logger, cfg config.CanopsisConf) alarm.Service {
	client := m.DepMongoClient(cfg)
	alarmAdapter := alarm.NewAdapter(client)
	return alarm.NewService(
		alarmAdapter,
		logger,
		cfg,
	)
}

func (m DependencyMaker) getDefaultStatsStore(cfg config.CanopsisConf) watcherweather.StatsStore {
	location, err := cfg.Timezone.GetLocation()
	if err != nil {
		panic(err)
	}
	dbClient := m.DepMongoClient(cfg)
	return watcherweather.NewStatsStore(dbClient, location)
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

func (m DependencyMaker) getIdleAlarmService(logger zerolog.Logger, cfg config.CanopsisConf) idlealarm.Service {
	client := m.DepMongoClient(cfg)

	service := idlealarm.NewService(
		idlerule.NewRuleAdapter(client),
		alarm.NewAdapter(client),
		json.NewEncoder(),
		logger,
	)

	return service
}

func (m DependencyMaker) depOperationExecutor(cfg config.CanopsisConf, logger zerolog.Logger) operation.Executor {
	entityAdapter := entity.NewAdapter(entity.DefaultCollection(m.DepMongoSession()))
	statsService := m.getDefaultStatsService(logger, cfg)

	container := operation.NewExecutorContainer()
	container.Set(types.EventTypeAck, executor.NewAckExecutor(cfg))
	container.Set(types.EventTypeAckremove, executor.NewAckRemoveExecutor(cfg))
	container.Set(types.EventTypeActivate, executor.NewActivateExecutor())
	container.Set(types.EventTypeAssocTicket, executor.NewAssocTicketExecutor())
	container.Set(types.EventTypeCancel, executor.NewCancelExecutor(cfg))
	container.Set(types.EventTypeChangestate, executor.NewChangeStateExecutor(cfg))
	container.Set(types.EventTypeComment, executor.NewCommentExecutor(cfg))
	container.Set(types.EventTypeDeclareTicket, executor.NewDeclareTicketExecutor())
	container.Set(types.EventTypeDeclareTicketWebhook, executor.NewDeclareTicketWebhookExecutor(cfg))
	container.Set(types.EventTypeDone, executor.NewDoneExecutor(cfg))
	container.Set(types.EventTypeKeepstate, executor.NewChangeStateExecutor(cfg))
	container.Set(types.EventTypePbhEnter, executor.NewPbhEnterExecutor(cfg))
	container.Set(types.EventTypePbhLeave, executor.NewPbhLeaveExecutor(cfg))
	container.Set(types.EventTypePbhLeaveAndEnter, executor.NewPbhLeaveAndEnterExecutor(cfg))
	container.Set(types.EventTypeResolveDone, executor.NewResolveStatExecutor(executor.NewResolveDoneExecutor(), entityAdapter, statsService))
	container.Set(types.EventTypeResolveCancel, executor.NewResolveStatExecutor(executor.NewResolveCancelExecutor(), entityAdapter, statsService))
	container.Set(types.EventTypeResolveClose, executor.NewResolveStatExecutor(executor.NewResolveCloseExecutor(), entityAdapter, statsService))
	container.Set(types.EventTypeSnooze, executor.NewSnoozeExecutor(cfg))
	container.Set(types.EventTypeUncancel, executor.NewUncancelExecutor(cfg))
	container.Set(types.EventTypeUnsnooze, executor.NewUnsnoozeExecutor())
	container.Set(types.EventTypeUpdateStatus, executor.NewUpdateStatusExecutor(cfg))
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

	return executor.NewMongoUpdateExecutor(
		executor.NewCombinedExecutor(container),
		alarm.NewAdapter(m.DepMongoClient(cfg)),
	)
}
