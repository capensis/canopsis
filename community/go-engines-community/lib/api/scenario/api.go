package scenario

import (
	"context"
	"encoding/json"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/action"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/valyala/fastjson"
)

type api struct {
	store        Store
	actionLogger logger.ActionLogger

	//todo: priority intervals with new requirements are looks weird now, should think about cleaner solution
	priorityIntervals action.PriorityIntervals
}

type API interface {
	GetMinimalPriority(c *gin.Context)
	CheckPriority(c *gin.Context)
	common.BulkCrudAPI
}

func NewApi(
	store Store,
	actionLogger logger.ActionLogger,
	intervals action.PriorityIntervals,
) API {
	return &api{
		store:             store,
		actionLogger:      actionLogger,
		priorityIntervals: intervals,
	}
}

// List
// @Success 200 {object} common.PaginatedListResponse{data=[]Scenario}
func (a *api) List(c *gin.Context) {
	var query FilteredQuery
	query.Query = pagination.GetDefaultQuery()

	if err := c.ShouldBind(&query); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, query))
		return
	}

	scenarios, err := a.store.Find(c.Request.Context(), query)
	if err != nil {
		panic(err)
	}

	res, err := common.NewPaginatedResponse(query.Query, scenarios)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, res)
}

// Get
// @Success 200 {object} Scenario
func (a *api) Get(c *gin.Context) {
	scenario, err := a.store.GetOneBy(c.Request.Context(), c.Param("id"))
	if err != nil {
		panic(err)
	}
	if scenario == nil {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, scenario)
}

// Create
// @Param body body EditRequest true "body"
// @Success 201 {object} Scenario
func (a *api) Create(c *gin.Context) {
	var request CreateRequest

	priority := a.priorityIntervals.GetMinimal()
	request.Priority = &priority

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	userId := c.MustGet(auth.UserKey).(string)

	scenario, err := a.store.Insert(c.Request.Context(), request)
	if err != nil {
		panic(err)
	}
	if scenario == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}
	a.priorityIntervals.Take(scenario.Priority)

	err = a.actionLogger.Action(context.Background(), userId, logger.LogEntry{
		Action:    logger.ActionCreate,
		ValueType: logger.ValueTypeScenario,
		ValueID:   scenario.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusCreated, scenario)
}

// Update
// @Param body body EditRequest true "body"
// @Success 200 {object} Scenario
func (a *api) Update(c *gin.Context) {
	request := UpdateRequest{
		ID: c.Param("id"),
	}

	oldScenario, err := a.store.GetOneBy(c.Request.Context(), request.ID)
	if err != nil {
		panic(err)
	}
	if oldScenario == nil {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	priority := a.priorityIntervals.GetMinimal()
	request.Priority = &priority

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	userId := c.MustGet(auth.UserKey).(string)

	newScenario, err := a.store.Update(c.Request.Context(), request)
	if err != nil {
		panic(err)
	}
	if newScenario == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	a.priorityIntervals.Restore(oldScenario.Priority)
	a.priorityIntervals.Take(newScenario.Priority)

	err = a.actionLogger.Action(context.Background(), userId, logger.LogEntry{
		Action:    logger.ActionUpdate,
		ValueType: logger.ValueTypeScenario,
		ValueID:   c.Param("id"),
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusOK, newScenario)
}

func (a *api) Delete(c *gin.Context) {
	id := c.Param("id")

	scenario, err := a.store.GetOneBy(c.Request.Context(), id)
	if err != nil {
		panic(err)
	}
	if scenario == nil {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	ok, err := a.store.Delete(c.Request.Context(), id)

	if err != nil {
		panic(err)
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	a.priorityIntervals.Restore(scenario.Priority)

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    logger.ActionDelete,
		ValueType: logger.ValueTypeScenario,
		ValueID:   c.Param("id"),
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.Status(http.StatusNoContent)
}

// GetMinimalPriority
// @Success 200 {object} GetMinimalPriorityResponse
func (a *api) GetMinimalPriority(c *gin.Context) {
	c.JSON(http.StatusOK, GetMinimalPriorityResponse{
		Priority: a.priorityIntervals.GetMinimal(),
	})
}

// CheckPriority
// @Param body body CheckPriorityRequest true "body"
// @Success 200 {object} CheckPriorityResponse
func (a *api) CheckPriority(c *gin.Context) {
	request := CheckPriorityRequest{}
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	valid, err := a.store.IsPriorityValid(c.Request.Context(), request.Priority)
	if err != nil {
		panic(err)
	}

	recommendedPriority := 0
	if !valid {
		recommendedPriority = a.priorityIntervals.GetMinimal()
	}

	c.JSON(http.StatusOK, CheckPriorityResponse{
		Valid:               valid,
		RecommendedPriority: recommendedPriority,
	})
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

		if request.Priority == nil {
			priority := a.priorityIntervals.GetMinimal()
			request.Priority = &priority
		}

		err = binding.Validator.ValidateStruct(request)
		if err != nil {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusBadRequest, rawObject, common.NewValidationErrorFastJsonValue(&ar, err, request)))
			continue
		}

		scenario, err := a.store.Insert(ctx, request)
		if err != nil {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusInternalServerError, rawObject, ar.NewString(err.Error())))
			continue
		}

		a.priorityIntervals.Take(scenario.Priority)
		response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, scenario.ID, http.StatusOK, rawObject, nil))

		err = a.actionLogger.Action(context.Background(), userId, logger.LogEntry{
			Action:    logger.ActionCreate,
			ValueType: logger.ValueTypeScenario,
			ValueID:   scenario.ID,
		})
		if err != nil {
			a.actionLogger.Err(err, "failed to log action")
		}
	}

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

		if request.Priority == nil {
			priority := a.priorityIntervals.GetMinimal()
			request.Priority = &priority
		}

		err = binding.Validator.ValidateStruct(request)
		if err != nil {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusBadRequest, rawObject, common.NewValidationErrorFastJsonValue(&ar, err, request)))
			continue
		}

		oldScenario, err := a.store.GetOneBy(c.Request.Context(), request.ID)
		if err != nil {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusInternalServerError, rawObject, ar.NewString(err.Error())))
			continue
		}

		if oldScenario == nil {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusNotFound, rawObject, ar.NewString("Not found")))
			continue
		}

		scenario, err := a.store.Update(ctx, UpdateRequest(request))
		if err != nil {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusInternalServerError, rawObject, ar.NewString(err.Error())))
			continue
		}

		if scenario == nil {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusNotFound, rawObject, ar.NewString("Not found")))
			continue
		}

		a.priorityIntervals.Restore(oldScenario.Priority)
		a.priorityIntervals.Take(scenario.Priority)
		response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, scenario.ID, http.StatusOK, rawObject, nil))

		err = a.actionLogger.Action(context.Background(), userId, logger.LogEntry{
			Action:    logger.ActionUpdate,
			ValueType: logger.ValueTypeScenario,
			ValueID:   scenario.ID,
		})
		if err != nil {
			a.actionLogger.Err(err, "failed to log action")
		}
	}

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

		scenario, err := a.store.GetOneBy(c.Request.Context(), request.ID)
		if err != nil {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusInternalServerError, rawObject, ar.NewString(err.Error())))
			continue
		}

		if scenario == nil {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusNotFound, rawObject, ar.NewString("Not found")))
			continue
		}

		ok, err := a.store.Delete(ctx, request.ID)
		if err != nil {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusInternalServerError, rawObject, ar.NewString(err.Error())))
			continue
		}

		if !ok {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusNotFound, rawObject, ar.NewString("Not found")))
			continue
		}

		a.priorityIntervals.Restore(scenario.Priority)
		response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, request.ID, http.StatusOK, rawObject, nil))

		err = a.actionLogger.Action(context.Background(), userId, logger.LogEntry{
			Action:    logger.ActionDelete,
			ValueType: logger.ValueTypeScenario,
			ValueID:   request.ID,
		})
		if err != nil {
			a.actionLogger.Err(err, "failed to log action")
		}
	}

	c.Data(http.StatusMultiStatus, gin.MIMEJSON, response.MarshalTo(nil))
}
