package widget

import (
	"context"
	"errors"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
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
	UpdateGridPositions(c *gin.Context)
}

type api struct {
	store       Store
	enforcer    security.Enforcer
	transformer *RequestTransformer
}

func NewApi(
	store Store,
	enforcer security.Enforcer,
	transformer *RequestTransformer,
) API {
	return &api{
		store:       store,
		enforcer:    enforcer,
		transformer: transformer,
	}
}

// Get
// @Success 200 {object} Response
func (a *api) Get(c *gin.Context) {
	widget, err := a.store.GetOneBy(c, c.Param("id"))
	if err != nil {
		panic(err)
	}

	if widget == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, widget)
}

// Create
// @Param body body EditRequest true "body"
// @Success 201 {object} Response
func (a *api) Create(c *gin.Context) {
	request := CreateRequest{}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	err := a.transformer.Transform(c, &request.EditRequest)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}
		panic(err)
	}

	widget, err := a.store.Insert(c, request)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}

		panic(err)
	}

	c.JSON(http.StatusCreated, widget)
}

// Update
// @Param body body EditRequest true "body"
// @Success 200 {object} Response
func (a *api) Update(c *gin.Context) {
	request := UpdateRequest{
		ID: c.Param("id"),
	}

	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	err := a.transformer.Transform(c, &request.EditRequest)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}
		panic(err)
	}

	widget, err := a.store.Update(c, request)
	if err != nil {
		panic(err)
	}

	if widget == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, widget)
}

func (a *api) Delete(c *gin.Context) {
	ok, err := a.store.Delete(c, c.Param("id"), c.MustGet(auth.UserKey).(string))
	if err != nil {
		panic(err)
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.Status(http.StatusNoContent)
}

// Copy
// @Param body body EditRequest true "body"
// @Success 201 {object} Response
func (a *api) Copy(c *gin.Context) {
	request := CreateRequest{}

	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	widget, err := a.store.Copy(c, c.Param("id"), request)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}

		panic(err)
	}

	if widget == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusCreated, widget)
}

// UpdateGridPositions
// @Param body body []EditGridPositionItemRequest true "body"
func (a *api) UpdateGridPositions(c *gin.Context) {
	request := EditGridPositionRequest{}

	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	userID := c.MustGet(auth.UserKey).(string)
	ids := make([]string, len(request.Items))
	for i, item := range request.Items {
		ids[i] = item.ID
	}
	ok, err := a.checkAccess(c, ids, userID, model.PermissionUpdate)
	if err != nil {
		panic(err)
	}

	if !ok {
		c.JSON(http.StatusForbidden, common.ForbiddenResponse)
		return
	}

	ok, err = a.store.UpdateGridPositions(c, request.Items)
	if err != nil {
		valErr := ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.ErrorResponse{Error: err.Error()})
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

func (a *api) checkAccess(ctx context.Context, ids []string, userID, perm string) (bool, error) {
	tabInfos, err := a.store.FindTabPrivacySettings(ctx, ids)
	if err != nil || len(tabInfos) != len(ids) {
		return false, err
	}

	for _, tabInfo := range tabInfos {
		if tabInfo.IsPrivate && tabInfo.Author == userID {
			continue
		}

		ok, err := a.enforcer.Enforce(userID, tabInfo.View, perm)
		if err != nil || !ok {
			return false, err
		}
	}

	return true, nil
}
