package che

import (
	"context"
	"errors"
	"fmt"
	"time"

	apisecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/contextgraph"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	libengine "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/externaldata"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/healthcheck"
	communityimport "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/importcontextgraph"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/statesetting"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/techmetrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/usernotification"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/che/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/depmake"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/bsm/redislock"
	"github.com/rs/zerolog"
)

type DependencyMaker struct {
	depmake.DependencyMaker
}

func NewEngine(
	ctx context.Context,
	options Options,
	primaryDbClient mongo.DbClient,
	secondaryDbClient mongo.DbClient,
	noTimeoutClient mongo.DbClient,
	cfg config.CanopsisConf,
	metricsSender metrics.Sender,
	metricsEntityMetaUpdater metrics.MetaUpdater,
	externalDataContainer *externaldata.GetterContainer,
	timezoneConfigProvider *config.BaseTimezoneConfigProvider,
	templateConfigProvider *config.BaseTemplateConfigProvider,
	logger zerolog.Logger,
) libengine.Engine {
	defer depmake.Catch(logger)

	m := DependencyMaker{}
	alarmConfigProvider := config.NewAlarmConfigProvider(cfg, logger)
	metricsConfigProvider := config.NewMetricsConfigProvider(cfg, logger)
	amqpConnection := m.DepAmqpConnection(logger, cfg)
	amqpChannel := m.DepAMQPChannelPub(amqpConnection)
	entityAdapter := entity.NewAdapter(primaryDbClient)
	redisSession := m.DepRedisSession(ctx, redis.EngineLockStorage, logger, cfg)
	runInfoRedisSession := m.DepRedisSession(ctx, redis.EngineRunInfo, logger, cfg)
	serviceRedisSession := m.DepRedisSession(ctx, redis.EntityServiceStorage, logger, cfg)
	periodicalLockClient := redis.NewLockClient(redisSession)
	templateExecutor := template.NewExecutor(templateConfigProvider, timezoneConfigProvider)
	stateSettingsService := statesetting.NewService(primaryDbClient, logger)
	contextGraphManager := contextgraph.NewManager(entityAdapter, primaryDbClient, contextgraph.NewEntityServiceStorage(primaryDbClient), stateSettingsService, logger)
	eventFilterEventCounter := eventfilter.NewEventCounter(primaryDbClient,
		utils.MinDuration(canopsis.DefaultFlushInterval, options.PeriodicalWaitTime), logger)
	notifStore := usernotification.NewStore(primaryDbClient, amqpChannel, json.NewEncoder(),
		canopsis.ApiNotificationExchangeName, "", canopsis.JsonContentType)
	eventFilterFailureService := eventfilter.NewFailureService(primaryDbClient, notifStore,
		utils.MinDuration(canopsis.DefaultFlushInterval, options.PeriodicalWaitTime), apisecurity.ObjEventFilter, logger)
	techMetricsConfigProvider := config.NewTechMetricsConfigProvider(cfg, logger)
	techMetricsSender := techmetrics.NewSender(canopsis.CheEngineName+"/"+utils.NewID(), techMetricsConfigProvider, canopsis.TechMetricsFlushInterval,
		cfg.Global.ReconnectRetries, cfg.Global.GetReconnectTimeout(), logger)

	ruleApplicatorContainer := eventfilter.NewRuleApplicatorContainer()
	ruleApplicatorContainer.Set(eventfilter.RuleTypeChangeEntity, eventfilter.NewChangeEntityApplicator(
		externalDataContainer,
		eventFilterFailureService,
		templateExecutor,
	))
	ruleApplicatorContainer.Set(eventfilter.RuleTypeEnrichment, eventfilter.NewEnrichmentApplicator(
		externalDataContainer,
		eventfilter.NewActionProcessor(alarmConfigProvider, eventFilterFailureService, templateExecutor),
		eventFilterFailureService,
	))
	ruleApplicatorContainer.Set(eventfilter.RuleTypeDrop, eventfilter.NewDropApplicator())
	ruleApplicatorContainer.Set(eventfilter.RuleTypeBreak, eventfilter.NewBreakApplicator())

	ruleAdapter := eventfilter.NewRuleAdapter(primaryDbClient)
	eventFilterService := eventfilter.NewRuleService(ruleAdapter, ruleApplicatorContainer, eventFilterEventCounter,
		eventFilterFailureService, templateExecutor, logger)

	healthCheckCfg, err := config.NewHealthCheckAdapter(primaryDbClient).GetConfig(ctx)
	if err != nil {
		panic(fmt.Errorf("cannot load healthcheck config: %w", err))
	}

	runInfoPeriodicalWorker := libengine.NewRunInfoPeriodicalWorker(
		healthCheckCfg.ParseUpdateInterval(logger),
		libengine.NewRunInfoManager(runInfoRedisSession),
		libengine.NewInstanceRunInfo(canopsis.CheEngineName, canopsis.CheQueuePrefix, canopsis.AxeQueuePrefix, []string{
			canopsis.CheExternalQueueName,
			canopsis.CheSystemQueueName,
			canopsis.CheUserQueueName,
		}),
		amqpChannel,
		logger,
	)

	infosDictLockedPeriodicalWorker := libengine.NewLockedPeriodicalWorker(
		periodicalLockClient,
		redis.CheEntityInfosDictionaryPeriodicalLockKey,
		NewInfosDictionaryPeriodicalWorker(secondaryDbClient, options.InfosDictionaryWaitTime, logger),
		logger,
	)

	eventfilterIntervalsWorker := NewEventfilterIntervalsWorker(primaryDbClient, timezoneConfigProvider, options.PeriodicalWaitTime, logger)
	eventfilterIntervalsPeriodicalWorker := libengine.NewLockedPeriodicalWorker(
		periodicalLockClient,
		redis.CheEventFiltersIntervalsPeriodicalLockKey,
		eventfilterIntervalsWorker,
		logger,
	)

	engine := libengine.New(
		func(ctx context.Context) error {
			runInfoPeriodicalWorker.Work(ctx)
			eventfilterIntervalsWorker.Work(ctx)

			// run in goroutine because it may take some time to process heavy dbs, don't want to slow down the engine startup
			go infosDictLockedPeriodicalWorker.Work(ctx)

			_, err := periodicalLockClient.Obtain(ctx, redis.ChePeriodicalLockKey,
				options.PeriodicalWaitTime, &redislock.Options{
					RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(1*time.Second), 1),
				})
			if err != nil {
				// Lock is set for options.PeriodicalWaitTime TTL, other instances should skip actions below
				if errors.Is(err, redislock.ErrNotObtained) {
					return nil
				}

				return fmt.Errorf("cannot obtain lock: %w", err)
			}

			// Below are actions locked with ChePeriodicalLockKey for multi-instance configuration
			err = contextGraphManager.UpdateImpactedServicesFromDependencies(ctx)
			if err != nil {
				logger.Warn().Err(err).Msg("error while recomputing impacted services for connectors")
			}

			return nil
		},
		func(ctx context.Context) {
			err := amqpConnection.Close()
			if err != nil {
				logger.Error().Err(err).Msg("failed to close amqp connection")
			}

			err = redisSession.Close()
			if err != nil {
				logger.Error().Err(err).Msg("failed to close redis connection")
			}

			err = serviceRedisSession.Close()
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

	engine.AddRoutine(func(ctx context.Context) error {
		eventFilterEventCounter.Run(ctx)

		return nil
	})
	engine.AddRoutine(func(ctx context.Context) error {
		eventFilterFailureService.Run(ctx)

		return nil
	})
	engine.AddRoutine(func(ctx context.Context) error {
		techMetricsSender.Run(ctx)

		return nil
	})

	eventProcessor := event.NewProcessorContainer()
	eventProcessor.Set(types.SourceTypeResource, event.NewResourceProcessor(primaryDbClient, contextGraphManager, eventFilterService, metricsSender, logger))
	eventProcessor.Set(types.SourceTypeComponent, event.NewComponentProcessor(primaryDbClient, contextGraphManager, eventFilterService, metricsSender, logger))
	eventProcessor.Set(types.SourceTypeConnector, event.NewConnectorProcessor(primaryDbClient, contextGraphManager, eventFilterService, metricsSender))
	eventProcessor.Set(types.SourceTypeService, event.NewServiceProcessor(primaryDbClient, contextGraphManager, eventFilterService, metricsSender))

	mainMessageProcessor := &messageProcessor{
		FeaturePrintEventOnError: options.PrintEventOnError,
		AlarmConfigProvider:      alarmConfigProvider,
		MetricsConfigProvider:    metricsConfigProvider,
		TechMetricsSender:        techMetricsSender,
		MetricsSender:            metricsSender,
		AmqpPublisher:            m.DepAMQPChannelPub(amqpConnection),
		MetaUpdater:              metricsEntityMetaUpdater,
		EntityCollection:         primaryDbClient.Collection(mongo.EntityMongoCollection),
		EventProcessorContainer:  eventProcessor,
		Encoder:                  json.NewEncoder(),
		Decoder:                  json.NewDecoder(),
		Logger:                   logger,
	}
	engine.AddConsumer(libengine.NewConcurrentConsumer(
		canopsis.CheExternalConsumerName,
		canopsis.CheExternalQueueName,
		cfg.Global.PrefetchCount,
		cfg.Global.PrefetchSize,
		options.Purge,
		canopsis.EngineExchangeName,
		canopsis.AxeExternalQueueName,
		options.FifoAckExchange,
		canopsis.FIFOAckQueueName,
		options.ExternalWorkers,
		false,
		amqpConnection,
		mainMessageProcessor,
		logger,
	))
	engine.AddConsumer(libengine.NewConcurrentConsumer(
		canopsis.CheSystemConsumerName,
		canopsis.CheSystemQueueName,
		cfg.Global.PrefetchCount,
		cfg.Global.PrefetchSize,
		options.Purge,
		canopsis.EngineExchangeName,
		canopsis.AxeSystemQueueName,
		options.FifoAckExchange,
		canopsis.FIFOAckQueueName,
		options.SystemWorkers,
		false,
		amqpConnection,
		mainMessageProcessor,
		logger,
	))
	engine.AddConsumer(libengine.NewConcurrentConsumer(
		canopsis.CheUserConsumerName,
		canopsis.CheUserQueueName,
		cfg.Global.PrefetchCount,
		cfg.Global.PrefetchSize,
		options.Purge,
		canopsis.EngineExchangeName,
		canopsis.AxeUserQueueName,
		options.FifoAckExchange,
		canopsis.FIFOAckQueueName,
		options.UserWorkers,
		false,
		amqpConnection,
		mainMessageProcessor,
		logger,
	))
	engine.AddPeriodicalWorker("soft_delete", libengine.NewLockedPeriodicalWorker(
		periodicalLockClient,
		redis.CheSoftDeletePeriodicalLockKey,
		&softDeletePeriodicalWorker{
			entityCollection:          primaryDbClient.Collection(mongo.EntityMongoCollection),
			serviceCountersCollection: primaryDbClient.Collection(mongo.EntityCountersCollection),
			periodicalInterval:        options.PeriodicalWaitTime,
			eventPublisher:            communityimport.NewEventPublisher(canopsis.DefaultExchangeName, canopsis.FIFOQueueName, json.NewEncoder(), canopsis.JsonContentType, amqpChannel),
			softDeleteWaitTime:        options.SoftDeleteWaitTime,
			logger:                    logger,
		},
		logger,
	))
	engine.AddPeriodicalWorker("eventfilter_intervals", eventfilterIntervalsPeriodicalWorker)
	engine.AddPeriodicalWorker("run_info", runInfoPeriodicalWorker)
	engine.AddPeriodicalWorker("config", libengine.NewLoadConfigPeriodicalWorker(
		options.PeriodicalWaitTime,
		config.NewAdapter(primaryDbClient),
		logger,
		alarmConfigProvider,
		metricsConfigProvider,
		timezoneConfigProvider,
		techMetricsConfigProvider,
		templateConfigProvider,
	))
	engine.AddPeriodicalWorker("impacted_services", libengine.NewLockedPeriodicalWorker(
		periodicalLockClient,
		redis.ChePeriodicalLockKey,
		&impactedServicesPeriodicalWorker{
			Manager:            contextGraphManager,
			PeriodicalInterval: options.PeriodicalWaitTime,
			Logger:             logger,
		},
		logger,
	))
	engine.AddPeriodicalWorker("entity_infos_dictionary", infosDictLockedPeriodicalWorker)
	engine.AddRoutine(func(ctx context.Context) error {
		w := eventfilter.NewRulesChangesWatcher(noTimeoutClient, eventFilterService)

		logger.Debug().Msg("Loading event filter rules")

		for {
			select {
			case <-ctx.Done():
				return nil
			default:
				err := w.Watch(ctx, []string{eventfilter.RuleTypeDrop, eventfilter.RuleTypeEnrichment, eventfilter.RuleTypeBreak})
				if err != nil {
					logger.Error().Err(err).Msg("failed to watch eventfilter collection")
				}
			}
		}
	})
	engine.AddRoutine(func(ctx context.Context) error {
		w := statesetting.NewRulesChangesWatcher(noTimeoutClient, stateSettingsService)

		logger.Debug().Msg("Loading state settings rules")

		err := w.Watch(ctx)
		if err != nil {
			logger.Error().Err(err).Msg("failed to state settings collection")
			return err
		}

		return nil
	})

	healthcheck.Start(ctx, healthcheck.NewChecker(
		"che",
		mainMessageProcessor,
		json.NewEncoder(),
		false,
		false,
	), logger)

	return engine
}
