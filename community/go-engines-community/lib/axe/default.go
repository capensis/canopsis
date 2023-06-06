package axe

import (
	"context"
	"flag"
	"fmt"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmstatus"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmtag"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	libengine "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice/statecounters"
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
	"github.com/bsm/redislock"
	"github.com/rs/zerolog"
)

type Options struct {
	Version                  bool
	FeaturePrintEventOnError bool
	ModeDebug                bool
	PublishToQueue           string
	FifoAckExchange          string
	PeriodicalWaitTime       time.Duration
	TagsPeriodicalWaitTime   time.Duration
	WithRemediation          bool
	RecomputeAllOnInit       bool
}

func ParseOptions() Options {
	opts := Options{}

	flag.BoolVar(&opts.ModeDebug, "d", false, "debug")
	flag.BoolVar(&opts.FeaturePrintEventOnError, "printEventOnError", false, "Print event on processing error")
	flag.StringVar(&opts.PublishToQueue, "publishQueue", canopsis.ServiceQueueName, "Publish event to this queue")
	flag.DurationVar(&opts.PeriodicalWaitTime, "periodicalWaitTime", canopsis.PeriodicalWaitTime, "Duration to wait between two run of periodical process")
	flag.StringVar(&opts.FifoAckExchange, "fifoAckExchange", canopsis.FIFOAckExchangeName, "Publish FIFO Ack event to this exchange.")
	flag.DurationVar(&opts.TagsPeriodicalWaitTime, "tagsPeriodicalWaitTime", 5*time.Second, "Duration to wait between two run of periodical process to update alarm tags")
	flag.BoolVar(&opts.WithRemediation, "withRemediation", false, "Start remediation instructions")
	flag.BoolVar(&opts.RecomputeAllOnInit, "recomputeAllOnInit", false, "Recompute entity services on init.")
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
	autoInstructionMatcher AutoInstructionMatcher,
	logger zerolog.Logger,
) libengine.Engine {
	defer depmake.Catch(logger)

	m := DependencyMaker{}
	alarmConfigProvider := config.NewAlarmConfigProvider(cfg, logger)
	timezoneConfigProvider := config.NewTimezoneConfigProvider(cfg, logger)
	dataStorageConfigProvider := config.NewDataStorageConfigProvider(cfg, logger)
	userInterfaceAdapter := config.NewUserInterfaceAdapter(dbClient)
	userInterfaceConfig, err := userInterfaceAdapter.GetConfig(ctx)
	if err != nil {
		panic(fmt.Errorf("dependency error: %s: %v", "can't get user interface config", err))
	}
	userInterfaceConfigProvider := config.NewUserInterfaceConfigProvider(userInterfaceConfig, logger)
	amqpConnection := m.DepAmqpConnection(logger, cfg)
	amqpChannel := m.DepAMQPChannelPub(amqpConnection)
	lockRedisClient := m.DepRedisSession(ctx, redis.EngineLockStorage, logger, cfg)
	pbhRedisClient := m.DepRedisSession(ctx, redis.PBehaviorLockStorage, logger, cfg)
	runInfoRedisClient := m.DepRedisSession(ctx, redis.EngineRunInfo, logger, cfg)
	initRedisLock := redis.NewLockClient(lockRedisClient)

	alarmStatusService := alarmstatus.NewService(flappingrule.NewAdapter(dbClient), alarmConfigProvider, logger)

	stateCountersService := statecounters.NewStateCountersService(
		dbClient,
		amqpChannel,
		canopsis.FIFOExchangeName,
		canopsis.FIFOQueueName,
		json.NewEncoder(),
		logger,
	)
	actionRpcClient := libengine.NewRPCClient(
		canopsis.AxeRPCConsumerName,
		canopsis.ActionAxeRPCClientQueueName,
		"",
		cfg.Global.PrefetchCount,
		cfg.Global.PrefetchSize,
		nil,
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

	idleSinceService := entityservice.NewService(
		entityservice.NewAdapter(dbClient),
		entity.NewAdapter(dbClient),
		logger,
	)

	entityAdapter := entity.NewAdapter(dbClient)
	alarmAdapter := alarm.NewAdapter(dbClient)

	pbhRpcClient := libengine.NewRPCClient(
		canopsis.AxeRPCConsumerName,
		canopsis.PBehaviorRPCQueueServerName,
		canopsis.AxePbehaviorRPCClientQueueName,
		cfg.Global.PrefetchCount,
		cfg.Global.PrefetchSize,
		&rpcPBehaviorClientMessageProcessor{
			DbClient:       dbClient,
			MetricsSender:  metricsSender,
			PublishCh:      amqpChannel,
			RemediationRpc: remediationRpcClient,
			Executor: m.depOperationExecutor(dbClient, alarmConfigProvider, userInterfaceConfigProvider,
				alarmStatusService, metricsSender),
			EntityAdapter:            entityAdapter,
			AlarmAdapter:             alarmAdapter,
			PbehaviorAdapter:         pbehavior.NewAdapter(dbClient),
			StateCountersService:     stateCountersService,
			Decoder:                  json.NewDecoder(),
			Encoder:                  json.NewEncoder(),
			Logger:                   logger,
			FeaturePrintEventOnError: options.FeaturePrintEventOnError,
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
			RemediationRpc:           remediationRpcClient,
			Executor: m.depOperationExecutor(dbClient, alarmConfigProvider, userInterfaceConfigProvider,
				alarmStatusService, metricsSender),
			EntityAdapter:    entity.NewAdapter(dbClient),
			PbehaviorAdapter: pbehavior.NewAdapter(dbClient),
			Decoder:          json.NewDecoder(),
			Encoder:          json.NewEncoder(),
			Logger:           logger,
		},
		amqpChannel,
		logger,
	)

	runInfoPeriodicalWorker := libengine.NewRunInfoPeriodicalWorker(
		options.PeriodicalWaitTime,
		libengine.NewRunInfoManager(runInfoRedisClient),
		libengine.NewInstanceRunInfo(canopsis.AxeEngineName, canopsis.AxeQueueName, options.PublishToQueue, nil, rpcPublishQueues),
		amqpChannel,
		logger,
	)

	metaAlarmEventProcessor := NewMetaAlarmEventProcessor(dbClient, alarm.NewAdapter(dbClient), correlation.NewRuleAdapter(dbClient),
		alarmStatusService, alarmConfigProvider, json.NewEncoder(), amqpChannel, canopsis.FIFOExchangeName, canopsis.FIFOQueueName,
		metricsSender, logger)

	tagUpdater := alarmtag.NewUpdater(dbClient)

	engineAxe := libengine.New(
		func(ctx context.Context) error {
			if options.RecomputeAllOnInit {
				_, err := initRedisLock.Obtain(ctx, redis.AxeEntityServiceStateLockKey,
					options.PeriodicalWaitTime, &redislock.Options{
						RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(1*time.Second), 1),
					})
				if err != nil {
					// Lock is set for options.PeriodicalWaitTime TTL, other instances should skip actions below
					if err != redislock.ErrNotObtained {
						return fmt.Errorf("cannot obtain lock: %w", err)
					}
				} else {
					logger.Info().Msg("started to send recompute entity service events")

					err = stateCountersService.RecomputeAllServices(ctx)
					if err != nil {
						return fmt.Errorf("failed to send recompute entity service events: %w", err)
					}

					logger.Info().Msg("finished to send recompute entity service events")
				}
			}

			runInfoPeriodicalWorker.Work(ctx)

			err := alarmStatusService.Load(ctx)
			if err != nil {
				return err
			}

			return autoInstructionMatcher.Load(ctx)
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
		options.FifoAckExchange,
		canopsis.FIFOAckQueueName,
		amqpConnection,
		&messageProcessor{
			FeaturePrintEventOnError: options.FeaturePrintEventOnError,
			TechMetricsSender:        techMetricsSender,
			EventProcessor: NewEventProcessor(
				dbClient,
				alarm.NewAdapter(dbClient),
				entity.NewAdapter(dbClient),
				correlation.NewRuleAdapter(dbClient),
				alarmConfigProvider,
				m.depOperationExecutor(dbClient, alarmConfigProvider, userInterfaceConfigProvider, alarmStatusService, metricsSender),
				alarmStatusService,
				metricsSender,
				metaAlarmEventProcessor,
				statistics.NewEventStatisticsSender(dbClient, logger, timezoneConfigProvider),
				stateCountersService,
				pbehavior.NewEntityTypeResolver(pbehavior.NewStore(pbhRedisClient, json.NewEncoder(), json.NewDecoder()), pbehavior.NewEntityMatcher(dbClient), logger),
				autoInstructionMatcher,
				logger,
			),
			RemediationRpcClient:   remediationRpcClient,
			TimezoneConfigProvider: timezoneConfigProvider,
			Encoder:                json.NewEncoder(),
			Decoder:                json.NewDecoder(),
			Logger:                 logger,
			PbehaviorAdapter:       pbehavior.NewAdapter(dbClient),
			TagUpdater:             tagUpdater,
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
			DbClient:       dbClient,
			MetricsSender:  metricsSender,
			EntityAdapter:  entityAdapter,
			AlarmAdapter:   alarmAdapter,
			RMQChannel:     amqpChannel,
			PbhRpc:         pbhRpcClient,
			RemediationRpc: remediationRpcClient,
			Executor: m.depOperationExecutor(dbClient, alarmConfigProvider, userInterfaceConfigProvider,
				alarmStatusService, metricsSender),
			ActionRpc:                actionRpcClient,
			MetaAlarmEventProcessor:  metaAlarmEventProcessor,
			StateCountersService:     stateCountersService,
			Decoder:                  json.NewDecoder(),
			Encoder:                  json.NewEncoder(),
			Logger:                   logger,
			FeaturePrintEventOnError: options.FeaturePrintEventOnError,
		},
		logger,
	))
	engineAxe.AddConsumer(pbhRpcClient)
	engineAxe.AddPeriodicalWorker("run info", runInfoPeriodicalWorker)
	engineAxe.AddPeriodicalWorker("local cache", &reloadLocalCachePeriodicalWorker{
		PeriodicalInterval:     options.PeriodicalWaitTime,
		AlarmStatusService:     alarmStatusService,
		AutoInstructionMatcher: autoInstructionMatcher,
		Logger:                 logger,
	})
	engineAxe.AddPeriodicalWorker("tags", &tagPeriodicalWorker{
		PeriodicalInterval: options.TagsPeriodicalWaitTime,
		TagUpdater:         tagUpdater,
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
			DataStorageConfigProvider: dataStorageConfigProvider,
			LimitConfigAdapter:        datastorage.NewAdapter(dbClient),
			Logger:                    logger,
		},
		logger,
	))
	engineAxe.AddPeriodicalWorker("idle since", libengine.NewLockedPeriodicalWorker(
		redis.NewLockClient(lockRedisClient),
		redis.ServiceIdleSincePeriodicalLockKey,
		&idleSincePeriodicalWorker{
			IdleSinceService:   idleSinceService,
			PeriodicalInterval: options.PeriodicalWaitTime,
			Logger:             logger,
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
		dataStorageConfigProvider,
	))

	engineAxe.AddPeriodicalWorker("user_interface_config", libengine.NewLoadUserInterfaceConfigPeriodicalWorker(
		options.PeriodicalWaitTime,
		userInterfaceAdapter,
		logger,
		userInterfaceConfigProvider,
	))

	return engineAxe
}

type DependencyMaker struct {
	depmake.DependencyMaker
}

func (m DependencyMaker) depOperationExecutor(
	dbClient mongo.DbClient,
	configProvider config.AlarmConfigProvider,
	userInterfaceConfigProvider config.UserInterfaceConfigProvider,
	alarmStatusService alarmstatus.Service,
	metricsSender metrics.Sender,
) operation.Executor {
	entityAdapter := entity.NewAdapter(dbClient)
	container := operation.NewExecutorContainer()
	container.Set(types.EventTypeAck, executor.NewAckExecutor(configProvider))
	container.Set(types.EventTypeAckremove, executor.NewAckRemoveExecutor(configProvider))
	container.Set(types.EventTypeActivate, executor.NewActivateExecutor())
	container.Set(types.EventTypeAssocTicket, executor.NewAssocTicketExecutor(metricsSender))
	container.Set(types.EventTypeCancel, executor.NewCancelExecutor(configProvider, alarmStatusService))
	container.Set(types.EventTypeChangestate, executor.NewChangeStateExecutor(configProvider, userInterfaceConfigProvider,
		alarmStatusService))
	container.Set(types.EventTypeComment, executor.NewCommentExecutor(configProvider))
	container.Set(types.EventTypeDeclareTicketWebhook, executor.NewDeclareTicketWebhookExecutor(configProvider, metricsSender))
	container.Set(types.EventTypePbhEnter, executor.NewPbhEnterExecutor(configProvider, metricsSender))
	container.Set(types.EventTypePbhLeave, executor.NewPbhLeaveExecutor(configProvider, metricsSender))
	container.Set(types.EventTypePbhLeaveAndEnter, executor.NewPbhLeaveAndEnterExecutor(configProvider, metricsSender))
	container.Set(types.EventTypeResolveCancel, executor.NewResolveStatExecutor(executor.NewResolveCancelExecutor(), entityAdapter, metricsSender))
	container.Set(types.EventTypeResolveClose, executor.NewResolveStatExecutor(executor.NewResolveCloseExecutor(), entityAdapter, metricsSender))
	container.Set(types.EventTypeResolveDeleted, executor.NewResolveStatExecutor(executor.NewResolveDeletedExecutor(), entityAdapter, metricsSender))
	container.Set(types.EventTypeEntityToggled, executor.NewResolveStatExecutor(executor.NewResolveDisabledExecutor(), entityAdapter, metricsSender))
	container.Set(types.EventTypeSnooze, executor.NewSnoozeExecutor(configProvider))
	container.Set(types.EventTypeUncancel, executor.NewUncancelExecutor(configProvider, alarmStatusService))
	container.Set(types.EventTypeUnsnooze, executor.NewUnsnoozeExecutor())
	container.Set(types.EventTypeUpdateStatus, executor.NewUpdateStatusExecutor(configProvider, alarmStatusService))
	container.Set(types.EventTypeWebhookStarted, executor.NewWebhookStartExecutor())
	container.Set(types.EventTypeWebhookCompleted, executor.NewWebhookCompleteExecutor(metricsSender))
	container.Set(types.EventTypeWebhookFailed, executor.NewWebhookFailExecutor())
	container.Set(types.EventTypeAutoWebhookStarted, executor.NewAutoWebhookStartExecutor())
	container.Set(types.EventTypeAutoWebhookCompleted, executor.NewAutoWebhookCompleteExecutor(metricsSender))
	container.Set(types.EventTypeAutoWebhookFailed, executor.NewAutoWebhookFailExecutor())
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
	container.Set(types.EventTypeInstructionJobFailed, executor.NewInstructionExecutor(metricsSender))
	container.Set(types.EventTypeAutoInstructionActivate, executor.NewAutoInstructionActivateExecutor())
	container.Set(types.EventTypeJunitTestSuiteUpdated, executor.NewJunitExecutor())
	container.Set(types.EventTypeJunitTestCaseUpdated, executor.NewJunitExecutor())

	return executor.NewMongoUpdateExecutor(
		executor.NewCombinedExecutor(container),
		alarm.NewAdapter(dbClient),
		entity.NewAdapter(dbClient),
	)
}
