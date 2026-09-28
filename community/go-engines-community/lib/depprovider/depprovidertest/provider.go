// Package depprovidertest decorates a depprovider.Provider and its dependencies
// so a test can run against real infrastructure while intercepting
// selected operations (e.g. to observe or block a specific query).
package depprovidertest

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/depprovider"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

// NewProvider returns a Provider that delegates to p, applying the hooks in opts.
// Any dependency without a hook is created by p unchanged.
func NewProvider(p depprovider.Provider, opts Options) depprovider.Provider {
	return &provider{
		realProvider: p,
		opts:         opts,
	}
}

type Options struct {
	// WrapMongoClient, if set, returns a decorator around the real client so a test
	// can intercept specific DB operations while still hitting the real database.
	WrapMongoClient func(c mongo.DbClient) mongo.DbClient
}

type provider struct {
	realProvider depprovider.Provider
	opts         Options
}

func (p *provider) MongoClient(ctx context.Context, clientOptions mongo.ClientOptions) (mongo.DbClient, error) {
	d, err := p.realProvider.MongoClient(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	if p.opts.WrapMongoClient != nil {
		d = p.opts.WrapMongoClient(d)
	}

	return d, nil
}

func (p *provider) Config(ctx context.Context, dbClient mongo.DbClient) (config.CanopsisConf, error) {
	return p.realProvider.Config(ctx, dbClient)
}

func (p *provider) AMQPConnection(logger zerolog.Logger, cfg config.CanopsisConf) (amqp.Connection, error) {
	return p.realProvider.AMQPConnection(logger, cfg)
}

func (p *provider) AMQPConsumeChannelPool(conn amqp.Connection) amqp.ChannelPool {
	return p.realProvider.AMQPConsumeChannelPool(conn)
}

func (p *provider) AMQPPubChannelPool(conn amqp.Connection) amqp.ChannelPool {
	return p.realProvider.AMQPPubChannelPool(conn)
}

func (p *provider) RedisClient(ctx context.Context, db int, logger zerolog.Logger, cfg config.CanopsisConf) (*redis.Client, error) {
	return p.realProvider.RedisClient(ctx, db, logger, cfg)
}
