// Package dataexport provides the generic http API to export a resource to a file, and the writers producing that file.
package dataexport

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/asyncjob"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/binding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/response"
	"github.com/rs/zerolog"
)

type API interface {
	Start(w http.ResponseWriter, r *http.Request)
	Status(w http.ResponseWriter, r *http.Request)
	Download(w http.ResponseWriter, r *http.Request)
}

func NewAPI[TJob Job, TCreate any](
	jobType string,
	store *asyncjob.Store[TJob],
	prepare func(context.Context, *TCreate) error,
	newJob func(TCreate) TJob,
	jobPublisher asyncjob.Publisher,
	dir string,
	binder binding.Binder,
	responder response.Responder,
	logger zerolog.Logger,
) API {
	return &api[TJob, TCreate]{
		jobType:      jobType,
		store:        store,
		prepare:      prepare,
		newJob:       newJob,
		jobPublisher: jobPublisher,
		dir:          dir,
		binder:       binder,
		responder:    responder,
		logger:       logger,
	}
}

type api[TJob Job, TCreate any] struct {
	jobType      string
	store        *asyncjob.Store[TJob]
	prepare      func(context.Context, *TCreate) error
	newJob       func(TCreate) TJob
	jobPublisher asyncjob.Publisher
	dir          string
	binder       binding.Binder
	responder    response.Responder
	logger       zerolog.Logger
}

func (a *api[TJob, TCreate]) Start(w http.ResponseWriter, r *http.Request) {
	var d TCreate
	if err := a.binder.Bind(r, &d); err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	if a.prepare != nil {
		err := a.prepare(r.Context(), &d)
		if err != nil {
			a.responder.RespondError(w, r, err)

			return
		}
	}

	job := a.newJob(d)
	err := a.store.Create(r.Context(), job)
	if err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	a.logger.Debug().Str("job_id", job.JobID()).Msgf("export job %s created", a.jobType)

	err = a.jobPublisher.Publish(r.Context(), a.jobType, job.JobID())
	if err != nil {
		a.responder.RespondError(w, r, fmt.Errorf("cannot signal export job %q: %w", job.JobID(), err))

		return
	}

	a.responder.Respond(w, r, http.StatusOK, job.Response())
}

func (a *api[TJob, TCreate]) Status(w http.ResponseWriter, r *http.Request) {
	id := a.binder.URLParam(r, "id")
	job, err := a.store.Get(r.Context(), id)
	if err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	a.responder.Respond(w, r, http.StatusOK, job.Response())
}

func (a *api[TJob, TCreate]) Download(w http.ResponseWriter, r *http.Request) {
	id := a.binder.URLParam(r, "id")
	job, err := a.store.Get(r.Context(), id)
	if err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	if job.JobStatus() != asyncjob.StatusDone {
		a.responder.RespondError(w, r, httperror.ErrNotFound)

		return
	}

	mt, format := job.MimeType()
	w.Header().Set("Content-Disposition", `attachment; filename*=UTF-8''`+url.PathEscape(job.DownloadName()))
	w.Header().Set("Content-Type", mt)
	http.ServeFile(w, r, Filepath(a.dir, job.JobID(), format))
}
