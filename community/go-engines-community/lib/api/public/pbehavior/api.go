package pbehavior

import (
	"context"
	"net/http"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/binding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/bulk"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/response"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

type API interface {
	bulk.API
	Patch(w http.ResponseWriter, r *http.Request)
	BulkConnectorCreate(w http.ResponseWriter, r *http.Request)
	BulkConnectorDelete(w http.ResponseWriter, r *http.Request)
	BulkConnectorEdit(w http.ResponseWriter, r *http.Request)
}

func NewAPI(
	store Store,
	binder binding.Binder,
	responder response.Responder,
	amqpPublisher amqp.Publisher,
	encoder encoding.Encoder,
	decoder encoding.Decoder,
	logger zerolog.Logger,
) API {
	return &api{
		store:         store,
		binder:        binder,
		responder:     responder,
		amqpPublisher: amqpPublisher,
		encoder:       encoder,
		decoder:       decoder,
		logger:        logger,
	}
}

type api struct {
	store         Store
	binder        binding.Binder
	responder     response.Responder
	amqpPublisher amqp.Publisher
	encoder       encoding.Encoder
	decoder       encoding.Decoder
	logger        zerolog.Logger
}

func (a *api) List(w http.ResponseWriter, r *http.Request) {
	var d ListRequest
	if err := a.binder.BindForm(r, &d); err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	res, err := a.store.Find(r.Context(), d)
	if err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	a.responder.Respond(w, r, http.StatusOK, res)
}

// Get
// @Success 200 {object} Response
func (a *api) Get(w http.ResponseWriter, r *http.Request) {
	id := a.binder.URLParam(r, "id")
	res, err := a.store.GetByID(r.Context(), id)
	if err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	a.responder.Respond(w, r, http.StatusOK, res)
}

// Create
// @Param body body EditRequest true "body"
// @Success 201 {object} Response
func (a *api) Create(w http.ResponseWriter, r *http.Request) {
	var d EditRequest
	if err := a.binder.Bind(r, &d); err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	res, err := a.store.Insert(r.Context(), d)
	if err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	a.sendRecomputeEvent(r.Context(), rpc.PbehaviorRecomputeEvent{Ids: []string{res.ID}, RecomputeInherited: res.Inherited})

	a.responder.Respond(w, r, http.StatusCreated, res)
}

// Update
// @Param body body EditRequest true "body"
// @Success 200 {object} Response
func (a *api) Update(w http.ResponseWriter, r *http.Request) {
	var d EditRequest
	if err := a.binder.Bind(r, &d); err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	id := a.binder.URLParam(r, "id")
	res, recomputeInherited, err := a.store.Update(r.Context(), id, d)
	if err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	a.sendRecomputeEvent(r.Context(), rpc.PbehaviorRecomputeEvent{Ids: []string{res.ID}, RecomputeInherited: recomputeInherited})

	a.responder.Respond(w, r, http.StatusOK, res)
}

// Patch
// @Param body body PatchRequest true "body"
// @Success 200 {object} Response
func (a *api) Patch(w http.ResponseWriter, r *http.Request) {
	var d PatchRequest
	if err := a.binder.Bind(r, &d); err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	id := a.binder.URLParam(r, "id")
	res, recomputeInherited, err := a.store.UpdateByPatch(r.Context(), id, d)
	if err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	a.sendRecomputeEvent(r.Context(), rpc.PbehaviorRecomputeEvent{Ids: []string{res.ID}, RecomputeInherited: recomputeInherited})

	a.responder.Respond(w, r, http.StatusOK, res)
}

func (a *api) Delete(w http.ResponseWriter, r *http.Request) {
	userID, err := authctx.GetUserKey(r.Context())
	if err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	id := a.binder.URLParam(r, "id")
	recomputeInherited, err := a.store.Delete(r.Context(), id, userID)
	if err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	a.sendRecomputeEvent(r.Context(), rpc.PbehaviorRecomputeEvent{Ids: []string{id}, RecomputeInherited: recomputeInherited})

	a.responder.Respond(w, r, http.StatusNoContent, nil)
}

// BulkCreate
// @Param body body []EditRequest true "body"
func (a *api) BulkCreate(w http.ResponseWriter, r *http.Request) {
	ids := make([]string, 0)
	recomputeInherited := false
	bulk.Handler(w, r, func(d EditRequest) (string, error) {
		res, err := a.store.Insert(r.Context(), d)
		if err != nil {
			return "", err
		}

		ids = append(ids, res.ID)
		recomputeInherited = recomputeInherited || res.Inherited

		return res.ID, nil
	}, a.responder, a.binder, a.decoder)

	if len(ids) > 0 {
		a.sendRecomputeEvent(r.Context(), rpc.PbehaviorRecomputeEvent{Ids: ids, RecomputeInherited: recomputeInherited})
	}
}

// BulkUpdate
// @Param body body []BulkUpdateRequestItem true "body"
func (a *api) BulkUpdate(w http.ResponseWriter, r *http.Request) {
	ids := make([]string, 0)
	recomputeInherited := false
	processed := make(map[string]bool)
	bulk.Handler(w, r, func(d BulkUpdateRequestItem) (string, error) {
		res, current, err := a.store.Update(r.Context(), d.ID, d.EditRequest)
		if err != nil {
			return "", err
		}

		if !processed[res.ID] {
			ids = append(ids, res.ID)
			processed[res.ID] = true
		}

		recomputeInherited = recomputeInherited || current

		return res.ID, nil
	}, a.responder, a.binder, a.decoder)

	if len(ids) > 0 {
		a.sendRecomputeEvent(r.Context(), rpc.PbehaviorRecomputeEvent{Ids: ids, RecomputeInherited: recomputeInherited})
	}
}

// BulkDelete
// @Param body body []bulk.DeleteRequestItem true "body"
func (a *api) BulkDelete(w http.ResponseWriter, r *http.Request) {
	userID, err := authctx.GetUserKey(r.Context())
	if err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	ids := make([]string, 0)
	recomputeInherited := false
	bulk.Handler(w, r, func(d bulk.DeleteRequestItem) (string, error) {
		current, err := a.store.Delete(r.Context(), d.ID, userID)
		if err != nil {
			return "", err
		}

		recomputeInherited = recomputeInherited || current

		return d.ID, nil
	}, a.responder, a.binder, a.decoder)

	if len(ids) > 0 {
		a.sendRecomputeEvent(r.Context(), rpc.PbehaviorRecomputeEvent{Ids: ids, RecomputeInherited: recomputeInherited})
	}
}

// BulkConnectorCreate
// @Param body body []BulkConnectorCreateRequestItem true "body"
func (a *api) BulkConnectorCreate(w http.ResponseWriter, r *http.Request) {
	idsByOrigin := make(map[string][]string)
	processed := make(map[string]bool)
	bulk.HandlerWithGrouping(
		w,
		r,
		func(prev, cur BulkConnectorCreateRequestItem) bool {
			return cur.Origin == prev.Origin &&
				cur.Comment == prev.Comment &&
				cur.Start.Unix() == prev.Start.Unix() &&
				cur.Stop.Unix() == prev.Stop.Unix()
		},
		func(merged, cur BulkConnectorCreateRequestItem) BulkConnectorCreateRequestItem {
			merged.Entities = append(merged.Entities, cur.Entities...)

			return merged
		},
		func(d BulkConnectorCreateRequestItem) (string, error) {
			pbh, err := a.store.ConnectorCreate(r.Context(), d)
			if err != nil {
				return "", err
			}

			if !processed[pbh.ID] {
				// skip add to idsByOrigin when pbh is not effective
				if pbh.Stop == nil || pbh.RRule != "" || pbh.Stop.Unix() >= time.Now().Unix() {
					idsByOrigin[d.Origin] = append(idsByOrigin[d.Origin], pbh.ID)
				}
				processed[pbh.ID] = true
			}

			return pbh.ID, nil
		},
		a.responder,
		a.binder,
		a.decoder,
	)

	for origin, ids := range idsByOrigin {
		a.sendRecomputeEvent(r.Context(), rpc.PbehaviorRecomputeEvent{
			Ids:       ids,
			Author:    origin,
			Initiator: types.InitiatorExternal,
		})
	}
}

// BulkConnectorDelete
// @Param body body []BulkConnectorDeleteRequestItem true "body"
func (a *api) BulkConnectorDelete(w http.ResponseWriter, r *http.Request) {
	idsByOrigin := make(map[string][]string)
	processed := make(map[string]bool)
	bulk.HandlerWithGrouping(
		w,
		r,
		func(prev, cur BulkConnectorDeleteRequestItem) bool {
			return cur.Origin == prev.Origin &&
				cur.Comment == prev.Comment &&
				cur.Start.Unix() == prev.Start.Unix() &&
				cur.Stop.Unix() == prev.Stop.Unix()
		},
		func(merged, cur BulkConnectorDeleteRequestItem) BulkConnectorDeleteRequestItem {
			merged.Entities = append(merged.Entities, cur.Entities...)

			return merged
		},
		func(d BulkConnectorDeleteRequestItem) (string, error) {
			id, err := a.store.ConnectorDelete(r.Context(), d)
			if err != nil {
				return "", err
			}

			if !processed[id] {
				idsByOrigin[d.Origin] = append(idsByOrigin[d.Origin], id)
				processed[id] = true
			}

			return id, nil
		},
		a.responder,
		a.binder,
		a.decoder,
	)

	for origin, ids := range idsByOrigin {
		a.sendRecomputeEvent(r.Context(), rpc.PbehaviorRecomputeEvent{
			Ids:       ids,
			Author:    origin,
			Initiator: types.InitiatorExternal,
		})
	}
}

// BulkConnectorEdit
// @Param body body []BulkConnectorEditRequestItem true "body"
func (a *api) BulkConnectorEdit(w http.ResponseWriter, r *http.Request) {
	idsByOrigin := make(map[string][]string)
	processed := make(map[string]bool)
	bulk.HandlerWithGrouping(
		w,
		r,
		func(prev, cur BulkConnectorEditRequestItem) bool {
			return cur.Action == prev.Action &&
				cur.Origin == prev.Origin &&
				cur.Comment == prev.Comment &&
				cur.Start.Unix() == prev.Start.Unix() &&
				cur.Stop.Unix() == prev.Stop.Unix()
		},
		func(merged, cur BulkConnectorEditRequestItem) BulkConnectorEditRequestItem {
			merged.Entities = append(merged.Entities, cur.Entities...)

			return merged
		},
		func(request BulkConnectorEditRequestItem) (string, error) {
			var id string
			var err error
			switch request.Action {
			case BulkConnectorActionCreate:
				var res Response
				res, err = a.store.ConnectorCreate(r.Context(), BulkConnectorCreateRequestItem{
					Author:   request.Author,
					Entities: request.Entities,
					Origin:   request.Origin,
					Start:    request.Start,
					Stop:     request.Stop,
					Comment:  request.Comment,
					Name:     request.Name,
					Reason:   request.Reason,
					Type:     request.Type,
					Color:    request.Color,
				})
				id = res.ID
			case BulkConnectorActionDelete:
				id, err = a.store.ConnectorDelete(r.Context(), BulkConnectorDeleteRequestItem{
					Author:   request.Author,
					Entities: request.Entities,
					Origin:   request.Origin,
					Start:    request.Start,
					Stop:     request.Stop,
					Comment:  request.Comment,
				})
			default:
				return "", validation.NewSingleErrorWithParam("oneof", "Action", "Action", BulkConnectorActionCreate+" "+BulkConnectorActionDelete, request)
			}

			if err != nil {
				return "", err
			}

			if !processed[id] {
				idsByOrigin[request.Origin] = append(idsByOrigin[request.Origin], id)
				processed[id] = true
			}

			return id, nil
		},
		a.responder,
		a.binder,
		a.decoder,
	)

	for origin, ids := range idsByOrigin {
		a.sendRecomputeEvent(r.Context(), rpc.PbehaviorRecomputeEvent{
			Ids:       ids,
			Author:    origin,
			Initiator: types.InitiatorExternal,
		})
	}
}

func (a *api) sendRecomputeEvent(ctx context.Context, event rpc.PbehaviorRecomputeEvent) {
	body, err := a.encoder.Encode(event)
	if err != nil {
		a.logger.Err(err).Msg("cannot encode event")
	}

	err = a.amqpPublisher.PublishWithContext(
		ctx,
		canopsis.DefaultExchangeName,
		canopsis.PBehaviorQueueRecomputeName,
		false,
		false,
		amqp091.Publishing{
			ContentType:  canopsis.JsonContentType,
			Body:         body,
			DeliveryMode: amqp091.Persistent,
		},
	)
	if err != nil {
		a.logger.Err(err).Msg("cannot send event")
	}
}
