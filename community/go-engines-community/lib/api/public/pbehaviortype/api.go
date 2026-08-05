package pbehaviortype

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/binding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/crud"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/response"
)

func NewAPI(store Store, binder binding.Binder, responder response.Responder) crud.API {
	return crud.NewAPI[ListRequest, CreateRequest, UpdateRequest, Response](store, binder, responder)
}
