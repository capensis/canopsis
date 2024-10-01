package export

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const fileNameTimeLayout = "2006-01-02T15-04-05-MST"

// TaskExecutor is used to implement export task executor.
type TaskExecutor interface {
	RegisterType(t string, fetch FetchData)
	// Execute receives tasks from channel and save its result to storage.
	Execute(ctx context.Context)
	// StartExecute creates new export task.
	StartExecute(ctx context.Context, t TaskParameters) (*Task, error)
	Get(ctx context.Context, id string) (*Task, error)
	SetFormatter(t string, f OutputFormatter)
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
	timezoneConfigProvider config.TimezoneConfigProvider,
	logger zerolog.Logger,
) TaskExecutor {
	return &taskExecutor{
		client:      client,
		collection:  client.Collection(mongo.ExportTaskMongoCollection),
		logger:      logger,
		workerCount: 10,

		abandonedInterval:         time.Minute,
		abandonedLaunchedInterval: 5 * time.Minute,
		removeInterval:            5 * time.Minute,

		fetches: make(map[string]FetchData),

		formatter:              &csvFormatter{}, // default formatter, can be changed by SetFormatter
		customFormatter:        make(map[string]OutputFormatter),
		timezoneConfigProvider: timezoneConfigProvider,
	}
}

type taskExecutor struct {
	client      mongo.DbClient
	collection  mongo.DbCollection
	workerCount int
	logger      zerolog.Logger

	fetches map[string]FetchData

	abandonedInterval         time.Duration
	abandonedLaunchedInterval time.Duration
	removeInterval            time.Duration

	taskCh   chan<- string
	taskChMx sync.Mutex

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

func (e *taskExecutor) Execute(parentCtx context.Context) {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	ch := make(chan string)
	e.taskChMx.Lock()
	e.taskCh = ch
	e.taskChMx.Unlock()

	defer func() {
		e.taskChMx.Lock()
		e.taskCh = nil
		close(ch)
		e.taskChMx.Unlock()
	}()

	wg := sync.WaitGroup{}
	// Run export workers
	for i := 0; i < e.workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case t, ok := <-ch:
					if !ok {
						return
					}

					err := e.executeTask(ctx, t)
					if err != nil {
						e.logger.Err(err).Msg("cannot execute export task")
					}
				}
			}
		}()
	}

	// Fetch tasks
	wg.Add(1)
	go func() {
		defer wg.Done()

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
	}()

	// Run delete worker
	wg.Add(1)
	go func() {
		defer wg.Done()

		ticker := time.NewTicker(time.Minute)
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
	}()

	wg.Wait()
}

func (e *taskExecutor) StartExecute(ctx context.Context, params TaskParameters) (*Task, error) {
	e.taskChMx.Lock()
	defer e.taskChMx.Unlock()

	if e.taskCh == nil {
		return nil, errors.New("execute is not running")
	}

	location := e.timezoneConfigProvider.Get().Location
	now := datetime.NewCpsTime().In(location)
	t := Task{
		ID:         utils.NewID(),
		Status:     TaskStatusRunning,
		Type:       params.Type,
		Parameters: params.Parameters,
		Fields:     params.Fields,
		Separator:  params.Separator,
		Filename: params.FilenamePrefix + "-" + now.Time.Format(fileNameTimeLayout) +
			e.getFormatter(params.Type).GetFileExtension(),
		User:    params.UserID,
		Created: now,
	}

	_, err := e.collection.InsertOne(ctx, t)
	if err != nil {
		return nil, err
	}

	select {
	case e.taskCh <- t.ID:
	default:
	}

	return &t, nil
}

func (e *taskExecutor) getFormatter(t string) OutputFormatter {
	formatter, ok := e.customFormatter[t]
	if !ok {
		formatter = e.formatter
	}

	return formatter
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

func (e *taskExecutor) executeTask(ctx context.Context, id string) error {
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
		_, err := e.collection.UpdateOne(ctx, updateFilter, bson.M{"$set": bson.M{
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

		return err
	}

	fileName, err := e.getFormatter(t.Type).DataFetcher(ctx, t, cursor)
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

func (e *taskExecutor) fetchTasks(ctx context.Context) error {
	e.taskChMx.Lock()
	defer e.taskChMx.Unlock()

	if e.taskCh == nil {
		return errors.New("execute is not running")
	}

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

		select {
		case <-ctx.Done():
			return nil
		case e.taskCh <- t.ID:
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

func (e *taskExecutor) SetFormatter(t string, f OutputFormatter) {
	e.customFormatter[t] = f
}
