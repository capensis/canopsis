package main

import (
	"context"
	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/ratelimit"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/scheduler"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/statistics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/depmake"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"github.com/rs/zerolog"
	"time"
)

type Options struct {
	PrintEventOnError         bool
	ModeDebug                 bool
	ConsumeFromQueue          string
	PublishToQueue            string
	LockTtl                   int
	EnableMetaAlarmProcessing bool
	EventsStatsFlushInterval  time.Duration
}

type References struct {
	Scheduler      scheduler.Scheduler
	RunInfoManager engine.RunInfoManager
	ChannelPub     libamqp.Channel
	ChannelSub     libamqp.Channel
	AckChanSub     libamqp.Channel
	AckChanPub     libamqp.Channel
	JSONDecoder    encoding.Decoder
	StatsSender    ratelimit.StatsSender
	StatsListener  statistics.StatsListener
	StatsCh        chan statistics.Message
	Logger         zerolog.Logger
}

type DependencyMaker struct {
	depmake.DependencyMaker
}

func (m DependencyMaker) GetDefaultReferences(ctx context.Context, options Options, logger zerolog.Logger) References {
	defer depmake.Catch(logger)

	cfg := m.DepConfig()

	dbClient, err := mongo.NewClient(cfg.Global.ReconnectRetries, cfg.Global.GetReconnectTimeout())
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot connect to mongodb")
		panic(err)
	}

	channelSub := m.DepAMQPChannelSub(m.DepAmqpConnection(logger, cfg), cfg.Global.PrefetchCount, cfg.Global.PrefetchSize)
	ackChanSub := m.DepAMQPChannelSub(m.DepAmqpConnection(logger, cfg), cfg.Global.PrefetchCount, cfg.Global.PrefetchSize)
	channelPub := m.DepAMQPChannelPub(m.DepAmqpConnection(logger, cfg))
	ackChanPub := m.DepAMQPChannelPub(m.DepAmqpConnection(logger, cfg))

	redisLockStorage := m.DepRedisSession(ctx, redis.LockStorage, logger, cfg)
	redisQueueStorage := m.DepRedisSession(ctx, redis.QueueStorage, logger, cfg)
	statsRedisClient := m.DepRedisSession(ctx, redis.FIFOMessageStatisticsStorage, logger, cfg)

	jsonDecoder := json.NewDecoder()

	eventScheduler := scheduler.NewSchedulerService(
		redisLockStorage,
		redisQueueStorage,
		channelPub, options.PublishToQueue,
		logger,
		options.LockTtl,
		jsonDecoder,
		options.EnableMetaAlarmProcessing,
	)

	statsCh := make(chan statistics.Message)
	statsSender := ratelimit.NewStatsSender(statsCh, logger)

	statsListener := statistics.NewStatsListener(
		dbClient,
		statsRedisClient,
		options.EventsStatsFlushInterval,
		map[string]int64{
			mongo.MessageRateStatsMinuteCollectionName: 1,  // 1 minute
			mongo.MessageRateStatsHourCollectionName:   60, // 60 minutes
		},
		logger,
	)

	return References{
		Scheduler:      eventScheduler,
		RunInfoManager: engine.NewRunInfoManager(m.DepRedisSession(ctx, redis.EngineRunInfo, logger, cfg)),
		ChannelSub:     channelSub,
		ChannelPub:     channelPub,
		AckChanSub:     ackChanSub,
		AckChanPub:     ackChanPub,
		JSONDecoder:    jsonDecoder,
		StatsSender:    statsSender,
		StatsListener:  statsListener,
		StatsCh:        statsCh,
		Logger:         logger,
	}
}

func NewEngineFIFO(options Options, references References) *EngineFIFO {
	defaultEngine := canopsis.NewDefaultEngine(
		canopsis.PeriodicalWaitTime,
		true,
		true,
		references.ChannelSub,
		references.Logger,
		references.RunInfoManager,
	)

	defaultEngine.Debug = options.ModeDebug

	return &EngineFIFO{
		DefaultEngine: defaultEngine,
		Options:       options,
		References:    references,
	}
}
