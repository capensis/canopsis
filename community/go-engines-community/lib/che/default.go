package che

import (
	"context"
	"fmt"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	libcontext "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	libengine "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
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
	metricsEntityMetaUpdater metrics.MetaUpdater,
	externalDataContainer *eventfilter.ExternalDataContainer,
	timezoneConfigProvider *config.BaseTimezoneConfigProvider,
	templateConfigProvider *config.BaseTemplateConfigProvider,
	logger zerolog.Logger,
) libengine.Engine {
	defer depmake.Catch(logger)

	m := DependencyMaker{}
	alarmConfigProvider := config.NewAlarmConfigProvider(cfg, logger)
	amqpConnection := m.DepAmqpConnection(logger, cfg)
	amqpChannel := m.DepAMQPChannelPub(amqpConnection)
	entityAdapter := entity.NewAdapter(mongoClient)
	entityServiceAdapter := entityservice.NewAdapter(mongoClient)
	redisSession := m.DepRedisSession(ctx, redis.EngineLockStorage, logger, cfg)
	runInfoRedisSession := m.DepRedisSession(ctx, redis.EngineRunInfo, logger, cfg)
	serviceRedisSession := m.DepRedisSession(ctx, redis.EntityServiceStorage, logger, cfg)
	periodicalLockClient := redis.NewLockClient(redisSession)
	templateExecutor := template.NewExecutor(templateConfigProvider, timezoneConfigProvider)

	enrichmentCenter := libcontext.NewEnrichmentCenter(
		entityAdapter,
		mongoClient,
		entityservice.NewManager(
			entityServiceAdapter,
			entityservice.NewStorage(entityServiceAdapter, serviceRedisSession, json.NewEncoder(), json.NewDecoder(), logger),
			logger,
		),
		metricsEntityMetaUpdater,
	)

	techMetricsConfigProvider := config.NewTechMetricsConfigProvider(cfg, logger)
	techMetricsSender := techmetrics.NewSender(techMetricsConfigProvider, canopsis.TechMetricsFlushInterval,
		cfg.Global.ReconnectRetries, cfg.Global.GetReconnectTimeout(), logger)

	ruleApplicatorContainer := eventfilter.NewRuleApplicatorContainer()
	ruleApplicatorContainer.Set(eventfilter.RuleTypeChangeEntity, eventfilter.NewChangeEntityApplicator(externalDataContainer, templateExecutor))
	ruleApplicatorContainer.Set(eventfilter.RuleTypeEnrichment, eventfilter.NewEnrichmentApplicator(externalDataContainer, eventfilter.NewActionProcessor(templateExecutor, techMetricsSender)))
	ruleApplicatorContainer.Set(eventfilter.RuleTypeDrop, eventfilter.NewDropApplicator())
	ruleApplicatorContainer.Set(eventfilter.RuleTypeBreak, eventfilter.NewBreakApplicator())

	ruleAdapter := eventfilter.NewRuleAdapter(mongoClient)

	eventfilterService := eventfilter.NewRuleService(ruleAdapter, ruleApplicatorContainer, logger)

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
				err := eventfilterService.LoadRules(ctx, []string{eventfilter.RuleTypeDrop, eventfilter.RuleTypeEnrichment, eventfilter.RuleTypeBreak})
				if err != nil {
					return fmt.Errorf("unable to load rules: %w", err)
				}
			}

			err := enrichmentCenter.LoadServices(ctx)
			if err != nil {
				return fmt.Errorf("unable to load services: %w", err)
			}

			_, err = periodicalLockClient.Obtain(ctx, redis.ChePeriodicalLockKey,
				options.PeriodicalWaitTime, &redislock.Options{
					RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(1*time.Second), 1),
				})
			if err != nil {
				// Lock is set for options.PeriodicalWaitTime TTL, other instances should skip actions below
				if err == redislock.ErrNotObtained {
					return nil
				}

				return fmt.Errorf("cannot obtain lock: %w", err)
			}

			// Below are actions locked with ChePeriodicalLockKey for multi-instance configuration
			err = enrichmentCenter.UpdateImpactedServices(ctx)
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
		techMetricsSender.Run(ctx)
		return nil
	})

	engine.AddConsumer(libengine.NewDefaultConsumer(
		canopsis.CheConsumerName,
		options.ConsumeFromQueue,
		cfg.Global.PrefetchCount,
		cfg.Global.PrefetchSize,
		options.Purge,
		"",
		options.PublishToQueue,
		options.FifoAckExchange,
		canopsis.FIFOAckQueueName,
		amqpConnection,
		&messageProcessor{
			FeaturePrintEventOnError: options.PrintEventOnError,
			FeatureEventProcessing:   options.FeatureEventProcessing,
			FeatureContextCreation:   options.FeatureContextCreation,

			AlarmConfigProvider: alarmConfigProvider,
			EventFilterService:  eventfilterService,
			EnrichmentCenter:    enrichmentCenter,
			TechMetricsSender:   techMetricsSender,
			AmqpPublisher:       m.DepAMQPChannelPub(amqpConnection),
			AlarmAdapter:        alarm.NewAdapter(mongoClient),
			Encoder:             json.NewEncoder(),
			Decoder:             json.NewDecoder(),
			Logger:              logger,
		},
		logger,
	))
	engine.AddPeriodicalWorker("local cache", &reloadLocalCachePeriodicalWorker{
		EventFilterService: eventfilterService,
		EnrichmentCenter:   enrichmentCenter,
		PeriodicalInterval: options.PeriodicalWaitTime,
		Logger:             logger,
		LoadRules:          !mongoClient.IsDistributed(),
	})
	engine.AddPeriodicalWorker("soft delete", libengine.NewLockedPeriodicalWorker(
		periodicalLockClient,
		redis.CheSoftDeletePeriodicalLockKey,
		&softDeletePeriodicalWorker{
			collection:         mongoClient.Collection(mongo.EntityMongoCollection),
			periodicalInterval: options.PeriodicalWaitTime,
			eventPublisher:     communityimport.NewEventPublisher(canopsis.FIFOExchangeName, canopsis.FIFOQueueName, json.NewEncoder(), canopsis.JsonContentType, amqpChannel),
			softDeleteWaitTime: options.SoftDeleteWaitTime,
			logger:             logger,
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
	))
	engine.AddPeriodicalWorker("impacted services", libengine.NewLockedPeriodicalWorker(
		periodicalLockClient,
		redis.ChePeriodicalLockKey,
		&impactedServicesPeriodicalWorker{
			EnrichmentCenter:   enrichmentCenter,
			PeriodicalInterval: options.PeriodicalWaitTime,
			Logger:             logger,
		},
		logger,
	))
	engine.AddPeriodicalWorker("entity infos dictionary", infosDictLockedPeriodicalWorker)
	if mongoClient.IsDistributed() {
		engine.AddRoutine(func(ctx context.Context) error {
			w := eventfilter.NewRulesChangesWatcher(mongoClient, eventfilterService)

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

	return engine
}
