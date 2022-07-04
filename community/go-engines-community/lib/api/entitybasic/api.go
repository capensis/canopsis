package entitybasic

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin/binding"
	"github.com/valyala/fastjson"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type API interface {
	Get(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	BulkUpdate(c *gin.Context)
}

type api struct {
	store                Store
	entityChangeListener chan<- entityservice.ChangeEntityMessage
	metricMetaUpdater    metrics.MetaUpdater
	actionLogger         logger.ActionLogger
	logger               zerolog.Logger
}

func NewApi(
	store Store,
	entityChangeListener chan<- entityservice.ChangeEntityMessage,
	metricMetaUpdater metrics.MetaUpdater,
	actionLogger logger.ActionLogger,
	logger zerolog.Logger,
) API {
	return &api{
		store:                store,
		entityChangeListener: entityChangeListener,
		logger:               logger,
		metricMetaUpdater:    metricMetaUpdater,
		actionLogger:         actionLogger,
	}
}

// Get
// @Success 200 {object} Entity
func (a *api) Get(c *gin.Context) {
	var request IdRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	entity, err := a.store.GetOneBy(c.Request.Context(), request.ID)
	if err != nil {
		panic(err)
	}
	if entity == nil {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, entity)
}

// Update
// @Param body body EditRequest true "body"
// @Success 200 {object} Entity
func (a *api) Update(c *gin.Context) {
	idRequest := IdRequest{}
	if err := c.ShouldBindQuery(&idRequest); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, idRequest))
		return
	}

	request := EditRequest{}

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	request.ID = idRequest.ID
	entity, isToggled, err := a.store.Update(c.Request.Context(), request)
	if err != nil {
		panic(err)
	}

	if entity == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	if entity.Enabled || isToggled {
		a.sendChangeMessage(entityservice.ChangeEntityMessage{
			ID:         entity.ID,
			EntityType: entity.Type,
			IsToggled:  isToggled,
		})
	}

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    logger.ActionUpdate,
		ValueType: logger.ValueTypeEntity,
		ValueID:   entity.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	a.metricMetaUpdater.UpdateById(c.Request.Context(), entity.ID)

	c.JSON(http.StatusOK, entity)
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
		userObject, err := rawObject.Object()
		if err != nil {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusBadRequest, rawObject, ar.NewString(err.Error())))
			continue
		}

		var request BulkUpdateRequestItem
		err = json.Unmarshal(userObject.MarshalTo(nil), &request)
		if err != nil {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusInternalServerError, rawObject, ar.NewString(err.Error())))
			continue
		}

		err = binding.Validator.ValidateStruct(request)
		if err != nil {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusBadRequest, rawObject, common.NewValidationErrorFastJsonValue(&ar, err, request)))
			continue
		}

		sReq := request.EditRequest
		sReq.ID = request.ID

		entity, isToggled, err := a.store.Update(ctx, sReq)
		if err != nil {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusInternalServerError, rawObject, ar.NewString(err.Error())))
			continue
		}

		if entity == nil {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusNotFound, rawObject, ar.NewString("Not found")))
			continue
		}

		if entity.Enabled || isToggled {
			a.sendChangeMessage(entityservice.ChangeEntityMessage{
				ID:         entity.ID,
				EntityType: entity.Type,
				IsToggled:  isToggled,
			})
		}

		response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, entity.ID, http.StatusOK, rawObject, nil))

		err = a.actionLogger.Action(context.Background(), userId, logger.LogEntry{
			Action:    logger.ActionUpdate,
			ValueType: logger.ValueTypeEntity,
			ValueID:   entity.ID,
		})
		if err != nil {
			a.actionLogger.Err(err, "failed to log action")
		}

		a.metricMetaUpdater.UpdateById(c.Request.Context(), entity.ID)
	}

	c.Data(http.StatusMultiStatus, gin.MIMEJSON, response.MarshalTo(nil))
}

func (a *api) Delete(c *gin.Context) {
	var request IdRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	ok, err := a.store.Delete(c.Request.Context(), request.ID)

	if err != nil {
		if err == ErrLinkedEntityToAlarm {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
			return
		}
		panic(err)
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    logger.ActionDelete,
		ValueType: logger.ValueTypeEntity,
		ValueID:   request.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	a.metricMetaUpdater.DeleteById(c.Request.Context(), request.ID)

	c.Status(http.StatusNoContent)
}

func (a *api) sendChangeMessage(msg entityservice.ChangeEntityMessage) {
	select {
	case a.entityChangeListener <- msg:
	default:
		a.logger.Err(errors.New("channel is full")).
			Str("entity", msg.ID).
			Msg("fail to send change entity message")
	}
}
