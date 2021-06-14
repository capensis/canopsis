package account

import (
	"git.canopsis.net/canopsis/go-engines/lib/api/auth"
	"git.canopsis.net/canopsis/go-engines/lib/api/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

type API interface {
	Me() gin.HandlerFunc
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
func (a *api) Me() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.Get(auth.UserKey)

		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.UnauthorizedResponse)
			return
		}

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
}
