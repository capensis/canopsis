package pbehavior

import (
	"context"
	"encoding/json"
	"errors"
	"io"
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
	common.CrudAPI
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

// Find all pbehaviors
// @Summary Find all pbehaviors
// @Description Get paginated list of pbehaviors
// @Tags pbehaviors
// @ID pbehaviors-find-all
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param page query integer true "current page"
// @Param limit query integer true "items per page"
// @Param search query string false "search query"
// @Param sort query string false "sort query"
// @Param sort_by query string false "sort query"
// @Success 200 {object} common.PaginatedListResponse{data=[]Response}
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /pbehaviors [get]
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

// Find pbehaviors by entity id
// @Summary Find pbehaviors by entity id
// @Description Get list of pbehaviors
// @Tags pbehaviors
// @ID pbehaviors-find-by-entity-id
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id query string true "Entity id"
// @Success 200 {array} Response
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /entities/pbehaviors [get]
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

// Get pbehavior by id
// @Summary Get pbehavior by id
// @Description Get pbehavior by id
// @Tags pbehaviors
// @ID pbehaviors-get-by-id
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "pbehavior id"
// @Success 200 {object} Response
// @Failure 404 {object} common.ErrorResponse
// @Router /pbehaviors/{id} [get]
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

// Find entities by pbehavior id
// @Summary Find entities by pbehavior id
// @Description Find entities by pbehavior id
// @Tags pbehaviors
// @ID pbehaviors-find-entities
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "pbehavior id"
// @Param page query integer true "current page"
// @Param limit query integer true "items per page"
// @Param search query string false "search query"
// @Param sort query string false "sort query"
// @Param sort_by query string false "sort query"
// @Success 200 {object} common.PaginatedListResponse{data=[]entity.Entity}
// @Failure 404 {object} common.ErrorResponse
// @Failure 400 {object} common.ErrorResponse
// @Router /pbehaviors/{id}/entities [get]
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

// Create pbehavior
// @Summary Create pbehavior
// @Description Create pbehavior
// @Tags pbehaviors
// @ID pbehaviors-create
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param body body EditRequest true "body"
// @Success 201 {object} Response
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /pbehaviors [post]
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

	err = a.actionLogger.Action(c, logger.LogEntry{
		Action:    logger.ActionCreate,
		ValueType: logger.ValueTypePbehavior,
		ValueID:   model.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	a.sendComputeTask(pbehavior.ComputeTask{
		PbehaviorID:   model.ID,
		OperationType: pbehavior.OperationCreate,
	})
	c.JSON(http.StatusCreated, model)
}

// Update behavior by id
// @Summary Update behavior by id
// @Description Update behavior by id
// @Tags pbehaviors
// @ID pbehaviors-update-by-id
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "pbehavior id"
// @Param body body EditRequest true "body"
// @Success 200 {object} Response
// @Failure 400 {object} common.ValidationErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /pbehaviors/{id} [put]
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

	err = a.actionLogger.Action(c, logger.LogEntry{
		Action:    logger.ActionUpdate,
		ValueType: logger.ValueTypePbehavior,
		ValueID:   model.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	a.sendComputeTask(pbehavior.ComputeTask{
		PbehaviorID:   model.ID,
		OperationType: pbehavior.OperationUpdate,
	})

	c.JSON(http.StatusOK, model)
}

// Patch partial set of behavior's attributes by id
// @Summary Patch partial set of behavior attributes by id
// @Description Patch partial set of behavior attributes by id
// @Tags pbehaviors
// @ID pbehaviors-patch-by-id
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "pbehavior id"
// @Param body body PatchRequest true "body"
// @Success 200 {object} Response
// @Failure 400 {object} common.ValidationErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /pbehaviors/{id} [patch]
func (a *api) Patch(c *gin.Context) {
	req := PatchRequest{}
	var err error
	var reqBody []byte
	reqBody, err = io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	err = json.Unmarshal(reqBody, &req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, req))
		return
	}

	keys, err := a.getProvidedKeys(reqBody)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, req))
		return
	}

	var model *Response
	if a.hasValidationFields(keys, []string{"tstart", "tstop", "type"}) {
		// Patching fields having constraint validation will retry
		// until snapshot is matching or retry count reached
		var updated = false
		retried := 0
		for !updated && retried < 5 {
			model, err = a.store.GetOneBy(c.Request.Context(), bson.M{"_id": c.Param("id")})
			if err != nil || model == nil {
				c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
				return
			}
			snapshot := bson.M{
				"_id":    model.ID,
				"tstart": model.Start,
				"tstop":  model.Stop,
				"type_":  model.Type.ID,
			}

			// Clear tstop field when tstop is defined as null in request body
			if a.hasValidationFields(keys, []string{"tstop"}) {
				model.Stop = nil
			}

			err = a.transformer.Patch(c.Request.Context(), req, model)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
				return
			}

			// Validation
			if model.Type.Type != pbehavior.TypePause && (model.Stop == nil ||
				!model.Stop.After(model.Start.Time)) {
				c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(errors.New("invalid fields start, stop, type")))
				return
			}
			updated, err = a.store.UpdateByFilter(c.Request.Context(), model, snapshot)
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

	err = a.actionLogger.Action(c, logger.LogEntry{
		Action:    logger.ActionUpdate,
		ValueType: logger.ValueTypePbehavior,
		ValueID:   model.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	a.sendComputeTask(pbehavior.ComputeTask{
		PbehaviorID:   model.ID,
		OperationType: pbehavior.OperationUpdate,
	})

	c.JSON(http.StatusOK, model)
}

// Delete pbehavior by id
// @Summary Delete pbehavior by id
// @Description Delete pbehavior by id
// @Tags pbehaviors
// @ID pbehavior-delete-by-id
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "pbehavior id"
// @Success 204
// @Failure 404 {object} common.ErrorResponse
// @Router /pbehaviors/{id} [delete]
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

	err = a.actionLogger.Action(c, logger.LogEntry{
		Action:    logger.ActionDelete,
		ValueType: logger.ValueTypePbehavior,
		ValueID:   c.Param("id"),
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	a.sendComputeTask(pbehavior.ComputeTask{
		PbehaviorID:   id,
		OperationType: pbehavior.OperationDelete,
	})
	c.JSON(http.StatusNoContent, nil)
}

// Count entities matching filter
// @Summary Count entities matching filter
// @Description Count entities matching filter
// @Tags pbehaviors
// @ID pbehavior-countfilter
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param body body FilterRequest true "body"
// @Success 200 {object} CountFilterResult
// @Failure 400 {object} common.ErrorResponse
// @Failure 408 {object} common.ErrorResponse
// @Router /pbehaviors/count [post]
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
			Str("pbehavior", task.PbehaviorID).
			Msg("fail to start pbehavior recompute")
	}
}

func (a *api) getProvidedKeys(data []byte) ([]string, error) {
	m := make(map[string]interface{})
	err := json.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}

	result := make([]string, 0)
	for key := range m {
		result = append(result, key)
	}
	return result, nil
}

func (a *api) hasValidationFields(keys []string, fields []string) bool {
	for _, key := range keys {
		for _, field := range fields {
			if key == field {
				return true
			}
		}
	}
	return false
}
