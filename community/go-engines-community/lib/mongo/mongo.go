package mongo

//go:generate go tool go.uber.org/mock/mockgen -destination=../../mocks/lib/mongo/mongo.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo DbCollection,DbClient,SingleResultHelper,Cursor,ChangeStream

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
	"go.mongodb.org/mongo-driver/v2/x/mongo/driver/description"
	"go.mongodb.org/mongo-driver/v2/x/mongo/driver/topology"
)

const (
	DefaultClientTimeout          = 15 * time.Second
	DefaultServerSelectionTimeout = 30 * time.Second

	disableRetries contextKey = "disable_retries"

	topologyCheckTimeout = 1 * time.Second

	ChangeStreamTypeInsert = "insert"
	ChangeStreamTypeUpdate = "update"
	ChangeStreamTypeDelete = "delete"

	DefaultGraphLookupMaxDepth = 100
)

type contextKey string

type SingleResultHelper interface {
	Decode(v interface{}) error
	Raw() (bson.Raw, error)
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
	Name() string
	Aggregate(ctx context.Context, pipeline interface{}, opts ...options.Lister[options.AggregateOptions]) (Cursor, error)
	BulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...options.Lister[options.BulkWriteOptions]) (*mongo.BulkWriteResult, error)
	CountDocuments(ctx context.Context, filter interface{}, opts ...options.Lister[options.CountOptions]) (int64, error)
	DeleteOne(ctx context.Context, filter interface{}, opts ...options.Lister[options.DeleteOneOptions]) (int64, error)
	DeleteMany(ctx context.Context, filter interface{}, opts ...options.Lister[options.DeleteManyOptions]) (int64, error)
	Distinct(ctx context.Context, fieldName string, filter interface{},
		opts ...options.Lister[options.DistinctOptions]) *mongo.DistinctResult
	Drop(ctx context.Context) error
	Find(ctx context.Context, filter interface{},
		opts ...options.Lister[options.FindOptions]) (Cursor, error)
	FindOne(ctx context.Context, filter interface{}, opts ...options.Lister[options.FindOneOptions]) SingleResultHelper
	FindOneAndDelete(ctx context.Context, filter interface{},
		opts ...options.Lister[options.FindOneAndDeleteOptions]) SingleResultHelper
	FindOneAndReplace(ctx context.Context, filter interface{}, replacement interface{},
		opts ...options.Lister[options.FindOneAndReplaceOptions]) SingleResultHelper
	FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{},
		opts ...options.Lister[options.FindOneAndUpdateOptions]) SingleResultHelper
	Indexes() mongo.IndexView
	InsertOne(ctx context.Context, document interface{},
		opts ...options.Lister[options.InsertOneOptions]) (interface{}, error)
	InsertMany(ctx context.Context, documents []interface{},
		opts ...options.Lister[options.InsertManyOptions]) ([]interface{}, error)
	ReplaceOne(ctx context.Context, filter interface{},
		replacement interface{}, opts ...options.Lister[options.ReplaceOptions]) (*mongo.UpdateResult, error)
	UpdateMany(ctx context.Context, filter interface{}, update interface{},
		opts ...options.Lister[options.UpdateManyOptions]) (*mongo.UpdateResult, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{},
		opts ...options.Lister[options.UpdateOneOptions]) (*mongo.UpdateResult, error)
	Watch(ctx context.Context, pipeline interface{},
		opts ...options.Lister[options.ChangeStreamOptions]) (ChangeStream, error)
}

// DbClient connected MongoDB client settings
type DbClient interface {
	Name() string
	Watch(ctx context.Context, pipeline interface{}, opts ...options.Lister[options.ChangeStreamOptions]) (*mongo.ChangeStream, error)
	Collection(string) DbCollection
	CreateCollection(ctx context.Context, name string, opts ...options.Lister[options.CreateCollectionOptions]) error
	Disconnect(ctx context.Context) error
	SetRetry(count int, timeout time.Duration)
	Ping(ctx context.Context, rp *readpref.ReadPref) error
	WithTransaction(ctx context.Context, f func(context.Context) error) error
	ListCollectionNames(ctx context.Context, filter interface{}, opts ...options.Lister[options.ListCollectionsOptions]) ([]string, error)
	RunCommand(ctx context.Context, runCommand any, opts ...options.Lister[options.RunCmdOptions]) SingleResultHelper
	RunAdminCommand(ctx context.Context, runCommand any, opts ...options.Lister[options.RunCmdOptions]) SingleResultHelper
}

type dbClient struct {
	Client          *mongo.Client
	Database        *mongo.Database
	RetryCount      int
	MinRetryTimeout time.Duration
}

// dbCollection
// nolint:wrapcheck
type dbCollection struct {
	mongoCollection *mongo.Collection
	retryCount      int
	minRetryTimeout time.Duration
}

type ClientOptions struct {
	RetryCount             int
	MinRetryTimeout        time.Duration
	ServerSelectionTimeout time.Duration
	ClientTimeout          time.Duration
	ReadPreference         *readpref.ReadPref

	// if NoClientTimeout set to true, then client timeout set to 0 and ClientTimeout value is ignored
	NoClientTimeout bool
}

func (c *dbCollection) Name() string {
	return c.mongoCollection.Name()
}

func (c *dbCollection) Aggregate(ctx context.Context, pipeline interface{},
	opts ...options.Lister[options.AggregateOptions]) (Cursor, error) {
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
	opts ...options.Lister[options.BulkWriteOptions]) (*mongo.BulkWriteResult, error) {
	var res *mongo.BulkWriteResult
	var err error
	retry(ctx, c.retryCount, c.minRetryTimeout, func(ctx context.Context) error {
		res, err = c.mongoCollection.BulkWrite(ctx, models, opts...)
		return err
	})

	return res, err
}

func (c *dbCollection) CountDocuments(ctx context.Context, filter interface{},
	opts ...options.Lister[options.CountOptions]) (int64, error) {
	var res int64
	var err error
	retry(ctx, c.retryCount, c.minRetryTimeout, func(ctx context.Context) error {
		res, err = c.mongoCollection.CountDocuments(ctx, filter, opts...)
		return err
	})

	return res, err
}

func (c *dbCollection) DeleteMany(ctx context.Context, filter interface{},
	opts ...options.Lister[options.DeleteManyOptions]) (int64, error) {
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
	opts ...options.Lister[options.DistinctOptions]) *mongo.DistinctResult {
	var res *mongo.DistinctResult

	retry(ctx, c.retryCount, c.minRetryTimeout, func(ctx context.Context) error {
		res = c.mongoCollection.Distinct(ctx, fieldName, filter, opts...)
		return res.Err()
	})

	return res
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
	opts ...options.Lister[options.FindOptions]) (Cursor, error) {
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
	opts ...options.Lister[options.FindOneOptions]) SingleResultHelper {
	var res *mongo.SingleResult
	retry(ctx, c.retryCount, c.minRetryTimeout, func(ctx context.Context) error {
		res = c.mongoCollection.FindOne(ctx, filter, opts...)
		return res.Err()
	})

	return res
}

func (c *dbCollection) FindOneAndDelete(ctx context.Context, filter interface{},
	opts ...options.Lister[options.FindOneAndDeleteOptions]) SingleResultHelper {
	var res *mongo.SingleResult
	retry(ctx, c.retryCount, c.minRetryTimeout, func(ctx context.Context) error {
		res = c.mongoCollection.FindOneAndDelete(ctx, filter, opts...)
		return res.Err()
	})

	return res
}

func (c *dbCollection) FindOneAndReplace(ctx context.Context, filter, replacement interface{},
	opts ...options.Lister[options.FindOneAndReplaceOptions]) SingleResultHelper {
	var res *mongo.SingleResult
	retry(ctx, c.retryCount, c.minRetryTimeout, func(ctx context.Context) error {
		res = c.mongoCollection.FindOneAndReplace(ctx, filter, replacement, opts...)
		return res.Err()
	})

	return res
}

func (c *dbCollection) FindOneAndUpdate(ctx context.Context, filter, update interface{},
	opts ...options.Lister[options.FindOneAndUpdateOptions]) SingleResultHelper {
	var res *mongo.SingleResult
	retry(ctx, c.retryCount, c.minRetryTimeout, func(ctx context.Context) error {
		res = c.mongoCollection.FindOneAndUpdate(ctx, filter, update, opts...)
		return res.Err()
	})

	return res
}

func (c *dbCollection) DeleteOne(ctx context.Context, filter interface{},
	opts ...options.Lister[options.DeleteOneOptions]) (int64, error) {
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
	opts ...options.Lister[options.InsertOneOptions]) (interface{}, error) {
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
	opts ...options.Lister[options.InsertManyOptions]) ([]interface{}, error) {
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
	opts ...options.Lister[options.ReplaceOptions]) (*mongo.UpdateResult, error) {
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
	opts ...options.Lister[options.UpdateManyOptions]) (*mongo.UpdateResult, error) {
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
	opts ...options.Lister[options.ChangeStreamOptions]) (ChangeStream, error) {
	return c.mongoCollection.Watch(ctx, pipeline, opts...)
}

func (c *dbCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{},
	opts ...options.Lister[options.UpdateOneOptions]) (*mongo.UpdateResult, error) {
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

func NewClient(ctx context.Context, opts ...ClientOptions) (DbClient, error) {
	var clientOptions ClientOptions
	if len(opts) == 1 {
		clientOptions = opts[0]
	} else if len(opts) > 1 {
		return nil, errors.New("only one ClientOptions is allowed")
	}

	mongoURL, dbName, err := getURL()
	if err != nil {
		return nil, err
	}

	if dbName == "*" {
		dbName = DB
	}

	bsonOpts := &options.BSONOptions{
		DefaultDocumentM:    true,
		ObjectIDAsHexString: true,
	}

	mongoClientOptions := options.Client().ApplyURI(mongoURL).SetBSONOptions(bsonOpts)

	if clientOptions.ServerSelectionTimeout <= 0 {
		if mongoClientOptions.ServerSelectionTimeout == nil {
			mongoClientOptions.SetServerSelectionTimeout(DefaultServerSelectionTimeout)
		}
	} else {
		mongoClientOptions.SetServerSelectionTimeout(clientOptions.ServerSelectionTimeout)
	}

	if clientOptions.NoClientTimeout {
		mongoClientOptions.SetTimeout(0)
	} else if clientOptions.ClientTimeout <= 0 {
		if mongoClientOptions.Timeout == nil {
			mongoClientOptions.SetTimeout(DefaultClientTimeout)
		}
	} else {
		mongoClientOptions.SetTimeout(clientOptions.ClientTimeout)
	}

	if clientOptions.ReadPreference == nil {
		// don't get readPreference from mongoClientOptions, it should be defined ONLY by clientOptions to avoid the readPreference misusing by the end user.
		mongoClientOptions.SetReadPreference(readpref.Primary())
	} else {
		mongoClientOptions.SetReadPreference(clientOptions.ReadPreference)
	}

	isDistributed, err := isMongoReplicaSetEnabled(ctx, mongoClientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to check if replica set is enabled: %w", err)
	}

	if !isDistributed {
		return nil, errors.New("replica set is required")
	}

	mongoClient, err := mongo.Connect(mongoClientOptions)
	if err != nil {
		return nil, err
	}

	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		_ = mongoClient.Disconnect(ctx)
		return nil, err
	}

	return &dbClient{
		Client:          mongoClient,
		Database:        mongoClient.Database(dbName),
		RetryCount:      clientOptions.RetryCount,
		MinRetryTimeout: clientOptions.MinRetryTimeout,
	}, nil
}

func (c *dbClient) Name() string {
	return c.Database.Name()
}

func (c *dbClient) Watch(ctx context.Context, pipeline interface{},
	opts ...options.Lister[options.ChangeStreamOptions]) (*mongo.ChangeStream, error) {
	return c.Database.Watch(ctx, pipeline, opts...)
}

func (c *dbClient) Collection(name string) DbCollection {
	return &dbCollection{
		mongoCollection: c.Database.Collection(name),
		retryCount:      c.RetryCount,
		minRetryTimeout: c.MinRetryTimeout,
	}
}

func (c *dbClient) CreateCollection(ctx context.Context, name string, opts ...options.Lister[options.CreateCollectionOptions]) error {
	return c.Database.CreateCollection(ctx, name, opts...)
}

func (c *dbClient) Disconnect(ctx context.Context) error {
	return c.Client.Disconnect(ctx)
}

func (c *dbClient) Ping(ctx context.Context, rp *readpref.ReadPref) error {
	return c.Client.Ping(ctx, rp)
}

func (c *dbClient) ListCollectionNames(ctx context.Context, filter interface{}, opts ...options.Lister[options.ListCollectionsOptions]) ([]string, error) {
	return c.Database.ListCollectionNames(ctx, filter, opts...)
}

func (c *dbClient) SetRetry(count int, timeout time.Duration) {
	c.RetryCount = count
	c.MinRetryTimeout = timeout
}

func (c *dbClient) WithTransaction(ctx context.Context, f func(context.Context) error) error {
	txnOpts := options.Transaction().SetReadPreference(readpref.Primary())
	sessOpts := options.Session().SetDefaultTransactionOptions(txnOpts)

	var session *mongo.Session
	var err error

	retry(ctx, c.RetryCount, c.MinRetryTimeout, func(ctx context.Context) error {
		session, err = c.Client.StartSession(sessOpts)
		if err != nil {
			return err
		}

		defer session.EndSession(ctx)

		_, err = session.WithTransaction(ctx, func(sessCtx context.Context) (interface{}, error) {
			return nil, f(context.WithValue(sessCtx, disableRetries, true))
		})

		return err
	})

	return err
}

func (c *dbClient) RunCommand(ctx context.Context, runCommand any, opts ...options.Lister[options.RunCmdOptions]) SingleResultHelper {
	return c.Database.RunCommand(ctx, runCommand, opts...)
}

func (c *dbClient) RunAdminCommand(ctx context.Context, runCommand any, opts ...options.Lister[options.RunCmdOptions]) SingleResultHelper {
	return c.Database.Client().Database("admin").RunCommand(ctx, runCommand, opts...)
}

func isMongoReplicaSetEnabled(ctx context.Context, clientOptions *options.ClientOptions) (bool, error) {
	cfg, err := topology.NewConfig(clientOptions, nil)
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
			case description.TopologyKindSharded,
				description.TopologyKindReplicaSet,
				description.TopologyKindReplicaSetNoPrimary,
				description.TopologyKindReplicaSetWithPrimary:
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
	var sse topology.ServerSelectionError

	return mongo.IsNetworkError(err) || errors.As(err, &sse)
}

func SecondaryPreferred(opts ...readpref.Option) *readpref.ReadPref {
	return readpref.SecondaryPreferred(opts...)
}
