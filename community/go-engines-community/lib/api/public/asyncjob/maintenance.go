package asyncjob

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"github.com/rs/zerolog"
)

const (
	// jobMaintenanceInterval is how often the interrupted jobs are requeued and the expired ones are deleted
	jobMaintenanceInterval = time.Minute
	// jobStaleAfter is how long a job can go without pinging before the worker which took it is considered dead.
	jobStaleAfter = time.Minute
	// jobRetention is how long a finished job and its files are kept.
	jobRetention = 24 * time.Hour
	// jobMaxAttempts is how many times a job is given to a worker before it's failed for good.
	jobMaxAttempts = 3
)

// Maintainer recovers and cleans up the jobs of a single type, see RunMaintenance.
type Maintainer interface {
	JobType() string
	// Requeue signals again the jobs left ongoing by a worker which died.
	Requeue(ctx context.Context, t datetime.CpsTime, maxAttempts int64) error
	// Clean deletes the jobs which finished before t, and the files they own.
	Clean(ctx context.Context, t datetime.CpsTime) error
}

func NewMaintainer[T Job](
	store *Store[T],
	jobType string,
	publisher Publisher,
	filepath func(job T) string,
) Maintainer {
	return &maintainer[T]{
		store:     store,
		jobType:   jobType,
		publisher: publisher,
		filepath:  filepath,
	}
}

type maintainer[T Job] struct {
	store     *Store[T]
	jobType   string
	publisher Publisher
	filepath  func(job T) string
}

func (m *maintainer[T]) JobType() string {
	return m.jobType
}

func (m *maintainer[T]) Requeue(ctx context.Context, t datetime.CpsTime, maxAttempts int64) error {
	// The jobs a dead worker left ongoing, at the cost of an attempt.
	err := m.signalClaimed(ctx, func() (string, bool, error) {
		return m.store.claimStale(ctx, t, maxAttempts)
	})
	if err != nil {
		return err
	}

	// The jobs whose creation signal was lost. Claimed after the stale ones, which are moved back to
	// pending, so a job isn't signaled twice in the same pass.
	err = m.signalClaimed(ctx, func() (string, bool, error) {
		return m.store.claimUnsignaled(ctx, t)
	})
	if err != nil {
		return err
	}

	return m.store.failStale(ctx, t, maxAttempts)
}

// signalClaimed signals every job claim hands out, until there's none left to claim.
func (m *maintainer[T]) signalClaimed(ctx context.Context, claim func() (string, bool, error)) error {
	for {
		id, ok, err := claim()
		if err != nil {
			return err
		}

		if !ok {
			return nil
		}

		err = m.publisher.Publish(ctx, m.jobType, id)
		if err != nil {
			return fmt.Errorf("cannot signal job %q: %w", id, err)
		}
	}
}

func (m *maintainer[T]) Clean(ctx context.Context, t datetime.CpsTime) error {
	jobs, err := m.store.FindCompletedBefore(ctx, t)
	if err != nil {
		return err
	}

	for _, job := range jobs {
		if m.filepath != nil {
			err = os.Remove(m.filepath(job))
			if err != nil && !errors.Is(err, os.ErrNotExist) {
				return err
			}
		}

		_, err = m.store.DeleteCompleted(ctx, job.JobID(), t)
		if err != nil {
			return err
		}
	}

	return nil
}

// RunMaintenance requeues the jobs interrupted by a dead worker and deletes the ones which have
// been finished for longer than the retention, every opts.Interval, until ctx is done.
func RunMaintenance(ctx context.Context, maintainers []Maintainer, logger zerolog.Logger) {
	ticker := time.NewTicker(jobMaintenanceInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			now := datetime.NewCpsTime()
			for _, m := range maintainers {
				err := m.Requeue(ctx, datetime.NewCpsTime(now.Add(-jobStaleAfter).Unix()), jobMaxAttempts)
				if err != nil {
					logger.Err(err).Str("type", m.JobType()).Msg("cannot requeue the interrupted jobs")
				}

				err = m.Clean(ctx, datetime.NewCpsTime(now.Add(-jobRetention).Unix()))
				if err != nil {
					logger.Err(err).Str("type", m.JobType()).Msg("cannot delete the old jobs")
				}
			}
		}
	}
}
