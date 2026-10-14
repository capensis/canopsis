// Package dataimport provides the generic http API to import a resource from an uploaded file.
package dataimport

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/asyncjob"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/binding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/response"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/rs/zerolog"
)

type API interface {
	Start(w http.ResponseWriter, r *http.Request)
	Status(w http.ResponseWriter, r *http.Request)
}

func NewAPI[TJob Job](
	jobType string,
	store *asyncjob.Store[TJob],
	newJob func(id, source, author string) TJob,
	jobPublisher asyncjob.Publisher,
	dir string,
	maxSize uint64,
	binder binding.Binder,
	responder response.Responder,
	logger zerolog.Logger,
) API {
	return &api[TJob]{
		jobType:      jobType,
		store:        store,
		newJob:       newJob,
		jobPublisher: jobPublisher,
		dir:          dir,
		maxSize:      maxSize,
		binder:       binder,
		responder:    responder,
		logger:       logger,
	}
}

type api[TJob Job] struct {
	jobType      string
	store        *asyncjob.Store[TJob]
	newJob       func(id, source, author string) TJob
	jobPublisher asyncjob.Publisher
	dir          string
	maxSize      uint64
	binder       binding.Binder
	responder    response.Responder
	logger       zerolog.Logger
}

func (a *api[TJob]) Start(w http.ResponseWriter, r *http.Request) {
	var d Query
	if err := a.binder.BindForm(r, &d); err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	body := r.Body
	if a.maxSize > 0 {
		body = http.MaxBytesReader(w, r.Body, int64(a.maxSize))
	}

	job, err := a.createJob(r.Context(), body, d.Source, d.Author)
	if err != nil {
		if _, ok := errors.AsType[*http.MaxBytesError](err); ok {
			err = httperror.ErrRequestEntityTooLarge
		}

		a.responder.RespondError(w, r, err)

		return
	}

	a.responder.Respond(w, r, http.StatusOK, job.Response())
}

func (a *api[TJob]) Status(w http.ResponseWriter, r *http.Request) {
	id := a.binder.URLParam(r, "id")
	job, err := a.store.Get(r.Context(), id)
	if err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	a.responder.Respond(w, r, http.StatusOK, job.Response())
}

func (a *api[TJob]) createJob(ctx context.Context, f io.Reader, source, author string) (job TJob, resErr error) {
	err := os.MkdirAll(a.dir, os.ModeDir|asyncjob.DirPerm)
	if err != nil {
		return job, err
	}

	jobID := utils.NewID()
	path := Filepath(a.dir, jobID)

	df, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, asyncjob.FilePerm)
	if err != nil {
		return job, fmt.Errorf("failed to open file %q: %w", path, err)
	}

	defer func() {
		err := df.Close()
		if err != nil && resErr == nil {
			resErr = fmt.Errorf("failed to close file %q: %w", path, err)
		}
	}()

	n, err := io.Copy(df, f)
	if err != nil {
		_ = os.Remove(path)

		return job, fmt.Errorf("failed to copy to file %q: %w", path, err)
	}

	if n == 0 {
		_ = os.Remove(path)

		return job, validation.NewSingleError("required", "body", "body", nil)
	}

	job = a.newJob(jobID, source, author)
	err = a.store.Create(ctx, job)
	if err != nil {
		_ = os.Remove(path)

		return job, err
	}

	a.logger.Debug().Str("job_id", job.JobID()).Msgf("import job %s created", a.jobType)

	err = a.jobPublisher.Publish(ctx, a.jobType, job.JobID())
	if err != nil {
		return job, fmt.Errorf("cannot signal import job %q: %w", job.JobID(), err)
	}

	return job, nil
}
