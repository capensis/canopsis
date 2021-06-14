package main

import (
	"context"
	libamqp "git.canopsis.net/canopsis/go-engines/lib/amqp"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/encoding/json"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/engine"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/heartbeat"
	"git.canopsis.net/canopsis/go-engines/lib/depmake"
	libredis "git.canopsis.net/canopsis/go-engines/lib/redis"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
)

type DependencyMaker struct {
	depmake.DependencyMaker
}

type References struct {
	Redis          *redis.Client
	ChannelPub     libamqp.Channel
	ChannelSub     libamqp.Channel
	Adapter        heartbeat.Adapter
	Logger         zerolog.Logger
	JSONEncoder    encoding.Encoder
	RunInfoManager engine.RunInfoManager
}

func (m DependencyMaker) GetDefaultReferences(ctx context.Context, logger zerolog.Logger) References {
	defer depmake.Catch(logger)

	cfg := m.DepConfig()

	amqpConnection := m.DepAmqpConnection(logger, cfg)
	channelPub := m.DepAMQPChannelPub(amqpConnection)
	channelSub := m.DepAMQPChannelSub(amqpConnection, cfg.Global.PrefetchCount, cfg.Global.PrefetchSize)
	redisClient := m.DepRedisSession(ctx, 0, logger, cfg)

	dbClient := m.DepMongoClient(cfg)
	heartbeatAdapter := heartbeat.NewAdapter(dbClient)

	return References{
		ChannelPub:     channelPub,
		ChannelSub:     channelSub,
		Redis:          redisClient,
		Adapter:        heartbeatAdapter,
		Logger:         logger,
		JSONEncoder:    json.NewEncoder(),
		RunInfoManager: engine.NewRunInfoManager(m.DepRedisSession(ctx, libredis.EngineRunInfo, logger, cfg)),
	}
}

func NewEngineHeartBeat(references References) *EngineHeartBeat {
	return &EngineHeartBeat{
		DefaultEngine: canopsis.NewDefaultEngine(canopsis.PeriodicalWaitTime, true, true, references.ChannelSub, references.Logger, references.RunInfoManager),
		References:    references,
	}
}
