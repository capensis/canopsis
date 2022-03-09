package entitycategory

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"github.com/gin-gonic/gin"
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

// Find all entity categories
// @Summary Find entity categories
// @Description Get paginated list of entity categories
// @Tags entity-categories
// @ID entity-categories-find-all
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param page query integer true "current page"
// @Param limit query integer true "items per page"
// @Param search query string false "search query"
// @Param sort query string false "sort query"
// @Param sort_by query string false "sort query"
// @Success 200 {object} common.PaginatedListResponse{data=[]Category}
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /entity-categories [get]
func (a *api) List(c *gin.Context) {
	var r ListRequest
	r.Query = pagination.GetDefaultQuery()

	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))
		return
	}

	categories, err := a.store.Find(c.Request.Context(), r)
	if err != nil {
		panic(err)
	}

	res, err := common.NewPaginatedResponse(r.Query, categories)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, res)
}

// Get entity category by id
// @Summary Get entity category by id
// @Description Get entity category by id
// @Tags entity-categories
// @ID entity-categories-get-by-id
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "category id"
// @Success 200 {object} Category
// @Failure 404 {object} common.ErrorResponse
// @Router /entity-categories/{id} [get]
func (a *api) Get(c *gin.Context) {
	category, err := a.store.GetOneBy(c.Request.Context(), c.Param("id"))
	if err != nil {
		panic(err)
	}
	if category == nil {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, category)
}

// Create entity category
// @Summary Create entity category
// @Description Create entity category
// @Tags entity-categories
// @ID entity-categories-create
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param body body EditRequest true "body"
// @Success 201 {object} Category
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /entity-categories [post]
func (a *api) Create(c *gin.Context) {
	var request EditRequest

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	category, err := a.store.Insert(c.Request.Context(), request)
	if err != nil {
		panic(err)
	}

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    logger.ActionCreate,
		ValueType: logger.ValueTypeEntityCategory,
		ValueID:   category.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusCreated, category)
}

// Update entity category by id
// @Summary Update entity category by id
// @Description Update entity category by id
// @Tags entity-categories
// @ID entity-categories-update-by-id
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "category id"
// @Param body body EditRequest true "body"
// @Success 200 {object} Category
// @Failure 400 {object} common.ValidationErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /entity-categories/{id} [put]
func (a *api) Update(c *gin.Context) {
	request := EditRequest{
		ID: c.Param("id"),
	}

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	category, err := a.store.Update(c.Request.Context(), request)
	if err != nil {
		panic(err)
	}

	if category == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    logger.ActionUpdate,
		ValueType: logger.ValueTypeEntityCategory,
		ValueID:   category.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusOK, category)
}

// Delete entity category by id
// @Summary Delete entity category by id
// @Description Delete entity category by id
// @Tags entity-categories
// @ID entity-categories-delete-by-id
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "category id"
// @Success 204
// @Failure 404 {object} common.ErrorResponse
// @Router /entity-categories/{id} [delete]
func (a *api) Delete(c *gin.Context) {
	id := c.Param("id")
	ok, err := a.store.Delete(c.Request.Context(), id)

	if err != nil {
		if err == ErrLinkedCategoryToEntity {
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
		ValueType: logger.ValueTypeEntityCategory,
		ValueID:   id,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.Status(http.StatusNoContent)
}
