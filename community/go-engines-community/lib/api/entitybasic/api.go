package entitybasic

import (
	"context"
	"errors"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type API interface {
	Get(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
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

	entity, err := a.store.GetOneBy(c, request.ID)
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
	entity, isToggled, err := a.store.Update(c, request)
	if err != nil {
		panic(err)
	}

	if entity == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	if entity.Enabled || isToggled {
		msg := entityservice.ChangeEntityMessage{
			ID:         entity.ID,
			EntityType: entity.Type,
			IsToggled:  isToggled,
		}

		if !entity.Enabled && entity.Type == types.EntityTypeComponent {
			msg.Resources = make([]string, len(entity.Resources))
			copy(msg.Resources, entity.Resources)
		}

		a.sendChangeMessage(msg)
	}

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    logger.ActionUpdate,
		ValueType: logger.ValueTypeEntity,
		ValueID:   entity.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	a.metricMetaUpdater.UpdateById(c, entity.ID)
	if isToggled && entity.Type == types.EntityTypeComponent {
		a.metricMetaUpdater.UpdateById(c, entity.Resources...)
	}

	c.JSON(http.StatusOK, entity)
}

func (a *api) Delete(c *gin.Context) {
	var request IdRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	ok, err := a.store.Delete(c, request.ID)

	if err != nil {
		if err == ErrLinkedEntityToAlarm || err == ErrComponent {
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

	a.metricMetaUpdater.DeleteById(c, request.ID)

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
