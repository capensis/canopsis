package entityimport

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/asyncjob"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/binding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/dataimport"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/response"
	"github.com/rs/zerolog"
)

// NewAPI
//
// @Param body body RequestItem[] true "body"
// @Success 200 {object} asyncjob.Response
func NewAPI(
	store *asyncjob.Store[Job],
	jobPublisher asyncjob.Publisher,
	dir string,
	maxSize uint64,
	binder binding.Binder,
	responder response.Responder,
	logger zerolog.Logger,
) dataimport.API {
	return dataimport.NewAPI[Job](
		asyncjob.TypeEntityImport,
		store,
		newJob,
		jobPublisher,
		dir,
		maxSize,
		binder,
		responder,
		logger,
	)
}

func newJob(id, source string) Job {
	return Job{
		Base:   asyncjob.NewBase(id),
		Source: source,
	}
}
