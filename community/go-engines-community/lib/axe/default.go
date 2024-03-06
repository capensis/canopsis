package axe

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/axe/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	libalarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmstatus"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmtag"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	libengine "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters/calculator"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/flappingrule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/healthcheck"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/idlealarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/idlerule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/resolverule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/statistics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/techmetrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
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
	SliPeriodicalWaitTime    time.Duration
	WithRemediation          bool
	WithDynamicInfos         bool
	RecomputeAllOnInit       bool
	Workers                  int
}

func ParseOptions() Options {
	opts := Options{}

	flag.BoolVar(&opts.ModeDebug, "d", false, "debug")
	flag.BoolVar(&opts.FeaturePrintEventOnError, "printEventOnError", false, "Print event on processing error")
	flag.StringVar(&opts.PublishToQueue, "publishQueue", canopsis.ActionQueueName, "Publish event to this queue")
	flag.DurationVar(&opts.PeriodicalWaitTime, "periodicalWaitTime", canopsis.PeriodicalWaitTime, "Duration to wait between two run of periodical process")
	flag.StringVar(&opts.FifoAckExchange, "fifoAckExchange", canopsis.FIFOAckExchangeName, "Publish FIFO Ack event to this exchange.")
	flag.DurationVar(&opts.TagsPeriodicalWaitTime, "tagsPeriodicalWaitTime", 5*time.Second, "Duration to wait between two run of periodical process to update alarm tags")
	flag.DurationVar(&opts.SliPeriodicalWaitTime, "sliPeriodicalWaitTime", 5*time.Minute, "Duration to wait between two run of periodical process to update SLI metrics")
	flag.BoolVar(&opts.WithRemediation, "withRemediation", false, "Start remediation instructions")
	flag.BoolVar(&opts.WithDynamicInfos, "withDynamicInfos", false, "Update dynamic infos on RPC changes")
	flag.BoolVar(&opts.RecomputeAllOnInit, "recomputeAllOnInit", false, "Recompute entity services on init.")
	flag.BoolVar(&opts.Version, "version", false, "Show the version information")
	flag.IntVar(&opts.Workers, "workers", 2, "Amount of workers to process main event flow")
	flag.Parse()

	return opts
}

func NewEngine(
	ctx context.Context,
	options Options,
	dbClient mongo.DbClient,
	cfg config.CanopsisConf,
	metricsSender metrics.Sender,
	autoInstructionMatcher event.AutoInstructionMatcher,
	logger zerolog.Logger,
) libengine.Engine {
	defer depmake.Catch(logger)

	m := DependencyMaker{}
	alarmConfigProvider := config.NewAlarmConfigProvider(cfg, logger)
	timezoneConfigProvider := config.NewTimezoneConfigProvider(cfg, logger)
	dataStorageConfigProvider := config.NewDataStorageConfigProvider(cfg, logger)
	templateConfigProvider := config.NewTemplateConfigProvider(cfg, logger)
	userInterfaceAdapter := config.NewUserInterfaceAdapter(dbClient)
	userInterfaceConfig, err := userInterfaceAdapter.GetConfig(ctx)
	if err != nil {
		panic(fmt.Errorf("dependency error: %s: %w", "can't get user interface config", err))
	}
	userInterfaceConfigProvider := config.NewUserInterfaceConfigProvider(userInterfaceConfig, logger)
	amqpConnection := m.DepAmqpConnection(logger, cfg)
	amqpChannel := m.DepAMQPChannelPub(amqpConnection)
	lockRedisClient := m.DepRedisSession(ctx, redis.EngineLockStorage, logger, cfg)
	pbhRedisClient := m.DepRedisSession(ctx, redis.PBehaviorLockStorage, logger, cfg)
	runInfoRedisClient := m.DepRedisSession(ctx, redis.EngineRunInfo, logger, cfg)
	initRedisLock := redis.NewLockClient(lockRedisClient)

	alarmStatusService := alarmstatus.NewService(flappingrule.NewAdapter(dbClient), alarmConfigProvider, logger)

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
	var dynamicInfosRpcClient libengine.RPCClient
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

	if options.WithDynamicInfos {
		dynamicInfosRpcClient = libengine.NewRPCClient(
			canopsis.AxeRPCConsumerName,
			canopsis.DynamicInfosRPCQueueServerName,
			canopsis.AxeDynamicInfosRPCClientQueueName,
			cfg.Global.PrefetchCount,
			cfg.Global.PrefetchSize,
			&rpcDynamicInfosClientMessageProcessor{
				FeaturePrintEventOnError: options.FeaturePrintEventOnError,
				PublishCh:                amqpChannel,
				Decoder:                  json.NewDecoder(),
				Encoder:                  json.NewEncoder(),
				Logger:                   logger,
			},
			amqpChannel,
			logger,
		)
		rpcPublishQueues = append(rpcPublishQueues, canopsis.DynamicInfosRPCQueueServerName)
	}

	idleSinceService := entityservice.NewService(
		entityservice.NewAdapter(dbClient),
		entity.NewAdapter(dbClient),
		logger,
	)

	entityAdapter := entity.NewAdapter(dbClient)
	alarmAdapter := libalarm.NewAdapter(dbClient)

	metaAlarmEventProcessor := NewMetaAlarmEventProcessor(dbClient, libalarm.NewAdapter(dbClient), correlation.NewRuleAdapter(dbClient),
		alarmStatusService, alarmConfigProvider, json.NewEncoder(), amqpChannel, metricsSender, correlation.NewMetaAlarmStateService(dbClient),
		template.NewExecutor(templateConfigProvider, timezoneConfigProvider), canopsis.AxeConnector, logger)

	externalTagUpdater := alarmtag.NewExternalUpdater(dbClient)
	internalTagAlarmMatcher := alarmtag.NewInternalTagAlarmMatcher(dbClient)

	eventsSender := entitycounters.NewEventSender(json.NewEncoder(), amqpChannel, canopsis.FIFOExchangeName, canopsis.FIFOQueueName, canopsis.AxeConnector)
	entityServiceCountersCalculator := calculator.NewEntityServiceCountersCalculator(dbClient, template.NewExecutor(templateConfigProvider, timezoneConfigProvider), eventsSender)
	componentCountersCalculator := calculator.NewComponentCountersCalculator(dbClient, eventsSender)

	eventProcessor := m.EventProcessor(
		dbClient,
		alarmConfigProvider,
		userInterfaceConfigProvider,
		alarmStatusService,
		pbehavior.NewEntityTypeResolver(pbehavior.NewStore(pbhRedisClient, json.NewEncoder(), json.NewDecoder()), logger),
		autoInstructionMatcher,
		entityServiceCountersCalculator,
		componentCountersCalculator,
		eventsSender,
		metaAlarmEventProcessor,
		metricsSender,
		statistics.NewEventStatisticsSender(dbClient, logger, timezoneConfigProvider),
		remediationRpcClient,
		externalTagUpdater,
		internalTagAlarmMatcher,
		amqpChannel,
		logger,
	)

	pbhRpcClient := libengine.NewRPCClient(
		canopsis.AxeRPCConsumerName,
		canopsis.PBehaviorRPCQueueServerName,
		canopsis.AxePbehaviorRPCClientQueueName,
		cfg.Global.PrefetchCount,
		cfg.Global.PrefetchSize,
		&rpcPBehaviorClientMessageProcessor{
			DbClient:                 dbClient,
			MetricsSender:            metricsSender,
			PublishCh:                amqpChannel,
			RemediationRpc:           remediationRpcClient,
			EventProcessor:           eventProcessor,
			EntityAdapter:            entityAdapter,
			AlarmAdapter:             alarmAdapter,
			PbehaviorAdapter:         pbehavior.NewAdapter(dbClient),
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
			EventProcessor:           eventProcessor,
			EntityAdapter:            entity.NewAdapter(dbClient),
			PbehaviorAdapter:         pbehavior.NewAdapter(dbClient),
			Decoder:                  json.NewDecoder(),
			Encoder:                  json.NewEncoder(),
			Logger:                   logger,
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

	engineAxe := libengine.New(
		func(ctx context.Context) error {
			if options.RecomputeAllOnInit {
				_, err := initRedisLock.Obtain(ctx, redis.AxeEntityServiceStateLockKey,
					options.PeriodicalWaitTime, &redislock.Options{
						RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(1*time.Second), 1),
					})
				if err != nil {
					// Lock is set for options.PeriodicalWaitTime TTL, other instances should skip actions below
					if !errors.Is(err, redislock.ErrNotObtained) {
						return fmt.Errorf("cannot obtain lock: %w", err)
					}
				} else {
					logger.Info().Msg("started to send recompute entity service events")

					err = entityServiceCountersCalculator.RecomputeAll(ctx)
					if err != nil {
						return fmt.Errorf("failed to send recompute entity service events: %w", err)
					}

					logger.Info().Msg("finished to send recompute entity service events")

					logger.Info().Msg("started to send recompute components events")

					err = componentCountersCalculator.RecomputeAll(ctx)
					if err != nil {
						return fmt.Errorf("failed to send recompute components events: %w", err)
					}

					logger.Info().Msg("finished to send recompute components events")
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

	mainMessageProcessor := &MessageProcessor{
		FeaturePrintEventOnError: options.FeaturePrintEventOnError,
		EventProcessor:           eventProcessor,
		Encoder:                  json.NewEncoder(),
		Decoder:                  json.NewDecoder(),
		TechMetricsSender:        techMetricsSender,
		AlarmCollection:          dbClient.Collection(mongo.AlarmMongoCollection),
		Logger:                   logger,
	}
	engineAxe.AddConsumer(libengine.NewConcurrentConsumer(
		canopsis.AxeConsumerName,
		canopsis.AxeQueueName,
		cfg.Global.PrefetchCount,
		cfg.Global.PrefetchSize,
		false,
		"",
		options.PublishToQueue,
		options.FifoAckExchange,
		canopsis.FIFOAckQueueName,
		options.Workers,
		amqpConnection,
		mainMessageProcessor,
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
			EventProcessor:           eventProcessor,
			ActionRpc:                actionRpcClient,
			PbhRpc:                   pbhRpcClient,
			DynamicInfosRpc:          dynamicInfosRpcClient,
			Encoder:                  json.NewEncoder(),
			Decoder:                  json.NewDecoder(),
			Logger:                   logger,
		},
		logger,
	))
	engineAxe.AddConsumer(pbhRpcClient)
	if dynamicInfosRpcClient != nil {
		engineAxe.AddConsumer(dynamicInfosRpcClient)
	}

	engineAxe.AddPeriodicalWorker("run info", runInfoPeriodicalWorker)
	engineAxe.AddPeriodicalWorker("local cache", &reloadLocalCachePeriodicalWorker{
		PeriodicalInterval:      options.PeriodicalWaitTime,
		AlarmStatusService:      alarmStatusService,
		AutoInstructionMatcher:  autoInstructionMatcher,
		InternalTagAlarmMatcher: internalTagAlarmMatcher,
		Logger:                  logger,
	})
	engineAxe.AddPeriodicalWorker("external_tags", &externalTagPeriodicalWorker{
		PeriodicalInterval: options.TagsPeriodicalWaitTime,
		ExternalTagUpdater: externalTagUpdater,
		Logger:             logger,
	})
	engineAxe.AddPeriodicalWorker("internal_tags", libengine.NewLockedPeriodicalWorker(
		redis.NewLockClient(lockRedisClient),
		redis.AxeInternalTagsPeriodicalLockKey,
		&internalTagPeriodicalWorker{
			PeriodicalInterval: options.PeriodicalWaitTime,
			Logger:             logger,
			TagCollection:      dbClient.Collection(mongo.AlarmTagCollection),
			AlarmCollection:    dbClient.Collection(mongo.AlarmMongoCollection),
		},
		logger,
	))
	engineAxe.AddPeriodicalWorker("alarms", libengine.NewLockedPeriodicalWorker(
		redis.NewLockClient(lockRedisClient),
		redis.AxePeriodicalLockKey,
		&periodicalWorker{
			TechMetricsSender:  techMetricsSender,
			PeriodicalInterval: options.PeriodicalWaitTime,
			ChannelPub:         amqpChannel,
			AlarmService:       libalarm.NewService(libalarm.NewAdapter(dbClient), resolverule.NewAdapter(dbClient), alarmStatusService, logger),
			AlarmAdapter:       libalarm.NewAdapter(dbClient),
			Encoder:            json.NewEncoder(),
			IdleAlarmService: idlealarm.NewService(
				idlerule.NewRuleAdapter(dbClient),
				libalarm.NewAdapter(dbClient),
				entity.NewAdapter(dbClient),
				pbhRpcClientForIdleRules,
				canopsis.AxeConnector,
				json.NewEncoder(),
				logger,
			),
			AlarmConfigProvider: alarmConfigProvider,
			Logger:              logger,
		},
		logger,
	))
	engineAxe.AddPeriodicalWorker("resolve_archiver", libengine.NewLockedPeriodicalWorker(
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
		redis.AxeIdleSincePeriodicalLockKey,
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

	healthcheck.Start(ctx, healthcheck.NewChecker(
		"axe",
		mainMessageProcessor,
		json.NewEncoder(),
		true,
		false,
	), logger)

	return engineAxe
}

type DependencyMaker struct {
	depmake.DependencyMaker
}

func (m DependencyMaker) EventProcessor(
	dbClient mongo.DbClient,
	alarmConfigProvider config.AlarmConfigProvider,
	userInterfaceConfigProvider config.UserInterfaceConfigProvider,
	alarmStatusService alarmstatus.Service,
	pbhTypeResolver pbehavior.EntityTypeResolver,
	autoInstructionMatcher event.AutoInstructionMatcher,
	entityServiceCountersCalculator calculator.EntityServiceCountersCalculator,
	componentCountersCalculator calculator.ComponentCountersCalculator,
	eventsSender entitycounters.EventsSender,
	metaAlarmEventProcessor libalarm.MetaAlarmEventProcessor,
	metricsSender metrics.Sender,
	eventStatisticsSender statistics.EventStatisticsSender,
	remediationRpcClient libengine.RPCClient,
	externalTagUpdater alarmtag.ExternalUpdater,
	internalTagAlarmMatcher alarmtag.InternalTagAlarmMatcher,
	amqpPublisher amqp.Publisher,
	logger zerolog.Logger,
) event.Processor {
	container := event.NewProcessorContainer()

	container.Set(types.EventTypeCheck, event.NewCheckProcessor(dbClient, alarmConfigProvider, alarmStatusService,
		pbhTypeResolver, autoInstructionMatcher, metaAlarmEventProcessor, metricsSender,
		eventStatisticsSender, remediationRpcClient, externalTagUpdater, internalTagAlarmMatcher, entityServiceCountersCalculator, componentCountersCalculator, eventsSender, json.NewEncoder(), logger))
	container.Set(types.EventTypeNoEvents, event.NewNoEventsProcessor(dbClient, alarmConfigProvider, alarmStatusService,
		pbhTypeResolver, autoInstructionMatcher, entityServiceCountersCalculator, componentCountersCalculator, eventsSender, metaAlarmEventProcessor, metricsSender,
		remediationRpcClient, json.NewEncoder(), logger))
	container.Set(types.EventTypeAck, event.NewAckProcessor(dbClient, alarmConfigProvider, entityServiceCountersCalculator, eventsSender, metaAlarmEventProcessor, metricsSender, logger))
	container.Set(types.EventTypeAckremove, event.NewAckRemoveProcessor(dbClient, alarmConfigProvider, entityServiceCountersCalculator, eventsSender, metaAlarmEventProcessor, metricsSender, logger))
	container.Set(types.EventTypeActivate, event.NewActivateProcessor(dbClient, autoInstructionMatcher, remediationRpcClient, json.NewEncoder(), logger))
	container.Set(types.EventTypeAssocTicket, event.NewAssocTicketProcessor(dbClient, metaAlarmEventProcessor, metricsSender, logger))
	container.Set(types.EventTypeCancel, event.NewCancelProcessor(dbClient, metaAlarmEventProcessor, logger))
	container.Set(types.EventTypeChangestate, event.NewChangeStateProcessor(dbClient, alarmConfigProvider, userInterfaceConfigProvider,
		alarmStatusService, autoInstructionMatcher, entityServiceCountersCalculator, componentCountersCalculator, eventsSender, metaAlarmEventProcessor, metricsSender, remediationRpcClient, json.NewEncoder(), logger))
	container.Set(types.EventTypeComment, event.NewCommentProcessor(dbClient, alarmConfigProvider, metaAlarmEventProcessor, logger))
	container.Set(types.EventTypePbhEnter, event.NewPbhEnterProcessor(dbClient, autoInstructionMatcher, entityServiceCountersCalculator, componentCountersCalculator, eventsSender, metricsSender, remediationRpcClient, json.NewEncoder(), logger))
	container.Set(types.EventTypePbhLeave, event.NewPbhLeaveProcessor(dbClient, autoInstructionMatcher, entityServiceCountersCalculator, componentCountersCalculator, eventsSender, metricsSender, remediationRpcClient, json.NewEncoder(), logger))
	container.Set(types.EventTypePbhLeaveAndEnter, event.NewPbhLeaveAndEnterProcessor(dbClient, autoInstructionMatcher, entityServiceCountersCalculator, componentCountersCalculator, eventsSender, metricsSender, remediationRpcClient, json.NewEncoder(), logger))
	container.Set(types.EventTypeDeclareTicketWebhook, event.NewDeclareTicketWebhookProcessor(dbClient, metricsSender, amqpPublisher, json.NewEncoder(), logger))
	container.Set(types.EventTypeResolveCancel, event.NewResolveCancelProcessor(dbClient, entityServiceCountersCalculator, componentCountersCalculator, eventsSender, metaAlarmEventProcessor, metricsSender, remediationRpcClient, json.NewEncoder(), logger))
	container.Set(types.EventTypeResolveClose, event.NewResolveCloseProcessor(dbClient, entityServiceCountersCalculator, componentCountersCalculator, eventsSender, metaAlarmEventProcessor, metricsSender, remediationRpcClient, json.NewEncoder(), logger))
	container.Set(types.EventTypeResolveDeleted, event.NewResolveDeletedProcessor(dbClient, entityServiceCountersCalculator, componentCountersCalculator, eventsSender, metaAlarmEventProcessor, metricsSender, remediationRpcClient, json.NewEncoder(), logger))
	container.Set(types.EventTypeEntityToggled, event.NewEntityToggledProcessor(dbClient, entityServiceCountersCalculator, componentCountersCalculator, eventsSender, metaAlarmEventProcessor, metricsSender, remediationRpcClient, json.NewEncoder(), logger))
	container.Set(types.EventTypeRecomputeEntityService, event.NewRecomputeEntityServiceProcessor(dbClient, entityServiceCountersCalculator, componentCountersCalculator, eventsSender, metaAlarmEventProcessor, metricsSender, remediationRpcClient, json.NewEncoder(), logger))
	container.Set(types.EventTypeEntityUpdated, event.NewEntityUpdatedProcessor(dbClient, entityServiceCountersCalculator, componentCountersCalculator, eventsSender))
	container.Set(types.EventTypeUpdateCounters, event.NewUpdateCountersProcessor(dbClient, entityServiceCountersCalculator, eventsSender))
	container.Set(types.EventTypeSnooze, event.NewSnoozeProcessor(dbClient, metaAlarmEventProcessor, logger))
	container.Set(types.EventTypeUncancel, event.NewUncancelProcessor(dbClient, alarmStatusService, metaAlarmEventProcessor, logger))
	container.Set(types.EventTypeUnsnooze, event.NewUnsnoozeProcessor(dbClient, autoInstructionMatcher, remediationRpcClient, json.NewEncoder(), logger))
	container.Set(types.EventTypeUpdateStatus, event.NewUpdateStatusProcessor(dbClient, alarmStatusService, alarmConfigProvider, metaAlarmEventProcessor, logger))
	container.Set(types.EventTypeWebhookStarted, event.NewWebhookStartProcessor(dbClient))
	container.Set(types.EventTypeWebhookCompleted, event.NewWebhookCompleteProcessor(dbClient, metaAlarmEventProcessor, metricsSender, amqpPublisher, json.NewEncoder(), logger))
	container.Set(types.EventTypeWebhookFailed, event.NewWebhookFailProcessor(dbClient))
	container.Set(types.EventTypeAutoWebhookStarted, event.NewAutoWebhookStartProcessor(dbClient))
	container.Set(types.EventTypeAutoWebhookCompleted, event.NewAutoWebhookCompleteProcessor(dbClient, metaAlarmEventProcessor, metricsSender, logger))
	container.Set(types.EventTypeAutoWebhookFailed, event.NewAutoWebhookFailProcessor(dbClient))
	instructionProcessor := event.NewInstructionProcessor(dbClient, metricsSender, amqpPublisher, json.NewEncoder(), logger)
	container.Set(types.EventTypeInstructionStarted, instructionProcessor)
	container.Set(types.EventTypeInstructionPaused, instructionProcessor)
	container.Set(types.EventTypeInstructionResumed, instructionProcessor)
	container.Set(types.EventTypeInstructionCompleted, instructionProcessor)
	container.Set(types.EventTypeInstructionAborted, instructionProcessor)
	container.Set(types.EventTypeInstructionFailed, instructionProcessor)
	container.Set(types.EventTypeAutoInstructionStarted, instructionProcessor)
	container.Set(types.EventTypeAutoInstructionCompleted, instructionProcessor)
	container.Set(types.EventTypeAutoInstructionFailed, instructionProcessor)
	container.Set(types.EventTypeInstructionJobStarted, instructionProcessor)
	container.Set(types.EventTypeInstructionJobCompleted, instructionProcessor)
	container.Set(types.EventTypeInstructionJobFailed, instructionProcessor)
	container.Set(types.EventTypeAutoInstructionActivate, event.NewAutoInstructionActivateProcessor(dbClient))
	container.Set(types.EventTypeMetaAlarmChildActivate, event.NewMetaAlarmChildActivateProcessor(dbClient))
	container.Set(types.EventTypeMetaAlarmChildDeactivate, event.NewMetaAlarmChildDeactivateProcessor(dbClient))
	container.Set(types.EventTypeJunitTestSuiteUpdated, event.NewJunitProcessor(dbClient))
	container.Set(types.EventTypeJunitTestCaseUpdated, event.NewJunitProcessor(dbClient))
	container.Set(types.EventTypeRunDelayedScenario, event.NewForwardWithAlarmProcessor(dbClient))
	container.Set(types.EventTypeMetaAlarm, event.NewMetaAlarmProcessor(metaAlarmEventProcessor, autoInstructionMatcher, metricsSender, remediationRpcClient, json.NewEncoder(), logger))
	container.Set(types.EventTypeMetaAlarmAttachChildren, event.NewMetaAlarmAttachProcessor(metaAlarmEventProcessor, metricsSender, json.NewEncoder(), amqpPublisher, logger))
	container.Set(types.EventTypeMetaAlarmDetachChildren, event.NewMetaAlarmDetachProcessor(metaAlarmEventProcessor))
	container.Set(types.EventTypeMetaAlarmUngroup, event.NewForwardProcessor())
	container.Set(types.EventTypeManualMetaAlarmGroup, event.NewForwardProcessor())
	container.Set(types.EventTypeManualMetaAlarmUngroup, event.NewForwardProcessor())
	container.Set(types.EventTypeManualMetaAlarmUpdate, event.NewForwardProcessor())
	container.Set(types.EventTypeTrigger, event.NewTriggerProcessor(dbClient))

	return event.NewCombinedProcessor(container)
}
