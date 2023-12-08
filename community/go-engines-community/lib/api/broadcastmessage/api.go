package broadcastmessage

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
	common.CrudAPI
	GetActive(c *gin.Context)
}

type api struct {
	store            Store
	onChangeListener chan<- bool
	actionLogger     logger.ActionLogger
}

// Create
// @Param body body BroadcastMessage true "body"
// @Success 201 {object} BroadcastMessage
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

// List
// @Success 200 {object} common.PaginatedListResponse{data=[]BroadcastMessage}
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

// Get
// @Success 200 {object} BroadcastMessage
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

// Update
// @Param body body Payload true "body"
// @Success 200 {object} BroadcastMessage
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

// GetActive
// @Success 200 {array} BroadcastMessage
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
