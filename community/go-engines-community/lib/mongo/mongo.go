package mongo

//go:generate mockgen -destination=../../mocks/lib/mongo/mongo.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo DbCollection,DbClient,SingleResultHelper,Cursor,ChangeStream

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/mongo/driver/topology"
)

const (
	defaultClientTimeout            = 15 * time.Second
	disableRetries       contextKey = "disable_retries"

	topologyCheckTimeout = 1 * time.Second
)

type contextKey string

type SingleResultHelper interface {
	Decode(v interface{}) error
	DecodeBytes() (bson.Raw, error)
	Err() error
}

type ChangeStream interface {
	ID() int64
	Decode(val interface{}) error
	Err() error
	Close(ctx context.Context) error
	ResumeToken() bson.Raw
	Next(ctx context.Context) bool
	TryNext(ctx context.Context) bool
}

type DbCollection interface {
	Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (Cursor, error)
	BulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error)
	CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error)
	DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (int64, error)
	DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (int64, error)
	Distinct(ctx context.Context, fieldName string, filter interface{},
		opts ...*options.DistinctOptions) ([]interface{}, error)
	Drop(ctx context.Context) error
	Find(ctx context.Context, filter interface{},
		opts ...*options.FindOptions) (Cursor, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) SingleResultHelper
	FindOneAndDelete(ctx context.Context, filter interface{},
		opts ...*options.FindOneAndDeleteOptions) SingleResultHelper
	FindOneAndReplace(ctx context.Context, filter interface{}, replacement interface{},
		opts ...*options.FindOneAndReplaceOptions) SingleResultHelper
	FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{},
		opts ...*options.FindOneAndUpdateOptions) SingleResultHelper
	Indexes() mongo.IndexView
	InsertOne(ctx context.Context, document interface{},
		opts ...*options.InsertOneOptions) (interface{}, error)
	InsertMany(ctx context.Context, documents []interface{},
		opts ...*options.InsertManyOptions) ([]interface{}, error)
	ReplaceOne(ctx context.Context, filter interface{},
		replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error)
	UpdateMany(ctx context.Context, filter interface{}, update interface{},
		opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{},
		opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	Watch(ctx context.Context, pipeline interface{},
		opts ...*options.ChangeStreamOptions) (ChangeStream, error)
}

// DbClient connected MongoDB client settings
type DbClient interface {
	Collection(string) DbCollection
	Disconnect(ctx context.Context) error
	SetRetry(count int, timeout time.Duration)
	Ping(ctx context.Context, rp *readpref.ReadPref) error
	WithTransaction(ctx context.Context, f func(context.Context) error) error
	ListCollectionNames(ctx context.Context, filter interface{}, opts ...*options.ListCollectionsOptions) ([]string, error)
	IsDistributed() bool
}

type dbClient struct {
	Client          *mongo.Client
	Database        *mongo.Database
	RetryCount      int
	MinRetryTimeout time.Duration

	isDistributed bool
}

// dbCollection
// nolint:wrapcheck
type dbCollection struct {
	mongoCollection *mongo.Collection
	retryCount      int
	minRetryTimeout time.Duration
}

func (c *dbCollection) Aggregate(ctx context.Context, pipeline interface{},
	opts ...*options.AggregateOptions) (Cursor, error) {
	var mongoCursor *mongo.Cursor
	var err error

	retry(ctx, c.retryCount, c.minRetryTimeout, func(ctx context.Context) error {
		mongoCursor, err = c.mongoCollection.Aggregate(ctx, pipeline, opts...)
		return err
	})

	if err != nil {
		return nil, err
	}

	return &cursor{mongoCursor: mongoCursor}, nil
}

func (c *dbCollection) BulkWrite(ctx context.Context, models []mongo.WriteModel,
	opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error) {
	var res *mongo.BulkWriteResult
	var err error
	retry(ctx, c.retryCount, c.minRetryTimeout, func(ctx context.Context) error {
		res, err = c.mongoCollection.BulkWrite(ctx, models, opts...)
		return err
	})

	return res, err
}

func (c *dbCollection) CountDocuments(ctx context.Context, filter interface{},
	opts ...*options.CountOptions) (int64, error) {
	var res int64
	var err error
	retry(ctx, c.retryCount, c.minRetryTimeout, func(ctx context.Context) error {
		res, err = c.mongoCollection.CountDocuments(ctx, filter, opts...)
		return err
	})

	return res, err
}

func (c *dbCollection) DeleteMany(ctx context.Context, filter interface{},
	opts ...*options.DeleteOptions) (int64, error) {
	var res *mongo.DeleteResult
	var err error
	retry(ctx, c.retryCount, c.minRetryTimeout, func(ctx context.Context) error {
		res, err = c.mongoCollection.DeleteMany(ctx, filter, opts...)
		return err
	})

	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}

func (c *dbCollection) Distinct(ctx context.Context, fieldName string, filter interface{},
	opts ...*options.DistinctOptions) ([]interface{}, error) {
	var res []interface{}
	var err error
	retry(ctx, c.retryCount, c.minRetryTimeout, func(ctx context.Context) error {
		res, err = c.mongoCollection.Distinct(ctx, fieldName, filter, opts...)
		return err
	})

	return res, err
}

func (c *dbCollection) Drop(ctx context.Context) error {
	var err error
	retry(ctx, c.retryCount, c.minRetryTimeout, func(ctx context.Context) error {
		err = c.mongoCollection.Drop(ctx)
		return err
	})

	return err
}

func (c *dbCollection) Find(ctx context.Context, filter interface{},
	opts ...*options.FindOptions) (Cursor, error) {
	var mongoCursor *mongo.Cursor
	var err error

	retry(ctx, c.retryCount, c.minRetryTimeout, func(ctx context.Context) error {
		mongoCursor, err = c.mongoCollection.Find(ctx, filter, opts...)
		return err
	})

	if err != nil {
		return nil, err
	}

	return &cursor{mongoCursor: mongoCursor}, nil
}

func (c *dbCollection) FindOne(ctx context.Context, filter interface{},
	opts ...*options.FindOneOptions) SingleResultHelper {
	var res *mongo.SingleResult
	retry(ctx, c.retryCount, c.minRetryTimeout, func(ctx context.Context) error {
		res = c.mongoCollection.FindOne(ctx, filter, opts...)
		return res.Err()
	})

	return res
}

func (c *dbCollection) FindOneAndDelete(ctx context.Context, filter interface{},
	opts ...*options.FindOneAndDeleteOptions) SingleResultHelper {
	var res *mongo.SingleResult
	retry(ctx, c.retryCount, c.minRetryTimeout, func(ctx context.Context) error {
		res = c.mongoCollection.FindOneAndDelete(ctx, filter, opts...)
		return res.Err()
	})

	return res
}

func (c *dbCollection) FindOneAndReplace(ctx context.Context, filter, replacement interface{},
	opts ...*options.FindOneAndReplaceOptions) SingleResultHelper {
	var res *mongo.SingleResult
	retry(ctx, c.retryCount, c.minRetryTimeout, func(ctx context.Context) error {
		res = c.mongoCollection.FindOneAndReplace(ctx, filter, replacement, opts...)
		return res.Err()
	})

	return res
}

func (c *dbCollection) FindOneAndUpdate(ctx context.Context, filter, update interface{},
	opts ...*options.FindOneAndUpdateOptions) SingleResultHelper {
	var res *mongo.SingleResult
	retry(ctx, c.retryCount, c.minRetryTimeout, func(ctx context.Context) error {
		res = c.mongoCollection.FindOneAndUpdate(ctx, filter, update, opts...)
		return res.Err()
	})

	return res
}

func (c *dbCollection) DeleteOne(ctx context.Context, filter interface{},
	opts ...*options.DeleteOptions) (int64, error) {
	var res *mongo.DeleteResult
	var err error
	retry(ctx, c.retryCount, c.minRetryTimeout, func(ctx context.Context) error {
		res, err = c.mongoCollection.DeleteOne(ctx, filter, opts...)
		return err
	})

	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}

func (c *dbCollection) Indexes() mongo.IndexView {
	return c.mongoCollection.Indexes()
}

func (c *dbCollection) InsertOne(ctx context.Context, document interface{},
	opts ...*options.InsertOneOptions) (interface{}, error) {
	var res *mongo.InsertOneResult
	var err error
	retry(ctx, c.retryCount, c.minRetryTimeout, func(ctx context.Context) error {
		res, err = c.mongoCollection.InsertOne(ctx, document, opts...)
		return err
	})

	if err != nil {
		return nil, err
	}

	return res.InsertedID, nil
}

func (c *dbCollection) InsertMany(ctx context.Context, documents []interface{},
	opts ...*options.InsertManyOptions) ([]interface{}, error) {
	var res *mongo.InsertManyResult
	var err error
	retry(ctx, c.retryCount, c.minRetryTimeout, func(ctx context.Context) error {
		res, err = c.mongoCollection.InsertMany(ctx, documents, opts...)
		return err
	})

	if err != nil {
		return nil, err
	}
	return res.InsertedIDs, nil
}

func (c *dbCollection) ReplaceOne(ctx context.Context, filter, replacement interface{},
	opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	var res *mongo.UpdateResult
	var err error
	retry(ctx, c.retryCount, c.minRetryTimeout, func(ctx context.Context) error {
		res, err = c.mongoCollection.ReplaceOne(ctx, filter, replacement, opts...)
		return err
	})

	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *dbCollection) UpdateMany(ctx context.Context, filter interface{}, update interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	var res *mongo.UpdateResult
	var err error
	retry(ctx, c.retryCount, c.minRetryTimeout, func(ctx context.Context) error {
		res, err = c.mongoCollection.UpdateMany(ctx, filter, update, opts...)
		return err
	})

	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *dbCollection) Watch(ctx context.Context, pipeline interface{},
	opts ...*options.ChangeStreamOptions) (ChangeStream, error) {
	return c.mongoCollection.Watch(ctx, pipeline, opts...)
}

func (c *dbCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	var res *mongo.UpdateResult
	var err error
	retry(ctx, c.retryCount, c.minRetryTimeout, func(ctx context.Context) error {
		res, err = c.mongoCollection.UpdateOne(ctx, filter, update, opts...)
		return err
	})

	if err != nil {
		return nil, err
	}
	return res, nil
}

// NewClient creates a new connection to the MongoDB database.
// It uses EnvURL as configuration source.
func NewClient(ctx context.Context, retryCount int, minRetryTimeout time.Duration, logger zerolog.Logger) (DbClient, error) {
	mongoURL, dbName, err := getURL()
	if err != nil {
		return nil, err
	}
	if dbName == "*" {
		dbName = DB
	}

	clientOptions := options.Client().ApplyURI(mongoURL)
	if clientOptions.Timeout == nil {
		clientOptions.SetTimeout(defaultClientTimeout)
	}
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	dbClient := &dbClient{
		Client:          client,
		Database:        client.Database(dbName),
		RetryCount:      retryCount,
		MinRetryTimeout: minRetryTimeout,
	}

	err = dbClient.checkTransactionEnabled(ctx, logger)
	if err != nil {
		return nil, err
	}

	return dbClient, nil
}

func NewClientWithOptions(
	ctx context.Context,
	retryCount int,
	minRetryTimeout time.Duration,
	serverSelectionTimeout time.Duration,
	clientTimeout time.Duration,
	logger zerolog.Logger,
) (DbClient, error) {
	mongoURL, dbName, err := getURL()
	if err != nil {
		return nil, err
	}
	if dbName == "*" {
		dbName = DB
	}

	clientOptions := options.Client().ApplyURI(mongoURL)
	if serverSelectionTimeout > 0 {
		clientOptions.SetServerSelectionTimeout(serverSelectionTimeout)
	}
	if clientTimeout > 0 {
		clientOptions.SetTimeout(clientTimeout)
	}
	if clientOptions.Timeout == nil {
		clientOptions.SetTimeout(defaultClientTimeout)
	}
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	dbClient := &dbClient{
		Client:          client,
		Database:        client.Database(dbName),
		RetryCount:      retryCount,
		MinRetryTimeout: minRetryTimeout,
	}

	err = dbClient.checkTransactionEnabled(ctx, logger)
	if err != nil {
		return nil, err
	}

	return dbClient, nil
}

func (c *dbClient) Collection(name string) DbCollection {
	return &dbCollection{
		mongoCollection: c.Database.Collection(name),
		retryCount:      c.RetryCount,
		minRetryTimeout: c.MinRetryTimeout,
	}
}

func (c *dbClient) Disconnect(ctx context.Context) error {
	return c.Client.Disconnect(ctx)
}

func (c *dbClient) Ping(ctx context.Context, rp *readpref.ReadPref) error {
	return c.Client.Ping(ctx, rp)
}

func (c *dbClient) ListCollectionNames(ctx context.Context, filter interface{}, opts ...*options.ListCollectionsOptions) ([]string, error) {
	return c.Database.ListCollectionNames(ctx, filter, opts...)
}

func (c *dbClient) SetRetry(count int, timeout time.Duration) {
	c.RetryCount = count
	c.MinRetryTimeout = timeout
}

func (c *dbClient) WithTransaction(ctx context.Context, f func(context.Context) error) error {
	if !c.isDistributed {
		return f(ctx)
	}

	opts := options.Session().SetDefaultReadPreference(readpref.Primary())

	var session mongo.Session
	var err error

	retry(ctx, c.RetryCount, c.MinRetryTimeout, func(ctx context.Context) error {
		session, err = c.Client.StartSession(opts)
		if err != nil {
			return err
		}

		defer session.EndSession(ctx)

		_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) { // nolint:contextcheck
			return nil, f(context.WithValue(sessCtx, disableRetries, true))
		})

		return err
	})

	return err
}

func (c *dbClient) checkTransactionEnabled(ctx context.Context, logger zerolog.Logger) error {
	var err error

	c.isDistributed, err = isMongoReplicaSetEnabled(ctx)
	if err != nil {
		return err
	}

	if c.isDistributed {
		logger.Info().Msg("replica set is detected, transactions are enabled")
	} else {
		logger.Warn().Msg("replica set is not detected, transactions are disabled")
	}

	return nil
}

// IsDistributed returns true if MongoDB is Replica Set or Sharded Cluster.
// Use to check feature availability : Transactions, Change Streams, etc.
func (c *dbClient) IsDistributed() bool {
	return c.isDistributed
}

func isMongoReplicaSetEnabled(ctx context.Context) (bool, error) {
	mongoURL, _, err := getURL()
	if err != nil {
		return false, fmt.Errorf("could not get mongo url: %w", err)
	}

	cfg, err := topology.NewConfig(options.Client().ApplyURI(mongoURL), nil)
	if err != nil {
		return false, fmt.Errorf("could not create topology config: %w", err)
	}

	top, err := topology.New(cfg)
	if err != nil {
		return false, fmt.Errorf("could not create topology: %w", err)
	}

	defer func() {
		_ = top.Disconnect(ctx)
	}()

	err = top.Connect()
	if err != nil {
		return false, fmt.Errorf("could not connect to topology: %w", err)
	}

	sub, err := top.Subscribe()
	if err != nil {
		return false, fmt.Errorf("could not subscribe to topology: %w", err)
	}

	defer func() {
		_ = top.Unsubscribe(sub)
	}()

	for {
		select {
		case <-ctx.Done():
			return false, nil
		case <-time.After(topologyCheckTimeout):
			return false, nil
		case desc, ok := <-sub.Updates:
			if !ok {
				return false, fmt.Errorf("topology subscription was closed: %w", err)
			}

			switch desc.Kind {
			case description.Unknown:
				continue
			case description.Sharded,
				description.ReplicaSet,
				description.ReplicaSetNoPrimary,
				description.ReplicaSetWithPrimary:
				return true, nil
			default:
				return false, nil
			}
		}
	}
}

// getURL parses URL value in EnvURL environment variable
func getURL() (mongoURL, dbName string, err error) {
	mongoURL = os.Getenv(EnvURL)
	if mongoURL == "" {
		return "", "", fmt.Errorf("environment variable %s empty", EnvURL)
	}
	parsed, err := url.ParseRequestURI(mongoURL)
	if err != nil {
		return "", "", err
	}
	dbName = strings.TrimPrefix(parsed.EscapedPath(), "/")
	return mongoURL, dbName, nil
}

func retry(ctx context.Context, retryCount int, retryTimeout time.Duration, f func(context.Context) error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	withoutRetries, _ := ctx.Value(disableRetries).(bool)
	if withoutRetries {
		_ = f(ctx)
		return
	}

	for i := 0; i <= retryCount; i++ {
		err := f(ctx)
		if err == nil {
			return
		}

		if retryCount == i || retryTimeout == 0 {
			return
		}

		if !IsConnectionError(err) && !mongo.IsDuplicateKeyError(err) {
			return
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(retryTimeout):
			retryTimeout *= 2
		}
	}
}

func IsConnectionError(err error) bool {
	return mongo.IsNetworkError(err) ||
		strings.Contains(err.Error(), "server selection error")
}
