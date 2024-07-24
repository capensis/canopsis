package viewgroup

import (
	"errors"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/middleware"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"github.com/gin-gonic/gin"
)

type api struct {
	store Store
}

func NewApi(
	store Store,
) common.CrudAPI {
	return &api{
		store: store,
	}
}

// List
// @Success 200 {object} common.PaginatedListResponse{data=[]ViewGroup}
func (a *api) List(c *gin.Context) {
	var query ListRequest
	query.Query = pagination.GetDefaultQuery()
	query.UserID = c.MustGet(auth.UserKey).(string)

	if err := c.ShouldBind(&query); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, query))
		return
	}

	var authorizedIds []string
	ids, ok := c.Get(middleware.AuthorizedIds)
	if ok {
		if s, ok := ids.([]string); ok {
			authorizedIds = s
		}
	}

	var ownedIds []string
	ids, ok = c.Get(middleware.OwnedIds)
	if ok {
		if s, ok := ids.([]string); ok {
			ownedIds = s
		}
	}

	viewgroups, err := a.store.Find(c, query, authorizedIds, ownedIds)
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
	viewgroup, err := a.store.GetOneBy(c, c.Param("id"))
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

	group, err := a.store.Insert(c, request)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusCreated, group)
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

	group, err := a.store.Update(c, request)
	if err != nil {
		panic(err)
	}

	if group == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, group)
}

func (a *api) Delete(c *gin.Context) {
	ok, err := a.store.Delete(c, c.Param("id"), c.MustGet(auth.UserKey).(string))
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

	c.Status(http.StatusNoContent)
}
