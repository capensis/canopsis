package crud

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/binding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/response"
)

type API interface {
	Create(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

func NewAPI[TList, TCreate, TUpdate, TResp any](
	store Store[TList, TCreate, TUpdate, TResp],
	binder binding.Binder,
	responder response.Responder,
) API {
	return &api[TList, TCreate, TUpdate, TResp]{
		store:     store,
		binder:    binder,
		responder: responder,
	}
}

type api[TList, TCreate, TUpdate, TResp any] struct {
	store     Store[TList, TCreate, TUpdate, TResp]
	binder    binding.Binder
	responder response.Responder
}

func (a *api[TList, TCreate, TUpdate, TResp]) List(w http.ResponseWriter, r *http.Request) {
	var d TList
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

func (a *api[TList, TCreate, TUpdate, TResp]) Get(w http.ResponseWriter, r *http.Request) {
	id := a.binder.URLParam(r, "id")
	res, err := a.store.GetByID(r.Context(), id)
	if err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	a.responder.Respond(w, r, http.StatusOK, res)
}

func (a *api[TList, TCreate, TUpdate, TResp]) Create(w http.ResponseWriter, r *http.Request) {
	var d TCreate
	if err := a.binder.Bind(r, &d); err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	res, err := a.store.Insert(r.Context(), d)
	if err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	a.responder.Respond(w, r, http.StatusCreated, res)
}

func (a *api[TList, TCreate, TUpdate, TResp]) Update(w http.ResponseWriter, r *http.Request) {
	var d TUpdate
	if err := a.binder.Bind(r, &d); err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	id := a.binder.URLParam(r, "id")
	res, err := a.store.Update(r.Context(), id, d)
	if err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	a.responder.Respond(w, r, http.StatusOK, res)
}

func (a *api[TList, TCreate, TUpdate, TResp]) Delete(w http.ResponseWriter, r *http.Request) {
	userID, err := authctx.GetUserKey(r.Context())
	if err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	id := a.binder.URLParam(r, "id")
	err = a.store.Delete(r.Context(), id, userID)
	if err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	a.responder.Respond(w, r, http.StatusNoContent, nil)
}
