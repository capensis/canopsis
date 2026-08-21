package entity

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/binding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/response"
)

type API interface {
	List(w http.ResponseWriter, r *http.Request)
	Search(w http.ResponseWriter, r *http.Request)
}

func NewAPI(
	store Store,
	binder binding.Binder,
	responder response.Responder,
) API {
	return &api{
		store:     store,
		binder:    binder,
		responder: responder,
	}
}

type api struct {
	store     Store
	binder    binding.Binder
	responder response.Responder
}

// List
//
// @Success 200 {object} Response
func (a *api) List(w http.ResponseWriter, r *http.Request) {
	var d PaginatedListRequest
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

// Search
//
// @Param body body SearchRequest true "body"
// @Success 200 {object} Response
func (a *api) Search(w http.ResponseWriter, r *http.Request) {
	var d PaginatedSearchRequest
	d.Query.SetDefaults()
	if err := a.binder.Bind(r, &d); err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	res, err := a.store.Search(r.Context(), d)
	if err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	a.responder.Respond(w, r, http.StatusOK, res)
}
