package serviceweather

import (
	"errors"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"github.com/gin-gonic/gin"
)

func NewApi(store Store, errorResponder httperror.Responder) API {
	return &api{
		store:          store,
		errorResponder: errorResponder,
	}
}

type API interface {
	List(c *gin.Context)
	EntityList(c *gin.Context)
}

type api struct {
	store          Store
	errorResponder httperror.Responder
}

// List
// @Success 200 {object} pagination.ListResponse{data=[]Service}
func (a *api) List(c *gin.Context) {
	var query ListRequest
	query.Query = pagination.GetDefaultQuery()

	if err := validation.Bind(c, &query); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	aggregationResult, err := a.store.Find(c, query)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}

		a.errorResponder.Respond(c, err)

		return
	}

	res := pagination.NewResponse(query.Query, aggregationResult)
	c.JSON(http.StatusOK, res)
}

// EntityList
// @Success 200 {object} pagination.ListResponse{data=[]Entity}
func (a *api) EntityList(c *gin.Context) {
	var query EntitiesListRequest
	query.Query = pagination.GetDefaultQuery()

	if err := validation.Bind(c, &query); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	id := c.Param("id")
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	aggregationResult, err := a.store.FindEntities(c, id, query, userID)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if aggregationResult == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	res := pagination.NewResponse(query.Query, aggregationResult)
	c.JSON(http.StatusOK, res)
}
