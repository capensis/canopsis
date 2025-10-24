package entityservice

import (
	"context"
	"errors"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/bulk"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
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
}

type api struct {
	store             Store
	metricMetaUpdater metrics.MetaUpdater
	transformer       common.PatternFieldsTransformer
	logger            zerolog.Logger

	serviceChangeListener chan<- entityservice.ChangeEntityMessage
}

func NewApi(
	store Store,
	serviceChangeListener chan<- entityservice.ChangeEntityMessage,
	metricMetaUpdater metrics.MetaUpdater,
	transformer common.PatternFieldsTransformer,
	logger zerolog.Logger,
) API {
	return &api{
		store:                 store,
		serviceChangeListener: serviceChangeListener,
		metricMetaUpdater:     metricMetaUpdater,
		transformer:           transformer,
		logger:                logger,
	}
}

// Get
// @Success 200 {object} Response
func (a *api) Get(c *gin.Context) {
	service, err := a.store.GetOneBy(c, c.Param("id"))
	if err != nil {
		panic(err)
	}
	if service == nil {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, service)
}

// GetDependencies
// @Success 200 {object} common.PaginatedListResponse{data=[]ContextGraphEntity}
func (a *api) GetDependencies(c *gin.Context) {
	var r ContextGraphRequest
	r.Query = pagination.GetDefaultQuery()

	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))
		return
	}

	userID := c.MustGet(auth.UserKey).(string)
	aggregationResult, err := a.store.GetDependencies(c, r, userID)
	if err != nil {
		panic(err)
	}

	if aggregationResult == nil {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	res, err := common.NewPaginatedResponse(r.Query, aggregationResult)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, res)
}

// GetImpacts
// @Success 200 {object} common.PaginatedListResponse{data=[]ContextGraphEntity}
func (a *api) GetImpacts(c *gin.Context) {
	var r ContextGraphRequest
	r.Query = pagination.GetDefaultQuery()

	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))
		return
	}

	userID := c.MustGet(auth.UserKey).(string)
	aggregationResult, err := a.store.GetImpacts(c, r, userID)
	if err != nil {
		panic(err)
	}

	if aggregationResult == nil {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	res, err := common.NewPaginatedResponse(r.Query, aggregationResult)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, res)
}

// Create
// @Param body body CreateRequest true "body"
// @Success 201 {object} Response
func (a *api) Create(c *gin.Context) {
	var request CreateRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	err := a.transformEditRequest(c, &request.EditRequest)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}

		panic(err)
	}

	service, err := a.store.Create(c, request)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}

		panic(err)
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
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	err := a.transformEditRequest(c, &request.EditRequest)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}

		panic(err)
	}

	service, serviceChanges, err := a.store.Update(c, request)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}

		panic(err)
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
	ok, err := a.store.Delete(c, id, c.MustGet(auth.UserKey).(string))
	if err != nil {
		panic(err)
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
		err := a.transformEditRequest(c, &request.EditRequest)
		if err != nil {
			return "", err
		}

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
	}, a.logger)
	a.metricMetaUpdater.UpdateById(c, serviceIDs...)
}

// BulkUpdate
// @Param body body []BulkUpdateRequestItem true "body"
func (a *api) BulkUpdate(c *gin.Context) {
	serviceIDs := make([]string, 0)
	bulk.Handler(c, func(request BulkUpdateRequestItem) (string, error) {
		err := a.transformEditRequest(c, &request.EditRequest)
		if err != nil {
			return "", err
		}

		service, serviceChanges, err := a.store.Update(c, UpdateRequest(request))
		if err != nil || service == nil {
			return "", err
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
	}, a.logger)
	a.metricMetaUpdater.UpdateById(c, serviceIDs...)
}

// BulkDelete
// @Param body body []BulkDeleteRequestItem true "body"
func (a *api) BulkDelete(c *gin.Context) {
	userID := c.MustGet(auth.UserKey).(string)

	serviceIDs := make([]string, 0)
	bulk.Handler(c, func(request BulkDeleteRequestItem) (string, error) {
		ok, err := a.store.Delete(c, request.ID, userID)
		if err != nil || !ok {
			return "", err
		}

		a.sendChangeMsg(entityservice.ChangeEntityMessage{
			ID:                      request.ID,
			EntityType:              types.EntityTypeService,
			IsServicePatternChanged: true,
		})

		serviceIDs = append(serviceIDs, request.ID)

		return request.ID, nil
	}, a.logger)

	a.metricMetaUpdater.DeleteById(c, serviceIDs...)
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

func (a *api) transformEditRequest(ctx context.Context, request *EditRequest) error {
	var err error
	request.EntityPatternFieldsRequest, err = a.transformer.TransformEntityPatternFieldsRequest(ctx, request.EntityPatternFieldsRequest)
	if err != nil {
		return err
	}

	return nil
}
