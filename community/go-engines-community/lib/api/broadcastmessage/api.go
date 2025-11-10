package broadcastmessage

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/crud"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/websocket"
	"github.com/gin-gonic/gin"
)

type API interface {
	crud.API
	GetActive(c *gin.Context)
	Read(c *gin.Context)
}

func NewAPI(
	store Store,
	onChangeListener chan<- bool,
	websocketHub websocket.Hub,
) API {
	return &api{
		store:            store,
		onChangeListener: onChangeListener,
		websocketHub:     websocketHub,
	}
}

type api struct {
	store            Store
	onChangeListener chan<- bool
	websocketHub     websocket.Hub
}

// Create
// @Param body body CreateRequest true "body"
// @Success 201 {object} Response
func (a *api) Create(c *gin.Context) {
	request := CreateRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	res, err := a.store.Insert(c, request)
	if err != nil {
		panic(err)
	}

	if res == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	a.sendOnChange()

	c.JSON(http.StatusCreated, res)
}

// List
// @Success 200 {object} pagination.ListResponse{data=[]Response}
func (a *api) List(c *gin.Context) {
	var query ListRequest
	query.Query = pagination.GetDefaultQuery()

	if err := c.ShouldBind(&query); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, query))
		return
	}

	aggregationResult, err := a.store.Find(c, query)
	if err != nil {
		panic(err)
	}

	res := pagination.NewResponse(query.Query, aggregationResult)
	c.JSON(http.StatusOK, res)
}

// Get
// @Success 200 {object} Response
func (a *api) Get(c *gin.Context) {
	bm, err := a.store.GetByID(c, c.Param("id"))
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
// @Param body body UpdateRequest true "body"
// @Success 200 {object} Response
func (a *api) Update(c *gin.Context) {
	request := UpdateRequest{
		ID: c.Param("id"),
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))

		return
	}

	res, err := a.store.Update(c, request)
	if err != nil {
		panic(err)
	}

	if res == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	a.sendOnChange()

	c.JSON(http.StatusOK, res)
}

func (a *api) Delete(c *gin.Context) {
	ok, err := a.store.Delete(c, c.Param("id"), c.MustGet(authctx.UserKey).(string))
	if err != nil {
		panic(err)
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)

		return
	}

	a.sendOnChange()

	c.JSON(http.StatusNoContent, nil)
}

// GetActive
// @Success 200 {array} Response
func (a *api) GetActive(c *gin.Context) {
	userID := ""
	if v, ok := c.Get(authctx.UserKey); ok {
		userID, _ = v.(string)
	}

	res, err := a.store.GetActive(c, userID)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, res)
}

func (a *api) Read(c *gin.Context) {
	userID := c.MustGet(authctx.UserKey).(string)
	ok, err := a.store.Read(c, c.Param("id"), userID)
	if err != nil {
		panic(err)
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)

		return
	}

	msgs, err := a.store.GetActive(c, userID)
	if err != nil {
		panic(err)
	}

	a.websocketHub.SendToUser(userID, websocket.RoomBroadcastMessages, msgs)

	c.Status(http.StatusNoContent)
}

func (a *api) sendOnChange() {
	select {
	case a.onChangeListener <- true:
	default:
	}
}
