package depprovidertest

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewMongoClient(c mongo.DbClient, opts MongoClientOptions) mongo.DbClient {
	return &hookedDbClient{
		DbClient: c,
		opts:     opts,
	}
}

func NewMongoCollection(coll mongo.DbCollection, opts MongoCollOptions) mongo.DbCollection {
	return &hookedDbCollection{
		DbCollection: coll,
		opts:         opts,
	}
}

func NewMongoCursor(cur mongo.Cursor, opts MongoCursorOptions) mongo.Cursor {
	return &hookedCursor{
		Cursor: cur,
		opts:   opts,
	}
}

type MongoClientOptions struct {
	WrapCollection func(coll mongo.DbCollection) mongo.DbCollection
}

type MongoCollOptions struct {
	WrapAggregate func(c mongo.Cursor) mongo.Cursor
}

type MongoCursorOptions struct {
	OnAll    func(results any)
	OnNext   func(bool)
	OnDecode func(result any)
}

type hookedDbClient struct {
	mongo.DbClient
	opts MongoClientOptions
}

func (c *hookedDbClient) Collection(name string) mongo.DbCollection {
	coll := c.DbClient.Collection(name)
	if c.opts.WrapCollection != nil {
		coll = c.opts.WrapCollection(coll)
	}

	return coll
}

type hookedDbCollection struct {
	mongo.DbCollection
	opts MongoCollOptions
}

func (c *hookedDbCollection) Aggregate(ctx context.Context, pipeline any, opts ...options.Lister[options.AggregateOptions]) (mongo.Cursor, error) {
	cur, err := c.DbCollection.Aggregate(ctx, pipeline, opts...)
	if err != nil {
		return cur, err
	}

	if c.opts.WrapAggregate != nil {
		cur = c.opts.WrapAggregate(cur)
	}

	return cur, nil
}

type hookedCursor struct {
	mongo.Cursor
	opts MongoCursorOptions
}

func (c *hookedCursor) All(ctx context.Context, results any) error {
	err := c.Cursor.All(ctx, results)
	if err != nil {
		return err
	}

	if c.opts.OnAll != nil {
		c.opts.OnAll(results)
	}

	return nil
}

func (c *hookedCursor) Next(ctx context.Context) bool {
	b := c.Cursor.Next(ctx)
	if c.opts.OnNext != nil {
		c.opts.OnNext(b)
	}

	return b
}

func (c *hookedCursor) Decode(result any) error {
	err := c.Cursor.Decode(result)
	if err != nil {
		return err
	}

	if c.opts.OnDecode != nil {
		c.opts.OnDecode(result)
	}

	return nil
}
