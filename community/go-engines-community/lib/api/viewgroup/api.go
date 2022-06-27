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
) common.BulkCrudAPI {
	return &api{
		store:        store,
		actionLogger: actionLogger,
	}
}

// List
// @Success 200 {object} common.PaginatedListResponse{data=[]ViewGroup}
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

// Get
// @Success 200 {object} ViewGroup
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

// Create
// @Param body body EditRequest true "body"
// @Success 201 {object} ViewGroup
func (a *api) Create(c *gin.Context) {
	request := EditRequest{}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	groups, err := a.store.Insert(c.Request.Context(), []EditRequest{request})
	if err != nil {
		panic(err)
	}

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    logger.ActionCreate,
		ValueType: logger.ValueTypeViewGroup,
		ValueID:   groups[0].ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusCreated, groups[0])
}

// Update
// @Param body body EditRequest true "body"
// @Success 200 {object} ViewGroup
func (a *api) Update(c *gin.Context) {
	request := EditRequest{
		ID: c.Param("id"),
	}

	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	groups, err := a.store.Update(c.Request.Context(), []BulkUpdateRequestItem{{
		ID:              request.ID,
		BaseEditRequest: request.BaseEditRequest,
	}})
	if err != nil {
		panic(err)
	}

	if len(groups) == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    logger.ActionUpdate,
		ValueType: logger.ValueTypeViewGroup,
		ValueID:   groups[0].ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusOK, groups[0])
}

func (a *api) Delete(c *gin.Context) {
	id := c.Param("id")
	ok, err := a.store.Delete(c.Request.Context(), []string{id})

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

func (a *api) BulkCreate(c *gin.Context) {
	var request BulkCreateRequest
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	groups, err := a.store.Insert(c.Request.Context(), request.Items)
	if err != nil {
		panic(err)
	}

	for _, group := range groups {
		err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
			Action:    logger.ActionCreate,
			ValueType: logger.ValueTypeViewGroup,
			ValueID:   group.ID,
		})
		if err != nil {
			a.actionLogger.Err(err, "failed to log action")
		}
	}

	c.JSON(http.StatusCreated, groups)
}

func (a *api) BulkUpdate(c *gin.Context) {
	request := BulkUpdateRequest{}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	groups, err := a.store.Update(c.Request.Context(), request.Items)
	if err != nil {
		panic(err)
	}

	if len(groups) == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	for _, group := range groups {
		err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
			Action:    logger.ActionUpdate,
			ValueType: logger.ValueTypeViewGroup,
			ValueID:   group.ID,
		})
		if err != nil {
			a.actionLogger.Err(err, "failed to log action")
		}
	}

	c.JSON(http.StatusOK, groups)
}

func (a *api) BulkDelete(c *gin.Context) {
	request := BulkDeleteRequest{}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	ok, err := a.store.Delete(c.Request.Context(), request.IDs)
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

	for _, id := range request.IDs {
		err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
			Action:    logger.ActionDelete,
			ValueType: logger.ValueTypeViewGroup,
			ValueID:   id,
		})
		if err != nil {
			a.actionLogger.Err(err, "failed to log action")
		}
	}

	c.Status(http.StatusNoContent)
}
