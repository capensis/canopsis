package viewgroup

import (
	"context"
	"errors"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/middleware"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"github.com/gin-gonic/gin"
	"net/http"
)

type api struct {
	store        Store
	actionLogger logger.ActionLogger
}

func NewApi(
	store Store,
	actionLogger logger.ActionLogger,
) common.CrudAPI {
	return &api{
		store:        store,
		actionLogger: actionLogger,
	}
}

// Find all view groups
// @Summary Find view groups
// @Description Get paginated list of view groups
// @Tags viewgroups
// @ID viewgroups-find-all
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param request query ListRequest true "request"
// @Success 200 {object} common.PaginatedListResponse{data=[]ViewGroup}
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /view-groups [get]
func (a *api) List(c *gin.Context) {
	var query ListRequest
	query.Query = pagination.GetDefaultQuery()

	if err := c.ShouldBind(&query); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, query))
		return
	}

	var authorizedIds []string
	ids, ok := c.Get(middleware.AuthorizedIds)
	if ok {
		authorizedIds = ids.([]string)
	}

	viewgroups, err := a.store.Find(c.Request.Context(), query, authorizedIds)
	if err != nil {
		panic(err)
	}

	res, err := common.NewPaginatedResponse(query.Query, viewgroups)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, res)
}

// Get view group by id
// @Summary Get view group by id
// @Description Get view group by id
// @Tags viewgroups
// @ID viewgroups-get-by-id
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "viewgroup id"
// @Success 200 {object} ViewGroup
// @Failure 404 {object} common.ErrorResponse
// @Router /view-groups/{id} [get]
func (a *api) Get(c *gin.Context) {
	viewgroup, err := a.store.GetOneBy(c.Request.Context(), c.Param("id"))
	if err != nil {
		panic(err)
	}
	if viewgroup == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, viewgroup)
}

// Create view group
// @Summary Create view group
// @Description Create view group
// @Tags viewgroups
// @ID viewgroups-create
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param body body EditRequest true "body"
// @Success 201 {object} ViewGroup
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /view-groups [post]
func (a *api) Create(c *gin.Context) {
	request := EditRequest{}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	group, err := a.store.Insert(c.Request.Context(), request)
	if err != nil {
		panic(err)
	}

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    logger.ActionCreate,
		ValueType: logger.ValueTypeViewGroup,
		ValueID:   group.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusCreated, group)
}

// Update view group by id
// @Summary Update view group by id
// @Description Update view group by id
// @Tags viewgroups
// @ID viewgroups-update-by-id
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "view group id"
// @Param body body EditRequest true "body"
// @Success 200 {object} ViewGroup
// @Failure 400 {object} common.ValidationErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /view-groups/{id} [put]
func (a *api) Update(c *gin.Context) {
	request := EditRequest{
		ID: c.Param("id"),
	}

	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	group, err := a.store.Update(c.Request.Context(), request)
	if err != nil {
		panic(err)
	}

	if group == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    logger.ActionUpdate,
		ValueType: logger.ValueTypeViewGroup,
		ValueID:   group.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusOK, group)
}

// Delete view group by id
// @Summary Delete view group by id
// @Description Delete view group by id
// @Tags viewgroups
// @ID viewgroups-delete-by-id
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "view group id"
// @Success 204
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /view-groups/{id} [delete]
func (a *api) Delete(c *gin.Context) {
	id := c.Param("id")
	ok, err := a.store.Delete(c.Request.Context(), id)

	if err != nil {
		if errors.Is(err, ErrLinkedToView) {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
			return
		}
		panic(err)
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    logger.ActionDelete,
		ValueType: logger.ValueTypeViewGroup,
		ValueID:   id,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.Status(http.StatusNoContent)
}
