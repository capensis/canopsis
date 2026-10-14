package entity

import (
	"context"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/binding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/bulk"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/response"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
)

type API interface {
	List(w http.ResponseWriter, r *http.Request)
	Search(w http.ResponseWriter, r *http.Request)
	Pause(w http.ResponseWriter, r *http.Request)
	Resume(w http.ResponseWriter, r *http.Request)
	BulkPause(w http.ResponseWriter, r *http.Request)
	BulkResume(w http.ResponseWriter, r *http.Request)
}

func NewAPI(
	store Store,
	binder binding.Binder,
	responder response.Responder,
	decoder encoding.Decoder,
) API {
	return &api{
		store:     store,
		binder:    binder,
		responder: responder,
		decoder:   decoder,
	}
}

type api struct {
	store     Store
	binder    binding.Binder
	responder response.Responder
	decoder   encoding.Decoder
}

// List
//
// @Success 200 {object} Response
func (a *api) List(w http.ResponseWriter, r *http.Request) {
	userID, err := authctx.GetUserKey(r.Context())
	if err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	var d PaginatedListRequest
	if err := a.binder.BindForm(r, &d); err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	res, err := a.store.Find(r.Context(), d, userID)
	if err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	a.responder.Respond(w, r, http.StatusOK, res)
}

// Search
//
// @Param body body SearchRequest true "body"
// @Success 200 {object} Response
func (a *api) Search(w http.ResponseWriter, r *http.Request) {
	userID, err := authctx.GetUserKey(r.Context())
	if err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	var d PaginatedSearchRequest
	d.Query.SetDefaults()
	if err := a.binder.Bind(r, &d); err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	res, err := a.store.Search(r.Context(), d, userID)
	if err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	a.responder.Respond(w, r, http.StatusOK, res)
}

// Pause
// @Param body body PauseRequest true "body"
func (a *api) Pause(w http.ResponseWriter, r *http.Request) {
	var d PauseRequest
	if err := a.binder.Bind(r, &d); err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	userID, err := authctx.GetUserKey(r.Context())
	if err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	username, err := authctx.GetUsername(r.Context())
	if err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	err = a.store.Pause(r.Context(), d, userID, username)
	if err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	a.responder.Respond(w, r, http.StatusNoContent, nil)
}

// Resume
// @Param body body ResumeRequest true "body"
func (a *api) Resume(w http.ResponseWriter, r *http.Request) {
	var d ResumeRequest
	if err := a.binder.Bind(r, &d); err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	userID, err := authctx.GetUserKey(r.Context())
	if err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	username, err := authctx.GetUsername(r.Context())
	if err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	err = a.store.Resume(r.Context(), d, userID, username)
	if err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	a.responder.Respond(w, r, http.StatusNoContent, nil)
}

// BulkPause
// @Param body body []PauseRequest true "body"
func (a *api) BulkPause(w http.ResponseWriter, r *http.Request) {
	userID, err := authctx.GetUserKey(r.Context())
	if err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	username, err := authctx.GetUsername(r.Context())
	if err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	bulk.Handler(w, r, func(ctx context.Context, d PauseRequest) (string, error) {
		err := a.store.Pause(ctx, d, userID, username)
		if err != nil {
			return "", err
		}

		return d.ID, nil
	}, a.responder, a.binder, a.decoder)
}

// BulkResume
// @Param body body []ResumeRequest true "body"
func (a *api) BulkResume(w http.ResponseWriter, r *http.Request) {
	userID, err := authctx.GetUserKey(r.Context())
	if err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	username, err := authctx.GetUsername(r.Context())
	if err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	bulk.Handler(w, r, func(ctx context.Context, d ResumeRequest) (string, error) {
		err := a.store.Resume(ctx, d, userID, username)
		if err != nil {
			return "", err
		}

		return d.ID, nil
	}, a.responder, a.binder, a.decoder)
}
