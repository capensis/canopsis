package scenario

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
	crud.BulkAPI
	DBExport(c *gin.Context)
	ValidateTemplates(c *gin.Context)
	GetTemplateVars(c *gin.Context)
	BulkEnable(c *gin.Context)
	BulkDisable(c *gin.Context)
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

// List
// @Success 200 {object} pagination.ListResponse{data=[]Scenario}
func (a *api) List(c *gin.Context) {
	var query FilteredQuery
	query.Query = pagination.GetDefaultQuery()
	if err := validation.Bind(c, &query); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	scenarios, err := a.store.Find(c, query)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	res := pagination.NewResponse(query.Query, scenarios)
	c.JSON(http.StatusOK, res)
}

// Get
// @Success 200 {object} Scenario
func (a *api) Get(c *gin.Context) {
	scenario, err := a.store.GetOneBy(c, c.Param("id"))
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if scenario == nil {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	c.JSON(http.StatusOK, scenario)
}

// Create
// @Param body body EditRequest true "body"
// @Success 201 {object} Scenario
func (a *api) Create(c *gin.Context) {
	var request CreateRequest
	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	scenario, err := a.store.Insert(c, request)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if scenario == nil {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	c.JSON(http.StatusCreated, scenario)
}

// Update
// @Param body body EditRequest true "body"
// @Success 200 {object} Scenario
func (a *api) Update(c *gin.Context) {
	request := UpdateRequest{
		ID: c.Param("id"),
	}

	oldScenario, err := a.store.GetOneBy(c, request.ID)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	if oldScenario == nil {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	newScenario, err := a.store.Update(c, request)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if newScenario == nil {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	c.JSON(http.StatusOK, newScenario)
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

	c.Status(http.StatusNoContent)
}

// BulkCreate
// @Param body body []CreateRequest true "body"
func (a *api) BulkCreate(c *gin.Context) {
	bulk.Handler(c, func(request CreateRequest) (string, error) {
		scenario, err := a.store.Insert(c, request)
		if err != nil {
			return "", err
		}

		return scenario.ID, nil
	}, a.errorResponder)
}

// BulkUpdate
// @Param body body []BulkUpdateRequestItem true "body"
func (a *api) BulkUpdate(c *gin.Context) {
	bulk.Handler(c, func(request BulkUpdateRequestItem) (string, error) {
		scenario, err := a.store.Update(c, UpdateRequest(request))
		if err != nil {
			return "", err
		}

		if scenario == nil {
			return "", httperror.ErrNotFound
		}

		return scenario.ID, nil
	}, a.errorResponder)
}

// BulkDelete
// @Param body body []BulkDeleteRequestItem true "body"
func (a *api) BulkDelete(c *gin.Context) {
	bulk.Handler(c, func(request BulkDeleteRequestItem) (string, error) {
		ok, err := a.store.Delete(c, request.ID, request.Author)
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

	b, err := a.mongoExporter.Export(c, mongo.ScenarioCollection, request)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	dbexport.AttachFile(c, mongo.ScenarioCollection, b)
}

// ValidateTemplates
// @Param body body TemplateRequest true "body"
// @Success 200 {object} template.ValidateResponse
func (a *api) ValidateTemplates(c *gin.Context) {
	var request TemplateRequest
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

// BulkEnable
// @Param body body []BulkToggleRequestItem true "body"
func (a *api) BulkEnable(c *gin.Context) {
	a.toggle(c, true)
}

// BulkDisable
// @Param body body []BulkToggleRequestItem true "body"
func (a *api) BulkDisable(c *gin.Context) {
	a.toggle(c, false)
}

func (a *api) toggle(c *gin.Context, enabled bool) {
	bulk.Handler(c, func(request BulkToggleRequestItem) (string, error) {
		found, err := a.store.Toggle(c, request, enabled)
		if err != nil {
			return "", err
		}

		if !found {
			return "", httperror.ErrNotFound
		}

		return request.ID, nil
	}, a.errorResponder)
}
