package main

import (
	"time"

	libamqp "git.canopsis.net/canopsis/go-engines/lib/amqp"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/config"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/context"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/encoding/json"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/engine"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/entity"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/watcher"
	"git.canopsis.net/canopsis/go-engines/lib/depmake"
	"git.canopsis.net/canopsis/go-engines/lib/mongo/bulk"
	"git.canopsis.net/canopsis/go-engines/lib/redis"
	"github.com/rs/zerolog"
)

type Options struct {
	FeatureEventProcessing     bool
	FeatureContextCreation     bool
	FeatureContextEnrich       bool
	FeatureAlwaysFlushEntities bool
	Purge                      bool
	PrintEventOnError          bool
	ModeDebug                  bool
	ConsumeFromQueue           string
	PublishToQueue             string
	PublishToExchange          string
	EnrichExclude              string
	EnrichInclude              string
	DataSourceDirectory        string
}

type References struct {
	EnrichFields       context.EnrichFields
	EventFilterService eventfilter.Service
	EnrichmentCenter   context.EnrichmentCenter
	RunInfoManager     engine.RunInfoManager
	ChannelPub         libamqp.Channel
	ChannelSub         libamqp.Channel
	JSONDecoder        encoding.Decoder
	JSONEncoder        encoding.Encoder
	Logger             zerolog.Logger
	Config             config.CanopsisConf
}

type DependencyMaker struct {
	depmake.DependencyMaker
}

func (m DependencyMaker) GetDefaultReferences(options Options, logger zerolog.Logger) References {
	defer depmake.Catch(logger)

	cfg := m.DepConfig()

	mongoSession := m.DepMongoSession()
	amqpSession := m.DepAmqpConnection(logger, cfg)
	channelSub := m.DepAMQPChannelSub(amqpSession, cfg.Global.PrefetchCount, cfg.Global.PrefetchSize)
	channelPub := m.DepAMQPChannelPub(amqpSession)

	entityCollection := entity.DefaultCollection(mongoSession)
	entityAdapter := entity.NewAdapter(entityCollection)

	eventFilterCollection := eventfilter.DefaultCollection(mongoSession)
	eventFilterAdapter := eventfilter.NewAdapter(eventFilterCollection)

	watcherCollection := watcher.DefaultCollection(mongoSession)
	watcherBulk := watcherCollection.NewBulk(bulk.BulkSizeMax)
	watcherAdapter := watcher.NewAdapter(
		watcherCollection,
		watcherBulk,
		entity.EntityCollectionName,
		alarm.AlarmCollectionName,
		logger)

	return References{
		Config:             cfg,
		EventFilterService: eventfilter.NewService(eventFilterAdapter, logger),
		EnrichmentCenter: context.NewSafeEnrichmentCenter(context.NewEnrichmentCenter(
			bulk.BulkSizeMax,
			options.FeatureContextEnrich,
			entityAdapter,
			watcherAdapter,
			logger,
		)),
		RunInfoManager: engine.NewRunInfoManager(m.DepRedisSession(redis.EngineRunInfo, logger, cfg)),
		ChannelSub:     channelSub,
		ChannelPub:     channelPub,
		JSONDecoder:    json.NewDecoder(),
		JSONEncoder:    json.NewEncoder(),
		Logger:         logger,
	}
}

// NewEngineCHE returns *EngineChe. Note that references.EnrichFields
// WILL be replaced by a new instance using options.EnrichExclude and
// options.EnrichInclude values, so creating it yourself is useless.
func NewEngineCHE(options Options, references References) *EngineChe {
	if references.EnrichmentCenter == nil {
		panic("enrichment center is nil")
	}
	if references.ChannelPub == nil {
		panic("publish channel is nil")
	}
	if references.ChannelSub == nil {
		panic("consume channel is nil")
	}
	if references.JSONDecoder == nil {
		panic("json decoder is nil")
	}
	if references.JSONEncoder == nil {
		panic("json encoder is nil")
	}

	references.EnrichFields = context.NewEnrichFields(
		options.EnrichInclude,
		options.EnrichExclude,
	)

	defaultEngine := canopsis.NewDefaultEngine(
		time.Minute,
		true,
		true,
		references.ChannelSub,
		references.Logger,
		references.RunInfoManager,
	)

	defaultEngine.Debug = options.ModeDebug

	return &EngineChe{
		DefaultEngine: defaultEngine,
		Options:       options,
		References:    references,
	}
}
