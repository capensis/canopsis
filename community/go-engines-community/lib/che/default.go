package che

import (
	"context"
	"errors"
	"fmt"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/contextgraph"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	libengine "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/healthcheck"
	communityimport "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/importcontextgraph"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/techmetrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/depmake"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"github.com/bsm/redislock"
	"github.com/rs/zerolog"
)

type DependencyMaker struct {
	depmake.DependencyMaker
}

func NewEngine(
	ctx context.Context,
	options Options,
	mongoClient mongo.DbClient,
	cfg config.CanopsisConf,
	metricsSender metrics.Sender,
	metricsEntityMetaUpdater metrics.MetaUpdater,
	externalDataContainer *eventfilter.ExternalDataContainer,
	timezoneConfigProvider *config.BaseTimezoneConfigProvider,
	templateConfigProvider *config.BaseTemplateConfigProvider,
	eventFilterEventCounter eventfilter.EventCounter,
	eventFilterFailureService eventfilter.FailureService,
	logger zerolog.Logger,
) libengine.Engine {
	defer depmake.Catch(logger)

	m := DependencyMaker{}
	alarmConfigProvider := config.NewAlarmConfigProvider(cfg, logger)
	amqpConnection := m.DepAmqpConnection(logger, cfg)
	amqpChannel := m.DepAMQPChannelPub(amqpConnection)
	entityAdapter := entity.NewAdapter(mongoClient)
	redisSession := m.DepRedisSession(ctx, redis.EngineLockStorage, logger, cfg)
	runInfoRedisSession := m.DepRedisSession(ctx, redis.EngineRunInfo, logger, cfg)
	serviceRedisSession := m.DepRedisSession(ctx, redis.EntityServiceStorage, logger, cfg)
	periodicalLockClient := redis.NewLockClient(redisSession)
	templateExecutor := template.NewExecutor(templateConfigProvider, timezoneConfigProvider)
	dataStorageConfigProvider := config.NewDataStorageConfigProvider(cfg, logger)
	contextGraphManager := contextgraph.NewManager(entityAdapter, mongoClient, contextgraph.NewEntityServiceStorage(mongoClient), metricsEntityMetaUpdater, logger)

	techMetricsConfigProvider := config.NewTechMetricsConfigProvider(cfg, logger)
	techMetricsSender := techmetrics.NewSender(techMetricsConfigProvider, canopsis.TechMetricsFlushInterval,
		cfg.Global.ReconnectRetries, cfg.Global.GetReconnectTimeout(), logger)

	ruleApplicatorContainer := eventfilter.NewRuleApplicatorContainer()
	ruleApplicatorContainer.Set(eventfilter.RuleTypeChangeEntity, eventfilter.NewChangeEntityApplicator(
		externalDataContainer,
		eventFilterFailureService,
		templateExecutor,
	))
	ruleApplicatorContainer.Set(eventfilter.RuleTypeEnrichment, eventfilter.NewEnrichmentApplicator(
		externalDataContainer,
		eventfilter.NewActionProcessor(alarmConfigProvider, eventFilterFailureService, templateExecutor, techMetricsSender),
		eventFilterFailureService,
	))
	ruleApplicatorContainer.Set(eventfilter.RuleTypeDrop, eventfilter.NewDropApplicator())
	ruleApplicatorContainer.Set(eventfilter.RuleTypeBreak, eventfilter.NewBreakApplicator())

	ruleAdapter := eventfilter.NewRuleAdapter(mongoClient)
	eventFilterService := eventfilter.NewRuleService(ruleAdapter, ruleApplicatorContainer, eventFilterEventCounter,
		eventFilterFailureService, templateExecutor, logger)

	runInfoPeriodicalWorker := libengine.NewRunInfoPeriodicalWorker(
		options.PeriodicalWaitTime,
		libengine.NewRunInfoManager(runInfoRedisSession),
		libengine.NewInstanceRunInfo(canopsis.CheEngineName, options.ConsumeFromQueue, options.PublishToQueue),
		amqpChannel,
		logger,
	)

	infosDictLockedPeriodicalWorker := libengine.NewLockedPeriodicalWorker(
		periodicalLockClient,
		redis.CheEntityInfosDictionaryPeriodicalLockKey,
		NewInfosDictionaryPeriodicalWorker(mongoClient, options.InfosDictionaryWaitTime, logger),
		logger,
	)

	eventfilterIntervalsWorker := NewEventfilterIntervalsWorker(mongoClient, timezoneConfigProvider, options.PeriodicalWaitTime, logger)
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

			if !mongoClient.IsDistributed() {
				logger.Debug().Msg("Loading event filter rules")
				err := eventFilterService.LoadRules(ctx, []string{eventfilter.RuleTypeDrop, eventfilter.RuleTypeEnrichment, eventfilter.RuleTypeBreak})
				if err != nil {
					return fmt.Errorf("unable to load rules: %w", err)
				}
			}

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

	mainMessageProcessor := &messageProcessor{
		FeaturePrintEventOnError: options.PrintEventOnError,
		DbClient:                 mongoClient,

		AlarmConfigProvider: alarmConfigProvider,
		EventFilterService:  eventFilterService,
		ContextGraphManager: contextGraphManager,
		TechMetricsSender:   techMetricsSender,
		MetricsSender:       metricsSender,
		AmqpPublisher:       m.DepAMQPChannelPub(amqpConnection),
		MetaUpdater:         metricsEntityMetaUpdater,
		EntityCollection:    mongoClient.Collection(mongo.EntityMongoCollection),
		Encoder:             json.NewEncoder(),
		Decoder:             json.NewDecoder(),
		Logger:              logger,
	}
	engine.AddConsumer(libengine.NewConcurrentConsumer(
		canopsis.CheConsumerName,
		options.ConsumeFromQueue,
		cfg.Global.PrefetchCount,
		cfg.Global.PrefetchSize,
		options.Purge,
		"",
		options.PublishToQueue,
		options.FifoAckExchange,
		canopsis.FIFOAckQueueName,
		options.Workers,
		amqpConnection,
		mainMessageProcessor,
		logger,
	))
	engine.AddPeriodicalWorker("local cache", &reloadLocalCachePeriodicalWorker{
		EventFilterService: eventFilterService,
		PeriodicalInterval: options.PeriodicalWaitTime,
		Logger:             logger,
		LoadRules:          !mongoClient.IsDistributed(),
	})
	engine.AddPeriodicalWorker("soft delete", libengine.NewLockedPeriodicalWorker(
		periodicalLockClient,
		redis.CheSoftDeletePeriodicalLockKey,
		&softDeletePeriodicalWorker{
			entityCollection:          mongoClient.Collection(mongo.EntityMongoCollection),
			serviceCountersCollection: mongoClient.Collection(mongo.EntityServiceCountersCollection),
			periodicalInterval:        options.PeriodicalWaitTime,
			eventPublisher:            communityimport.NewEventPublisher(canopsis.FIFOExchangeName, canopsis.FIFOQueueName, json.NewEncoder(), canopsis.JsonContentType, amqpChannel),
			softDeleteWaitTime:        options.SoftDeleteWaitTime,
			logger:                    logger,
		},
		logger,
	))
	engine.AddPeriodicalWorker("eventfilter intervals", eventfilterIntervalsPeriodicalWorker)
	engine.AddPeriodicalWorker("run info", runInfoPeriodicalWorker)
	engine.AddPeriodicalWorker("config", libengine.NewLoadConfigPeriodicalWorker(
		options.PeriodicalWaitTime,
		config.NewAdapter(mongoClient),
		logger,
		alarmConfigProvider,
		timezoneConfigProvider,
		techMetricsConfigProvider,
		templateConfigProvider,
		dataStorageConfigProvider,
	))
	engine.AddPeriodicalWorker("impacted services", libengine.NewLockedPeriodicalWorker(
		periodicalLockClient,
		redis.ChePeriodicalLockKey,
		&impactedServicesPeriodicalWorker{
			Manager:            contextGraphManager,
			PeriodicalInterval: options.PeriodicalWaitTime,
			Logger:             logger,
		},
		logger,
	))
	engine.AddPeriodicalWorker("entity infos dictionary", infosDictLockedPeriodicalWorker)
	if mongoClient.IsDistributed() {
		engine.AddRoutine(func(ctx context.Context) error {
			w := eventfilter.NewRulesChangesWatcher(mongoClient, eventFilterService)

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
	}
	engine.AddPeriodicalWorker("clean", &cleanPeriodicalWorker{
		PeriodicalInterval:        time.Hour,
		TimezoneConfigProvider:    timezoneConfigProvider,
		DataStorageConfigProvider: dataStorageConfigProvider,
		LimitConfigAdapter:        datastorage.NewAdapter(mongoClient),
		Logger:                    logger,
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
