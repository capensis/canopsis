package pattern

import (
	"context"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	apisecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"github.com/gin-gonic/gin"
)

type api struct {
	store        Store
	enforcer     security.Enforcer
	actionLogger logger.ActionLogger
}

func NewApi(
	store Store,
	enforcer security.Enforcer,
	actionLogger logger.ActionLogger,
) common.CrudAPI {
	return &api{
		store:        store,
		enforcer:     enforcer,
		actionLogger: actionLogger,
	}
}

// Create creates new pattern.
// @Summary Create saved pattern
// @Description Create saved pattern
// @Tags saved-patterns
// @ID saved-patterns-create
// @Accept json
// @Produce json
// @Security JWTAuth
// @Security BasicAuth
// @Param body body EditRequest true "body"
// @Success 201 {object} Response
// @Failure 400 {object} common.ErrorResponse
// @Router /patterns [post]
func (a *api) Create(c *gin.Context) {
	userId := c.MustGet(auth.UserKey).(string)
	request := EditRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	if *request.IsCorporate {
		ok, err := a.enforcer.Enforce(userId, apisecurity.PermCorporatePattern, model.PermissionCan)
		if err != nil {
			panic(err)
		}

		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, common.ForbiddenResponse)
			return
		}
	}

	pattern, err := a.store.Insert(c.Request.Context(), request)
	if err != nil {
		panic(err)
	}

	err = a.actionLogger.Action(context.Background(), userId, logger.LogEntry{
		Action:    logger.ActionCreate,
		ValueType: logger.ValueTypePattern,
		ValueID:   pattern.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusCreated, pattern)
}

// List finds saved patterns.
// @Summary Find all saved patterns
// @Description Get paginated list of saved patterns
// @Tags saved-patterns
// @ID saved-patterns-find-all
// @Accept json
// @Produce json
// @Security JWTAuth
// @Security BasicAuth
// @Param request query ListRequest true "request"
// @Success 200 {object} common.PaginatedListResponse{data=[]Response}
// @Failure 400 {object} common.ErrorResponse
// @Router /patterns [get]
func (a *api) List(c *gin.Context) {
	var request ListRequest
	request.Query = pagination.GetDefaultQuery()

	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	aggregationResult, err := a.store.Find(c.Request.Context(), request, c.MustGet(auth.UserKey).(string))
	if err != nil {
		panic(err)
	}

	res, err := common.NewPaginatedResponse(request.Query, aggregationResult)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, res)
}

// Get saved pattern by id.
// @Summary Get saved pattern by id
// @Description Get saved pattern by id
// @Tags saved-patterns
// @ID saved-patterns-get-by-id
// @Produce json
// @Security JWTAuth
// @Security BasicAuth
// @Param id path string true "pattern id"
// @Success 200 {object} Response
// @Failure 404 {object} common.ErrorResponse
// @Router /patterns/{id} [get]
func (a *api) Get(c *gin.Context) {
	pattern, err := a.store.GetById(c.Request.Context(), c.Param("id"), c.MustGet(auth.UserKey).(string))
	if err != nil {
		panic(err)
	}

	if pattern == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, pattern)
}

// Update saved pattern by id.
// @Summary Update saved pattern by id
// @Description Update saved pattern by id
// @Tags saved-patterns
// @ID saved-patterns-update-by-id
// @Accept json
// @Produce json
// @Security JWTAuth
// @Security BasicAuth
// @Param id path string true "pattern id"
// @Param body body EditRequest true "body"
// @Success 200 {object} Response
// @Failure 400 {object} common.ValidationErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /patterns/{id} [put]
func (a *api) Update(c *gin.Context) {
	userId := c.MustGet(auth.UserKey).(string)
	request := EditRequest{
		ID: c.Param("id"),
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	pattern, err := a.store.GetById(c.Request.Context(), request.ID, userId)
	if err != nil {
		panic(err)
	}

	if pattern == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	if pattern.Type != request.Type {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.ValidationErrorResponse{Errors: map[string]string{"type": "Type cannot be changed"}})
		return
	}

	if pattern.IsCorporate != *request.IsCorporate {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.ValidationErrorResponse{Errors: map[string]string{"is_corporate": "IsCorporate cannot be changed"}})
		return
	}

	if pattern.IsCorporate {
		ok, err := a.enforcer.Enforce(userId, apisecurity.PermCorporatePattern, model.PermissionCan)
		if err != nil {
			panic(err)
		}

		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, common.ForbiddenResponse)
			return
		}
	}

	pattern, err = a.store.Update(c.Request.Context(), request)
	if err != nil {
		panic(err)
	}

	if pattern == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	err = a.actionLogger.Action(context.Background(), userId, logger.LogEntry{
		Action:    logger.ActionUpdate,
		ValueType: logger.ValueTypePattern,
		ValueID:   pattern.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusOK, pattern)
}

// Delete saved pattern by id
// @Summary Delete saved pattern by id
// @Description Delete saved pattern by id
// @Tags saved-patterns
// @ID saved-patterns-delete-by-id
// @Security JWTAuth
// @Security BasicAuth
// @Param id path string true "pattern id"
// @Success 204
// @Failure 404 {object} common.ErrorResponse
// @Router /patterns/{id} [delete]
func (a *api) Delete(c *gin.Context) {
	userId := c.MustGet(auth.UserKey).(string)
	id := c.Param("id")
	pattern, err := a.store.GetById(c.Request.Context(), id, userId)
	if err != nil {
		panic(err)
	}

	if pattern == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	if pattern.IsCorporate {
		ok, err := a.enforcer.Enforce(userId, apisecurity.PermCorporatePattern, model.PermissionCan)
		if err != nil {
			panic(err)
		}

		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, common.ForbiddenResponse)
			return
		}
	}

	ok, err := a.store.Delete(c.Request.Context(), id)
	if err != nil {
		panic(err)
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	err = a.actionLogger.Action(context.Background(), userId, logger.LogEntry{
		Action:    logger.ActionDelete,
		ValueType: logger.ValueTypePattern,
		ValueID:   id,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusNoContent, nil)
}
