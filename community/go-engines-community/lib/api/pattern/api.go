package pattern

import (
	"context"
	"net/http"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/bulk"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	apisecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type API interface {
	common.CrudAPI
	BulkDelete(c *gin.Context)
	CountAlarms(c *gin.Context)
	CountEntities(c *gin.Context)
}

type api struct {
	store    Store
	enforcer security.Enforcer
	logger   zerolog.Logger

	configProvider config.UserInterfaceConfigProvider
}

func NewApi(
	store Store,
	configProvider config.UserInterfaceConfigProvider,
	enforcer security.Enforcer,
	logger zerolog.Logger,
) API {
	return &api{
		store:    store,
		enforcer: enforcer,
		logger:   logger,

		configProvider: configProvider,
	}
}

// Create
// @Param body body EditRequest true "body"
// @Success 201 {object} Response
func (a *api) Create(c *gin.Context) {
	request := EditRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	if *request.IsCorporate {
		ok, err := a.enforcer.Enforce(c.MustGet(auth.UserKey).(string), apisecurity.PermCorporatePattern, model.PermissionCan)
		if err != nil {
			panic(err)
		}

		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, common.ForbiddenResponse)
			return
		}
	}

	pattern, err := a.store.Insert(c, request)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusCreated, pattern)
}

// List
// @Success 200 {object} common.PaginatedListResponse{data=[]Response}
func (a *api) List(c *gin.Context) {
	var request ListRequest
	request.Query = pagination.GetDefaultQuery()

	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	aggregationResult, err := a.store.Find(c, request, c.MustGet(auth.UserKey).(string))
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

// Get
// @Success 200 {object} Response
func (a *api) Get(c *gin.Context) {
	pattern, err := a.store.GetByID(c, c.Param("id"), c.MustGet(auth.UserKey).(string))
	if err != nil {
		panic(err)
	}

	if pattern == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, pattern)
}

// Update
// @Param body body EditRequest true "body"
// @Success 200 {object} Response
func (a *api) Update(c *gin.Context) {
	userID := c.MustGet(auth.UserKey).(string)
	request := EditRequest{
		ID: c.Param("id"),
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	pattern, err := a.store.GetByID(c, request.ID, userID)
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
		ok, err := a.enforcer.Enforce(userID, apisecurity.PermCorporatePattern, model.PermissionCan)
		if err != nil {
			panic(err)
		}

		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, common.ForbiddenResponse)
			return
		}
	}

	pattern, err = a.store.Update(c, request)
	if err != nil {
		panic(err)
	}

	if pattern == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, pattern)
}

func (a *api) Delete(c *gin.Context) {
	userID := c.MustGet(auth.UserKey).(string)

	pattern, err := a.store.GetByID(c, c.Param("id"), userID)
	if err != nil {
		panic(err)
	}

	if pattern == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	if pattern.IsCorporate {
		ok, err := a.enforcer.Enforce(userID, apisecurity.PermCorporatePattern, model.PermissionCan)
		if err != nil {
			panic(err)
		}

		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, common.ForbiddenResponse)
			return
		}
	}

	ok, err := a.store.Delete(c, *pattern, userID)
	if err != nil {
		panic(err)
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// BulkDelete
// @Param body body []BulkDeleteRequestItem true "body"
func (a *api) BulkDelete(c *gin.Context) {
	userID := c.MustGet(auth.UserKey).(string)

	canDeleteCorporate, err := a.enforcer.Enforce(userID, apisecurity.PermCorporatePattern, model.PermissionCan)
	if err != nil {
		panic(err)
	}

	bulk.Handler(c, func(request BulkDeleteRequestItem) (string, error) {
		pattern, err := a.store.GetByID(c, request.ID, userID)
		if err != nil || pattern == nil {
			return "", err
		}

		if pattern.IsCorporate && !canDeleteCorporate {
			return "", bulk.ErrUnauthorized
		}

		ok, err := a.store.Delete(c, *pattern, userID)
		if err != nil || !ok {
			return "", err
		}

		return pattern.ID, nil
	}, a.logger)
}

// CountAlarms
// @Param body body CountRequest true "body"
// @Success 200 {object} CountAlarmsResponse
func (a *api) CountAlarms(c *gin.Context) {
	request := CountRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	conf := a.configProvider.Get()
	ctx, cancel := context.WithTimeout(c, time.Duration(conf.CheckCountRequestTimeout)*time.Second)
	defer cancel()

	res, err := a.store.CountAlarms(ctx, request, int64(conf.MaxMatchedItems))
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, res)
}

// CountEntities
// @Param body body CountRequest true "body"
// @Success 200 {object} CountEntitiesResponse
func (a *api) CountEntities(c *gin.Context) {
	request := CountRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	conf := a.configProvider.Get()
	ctx, cancel := context.WithTimeout(c, time.Duration(conf.CheckCountRequestTimeout)*time.Second)
	defer cancel()

	res, err := a.store.CountEntities(ctx, request, int64(conf.MaxMatchedItems))
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, res)
}
