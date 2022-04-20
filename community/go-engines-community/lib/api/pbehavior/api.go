package pbehavior

import (
	"context"
	"encoding/json"
	"errors"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"github.com/gin-gonic/gin/binding"
	"github.com/valyala/fastjson"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pbehaviorexception"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
)

type API interface {
	common.BulkCrudAPI
	Patch(c *gin.Context)
	ListByEntityID(c *gin.Context)
	ListEntities(c *gin.Context)
	CountFilter(c *gin.Context)
}

type api struct {
	transformer  ModelTransformer
	store        Store
	computeChan  chan<- pbehavior.ComputeTask
	conf         config.UserInterfaceConfigProvider
	actionLogger logger.ActionLogger
	logger       zerolog.Logger
}

func NewApi(
	transformer ModelTransformer,
	store Store,
	computeChan chan<- pbehavior.ComputeTask,
	conf config.UserInterfaceConfigProvider,
	actionLogger logger.ActionLogger,
	logger zerolog.Logger,
) API {
	return &api{
		transformer:  transformer,
		store:        store,
		computeChan:  computeChan,
		conf:         conf,
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

	res, err := a.store.FindByEntityID(c.Request.Context(), r.ID)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, res)
}

// Get
// @Success 200 {object} Response
func (a *api) Get(c *gin.Context) {
	pbh, err := a.store.GetOneBy(c.Request.Context(), bson.M{"_id": c.Param("id")})
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

	model, err := a.transformer.TransformCreateRequestToModel(c.Request.Context(), request)
	if err != nil {
		if err == ErrReasonNotExists || err == ErrExceptionNotExists || err == pbehaviorexception.ErrTypeNotExists {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		} else {
			panic(err)
		}

		return
	}

	err = a.store.Insert(c.Request.Context(), model)
	if err != nil {
		panic(err)
	}

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    logger.ActionCreate,
		ValueType: logger.ValueTypePbehavior,
		ValueID:   model.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	a.sendComputeTask(pbehavior.ComputeTask{
		PbehaviorIds: []string{model.ID},
	})

	c.JSON(http.StatusCreated, model)
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

	model, err := a.transformer.TransformUpdateRequestToModel(c.Request.Context(), request)
	if err != nil {
		if err == ErrReasonNotExists || err == ErrExceptionNotExists || err == pbehaviorexception.ErrTypeNotExists {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		} else {
			panic(err)
		}

		return
	}

	ok, err := a.store.Update(c.Request.Context(), model)
	if err != nil {
		panic(err)
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    logger.ActionUpdate,
		ValueType: logger.ValueTypePbehavior,
		ValueID:   model.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	a.sendComputeTask(pbehavior.ComputeTask{
		PbehaviorIds: []string{model.ID},
	})

	c.JSON(http.StatusOK, model)
}

// Patch
// @Param body body PatchRequest true "body"
// @Success 200 {object} Response
func (a *api) Patch(c *gin.Context) {
	req := PatchRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, req))
		return
	}
	ctx := c.Request.Context()

	var model *Response
	if req.Start != nil || req.Stop.isSet || req.Type != nil {
		// Patching fields having constraint validation will retry
		// until snapshot is matching or retry count reached
		updated := false
		retried := 0
		for !updated && retried < 5 {
			model, err = a.store.GetOneBy(ctx, bson.M{"_id": c.Param("id")})
			if err != nil || model == nil {
				c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
				return
			}
			if model.Stop != nil && model.Stop.IsZero() {
				model.Stop = nil
			}
			snapshot := bson.M{
				"_id":    model.ID,
				"tstart": model.Start,
				"tstop":  model.Stop,
				"type_":  model.Type.ID,
			}

			// Clear tstop field when tstop is defined as null in request body
			if req.Stop.isSet && req.Stop.CpsTime == nil {
				model.Stop = nil
			}

			err = a.transformer.Patch(ctx, req, model)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
				return
			}

			// Validation
			if model.Type.Type != pbehavior.TypePause && model.Stop == nil ||
				(model.Stop != nil && model.Stop.Before(*model.Start)) {
				c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(errors.New("invalid fields start, stop, type")))
				return
			}
			updated, err = a.store.UpdateByFilter(ctx, model, snapshot)
			if err != nil {
				panic(err)
			}
			if updated {
				break
			}
			retried++
		}

		if !updated {
			c.AbortWithStatusJSON(http.StatusNotFound, common.NewErrorResponse(errors.New("max update retry reached")))
			return
		}
	} else {
		// Patching fields that doesn't need to be validated will be executed once
		var ok bool
		model, err = a.store.GetOneBy(c.Request.Context(), bson.M{"_id": c.Param("id")})
		if err != nil || model == nil {
			c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
			return
		}

		err = a.transformer.Patch(c.Request.Context(), req, model)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
			return
		}

		ok, err = a.store.Update(c.Request.Context(), model)
		if err != nil {
			panic(err)
		}

		if !ok {
			c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
			return
		}
	}

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    logger.ActionUpdate,
		ValueType: logger.ValueTypePbehavior,
		ValueID:   model.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	a.sendComputeTask(pbehavior.ComputeTask{
		PbehaviorIds: []string{model.ID},
	})

	c.JSON(http.StatusOK, model)
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

// BulkCreate
// @Param body body []CreateRequest true "body"
// @Success 207 {array} []BulkCreateResponseItem
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

		model, err := a.transformer.TransformCreateRequestToModel(ctx, request)
		if err != nil {
			if err == ErrReasonNotExists || err == ErrExceptionNotExists || err == pbehaviorexception.ErrTypeNotExists {
				response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusBadRequest, rawObject, ar.NewString(err.Error())))
				continue
			}

			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusInternalServerError, rawObject, ar.NewString(err.Error())))
			continue
		}

		err = a.store.Insert(ctx, model)
		if err != nil {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusInternalServerError, rawObject, ar.NewString(err.Error())))
			continue
		}

		response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, model.ID, http.StatusOK, rawObject, nil))

		err = a.actionLogger.Action(context.Background(), userId, logger.LogEntry{
			Action:    logger.ActionCreate,
			ValueType: logger.ValueTypePbehavior,
			ValueID:   model.ID,
		})
		if err != nil {
			a.actionLogger.Err(err, "failed to log action")
		}

		ids = append(ids, model.ID)
	}

	a.sendComputeTask(pbehavior.ComputeTask{
		PbehaviorIds: ids,
	})

	c.Data(http.StatusMultiStatus, gin.MIMEJSON, response.MarshalTo(nil))
}

// BulkUpdate
// @Param body body []BulkUpdateRequestItem true "body"
// @Success 207 {array} []BulkUpdateResponseItem
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

		model, err := a.transformer.TransformUpdateRequestToModel(ctx, UpdateRequest(request))
		if err != nil {
			if err == ErrReasonNotExists || err == ErrExceptionNotExists || err == pbehaviorexception.ErrTypeNotExists {
				response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusBadRequest, rawObject, ar.NewString(err.Error())))
				continue
			}

			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusInternalServerError, rawObject, ar.NewString(err.Error())))
			continue
		}

		ok, err := a.store.Update(ctx, model)
		if err != nil {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusInternalServerError, rawObject, ar.NewString(err.Error())))
			continue
		}

		if !ok {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusNotFound, rawObject, ar.NewString("Not found")))
			continue
		}

		response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, model.ID, http.StatusOK, rawObject, nil))

		err = a.actionLogger.Action(context.Background(), userId, logger.LogEntry{
			Action:    logger.ActionUpdate,
			ValueType: logger.ValueTypePbehavior,
			ValueID:   model.ID,
		})
		if err != nil {
			a.actionLogger.Err(err, "failed to log action")
		}

		ids = append(ids, model.ID)
	}

	a.sendComputeTask(pbehavior.ComputeTask{
		PbehaviorIds: ids,
	})

	c.Data(http.StatusMultiStatus, gin.MIMEJSON, response.MarshalTo(nil))
}

// BulkDelete
// @Param body body []BulkDeleteRequestItem true "body"
// @Success 207 {array} []BulkDeleteResponseItem
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
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusBadRequest, rawObject, ar.NewString(err.Error())))
			continue
		}

		if !ok {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusNotFound, rawObject, ar.NewString("Not found")))
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

// CountFilter
// @Param body body FilterRequest true "body"
// @Success 200 {object} CountFilterResult
func (a api) CountFilter(c *gin.Context) {
	var request FilterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))

		return
	}

	data, err := a.store.Count(c.Request.Context(), NewFilter(request.Filter), a.conf.Get().CheckCountRequestTimeout)
	if errors.Is(err, context.DeadlineExceeded) {
		c.AbortWithStatusJSON(http.StatusRequestTimeout, common.ErrTimeoutResponse)
		return
	} else if err != nil {
		panic(err)
	}
	data.OverLimit = int(data.GetTotal()) > a.conf.Get().MaxMatchedItems

	c.JSON(http.StatusOK, data)
}

func (a *api) sendComputeTask(task pbehavior.ComputeTask) {
	select {
	case a.computeChan <- task:
	default:
		a.logger.Err(errors.New("channel is full")).
			Strs("pbehavior", task.PbehaviorIds).
			Msg("fail to start pbehavior recompute")
	}
}
