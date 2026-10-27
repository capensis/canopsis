package alarmexport

import (
	"cmp"
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/patternfields"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/asyncjob"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/binding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/dataexport"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/response"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/sortparam"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

type api struct {
	dataexport.API
	transformer patternfields.Transformer
}

// NewAPI
//
// @Param body body Request true "body"
// @Success 200 {object} asyncjob.Response
func NewAPI(
	store *asyncjob.Store[Job],
	jobPublisher asyncjob.Publisher,
	dir string,
	transformer patternfields.Transformer,
	binder binding.Binder,
	responder response.Responder,
	logger zerolog.Logger,
) dataexport.API {
	a := &api{
		transformer: transformer,
	}
	a.API = dataexport.NewAPI(asyncjob.TypeAlarmExport, store, a.prepare, a.newJob, jobPublisher, dir, binder, responder, logger)

	return a
}

func (a *api) prepare(ctx context.Context, r *Request) error {
	if r.Sort != "" {
		_, err := sortparam.ValidateResolveFunc(cmp.Or(r.Sort, alarm.DefaultSort), func(f string) (string, bool) {
			rf, _, isSortable := alarm.ResolveSortableField(f, "")

			return rf, isSortable
		}, alarm.SortErrParam())
		if err != nil {
			return err
		}
	}

	if len(r.EntityPattern) > 0 {
		patternAliases := patternfields.GetAliases(r.EntityPattern)
		if len(patternAliases) != 0 {
			aliases, err := a.transformer.FetchAliases(ctx, patternAliases)
			if err != nil {
				return err
			}

			var valErrs validator.ValidationErrors
			r.EntityPattern, _, valErrs = a.transformer.ApplyAliases(r.EntityPattern, aliases, "EntityPattern")
			if len(valErrs) > 0 {
				return validation.NewError(valErrs, r)
			}
		}
	}

	return nil
}

func (a *api) newJob(r Request) Job {
	return Job{
		Base:   asyncjob.NewBase(utils.NewID(), r.Author),
		Params: r,
	}
}
