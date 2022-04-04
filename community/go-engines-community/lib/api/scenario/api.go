package scenario

import (
	"context"
	"encoding/json"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/action"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/valyala/fastjson"
	"net/http"
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

// Find all scenarios
// @Summary Find scenarios
// @Description Get paginated list of scenarios
// @Tags scenarios
// @ID scenarios-find-all
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param page query integer true "current page"
// @Param limit query integer true "items per page"
// @Param search query string false "search query"
// @Param sort query string false "sort query"
// @Param sort_by query string false "sort query"
// @Success 200 {object} common.PaginatedListResponse{data=[]Scenario}
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /scenarios [get]
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

// Get scenario by id
// @Summary Get scenario by id
// @Description Get scenario by id
// @Tags scenarios
// @ID scenarios-get-by-id
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "scenario id"
// @Success 200 {object} Scenario
// @Failure 404 {object} common.ErrorResponse
// @Router /scenarios/{id} [get]
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

// Create scenario
// @Summary Create scenario
// @Description Create scenario
// @Tags scenarios
// @ID scenarios-create
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param body body EditRequest true "body"
// @Success 201 {object} Scenario
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /scenarios [post]
func (a *api) Create(c *gin.Context) {
	var request CreateRequest

	priority := a.priorityIntervals.GetMinimal()
	request.Priority = &priority

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	userId := c.MustGet(auth.UserKey).(string)
	author := c.MustGet(auth.Username).(string)
	setActionParameterAuthorAndUserID(&request.EditRequest, author, userId)

	scenario, err := a.store.Insert(c.Request.Context(), request)
	if err != nil {
		panic(err)
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

// Update scenario by id
// @Summary Update scenario by id
// @Description Update scenario by id
// @Tags scenarios
// @ID scenarios-update-by-id
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "scenario id"
// @Param body body EditRequest true "body"
// @Success 200 {object} Scenario
// @Failure 400 {object} common.ValidationErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /scenarios/{id} [put]
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
	author := c.MustGet(auth.Username).(string)
	setActionParameterAuthorAndUserID(&request.EditRequest, author, userId)

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

// Delete scenario by id
// @Summary Delete scenario by id
// @Description Delete scenario by id
// @Tags scenarios
// @ID scenarios-delete-by-id
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "scenario id"
// @Success 204
// @Failure 404 {object} common.ErrorResponse
// @Router /scenarios/{id} [delete]
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

// Get minimal priority
// @Summary Get minimal priority
// @Description Get minimal priority
// @Tags scenarios
// @ID scenarios-get-minimal-priority
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Success 200 {object} GetMinimalPriorityResponse
// @Router /scenarios/minimal-priority [get]
func (a *api) GetMinimalPriority(c *gin.Context) {
	c.JSON(http.StatusOK, GetMinimalPriorityResponse{
		Priority: a.priorityIntervals.GetMinimal(),
	})
}

// Check priority
// @Summary Check priority
// @Description Check priority
// @Tags scenarios
// @ID scenarios-check-priority
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param body body CheckPriorityRequest true "body"
// @Success 200 {object} CheckPriorityResponse
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /scenarios/check-priority [post]
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

// Bulk create scenarios
// @Summary Bulk create scenarios
// @Description Bulk create scenarios
// @Tags scenarios
// @ID scenarios-bulk-create
// @Accept json
// @Produce json
// @Security JWTAuth
// @Security BasicAuth
// @Param body body []CreateRequest true "body"
// @Success 207 {array} []BulkCreateResponseItem
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /bulk/scenarios [post]
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

		setActionParameterAuthorAndUserID(&request.EditRequest, c.MustGet(auth.Username).(string), userId)

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

// Bulk update scenarios
// @Summary Bulk update scenarios
// @Description Bulk update scenarios
// @Tags scenarios
// @ID scenarios-bulk-update
// @Accept json
// @Produce json
// @Security JWTAuth
// @Security BasicAuth
// @Param body body []BulkUpdateRequestItem true "body"
// @Success 207 {array} []BulkUpdateResponseItem
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /bulk/scenarios [put]
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

		setActionParameterAuthorAndUserID(&request.EditRequest, c.MustGet(auth.Username).(string), userId)

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

// Bulk delete scenarios
// @Summary Bulk delete scenarios
// @Description Bulk delete scenarios
// @Tags scenarios
// @ID scenarios-bulk-delete
// @Accept json
// @Produce json
// @Security JWTAuth
// @Security BasicAuth
// @Param body body []BulkDeleteRequestItem true "body"
// @Success 207 {array} []BulkDeleteResponseItem
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /bulk/scenarios [delete]
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

func setActionParameterAuthorAndUserID(request *EditRequest, author, userID string) {
	for i, action := range request.Actions {
		switch v := action.Parameters.(type) {
		case SnoozeParametersRequest:
			if v.Author == "" {
				v.Author = author
			}
			v.User = userID
			request.Actions[i].Parameters = v
		case ChangeStateParametersRequest:
			if v.Author == "" {
				v.Author = author
			}
			v.User = userID
			request.Actions[i].Parameters = v
		case AssocTicketParametersRequest:
			if v.Author == "" {
				v.Author = author
			}
			v.User = userID
			request.Actions[i].Parameters = v
		case PbehaviorParametersRequest:
			v.Author = author
			v.User = userID
			request.Actions[i].Parameters = v
		case ParametersRequest:
			if v.Author == "" {
				v.Author = author
			}
			v.User = userID
			request.Actions[i].Parameters = v
		}
	}
}
