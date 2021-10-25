package axe

import (
	"context"
	"flag"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/serviceweather"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	libengine "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/idlealarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/idlerule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation/executor"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/statsng"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/depmake"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
	"time"
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

func ParseOptions() Options {
	opts := Options{}

	flag.BoolVar(&opts.ModeDebug, "d", false, "debug")
	flag.BoolVar(&opts.FeaturePrintEventOnError, "printEventOnError", false, "Print event on processing error")
	flag.BoolVar(&opts.FeatureHideResources, "featureHideResources", false, "Enable Hide Resources Management - deprecated")
	flag.BoolVar(&opts.FeatureStatEvents, "featureStatEvents", false, "Send statistic events")
	flag.StringVar(&opts.PublishToQueue, "publishQueue", canopsis.ServiceQueueName, "Publish event to this queue")
	flag.StringVar(&opts.PostProcessorsDirectory, "postProcessorsDirectory", ".", "The path of the directory containing the post-processing plugins.")
	flag.BoolVar(&opts.IgnoreDefaultTomlConfig, "ignoreDefaultTomlConfig", false, "load toml file values into database. - deprecated")
	flag.DurationVar(&opts.PeriodicalWaitTime, "periodicalWaitTime", canopsis.PeriodicalWaitTime, "Duration to wait between two run of periodical process")
	flag.BoolVar(&opts.WithRemediation, "withRemediation", false, "Start remediation instructions")
	flag.Parse()

	flagVersion := flag.Bool("version", false, "version infos")
	if *flagVersion {
		canopsis.PrintVersionExit()
	}

	return opts
}

func Default(ctx context.Context, options Options, metricsSender metrics.Sender, postgresPool *pgxpool.Pool, logger zerolog.Logger) libengine.Engine {
	defer depmake.Catch(logger)

	m := DependencyMaker{}
	dbClient := m.DepMongoClient(ctx)
	cfg := m.DepConfig(ctx, dbClient)
	config.SetDbClientRetry(dbClient, cfg)
	alarmConfigProvider := config.NewAlarmConfigProvider(cfg, logger)
	timezoneConfigProvider := config.NewTimezoneConfigProvider(cfg, logger)
	amqpConnection := m.DepAmqpConnection(logger, cfg)
	amqpChannel := m.DepAMQPChannelPub(amqpConnection)
	lockRedisClient := m.DepRedisSession(ctx, redis.EngineLockStorage, logger, cfg)
	corrRedisClient := m.DepRedisSession(ctx, redis.CorrelationLockStorage, logger, cfg)
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

	statsService := statsng.NewService(
		m.DepAMQPChannelPub(amqpConnection),
		canopsis.StatsngExchangeName,
		canopsis.StatsngQueueName,
		json.NewEncoder(),
		serviceweather.NewStatsStore(dbClient),
		logger,
	)

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
			Executor:                 m.depOperationExecutor(dbClient, alarmConfigProvider, statsService, metricsSender),
			Decoder:                  json.NewDecoder(),
			Encoder:                  json.NewEncoder(),
			Logger:                   logger,
		},
		amqpChannel,
		logger,
	)

	rpcPublishQueues := make([]string, 0)
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

	engineAxe := libengine.New(
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

			if postgresPool != nil {
				postgresPool.Close()
			}
		},
		logger,
	)
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
			FeatureStatEvents:        options.FeatureStatEvents,
			EventProcessor: alarm.NewEventProcessor(
				alarm.NewAdapter(dbClient),
				entity.NewAdapter(dbClient),
				correlation.NewRuleAdapter(dbClient),
				alarmConfigProvider,
				m.depOperationExecutor(dbClient, alarmConfigProvider, statsService, metricsSender),
				redis.NewLockClient(corrRedisClient),
				metricsSender,
				logger,
			),
			StatsService:           statsService,
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
			PbhRpc:                   pbhRpcClient,
			RemediationRpc:           remediationRpcClient,
			AlarmAdapter:             alarm.NewAdapter(dbClient),
			Executor:                 m.depOperationExecutor(dbClient, alarmConfigProvider, statsService, metricsSender),
			Encoder:                  json.NewEncoder(),
			Decoder:                  json.NewDecoder(),
			Logger:                   logger,
		},
		logger,
	))
	engineAxe.AddConsumer(serviceRpcClient)
	engineAxe.AddConsumer(pbhRpcClient)
	engineAxe.AddPeriodicalWorker(libengine.NewRunInfoPeriodicalWorker(
		options.PeriodicalWaitTime,
		libengine.NewRunInfoManager(runInfoRedisClient),
		libengine.NewInstanceRunInfo(canopsis.AxeEngineName, canopsis.AxeQueueName, options.PublishToQueue, nil, rpcPublishQueues),
		amqpChannel,
		logger,
	))
	engineAxe.AddPeriodicalWorker(libengine.NewLockedPeriodicalWorker(
		redis.NewLockClient(lockRedisClient),
		redis.AxePeriodicalLockKey,
		&periodicalWorker{
			PeriodicalInterval: options.PeriodicalWaitTime,
			ChannelPub:         amqpChannel,
			AlarmService:       alarm.NewService(alarm.NewAdapter(dbClient), logger),
			AlarmAdapter:       alarm.NewAdapter(dbClient),
			Encoder:            json.NewEncoder(),
			IdleAlarmService: idlealarm.NewService(
				idlerule.NewRuleAdapter(dbClient),
				alarm.NewAdapter(dbClient),
				entity.NewAdapter(dbClient),
				pbehavior.NewEntityTypeResolver(
					pbehavior.NewStore(pbhRedisClient, json.NewEncoder(), json.NewDecoder()),
					pbehavior.NewComputedEntityMatcher(dbClient, pbhRedisClient, json.NewEncoder(), json.NewDecoder()),
				),
				json.NewEncoder(),
				logger,
			),
			AlarmConfigProvider: alarmConfigProvider,
			Logger:              logger,
		},
		logger,
	))
	engineAxe.AddPeriodicalWorker(libengine.NewLockedPeriodicalWorker(
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
	engineAxe.AddPeriodicalWorker(libengine.NewLoadConfigPeriodicalWorker(
		options.PeriodicalWaitTime,
		config.NewAdapter(dbClient),
		alarmConfigProvider,
		logger,
	))
	engineAxe.AddPeriodicalWorker(libengine.NewLoadConfigPeriodicalWorker(
		options.PeriodicalWaitTime,
		config.NewAdapter(dbClient),
		timezoneConfigProvider,
		logger,
	))

	return engineAxe
}

type DependencyMaker struct {
	depmake.DependencyMaker
}

func (m DependencyMaker) depOperationExecutor(
	dbClient mongo.DbClient,
	configProvider config.AlarmConfigProvider,
	statsService statsng.Service,
	metricsSender metrics.Sender,
) operation.Executor {
	entityAdapter := entity.NewAdapter(dbClient)
	container := operation.NewExecutorContainer()
	container.Set(types.EventTypeAck, executor.NewAckExecutor(metricsSender, configProvider))
	container.Set(types.EventTypeAckremove, executor.NewAckRemoveExecutor(metricsSender, configProvider))
	container.Set(types.EventTypeActivate, executor.NewActivateExecutor())
	container.Set(types.EventTypeAssocTicket, executor.NewAssocTicketExecutor(metricsSender))
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
	container.Set(types.EventTypeResolveDone, executor.NewResolveStatExecutor(executor.NewResolveDoneExecutor(), entityAdapter, statsService, metricsSender))
	container.Set(types.EventTypeResolveCancel, executor.NewResolveStatExecutor(executor.NewResolveCancelExecutor(), entityAdapter, statsService, metricsSender))
	container.Set(types.EventTypeResolveClose, executor.NewResolveStatExecutor(executor.NewResolveCloseExecutor(), entityAdapter, statsService, metricsSender))
	container.Set(types.EventTypeEntityToggled, executor.NewResolveStatExecutor(executor.NewResolveDisabledExecutor(), entityAdapter, statsService, metricsSender))
	container.Set(types.EventTypeSnooze, executor.NewSnoozeExecutor(configProvider))
	container.Set(types.EventTypeUncancel, executor.NewUncancelExecutor(configProvider))
	container.Set(types.EventTypeUnsnooze, executor.NewUnsnoozeExecutor())
	container.Set(types.EventTypeUpdateStatus, executor.NewUpdateStatusExecutor(configProvider))
	container.Set(types.EventTypeInstructionStarted, executor.NewInstructionExecutor(metricsSender))
	container.Set(types.EventTypeInstructionPaused, executor.NewInstructionExecutor(metricsSender))
	container.Set(types.EventTypeInstructionResumed, executor.NewInstructionExecutor(metricsSender))
	container.Set(types.EventTypeInstructionCompleted, executor.NewInstructionExecutor(metricsSender))
	container.Set(types.EventTypeInstructionAborted, executor.NewInstructionExecutor(metricsSender))
	container.Set(types.EventTypeInstructionFailed, executor.NewInstructionExecutor(metricsSender))
	container.Set(types.EventTypeAutoInstructionStarted, executor.NewInstructionExecutor(metricsSender))
	container.Set(types.EventTypeAutoInstructionCompleted, executor.NewInstructionExecutor(metricsSender))
	container.Set(types.EventTypeAutoInstructionFailed, executor.NewInstructionExecutor(metricsSender))
	container.Set(types.EventTypeAutoInstructionAlreadyRunning, executor.NewInstructionExecutor(metricsSender))
	container.Set(types.EventTypeInstructionJobStarted, executor.NewInstructionExecutor(metricsSender))
	container.Set(types.EventTypeInstructionJobCompleted, executor.NewInstructionExecutor(metricsSender))
	container.Set(types.EventTypeInstructionJobAborted, executor.NewInstructionExecutor(metricsSender))
	container.Set(types.EventTypeInstructionJobFailed, executor.NewInstructionExecutor(metricsSender))
	container.Set(types.EventTypeJunitTestSuiteUpdated, executor.NewJunitExecutor())
	container.Set(types.EventTypeJunitTestCaseUpdated, executor.NewJunitExecutor())

	return executor.NewMongoUpdateExecutor(
		executor.NewCombinedExecutor(container),
		alarm.NewAdapter(dbClient),
	)
}
