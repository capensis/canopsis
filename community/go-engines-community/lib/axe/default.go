package axe

import (
	"context"
	"flag"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmstatus"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	libengine "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/flappingrule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/idlealarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/idlerule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation/executor"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/resolverule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/statistics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/techmetrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/depmake"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"github.com/rs/zerolog"
)

type Options struct {
	Version                  bool
	FeatureHideResources     bool
	FeaturePrintEventOnError bool
	ModeDebug                bool
	PublishToQueue           string
	PostProcessorsDirectory  string
	IgnoreDefaultTomlConfig  bool
	PeriodicalWaitTime       time.Duration
	WithRemediation          bool
}

func ParseOptions() Options {
	opts := Options{}

	flag.BoolVar(&opts.ModeDebug, "d", false, "debug")
	flag.BoolVar(&opts.FeaturePrintEventOnError, "printEventOnError", false, "Print event on processing error")
	flag.BoolVar(&opts.FeatureHideResources, "featureHideResources", false, "Enable Hide Resources Management - deprecated")
	flag.StringVar(&opts.PublishToQueue, "publishQueue", canopsis.ServiceQueueName, "Publish event to this queue")
	flag.StringVar(&opts.PostProcessorsDirectory, "postProcessorsDirectory", ".", "The path of the directory containing the post-processing plugins.")
	flag.BoolVar(&opts.IgnoreDefaultTomlConfig, "ignoreDefaultTomlConfig", false, "load toml file values into database. - deprecated")
	flag.DurationVar(&opts.PeriodicalWaitTime, "periodicalWaitTime", canopsis.PeriodicalWaitTime, "Duration to wait between two run of periodical process")
	flag.BoolVar(&opts.WithRemediation, "withRemediation", false, "Start remediation instructions")
	flag.BoolVar(&opts.Version, "version", false, "Show the version information")
	flag.Parse()

	return opts
}

func NewEngine(
	ctx context.Context,
	options Options,
	dbClient mongo.DbClient,
	cfg config.CanopsisConf,
	metricsSender metrics.Sender,
	logger zerolog.Logger,
) libengine.Engine {
	defer depmake.Catch(logger)

	m := DependencyMaker{}
	alarmConfigProvider := config.NewAlarmConfigProvider(cfg, logger)
	timezoneConfigProvider := config.NewTimezoneConfigProvider(cfg, logger)
	amqpConnection := m.DepAmqpConnection(logger, cfg)
	amqpChannel := m.DepAMQPChannelPub(amqpConnection)
	lockRedisClient := m.DepRedisSession(ctx, redis.EngineLockStorage, logger, cfg)
	pbhRedisClient := m.DepRedisSession(ctx, redis.PBehaviorLockStorage, logger, cfg)
	runInfoRedisClient := m.DepRedisSession(ctx, redis.EngineRunInfo, logger, cfg)

	serviceRpcClient := libengine.NewRPCClient(
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

	alarmStatusService := alarmstatus.NewService(flappingrule.NewAdapter(dbClient), alarmConfigProvider, logger)

	pbhRpcClient := libengine.NewRPCClient(
		canopsis.AxeRPCConsumerName,
		canopsis.PBehaviorRPCQueueServerName,
		canopsis.AxePbehaviorRPCClientQueueName,
		cfg.Global.PrefetchCount,
		cfg.Global.PrefetchSize,
		&rpcPBehaviorClientMessageProcessor{
			FeaturePrintEventOnError: options.FeaturePrintEventOnError,
			PublishCh:                amqpChannel,
			ServiceRpc:               serviceRpcClient,
			Executor:                 m.depOperationExecutor(dbClient, alarmConfigProvider, alarmStatusService, metricsSender),
			EntityAdapter:            entity.NewAdapter(dbClient),
			PbehaviorAdapter:         pbehavior.NewAdapter(dbClient),
			Decoder:                  json.NewDecoder(),
			Encoder:                  json.NewEncoder(),
			Logger:                   logger,
		},
		amqpChannel,
		logger,
	)
	pbhRpcClientForIdleRules := libengine.NewRPCClient(
		canopsis.AxeRPCConsumerName,
		canopsis.PBehaviorRPCQueueServerName,
		"",
		cfg.Global.PrefetchCount,
		cfg.Global.PrefetchSize,
		&rpcPBehaviorClientMessageProcessor{
			FeaturePrintEventOnError: options.FeaturePrintEventOnError,
			PublishCh:                amqpChannel,
			ServiceRpc:               serviceRpcClient,
			Executor:                 m.depOperationExecutor(dbClient, alarmConfigProvider, alarmStatusService, metricsSender),
			EntityAdapter:            entity.NewAdapter(dbClient),
			PbehaviorAdapter:         pbehavior.NewAdapter(dbClient),
			Decoder:                  json.NewDecoder(),
			Encoder:                  json.NewEncoder(),
			Logger:                   logger,
		},
		amqpChannel,
		logger,
	)

	rpcPublishQueues := []string{canopsis.PBehaviorRPCQueueServerName}
	var remediationRpcClient libengine.RPCClient
	if options.WithRemediation {
		remediationRpcClient = libengine.NewRPCClient(
			canopsis.AxeRPCConsumerName,
			canopsis.RemediationRPCQueueServerName,
			"",
			cfg.Global.PrefetchCount,
			cfg.Global.PrefetchSize,
			nil,
			amqpChannel,
			logger,
		)
		rpcPublishQueues = append(rpcPublishQueues, canopsis.RemediationRPCQueueServerName)
	}

	runInfoPeriodicalWorker := libengine.NewRunInfoPeriodicalWorker(
		options.PeriodicalWaitTime,
		libengine.NewRunInfoManager(runInfoRedisClient),
		libengine.NewInstanceRunInfo(canopsis.AxeEngineName, canopsis.AxeQueueName, options.PublishToQueue, nil, rpcPublishQueues),
		amqpChannel,
		logger,
	)

	metaAlarmEventProcessor := alarm.NewMetaAlarmEventProcessor(dbClient, alarm.NewAdapter(dbClient), correlation.NewRuleAdapter(dbClient),
		alarmStatusService, alarmConfigProvider, json.NewEncoder(), amqpChannel, canopsis.FIFOExchangeName, canopsis.FIFOQueueName,
		metricsSender, logger)

	engineAxe := libengine.New(
		func(ctx context.Context) error {
			runInfoPeriodicalWorker.Work(ctx)

			return alarmStatusService.Load(ctx)
		},
		func(ctx context.Context) {
			err := amqpConnection.Close()
			if err != nil {
				logger.Error().Err(err).Msg("failed to close amqp connection")
			}

			err = lockRedisClient.Close()
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

	techMetricsConfigProvider := config.NewTechMetricsConfigProvider(cfg, logger)
	techMetricsSender := techmetrics.NewSender(techMetricsConfigProvider, canopsis.TechMetricsFlushInterval,
		cfg.Global.ReconnectRetries, cfg.Global.GetReconnectTimeout(), logger)

	engineAxe.AddRoutine(func(ctx context.Context) error {
		techMetricsSender.Run(ctx)
		return nil
	})

	engineAxe.AddConsumer(libengine.NewDefaultConsumer(
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

			TechMetricsSender: techMetricsSender,
			EventProcessor: alarm.NewEventProcessor(
				dbClient,
				alarm.NewAdapter(dbClient),
				entity.NewAdapter(dbClient),
				correlation.NewRuleAdapter(dbClient),
				alarmConfigProvider,
				m.depOperationExecutor(dbClient, alarmConfigProvider, alarmStatusService, metricsSender),
				alarmStatusService,
				metricsSender,
				metaAlarmEventProcessor,
				statistics.NewEventStatisticsSender(dbClient, logger, timezoneConfigProvider),
				pbehavior.NewEntityTypeResolver(pbehavior.NewStore(pbhRedisClient, json.NewEncoder(), json.NewDecoder()), pbehavior.NewEntityMatcher(dbClient), logger),
				logger,
			),
			RemediationRpcClient:   remediationRpcClient,
			TimezoneConfigProvider: timezoneConfigProvider,
			Encoder:                json.NewEncoder(),
			Decoder:                json.NewDecoder(),
			Logger:                 logger,
			PbehaviorAdapter:       pbehavior.NewAdapter(dbClient),
		},
		logger,
	))
	engineAxe.AddConsumer(libengine.NewRPCServer(
		canopsis.AxeRPCConsumerName,
		canopsis.AxeRPCQueueServerName,
		cfg.Global.PrefetchCount,
		cfg.Global.PrefetchSize,
		amqpConnection,
		&rpcMessageProcessor{
			FeaturePrintEventOnError: options.FeaturePrintEventOnError,
			ServiceRpc:               serviceRpcClient,
			RMQChannel:               amqpChannel,
			PbhRpc:                   pbhRpcClient,
			RemediationRpc:           remediationRpcClient,
			MetaAlarmEventProcessor:  metaAlarmEventProcessor,
			Executor:                 m.depOperationExecutor(dbClient, alarmConfigProvider, alarmStatusService, metricsSender),
			Encoder:                  json.NewEncoder(),
			Decoder:                  json.NewDecoder(),
			Logger:                   logger,
		},
		logger,
	))
	engineAxe.AddConsumer(serviceRpcClient)
	engineAxe.AddConsumer(pbhRpcClient)
	engineAxe.AddPeriodicalWorker("run info", runInfoPeriodicalWorker)
	engineAxe.AddPeriodicalWorker("local cache", &reloadLocalCachePeriodicalWorker{
		PeriodicalInterval: options.PeriodicalWaitTime,
		AlarmStatusService: alarmStatusService,
		Logger:             logger,
	})
	engineAxe.AddPeriodicalWorker("alarms", libengine.NewLockedPeriodicalWorker(
		redis.NewLockClient(lockRedisClient),
		redis.AxePeriodicalLockKey,
		&periodicalWorker{
			TechMetricsSender:  techMetricsSender,
			PeriodicalInterval: options.PeriodicalWaitTime,
			ChannelPub:         amqpChannel,
			AlarmService:       alarm.NewService(alarm.NewAdapter(dbClient), resolverule.NewAdapter(dbClient), alarmStatusService, logger),
			AlarmAdapter:       alarm.NewAdapter(dbClient),
			Encoder:            json.NewEncoder(),
			IdleAlarmService: idlealarm.NewService(
				idlerule.NewRuleAdapter(dbClient),
				alarm.NewAdapter(dbClient),
				entity.NewAdapter(dbClient),
				pbhRpcClientForIdleRules,
				json.NewEncoder(),
				logger,
			),
			AlarmConfigProvider: alarmConfigProvider,
			Logger:              logger,
		},
		logger,
	))
	engineAxe.AddPeriodicalWorker("resolve archiver", libengine.NewLockedPeriodicalWorker(
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
	engineAxe.AddPeriodicalWorker("config", libengine.NewLoadConfigPeriodicalWorker(
		options.PeriodicalWaitTime,
		config.NewAdapter(dbClient),
		logger,
		alarmConfigProvider,
		timezoneConfigProvider,
		techMetricsConfigProvider,
	))

	return engineAxe
}

type DependencyMaker struct {
	depmake.DependencyMaker
}

func (m DependencyMaker) depOperationExecutor(
	dbClient mongo.DbClient,
	configProvider config.AlarmConfigProvider,
	alarmStatusService alarmstatus.Service,
	metricsSender metrics.Sender,
) operation.Executor {
	entityAdapter := entity.NewAdapter(dbClient)
	container := operation.NewExecutorContainer()
	container.Set(types.EventTypeAck, executor.NewAckExecutor(metricsSender, configProvider))
	container.Set(types.EventTypeAckremove, executor.NewAckRemoveExecutor(metricsSender, configProvider))
	container.Set(types.EventTypeActivate, executor.NewActivateExecutor())
	container.Set(types.EventTypeAssocTicket, executor.NewAssocTicketExecutor(metricsSender))
	container.Set(types.EventTypeCancel, executor.NewCancelExecutor(configProvider, alarmStatusService))
	container.Set(types.EventTypeChangestate, executor.NewChangeStateExecutor(configProvider, alarmStatusService, metricsSender))
	container.Set(types.EventTypeComment, executor.NewCommentExecutor(configProvider))
	container.Set(types.EventTypeDeclareTicket, executor.NewDeclareTicketExecutor())
	container.Set(types.EventTypeDeclareTicketWebhook, executor.NewDeclareTicketWebhookExecutor(configProvider, metricsSender))
	container.Set(types.EventTypeDone, executor.NewDoneExecutor(configProvider))
	container.Set(types.EventTypeKeepstate, executor.NewChangeStateExecutor(configProvider, alarmStatusService, metricsSender))
	container.Set(types.EventTypePbhEnter, executor.NewPbhEnterExecutor(configProvider, metricsSender))
	container.Set(types.EventTypePbhLeave, executor.NewPbhLeaveExecutor(configProvider, metricsSender))
	container.Set(types.EventTypePbhLeaveAndEnter, executor.NewPbhLeaveAndEnterExecutor(configProvider, metricsSender))
	container.Set(types.EventTypeResolveDone, executor.NewResolveStatExecutor(executor.NewResolveDoneExecutor(), entityAdapter, metricsSender))
	container.Set(types.EventTypeResolveCancel, executor.NewResolveStatExecutor(executor.NewResolveCancelExecutor(), entityAdapter, metricsSender))
	container.Set(types.EventTypeResolveClose, executor.NewResolveStatExecutor(executor.NewResolveCloseExecutor(), entityAdapter, metricsSender))
	container.Set(types.EventTypeEntityToggled, executor.NewResolveStatExecutor(executor.NewResolveDisabledExecutor(), entityAdapter, metricsSender))
	container.Set(types.EventTypeSnooze, executor.NewSnoozeExecutor(configProvider))
	container.Set(types.EventTypeUncancel, executor.NewUncancelExecutor(configProvider, alarmStatusService))
	container.Set(types.EventTypeUnsnooze, executor.NewUnsnoozeExecutor())
	container.Set(types.EventTypeUpdateStatus, executor.NewUpdateStatusExecutor(configProvider, alarmStatusService))
	container.Set(types.EventTypeInstructionStarted, executor.NewInstructionExecutor(metricsSender))
	container.Set(types.EventTypeInstructionPaused, executor.NewInstructionExecutor(metricsSender))
	container.Set(types.EventTypeInstructionResumed, executor.NewInstructionExecutor(metricsSender))
	container.Set(types.EventTypeInstructionCompleted, executor.NewInstructionExecutor(metricsSender))
	container.Set(types.EventTypeInstructionAborted, executor.NewInstructionExecutor(metricsSender))
	container.Set(types.EventTypeInstructionFailed, executor.NewInstructionExecutor(metricsSender))
	container.Set(types.EventTypeAutoInstructionStarted, executor.NewInstructionExecutor(metricsSender))
	container.Set(types.EventTypeAutoInstructionCompleted, executor.NewInstructionExecutor(metricsSender))
	container.Set(types.EventTypeAutoInstructionFailed, executor.NewInstructionExecutor(metricsSender))
	container.Set(types.EventTypeInstructionJobStarted, executor.NewInstructionExecutor(metricsSender))
	container.Set(types.EventTypeInstructionJobCompleted, executor.NewInstructionExecutor(metricsSender))
	container.Set(types.EventTypeInstructionJobAborted, executor.NewInstructionExecutor(metricsSender))
	container.Set(types.EventTypeInstructionJobFailed, executor.NewInstructionExecutor(metricsSender))
	container.Set(types.EventTypeJunitTestSuiteUpdated, executor.NewJunitExecutor())
	container.Set(types.EventTypeJunitTestCaseUpdated, executor.NewJunitExecutor())

	return executor.NewMongoUpdateExecutor(
		executor.NewCombinedExecutor(container),
		alarm.NewAdapter(dbClient),
		entity.NewAdapter(dbClient),
	)
}
