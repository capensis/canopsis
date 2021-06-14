package account

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"github.com/gin-gonic/gin"
)

type API interface {
	Me(c *gin.Context)
}

func NewApi(store Store) API {
	return &api{
		store: store,
	}
}

type api struct {
	store Store
}

// Get account
// @Summary Get account
// @Description Get account
// @Tags account
// @ID account-get
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Success 200 {object} User
// @Failure 401 {object} common.ErrorResponse
// @Router /account/me [get]
func (a *api) Me(c *gin.Context) {
	userID := c.MustGet(auth.UserKey)

	user, err := a.store.GetOneBy(userID.(string))
	if err != nil {
		panic(err)
	}

	if user == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, common.UnauthorizedResponse)
		return
	}

	c.JSON(http.StatusOK, user)
}
