package techmetrics

import (
	"context"
	"errors"
	"os"
	"sync"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/postgres"
	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog"
)

const (
	TaskStatusNone = iota
	TaskStatusRunning
	TaskStatusSucceeded
	TaskStatusFailed
)

const (
	filenameDumpPattern = "cps_tech_metrics_*.bak"
)

// TaskExecutor is used to implement export task executor.
type TaskExecutor interface {
	Run(ctx context.Context)
	// StartExecute creates new export task.
	StartExecute(ctx context.Context) (Task, error)
	// GetStatus returns export task status.
	GetStatus(ctx context.Context) (Task, error)
}

type Task struct {
	ID        int
	Status    int
	Filepath  string
	Created   time.Time
	Started   *time.Time
	Completed *time.Time
}

func NewTaskExecutor(
	pgConnStr string,
	configProvider config.TechMetricsConfigProvider,
	logger zerolog.Logger,
) TaskExecutor {
	return &taskExecutor{
		pgConnStr:      pgConnStr,
		configProvider: configProvider,
		logger:         logger,
	}
}

type taskExecutor struct {
	pgConnStr      string
	logger         zerolog.Logger
	configProvider config.TechMetricsConfigProvider

	pgPoolMx     sync.Mutex
	pgPool       postgres.Pool
	pgPoolClosed bool

	executeChMx sync.Mutex
	executeCh   chan struct{}
}

func (e *taskExecutor) Run(ctx context.Context) {
	e.pgPoolMx.Lock()
	e.pgPoolClosed = false
	e.pgPoolMx.Unlock()
	ch := make(chan struct{}, 1)
	defer close(ch)
	e.executeChMx.Lock()
	e.executeCh = ch
	e.executeChMx.Unlock()

	defer e.closePgPool()

	e.executeLastTask(ctx)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case _, ok := <-ch:
				if !ok {
					return
				}

				e.executeLastTask(ctx)
			}
		}
	}()

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

func (e *taskExecutor) StartExecute(ctx context.Context) (Task, error) {
	pgPool, err := e.getPgPool(ctx)
	if err != nil {
		return Task{}, err
	}

	now := time.Now().In(time.UTC)
	var task Task
	err = pgPool.WithTransaction(ctx, func(ctx context.Context, tx pgx.Tx) error {
		task = Task{}
		var runningTaskId int
		res := tx.QueryRow(ctx, "SELECT id FROM export WHERE status = $1", TaskStatusRunning)
		err = res.Scan(&runningTaskId)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
		if runningTaskId > 0 {
			return nil
		}

		task = Task{Status: TaskStatusRunning, Created: now}
		res = tx.QueryRow(ctx, "INSERT INTO export (status, created, filepath) VALUES ($1, $2, '') RETURNING id", task.Status, task.Created)
		err = res.Scan(&task.ID)
		return err
	})

	if err != nil || task.ID == 0 {
		return task, err
	}

	e.executeChMx.Lock()
	if e.executeCh != nil {
		e.executeCh <- struct{}{}
	}
	e.executeChMx.Unlock()

	return task, nil
}

func (e *taskExecutor) GetStatus(ctx context.Context) (Task, error) {
	pgPool, err := e.getPgPool(ctx)
	if err != nil {
		return Task{}, err
	}
	res := pgPool.QueryRow(ctx, "SELECT id, filepath, status, created, started, completed FROM export ORDER BY created DESC LIMIT 1")
	task := Task{}
	err = res.Scan(&task.ID, &task.Filepath, &task.Status, &task.Created, &task.Started, &task.Completed)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return task, err
	}
	return task, nil
}

func (e *taskExecutor) executeLastTask(ctx context.Context) {
	pgPool, err := e.getPgPool(ctx)
	if err != nil {
		e.logger.Err(err).Msg("cannot connect to postgres")
		return
	}

	now := time.Now().In(time.UTC)
	var lastTaskId int
	res := pgPool.QueryRow(ctx, "UPDATE export SET started = $2 WHERE status = $1 RETURNING id", TaskStatusRunning, now)
	err = res.Scan(&lastTaskId)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		e.logger.Err(err).Msg("cannot update status")
		return
	}
	if lastTaskId == 0 {
		return
	}

	filepath, err := e.dumpDb(ctx)
	status := TaskStatusSucceeded
	if err != nil {
		e.logger.Err(err).Msg("cannot dump tech metrics")
		status = TaskStatusFailed
	}

	_, err = pgPool.Exec(ctx, "UPDATE export SET status = $2, completed = $3, filepath = $4 WHERE id = $1", lastTaskId, status, now, filepath)
	if err != nil {
		e.logger.Err(err).Msg("cannot update status")
	}
}

func (e *taskExecutor) deleteTasks(ctx context.Context) {
	pgPool, err := e.getPgPool(ctx)
	if err != nil {
		e.logger.Err(err).Msg("cannot connect to postgres")
		return
	}

	now := time.Now().In(time.UTC)
	date := now.Add(-e.configProvider.Get().DumpKeepInterval)
	rows, err := pgPool.Query(ctx, "SELECT id, filepath FROM export WHERE status != $1 AND completed < $2", TaskStatusRunning, date)
	if err != nil {
		e.logger.Err(err).Msg("cannot get tasks")
		return
	}
	defer rows.Close()

	ids := make([]int, 0)
	for rows.Next() {
		task := Task{}
		err := rows.Scan(&task.ID, &task.Filepath)
		if err != nil {
			e.logger.Err(err).Msg("cannot scan task")
			return
		}
		ids = append(ids, task.ID)

		if task.Filepath != "" {
			err = os.Remove(task.Filepath)
			if err != nil && !os.IsNotExist(err) {
				e.logger.Err(err).Msg("cannot remove dump")
				return
			}
		}
	}

	if len(ids) > 0 {
		_, err = pgPool.Exec(ctx, "DELETE FROM export WHERE id = ANY($1)", ids)
		if err != nil {
			e.logger.Err(err).Msg("cannot remove tasks")
			return
		}
	}
}

func (e *taskExecutor) getPgPool(ctx context.Context) (postgres.Pool, error) {
	e.pgPoolMx.Lock()
	defer e.pgPoolMx.Unlock()

	if e.pgPoolClosed {
		return nil, errors.New("postgres client is closed")
	}

	if e.pgPool == nil {
		var err error
		e.pgPool, err = postgres.NewPoolByConnStr(ctx, e.pgConnStr, 0, 0)
		if err != nil {
			return nil, err
		}
	}

	return e.pgPool, nil
}

func (e *taskExecutor) closePgPool() {
	e.pgPoolMx.Lock()
	defer e.pgPoolMx.Unlock()

	if e.pgPool != nil {
		e.pgPool.Close()
		e.pgPool = nil
		e.pgPoolClosed = true
	}
}

func (e *taskExecutor) dumpDb(ctx context.Context) (string, error) {
	var err error
	done := make(chan struct{})
	var dumpFilepath string
	go func() {
		defer close(done)

		var f *os.File
		f, err = os.CreateTemp("", filenameDumpPattern)
		if err != nil {
			return
		}
		dumpFilepath = f.Name()
		err = f.Close()
		if err != nil {
			return
		}

		err = postgres.Dump(e.pgConnStr, dumpFilepath)
		if err != nil {
			return
		}
	}()

	select {
	case <-ctx.Done():
		return "", nil
	case <-done:
		if err != nil {
			return "", err
		}
		return dumpFilepath, nil
	}
}
