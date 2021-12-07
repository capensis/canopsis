package user

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"github.com/gin-gonic/gin"
	"net/http"
)

type api struct {
	store        Store
	actionLogger logger.ActionLogger

	metricMetaUpdater metrics.MetaUpdater
}

func NewApi(
	store Store,
	actionLogger logger.ActionLogger,
	metricMetaUpdater metrics.MetaUpdater,
) common.BulkCrudAPI {
	return &api{
		store:        store,
		actionLogger: actionLogger,

		metricMetaUpdater: metricMetaUpdater,
	}
}

// Find all users
// @Summary Find users
// @Description Get paginated list of users
// @Tags users
// @ID users-find-all
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param page query integer true "current page"
// @Param limit query integer true "items per page"
// @Param search query string false "search query"
// @Param sort query string false "sort query"
// @param permission query string false "role permission"
// @Param sort_by query string false "sort query"
// @Success 200 {object} common.PaginatedListResponse{data=[]User}
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /users [get]
func (a *api) List(c *gin.Context) {
	var query ListRequest
	query.Query = pagination.GetDefaultQuery()

	if err := c.ShouldBind(&query); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, query))
		return
	}

	users, err := a.store.Find(c.Request.Context(), query)
	if err != nil {
		panic(err)
	}

	res, err := common.NewPaginatedResponse(query.Query, users)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, res)
}

// Get user by id
// @Summary Get user by id
// @Description Get user by id
// @Tags users
// @ID users-get-by-id
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "user id"
// @Success 200 {object} User
// @Failure 404 {object} common.ErrorResponse
// @Router /users/{id} [get]
func (a *api) Get(c *gin.Context) {
	user, err := a.store.GetOneBy(c.Request.Context(), c.Param("id"))
	if err != nil {
		panic(err)
	}
	if user == nil {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, user)
}

// Create user
// @Summary Create user
// @Description Create user
// @Tags users
// @ID users-create
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param body body EditRequest true "body"
// @Success 201 {object} User
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /users [post]
func (a *api) Create(c *gin.Context) {
	var request Request
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	user, err := a.store.Insert(c.Request.Context(), request)
	if err != nil {
		panic(err)
	}

	err = a.actionLogger.Action(c, logger.LogEntry{
		Action:    logger.ActionCreate,
		ValueType: logger.ValueTypeUser,
		ValueID:   user.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	a.metricMetaUpdater.UpdateById(c.Request.Context(), user.ID)

	c.JSON(http.StatusCreated, user)
}

// Update user by id
// @Summary Update user by id
// @Description Update user by id
// @Tags users
// @ID users-update-by-id
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "user id"
// @Param body body EditRequest true "body"
// @Success 200 {object} User
// @Failure 400 {object} common.ValidationErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /users/{id} [put]
func (a *api) Update(c *gin.Context) {
	request := Request{
		ID: c.Param("id"),
	}

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	user, err := a.store.Update(c.Request.Context(), request)
	if err != nil {
		panic(err)
	}

	if user == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	err = a.actionLogger.Action(c, logger.LogEntry{
		Action:    logger.ActionUpdate,
		ValueType: logger.ValueTypeUser,
		ValueID:   user.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	a.metricMetaUpdater.UpdateById(c.Request.Context(), user.ID)

	c.JSON(http.StatusOK, user)
}

// Delete user by id
// @Summary Delete user by id
// @Description Delete user by id
// @Tags users
// @ID users-delete-by-id
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "user id"
// @Success 204
// @Failure 404 {object} common.ErrorResponse
// @Router /users/{id} [delete]
func (a *api) Delete(c *gin.Context) {
	id := c.Param("id")
	ok, err := a.store.Delete(c.Request.Context(), id)

	if err != nil {
		panic(err)
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	err = a.actionLogger.Action(c, logger.LogEntry{
		Action:    logger.ActionDelete,
		ValueType: logger.ValueTypeUser,
		ValueID:   id,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	a.metricMetaUpdater.DeleteById(c.Request.Context(), id)

	c.Status(http.StatusNoContent)
}

func (a *api) BulkCreate(c *gin.Context) {
	var request BulkCreateRequest

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	ctx := c.Request.Context()

	err := a.store.BulkInsert(ctx, request.Items)
	if err != nil {
		panic(err)
	}

	for _, item := range request.Items {
		err = a.actionLogger.Action(c, logger.LogEntry{
			Action:    logger.ActionCreate,
			ValueType: logger.ValueTypeUser,
			ValueID:   item.Name,
		})
		if err != nil {
			a.actionLogger.Err(err, "failed to log action")
		}
	}

	c.Status(http.StatusNoContent)
}

func (a *api) BulkUpdate(c *gin.Context) {
	var request BulkUpdateRequest

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	ctx := c.Request.Context()

	err := a.store.BulkUpdate(ctx, request.Items)
	if err != nil {
		panic(err)
	}

	for _, v := range request.Items {
		err = a.actionLogger.Action(c, logger.LogEntry{
			Action:    logger.ActionUpdate,
			ValueType: logger.ValueTypeUser,
			ValueID:   v.ID,
		})
		if err != nil {
			a.actionLogger.Err(err, "failed to log action")
		}
	}

	c.Status(http.StatusNoContent)
}

func (a *api) BulkDelete(c *gin.Context) {
	request := BulkDeleteRequest{}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	err := a.store.BulkDelete(c.Request.Context(), request.IDs)
	if err != nil {
		panic(err)
	}

	for _, id := range request.IDs {
		err = a.actionLogger.Action(c, logger.LogEntry{
			Action:    logger.ActionDelete,
			ValueType: logger.ValueTypeUser,
			ValueID:   id,
		})
		if err != nil {
			a.actionLogger.Err(err, "failed to log action")
		}
	}

	c.Status(http.StatusNoContent)
}
