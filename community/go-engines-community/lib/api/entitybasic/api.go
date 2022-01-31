package entitybasic

import (
	"errors"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
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
	actionLogger         logger.ActionLogger
	logger               zerolog.Logger
}

func NewApi(
	store Store,
	entityChangeListener chan<- entityservice.ChangeEntityMessage,
	actionLogger logger.ActionLogger,
	logger zerolog.Logger,
) API {
	return &api{
		store:                store,
		entityChangeListener: entityChangeListener,
		logger:               logger,
		actionLogger:         actionLogger,
	}
}

// Get entity by id
// @Summary Get entity by id
// @Description Get entity by id
// @Tags entitybasics
// @ID entitybasics-get-by-id
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param _id query string true "Entity id"
// @Success 200 {object} Entity
// @Failure 404 {object} common.ErrorResponse
// @Router /entitybasics [get]
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

// Update entity by id
// @Summary Update entity by id
// @Description Update entity by id
// @Tags entitybasics
// @ID entitybasics-update-by-id
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param _id query string true "Entity id"
// @Param body body EditRequest true "body"
// @Success 200 {object} Entity
// @Failure 400 {object} common.ValidationErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /entitybasics [put]
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
			ID:        entity.ID,
			IsToggled: isToggled,
		})
	}

	err = a.actionLogger.Action(c, logger.LogEntry{
		Action:    logger.ActionUpdate,
		ValueType: logger.ValueTypeEntity,
		ValueID:   entity.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusOK, entity)
}

// Delete entity by id
// @Summary Delete entity by id
// @Description Delete entity by id
// @Tags entitybasics
// @ID entitybasics-delete-by-id
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param _id query string true "Entity id"
// @Success 204
// @Failure 400 {object} common.ValidationErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /entitybasics [delete]
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

	err = a.actionLogger.Action(c, logger.LogEntry{
		Action:    logger.ActionDelete,
		ValueType: logger.ValueTypeEntity,
		ValueID:   request.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

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
