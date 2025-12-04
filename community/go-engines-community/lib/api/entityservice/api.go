package entityservice

import (
	"errors"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/bulk"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type API interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	GetDependencies(c *gin.Context)
	GetImpacts(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	BulkCreate(c *gin.Context)
	BulkUpdate(c *gin.Context)
	BulkDelete(c *gin.Context)
	ValidateTemplates(c *gin.Context)
	GetTemplateVars(c *gin.Context)
}

type api struct {
	store             Store
	metricMetaUpdater metrics.MetaUpdater
	errorResponder    httperror.Responder
	logger            zerolog.Logger

	serviceChangeListener chan<- entityservice.ChangeEntityMessage
}

func NewApi(
	store Store,
	serviceChangeListener chan<- entityservice.ChangeEntityMessage,
	metricMetaUpdater metrics.MetaUpdater,
	errorResponder httperror.Responder,
	logger zerolog.Logger,
) API {
	return &api{
		store:                 store,
		serviceChangeListener: serviceChangeListener,
		metricMetaUpdater:     metricMetaUpdater,
		errorResponder:        errorResponder,
		logger:                logger,
	}
}

// Get
// @Success 200 {object} Response
func (a *api) Get(c *gin.Context) {
	service, err := a.store.GetOneBy(c, c.Param("id"))
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	if service == nil {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, service)
}

// GetDependencies
// @Success 200 {object} pagination.ListResponse{data=[]ContextGraphEntity}
func (a *api) GetDependencies(c *gin.Context) {
	var r ContextGraphRequest
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
	aggregationResult, err := a.store.GetDependencies(c, r, userID)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if aggregationResult == nil {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	res := pagination.NewResponse(r.Query, aggregationResult)
	c.JSON(http.StatusOK, res)
}

// GetImpacts
// @Success 200 {object} pagination.ListResponse{data=[]ContextGraphEntity}
func (a *api) GetImpacts(c *gin.Context) {
	var r ContextGraphRequest
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
	aggregationResult, err := a.store.GetImpacts(c, r, userID)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if aggregationResult == nil {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	res := pagination.NewResponse(r.Query, aggregationResult)
	c.JSON(http.StatusOK, res)
}

// Create
// @Param body body CreateRequest true "body"
// @Success 201 {object} Response
func (a *api) Create(c *gin.Context) {
	var request CreateRequest
	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	service, err := a.store.Create(c, request)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if service.Enabled {
		a.sendChangeMsg(entityservice.ChangeEntityMessage{
			ID:                      service.ID,
			EntityType:              service.Type,
			IsServicePatternChanged: true,
		})
	}

	a.metricMetaUpdater.UpdateById(c, service.ID)

	c.JSON(http.StatusCreated, service)
}

// Update
// @Param body body UpdateRequest true "body"
// @Success 200 {object} Response
func (a *api) Update(c *gin.Context) {
	request := UpdateRequest{
		ID: c.Param("id"),
	}
	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	service, serviceChanges, err := a.store.Update(c, request)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if service == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	if service.Enabled || serviceChanges.IsToggled {
		a.sendChangeMsg(entityservice.ChangeEntityMessage{
			ID:                      service.ID,
			EntityType:              service.Type,
			IsServicePatternChanged: serviceChanges.IsPatternChanged,
			IsToggled:               serviceChanges.IsToggled,
		})
	}

	a.metricMetaUpdater.UpdateById(c, service.ID)

	c.JSON(http.StatusOK, service)
}

func (a *api) Delete(c *gin.Context) {
	id := c.Param("id")
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	ok, err := a.store.Delete(c, id, userID)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	a.sendChangeMsg(entityservice.ChangeEntityMessage{
		ID:         id,
		EntityType: types.EntityTypeService,
		IsDeleted:  true,
	})

	a.metricMetaUpdater.DeleteById(c, id)

	c.Status(http.StatusNoContent)
}

// BulkCreate
// @Param body body []CreateRequest true "body"
func (a *api) BulkCreate(c *gin.Context) {
	serviceIDs := make([]string, 0)
	bulk.Handler(c, func(request CreateRequest) (string, error) {
		service, err := a.store.Create(c, request)
		if err != nil {
			return "", err
		}

		if service.Enabled {
			a.sendChangeMsg(entityservice.ChangeEntityMessage{
				ID:                      service.ID,
				EntityType:              service.Type,
				IsServicePatternChanged: true,
			})
		}

		serviceIDs = append(serviceIDs, service.ID)

		return service.ID, nil
	}, a.errorResponder)
	a.metricMetaUpdater.UpdateById(c, serviceIDs...)
}

// BulkUpdate
// @Param body body []BulkUpdateRequestItem true "body"
func (a *api) BulkUpdate(c *gin.Context) {
	serviceIDs := make([]string, 0)
	bulk.Handler(c, func(request BulkUpdateRequestItem) (string, error) {
		service, serviceChanges, err := a.store.Update(c, UpdateRequest(request))
		if err != nil {
			return "", err
		}

		if service == nil {
			return "", httperror.ErrNotFound
		}

		if service.Enabled || serviceChanges.IsToggled {
			a.sendChangeMsg(entityservice.ChangeEntityMessage{
				ID:                      service.ID,
				EntityType:              service.Type,
				IsServicePatternChanged: serviceChanges.IsPatternChanged,
				IsToggled:               serviceChanges.IsToggled,
			})
		}

		serviceIDs = append(serviceIDs, service.ID)

		return service.ID, nil
	}, a.errorResponder)
	a.metricMetaUpdater.UpdateById(c, serviceIDs...)
}

// BulkDelete
// @Param body body []BulkDeleteRequestItem true "body"
func (a *api) BulkDelete(c *gin.Context) {
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	serviceIDs := make([]string, 0)
	bulk.Handler(c, func(request BulkDeleteRequestItem) (string, error) {
		ok, err := a.store.Delete(c, request.ID, userID)
		if err != nil {
			return "", err
		}

		if !ok {
			return "", httperror.ErrNotFound
		}

		a.sendChangeMsg(entityservice.ChangeEntityMessage{
			ID:                      request.ID,
			EntityType:              types.EntityTypeService,
			IsServicePatternChanged: true,
		})

		serviceIDs = append(serviceIDs, request.ID)

		return request.ID, nil
	}, a.errorResponder)

	a.metricMetaUpdater.DeleteById(c, serviceIDs...)
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

	response, err := a.store.ValidateTemplates(request)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	c.JSON(http.StatusOK, response)
}

// GetTemplateVars
// @Success 200 {array} TemplateVarsResponse
func (a *api) GetTemplateVars(c *gin.Context) {
	c.JSON(http.StatusOK, a.store.GetTemplateVars())
}

func (a *api) sendChangeMsg(msg entityservice.ChangeEntityMessage) {
	select {
	case a.serviceChangeListener <- msg:
	default:
		a.logger.Err(errors.New("channel is full")).
			Str("service_id", msg.ID).
			Msg("fail to send change message")
	}
}
