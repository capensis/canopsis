package viewstats

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/session/stats"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/gin-gonic/gin"
)

type API interface {
	Create(c *gin.Context)
	Update(c *gin.Context)
	List(c *gin.Context)
}

func NewApi(manager stats.Manager) API {
	return &api{
		manager: manager,
	}
}

type api struct {
	manager stats.Manager
}

// Create view stats
// @Summary Create view stats
// @Description Create view stats
// @Tags view-stats
// @ID view-stats-create
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Success 201 {object} stats.Stats
// @Router /view-stats [post]
func (a *api) Create(c *gin.Context) {
	s, err := a.manager.Ping(
		c.Request.Context(),
		stats.SessionData{
			SessionID: utils.NewID(),
			UserID:    c.MustGet(auth.UserKey).(string),
		},
		stats.PathData{},
	)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusCreated, s)
}

// Update view stats by id
// @Summary Update view stats by id
// @Description Update view stats by id
// @Tags view-stats
// @ID view-stats-update-by-id
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "view stats id"
// @Param body body pingRequest true "body"
// @Success 200 {object} stats.Stats
// @Failure 400 {object} common.ValidationErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /view-stats/{id} [put]
func (a *api) Update(c *gin.Context) {
	var r pingRequest

	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))
		return
	}

	viewID := r.Path[0]
	tabID := ""
	if len(r.Path) > 1 {
		tabID = r.Path[1]
	}

	id := c.Param("id")
	ok, err := a.manager.Exists(c.Request.Context(), id)
	if err != nil {
		panic(err)
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	s, err := a.manager.Ping(
		c.Request.Context(),
		stats.SessionData{
			SessionID: id,
			UserID:    c.MustGet(auth.UserKey).(string),
		},
		stats.PathData{
			ViewID:  viewID,
			TabID:   tabID,
			Visible: r.Visible != nil && *r.Visible,
		},
	)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, s)
}

// Find view stats
// @Summary Find view stats
// @Description Get list of view stats
// @Tags view-stats
// @ID view-stats-find-all
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param query query listRequest true "request"
// @Success 200 {object} listResponse
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /view-stats [get]
func (a *api) List(c *gin.Context) {
	var r listRequest

	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))
		return
	}

	s, err := a.manager.Find(c.Request.Context(), stats.Filter{
		IsActive:      r.IsActive,
		Usernames:     r.Usernames,
		StartedAfter:  r.StartedAfter,
		StoppedBefore: r.StoppedBefore,
	})
	if err != nil {
		panic(err)
	}

	res := listResponse{
		Data: s,
	}
	c.JSON(http.StatusOK, res)
}
