package broadcastmessage

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"github.com/gin-gonic/gin"
)

type API interface {
	common.CrudAPI
	GetActive(c *gin.Context)
}

type api struct {
	store            Store
	onChangeListener chan<- bool
	actionLogger     logger.ActionLogger
}

// Create broadcast-message
// @Summary Create broadcast-message
// @Description Create broadcast-message
// @Tags broadcast-messages
// @ID broadcast-messages-create
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param body body BroadcastMessage true "body"
// @Success 201 {object} BroadcastMessage
// @Failure 400 {object} common.ErrorResponse
// @Router /broadcast-message [post]
func (a api) Create(c *gin.Context) {
	request := BroadcastMessage{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	err := a.store.Insert(c.Request.Context(), &request)
	if err != nil {
		panic(err)
	}

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    logger.ActionCreate,
		ValueType: logger.ValueTypeBroadcastMessage,
		ValueID:   request.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	a.sendOnChange()

	c.JSON(http.StatusCreated, request)
}

// Find all broadcast-message
// @Summary Find all broadcast-message
// @Description Get paginated list of broadcast-message
// @Tags broadcast-messages
// @ID broadcast-messages-find-all
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param page query integer true "current page"
// @Param limit query integer true "items per page"
// @Param search query string false "search query"
// @Success 200 {object} common.PaginatedListResponse{data=[]BroadcastMessage}
// @Failure 400 {object} common.ErrorResponse
// @Router /broadcast-message [get]
func (a api) List(c *gin.Context) {
	var query FilteredQuery
	query.Query = pagination.GetDefaultQuery()

	if err := c.ShouldBind(&query); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, query))
		return
	}

	aggregationResult, err := a.store.Find(c.Request.Context(), query)
	if err != nil {
		panic(err)
	}

	res, err := common.NewPaginatedResponse(query.Query, aggregationResult)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, res)
}

// Get broadcast-message by id
// @Summary Get broadcast-message by id
// @Description Get broadcast-message by id
// @Tags broadcast-messages
// @ID broadcast-messages-get-by-id
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "broadcast-message id"
// @Success 200 {object} BroadcastMessage
// @Failure 404 {object} common.ErrorResponse
// @Router /broadcast-message/{id} [get]
func (a api) Get(c *gin.Context) {
	bm, err := a.store.GetById(c.Request.Context(), c.Param("id"))
	if err != nil {
		panic(err)
	}

	if bm == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, bm)
}

// Update broadcast-message by id
// @Summary Update broadcast-message by id
// @Description Update broadcast-message by id
// @Tags broadcast-messages
// @ID broadcast-messages-update-by-id
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "broadcast-message id"
// @Param body body Payload true "body"
// @Success 200 {object} BroadcastMessage
// @Failure 400 {object} common.ValidationErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /broadcast-message/{id} [put]
func (a api) Update(c *gin.Context) {
	var request Payload
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))

		return
	}

	var data BroadcastMessage
	data.Payload = request
	data.ID = c.Param("id")
	ok, _ := a.store.Update(c.Request.Context(), &data)

	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	err := a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    logger.ActionUpdate,
		ValueType: logger.ValueTypeBroadcastMessage,
		ValueID:   data.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	a.sendOnChange()

	c.JSON(http.StatusOK, data)
}

// Delete broadcast-message by id
// @Summary Delete broadcast-message by id
// @Description Delete broadcast-message by id
// @Tags broadcast-messages
// @ID broadcast-messages-delete-by-id
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "broadcast-message id"
// @Success 204
// @Failure 404 {object} common.ErrorResponse
// @Router /broadcast-message/{id} [delete]
func (a api) Delete(c *gin.Context) {
	ok, err := a.store.Delete(c.Request.Context(), c.Param("id"))
	if err != nil {
		panic(err)
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    logger.ActionDelete,
		ValueType: logger.ValueTypeBroadcastMessage,
		ValueID:   c.Param("id"),
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	a.sendOnChange()

	c.JSON(http.StatusNoContent, nil)
}

// Get all active broadcast-message
// @Summary Get all active broadcast-message
// @Description Get all active broadcast-message
// @Tags broadcast-messages
// @ID broadcast-messages-get-active
// @Success 200 {object} []BroadcastMessage
// @Router /active-broadcast-message [get]
func (a api) GetActive(c *gin.Context) {
	actives, err := a.store.GetActive(c.Request.Context())
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, actives)
}

func NewApi(
	store Store,
	onChangeListener chan<- bool,
	actionLogger logger.ActionLogger,
) API {
	return &api{
		store:            store,
		onChangeListener: onChangeListener,
		actionLogger:     actionLogger,
	}
}

func (a *api) sendOnChange() {
	select {
	case a.onChangeListener <- true:
	default:
	}
}
