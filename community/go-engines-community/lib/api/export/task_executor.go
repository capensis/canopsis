package export

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/workers"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const FileNameTimeLayout = "2006-01-02T15-04-05-MST"

// TaskCreator is used to create export task.
type TaskCreator interface {
	Create(ctx context.Context, t TaskParameters) (*Task, error)
	Get(ctx context.Context, id string) (*Task, error)
	SetFormatter(t string, f OutputFormatter)
	GetFormatter(t string) OutputFormatter
}

type TaskExecutor interface {
	TaskCreator
	RegisterType(t string, fetch FetchData)
	ExecuteTask(ctx context.Context, id string) error
	ProcessAbandonedTasks(ctx context.Context)
	DeleteOldTasks(ctx context.Context)
}

type FetchData func(ctx context.Context, t Task) (DataCursor, error)

type DataCursor interface {
	Next(ctx context.Context) bool
	Scan(*map[string]any) error
	Close(ctx context.Context) error
}

type OutputFormatter interface {
	GetFileExtension() string
	DataFetcher(context.Context, FieldsSeparatorGetter, DataCursor) (string, error)
}

type FieldsSeparatorGetter interface {
	GetFields() Fields
	GetSeparator() rune
}

func NewTaskExecutor(
	client mongo.DbClient,
	jobPublisher workers.JobPublisher,
	timezoneConfigProvider config.TimezoneConfigProvider,
	dir string,
	logger zerolog.Logger,
) TaskExecutor {
	return &taskExecutor{
		client:       client,
		collection:   client.Collection(mongo.ExportTaskMongoCollection),
		jobPublisher: jobPublisher,
		logger:       logger,
		workerCount:  10,

		abandonedInterval:         time.Minute,
		abandonedLaunchedInterval: 5 * time.Minute,
		removeInterval:            time.Hour,

		fetches: make(map[string]FetchData),

		formatter:              &csvFormatter{dir: dir}, // default formatter, can be changed by SetFormatter
		customFormatter:        make(map[string]OutputFormatter),
		timezoneConfigProvider: timezoneConfigProvider,
	}
}

type taskExecutor struct {
	client       mongo.DbClient
	collection   mongo.DbCollection
	jobPublisher workers.JobPublisher
	workerCount  int
	logger       zerolog.Logger

	fetches map[string]FetchData

	abandonedInterval         time.Duration
	abandonedLaunchedInterval time.Duration
	removeInterval            time.Duration

	formatter              OutputFormatter
	customFormatter        map[string]OutputFormatter
	timezoneConfigProvider config.TimezoneConfigProvider
}

func (e *taskExecutor) RegisterType(t string, fetch FetchData) {
	if _, ok := e.fetches[t]; ok {
		panic(fmt.Errorf("type %q is already registered", t))
	}

	e.fetches[t] = fetch
}

func (e *taskExecutor) Create(ctx context.Context, params TaskParameters) (*Task, error) {
	location := e.timezoneConfigProvider.Get().Location
	now := datetime.NewCpsTime().In(location)
	t := Task{
		ID:         utils.NewID(),
		Status:     TaskStatusRunning,
		Type:       params.Type,
		Parameters: params.Parameters,
		Fields:     params.Fields,
		Separator:  params.Separator,
		Filename: params.FilenamePrefix + "-" + now.Time.Format(FileNameTimeLayout) +
			e.GetFormatter(params.Type).GetFileExtension(),
		User:       params.UserID,
		TimeFormat: params.TimeFormat,
		Created:    now,
	}

	_, err := e.collection.InsertOne(ctx, t)
	if err != nil {
		return nil, err
	}

	err = e.jobPublisher.Publish(ctx, t.ID)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (e *taskExecutor) Get(ctx context.Context, id string) (*Task, error) {
	res := e.collection.FindOne(ctx, bson.M{"_id": id})
	if err := res.Err(); err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, nil
		}

		return nil, err
	}

	t := Task{}
	err := res.Decode(&t)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (e *taskExecutor) ExecuteTask(ctx context.Context, id string) error {
	t := Task{}
	err := e.collection.FindOneAndUpdate(
		ctx,
		bson.M{
			"_id":    id,
			"status": TaskStatusRunning,
		},
		bson.M{"$set": bson.M{
			"launched": datetime.NewCpsTime(),
		}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&t)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil
		}
		return err
	}

	updateFilter := bson.M{
		"_id":    id,
		"status": TaskStatusRunning,
	}
	fetch := e.fetches[t.Type]
	if fetch == nil {
		_, err = e.collection.UpdateOne(ctx, updateFilter, bson.M{"$set": bson.M{
			"status":      TaskStatusFailed,
			"completed":   datetime.NewCpsTime(),
			"fail_reason": "unknown type: " + t.Type,
		}})
		return err
	}

	cursor, err := fetch(ctx, t)
	if err != nil {
		_, updateErr := e.collection.UpdateOne(ctx, updateFilter, bson.M{"$set": bson.M{
			"status":      TaskStatusFailed,
			"completed":   datetime.NewCpsTime(),
			"fail_reason": "cannot fetch data: " + err.Error(),
		}})
		if updateErr != nil {
			e.logger.Err(updateErr).Msg("cannot update export task")
		}

		valErr := &validation.Error{}
		if errors.As(err, &valErr) {
			err = fmt.Errorf("invalid params %s: %w", t.Parameters, err)
		}

		return err
	}

	fileName, err := e.GetFormatter(t.Type).DataFetcher(ctx, t, cursor)
	if err != nil {
		_, updateErr := e.collection.UpdateOne(ctx, updateFilter, bson.M{"$set": bson.M{
			"status":      TaskStatusFailed,
			"completed":   datetime.NewCpsTime(),
			"fail_reason": "cannot fetch data: " + err.Error(),
		}})
		if updateErr != nil {
			e.logger.Err(updateErr).Msg("cannot update export task")
		}

		return fmt.Errorf("cannot export data: %w", err)
	}

	_, err = e.collection.UpdateOne(ctx, updateFilter, bson.M{"$set": bson.M{
		"status":    TaskStatusSucceeded,
		"file":      fileName,
		"completed": datetime.NewCpsTime(),
	}})
	if err != nil {
		return fmt.Errorf("cannot update export task: %w", err)
	}

	return nil
}

func (e *taskExecutor) ProcessAbandonedTasks(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := e.fetchTasks(ctx)
			if err != nil {
				e.logger.Err(err).Msg("cannot fetch export tasks")
			}
		}
	}
}

func (e *taskExecutor) DeleteOldTasks(ctx context.Context) {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := e.deleteTasks(ctx)
			if err != nil {
				e.logger.Err(err).Msg("cannot delete export tasks")
			}
		}
	}
}

func (e *taskExecutor) SetFormatter(t string, f OutputFormatter) {
	e.customFormatter[t] = f
}

func (e *taskExecutor) GetFormatter(t string) OutputFormatter {
	formatter, ok := e.customFormatter[t]
	if !ok {
		formatter = e.formatter
	}

	return formatter
}

func (e *taskExecutor) fetchTasks(ctx context.Context) error {
	cursor, err := e.collection.Find(ctx, bson.M{"$or": []bson.M{
		{
			"status":   TaskStatusRunning,
			"launched": nil,
			"started":  bson.M{"$lte": datetime.CpsTime{Time: time.Now().Add(-e.abandonedInterval)}},
		},
		{
			"status": TaskStatusRunning,
			"launched": bson.M{
				"$gt":  0,
				"$lte": datetime.CpsTime{Time: time.Now().Add(-e.abandonedLaunchedInterval)},
			},
		},
	}})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		t := Task{}
		err = cursor.Decode(&t)
		if err != nil {
			return err
		}

		err = e.jobPublisher.Publish(ctx, t.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *taskExecutor) deleteTasks(ctx context.Context) error {
	cursor, err := e.collection.Find(ctx, bson.M{
		"status":    bson.M{"$in": bson.A{TaskStatusSucceeded, TaskStatusFailed}},
		"completed": bson.M{"$lte": datetime.CpsTime{Time: time.Now().Add(-e.removeInterval)}},
	})
	if err != nil {
		return fmt.Errorf("cannot find export tasks to delete: %w", err)
	}

	defer cursor.Close(ctx)

	ids := make([]string, 0)
	for cursor.Next(ctx) {
		t := struct {
			ID   string `bson:"_id"`
			File string `bson:"file"`
		}{}

		err := cursor.Decode(&t)
		if err != nil {
			return fmt.Errorf("cannot decode export task: %w", err)
		}

		if t.File != "" {
			err = os.Remove(t.File)
			if err != nil && !os.IsNotExist(err) {
				e.logger.Err(err).Msg("cannot remove export file")
				continue
			}
		}

		ids = append(ids, t.ID)
	}

	if len(ids) == 0 {
		return nil
	}

	_, err = e.collection.DeleteMany(ctx, bson.M{
		"_id": bson.M{"$in": ids},
	})
	if err != nil {
		return fmt.Errorf("cannot delete export tasks: %w", err)
	}

	return nil
}
