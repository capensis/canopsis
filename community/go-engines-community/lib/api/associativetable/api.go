package associativetable

import (
	"context"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
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
	actionLogger logger.ActionLogger
	store        Store
}

func NewApi(
	store Store,
	actionLogger logger.ActionLogger,
) API {
	return &api{
		actionLogger: actionLogger,
		store:        store,
	}
}

// Update
// @Param body body AssociativeTable true "body"
// @Success 200 {object} AssociativeTable
func (a api) Update(c *gin.Context) {
	request := AssociativeTable{}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	isNew, err := a.store.Update(c.Request.Context(), &request)
	if err != nil {
		panic(err)
	}

	action := logger.ActionUpdate
	if isNew {
		action = logger.ActionCreate
	}

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    action,
		ValueType: logger.ValueAssociativeTable,
		ValueID:   request.Name,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusOK, request)
}

// Get
// @Success 200 {object} AssociativeTable
func (a api) Get(c *gin.Context) {
	request := GetRequest{}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	at, err := a.store.GetByName(c.Request.Context(), request.Name)
	if err != nil {
		panic(err)
	}

	if at == nil {
		at = &AssociativeTable{Name: request.Name}
	}

	c.JSON(http.StatusOK, at)
}

func (a api) Delete(c *gin.Context) {
	request := GetRequest{}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	ok, err := a.store.Delete(c.Request.Context(), request.Name)
	if err != nil {
		panic(err)
	}

	if ok {
		err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
			Action:    logger.ActionDelete,
			ValueType: logger.ValueAssociativeTable,
			ValueID:   request.Name,
		})
		if err != nil {
			a.actionLogger.Err(err, "failed to log action")
		}
	}

	c.JSON(http.StatusNoContent, nil)
}
