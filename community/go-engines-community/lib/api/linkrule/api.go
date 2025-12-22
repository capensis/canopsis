package linkrule

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/bulk"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/crud"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/dbexport"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/gin-gonic/gin"
)

type API interface {
	crud.API
	GetCategories(*gin.Context)
	BulkDelete(c *gin.Context)
	DBExport(c *gin.Context)
	ValidateTemplates(c *gin.Context)
	GetTemplateVars(c *gin.Context)
}

type api struct {
	store          Store
	mongoExporter  dbexport.Exporter
	errorResponder httperror.Responder
}

func NewApi(
	store Store,
	mongoExporter dbexport.Exporter,
	errorResponder httperror.Responder,
) API {
	return &api{
		store:          store,
		mongoExporter:  mongoExporter,
		errorResponder: errorResponder,
	}
}

// Create
// @Param body body EditRequest true "body"
// @Success 201 {object} Response
func (a *api) Create(c *gin.Context) {
	request := EditRequest{}
	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	rule, err := a.store.Insert(c, request)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	c.JSON(http.StatusCreated, rule)
}

// List
// @Success 200 {object} pagination.ListResponse{data=[]Response}
func (a *api) List(c *gin.Context) {
	var request ListRequest
	request.Query = pagination.GetDefaultQuery()

	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	aggregationResult, err := a.store.Find(c, request)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	res := pagination.NewResponse(request.Query, aggregationResult)
	c.JSON(http.StatusOK, res)
}

// Get
// @Success 200 {object} Response
func (a *api) Get(c *gin.Context) {
	rule, err := a.store.GetByID(c, c.Param("id"))
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if rule == nil {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	c.JSON(http.StatusOK, rule)
}

// GetCategories
// @Success 200 {object} CategoryResponse
func (a *api) GetCategories(c *gin.Context) {
	var r CategoriesRequest
	if err := validation.Bind(c, &r); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	categories, err := a.store.GetCategories(c, r)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	c.JSON(http.StatusOK, categories)
}

// Update
// @Param body body EditRequest true "body"
// @Success 200 {object} Response
func (a *api) Update(c *gin.Context) {
	request := EditRequest{
		ID: c.Param("id"),
	}
	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	rule, err := a.store.Update(c, request)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if rule == nil {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	c.JSON(http.StatusOK, rule)
}

func (a *api) Delete(c *gin.Context) {
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	ok, err := a.store.Delete(c, c.Param("id"), userID)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if !ok {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// BulkDelete
// @Param body body []BulkDeleteRequestItem true "body"
func (a *api) BulkDelete(c *gin.Context) {
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	bulk.Handler(c, func(request BulkDeleteRequestItem) (string, error) {
		ok, err := a.store.Delete(c, request.ID, userID)
		if err != nil {
			return "", err
		}

		if !ok {
			return "", httperror.ErrNotFound
		}

		return request.ID, nil
	}, a.errorResponder)
}

// DBExport
// @Param body body dbexport.Request true "body"
func (a *api) DBExport(c *gin.Context) {
	request := dbexport.Request{}

	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	b, err := a.mongoExporter.Export(c, mongo.LinkRuleMongoCollection, request)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	dbexport.AttachFile(c, mongo.LinkRuleMongoCollection, b)
}

// ValidateTemplates
// @Param body body TemplateRequest true "body"
// @Success 200 {object} template.ValidateResponse
func (a *api) ValidateTemplates(c *gin.Context) {
	request := TemplateRequest{}
	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	response, err := a.store.ValidateTemplates(c, request)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	c.JSON(http.StatusOK, response)
}

// GetTemplateVars
// @Success 200 {array} TemplateVarsResponse
func (a *api) GetTemplateVars(c *gin.Context) {
	vars, err := a.store.GetTemplateVars(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	c.JSON(http.StatusOK, vars)
}
