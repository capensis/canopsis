package heartbeat

import (
	"git.canopsis.net/canopsis/go-engines/lib/api/common"
	"git.canopsis.net/canopsis/go-engines/lib/api/logger"
	"git.canopsis.net/canopsis/go-engines/lib/api/pagination"
	"github.com/gin-gonic/gin"
	"net/http"
)

type api struct {
	store        Store
	transformer  ModelTransformer
	actionLogger logger.ActionLogger
}

func NewApi(
	store 		 Store,
	transformer  ModelTransformer,
	actionLogger logger.ActionLogger,
) common.BulkCrudAPI {
	return &api{
		store:        store,
		transformer:  transformer,
		actionLogger: actionLogger,
	}
}

// Find all heartbeats
// @Summary Find heartbeats
// @Description Get paginated list of heartbeats
// @Tags heartbeats
// @ID heartbeats-find-all
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param page query integer true "current page"
// @Param limit query integer true "items per page"
// @Param search query string false "search query"
// @Param sort query string false "sort query"
// @Param sort_by query string false "sort query"
// @Success 200 {object} common.PaginatedListResponse{data=[]Heartbeat}
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /heartbeats [get]
func (a *api) List(c *gin.Context) {
	var query FilteredQuery
	query.Query = pagination.GetDefaultQuery()

	if err := c.ShouldBind(&query); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, query))
		return
	}

	heartbeats, err := a.store.Find(query)
	if err != nil {
		panic(err)
	}

	res, err := common.NewPaginatedResponse(query.Query, heartbeats)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, res)
}

// Get heartbeat by id
// @Summary Get heartbeat by id
// @Description Get heartbeat by id
// @Tags heartbeats
// @ID heartbeats-get-by-id
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "heartbeat id"
// @Success 200 {object} Heartbeat
// @Failure 404 {object} common.ErrorResponse
// @Router /heartbeats/{id} [get]
func (a *api) Get(c *gin.Context) {
	heartbeat, err := a.store.GetOneBy(c.Param("id"))
	if err != nil {
		panic(err)
	}
	if heartbeat == nil {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, heartbeat)
}

// Create heartbeat
// @Summary Create heartbeat
// @Description Create heartbeat
// @Tags heartbeats
// @ID heartbeats-create
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param body body CreateRequest true "body"
// @Success 201 {object} Heartbeat
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /heartbeats [post]
func (a *api) Create(c *gin.Context) {
	var request CreateRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	heartbeat := a.transformer.TransformCreateRequestToModel(request)
	err := a.store.Insert([]*Heartbeat{heartbeat})
	if err != nil {
		panic(err)
	}

	err = a.actionLogger.Action(c, logger.LogEntry{
		Action:    logger.ActionCreate,
		ValueType: logger.ValueTypeHeartbeat,
		ValueID:   heartbeat.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusCreated, heartbeat)
}

// Update heartbeat by id
// @Summary Update heartbeat by id
// @Description Update heartbeat by id
// @Tags heartbeats
// @ID heartbeats-update-by-id
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "heartbeat id"
// @Param body body UpdateRequest true "body"
// @Success 200 {object} Heartbeat
// @Failure 400 {object} common.ValidationErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /heartbeats/{id} [put]
func (a *api) Update(c *gin.Context) {
	request := UpdateRequest{
		ID: c.Param("id"),
	}

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	heartbeat := a.transformer.TransformUpdateRequestToModel(request)
	err := a.store.Update([]*Heartbeat{heartbeat})
	if err != nil {
		if _, ok := err.(NotFoundError); ok {
			c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
			return
		}

		panic(err)
	}

	err = a.actionLogger.Action(c, logger.LogEntry{
		Action:    logger.ActionUpdate,
		ValueType: logger.ValueTypeHeartbeat,
		ValueID:   heartbeat.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusOK, heartbeat)
}

// Delete heartbeat by id
// @Summary Delete heartbeat by id
// @Description Delete heartbeat by id
// @Tags heartbeats
// @ID heartbeats-delete-by-id
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "heartbeat id"
// @Success 204
// @Failure 404 {object} common.ErrorResponse
// @Router /heartbeats/{id} [delete]
func (a *api) Delete(c *gin.Context) {
	id := c.Param("id")
	err := a.store.Delete([]string{id})

	if err != nil {
		if _, ok := err.(NotFoundError); ok {
			c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
			return
		}

		panic(err)
	}

	err = a.actionLogger.Action(c, logger.LogEntry{
		Action:    logger.ActionDelete,
		ValueType: logger.ValueTypeHeartbeat,
		ValueID:   c.Param("id"),
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.Status(http.StatusNoContent)
}

// Bulk create heartbeats
// @Summary Bulk create heartbeats
// @Description Bulk create heartbeats
// @Tags heartbeats
// @ID heartbeats-bulk-create
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param body body []CreateRequest true "body"
// @Success 201 {array} Heartbeat
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /bulk/heartbeats [post]
func (a *api) BulkCreate(c *gin.Context) {
	var request BulkCreateRequest

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	heartbeats := a.transformer.TransformBulkCreateRequestToModels(request)
	err := a.store.Insert(heartbeats)
	if err != nil {
		panic(err)
	}

	for _, v := range heartbeats {
		err = a.actionLogger.Action(c, logger.LogEntry{
			Action:    logger.ActionCreate,
			ValueType: logger.ValueTypeHeartbeat,
			ValueID:   v.ID,
		})
		if err != nil {
			a.actionLogger.Err(err, "failed to log action")
		}
	}

	c.JSON(http.StatusCreated, heartbeats)
}

// Bulk update heartbeats by id
// @Summary Bulk update heartbeats by id
// @Description Bulk update heartbeats by id
// @Tags heartbeats
// @ID heartbeats-bulk-update-by-id
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param body body []BulkUpdateRequestItem true "body"
// @Success 200 {array} Heartbeat
// @Failure 400 {object} common.ValidationErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /bulk/heartbeats [put]
func (a *api) BulkUpdate(c *gin.Context) {
	request := BulkUpdateRequest{}

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	heartbeats := a.transformer.TransformBulkUpdateRequestToModels(request)
	err := a.store.Update(heartbeats)
	if err != nil {
		if _, ok := err.(NotFoundError); ok {
			c.AbortWithStatusJSON(http.StatusNotFound, common.ErrorResponse{Error: err.Error()})
			return
		}

		panic(err)
	}

	for _, v := range heartbeats {
		err = a.actionLogger.Action(c, logger.LogEntry{
			Action:    logger.ActionUpdate,
			ValueType: logger.ValueTypeHeartbeat,
			ValueID:   v.ID,
		})
		if err != nil {
			a.actionLogger.Err(err, "failed to log action")
		}
	}

	c.JSON(http.StatusOK, heartbeats)
}

// Bulk delete heartbeats by id
// @Summary Bulk delete heartbeats by id
// @Description Bulk delete heartbeats by id
// @Tags heartbeats
// @ID heartbeats-bulk-delete-by-id
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param request query BulkDeleteRequest true "request"
// @Success 204
// @Failure 404 {object} common.ErrorResponse
// @Router /bulk/heartbeats [delete]
func (a *api) BulkDelete(c *gin.Context) {
	request := BulkDeleteRequest{}

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	err := a.store.Delete(request.IDs)
	if err != nil {
		if _, ok := err.(NotFoundError); ok {
			c.AbortWithStatusJSON(http.StatusNotFound, common.ErrorResponse{Error: err.Error()})
			return
		}

		panic(err)
	}

	for _, v := range request.IDs {
		err = a.actionLogger.Action(c, logger.LogEntry{
			Action:    logger.ActionDelete,
			ValueType: logger.ValueTypeHeartbeat,
			ValueID:   v,
		})
		if err != nil {
			a.actionLogger.Err(err, "failed to log action")
		}
	}

	c.Status(http.StatusNoContent)
}
