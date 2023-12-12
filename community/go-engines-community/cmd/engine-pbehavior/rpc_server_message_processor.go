package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	libpbehavior "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type rpcServerMessageProcessor struct {
	FeaturePrintEventOnError bool
	DbClient                 mongo.DbClient
	PbhService               libpbehavior.Service
	EventManager             libpbehavior.EventManager
	TimezoneConfigProvider   config.TimezoneConfigProvider
	Decoder                  encoding.Decoder
	Encoder                  encoding.Encoder
	Logger                   zerolog.Logger
}

func (p *rpcServerMessageProcessor) Process(ctx context.Context, d amqp.Delivery) ([]byte, error) {
	msg := d.Body
	var event rpc.PbehaviorEvent
	err := p.Decoder.Decode(msg, &event)
	if err != nil || event.Alarm == nil || event.Entity == nil {
		p.logError(err, "invalid event", msg)

		return p.getErrRpcEvent(errors.New("invalid event")), nil
	}

	var pbhEvent *types.Event
	if !event.Healthcheck {
		pbhEvent, err = p.processCreatePbhEvent(
			ctx,
			*event.Alarm,
			*event.Entity,
			event.Params,
		)
		if err != nil {
			if engine.IsConnectionError(err) {
				return nil, err
			}

			p.logError(err, "cannot process event", msg)
			return p.getErrRpcEvent(err), nil
		}
	}

	if pbhEvent == nil {
		pbhEvent = &types.Event{}
	} else {
		pbhEvent.Entity = event.Entity
	}

	return p.getRpcEvent(rpc.PbehaviorResultEvent{
		Alarm:    event.Alarm,
		Entity:   event.Entity,
		PbhEvent: *pbhEvent,
	})
}

func (p *rpcServerMessageProcessor) processCreatePbhEvent(
	ctx context.Context,
	alarm types.Alarm,
	entity types.Entity,
	params rpc.PbehaviorParameters,
) (*types.Event, error) {
	pbehavior, err := p.createPbehavior(ctx, params, entity)
	if err != nil {
		return nil, err
	}

	resolver, err := p.PbhService.RecomputeByIds(ctx, []string{pbehavior.ID})
	if err != nil {
		return nil, fmt.Errorf("pbehavior recompute failed: %w", err)
	}

	p.Logger.Debug().Str("pbehavior", pbehavior.ID).Msg("pbehavior recomputed")

	resolveResult, err := p.getResolveResult(ctx, entity, resolver)
	if err != nil {
		return nil, err
	}

	pbhEvent := p.EventManager.GetEvent(resolveResult, alarm, time.Now())
	if pbhEvent.EventType == "" {
		return nil, nil
	}

	return &pbhEvent, nil
}

func (p *rpcServerMessageProcessor) createPbehavior(
	ctx context.Context,
	params rpc.PbehaviorParameters,
	entity types.Entity,
) (*libpbehavior.PBehavior, error) {
	typeCollection := p.DbClient.Collection(mongo.PbehaviorTypeMongoCollection)
	err := typeCollection.FindOne(ctx, bson.M{"_id": params.Type}).Err()
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, fmt.Errorf("pbehavior type not exist: %q", params.Type)
		}

		return nil, fmt.Errorf("cannot get pbehavior type: %w", err)
	}

	reasonCollection := p.DbClient.Collection(mongo.PbehaviorReasonMongoCollection)
	err = reasonCollection.FindOne(ctx, bson.M{"_id": params.Reason}).Err()
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, fmt.Errorf("pbehavior reason not exist: %q", params.Reason)
		}

		return nil, fmt.Errorf("cannot get pbehavior reason: %w", err)
	}

	now := datetime.NewCpsTime()
	var start, stop datetime.CpsTime
	if params.Tstart != nil && params.Tstop != nil {
		start = *params.Tstart
		stop = *params.Tstop
	} else if params.StartOnTrigger != nil && *params.StartOnTrigger &&
		params.Duration != nil && params.Duration.Value > 0 {
		start = now
		stop = params.Duration.AddTo(now)
	}

	if start.IsZero() {
		return nil, fmt.Errorf("invalid action parameters, tstart with tstop or start_on_trigger with duration must be defined: %+v", params)
	}

	pbehavior := libpbehavior.PBehavior{
		ID:      utils.NewID(),
		Enabled: true,
		Name:    params.Name,
		Reason:  params.Reason,
		RRule:   params.RRule,
		Start:   &start,
		Stop:    &stop,
		Type:    params.Type,
		Color:   types.ActionPbehaviorColor,
		Created: &now,
		Updated: &now,

		EntityPatternFields: savedpattern.EntityPatternFields{
			EntityPattern: pattern.Entity{
				{
					{
						Field: "_id",
						Condition: pattern.Condition{
							Type:  pattern.ConditionEqual,
							Value: entity.ID,
						},
					},
				},
			},
		},
	}

	collection := p.DbClient.Collection(mongo.PbehaviorMongoCollection)
	_, err = collection.InsertOne(ctx, pbehavior)
	if err != nil {
		return nil, fmt.Errorf("create new pbehavior failed: %w", err)
	}

	p.Logger.Info().Str("pbehavior", pbehavior.ID).Msg("create pbehavior")
	return &pbehavior, nil
}

func (p *rpcServerMessageProcessor) getResolveResult(
	ctx context.Context,
	entity types.Entity,
	resolver libpbehavior.ComputedEntityTypeResolver,
) (libpbehavior.ResolveResult, error) {
	location := p.TimezoneConfigProvider.Get().Location
	now := time.Now().In(location)
	resolveResult, err := resolver.Resolve(ctx, entity, now)
	if err != nil {
		return libpbehavior.ResolveResult{}, fmt.Errorf("resolve an entity failed: %w", err)
	}

	return resolveResult, nil
}

func (p *rpcServerMessageProcessor) getErrRpcEvent(err error) []byte {
	msg, _ := p.getRpcEvent(rpc.PbehaviorResultEvent{Error: &rpc.Error{Error: err}})
	return msg
}

func (p *rpcServerMessageProcessor) getRpcEvent(event rpc.PbehaviorResultEvent) ([]byte, error) {
	msg, err := p.Encoder.Encode(event)
	if err != nil {
		p.Logger.Err(err).Msg("cannot encode event")
		return nil, nil
	}

	return msg, nil
}

func (p *rpcServerMessageProcessor) logError(err error, errMsg string, msg []byte) {
	if p.FeaturePrintEventOnError {
		p.Logger.Err(err).Str("event", string(msg)).Msg(errMsg)
	} else {
		p.Logger.Err(err).Msg(errMsg)
	}
}
