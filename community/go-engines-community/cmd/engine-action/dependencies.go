package main

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/action"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/healthcheck"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/techmetrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/depmake"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"github.com/rs/zerolog"
)

type Options struct {
	ModeDebug                bool
	FifoAckQueue             string
	FifoAckExchange          string
	FeaturePrintEventOnError bool
	PeriodicalWaitTime       time.Duration
	WorkerPoolSize           int
	WithWebhook              bool
	LastRetryInterval        time.Duration
}

// DependencyMaker can be created with DependencyMaker{}
type DependencyMaker struct {
	depmake.DependencyMaker
}

// NewEngineAction returns the default Action engine with default connections.
func NewEngineAction(ctx context.Context, options Options, logger zerolog.Logger) engine.Engine {
	m := DependencyMaker{}

	mongoClient := m.DepMongoClient(ctx, logger)
	cfg := m.DepConfig(ctx, mongoClient)
	config.SetDbClientRetry(mongoClient, cfg)
	templateConfigProvider := config.NewTemplateConfigProvider(cfg, logger)
	timezoneConfigProvider := config.NewTimezoneConfigProvider(cfg, logger)
	amqpConnection := m.DepAmqpConnection(logger, cfg)
	amqpChannel := m.DepAMQPChannelPub(amqpConnection)
	actionAdapter := action.NewAdapter(mongoClient)
	alarmAdapter := alarm.NewAdapter(mongoClient)
	actionRedisClient := m.DepRedisSession(ctx, redis.ActionScenarioStorage, logger, cfg)
	runInfoRedisClient := m.DepRedisSession(ctx, redis.EngineRunInfo, logger, cfg)
	lockRedisClient := m.DepRedisSession(ctx, redis.EngineLockStorage, logger, cfg)
	delayedScenarioManager := action.NewDelayedScenarioManager(actionAdapter, alarmAdapter,
		action.NewRedisDelayedScenarioStorage(redis.ActionDelayedScenarioKey, actionRedisClient, json.NewEncoder(), json.NewDecoder()),
		options.PeriodicalWaitTime, logger)
	scenarioExecChan := make(chan action.ExecuteScenariosTask)
	storage := action.NewRedisScenarioExecutionStorage(redis.ActionScenarioExecutionKey, actionRedisClient, json.NewEncoder(),
		json.NewDecoder(), options.LastRetryInterval, logger)
	actionScenarioStorage := action.NewScenarioStorage(actionAdapter, delayedScenarioManager, logger)
	actionService := action.NewService(alarmAdapter, scenarioExecChan,
		delayedScenarioManager, storage, json.NewEncoder(), json.NewDecoder(), amqpChannel,
		options.FifoAckExchange, options.FifoAckQueue,
		alarm.NewActivationService(json.NewEncoder(), amqpChannel, canopsis.CheQueueName, logger),
		logger)
	templateExecutor := template.NewExecutor(templateConfigProvider, timezoneConfigProvider)

	rpcResultChannel := make(chan action.RpcResult)

	axeRpcClient := engine.NewRPCClient(
		canopsis.ActionRPCConsumerName,
		canopsis.AxeRPCQueueServerName,
		canopsis.ActionAxeRPCClientQueueName,
		cfg.Global.PrefetchCount,
		cfg.Global.PrefetchSize,
		&axeRpcClientMessageProcessor{
			FeaturePrintEventOnError: options.FeaturePrintEventOnError,
			Decoder:                  json.NewDecoder(),
			Logger:                   logger,
			ResultChannel:            rpcResultChannel,
		},
		amqpChannel,
		logger,
	)
	var webhookRpcClient engine.RPCClient
	if options.WithWebhook {
		webhookRpcClient = engine.NewRPCClient(
			canopsis.ActionRPCConsumerName,
			canopsis.WebhookRPCQueueServerName,
			"",
			cfg.Global.PrefetchCount,
			cfg.Global.PrefetchSize,
			nil,
			amqpChannel,
			logger,
		)
	}
	rpcPublishQueues := make([]string, 0)
	if webhookRpcClient != nil {
		rpcPublishQueues = append(rpcPublishQueues, canopsis.WebhookRPCQueueServerName)
	}
	runInfoPeriodicalWorker := engine.NewRunInfoPeriodicalWorker(
		options.PeriodicalWaitTime,
		engine.NewRunInfoManager(runInfoRedisClient),
		engine.NewInstanceRunInfo(canopsis.ActionEngineName, canopsis.ActionQueueName, "", nil, rpcPublishQueues),
		amqpChannel,
		logger,
	)

	engineAction := engine.New(
		func(ctx context.Context) error {
			runInfoPeriodicalWorker.Work(ctx)
			manager := action.NewTaskManager(
				action.NewWorkerPool(options.WorkerPoolSize, mongoClient, axeRpcClient, webhookRpcClient, json.NewEncoder(), logger, templateExecutor),
				storage,
				actionScenarioStorage,
				logger,
			)
			scenarioResultChannel, err := manager.Run(ctx, rpcResultChannel, scenarioExecChan)
			if err != nil {
				logger.Error().Err(err).Msg("Initialize: failed to run task manager! Engine will be stopped.")
				return err
			}

			err = actionScenarioStorage.ReloadScenarios(ctx)
			if err != nil {
				logger.Error().Err(err).Msg("Initialize: failed to load actions! Engine will be stopped.")
				return err
			}

			actionService.ListenScenarioFinish(ctx, scenarioResultChannel)

			delayedActionCh, err := delayedScenarioManager.Run(ctx)
			if err != nil {
				logger.Error().Err(err).Msg("Initialize: run delayed scenario manager")
				return err
			}

			listener := &delayedScenarioListener{
				PeriodicalInterval:     options.PeriodicalWaitTime,
				DelayedScenarioManager: delayedScenarioManager,
				AmqpChannel:            amqpChannel,
				Queue:                  canopsis.FIFOQueueName,
				Encoder:                json.NewEncoder(),
				Logger:                 logger,
			}
			go listener.Listen(ctx, delayedActionCh)

			return nil
		},
		func(ctx context.Context) {
			close(scenarioExecChan)
			close(rpcResultChannel)

			err := mongoClient.Disconnect(ctx)
			if err != nil {
				logger.Error().Err(err).Msg("failed to close mongo connection")
			}

			err = amqpConnection.Close()
			if err != nil {
				logger.Error().Err(err).Msg("failed to close amqp connection")
			}

			err = actionRedisClient.Close()
			if err != nil {
				logger.Error().Err(err).Msg("failed to close redis connection")
			}

			err = lockRedisClient.Close()
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

	engineAction.AddRoutine(func(ctx context.Context) error {
		techMetricsSender.Run(ctx)
		return nil
	})

	mainMessageProcessor := &messageProcessor{
		TechMetricsSender:        techMetricsSender,
		FeaturePrintEventOnError: options.FeaturePrintEventOnError,
		ActionService:            actionService,
		Decoder:                  json.NewDecoder(),
		Logger:                   logger,
	}
	engineAction.AddConsumer(engine.NewDefaultConsumer(
		canopsis.ActionConsumerName,
		canopsis.ActionQueueName,
		cfg.Global.PrefetchCount,
		cfg.Global.PrefetchSize,
		false,
		"",
		"",
		"",
		"",
		amqpConnection,
		mainMessageProcessor,
		logger,
	))
	engineAction.AddConsumer(axeRpcClient)
	engineAction.AddPeriodicalWorker("run info", runInfoPeriodicalWorker)
	engineAction.AddPeriodicalWorker("local cache", &reloadLocalCachePeriodicalWorker{
		PeriodicalInterval:    options.PeriodicalWaitTime,
		ActionScenarioStorage: actionScenarioStorage,
		Logger:                logger,
	})
	engineAction.AddPeriodicalWorker("abandon executions", engine.NewLockedPeriodicalWorker(
		redis.NewLockClient(lockRedisClient),
		redis.ActionPeriodicalLockKey,
		&scenarioPeriodicalWorker{
			PeriodicalInterval: options.PeriodicalWaitTime,
			ActionService:      actionService,
			Logger:             logger,
		},
		logger,
	))
	engineAction.AddPeriodicalWorker("config", engine.NewLoadConfigPeriodicalWorker(
		options.PeriodicalWaitTime,
		config.NewAdapter(mongoClient),
		logger,
		timezoneConfigProvider,
		techMetricsConfigProvider,
		templateConfigProvider,
	))

	healthcheck.Start(ctx, healthcheck.NewChecker(
		"action",
		mainMessageProcessor,
		json.NewEncoder(),
		true,
		true,
	), logger)

	return engineAction
}
