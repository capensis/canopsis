package pbehavior

import (
	"context"
	"errors"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/bulk"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
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
	store       Store
	computeChan chan<- rpc.PbehaviorRecomputeEvent
	logger      zerolog.Logger

	transformer common.PatternFieldsTransformer
}

func NewApi(
	store Store,
	computeChan chan<- rpc.PbehaviorRecomputeEvent,
	transformer common.PatternFieldsTransformer,
	logger zerolog.Logger,
) API {
	return &api{
		store:       store,
		computeChan: computeChan,
		transformer: transformer,
		logger:      logger,
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

	aggregationResult, err := a.store.Find(c, r)
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

	entity, err := a.store.FindEntity(c, r.ID)
	if err != nil {
		panic(err)
	}
	if entity == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	res, err := a.store.FindByEntityID(c, *entity, r)
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

	entity, err := a.store.FindEntity(c, r.ID)
	if err != nil {
		panic(err)
	}
	if entity == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	res, err := a.store.CalendarByEntityID(c, *entity, r)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, res)
}

// Get
// @Success 200 {object} Response
func (a *api) Get(c *gin.Context) {
	pbh, err := a.store.GetOneBy(c, c.Param("id"))
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

	aggregationResult, err := a.store.FindEntities(c, c.Param("id"), r)
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

	err := a.transformEditRequest(c, &request.EditRequest)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}
		panic(err)
	}

	pbh, err := a.store.Insert(c, request)
	if err != nil {
		validationErr := common.ValidationError{}
		if errors.As(err, &validationErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, validationErr.ValidationErrorResponse())
			return
		}
		panic(err)
	}

	a.sendComputeTask(rpc.PbehaviorRecomputeEvent{Ids: []string{pbh.ID}})

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

	err := a.transformEditRequest(c, &request.EditRequest)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}
		panic(err)
	}

	pbh, err := a.store.Update(c, request)
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

	a.sendComputeTask(rpc.PbehaviorRecomputeEvent{Ids: []string{pbh.ID}})

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
		r, err := a.transformer.TransformEntityPatternFieldsRequest(c, common.EntityPatternFieldsRequest{
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

	pbh, err := a.store.UpdateByPatch(c, request)
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

	a.sendComputeTask(rpc.PbehaviorRecomputeEvent{Ids: []string{pbh.ID}})

	c.JSON(http.StatusOK, pbh)
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

	a.sendComputeTask(rpc.PbehaviorRecomputeEvent{Ids: []string{id}})
	c.JSON(http.StatusNoContent, nil)
}

func (a *api) DeleteByName(c *gin.Context) {
	request := DeleteByNameRequest{}

	if err := c.ShouldBindQuery(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))

		return
	}

	id, err := a.store.DeleteByName(c, request.Name, c.MustGet(auth.UserKey).(string))
	if err != nil {
		panic(err)
	}

	if id == "" {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	a.sendComputeTask(rpc.PbehaviorRecomputeEvent{Ids: []string{id}})
	c.JSON(http.StatusNoContent, nil)
}

// BulkCreate
// @Param body body []CreateRequest true "body"
func (a *api) BulkCreate(c *gin.Context) {
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

		ids = append(ids, pbh.ID)

		return pbh.ID, nil
	}, a.logger)

	if len(ids) > 0 {
		a.sendComputeTask(rpc.PbehaviorRecomputeEvent{Ids: ids})
	}
}

// BulkUpdate
// @Param body body []BulkUpdateRequestItem true "body"
func (a *api) BulkUpdate(c *gin.Context) {
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

		if _, ok := exists[pbh.ID]; !ok {
			ids = append(ids, pbh.ID)
			exists[pbh.ID] = struct{}{}
		}

		return pbh.ID, nil
	}, a.logger)

	if len(ids) > 0 {
		a.sendComputeTask(rpc.PbehaviorRecomputeEvent{Ids: ids})
	}
}

// BulkDelete
// @Param body body []BulkDeleteRequestItem true "body"
func (a *api) BulkDelete(c *gin.Context) {
	userId := c.MustGet(auth.UserKey).(string)

	ids := make([]string, 0)
	bulk.Handler(c, func(request BulkDeleteRequestItem) (string, error) {
		ok, err := a.store.Delete(c, request.ID, userId)
		if err != nil || !ok {
			return "", err
		}

		ids = append(ids, request.ID)

		return request.ID, nil
	}, a.logger)

	if len(ids) > 0 {
		a.sendComputeTask(rpc.PbehaviorRecomputeEvent{Ids: ids})
	}
}

// BulkEntityCreate
// @Param body body []BulkEntityCreateRequestItem true "body"
func (a *api) BulkEntityCreate(c *gin.Context) {
	userId := c.MustGet(auth.UserKey).(string)
	username := c.MustGet(auth.Username).(string)
	ids := make([]string, 0)
	bulk.Handler(c, func(request BulkEntityCreateRequestItem) (string, error) {
		pbh, err := a.store.EntityInsert(c, request)
		if err != nil || pbh == nil {
			return "", err
		}

		ids = append(ids, pbh.ID)

		return pbh.ID, nil
	}, a.logger)

	if len(ids) > 0 {
		a.sendComputeTask(rpc.PbehaviorRecomputeEvent{
			Ids:       ids,
			Author:    username,
			UserID:    userId,
			Initiator: types.InitiatorUser,
		})
	}
}

// BulkEntityDelete
// @Param body body []BulkEntityDeleteRequestItem true "body"
func (a *api) BulkEntityDelete(c *gin.Context) {
	userId := c.MustGet(auth.UserKey).(string)
	username := c.MustGet(auth.Username).(string)
	ids := make([]string, 0)
	bulk.Handler(c, func(request BulkEntityDeleteRequestItem) (string, error) {
		id, err := a.store.EntityDelete(c, request)
		if err != nil || id == "" {
			return "", err
		}

		ids = append(ids, id)

		return id, nil
	}, a.logger)

	if len(ids) > 0 {
		a.sendComputeTask(rpc.PbehaviorRecomputeEvent{
			Ids:       ids,
			Author:    username,
			UserID:    userId,
			Initiator: types.InitiatorUser,
		})
	}
}

// BulkConnectorCreate
// @Param body body []BulkConnectorCreateRequestItem true "body"
func (a *api) BulkConnectorCreate(c *gin.Context) {
	idsByOrigin := make(map[string][]string)
	exists := make(map[string]struct{})
	bulk.Handler(c, func(request BulkConnectorCreateRequestItem) (string, error) {
		pbh, err := a.store.ConnectorCreate(c, request)
		if err != nil || pbh == nil {
			return "", err
		}

		if _, ok := exists[pbh.ID]; !ok {
			idsByOrigin[request.Origin] = append(idsByOrigin[request.Origin], pbh.ID)
			exists[pbh.ID] = struct{}{}
		}

		return pbh.ID, nil
	}, a.logger)

	for origin, ids := range idsByOrigin {
		a.sendComputeTask(rpc.PbehaviorRecomputeEvent{
			Ids:       ids,
			Author:    origin,
			Initiator: types.InitiatorExternal,
		})
	}
}

// BulkConnectorDelete
// @Param body body []BulkConnectorDeleteRequestItem true "body"
func (a *api) BulkConnectorDelete(c *gin.Context) {
	idsByOrigin := make(map[string][]string)
	exists := make(map[string]struct{})
	bulk.Handler(c, func(request BulkConnectorDeleteRequestItem) (string, error) {
		id, err := a.store.ConnectorDelete(c, request)
		if err != nil || id == "" {
			return "", err
		}

		if _, ok := exists[id]; !ok {
			idsByOrigin[request.Origin] = append(idsByOrigin[request.Origin], id)
			exists[id] = struct{}{}
		}

		return id, nil
	}, a.logger)

	for origin, ids := range idsByOrigin {
		a.sendComputeTask(rpc.PbehaviorRecomputeEvent{
			Ids:       ids,
			Author:    origin,
			Initiator: types.InitiatorExternal,
		})
	}
}

func (a *api) sendComputeTask(event rpc.PbehaviorRecomputeEvent) {
	a.computeChan <- event
}

func (a *api) transformEditRequest(ctx context.Context, request *EditRequest) error {
	var err error
	request.EntityPatternFieldsRequest, err = a.transformer.TransformEntityPatternFieldsRequest(ctx, request.EntityPatternFieldsRequest)
	if err != nil {
		return err
	}

	return nil
}
