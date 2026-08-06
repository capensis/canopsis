package bulk

import (
	"context"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/binding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/crud"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/response"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
)

type API interface {
	crud.API
	BulkCreate(w http.ResponseWriter, r *http.Request)
	BulkUpdate(w http.ResponseWriter, r *http.Request)
	BulkDelete(w http.ResponseWriter, r *http.Request)
}

type Store[TList, TCreate, TUpdate, TBulkUpdate, TBulkDelete, TResp any] interface {
	crud.Store[TList, TCreate, TUpdate, TResp]
	BulkInsert(ctx context.Context, r TCreate) (string, error)
	BulkUpdate(ctx context.Context, r TBulkUpdate) (string, error)
	BulkDelete(ctx context.Context, r TBulkDelete) (string, error)
}

func NewAPI[TList, TCreate, TUpdate, TBulkUpdate, TBulkDelete, TResp any](
	store Store[TList, TCreate, TUpdate, TBulkUpdate, TBulkDelete, TResp],
	binder binding.Binder,
	responder response.Responder,
	decoder encoding.Decoder,
) API {
	return &api[TList, TCreate, TUpdate, TBulkUpdate, TBulkDelete, TResp]{
		API:       crud.NewAPI(store, binder, responder),
		store:     store,
		binder:    binder,
		responder: responder,
		decoder:   decoder,
	}
}

type api[TList, TCreate, TUpdate, TBulkUpdate, TBulkDelete, TResp any] struct {
	crud.API
	store     Store[TList, TCreate, TUpdate, TBulkUpdate, TBulkDelete, TResp]
	binder    binding.Binder
	responder response.Responder
	decoder   encoding.Decoder
}

func (a *api[TList, TCreate, TUpdate, TBulkUpdate, TBulkDelete, TResp]) BulkCreate(w http.ResponseWriter, r *http.Request) {
	Handler(w, r, func(d TCreate) (string, error) {
		return a.store.BulkInsert(r.Context(), d)
	}, a.responder, a.binder, a.decoder)
}

func (a *api[TList, TCreate, TUpdate, TBulkUpdate, TBulkDelete, TResp]) BulkUpdate(w http.ResponseWriter, r *http.Request) {
	Handler(w, r, func(d TBulkUpdate) (string, error) {
		return a.store.BulkUpdate(r.Context(), d)
	}, a.responder, a.binder, a.decoder)
}

func (a *api[TList, TCreate, TUpdate, TBulkUpdate, TBulkDelete, TResp]) BulkDelete(w http.ResponseWriter, r *http.Request) {
	Handler(w, r, func(d TBulkDelete) (string, error) {
		return a.store.BulkDelete(r.Context(), d)
	}, a.responder, a.binder, a.decoder)
}
