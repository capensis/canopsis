package template

import (
	"errors"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/bulk"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type API interface {
	GetEnvVars(c *gin.Context)

	CreateData(c *gin.Context)
	ListData(c *gin.Context)
	GetData(c *gin.Context)
	UpdateData(c *gin.Context)
	DeleteData(c *gin.Context)
	BulkDeleteData(c *gin.Context)
}

type api struct {
	store                  Store
	templateConfigProvider config.TemplateConfigProvider
	logger                 zerolog.Logger
}

func NewAPI(store Store, templateConfigProvider config.TemplateConfigProvider, logger zerolog.Logger) API {
	return &api{
		store:                  store,
		templateConfigProvider: templateConfigProvider,
		logger:                 logger,
	}
}

func (a *api) GetEnvVars(c *gin.Context) {
	c.JSON(http.StatusOK, a.templateConfigProvider.Get().Vars)
}

// CreateData
// @Param body body EditDataRequest true "body"
// @Success 200 {object} DataResponse
func (a *api) CreateData(c *gin.Context) {
	r := EditDataRequest{}
	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))

		return
	}

	res, err := a.store.CreateData(c, r)
	if err != nil {
		validationError := common.ValidationError{}
		if errors.As(err, &validationError) {
			c.AbortWithStatusJSON(http.StatusBadRequest, validationError.ValidationErrorResponse())

			return
		}

		panic(err)
	}

	c.JSON(http.StatusCreated, res)
}

// ListData
// @Success 200 {object} common.PaginatedListResponse{data=[]DataResponse}
func (a *api) ListData(c *gin.Context) {
	var r ListDataRequest
	r.Query = pagination.GetDefaultQuery()
	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))

		return
	}

	aggregationResult, err := a.store.FindData(c, r)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())

			return
		}

		panic(err)
	}

	res, err := common.NewPaginatedResponse(r.Query, &aggregationResult)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))

		return
	}

	c.JSON(http.StatusOK, res)
}

// GetData
// @Success 200 {object} DataResponse
func (a *api) GetData(c *gin.Context) {
	res, err := a.store.GetData(c, c.Param("id"))
	if err != nil {
		panic(err)
	}

	if res.ID == "" {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)

		return
	}

	c.JSON(http.StatusOK, res)
}

// UpdateData
// @Param body body EditDataRequest true "body"
// @Success 200 {object} DataResponse
func (a *api) UpdateData(c *gin.Context) {
	r := EditDataRequest{
		ID: c.Param("id"),
	}
	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))

		return
	}

	res, err := a.store.UpdateData(c, r)
	if err != nil {
		validationError := common.ValidationError{}
		if errors.As(err, &validationError) {
			c.AbortWithStatusJSON(http.StatusBadRequest, validationError.ValidationErrorResponse())

			return
		}

		panic(err)
	}

	if res.ID == "" {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)

		return
	}

	c.JSON(http.StatusOK, res)
}

func (a *api) DeleteData(c *gin.Context) {
	userID := c.MustGet(auth.UserKey).(string)
	ok, err := a.store.DeleteData(c, c.Param("id"), userID)
	if err != nil {
		panic(err)
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)

		return
	}

	c.Status(http.StatusNoContent)
}

// BulkDeleteData
// @Param body body []BulkDeleteRequestItem true "body"
func (a *api) BulkDeleteData(c *gin.Context) {
	userID := c.MustGet(auth.UserKey).(string)
	bulk.Handler(c, func(r BulkDeleteRequestItem) (string, error) {
		ok, err := a.store.DeleteData(c, r.ID, userID)
		if err != nil || !ok {
			return "", err
		}

		return r.ID, nil
	}, a.logger)
}
