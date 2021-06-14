package depmake

import (
	"fmt"

	"git.canopsis.net/canopsis/go-engines/lib/amqp"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/config"
	"git.canopsis.net/canopsis/go-engines/lib/influx"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	"git.canopsis.net/canopsis/go-engines/lib/redis"
	"github.com/globalsign/mgo"
	redismod "github.com/go-redis/redis/v7"
	influxmod "github.com/influxdata/influxdb/client/v2"
	"github.com/rs/zerolog"
	amqpmod "github.com/streadway/amqp"
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

// DepMongoSession opens a mongo session.
// Use DepMongoClient to migrate to mongo-driver.
func (m DependencyMaker) DepMongoSession() *mgo.Session {
	s, err := mongo.NewSession(mongo.Timeout)
	Panic("mongo session", err)
	return s
}

// DepMongoClient opens a mongo session.
func (m DependencyMaker) DepMongoClient(cfg config.CanopsisConf) mongo.DbClient {
	c, err := mongo.NewClient(
		cfg.Global.ReconnectRetries,
		cfg.Global.GetReconnectTimeout(),
	)
	Panic("mongo session", err)
	return c
}

// DepConfig gets a config from mongodb
func (m DependencyMaker) DepConfig() config.CanopsisConf {
	mongoSession := m.DepMongoSession()
	defer mongoSession.Close()

	configAdapter := config.NewLegacyAdapter(mongoSession)

	cfg, err := configAdapter.GetConfig()
	if err != nil {
		panic(fmt.Errorf("dependency error: %s: %v", "can't get the config", err))
	}

	return cfg
}

// DepAmqpSession opens an amqp session.
func (m DependencyMaker) DepAmqpSession() *amqpmod.Connection {
	s, err := amqp.NewSession()
	Panic("amqp session", err)
	return s
}

// DepAmqpConnection opens an amqp session.
func (m DependencyMaker) DepAmqpConnection(logger zerolog.Logger, cfg config.CanopsisConf) amqp.Connection {
	c, err := amqp.NewConnection(logger, cfg.Global.ReconnectRetries, cfg.Global.GetReconnectTimeout())
	Panic("amqp session", err)
	return c
}

// DepAMQPChannelSub opens a channel from a given session, and apply Qos on it.
func (m DependencyMaker) DepAMQPChannelSub(session amqp.Connection, prefetchCount, prefetchSize int) amqp.Channel {
	channel, err := session.Channel()
	Panic("amqp consume channel", err)

	err = channel.Qos(prefetchCount, prefetchSize, true)
	Panic("amqp consume channel qos", err)

	return channel
}

// DepAMQPChannelPub opens a channel from a given session, to be used for publishing messages.
func (m DependencyMaker) DepAMQPChannelPub(session amqp.Connection) amqp.Channel {
	channel, err := session.Channel()
	Panic("amqp publish channel", err)
	return channel
}

// DepRedisSession opens a redis session.
func (m DependencyMaker) DepRedisSession(db int, logger zerolog.Logger, cfg config.CanopsisConf) *redismod.Client {
	s, err := redis.NewSession(db, logger, cfg.Global.ReconnectRetries,
		cfg.Global.GetReconnectTimeout())
	Panic("redis", err)
	return s
}

// DepInfluxSession opens an influx session.
func (m DependencyMaker) DepInfluxSession() influxmod.Client {
	s, err := influx.NewSession()
	Panic("influx", err)
	return s
}

// Panic if err != nil, appending a message. Use this only in main programs.
func Panic(msg string, err error) {
	if err != nil {
		panic(fmt.Errorf("dependency error: %s: %v", msg, err))
	}
}

// Catch will recover from a panic, designed to be used with Panic().
func Catch(logger zerolog.Logger) {
	if r := recover(); r != nil {
		logger.Fatal().Msgf("fatal error: %v", r)
	}
}
