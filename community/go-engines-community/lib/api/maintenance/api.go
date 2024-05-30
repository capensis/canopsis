package maintenance

import (
	"errors"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"github.com/gin-gonic/gin"
)

type API interface {
	Maintenance(c *gin.Context)
}

func NewApi(store Store) API {
	return &api{store: store}
}

type api struct {
	store Store
}

// Maintenance
// @Param body body Request true "body"
// @Success 204
func (a *api) Maintenance(c *gin.Context) {
	var err error

	r := Request{}
	if err = c.ShouldBindJSON(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))
		return
	}

	userId := c.MustGet(auth.UserKey).(string)

	// can be sure that enabled is not nil after ShouldBindJSON, because of binding=required
	if *r.Enabled {
		if r.Message == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationError("message", "message is required").ValidationErrorResponse())
			return
		}

		err = a.store.Enable(c, r.Message, r.Color, userId)
	} else {
		err = a.store.Disable(c, userId)
	}

	if err != nil {
		if errors.Is(err, ErrEnabled) || errors.Is(err, ErrDisabled) {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
			return
		}

		panic(err)
	}

	c.Status(http.StatusNoContent)
}
