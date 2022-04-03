package mongo

//go:generate mockgen -destination=../../mocks/lib/mongo/mongo.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo DbCollection,DbClient,SingleResultHelper,Cursor

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	DB                                = "canopsis"
	ConfigurationMongoCollection      = "configuration"
	ObjectMongoCollection             = "object"
	RightsMongoCollection             = "default_rights"
	SessionMongoCollection            = "session"
	AlarmMongoCollection              = "periodical_alarm"
	EntityMongoCollection             = "default_entities"
	PbehaviorMongoCollection          = "pbehavior"
	PbehaviorTypeMongoCollection      = "pbehavior_type"
	PbehaviorReasonMongoCollection    = "pbehavior_reason"
	PbehaviorExceptionMongoCollection = "pbehavior_exception"
	ScenarioMongoCollection           = "action_scenario"
	FileMongoCollection               = "files"
	MetaAlarmRulesMongoCollection     = "meta_alarm_rules"
	IdleRuleMongoCollection           = "idle_rule"
	ExportTaskMongoCollection         = "export_task"
	ActionLogMongoCollection          = "action_log"
	EventFilterRulesMongoCollection   = "eventfilter"
	DynamicInfosRulesMongoCollection  = "dynamic_infos"
	EntityCategoryMongoCollection     = "entity_category"
	ImportJobMongoCollection          = "default_importgraph"
	JunitTestSuiteMongoCollection     = "junit_test_suite"
	JunitTestCaseMediaMongoCollection = "junit_test_case_media"
	PlaylistMongoCollection           = "view_playlist"
	StateSettingsMongoCollection      = "state_settings"
	BroadcastMessageMongoCollection   = "broadcast_message"
	AssociativeTableCollection        = "default_associativetable"
	NotificationMongoCollection       = "notification"

	ViewMongoCollection          = "views"
	ViewTabMongoCollection       = "viewtabs"
	WidgetMongoCollection        = "widgets"
	WidgetFiltersMongoCollection = "widget_filters"
	ViewGroupMongoCollection     = "viewgroups"

	// Following collections are used for event statistics.
	MessageRateStatsMinuteCollectionName = "message_rate_statistic_minute"
	MessageRateStatsHourCollectionName   = "message_rate_statistic_hour"

	// Collection for ok/ko event statistics
	EventStatistics = "event_statistics"

	// Remediation collections
	InstructionMongoCollection          = "instruction"
	InstructionExecutionMongoCollection = "instruction_execution"
	InstructionRatingMongoCollection    = "instruction_rating"
	JobConfigMongoCollection            = "job_config"
	JobMongoCollection                  = "job"
	JobHistoryMongoCollection           = "job_history"
	// Instruction statistics collections
	InstructionWeekStatsMongoCollection = "instruction_week_stats"
	InstructionModStatsMongoCollection  = "instruction_mod_stats"
	// Data storage alarm collections
	ResolvedAlarmMongoCollection = "resolved_alarms"
	ArchivedAlarmMongoCollection = "archived_alarms"
	// Data storage entity collections
	ArchivedEntitiesMongoCollection = "archived_entities"

	TokenMongoCollection = "token"

	ResolveRuleMongoCollection  = "resolve_rule"
	FlappingRuleMongoCollection = "flapping_rule"

	UserPreferencesMongoCollection = "userpreferences"

	FilterMongoCollection = "filter"

	PatternMongoCollection = "pattern"

	EntityServiceCountersMongoCollection = "entity_service_counters"
)

const (
	transactionReplicaSetVersion     = "4.0"
	transactionShardedClusterVersion = "4.2"
)

type SingleResultHelper interface {
	Decode(v interface{}) error
	DecodeBytes() (bson.Raw, error)
	Err() error
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
}

// DbClient connected MongoDB client settings
type DbClient interface {
	Collection(string) DbCollection
	Disconnect(ctx context.Context) error
	SetRetry(count int, timeout time.Duration)
	Ping(ctx context.Context, rp *readpref.ReadPref) error
	WithTransaction(ctx context.Context, f func(context.Context) error) error
	ListCollectionNames(ctx context.Context, filter interface{}, opts ...*options.ListCollectionsOptions) ([]string, error)
}

type dbClient struct {
	Client          *mongo.Client
	Database        *mongo.Database
	RetryCount      int
	MinRetryTimeout time.Duration

	TransactionEnabled bool
}

type dbCollection struct {
	mongoCollection *mongo.Collection
	retryCount      int
	minRetryTimeout time.Duration
}

func (c *dbCollection) Aggregate(ctx context.Context, pipeline interface{},
	opts ...*options.AggregateOptions) (Cursor, error) {
	var mongoCursor *mongo.Cursor
	var err error

	c.retry(ctx, func(ctx context.Context) error {
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
	c.retry(ctx, func(ctx context.Context) error {
		res, err = c.mongoCollection.BulkWrite(ctx, models, opts...)
		return err
	})

	return res, err
}

func (c *dbCollection) CountDocuments(ctx context.Context, filter interface{},
	opts ...*options.CountOptions) (int64, error) {
	var res int64
	var err error
	c.retry(ctx, func(ctx context.Context) error {
		res, err = c.mongoCollection.CountDocuments(ctx, filter, opts...)
		return err
	})

	return res, err
}

func (c *dbCollection) DeleteMany(ctx context.Context, filter interface{},
	opts ...*options.DeleteOptions) (int64, error) {
	var res *mongo.DeleteResult
	var err error
	c.retry(ctx, func(ctx context.Context) error {
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
	c.retry(ctx, func(ctx context.Context) error {
		res, err = c.mongoCollection.Distinct(ctx, fieldName, filter, opts...)
		return err
	})

	return res, err
}

func (c *dbCollection) Drop(ctx context.Context) error {
	var err error
	c.retry(ctx, func(ctx context.Context) error {
		err = c.mongoCollection.Drop(ctx)
		return err
	})

	return err
}

func (c *dbCollection) Find(ctx context.Context, filter interface{},
	opts ...*options.FindOptions) (Cursor, error) {
	var mongoCursor *mongo.Cursor
	var err error

	c.retry(ctx, func(ctx context.Context) error {
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
	c.retry(ctx, func(ctx context.Context) error {
		res = c.mongoCollection.FindOne(ctx, filter, opts...)
		return res.Err()
	})

	return res
}

func (c *dbCollection) FindOneAndDelete(ctx context.Context, filter interface{},
	opts ...*options.FindOneAndDeleteOptions) SingleResultHelper {
	var res *mongo.SingleResult
	c.retry(ctx, func(ctx context.Context) error {
		res = c.mongoCollection.FindOneAndDelete(ctx, filter, opts...)
		return res.Err()
	})

	return res
}

func (c *dbCollection) FindOneAndReplace(ctx context.Context, filter, replacement interface{},
	opts ...*options.FindOneAndReplaceOptions) SingleResultHelper {
	var res *mongo.SingleResult
	c.retry(ctx, func(ctx context.Context) error {
		res = c.mongoCollection.FindOneAndReplace(ctx, filter, replacement, opts...)
		return res.Err()
	})

	return res
}

func (c *dbCollection) FindOneAndUpdate(ctx context.Context, filter, update interface{},
	opts ...*options.FindOneAndUpdateOptions) SingleResultHelper {
	var res *mongo.SingleResult
	c.retry(ctx, func(ctx context.Context) error {
		res = c.mongoCollection.FindOneAndUpdate(ctx, filter, update, opts...)
		return res.Err()
	})

	return res
}

func (c *dbCollection) DeleteOne(ctx context.Context, filter interface{},
	opts ...*options.DeleteOptions) (int64, error) {
	var res *mongo.DeleteResult
	var err error
	c.retry(ctx, func(ctx context.Context) error {
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
	c.retry(ctx, func(ctx context.Context) error {
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
	c.retry(ctx, func(ctx context.Context) error {
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
	c.retry(ctx, func(ctx context.Context) error {
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
	c.retry(ctx, func(ctx context.Context) error {
		res, err = c.mongoCollection.UpdateMany(ctx, filter, update, opts...)
		return err
	})

	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *dbCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	var res *mongo.UpdateResult
	var err error
	c.retry(ctx, func(ctx context.Context) error {
		res, err = c.mongoCollection.UpdateOne(ctx, filter, update, opts...)
		return err
	})

	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *dbCollection) retry(ctx context.Context, f func(context.Context) error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	timeout := c.minRetryTimeout

	for i := 0; i <= c.retryCount; i++ {
		err := f(ctx)
		if err == nil {
			return
		}

		if c.retryCount == i || timeout == 0 {
			return
		}

		if !IsConnectionError(err) {
			return
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(timeout):
			timeout *= 2
		}
	}
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
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	db := client.Database(dbName)

	dbClient := &dbClient{
		Client:          client,
		Database:        db,
		RetryCount:      retryCount,
		MinRetryTimeout: minRetryTimeout,
	}

	dbClient.checkTransactionEnabled(ctx, logger)

	return dbClient, nil
}

func NewClientWithOptions(
	ctx context.Context,
	retryCount int,
	minRetryTimeout time.Duration,
	serverSelectionTimeout time.Duration,
	socketTimeout time.Duration,
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
	clientOptions.SetServerSelectionTimeout(serverSelectionTimeout)
	clientOptions.SetSocketTimeout(socketTimeout)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	db := client.Database(dbName)

	dbClient := &dbClient{
		Client:          client,
		Database:        db,
		RetryCount:      retryCount,
		MinRetryTimeout: minRetryTimeout,
	}

	dbClient.checkTransactionEnabled(ctx, logger)

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
	if !c.TransactionEnabled {
		return f(ctx)
	}

	session, err := c.Client.StartSession()
	if err != nil {
		return err
	}

	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		return nil, f(sessCtx)
	})

	return err
}

func (c *dbClient) checkTransactionEnabled(ctx context.Context, logger zerolog.Logger) {
	res, err := c.Database.RunCommand(ctx, bson.D{{"hello", 1}}). //nolint:govet
									DecodeBytes()
	if err != nil {
		logger.Err(err).Msg("cannot determine MongoDB version, transactions are disabled")
		return
	}

	helloResult := helloCommandResult{}
	err = bson.Unmarshal(res, &helloResult)
	if err != nil {
		logger.Err(err).Msg("cannot determine MongoDB version, transactions are disabled")
		return
	}

	if !helloResult.IsReplicaSet() && !helloResult.IsShardedCluster() {
		logger.Warn().Msg("MongoDB version does not support transactions, transactions are disabled")
		return
	}

	res, err = c.Client.Database("admin").RunCommand(ctx, bson.D{{"getParameter", 1}, {"featureCompatibilityVersion", 1}}). //nolint:govet
																DecodeBytes()
	if err != nil {
		logger.Err(err).Msg("cannot determine MongoDB version, transactions are disabled")
		return
	}

	fCVResult := struct {
		FeatureCompatibilityVersion struct {
			Version string `bson:"version"`
		} `bson:"featureCompatibilityVersion"`
	}{}

	err = bson.Unmarshal(res, &fCVResult)
	if err != nil {
		logger.Err(err).Msg("cannot determine MongoDB version, transactions are disabled")
		return
	}

	version := fCVResult.FeatureCompatibilityVersion.Version
	if helloResult.IsReplicaSet() && !isVersionGte(version, transactionReplicaSetVersion) ||
		helloResult.IsShardedCluster() && !isVersionGte(version, transactionShardedClusterVersion) {
		logger.Warn().Msg("MongoDB version does not support transactions, transactions are disabled")

		return
	}

	logger.Info().Msg("MongoDB version supports transactions, transactions are enabled")
	c.TransactionEnabled = true
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

func IsConnectionError(err error) bool {
	return mongo.IsNetworkError(err) ||
		strings.Contains(err.Error(), "server selection error")
}

type helloCommandResult struct {
	SetName string `bson:"setName"`
	Msg     string `bson:"msg"`
}

func (r helloCommandResult) IsReplicaSet() bool {
	return r.SetName != ""
}

func (r helloCommandResult) IsShardedCluster() bool {
	return r.Msg == "isdbgrid"
}

func isVersionGte(version, expectedVersion string) bool {
	if version == "" || expectedVersion == "" {
		return false
	}

	versionParts := strings.Split(version, ".")
	expectedVersionParts := strings.Split(expectedVersion, ".")

	for i, ev := range expectedVersionParts {
		if len(versionParts) <= i {
			return true
		}

		v := versionParts[i]
		vi, err := strconv.Atoi(v)
		if err != nil {
			return false
		}
		evi, err := strconv.Atoi(ev)
		if err != nil {
			return false
		}

		if vi < evi {
			return false
		}
	}

	return true
}
