package template

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
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

	CreateTest(c *gin.Context)
	ListTest(c *gin.Context)
	GetTest(c *gin.Context)
	UpdateTest(c *gin.Context)
	DeleteTest(c *gin.Context)
}

type api struct {
	store                  Store
	templateConfigProvider config.TemplateConfigProvider
	errorResponder         httperror.Responder
	logger                 zerolog.Logger
}

func NewAPI(
	store Store,
	templateConfigProvider config.TemplateConfigProvider,
	errorResponder httperror.Responder,
	logger zerolog.Logger,
) API {
	return &api{
		store:                  store,
		templateConfigProvider: templateConfigProvider,
		errorResponder:         errorResponder,
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
	if err := validation.Bind(c, &r); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	res, err := a.store.CreateData(c, r)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	c.JSON(http.StatusCreated, res)
}

// ListData
// @Success 200 {object} pagination.ListResponse{data=[]DataResponse}
func (a *api) ListData(c *gin.Context) {
	var r ListDataRequest
	r.Query = pagination.GetDefaultQuery()
	if err := validation.Bind(c, &r); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	aggregationResult, err := a.store.FindData(c, r)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	res := pagination.NewResponse(r.Query, &aggregationResult)
	c.JSON(http.StatusOK, res)
}

// GetData
// @Success 200 {object} DataResponse
func (a *api) GetData(c *gin.Context) {
	res, err := a.store.GetData(c, c.Param("id"))
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if res.ID == "" {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

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
	if err := validation.Bind(c, &r); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	res, err := a.store.UpdateData(c, r)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if res.ID == "" {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	c.JSON(http.StatusOK, res)
}

func (a *api) DeleteData(c *gin.Context) {
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	ok, err := a.store.DeleteData(c, c.Param("id"), userID)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	if !ok {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	c.Status(http.StatusNoContent)
}

// CreateTest
// @Param body body EditTestRequest true "body"
// @Success 200 {object} TestResponse
func (a *api) CreateTest(c *gin.Context) {
	r := EditTestRequest{}
	if err := validation.Bind(c, &r); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	res, err := a.store.CreateTest(c, r)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	c.JSON(http.StatusCreated, res)
}

// ListTest
// @Success 200 {object} pagination.ListResponse{data=[]TestResponse}
func (a *api) ListTest(c *gin.Context) {
	var r ListTestRequest
	r.Query = pagination.GetDefaultQuery()
	if err := validation.Bind(c, &r); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	aggregationResult, err := a.store.FindTest(c, r, userID)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	res := pagination.NewResponse(r.Query, &aggregationResult)
	c.JSON(http.StatusOK, res)
}

// GetTest
// @Success 200 {object} TestResponse
func (a *api) GetTest(c *gin.Context) {
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	res, err := a.store.GetTest(c, c.Param("id"), userID)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if res.ID == "" {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	c.JSON(http.StatusOK, res)
}

// UpdateTest
// @Param body body EditTestRequest true "body"
// @Success 200 {object} TestResponse
func (a *api) UpdateTest(c *gin.Context) {
	r := EditTestRequest{
		ID: c.Param("id"),
	}
	if err := validation.Bind(c, &r); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	res, err := a.store.UpdateTest(c, r)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if res.ID == "" {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	c.JSON(http.StatusOK, res)
}

func (a *api) DeleteTest(c *gin.Context) {
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	ok, err := a.store.DeleteTest(c, c.Param("id"), userID)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if !ok {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	c.Status(http.StatusNoContent)
}
