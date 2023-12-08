package statesettings

import (
	"context"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"github.com/gin-gonic/gin"
)

type API interface {
	Update(c *gin.Context)
	List(c *gin.Context)
}

type api struct {
	store        Store
	actionLogger logger.ActionLogger
}

func NewApi(
	store Store,
	actionLogger logger.ActionLogger,
) API {
	return &api{
		store:        store,
		actionLogger: actionLogger,
	}
}

// List
// @Success 200 {object} common.PaginatedListResponse{data=[]StateSetting}
func (a *api) List(c *gin.Context) {
	var query pagination.Query

	if err := pagination.BindQuery(c, &query); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, query))
		return
	}

	aggregationResult, err := a.store.Find(c.Request.Context(), query)
	if err != nil {
		panic(err)
	}

	res, err := common.NewPaginatedResponse(query, &aggregationResult)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, res)
}

// Update
// @Param body body StateSettingRequest true "body"
// @Success 200 {object} StateSetting
func (a *api) Update(c *gin.Context) {
	request := StateSettingRequest{
		ID: c.Param("id"),
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	stateSetting, err := a.store.Update(c.Request.Context(), request)
	if err != nil {
		panic(err)
	}

	if stateSetting == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    logger.ActionUpdate,
		ValueType: logger.ValueTypeStateSetting,
		ValueID:   stateSetting.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusOK, stateSetting)
}
