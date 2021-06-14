package main

import (
	"time"

	libamqp "git.canopsis.net/canopsis/go-engines/lib/amqp"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/encoding/json"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/engine"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/scheduler"
	"git.canopsis.net/canopsis/go-engines/lib/depmake"
	"git.canopsis.net/canopsis/go-engines/lib/redis"
	"github.com/rs/zerolog"
)

type Options struct {
	PrintEventOnError         bool
	ModeDebug                 bool
	ConsumeFromQueue          string
	PublishToQueue            string
	LockTtl                   int
	EnableMetaAlarmProcessing bool
}

type References struct {
	Scheduler      scheduler.Scheduler
	RunInfoManager engine.RunInfoManager
	ChannelPub     libamqp.Channel
	ChannelSub     libamqp.Channel
	AckChanSub     libamqp.Channel
	AckChanPub     libamqp.Channel
	JSONDecoder    encoding.Decoder
	Logger         zerolog.Logger
}

type DependencyMaker struct {
	depmake.DependencyMaker
}

func (m DependencyMaker) GetDefaultReferences(options Options, logger zerolog.Logger) References {
	defer depmake.Catch(logger)

	cfg := m.DepConfig()

	channelSub := m.DepAMQPChannelSub(m.DepAmqpConnection(logger, cfg), cfg.Global.PrefetchCount, cfg.Global.PrefetchSize)
	ackChanSub := m.DepAMQPChannelSub(m.DepAmqpConnection(logger, cfg), cfg.Global.PrefetchCount, cfg.Global.PrefetchSize)
	channelPub := m.DepAMQPChannelPub(m.DepAmqpConnection(logger, cfg))
	ackChanPub := m.DepAMQPChannelPub(m.DepAmqpConnection(logger, cfg))

	redisLockStorage := m.DepRedisSession(redis.LockStorage, logger, cfg)
	redisQueueStorage := m.DepRedisSession(redis.QueueStorage, logger, cfg)

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

	return References{
		Scheduler:      eventScheduler,
		RunInfoManager: engine.NewRunInfoManager(m.DepRedisSession(redis.EngineRunInfo, logger, cfg)),
		ChannelSub:     channelSub,
		ChannelPub:     channelPub,
		AckChanSub:     ackChanSub,
		AckChanPub:     ackChanPub,
		JSONDecoder:    jsonDecoder,
		Logger:         logger,
	}
}

func NewEngineFIFO(options Options, references References) *EngineFIFO {
	defaultEngine := canopsis.NewDefaultEngine(
		time.Second*60,
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
