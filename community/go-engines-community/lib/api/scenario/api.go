package scenario

import (
	"errors"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/bulk"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/dbexport"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type API interface {
	DBExport(c *gin.Context)
	ValidateTemplates(c *gin.Context)
	GetTemplateVars(c *gin.Context)
	common.BulkCrudAPI
}

type api struct {
	store         Store
	mongoExporter dbexport.Exporter
	logger        zerolog.Logger
}

func NewApi(
	store Store,
	mongoExporter dbexport.Exporter,
	logger zerolog.Logger,
) API {
	return &api{
		store:         store,
		mongoExporter: mongoExporter,
		logger:        logger,
	}
}

// List
// @Success 200 {object} common.PaginatedListResponse{data=[]Scenario}
func (a *api) List(c *gin.Context) {
	var query FilteredQuery
	query.Query = pagination.GetDefaultQuery()

	if err := c.ShouldBind(&query); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, query))
		return
	}

	scenarios, err := a.store.Find(c, query)
	if err != nil {
		panic(err)
	}

	res, err := common.NewPaginatedResponse(query.Query, scenarios)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, res)
}

// Get
// @Success 200 {object} Scenario
func (a *api) Get(c *gin.Context) {
	scenario, err := a.store.GetOneBy(c, c.Param("id"))
	if err != nil {
		panic(err)
	}
	if scenario == nil {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, scenario)
}

// Create
// @Param body body EditRequest true "body"
// @Success 201 {object} Scenario
func (a *api) Create(c *gin.Context) {
	var request CreateRequest
	var err error

	if err = c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	scenario, err := a.store.Insert(c, request)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}
		panic(err)
	}

	if scenario == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
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
		panic(err)
	}
	if oldScenario == nil {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	newScenario, err := a.store.Update(c, request)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}
		panic(err)
	}

	if newScenario == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, newScenario)
}

func (a *api) Delete(c *gin.Context) {
	id := c.Param("id")

	scenario, err := a.store.GetOneBy(c, id)
	if err != nil {
		panic(err)
	}

	if scenario == nil {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	ok, err := a.store.Delete(c, id, c.MustGet(auth.UserKey).(string))
	if err != nil {
		panic(err)
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
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
	}, a.logger)
}

// BulkUpdate
// @Param body body []BulkUpdateRequestItem true "body"
func (a *api) BulkUpdate(c *gin.Context) {
	bulk.Handler(c, func(request BulkUpdateRequestItem) (string, error) {
		oldScenario, err := a.store.GetOneBy(c, request.ID)
		if err != nil || oldScenario == nil {
			return "", err
		}

		scenario, err := a.store.Update(c, UpdateRequest(request))
		if err != nil || scenario == nil {
			return "", err
		}

		return scenario.ID, nil
	}, a.logger)
}

// BulkDelete
// @Param body body []BulkDeleteRequestItem true "body"
func (a *api) BulkDelete(c *gin.Context) {
	userID := c.MustGet(auth.UserKey).(string)

	bulk.Handler(c, func(request BulkDeleteRequestItem) (string, error) {
		scenario, err := a.store.GetOneBy(c, request.ID)
		if err != nil || scenario == nil {
			return "", err
		}

		ok, err := a.store.Delete(c, request.ID, userID)
		if err != nil || !ok {
			return "", err
		}

		return scenario.ID, nil
	}, a.logger)
}

// DBExport
// @Param body body dbexport.Request true "body"
func (a *api) DBExport(c *gin.Context) {
	request := dbexport.Request{}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	b, err := a.mongoExporter.Export(c, mongo.ScenarioCollection, request)
	if err != nil {
		panic(err)
	}

	dbexport.AttachFile(c, mongo.ScenarioCollection, b)
}

// ValidateTemplates
// @Param body body TemplateRequest true "body"
// @Success 200 {object} template.ValidateResponse
func (a *api) ValidateTemplates(c *gin.Context) {
	var request TemplateRequest
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))

		return
	}

	response, err := a.store.ValidateTemplates(c, request)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())

			return
		}

		panic(err)
	}

	c.JSON(http.StatusOK, response)
}

// GetTemplateVars
// @Success 200 {array} TemplateVarsResponse
func (a *api) GetTemplateVars(c *gin.Context) {
	vars, err := a.store.GetTemplateVars(c)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, vars)
}
