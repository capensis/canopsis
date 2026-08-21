package entityimport

import (
	"context"
	"errors"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/asyncjob"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/dataimport"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/importcontextgraph"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type handler struct {
	store  *asyncjob.Store[Job]
	dir    string
	worker importcontextgraph.Worker
	logger zerolog.Logger
}

// NewHandler returns the asyncjob.Handler for asyncjob.TypeEntityImport jobs.
func NewHandler(
	store *asyncjob.Store[Job],
	dir string,
	worker importcontextgraph.Worker,
	logger zerolog.Logger,
) asyncjob.Handler {
	return &handler{
		store:  store,
		dir:    dir,
		worker: worker,
		logger: logger,
	}
}

func (h *handler) Handle(ctx context.Context, _ string) error {
	// Only one import job can be active at a time, so the jobs signaled while another one was
	// running are drained here instead.
	id, err := h.handleNext(ctx)
	if err != nil {
		h.logger.Err(err).Str("job_id", id).Msg("cannot process import job")
	}

	if id == "" {
		return nil
	}

	if ctx.Err() != nil {
		return ctx.Err()
	}

	return h.Handle(ctx, "")
}

// handleNext runs the next import job to process.
func (h *handler) handleNext(ctx context.Context) (string, error) {
	job, err := h.store.GetUncompleted(ctx)
	if err != nil {
		if errors.Is(err, httperror.ErrNotFound) {
			return "", nil
		}

		return "", err
	}

	if job.Status == asyncjob.StatusOngoing {
		return "", nil
	}

	ok, err := h.store.SetOngoing(ctx, job.ID)
	if err != nil {
		return job.ID, err
	}

	if !ok {
		// The job isn't pending anymore: it's already been taken by another worker.
		h.logger.Warn().Str("job_id", job.ID).Msg("import job is not pending, skipped")

		return "", nil
	}

	start := time.Now()
	stats, err := h.do(ctx, job)
	if err != nil {
		if ctx.Err() != nil {
			return job.ID, err
		}

		if updErr := h.store.SetFailed(ctx, job.ID, err.Error(), time.Since(start)); updErr != nil {
			h.logger.Err(updErr).Str("job_id", job.ID).Msg("cannot mark import job as failed")
		}

		return job.ID, err
	}

	err = h.store.SetDone(ctx, job.ID, time.Since(start), bson.M{
		"stats": Stats{
			Updated: stats.Updated,
			Deleted: stats.Deleted,
		},
	})
	if err != nil {
		return job.ID, err
	}

	return job.ID, nil
}

func (h *handler) do(ctx context.Context, job Job) (stats importcontextgraph.Stats, err error) {
	path := dataimport.Filepath(h.dir, job.ID)
	err = asyncjob.RunWithPing(
		ctx,
		func(ctx context.Context) error {
			return h.store.Ping(ctx, job.ID)
		},
		func(ctx context.Context) error {
			var err error
			stats, err = h.worker.WorkPartial(ctx, path, job.Source)

			return err
		},
	)

	return stats, err
}
