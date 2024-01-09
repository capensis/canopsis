package depmake

import (
	"context"
	"fmt"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	redismod "github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

// DependencyMaker is just a handling struct and can be initialized empty.
// DO NOT use this type in any other place than a main package with a main()
// func, under any circumstances.
//
// Every single function available MUST call depmake.Panic() instead of
// managing manually the error.
//
// The idea behind that is simple: what you can't get makes you stop right now.
type DependencyMaker struct {
	// This struct's methods MUST be like class static methods:
	// no internal state or members.
	// Final structs using this one MAY do so, while not recommended.
}

// DepMongoClient opens a mongo session.
func (m DependencyMaker) DepMongoClient(ctx context.Context, logger zerolog.Logger) mongo.DbClient {
	c, err := mongo.NewClient(ctx, 0, 0, logger)
	Panic("mongo session", err)
	return c
}

// DepConfig gets a config from mongodb
func (m DependencyMaker) DepConfig(ctx context.Context, dbClient mongo.DbClient) config.CanopsisConf {
	cfg, err := config.NewAdapter(dbClient).GetConfig(ctx)
	if err != nil {
		panic(fmt.Errorf("dependency error: %s: %w", "can't get the config", err))
	}

	return cfg
}

// DepAmqpConnection opens an amqp session.
func (m DependencyMaker) DepAmqpConnection(logger zerolog.Logger, cfg config.CanopsisConf) libamqp.Connection {
	c, err := libamqp.NewConnection(logger, cfg.Global.ReconnectRetries, cfg.Global.GetReconnectTimeout())
	Panic("amqp session", err)
	return c
}

// DepAMQPChannelSub opens a channel from a given session, and apply Qos on it.
func (m DependencyMaker) DepAMQPChannelSub(session libamqp.Connection, prefetchCount, prefetchSize int) libamqp.Channel {
	channel, err := session.Channel()
	Panic("amqp consume channel", err)

	err = channel.Qos(prefetchCount, prefetchSize, true)
	Panic("amqp consume channel qos", err)

	return channel
}

// DepAMQPChannelPub opens a channel from a given session, to be used for publishing messages.
func (m DependencyMaker) DepAMQPChannelPub(session libamqp.Connection) libamqp.Channel {
	channel, err := session.Channel()
	Panic("amqp publish channel", err)
	return channel
}

// DepRedisSession opens a redis session.
func (m DependencyMaker) DepRedisSession(ctx context.Context, db int, logger zerolog.Logger, cfg config.CanopsisConf) *redismod.Client {
	s, err := redis.NewSession(ctx, db, logger, cfg.Global.ReconnectRetries,
		cfg.Global.GetReconnectTimeout())
	Panic("redis", err)
	return s
}

// Panic if err != nil, appending a message. Use this only in main programs.
func Panic(msg string, err error) {
	if err != nil {
		panic(fmt.Errorf("dependency error: %s: %w", msg, err))
	}
}

// Catch will recover from a panic, designed to be used with Panic().
func Catch(logger zerolog.Logger) {
	if r := recover(); r != nil {
		logger.Fatal().Msgf("fatal error: %v", r)
	}
}
