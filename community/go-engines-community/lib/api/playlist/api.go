package playlist

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/middleware"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/viewtab"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type api struct {
	store        Store
	tabStore     viewtab.Store
	enforcer     security.Enforcer
	actionLogger logger.ActionLogger
}

func NewApi(
	store Store,
	tabStore viewtab.Store,
	enforcer security.Enforcer,
	actionLogger logger.ActionLogger,
) common.CrudAPI {
	return &api{
		store:        store,
		tabStore:     tabStore,
		enforcer:     enforcer,
		actionLogger: actionLogger,
	}
}

// Find all playlist
// @Summary Find all playlist
// @Description Get paginated list of playlist
// @Tags playlists
// @ID playlists-find-all
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param page query integer true "current page"
// @Param limit query integer true "items per page"
// @Param search query string false "search query"
// @Param sort query string false "sort query"
// @Param sort_by query string false "sort query"
// @Success 200 {object} common.PaginatedListResponse{data=[]Playlist}
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /playlists [get]
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

// Get playlist by id
// @Summary Get playlist by id
// @Description Get playlist by id
// @Tags playlists
// @ID playlists-get-by-id
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "playlist id"
// @Success 200 {object} Playlist
// @Failure 404 {object} common.ErrorResponse
// @Router /playlists/{id} [get]
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

// Create playlist
// @Summary Create playlist
// @Description Create playlist
// @Tags playlists
// @ID playlists-create
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param body body EditRequest true "body"
// @Success 201 {object} Playlist
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /playlists [post]
func (a *api) Create(c *gin.Context) {
	request := EditRequest{}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	userId := c.MustGet(auth.UserKey).(string)
	ok, err := a.checkAccess(c.Request.Context(), request.TabsList, userId)
	if err != nil {
		panic(err)
	}
	if !ok {
		c.AbortWithStatusJSON(http.StatusForbidden, common.ForbiddenResponse)
		return
	}

	playlist, err := a.store.Insert(c.Request.Context(), userId, request)
	if err != nil {
		panic(err)
	}

	err = a.actionLogger.Action(c, logger.LogEntry{
		Action:    logger.ActionCreate,
		ValueType: logger.ValueTypePlayList,
		ValueID:   playlist.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusCreated, playlist)
}

// Update playlist by id
// @Summary Update playlist by id
// @Description Update playlist by id
// @Tags playlists
// @ID playlists-update-by-id
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "playlist id"
// @Param body body Playlist true "body"
// @Success 200 {object} Playlist
// @Failure 400 {object} common.ValidationErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /playlists/{id} [put]
func (a *api) Update(c *gin.Context) {
	request := EditRequest{
		ID: c.Param("id"),
	}

	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	userId := c.MustGet(auth.UserKey).(string)
	ok, err := a.checkAccess(c.Request.Context(), request.TabsList, userId)
	if err != nil {
		panic(err)
	}
	if !ok {
		c.AbortWithStatusJSON(http.StatusForbidden, common.ForbiddenResponse)
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

	err = a.actionLogger.Action(c, logger.LogEntry{
		Action:    logger.ActionUpdate,
		ValueType: logger.ValueTypePlayList,
		ValueID:   playlist.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusOK, playlist)
}

// Delete playlist by id
// @Summary Delete playlist by id
// @Description Delete playlist by id
// @Tags playlists
// @ID playlists-delete-by-id
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "playlist id"
// @Success 204
// @Failure 404 {object} common.ErrorResponse
// @Router /playlists/{id} [delete]
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
		ValueType: logger.ValueTypePlayList,
		ValueID:   id,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.Status(http.StatusNoContent)
}

func (a *api) checkAccess(ctx context.Context, tabIds []string, userId string) (bool, error) {
	tabs, err := a.tabStore.Find(ctx, tabIds)
	if err != nil || len(tabs) != len(tabIds) {
		return false, err
	}

	for _, tab := range tabs {
		ok, err := a.enforcer.Enforce(userId, tab.View, model.PermissionRead)
		if err != nil || !ok {
			return false, err
		}
	}

	return true, nil
}
