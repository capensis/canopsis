package che

import (
	"context"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/neweventfilter"
	"path/filepath"
	"plugin"
	"strings"
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
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/depmake"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/postgres"
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
	pgPool postgres.Pool,
	metricsEntityMetaUpdater metrics.MetaUpdater,
	logger zerolog.Logger,
) libengine.Engine {
	defer depmake.Catch(logger)

	m := DependencyMaker{}
	cfg := m.DepConfig(ctx, mongoClient)
	config.SetDbClientRetry(mongoClient, cfg)
	if pgPool != nil {
		config.SetPgPoolRetry(pgPool, cfg)
	}
	alarmConfigProvider := config.NewAlarmConfigProvider(cfg, logger)
	timezoneConfigProvider := config.NewTimezoneConfigProvider(cfg, logger)
	amqpConnection := m.DepAmqpConnection(logger, cfg)
	amqpChannel := m.DepAMQPChannelPub(amqpConnection)
	entityAdapter := entity.NewAdapter(mongoClient)
	eventFilterAdapter := eventfilter.NewAdapter(mongoClient)
	entityServiceAdapter := entityservice.NewAdapter(mongoClient)
	redisSession := m.DepRedisSession(ctx, redis.EngineLockStorage, logger, cfg)
	runInfoRedisSession := m.DepRedisSession(ctx, redis.EngineRunInfo, logger, cfg)
	serviceRedisSession := m.DepRedisSession(ctx, redis.EntityServiceStorage, logger, cfg)
	periodicalLockClient := redis.NewLockClient(redisSession)

	eventFilterService := eventfilter.NewService(eventFilterAdapter, timezoneConfigProvider, logger)
	enrichmentCenter := libcontext.NewEnrichmentCenter(
		entityAdapter,
		options.FeatureContextEnrich,
		entityservice.NewManager(
			entityServiceAdapter,
			entityAdapter,
			entityservice.NewStorage(serviceRedisSession, json.NewEncoder(), json.NewDecoder(), logger),
			logger,
		),
		metricsEntityMetaUpdater,
		logger,
	)
	enrichFields := libcontext.NewEnrichFields(options.EnrichInclude, options.EnrichExclude)

	logger.Debug().Msg("Loading event filter data sources")
	factories, err := LoadDataSourceFactories(options.DataSourceDirectory)
	if err != nil {
		panic(fmt.Errorf("unable to load data sources: %w", err))
	}

	ruleAdapter := neweventfilter.NewRuleAdapter(mongoClient)

	ruleApplicatorContainer := neweventfilter.NewRuleApplicatorContainer()
	ruleApplicatorContainer.Set(neweventfilter.RuleTypeChangeEntity, neweventfilter.NewChangeEntityApplicator(factories))
	ruleApplicatorContainer.Set(neweventfilter.RuleTypeEnrichment, neweventfilter.NewEnrichmentApplicator(factories, neweventfilter.NewActionProcessor()))
	ruleApplicatorContainer.Set(neweventfilter.RuleTypeDrop, neweventfilter.NewDropApplicator())
	ruleApplicatorContainer.Set(neweventfilter.RuleTypeBreak, neweventfilter.NewBreakApplicator())

	newEventFilterService := neweventfilter.NewRuleService(ruleAdapter, ruleApplicatorContainer, logger)

	engine := libengine.New(
		func(ctx context.Context) error {
			logger.Debug().Msg("Loading event filter rules")
			err = newEventFilterService.LoadRules(ctx)
			if err != nil {
				return fmt.Errorf("unable to load rules: %v", err)
			}

			logger.Debug().Msg("Loading services")
			err = enrichmentCenter.LoadServices(ctx)
			if err != nil {
				logger.Error().Err(err).Msg("unable to load services")
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

				logger.Error().Err(err).Msg("cannot obtain lock")
				return err
			}
			// Below are actions locked with ChePeriodicalLockKey for multi-instance configuration

			logger.Debug().Msg("Recompute impacted services for connectors")
			err = enrichmentCenter.UpdateImpactedServices(ctx)
			if err != nil {
				logger.Warn().Err(err).Msg("error while recomputing impacted services for connectors")
			}
			logger.Debug().Msg("Recompute impacted services for connectors finished")
			return nil
		},
		func(ctx context.Context) {
			err := mongoClient.Disconnect(ctx)
			if err != nil {
				logger.Error().Err(err).Msg("failed to close mongo connection")
			}

			err = amqpConnection.Close()
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

			if pgPool != nil {
				pgPool.Close()
			}
		},
		logger,
	)
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
			FeaturePrintEventOnError: false,
			FeatureEventProcessing:   options.FeatureEventProcessing,
			FeatureContextCreation:   options.FeatureContextCreation,
			AlarmConfigProvider:      alarmConfigProvider,
			EventFilterService:       eventFilterService,
			NewEventFilterService:    newEventFilterService,
			EnrichmentCenter:         enrichmentCenter,
			EnrichFields:             enrichFields,
			AmqpPublisher:            m.DepAMQPChannelPub(amqpConnection),
			AlarmAdapter:             alarm.NewAdapter(mongoClient),
			Encoder:                  json.NewEncoder(),
			Decoder:                  json.NewDecoder(),
			Logger:                   logger,
		},
		logger,
	))
	engine.AddPeriodicalWorker(&reloadLocalCachePeriodicalWorker{
		EventFilterService: newEventFilterService,
		EnrichmentCenter:   enrichmentCenter,
		PeriodicalInterval: options.PeriodicalWaitTime,
		Logger:             logger,
	})
	engine.AddPeriodicalWorker(libengine.NewRunInfoPeriodicalWorker(
		options.PeriodicalWaitTime,
		libengine.NewRunInfoManager(runInfoRedisSession),
		libengine.NewInstanceRunInfo(canopsis.CheEngineName, options.ConsumeFromQueue, options.PublishToQueue),
		amqpChannel,
		logger,
	))
	engine.AddPeriodicalWorker(libengine.NewLoadConfigPeriodicalWorker(
		options.PeriodicalWaitTime,
		config.NewAdapter(mongoClient),
		alarmConfigProvider,
		logger,
	))
	engine.AddPeriodicalWorker(libengine.NewLoadConfigPeriodicalWorker(
		options.PeriodicalWaitTime,
		config.NewAdapter(mongoClient),
		timezoneConfigProvider,
		logger,
	))
	engine.AddPeriodicalWorker(libengine.NewLockedPeriodicalWorker(
		periodicalLockClient,
		redis.ChePeriodicalLockKey,
		&impactedServicesPeriodicalWorker{
			EnrichmentCenter:   enrichmentCenter,
			PeriodicalInterval: options.PeriodicalWaitTime,
			Logger:             logger,
		},
		logger,
	))

	return engine
}

//TODO: copy from eventfilter package, all mongo plugin feature should be refactored
func LoadDataSourceFactories(dataSourceDirectory string) (map[string]eventfilter.DataSourceFactory, error) {
	factories := make(map[string]eventfilter.DataSourceFactory)

	files, err := filepath.Glob(filepath.Join(dataSourceDirectory, fmt.Sprintf("*%s", canopsis.PluginExtension)))
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if strings.HasSuffix(file, canopsis.PluginExtension) {
			sourceName := strings.TrimSuffix(filepath.Base(file), canopsis.PluginExtension)
			plug, err := plugin.Open(file)
			if err != nil {
				return nil, fmt.Errorf("unable to open plugin: %w", err)
			}

			factorySymbol, err := plug.Lookup("DataSourceFactory")
			if err != nil {
				return nil, fmt.Errorf("unable to load plugin: %w", err)
			}

			factory, isFactory := factorySymbol.(eventfilter.DataSourceFactory)
			if !isFactory {
				return nil, fmt.Errorf("the plugin does not define a valid data source")
			}

			factories[sourceName] = factory
		}
	}

	return factories, nil
}
