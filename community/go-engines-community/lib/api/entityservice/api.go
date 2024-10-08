package entityservice

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/rs/zerolog"
	"github.com/valyala/fastjson"
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
	actionLogger      logger.ActionLogger
	logger            zerolog.Logger

	serviceChangeListener chan<- entityservice.ChangeEntityMessage
}

func NewApi(
	store Store,
	serviceChangeListener chan<- entityservice.ChangeEntityMessage,
	metricMetaUpdater metrics.MetaUpdater,
	transformer common.PatternFieldsTransformer,
	actionLogger logger.ActionLogger,
	logger zerolog.Logger,
) API {
	return &api{
		store:                 store,
		serviceChangeListener: serviceChangeListener,
		metricMetaUpdater:     metricMetaUpdater,
		transformer:           transformer,
		actionLogger:          actionLogger,
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
// @Success 200 {object} common.PaginatedListResponse{data=[]entity.Entity}
func (a *api) GetDependencies(c *gin.Context) {
	var r ContextGraphRequest
	r.Query = pagination.GetDefaultQuery()

	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))
		return
	}

	aggregationResult, err := a.store.GetDependencies(c, r)
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
// @Success 200 {object} common.PaginatedListResponse{data=[]entity.Entity}
func (a *api) GetImpacts(c *gin.Context) {
	var r ContextGraphRequest
	r.Query = pagination.GetDefaultQuery()

	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))
		return
	}

	aggregationResult, err := a.store.GetImpacts(c, r)
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
		panic(err)
	}

	if service.Enabled {
		a.sendChangeMsg(entityservice.ChangeEntityMessage{
			ID:                      service.ID,
			EntityType:              service.Type,
			IsServicePatternChanged: true,
		})
	}

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    logger.ActionCreate,
		ValueType: logger.ValueTypeEntityService,
		ValueID:   service.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
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

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    logger.ActionUpdate,
		ValueType: logger.ValueTypeEntityService,
		ValueID:   service.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	a.metricMetaUpdater.UpdateById(c, service.ID)

	c.JSON(http.StatusOK, service)
}

func (a *api) Delete(c *gin.Context) {
	id := c.Param("id")
	ok, err := a.store.Delete(c, id)

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

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    logger.ActionDelete,
		ValueType: logger.ValueTypeEntityService,
		ValueID:   id,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	a.metricMetaUpdater.DeleteById(c, id)

	c.Status(http.StatusNoContent)
}

// BulkCreate
// @Param body body []CreateRequest true "body"
func (a *api) BulkCreate(c *gin.Context) {
	userId := c.MustGet(auth.UserKey).(string)

	var ar fastjson.Arena

	raw, err := c.GetRawData()
	if err != nil {
		panic(err)
	}

	jsonValue, err := fastjson.ParseBytes(raw)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	rawObjects, err := jsonValue.Array()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	response := ar.NewArray()
	serviceIDs := make([]string, 0, len(rawObjects))

	for idx, rawObject := range rawObjects {
		object, err := rawObject.Object()
		if err != nil {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusBadRequest, rawObject, ar.NewString(err.Error())))
			continue
		}

		var request CreateRequest
		err = json.Unmarshal(object.MarshalTo(nil), &request)
		if err != nil {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusBadRequest, rawObject, ar.NewString(err.Error())))
			continue
		}

		err = binding.Validator.ValidateStruct(request)
		if err != nil {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusBadRequest, rawObject, common.NewValidationErrorFastJsonValue(&ar, err, request)))
			continue
		}

		err = a.transformEditRequest(c, &request.EditRequest)
		if err != nil {
			valErr := common.ValidationError{}
			if errors.As(err, &valErr) {
				response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusBadRequest, rawObject, common.NewValidationErrorFastJsonValue(&ar, valErr, request)))
				continue
			}

			a.logger.Err(err).Msg("cannot create entity service")
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusInternalServerError, rawObject, ar.NewString(common.InternalServerErrorResponse.Error)))
			continue
		}

		service, err := a.store.Create(c, request)
		if err != nil {
			a.logger.Err(err).Msg("cannot create entity service")
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusInternalServerError, rawObject, ar.NewString(common.InternalServerErrorResponse.Error)))
			continue
		}

		if service.Enabled {
			a.sendChangeMsg(entityservice.ChangeEntityMessage{
				ID:                      service.ID,
				EntityType:              service.Type,
				IsServicePatternChanged: true,
			})
		}

		response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, service.ID, http.StatusOK, rawObject, nil))

		err = a.actionLogger.Action(context.Background(), userId, logger.LogEntry{
			Action:    logger.ActionCreate,
			ValueType: logger.ValueTypeEntityService,
			ValueID:   service.ID,
		})
		if err != nil {
			a.actionLogger.Err(err, "failed to log action")
		}

		serviceIDs = append(serviceIDs, service.ID)
	}

	a.metricMetaUpdater.UpdateById(c, serviceIDs...)

	c.Data(http.StatusMultiStatus, gin.MIMEJSON, response.MarshalTo(nil))
}

// BulkUpdate
// @Param body body []BulkUpdateRequestItem true "body"
func (a *api) BulkUpdate(c *gin.Context) {
	userId := c.MustGet(auth.UserKey).(string)

	var ar fastjson.Arena

	raw, err := c.GetRawData()
	if err != nil {
		panic(err)
	}

	jsonValue, err := fastjson.ParseBytes(raw)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	rawObjects, err := jsonValue.Array()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	response := ar.NewArray()
	serviceIDs := make([]string, 0, len(rawObjects))

	for idx, rawObject := range rawObjects {
		object, err := rawObject.Object()
		if err != nil {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusBadRequest, rawObject, ar.NewString(err.Error())))
			continue
		}

		var request BulkUpdateRequestItem
		err = json.Unmarshal(object.MarshalTo(nil), &request)
		if err != nil {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusBadRequest, rawObject, ar.NewString(err.Error())))
			continue
		}

		err = binding.Validator.ValidateStruct(request)
		if err != nil {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusBadRequest, rawObject, common.NewValidationErrorFastJsonValue(&ar, err, request)))
			continue
		}

		err = a.transformEditRequest(c, &request.EditRequest)
		if err != nil {
			valErr := common.ValidationError{}
			if errors.As(err, &valErr) {
				response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusBadRequest, rawObject, common.NewValidationErrorFastJsonValue(&ar, valErr, request)))
				continue
			}

			a.logger.Err(err).Msg("cannot update entity service")
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusInternalServerError, rawObject, ar.NewString(common.InternalServerErrorResponse.Error)))
			continue
		}

		service, serviceChanges, err := a.store.Update(c, UpdateRequest(request))
		if err != nil {
			a.logger.Err(err).Msg("cannot update entity service")
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusInternalServerError, rawObject, ar.NewString(common.InternalServerErrorResponse.Error)))
			continue
		}

		if service == nil {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusNotFound, rawObject, ar.NewString(common.NotFoundResponse.Error)))
			continue
		}

		if service.Enabled || serviceChanges.IsToggled {
			a.sendChangeMsg(entityservice.ChangeEntityMessage{
				ID:                      service.ID,
				EntityType:              service.Type,
				IsServicePatternChanged: serviceChanges.IsPatternChanged,
				IsToggled:               serviceChanges.IsToggled,
			})
		}

		response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, service.ID, http.StatusOK, rawObject, nil))

		err = a.actionLogger.Action(context.Background(), userId, logger.LogEntry{
			Action:    logger.ActionUpdate,
			ValueType: logger.ValueTypeEntityService,
			ValueID:   service.ID,
		})
		if err != nil {
			a.actionLogger.Err(err, "failed to log action")
		}

		serviceIDs = append(serviceIDs, service.ID)
	}

	a.metricMetaUpdater.UpdateById(c, serviceIDs...)

	c.Data(http.StatusMultiStatus, gin.MIMEJSON, response.MarshalTo(nil))
}

// BulkDelete
// @Param body body []BulkDeleteRequestItem true "body"
func (a *api) BulkDelete(c *gin.Context) {
	userId := c.MustGet(auth.UserKey).(string)

	var ar fastjson.Arena

	raw, err := c.GetRawData()
	if err != nil {
		panic(err)
	}

	jsonValue, err := fastjson.ParseBytes(raw)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	rawObjects, err := jsonValue.Array()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	response := ar.NewArray()
	serviceIDs := make([]string, 0, len(rawObjects))

	for idx, rawObject := range rawObjects {
		object, err := rawObject.Object()
		if err != nil {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusBadRequest, rawObject, ar.NewString(err.Error())))
			continue
		}

		var request BulkDeleteRequestItem
		err = json.Unmarshal(object.MarshalTo(nil), &request)
		if err != nil {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusBadRequest, rawObject, ar.NewString(err.Error())))
			continue
		}

		err = binding.Validator.ValidateStruct(request)
		if err != nil {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusBadRequest, rawObject, common.NewValidationErrorFastJsonValue(&ar, err, request)))
			continue
		}

		ok, err := a.store.Delete(c, request.ID)
		if err != nil {
			a.logger.Err(err).Msg("cannot delete entity service")
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusInternalServerError, rawObject, ar.NewString(common.InternalServerErrorResponse.Error)))
			continue
		}

		if !ok {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusNotFound, rawObject, ar.NewString(common.NotFoundResponse.Error)))
			continue
		}

		a.sendChangeMsg(entityservice.ChangeEntityMessage{
			ID:                      request.ID,
			EntityType:              types.EntityTypeService,
			IsServicePatternChanged: true,
		})

		response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, request.ID, http.StatusOK, rawObject, nil))

		err = a.actionLogger.Action(context.Background(), userId, logger.LogEntry{
			Action:    logger.ActionDelete,
			ValueType: logger.ValueTypeEntityService,
			ValueID:   request.ID,
		})
		if err != nil {
			a.actionLogger.Err(err, "failed to log action")
		}

		serviceIDs = append(serviceIDs, request.ID)
	}

	a.metricMetaUpdater.DeleteById(c, serviceIDs...)

	c.Data(http.StatusMultiStatus, gin.MIMEJSON, response.MarshalTo(nil))
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
