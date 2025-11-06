package maintenance

import (
	"errors"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"github.com/gin-gonic/gin"
)

type API interface {
	Maintenance(c *gin.Context)
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

// Maintenance
// @Param body body Request true "body"
// @Success 204
func (a *api) Maintenance(c *gin.Context) {
	var err error

	r := Request{}
	if err = validation.Bind(c, &r); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	userID := c.MustGet(authctx.UserKey).(string)

	// can be sure that enabled is not nil after ShouldBindJSON, because of binding=required
	if *r.Enabled {
		if r.Message == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationError("message", "message is required").ValidationErrorResponse())
			return
		}

		err = a.store.Enable(c, r.Message, r.Color, userID)
	} else {
		err = a.store.Disable(c, userID)
	}

	if err != nil {
		if errors.Is(err, ErrEnabled) || errors.Is(err, ErrDisabled) {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
			return
		}

		a.errorResponder.Respond(c, err)

		return
	}

	c.Status(http.StatusNoContent)
}
