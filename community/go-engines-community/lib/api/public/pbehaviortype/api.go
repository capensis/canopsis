package pbehaviortype

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/binding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/crud"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/response"
)

// NewAPI
//
// @Param body body CreateRequest true "body"
// @Param body body UpdateRequest true "body"
// @Success 200 {object} Response
func NewAPI(store Store, binder binding.Binder, responder response.Responder) crud.API {
	return crud.NewAPI[ListRequest, CreateRequest, UpdateRequest, Response](store, binder, responder)
}
