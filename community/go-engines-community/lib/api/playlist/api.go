package playlist

import (
	"context"
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

// List
// @Success 200 {object} common.PaginatedListResponse{data=[]Playlist}
func (a *api) List(c *gin.Context) {
	var r ListRequest
	r.Query = pagination.GetDefaultQuery()

	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))
		return
	}

	ids, ok := c.Get(middleware.AuthorizedIds)
	if ok {
		r.Ids = ids.([]string)
	}

	playlists, err := a.store.Find(c.Request.Context(), r)
	if err != nil {
		panic(err)
	}

	res, err := common.NewPaginatedResponse(r.Query, playlists)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, res)
}

// Get
// @Success 200 {object} Playlist
func (a *api) Get(c *gin.Context) {
	playlist, err := a.store.GetById(c.Request.Context(), c.Param("id"))
	if err != nil {
		panic(err)
	}
	if playlist == nil {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, playlist)
}

// Create
// @Param body body EditRequest true "body"
// @Success 201 {object} Playlist
func (a *api) Create(c *gin.Context) {
	request := EditRequest{}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	userID := c.MustGet(auth.UserKey).(string)

	playlist, err := a.store.Insert(c.Request.Context(), userID, request)
	if err != nil {
		panic(err)
	}

	err = a.actionLogger.Action(context.Background(), userID, logger.LogEntry{
		Action:    logger.ActionCreate,
		ValueType: logger.ValueTypePlayList,
		ValueID:   playlist.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusCreated, playlist)
}

// Update
// @Param body body Playlist true "body"
// @Success 200 {object} Playlist
func (a *api) Update(c *gin.Context) {
	request := EditRequest{
		ID: c.Param("id"),
	}

	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	playlist, err := a.store.Update(c.Request.Context(), request)
	if err != nil {
		panic(err)
	}

	if playlist == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    logger.ActionUpdate,
		ValueType: logger.ValueTypePlayList,
		ValueID:   playlist.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusOK, playlist)
}

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

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    logger.ActionDelete,
		ValueType: logger.ValueTypePlayList,
		ValueID:   id,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.Status(http.StatusNoContent)
}
