package view

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"github.com/gin-gonic/gin"
)

type API interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	Copy(c *gin.Context)
	UpdatePositions(c *gin.Context)
	Import(c *gin.Context)
	Export(c *gin.Context)
}

type api struct {
	store        Store
	enforcer     security.Enforcer
	actionLogger logger.ActionLogger
}

func NewApi(
	store Store,
	enforcer security.Enforcer,
	actionLogger logger.ActionLogger,
) API {
	return &api{
		store:        store,
		enforcer:     enforcer,
		actionLogger: actionLogger,
	}
}

// Get
// @Success 200 {object} Response
func (a *api) Get(c *gin.Context) {
	view, err := a.store.GetOneBy(c, c.Param("id"))
	if err != nil {
		panic(err)
	}

	if view == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, view)
}

// Create
// @Param body body EditRequest true "body"
// @Success 201 {object} Response
func (a *api) Create(c *gin.Context) {
	request := EditRequest{}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	view, err := a.store.Insert(c, request, true)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}

		panic(err)
	}

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    logger.ActionCreate,
		ValueType: logger.ValueTypeView,
		ValueID:   view.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusCreated, view)
}

// Update
// @Param body body EditRequest true "body"
// @Success 200 {object} Response
func (a *api) Update(c *gin.Context) {
	request := EditRequest{
		ID: c.Param("id"),
	}

	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	view, err := a.store.Update(c, request)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}

		panic(err)
	}

	if view == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    logger.ActionUpdate,
		ValueType: logger.ValueTypeView,
		ValueID:   view.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusOK, view)
}

func (a *api) Delete(c *gin.Context) {
	id := c.Param("id")

	ok, err := a.store.Delete(c, id)
	if err != nil {
		panic(err)
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    logger.ActionDelete,
		ValueType: logger.ValueTypeView,
		ValueID:   id,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.Status(http.StatusNoContent)
}

// Copy
// @Param body body EditRequest true "body"
// @Success 201 {object} Response
func (a *api) Copy(c *gin.Context) {
	request := EditRequest{}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	view, err := a.store.Copy(c, c.Param("id"), request)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}

		panic(err)
	}

	if view == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    logger.ActionCreate,
		ValueType: logger.ValueTypeView,
		ValueID:   view.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusCreated, view)
}

// UpdatePositions
// @Param body body []EditPositionItemRequest true "body"
func (a *api) UpdatePositions(c *gin.Context) {
	userId := c.MustGet(auth.UserKey).(string)
	request := EditPositionRequest{}

	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	for _, item := range request.Items {
		for _, view := range item.Views {
			ok, err := a.enforcer.Enforce(userId, view, model.PermissionUpdate)
			if err != nil {
				panic(err)
			}
			if !ok {
				c.AbortWithStatusJSON(http.StatusForbidden, common.ForbiddenResponse)
				return
			}
		}
	}

	ok, err := a.store.UpdatePositions(c, request)
	if err != nil {
		panic(err)
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.Status(http.StatusNoContent)
}

// Import
// @Param body body []ImportItemRequest true "body"
func (a *api) Import(c *gin.Context) {
	userId := c.MustGet(auth.UserKey).(string)
	request := ImportRequest{}

	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	for _, group := range request.Items {
		if group.Views == nil {
			continue
		}
		for _, view := range group.Views {
			if view.ID == "" {
				continue
			}
			ok, err := a.enforcer.Enforce(userId, view.ID, model.PermissionUpdate)
			if err != nil {
				panic(err)
			}
			if !ok {
				c.AbortWithStatusJSON(http.StatusForbidden, common.ForbiddenResponse)
				return
			}
		}
	}

	err := a.store.Import(c, request, userId)
	if err != nil {
		valError := ValidationError{}
		if errors.As(err, &valError) {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.ValidationErrorResponse{
				Errors: map[string]string{
					valError.field: valError.Error(),
				},
			})
			return
		}
		panic(err)
	}

	c.Status(http.StatusNoContent)
}

// Export
// @Param body body ExportRequest true "body"
// @Success 200 {object} ExportResponse
func (a *api) Export(c *gin.Context) {
	userId := c.MustGet(auth.UserKey).(string)
	request := ExportRequest{}

	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	for _, group := range request.Groups {
		for _, view := range group.Views {
			ok, err := a.enforcer.Enforce(userId, view, model.PermissionRead)
			if err != nil {
				panic(err)
			}
			if !ok {
				c.AbortWithStatusJSON(http.StatusForbidden, common.ForbiddenResponse)
				return
			}
		}
	}
	for _, view := range request.Views {
		ok, err := a.enforcer.Enforce(userId, view, model.PermissionRead)
		if err != nil {
			panic(err)
		}
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, common.ForbiddenResponse)
			return
		}
	}

	response, err := a.store.Export(c, request)
	if err != nil {
		valError := ValidationError{}
		if errors.As(err, &valError) {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.ValidationErrorResponse{
				Errors: map[string]string{
					valError.field: valError.Error(),
				},
			})
			return
		}
		panic(err)
	}

	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fmt.Sprintf("views-%s.json", time.Now().Format("2006-01-02T15-04-05"))))

	c.JSON(http.StatusOK, response)
}
