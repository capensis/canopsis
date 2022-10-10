package export

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

const (
	TaskStatusRunning = iota
	TaskStatusSucceeded
	TaskStatusFailed
)

const fileNameTimeLayout = "2006-01-02T15-04-05-MST"

// TaskExecutor is used to implement export task executor.
type TaskExecutor interface {
	// StartExecute creates new export task.
	StartExecute(ctx context.Context, t Task) (string, error)
	// Execute receives tasks from channel and save its result to storage.
	Execute(ctx context.Context)
	// GetStatus returns export task status.
	GetStatus(ctx context.Context, id string) (*TaskStatus, error)
}

type Task struct {
	Filename     string
	ExportFields Fields
	Separator    rune
	DataFetcher  DataFetcher
	DataCursor   DataCursor
}

type Fields []Field

type Field struct {
	Name  string `json:"name"`
	Label string `json:"label"`
}

func (f *Fields) Fields() []string {
	fields := make([]string, len(*f))
	for i, field := range *f {
		fields[i] = field.Name
	}

	return fields
}

func (f *Fields) Labels() []string {
	labels := make([]string, len(*f))
	for i, field := range *f {
		labels[i] = field.Label
	}

	return labels
}

type DataFetcher func(ctx context.Context, page, limit int64) ([]map[string]string, int64, error)

type DataCursor interface {
	Next(ctx context.Context) bool
	Scan(*map[string]string) error
	Close()
}

type taskWithID struct {
	Task
	ID string
}

type TaskStatus struct {
	Status   int    `bson:"status"`
	File     string `bson:"file"`
	Filename string `bson:"filename"`
}

func NewTaskExecutor(
	client mongo.DbClient,
	timezoneConfigProvider config.TimezoneConfigProvider,
	logger zerolog.Logger,
) TaskExecutor {
	return &taskExecutor{
		client:         client,
		collection:     client.Collection(mongo.ExportTaskMongoCollection),
		logger:         logger,
		workerCount:    10,
		removeInterval: 5 * time.Minute,

		timezoneConfigProvider: timezoneConfigProvider,
	}
}

type taskExecutor struct {
	client         mongo.DbClient
	collection     mongo.DbCollection
	logger         zerolog.Logger
	workerCount    int
	removeInterval time.Duration
	taskCh         chan<- taskWithID
	taskChMx       sync.Mutex

	timezoneConfigProvider config.TimezoneConfigProvider
}

func (e *taskExecutor) Execute(parentCtx context.Context) {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	ch := make(chan taskWithID)
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

					e.executeTask(ctx, t)
				}
			}
		}()
	}

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
				e.deleteTasks(ctx)
			}
		}
	}()

	wg.Wait()
}

func (e *taskExecutor) StartExecute(ctx context.Context, t Task) (string, error) {
	e.taskChMx.Lock()
	defer e.taskChMx.Unlock()

	if e.taskCh == nil {
		return "", errors.New("execute is not running")
	}

	id := utils.NewID()
	now := time.Now().UTC()
	location := e.timezoneConfigProvider.Get().Location
	filename := fmt.Sprintf("%s-%s.csv", t.Filename, now.In(location).Format(fileNameTimeLayout))
	_, err := e.collection.InsertOne(ctx, bson.M{
		"_id":      id,
		"status":   TaskStatusRunning,
		"created":  types.CpsTime{Time: now},
		"filename": filename,
	})
	if err != nil {
		e.logger.Err(err).Msg("cannot save export task")
		return "", err
	}

	e.taskCh <- taskWithID{Task: t, ID: id}

	return id, nil
}

func (e *taskExecutor) GetStatus(ctx context.Context, id string) (*TaskStatus, error) {
	res := e.collection.FindOne(ctx, bson.M{"_id": id})
	if err := res.Err(); err != nil {
		if err == mongodriver.ErrNoDocuments {
			return nil, nil
		}

		return nil, err
	}

	t := TaskStatus{}
	err := res.Decode(&t)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (e *taskExecutor) executeTask(ctx context.Context, t taskWithID) {
	var err error
	var fileName string

	if t.DataFetcher != nil {
		fileName, err = ExportCsv(ctx, t.ExportFields, t.Separator, t.DataFetcher)
	} else {
		fileName, err = ExportCsvByCursor(ctx, t.ExportFields, t.Separator, t.DataCursor)
	}

	if err != nil {
		e.logger.Err(err).Msg("cannot export data")

		_, err := e.collection.UpdateOne(ctx, bson.M{"_id": t.ID}, bson.M{"$set": bson.M{
			"status":    TaskStatusFailed,
			"completed": types.CpsTime{Time: time.Now()},
		}})
		if err != nil {
			e.logger.Err(err).Msg("cannot update export task")
			return
		}

		return
	}

	_, err = e.collection.UpdateOne(ctx, bson.M{"_id": t.ID}, bson.M{"$set": bson.M{
		"status": TaskStatusSucceeded,
		"file":   fileName, "completed": types.CpsTime{Time: time.Now()},
	}})
	if err != nil {
		e.logger.Err(err).Msg("cannot update export task")
		return
	}
}

func (e *taskExecutor) deleteTasks(ctx context.Context) {
	cursor, err := e.collection.Find(ctx, bson.M{
		"completed": bson.M{"$lte": types.CpsTime{Time: time.Now().Add(-e.removeInterval)}},
	})
	if err != nil {
		e.logger.Err(err).Msg("cannot find export tasks to delete")
		return
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
			e.logger.Err(err).Msg("cannot decode export task")
			return
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
		return
	}

	_, err = e.collection.DeleteMany(ctx, bson.M{
		"_id": bson.M{"$in": ids},
	})
	if err != nil {
		e.logger.Err(err).Msg("cannot delete export tasks")
		return
	}
}
