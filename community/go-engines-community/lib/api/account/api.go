package account

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"github.com/gin-gonic/gin"
)

type API interface {
	Me(c *gin.Context)
	Update(c *gin.Context)
}

func NewApi(store Store) API {
	return &api{
		store: store,
	}
}

type api struct {
	store Store
}

// Me
// @Success 200 {object} User
func (a *api) Me(c *gin.Context) {
	userID := c.MustGet(auth.UserKey).(string)

	user, err := a.store.GetOneBy(c, userID)
	if err != nil {
		panic(err)
	}

	if user == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, common.UnauthorizedResponse)
		return
	}

	c.JSON(http.StatusOK, user)
}

// Update
// @Param body body EditRequest true "body"
// @Success 200 {object} User
func (a *api) Update(c *gin.Context) {
	userID := c.MustGet(auth.UserKey).(string)
	request := EditRequest{
		ID: userID,
		// author is needed for action logs, in that case the user modifies himself, so he's the author.
		Author: userID,
	}

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	user, err := a.store.Update(c, request)
	if err != nil {
		panic(err)
	}

	if user == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, common.UnauthorizedResponse)
		return
	}

	c.JSON(http.StatusOK, user)
}
