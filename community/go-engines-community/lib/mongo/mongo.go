package mongo

//go:generate mockgen -destination=../../mocks/lib/mongo/mongo.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo DbCollection,DbClient,SingleResultHelper,Cursor

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
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	DB                                = "canopsis"
	ConfigurationMongoCollection      = "configuration"
	RightsMongoCollection             = "default_rights"
	SessionMongoCollection            = "session"
	AlarmMongoCollection              = "periodical_alarm"
	EntityMongoCollection             = "default_entities"
	PbehaviorMongoCollection          = "pbehavior"
	PbehaviorTypeMongoCollection      = "pbehavior_type"
	PbehaviorReasonMongoCollection    = "pbehavior_reason"
	PbehaviorExceptionMongoCollection = "pbehavior_exception"
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

	ViewMongoCollection           = "views"
	ViewTabMongoCollection        = "viewtabs"
	WidgetMongoCollection         = "widgets"
	WidgetFiltersMongoCollection  = "widget_filters"
	WidgetTemplateMongoCollection = "widget_templates"
	ViewGroupMongoCollection      = "viewgroups"

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
	// InstructionWeekStatsMongoCollection
	// Deprecated : keep for backward compatibility, remove in the next release
	InstructionWeekStatsMongoCollection = "instruction_week_stats"
	// Data storage alarm collections
	ResolvedAlarmMongoCollection = "resolved_alarms"
	ArchivedAlarmMongoCollection = "archived_alarms"
	// Data storage entity collections
	ArchivedEntitiesMongoCollection = "archived_entities"

	TokenMongoCollection      = "token"
	ShareTokenMongoCollection = "share_token"

	ResolveRuleMongoCollection  = "resolve_rule"
	FlappingRuleMongoCollection = "flapping_rule"

	UserPreferencesMongoCollection = "userpreferences"

	KpiFilterMongoCollection = "kpi_filter"

	PatternMongoCollection = "pattern"

	EntityInfosDictionaryCollection  = "entity_infos_dictionary"
	DynamicInfosDictionaryCollection = "dynamic_infos_dictionary"

	MapMongoCollection = "map"

	AlarmTagCollection = "alarm_tag"

	ScenarioMongoCollection          = "action_scenario"
	DeclareTicketRuleMongoCollection = "declare_ticket_rule"
	WebhookHistoryMongoCollection    = "webhook_history"

	LinkRuleMongoCollection = "link_rule"

	OcwsNocChoiceCollection    = "ocws_noc_snow_sys_choice_new"
	OcwsNocContractCollection  = "ocws_noc_snow_service_contract_new"
	OcwsNocLocationCollection  = "ocws_noc_snow_location_group_new"
	OcwsNocGroupCollection     = "ocws_noc_snow_group_new"
	OcwsNocUserCollection      = "ocws_noc_snow_user_new"
	OcwsNocUserGroupCollection = "ocws_noc_snow_user_group_new"
)

const (
	defaultClientTimeout = 15 * time.Second
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
	Watch(ctx context.Context, pipeline interface{},
		opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error)
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

func (c *dbCollection) Watch(ctx context.Context, pipeline interface{},
	opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error) {
	return c.mongoCollection.Watch(ctx, pipeline, opts...)
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
	if !c.isDistributed {
		return f(ctx)
	}

	opts := options.Session().SetDefaultReadPreference(readpref.Primary())
	session, err := c.Client.StartSession(opts)
	if err != nil {
		return err
	}

	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		return nil, f(sessCtx)
	})

	return err
}

func (c *dbClient) checkTransactionEnabled(pCtx context.Context, logger zerolog.Logger) {
	ctx, cancel := context.WithTimeout(pCtx, time.Second)
	defer cancel()

	session, err := c.Client.StartSession()
	if err != nil {
		logger.Err(err).Msg("cannot determine MongoDB version, transactions are disabled")
		return
	}

	defer session.EndSession(ctx)

	collection := c.Collection("test_collection")
	_, err = collection.InsertOne(ctx, bson.M{})
	if err != nil {
		logger.Err(err).Msg("cannot determine MongoDB version, transactions are disabled")
		return
	}

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		return nil, func(ctx mongo.SessionContext) error {
			_, err = collection.DeleteMany(ctx, bson.M{})
			return err
		}(sessCtx)
	})

	_ = collection.Drop(ctx)

	if err != nil {
		logger.Warn().Err(err).Msg("cannot determine MongoDB version, transactions are disabled")
		return
	}

	logger.Info().Msg("MongoDB version supports transactions, transactions are enabled")
	c.isDistributed = true
}

// IsDistributed returns true if MongoDB is Replica Set or Sharded Cluster.
// Use to check feature availability : Transactions, Change Streams, etc.
func (c *dbClient) IsDistributed() bool {
	return c.isDistributed
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
