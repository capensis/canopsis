package pbehavior

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/rs/zerolog"
	"github.com/valyala/fastjson"
)

type API interface {
	common.BulkCrudAPI
	Patch(c *gin.Context)
	DeleteByName(c *gin.Context)
	ListByEntityID(c *gin.Context)
	CalendarByEntityID(c *gin.Context)
	ListEntities(c *gin.Context)
	BulkEntityCreate(c *gin.Context)
	BulkEntityDelete(c *gin.Context)
}

type api struct {
	store        Store
	computeChan  chan<- pbehavior.ComputeTask
	actionLogger logger.ActionLogger
	logger       zerolog.Logger

	transformer common.PatternFieldsTransformer
}

func NewApi(
	store Store,
	computeChan chan<- pbehavior.ComputeTask,
	transformer common.PatternFieldsTransformer,
	actionLogger logger.ActionLogger,
	logger zerolog.Logger,
) API {
	return &api{
		store:        store,
		computeChan:  computeChan,
		transformer:  transformer,
		actionLogger: actionLogger,
		logger:       logger,
	}
}

// List
// @Success 200 {object} common.PaginatedListResponse{data=[]Response}
func (a *api) List(c *gin.Context) {
	var r ListRequest
	r.Query = pagination.GetDefaultQuery()

	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))

		return
	}

	aggregationResult, err := a.store.Find(c.Request.Context(), r)
	if err != nil {
		panic(err)
	}

	res, err := common.NewPaginatedResponse(r.Query, aggregationResult)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, res)
}

// ListByEntityID
// @Success 200 {array} Response
func (a *api) ListByEntityID(c *gin.Context) {
	var r FindByEntityIDRequest
	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))

		return
	}

	entity, err := a.store.FindEntity(c.Request.Context(), r.ID)
	if err != nil {
		panic(err)
	}
	if entity == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	res, err := a.store.FindByEntityID(c.Request.Context(), *entity, r)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, res)
}

// CalendarByEntityID
// @Success 200 {array} CalendarResponse
func (a *api) CalendarByEntityID(c *gin.Context) {
	var r CalendarByEntityIDRequest
	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))

		return
	}

	entity, err := a.store.FindEntity(c.Request.Context(), r.ID)
	if err != nil {
		panic(err)
	}
	if entity == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	res, err := a.store.CalendarByEntityID(c.Request.Context(), *entity, r)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, res)
}

// Get
// @Success 200 {object} Response
func (a *api) Get(c *gin.Context) {
	pbh, err := a.store.GetOneBy(c.Request.Context(), c.Param("id"))
	if err != nil {
		panic(err)
	}

	if pbh == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, pbh)
}

// ListEntities
// @Success 200 {object} common.PaginatedListResponse{data=[]entity.Entity}
func (a *api) ListEntities(c *gin.Context) {
	var r EntitiesListRequest
	r.Query = pagination.GetDefaultQuery()

	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))

		return
	}

	aggregationResult, err := a.store.FindEntities(c.Request.Context(), c.Param("id"), r)
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

// Create
// @Param body body CreateRequest true "body"
// @Success 201 {object} Response
func (a *api) Create(c *gin.Context) {
	var request CreateRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))

		return
	}

	err := a.transformEditRequest(c.Request.Context(), &request.EditRequest)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}
		panic(err)
	}

	pbh, err := a.store.Insert(c.Request.Context(), request)
	if err != nil {
		validationErr := common.ValidationError{}
		if errors.As(err, &validationErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, validationErr.ValidationErrorResponse())
			return
		}
		panic(err)
	}

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    logger.ActionCreate,
		ValueType: logger.ValueTypePbehavior,
		ValueID:   pbh.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	a.sendComputeTask(pbehavior.ComputeTask{
		PbehaviorIds: []string{pbh.ID},
	})

	c.JSON(http.StatusCreated, pbh)
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

	err := a.transformEditRequest(c.Request.Context(), &request.EditRequest)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}
		panic(err)
	}

	pbh, err := a.store.Update(c.Request.Context(), request)
	if err != nil {
		validationErr := common.ValidationError{}
		if errors.As(err, &validationErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, validationErr.ValidationErrorResponse())
			return
		}
		panic(err)
	}

	if pbh == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    logger.ActionUpdate,
		ValueType: logger.ValueTypePbehavior,
		ValueID:   pbh.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	a.sendComputeTask(pbehavior.ComputeTask{
		PbehaviorIds: []string{pbh.ID},
	})

	c.JSON(http.StatusOK, pbh)
}

// Patch
// @Param body body PatchRequest true "body"
// @Success 200 {object} Response
func (a *api) Patch(c *gin.Context) {
	request := PatchRequest{
		ID: c.Param("id"),
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	if request.CorporateEntityPattern != nil {
		r, err := a.transformer.TransformEntityPatternFieldsRequest(c.Request.Context(), common.EntityPatternFieldsRequest{
			CorporateEntityPattern: *request.CorporateEntityPattern,
		})
		if err != nil {
			valErr := common.ValidationError{}
			if errors.As(err, &valErr) {
				c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
				return
			}
			panic(err)
		}
		if r.CorporatePattern.ID != "" {
			request.CorporatePattern = &r.CorporatePattern
		}
	}

	pbh, err := a.store.UpdateByPatch(c.Request.Context(), request)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}
		panic(err)
	}
	if pbh == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    logger.ActionUpdate,
		ValueType: logger.ValueTypePbehavior,
		ValueID:   pbh.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	a.sendComputeTask(pbehavior.ComputeTask{
		PbehaviorIds: []string{pbh.ID},
	})

	c.JSON(http.StatusOK, pbh)
}

func (a *api) Delete(c *gin.Context) {
	id := c.Param("id")
	ok, err := a.store.Delete(c.Request.Context(), id)
	if err != nil {
		panic(err)
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    logger.ActionDelete,
		ValueType: logger.ValueTypePbehavior,
		ValueID:   id,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	a.sendComputeTask(pbehavior.ComputeTask{
		PbehaviorIds: []string{id},
	})
	c.JSON(http.StatusNoContent, nil)
}

func (a *api) DeleteByName(c *gin.Context) {
	request := DeleteByNameRequest{}

	if err := c.ShouldBindQuery(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))

		return
	}

	id, err := a.store.DeleteByName(c.Request.Context(), request.Name)
	if err != nil {
		panic(err)
	}

	if id == "" {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    logger.ActionDelete,
		ValueType: logger.ValueTypePbehavior,
		ValueID:   id,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	a.sendComputeTask(pbehavior.ComputeTask{
		PbehaviorIds: []string{id},
	})
	c.JSON(http.StatusNoContent, nil)
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

	ctx := c.Request.Context()
	response := ar.NewArray()
	ids := make([]string, 0, len(rawObjects))

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

		err = a.transformEditRequest(ctx, &request.EditRequest)
		if err != nil {
			valErr := common.ValidationError{}
			if errors.As(err, &valErr) {
				response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusBadRequest, rawObject, common.NewValidationErrorFastJsonValue(&ar, valErr, request)))
				continue
			}

			a.logger.Err(err).Msg("cannot create pbehavior")
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusInternalServerError, rawObject, ar.NewString(common.InternalServerErrorResponse.Error)))
			continue
		}

		pbh, err := a.store.Insert(ctx, request)
		if err != nil {
			valErr := common.ValidationError{}
			if errors.As(err, &valErr) {
				response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusBadRequest, rawObject, common.NewValidationErrorFastJsonValue(&ar, valErr, request)))
				continue
			}

			a.logger.Err(err).Msg("cannot create pbehavior")
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusInternalServerError, rawObject, ar.NewString(common.InternalServerErrorResponse.Error)))
			continue
		}

		response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, pbh.ID, http.StatusOK, rawObject, nil))

		err = a.actionLogger.Action(context.Background(), userId, logger.LogEntry{
			Action:    logger.ActionCreate,
			ValueType: logger.ValueTypePbehavior,
			ValueID:   pbh.ID,
		})
		if err != nil {
			a.actionLogger.Err(err, "failed to log action")
		}

		ids = append(ids, pbh.ID)
	}

	a.sendComputeTask(pbehavior.ComputeTask{
		PbehaviorIds: ids,
	})

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

	ctx := c.Request.Context()
	response := ar.NewArray()
	ids := make([]string, 0, len(rawObjects))

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

		err = a.transformEditRequest(ctx, &request.EditRequest)
		if err != nil {
			valErr := common.ValidationError{}
			if errors.As(err, &valErr) {
				response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusBadRequest, rawObject, common.NewValidationErrorFastJsonValue(&ar, valErr, request)))
				continue
			}

			a.logger.Err(err).Msg("cannot update pbehavior")
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusInternalServerError, rawObject, ar.NewString(common.InternalServerErrorResponse.Error)))
			continue
		}

		pbh, err := a.store.Update(ctx, UpdateRequest(request))
		if err != nil {
			valErr := common.ValidationError{}
			if errors.As(err, &valErr) {
				response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusBadRequest, rawObject, common.NewValidationErrorFastJsonValue(&ar, valErr, request)))
				continue
			}

			a.logger.Err(err).Msg("cannot update pbehavior")
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusInternalServerError, rawObject, ar.NewString(common.InternalServerErrorResponse.Error)))
			continue
		}

		if pbh == nil {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusNotFound, rawObject, ar.NewString(common.NotFoundResponse.Error)))
			continue
		}

		response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, pbh.ID, http.StatusOK, rawObject, nil))

		err = a.actionLogger.Action(context.Background(), userId, logger.LogEntry{
			Action:    logger.ActionUpdate,
			ValueType: logger.ValueTypePbehavior,
			ValueID:   pbh.ID,
		})
		if err != nil {
			a.actionLogger.Err(err, "failed to log action")
		}

		ids = append(ids, pbh.ID)
	}

	a.sendComputeTask(pbehavior.ComputeTask{
		PbehaviorIds: ids,
	})

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

	ctx := c.Request.Context()
	response := ar.NewArray()
	ids := make([]string, 0, len(rawObjects))

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

		ok, err := a.store.Delete(ctx, request.ID)
		if err != nil {
			a.logger.Err(err).Msg("cannot delete pbehavior")
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusInternalServerError, rawObject, ar.NewString(common.InternalServerErrorResponse.Error)))
			continue
		}

		if !ok {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusNotFound, rawObject, ar.NewString(common.NotFoundResponse.Error)))
			continue
		}

		response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, request.ID, http.StatusOK, rawObject, nil))

		err = a.actionLogger.Action(context.Background(), userId, logger.LogEntry{
			Action:    logger.ActionDelete,
			ValueType: logger.ValueTypePbehavior,
			ValueID:   request.ID,
		})
		if err != nil {
			a.actionLogger.Err(err, "failed to log action")
		}

		ids = append(ids, request.ID)
	}

	a.sendComputeTask(pbehavior.ComputeTask{
		PbehaviorIds: ids,
	})

	c.Data(http.StatusMultiStatus, gin.MIMEJSON, response.MarshalTo(nil))
}

// BulkEntityCreate
// @Param body body []BulkEntityCreateRequestItem true "body"
func (a *api) BulkEntityCreate(c *gin.Context) {
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

	ctx := c.Request.Context()
	response := ar.NewArray()
	ids := make([]string, 0, len(rawObjects))

	for idx, rawObject := range rawObjects {
		object, err := rawObject.Object()
		if err != nil {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusBadRequest, rawObject, ar.NewString(err.Error())))
			continue
		}

		var request BulkEntityCreateRequestItem
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

		pbh, err := a.store.EntityInsert(ctx, request)
		if err != nil {
			valErr := common.ValidationError{}
			if errors.As(err, &valErr) {
				response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusBadRequest, rawObject, common.NewValidationErrorFastJsonValue(&ar, valErr, request)))
				continue
			}

			a.logger.Err(err).Msg("cannot create pbehavior")
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusInternalServerError, rawObject, ar.NewString(common.InternalServerErrorResponse.Error)))
			continue
		}

		if pbh == nil {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusNotFound, rawObject, ar.NewString(common.NotFoundResponse.Error)))
			continue
		}

		response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, pbh.ID, http.StatusOK, rawObject, nil))

		err = a.actionLogger.Action(context.Background(), userId, logger.LogEntry{
			Action:    logger.ActionCreate,
			ValueType: logger.ValueTypePbehavior,
			ValueID:   pbh.ID,
		})
		if err != nil {
			a.actionLogger.Err(err, "failed to log action")
		}

		ids = append(ids, pbh.ID)
	}

	a.sendComputeTask(pbehavior.ComputeTask{
		PbehaviorIds: ids,
	})

	c.Data(http.StatusMultiStatus, gin.MIMEJSON, response.MarshalTo(nil))
}

// BulkEntityDelete
// @Param body body []BulkEntityDeleteRequestItem true "body"
func (a *api) BulkEntityDelete(c *gin.Context) {
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

	ctx := c.Request.Context()
	response := ar.NewArray()
	ids := make([]string, 0, len(rawObjects))

	for idx, rawObject := range rawObjects {
		object, err := rawObject.Object()
		if err != nil {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusBadRequest, rawObject, ar.NewString(err.Error())))
			continue
		}

		var request BulkEntityDeleteRequestItem
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

		id, err := a.store.EntityDelete(ctx, request)
		if err != nil {
			a.logger.Err(err).Msg("cannot delete pbehavior")
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusInternalServerError, rawObject, ar.NewString(common.InternalServerErrorResponse.Error)))
			continue
		}

		if id == "" {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusNotFound, rawObject, ar.NewString(common.NotFoundResponse.Error)))
			continue
		}

		response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, id, http.StatusOK, rawObject, nil))

		err = a.actionLogger.Action(context.Background(), userId, logger.LogEntry{
			Action:    logger.ActionDelete,
			ValueType: logger.ValueTypePbehavior,
			ValueID:   id,
		})
		if err != nil {
			a.actionLogger.Err(err, "failed to log action")
		}

		ids = append(ids, id)
	}

	a.sendComputeTask(pbehavior.ComputeTask{
		PbehaviorIds: ids,
	})

	c.Data(http.StatusMultiStatus, gin.MIMEJSON, response.MarshalTo(nil))
}

func (a *api) sendComputeTask(task pbehavior.ComputeTask) {
	a.computeChan <- task
}

func (a *api) transformEditRequest(ctx context.Context, request *EditRequest) error {
	var err error
	request.EntityPatternFieldsRequest, err = a.transformer.TransformEntityPatternFieldsRequest(ctx, request.EntityPatternFieldsRequest)
	if err != nil {
		return err
	}

	return nil
}
