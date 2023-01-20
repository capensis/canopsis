package account

import (
	"context"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"github.com/gin-gonic/gin"
)

type API interface {
	Me(c *gin.Context)
	Update(c *gin.Context)
}

func NewApi(store Store, actionLogger logger.ActionLogger) API {
	return &api{
		store:        store,
		actionLogger: actionLogger,
	}
}

type api struct {
	store        Store
	actionLogger logger.ActionLogger
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

	err = a.actionLogger.Action(context.Background(), userID, logger.LogEntry{
		Action:    logger.ActionUpdate,
		ValueType: logger.ValueTypeUser,
		ValueID:   userID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusOK, user)
}
