package idlerule

import (
	"errors"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/bulk"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/crud"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/dbexport"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type API interface {
	crud.BulkAPI
	DBExport(c *gin.Context)
}

type api struct {
	store          Store
	mongoExporter  dbexport.Exporter
	errorResponder httperror.Responder
	logger         zerolog.Logger
}

func NewApi(
	store Store,
	mongoExporter dbexport.Exporter,
	errorResponder httperror.Responder,
	logger zerolog.Logger,
) API {
	return &api{
		store:          store,
		mongoExporter:  mongoExporter,
		errorResponder: errorResponder,
		logger:         logger,
	}
}

// List
// @Success 200 {object} pagination.ListResponse{data=[]idlerule.Rule}
func (a *api) List(c *gin.Context) {
	var query FilteredQuery
	query.Query = pagination.GetDefaultQuery()

	if err := validation.Bind(c, &query); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	rules, err := a.store.Find(c, query)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	res := pagination.NewResponse(query.Query, rules)
	c.JSON(http.StatusOK, res)
}

// Get
// @Success 200 {object} idlerule.Rule
func (a *api) Get(c *gin.Context) {
	rule, err := a.store.GetOneBy(c, c.Param("id"))
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	if rule == nil {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, rule)
}

// Create
// @Param body body EditRequest true "body"
// @Success 201 {object} idlerule.Rule
func (a *api) Create(c *gin.Context) {
	var request CreateRequest
	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	rule, err := a.store.Insert(c, request)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}

		a.errorResponder.Respond(c, err)

		return
	}

	if rule == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusCreated, rule)
}

// Update
// @Param body body EditRequest true "body"
// @Success 200 {object} Rule
func (a *api) Update(c *gin.Context) {
	request := UpdateRequest{
		ID: c.Param("id"),
	}

	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	rule, err := a.store.Update(c, request)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}

		a.errorResponder.Respond(c, err)

		return
	}

	if rule == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, rule)
}

func (a *api) Delete(c *gin.Context) {
	ok, err := a.store.Delete(c, c.Param("id"), c.MustGet(authctx.UserKey).(string))
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
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
		rule, err := a.store.Insert(c, request)
		if err != nil {
			return "", err
		}

		return rule.ID, nil
	}, a.logger)
}

// BulkUpdate
// @Param body body []BulkUpdateRequestItem true "body"
func (a *api) BulkUpdate(c *gin.Context) {
	bulk.Handler(c, func(request BulkUpdateRequestItem) (string, error) {
		rule, err := a.store.Update(c, UpdateRequest(request))
		if err != nil || rule == nil {
			return "", err
		}

		return rule.ID, nil
	}, a.logger)
}

// BulkDelete
// @Param body body []BulkDeleteRequestItem true "body"
func (a *api) BulkDelete(c *gin.Context) {
	userID := c.MustGet(authctx.UserKey).(string)

	bulk.Handler(c, func(request BulkDeleteRequestItem) (string, error) {
		ok, err := a.store.Delete(c, request.ID, userID)
		if err != nil || !ok {
			return "", err
		}

		return request.ID, nil
	}, a.logger)
}

// DBExport
// @Param body body dbexport.Request true "body"
func (a *api) DBExport(c *gin.Context) {
	request := dbexport.Request{}

	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	b, err := a.mongoExporter.Export(c, mongo.IdleRuleMongoCollection, request)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	dbexport.AttachFile(c, mongo.IdleRuleMongoCollection, b)
}
