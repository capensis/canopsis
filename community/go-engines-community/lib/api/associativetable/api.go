package associativetable

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"github.com/gin-gonic/gin"
)

// API represents associative table modification actions.
// Associative table API is not REST since it doesn't return error if model doesn't exist :
// - Update - creates model if not exist or updates model if exist
// - Get - returns empty model if not exist or returns model if exist
// - Delete - returns does nothing if not exist or deletes model if exist
type API interface {
	Update(c *gin.Context)
	Get(c *gin.Context)
	Delete(c *gin.Context)
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

// Update
// @Param body body AssociativeTable true "body"
// @Success 200 {object} AssociativeTable
func (a *api) Update(c *gin.Context) {
	request := AssociativeTable{}
	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	err := a.store.Update(c, &request)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	c.JSON(http.StatusOK, request)
}

// Get
// @Success 200 {object} AssociativeTable
func (a *api) Get(c *gin.Context) {
	request := GetRequest{}
	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	at, err := a.store.GetByName(c, request.Name)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if at == nil {
		at = &AssociativeTable{Name: request.Name}
	}

	c.JSON(http.StatusOK, at)
}

func (a *api) Delete(c *gin.Context) {
	request := GetRequest{}
	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	err := a.store.Delete(c, request.Name)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	c.JSON(http.StatusNoContent, nil)
}
