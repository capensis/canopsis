package datastorage

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"github.com/gin-gonic/gin"
)

type API interface {
	Get(c *gin.Context)

	Update(c *gin.Context)
}

func NewApi(store Store, errorResponder httperror.Responder) API {
	return &api{
		store:          store,
		errorResponder: errorResponder,
	}
}

type api struct {
	store          Store
	errorResponder httperror.Responder
}

// Get
// @Success 200 {object} DataStorage
func (a *api) Get(c *gin.Context) {
	data, err := a.store.Get(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	c.JSON(http.StatusOK, data)
}

// Update
// @Param body body datastorage.Config true "body"
// @Success 200 {object} DataStorage
func (a *api) Update(c *gin.Context) {
	conf := datastorage.Config{}
	if err := validation.Bind(c, &conf); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	data, err := a.store.Update(c, conf)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	c.JSON(http.StatusOK, data)
}
