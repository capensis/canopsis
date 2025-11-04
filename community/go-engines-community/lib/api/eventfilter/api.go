package eventfilter

import (
	"errors"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/bulk"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/crud"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/dbexport"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
)

type API interface {
	crud.BulkAPI
	ListFailures(c *gin.Context)
	ReadFailures(c *gin.Context)
	DBExport(c *gin.Context)
	ValidateTemplates(c *gin.Context)
	GetTemplateVars(c *gin.Context)
	GetCopyVars(c *gin.Context)
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

// Create
// @Param body body CreateRequest true "body"
// @Success 201 {object} Response
func (a *api) Create(c *gin.Context) {
	var request CreateRequest
	var err error

	if err = c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	eventfilter, err := a.store.Insert(c, request)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())

			return
		}

		panic(err)
	}

	c.JSON(http.StatusCreated, eventfilter)
}

// List
// @Success 200 {object} common.PaginatedListResponse{data=[]Response}
func (a *api) List(c *gin.Context) {
	var query FilteredQuery
	query.Query = pagination.GetDefaultQuery()

	if err := c.ShouldBind(&query); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, query))
		return
	}

	aggregationResult, err := a.store.Find(c, query)
	if err != nil {
		panic(err)
	}

	res, err := common.NewPaginatedResponse(query.Query, aggregationResult)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, res)
}

// Get
// @Success 200 {object} Response
func (a *api) Get(c *gin.Context) {
	evf, err := a.store.GetByID(c, c.Param("id"))

	if errors.Is(err, mongodriver.ErrNoDocuments) || evf == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, evf)
}

// Update
// @Param body body UpdateRequest true "body"
// @Success 200 {object} Response
func (a *api) Update(c *gin.Context) {
	request := UpdateRequest{
		ID: c.Param("id"),
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))

		return
	}

	eventfilter, err := a.store.Update(c, request)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())

			return
		}

		panic(err)
	}

	if eventfilter == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, eventfilter)
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

	c.JSON(http.StatusNoContent, nil)
}

// BulkCreate
// @Param body body []CreateRequest true "body"
func (a *api) BulkCreate(c *gin.Context) {
	bulk.Handler(c, func(request CreateRequest) (string, error) {
		eventfilter, err := a.store.Insert(c, request)
		if err != nil {
			return "", err
		}

		return eventfilter.ID, nil
	}, a.logger)
}

// BulkUpdate
// @Param body body []BulkUpdateRequestItem true "body"
func (a *api) BulkUpdate(c *gin.Context) {
	bulk.Handler(c, func(request BulkUpdateRequestItem) (string, error) {
		eventfilter, err := a.store.Update(c, UpdateRequest(request))
		if err != nil || eventfilter == nil {
			return "", err
		}

		return eventfilter.ID, nil
	}, a.logger)
}

// BulkDelete
// @Param body body []BulkDeleteRequestItem true "body"
func (a *api) BulkDelete(c *gin.Context) {
	userID := c.MustGet(auth.UserKey).(string)

	bulk.Handler(c, func(request BulkDeleteRequestItem) (string, error) {
		ok, err := a.store.Delete(c, request.ID, userID)
		if err != nil || !ok {
			return "", err
		}

		return request.ID, nil
	}, a.logger)
}

// ListFailures
// @Success 200 {object} common.PaginatedListResponse{data=[]FailureResponse}
func (a *api) ListFailures(c *gin.Context) {
	r := FailureRequest{}
	r.Query = pagination.GetDefaultQuery()
	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))
		return
	}

	aggregationResult, err := a.store.FindFailures(c, c.Param("id"), r)
	if err != nil {
		panic(err)
	}

	if aggregationResult == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	res, err := common.NewPaginatedResponse(r.Query, aggregationResult)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, res)
}

func (a *api) ReadFailures(c *gin.Context) {
	exists, err := a.store.ReadFailures(c, c.Param("id"))
	if err != nil {
		panic(err)
	}

	if !exists {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.Status(http.StatusNoContent)
}

// DBExport
// @Param body body dbexport.Request true "body"
func (a *api) DBExport(c *gin.Context) {
	request := dbexport.Request{}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	b, err := a.mongoExporter.Export(c, mongo.EventFilterRuleCollection, request)
	if err != nil {
		panic(err)
	}

	dbexport.AttachFile(c, mongo.EventFilterRuleCollection, b)
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
	c.JSON(http.StatusOK, a.store.GetTemplateVars())
}

// GetCopyVars
// @Success 200 {array} CopyVarsResponse
func (a *api) GetCopyVars(c *gin.Context) {
	c.JSON(http.StatusOK, a.store.GetCopyVars())
}
