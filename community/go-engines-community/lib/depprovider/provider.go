// Package depprovider provides Canopsis service dependencies.
package depprovider

import (
	"context"
	"fmt"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	redismod "github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

// Provider builds the dependencies shared by Canopsis services.
type Provider interface {
	MongoClient(ctx context.Context, clientOptions mongo.ClientOptions) (mongo.DbClient, error)
	Config(ctx context.Context, dbClient mongo.DbClient) (config.CanopsisConf, error)
	AMQPConnection(logger zerolog.Logger, cfg config.CanopsisConf) (libamqp.Connection, error)
	AMQPConsumeChannelPool(conn libamqp.Connection) libamqp.ChannelPool
	AMQPPubChannelPool(conn libamqp.Connection) libamqp.ChannelPool
	RedisClient(ctx context.Context, db int, logger zerolog.Logger, cfg config.CanopsisConf) (*redismod.Client, error)
}

func NewProvider() Provider {
	return provider{}
}

type provider struct{}

func (provider) MongoClient(ctx context.Context, clientOptions mongo.ClientOptions) (mongo.DbClient, error) {
	c, err := mongo.NewClient(ctx, clientOptions)
	if err != nil {
		return nil, wrap("mongodb connection", err)
	}

	return c, nil
}

func (provider) Config(ctx context.Context, dbClient mongo.DbClient) (config.CanopsisConf, error) {
	cfg, err := config.NewAdapter(dbClient).GetConfig(ctx)
	if err != nil {
		return config.CanopsisConf{}, wrap("config", err)
	}

	return cfg, nil
}

func (provider) AMQPConnection(logger zerolog.Logger, cfg config.CanopsisConf) (libamqp.Connection, error) {
	c, err := libamqp.New(cfg.Global.ReconnectRetries, cfg.Global.GetReconnectTimeout(), logger)
	if err != nil {
		return nil, wrap("amqp connection", err)
	}

	return c, nil
}

func (provider) AMQPConsumeChannelPool(conn libamqp.Connection) libamqp.ChannelPool {
	return libamqp.NewChannelPool(conn, 0)
}

func (provider) AMQPPubChannelPool(conn libamqp.Connection) libamqp.ChannelPool {
	return libamqp.NewChannelPool(conn, canopsis.DefaultAMQPPublishPoolSize)
}

func (provider) RedisClient(ctx context.Context, db int, logger zerolog.Logger, cfg config.CanopsisConf) (*redismod.Client, error) {
	s, err := redis.NewSession(ctx, db, logger, cfg.Global.ReconnectRetries,
		cfg.Global.GetReconnectTimeout())
	if err != nil {
		return nil, wrap("redis", err)
	}

	return s, nil
}

func wrap(msg string, err error) error {
	return fmt.Errorf("cannot create dependency: %s: %w", msg, err)
}
