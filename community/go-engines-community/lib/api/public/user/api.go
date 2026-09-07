package user

import (
	"context"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/binding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/bulk"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/response"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
)

type API interface {
	bulk.API
	Patch(w http.ResponseWriter, r *http.Request)
	BulkPatch(w http.ResponseWriter, r *http.Request)
}

func NewAPI(
	store Store,
	metricMetaUpdater metrics.MetaUpdater,
	binder binding.Binder,
	responder response.Responder,
	decoder encoding.Decoder,
) API {
	return &api{
		store:             store,
		metricMetaUpdater: metricMetaUpdater,
		binder:            binder,
		responder:         responder,
		decoder:           decoder,
	}
}

type api struct {
	bulk.API
	store             Store
	binder            binding.Binder
	responder         response.Responder
	metricMetaUpdater metrics.MetaUpdater
	decoder           encoding.Decoder
}

func (a *api) List(w http.ResponseWriter, r *http.Request) {
	var d ListRequest
	if err := a.binder.BindForm(r, &d); err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	res, err := a.store.Find(r.Context(), d)
	if err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	a.responder.Respond(w, r, http.StatusOK, res)
}

// Get
// @Success 200 {object} Response
func (a *api) Get(w http.ResponseWriter, r *http.Request) {
	id := a.binder.URLParam(r, "id")
	res, err := a.store.GetByID(r.Context(), id)
	if err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	a.responder.Respond(w, r, http.StatusOK, res)
}

// Create
// @Param body body EditRequest true "body"
// @Success 201 {object} Response
func (a *api) Create(w http.ResponseWriter, r *http.Request) {
	var d CreateRequest
	if err := a.binder.Bind(r, &d); err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	roleIDs, err := authctx.GetRoles(r.Context())
	if err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	res, err := a.store.Insert(r.Context(), d, roleIDs)
	if err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	a.metricMetaUpdater.UpdateById(r.Context(), res.ID)

	a.responder.Respond(w, r, http.StatusCreated, res)
}

// Update
// @Param body body EditRequest true "body"
// @Success 200 {object} Response
func (a *api) Update(w http.ResponseWriter, r *http.Request) {
	var d EditRequest
	if err := a.binder.Bind(r, &d); err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	roleIDs, err := authctx.GetRoles(r.Context())
	if err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	id := a.binder.URLParam(r, "id")
	res, err := a.store.Update(r.Context(), id, d, roleIDs)
	if err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	a.metricMetaUpdater.UpdateById(r.Context(), res.ID)

	a.responder.Respond(w, r, http.StatusOK, res)
}

// Patch
// @Param body body PatchRequest true "body"
// @Success 200 {object} Response
func (a *api) Patch(w http.ResponseWriter, r *http.Request) {
	roleIDs, err := authctx.GetRoles(r.Context())
	if err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	var d PatchRequest
	if err := a.binder.Bind(r, &d); err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	id := a.binder.URLParam(r, "id")
	res, err := a.store.UpdateByPatch(r.Context(), id, d, roleIDs)
	if err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	a.metricMetaUpdater.UpdateById(r.Context(), res.ID)

	a.responder.Respond(w, r, http.StatusOK, res)
}

func (a *api) Delete(w http.ResponseWriter, r *http.Request) {
	userID, err := authctx.GetUserKey(r.Context())
	if err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	roleIDs, err := authctx.GetRoles(r.Context())
	if err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	id := a.binder.URLParam(r, "id")
	err = a.store.Delete(r.Context(), id, userID, roleIDs)
	if err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	a.metricMetaUpdater.DeleteById(r.Context(), id)

	a.responder.Respond(w, r, http.StatusNoContent, nil)
}

// BulkCreate
// @Param body body []CreateRequest true "body"
func (a *api) BulkCreate(w http.ResponseWriter, r *http.Request) {
	roleIDs, err := authctx.GetRoles(r.Context())
	if err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	ids := make([]string, 0)
	bulk.Handler(w, r, func(ctx context.Context, d CreateRequest) (string, error) {
		res, err := a.store.Insert(ctx, d, roleIDs)
		if err != nil {
			return "", err
		}

		ids = append(ids, res.ID)

		return res.ID, nil
	}, a.responder, a.binder, a.decoder)

	if len(ids) > 0 {
		a.metricMetaUpdater.UpdateById(r.Context(), ids...)
	}
}

// BulkUpdate
// @Param body body []BulkUpdateRequestItem true "body"
func (a *api) BulkUpdate(w http.ResponseWriter, r *http.Request) {
	roleIDs, err := authctx.GetRoles(r.Context())
	if err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	ids := make([]string, 0)
	processed := make(map[string]bool)
	bulk.Handler(w, r, func(ctx context.Context, d BulkUpdateRequestItem) (string, error) {
		res, err := a.store.Update(ctx, d.ID, d.EditRequest, roleIDs)
		if err != nil {
			return d.ID, err
		}

		if !processed[res.ID] {
			ids = append(ids, res.ID)
			processed[res.ID] = true
		}

		return d.ID, nil
	}, a.responder, a.binder, a.decoder)

	if len(ids) > 0 {
		a.metricMetaUpdater.UpdateById(r.Context(), ids...)
	}
}

// BulkPatch
// @Param body body []BulkPatchRequestItem true "body"
func (a *api) BulkPatch(w http.ResponseWriter, r *http.Request) {
	roleIDs, err := authctx.GetRoles(r.Context())
	if err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	ids := make([]string, 0)
	processed := make(map[string]bool)
	bulk.Handler(w, r, func(ctx context.Context, d BulkPatchRequestItem) (string, error) {
		res, err := a.store.UpdateByPatch(ctx, d.ID, d.PatchRequest, roleIDs)
		if err != nil {
			return d.ID, err
		}

		if !processed[res.ID] {
			ids = append(ids, res.ID)
			processed[res.ID] = true
		}

		return d.ID, nil
	}, a.responder, a.binder, a.decoder)

	if len(ids) > 0 {
		a.metricMetaUpdater.UpdateById(r.Context(), ids...)
	}
}

// BulkDelete
// @Param body body []bulk.DeleteRequestItem true "body"
func (a *api) BulkDelete(w http.ResponseWriter, r *http.Request) {
	roleIDs, err := authctx.GetRoles(r.Context())
	if err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	ids := make([]string, 0)
	processed := make(map[string]bool)
	bulk.Handler(w, r, func(ctx context.Context, d bulk.DeleteRequestItem) (string, error) {
		err := a.store.Delete(ctx, d.ID, d.Author, roleIDs)
		if err != nil {
			return d.ID, err
		}

		if !processed[d.ID] {
			ids = append(ids, d.ID)
			processed[d.ID] = true
		}

		return d.ID, nil
	}, a.responder, a.binder, a.decoder)

	if len(ids) > 0 {
		a.metricMetaUpdater.DeleteById(r.Context(), ids...)
	}
}
