package scenario

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/bulk"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type api struct {
	store        Store
	actionLogger logger.ActionLogger
	transformer  common.PatternFieldsTransformer
	logger       zerolog.Logger
}

func NewApi(
	store Store,
	actionLogger logger.ActionLogger,
	transformer common.PatternFieldsTransformer,
	logger zerolog.Logger,
) common.BulkCrudAPI {
	return &api{
		store:        store,
		actionLogger: actionLogger,
		logger:       logger,
		transformer:  transformer,
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

	scenarios, err := a.store.Find(c, query)
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
	scenario, err := a.store.GetOneBy(c, c.Param("id"))
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
	var err error

	if err = c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	err = a.transformEditRequest(c, &request.EditRequest)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}
		panic(err)
	}

	scenario, err := a.store.Insert(c, request)
	if err != nil {
		panic(err)
	}
	if scenario == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
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

	oldScenario, err := a.store.GetOneBy(c, request.ID)
	if err != nil {
		panic(err)
	}
	if oldScenario == nil {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	err = a.transformEditRequest(c, &request.EditRequest)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}
		panic(err)
	}

	newScenario, err := a.store.Update(c, request)
	if err != nil {
		panic(err)
	}
	if newScenario == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
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

	scenario, err := a.store.GetOneBy(c, id)
	if err != nil {
		panic(err)
	}
	if scenario == nil {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	ok, err := a.store.Delete(c, id)

	if err != nil {
		panic(err)
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

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

// BulkCreate
// @Param body body []CreateRequest true "body"
func (a *api) BulkCreate(c *gin.Context) {
	userId := c.MustGet(auth.UserKey).(string)
	bulk.Handler(c, func(request CreateRequest) (string, error) {
		err := a.transformEditRequest(c, &request.EditRequest)
		if err != nil {
			return "", err
		}

		scenario, err := a.store.Insert(c, request)
		if err != nil {
			return "", err
		}

		err = a.actionLogger.Action(context.Background(), userId, logger.LogEntry{
			Action:    logger.ActionCreate,
			ValueType: logger.ValueTypeScenario,
			ValueID:   scenario.ID,
		})
		if err != nil {
			a.actionLogger.Err(err, "failed to log action")
		}

		return scenario.ID, nil
	}, a.logger)
}

// BulkUpdate
// @Param body body []BulkUpdateRequestItem true "body"
func (a *api) BulkUpdate(c *gin.Context) {
	userId := c.MustGet(auth.UserKey).(string)
	bulk.Handler(c, func(request BulkUpdateRequestItem) (string, error) {
		oldScenario, err := a.store.GetOneBy(c, request.ID)
		if err != nil || oldScenario == nil {
			return "", err
		}

		err = a.transformEditRequest(c, &request.EditRequest)
		if err != nil {
			return "", err
		}

		scenario, err := a.store.Update(c, UpdateRequest(request))
		if err != nil || scenario == nil {
			return "", err
		}

		err = a.actionLogger.Action(context.Background(), userId, logger.LogEntry{
			Action:    logger.ActionUpdate,
			ValueType: logger.ValueTypeScenario,
			ValueID:   scenario.ID,
		})
		if err != nil {
			a.actionLogger.Err(err, "failed to log action")
		}

		return scenario.ID, nil
	}, a.logger)
}

// BulkDelete
// @Param body body []BulkDeleteRequestItem true "body"
func (a *api) BulkDelete(c *gin.Context) {
	userId := c.MustGet(auth.UserKey).(string)
	bulk.Handler(c, func(request BulkDeleteRequestItem) (string, error) {
		scenario, err := a.store.GetOneBy(c, request.ID)
		if err != nil || scenario == nil {
			return "", err
		}

		ok, err := a.store.Delete(c, request.ID)
		if err != nil || !ok {
			return "", err
		}

		err = a.actionLogger.Action(context.Background(), userId, logger.LogEntry{
			Action:    logger.ActionDelete,
			ValueType: logger.ValueTypeScenario,
			ValueID:   request.ID,
		})
		if err != nil {
			a.actionLogger.Err(err, "failed to log action")
		}

		return scenario.ID, nil
	}, a.logger)
}

func (a *api) transformEditRequest(ctx context.Context, request *EditRequest) error {
	var err error

	for idx, actionRequest := range request.Actions {
		actionRequest.AlarmPatternFieldsRequest, err = a.transformer.TransformAlarmPatternFieldsRequest(ctx, actionRequest.AlarmPatternFieldsRequest)
		if err != nil {
			if errors.Is(err, common.ErrNotExistCorporateAlarmPattern) {
				return common.NewValidationError(fmt.Sprintf("actions.%d.corporate_alarm_pattern", idx), err.Error())
			}
			return err
		}

		actionRequest.EntityPatternFieldsRequest, err = a.transformer.TransformEntityPatternFieldsRequest(ctx, actionRequest.EntityPatternFieldsRequest)
		if err != nil {
			if errors.Is(err, common.ErrNotExistCorporateEntityPattern) {
				return common.NewValidationError(fmt.Sprintf("actions.%d.corporate_entity_pattern", idx), err.Error())
			}
			return err
		}

		request.Actions[idx] = actionRequest
	}

	return nil
}
