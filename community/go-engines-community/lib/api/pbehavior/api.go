package pbehavior

import (
	"context"
	"errors"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/bulk"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
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
	BulkConnectorCreate(c *gin.Context)
	BulkConnectorDelete(c *gin.Context)
}

type api struct {
	store        Store
	computeChan  chan<- []string
	actionLogger logger.ActionLogger
	logger       zerolog.Logger

	transformer common.PatternFieldsTransformer
}

func NewApi(
	store Store,
	computeChan chan<- []string,
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

	a.sendComputeTask([]string{pbh.ID})

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

	a.sendComputeTask([]string{pbh.ID})

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

	a.sendComputeTask([]string{pbh.ID})

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

	a.sendComputeTask([]string{id})
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

	a.sendComputeTask([]string{id})
	c.JSON(http.StatusNoContent, nil)
}

// BulkCreate
// @Param body body []CreateRequest true "body"
func (a *api) BulkCreate(c *gin.Context) {
	userId := c.MustGet(auth.UserKey).(string)
	ids := make([]string, 0)
	bulk.Handler(c, func(request CreateRequest) (string, error) {
		err := a.transformEditRequest(c, &request.EditRequest)
		if err != nil {
			return "", err
		}

		pbh, err := a.store.Insert(c, request)
		if err != nil {
			return "", err
		}

		err = a.actionLogger.Action(context.Background(), userId, logger.LogEntry{
			Action:    logger.ActionCreate,
			ValueType: logger.ValueTypePbehavior,
			ValueID:   pbh.ID,
		})
		if err != nil {
			a.actionLogger.Err(err, "failed to log action")
		}

		ids = append(ids, pbh.ID)

		return pbh.ID, nil
	}, a.logger)

	if len(ids) > 0 {
		a.sendComputeTask(ids)
	}
}

// BulkUpdate
// @Param body body []BulkUpdateRequestItem true "body"
func (a *api) BulkUpdate(c *gin.Context) {
	userId := c.MustGet(auth.UserKey).(string)
	ids := make([]string, 0)
	exists := make(map[string]struct{})
	bulk.Handler(c, func(request BulkUpdateRequestItem) (string, error) {
		err := a.transformEditRequest(c, &request.EditRequest)
		if err != nil {
			return "", err
		}

		pbh, err := a.store.Update(c, UpdateRequest(request))
		if err != nil || pbh == nil {
			return "", err
		}

		err = a.actionLogger.Action(context.Background(), userId, logger.LogEntry{
			Action:    logger.ActionUpdate,
			ValueType: logger.ValueTypePbehavior,
			ValueID:   pbh.ID,
		})
		if err != nil {
			a.actionLogger.Err(err, "failed to log action")
		}

		if _, ok := exists[pbh.ID]; !ok {
			ids = append(ids, pbh.ID)
			exists[pbh.ID] = struct{}{}
		}

		return pbh.ID, nil
	}, a.logger)

	if len(ids) > 0 {
		a.sendComputeTask(ids)
	}
}

// BulkDelete
// @Param body body []BulkDeleteRequestItem true "body"
func (a *api) BulkDelete(c *gin.Context) {
	userId := c.MustGet(auth.UserKey).(string)
	ids := make([]string, 0)
	bulk.Handler(c, func(request BulkDeleteRequestItem) (string, error) {
		ok, err := a.store.Delete(c, request.ID)
		if err != nil || !ok {
			return "", err
		}

		err = a.actionLogger.Action(context.Background(), userId, logger.LogEntry{
			Action:    logger.ActionDelete,
			ValueType: logger.ValueTypePbehavior,
			ValueID:   request.ID,
		})
		if err != nil {
			a.actionLogger.Err(err, "failed to log action")
		}

		ids = append(ids, request.ID)

		return request.ID, nil
	}, a.logger)

	if len(ids) > 0 {
		a.sendComputeTask(ids)
	}
}

// BulkEntityCreate
// @Param body body []BulkEntityCreateRequestItem true "body"
func (a *api) BulkEntityCreate(c *gin.Context) {
	userId := c.MustGet(auth.UserKey).(string)
	ids := make([]string, 0)
	bulk.Handler(c, func(request BulkEntityCreateRequestItem) (string, error) {
		pbh, err := a.store.EntityInsert(c, request)
		if err != nil || pbh == nil {
			return "", err
		}

		err = a.actionLogger.Action(context.Background(), userId, logger.LogEntry{
			Action:    logger.ActionCreate,
			ValueType: logger.ValueTypePbehavior,
			ValueID:   pbh.ID,
		})
		if err != nil {
			a.actionLogger.Err(err, "failed to log action")
		}

		ids = append(ids, pbh.ID)

		return pbh.ID, nil
	}, a.logger)

	if len(ids) > 0 {
		a.sendComputeTask(ids)
	}
}

// BulkEntityDelete
// @Param body body []BulkEntityDeleteRequestItem true "body"
func (a *api) BulkEntityDelete(c *gin.Context) {
	userId := c.MustGet(auth.UserKey).(string)
	ids := make([]string, 0)
	bulk.Handler(c, func(request BulkEntityDeleteRequestItem) (string, error) {
		id, err := a.store.EntityDelete(c, request)
		if err != nil || id == "" {
			return "", err
		}

		err = a.actionLogger.Action(context.Background(), userId, logger.LogEntry{
			Action:    logger.ActionDelete,
			ValueType: logger.ValueTypePbehavior,
			ValueID:   id,
		})
		if err != nil {
			a.actionLogger.Err(err, "failed to log action")
		}

		ids = append(ids, id)

		return id, nil
	}, a.logger)

	if len(ids) > 0 {
		a.sendComputeTask(ids)
	}
}

// BulkConnectorCreate
// @Param body body []BulkConnectorCreateRequestItem true "body"
func (a *api) BulkConnectorCreate(c *gin.Context) {
	userId := c.MustGet(auth.UserKey).(string)
	ids := make([]string, 0)
	exists := make(map[string]struct{})
	bulk.Handler(c, func(request BulkConnectorCreateRequestItem) (string, error) {
		pbh, err := a.store.ConnectorCreate(c, request)
		if err != nil || pbh == nil {
			return "", err
		}

		err = a.actionLogger.Action(context.Background(), userId, logger.LogEntry{
			Action:    logger.ActionCreate,
			ValueType: logger.ValueTypePbehavior,
			ValueID:   pbh.ID,
		})
		if err != nil {
			a.actionLogger.Err(err, "failed to log action")
		}

		if _, ok := exists[pbh.ID]; !ok {
			ids = append(ids, pbh.ID)
			exists[pbh.ID] = struct{}{}
		}

		return pbh.ID, nil
	}, a.logger)

	if len(ids) > 0 {
		a.sendComputeTask(ids)
	}
}

// BulkConnectorDelete
// @Param body body []BulkConnectorDeleteRequestItem true "body"
func (a *api) BulkConnectorDelete(c *gin.Context) {
	userId := c.MustGet(auth.UserKey).(string)
	ids := make([]string, 0)
	exists := make(map[string]struct{})
	bulk.Handler(c, func(request BulkConnectorDeleteRequestItem) (string, error) {
		id, err := a.store.ConnectorDelete(c, request)
		if err != nil || id == "" {
			return "", err
		}

		err = a.actionLogger.Action(context.Background(), userId, logger.LogEntry{
			Action:    logger.ActionDelete,
			ValueType: logger.ValueTypePbehavior,
			ValueID:   id,
		})
		if err != nil {
			a.actionLogger.Err(err, "failed to log action")
		}

		if _, ok := exists[id]; !ok {
			ids = append(ids, id)
			exists[id] = struct{}{}
		}

		return id, nil
	}, a.logger)

	if len(ids) > 0 {
		a.sendComputeTask(ids)
	}
}

func (a *api) sendComputeTask(ids []string) {
	a.computeChan <- ids
}

func (a *api) transformEditRequest(ctx context.Context, request *EditRequest) error {
	var err error
	request.EntityPatternFieldsRequest, err = a.transformer.TransformEntityPatternFieldsRequest(ctx, request.EntityPatternFieldsRequest)
	if err != nil {
		return err
	}

	return nil
}
