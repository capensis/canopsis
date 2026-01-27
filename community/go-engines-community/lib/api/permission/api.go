package permission

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"github.com/gin-gonic/gin"
)

type API interface {
	List(c *gin.Context)
}

type api struct {
	store          Store
	errorResponder httperror.Responder
}

func NewApi(store Store, errorResponder httperror.Responder) API {
	return &api{
		store:          store,
		errorResponder: errorResponder,
	}
}

// List
// @Success 200 {object} pagination.ListResponse{data=[]Permission}
func (a *api) List(c *gin.Context) {
	var query ListRequest
	query.Query = pagination.GetDefaultQuery()

	if err := validation.Bind(c, &query); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	permissions, err := a.store.Find(c, query)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	res := pagination.NewResponse(query.Query, permissions)
	c.JSON(http.StatusOK, res)
}
