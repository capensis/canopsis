package entityinfodictionary

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type API interface {
	ListKeys(c *gin.Context)
	ListValues(c *gin.Context)
}

type api struct {
	store          Store
	errorResponder httperror.Responder
	logger         zerolog.Logger
}

func NewApi(
	store Store,
	errorResponder httperror.Responder,
	logger zerolog.Logger,
) API {
	return &api{
		store:          store,
		errorResponder: errorResponder,
		logger:         logger,
	}
}

// ListKeys
// List info dictionary keys
// @Success 200 {object} pagination.ListResponse{data=[]Result}
func (a *api) ListKeys(c *gin.Context) {
	var request ListKeysRequest
	request.Query = pagination.GetDefaultQuery()
	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	findResult, err := a.store.FindKeys(c, request)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	res := pagination.NewResponse(request.Query, findResult)
	c.JSON(http.StatusOK, res)
}

// ListValues
// List info dictionary values
// @Success 200 {object} pagination.ListResponse{data=[]Result}
func (a *api) ListValues(c *gin.Context) {
	var request ListValuesRequest
	request.Query = pagination.GetDefaultQuery()
	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	findResult, err := a.store.FindValues(c, request)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	res := pagination.NewResponse(request.Query, findResult)
	c.JSON(http.StatusOK, res)
}
