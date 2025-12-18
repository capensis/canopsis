package entitybasic

import (
	"errors"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
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
	errorResponder       httperror.Responder
	logger               zerolog.Logger
}

func NewApi(
	store Store,
	entityChangeListener chan<- entityservice.ChangeEntityMessage,
	metricMetaUpdater metrics.MetaUpdater,
	errorResponder httperror.Responder,
	logger zerolog.Logger,
) API {
	return &api{
		store:                store,
		entityChangeListener: entityChangeListener,
		metricMetaUpdater:    metricMetaUpdater,
		errorResponder:       errorResponder,
		logger:               logger,
	}
}

// Get
// @Success 200 {object} Entity
func (a *api) Get(c *gin.Context) {
	var request IdRequest
	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	entity, err := a.store.GetOneBy(c, request.ID)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	if entity == nil {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	c.JSON(http.StatusOK, entity)
}

// Update
// @Param body body EditRequest true "body"
// @Success 200 {object} Entity
func (a *api) Update(c *gin.Context) {
	idRequest := IdRequest{}
	if err := validation.BindQuery(c, &idRequest); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	request := EditRequest{}

	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	request.ID = idRequest.ID
	entity, isToggled, err := a.store.Update(c, request)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if entity == nil {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	if entity.Enabled || isToggled {
		msg := entityservice.ChangeEntityMessage{
			ID:         entity.ID,
			Name:       entity.Name,
			Component:  entity.Component,
			EntityType: entity.Type,
			IsToggled:  isToggled,
		}

		if !entity.Enabled && entity.Type == types.EntityTypeComponent {
			msg.Resources = make([]string, len(entity.Resources))
			copy(msg.Resources, entity.Resources)
		}

		a.sendChangeMessage(msg)
	}

	a.metricMetaUpdater.UpdateById(c, entity.ID)
	if isToggled && entity.Type == types.EntityTypeComponent {
		a.metricMetaUpdater.UpdateById(c, entity.Resources...)
	}

	c.JSON(http.StatusOK, entity)
}

func (a *api) Delete(c *gin.Context) {
	var request IdRequest
	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	entity, err := a.store.Delete(c, request.ID, userID)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if entity == nil {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	a.metricMetaUpdater.DeleteById(c, request.ID)
	a.sendChangeMessage(entityservice.ChangeEntityMessage{
		ID:         entity.ID,
		Name:       entity.Name,
		Component:  entity.Component,
		EntityType: entity.Type,
		IsDeleted:  true,
	})

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
