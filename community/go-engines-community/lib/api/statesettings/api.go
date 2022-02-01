package statesettings

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"github.com/gin-gonic/gin"
	"net/http"
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

// Find all state settings for a widget
// @Summary Find all state settings for a widget
// @Description Get paginated list of state settings
// @Tags state-settings
// @ID state-settings-list
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param page query integer true "current page"
// @Param limit query integer true "items per page"
// @Success 200 {object} common.PaginatedListResponse{data=[]StateSetting}
// @Router /state-settings [get]
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

// Update state setting by id
// @Summary Update state setting type by id
// @Description Update state setting type by id
// @Tags state-settings
// @ID state-settings-update-by-id
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "type id"
// @Param body body StateSettingRequest true "body"
// @Success 200 {object} StateSetting
// @Failure 400 {object} common.ValidationErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /state-settings/{id} [put]
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
