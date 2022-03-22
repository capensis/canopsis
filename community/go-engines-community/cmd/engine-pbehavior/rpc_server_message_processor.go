package main

import (
	"context"
	"errors"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	libpbehavior "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"time"
)

type rpcServerMessageProcessor struct {
	FeaturePrintEventOnError bool
	Processor                createPbehaviorMessageProcessor
	Decoder                  encoding.Decoder
	Encoder                  encoding.Encoder
	Logger                   zerolog.Logger
}

func (p *rpcServerMessageProcessor) Process(ctx context.Context, d amqp.Delivery) ([]byte, error) {
	msg := d.Body
	var event types.RPCPBehaviorEvent
	err := p.Decoder.Decode(msg, &event)
	if err != nil || event.Alarm == nil {
		p.logError(err, "invalid event", msg)

		return p.getErrRpcEvent(errors.New("invalid event")), nil
	}

	pbhEvent, err := p.Processor.Process(
		ctx,
		event.Alarm,
		event.Entity,
		event.Params,
		msg,
	)
	if err != nil {
		if engine.IsConnectionError(err) {
			return nil, err
		}

		p.logError(err, "cannot process event", msg)
		return p.getErrRpcEvent(err), nil
	}

	if pbhEvent == nil {
		pbhEvent = &types.Event{}
	}

	return p.getRpcEvent(types.RPCPBehaviorResultEvent{
		Alarm:    event.Alarm,
		Entity:   event.Entity,
		PbhEvent: *pbhEvent,
		Error:    nil,
	})
}

type createPbehaviorMessageProcessor struct {
	FeaturePrintEventOnError bool
	DbClient                 mongo.DbClient
	PbhService               libpbehavior.Service
	EventManager             libpbehavior.EventManager
	AlarmAdapter             alarm.Adapter
	TimezoneConfigProvider   config.TimezoneConfigProvider
	Logger                   zerolog.Logger
}

func (p *createPbehaviorMessageProcessor) Process(
	ctx context.Context,
	alarm *types.Alarm,
	entity *types.Entity,
	params types.ActionPBehaviorParameters,
	_ []byte,
) (*types.Event, error) {
	pbehavior, err := p.createPbehavior(ctx, params, entity)
	if err != nil {
		return nil, err
	}

	err = p.recomputePbehavior(ctx, pbehavior.ID)
	if err != nil {
		return nil, err
	}

	resolveResult, err := p.getResolveResult(ctx, entity)
	if err != nil {
		return nil, err
	}

	if alarm == nil {
		alarms := make([]types.Alarm, 0)
		err := p.AlarmAdapter.GetOpenedAlarmsByIDs(ctx, []string{entity.ID}, &alarms)
		if err != nil {
			return nil, fmt.Errorf("failed to find alarm: %w", err)
		}

		if len(alarms) == 0 {
			return nil, nil
		}

		alarm = &alarms[0]
	}

	pbhEvent := p.EventManager.GetEvent(resolveResult, *alarm, time.Now())
	if pbhEvent.EventType == "" {
		return nil, nil
	}

	pbhEvent.Entity = entity

	return &pbhEvent, nil
}

func (p *createPbehaviorMessageProcessor) createPbehavior(
	ctx context.Context,
	params types.ActionPBehaviorParameters,
	entity *types.Entity,
) (*libpbehavior.PBehavior, error) {
	typeCollection := p.DbClient.Collection(mongo.PbehaviorTypeMongoCollection)
	res := typeCollection.FindOne(ctx, bson.M{"_id": params.Type})
	if err := res.Err(); err != nil {
		if err == mongodriver.ErrNoDocuments {
			return nil, fmt.Errorf("pbehavior type not exist: %q", params.Type)
		} else {
			return nil, fmt.Errorf("cannot get pbehavior type: %w", err)
		}
	}

	reasonCollection := p.DbClient.Collection(mongo.PbehaviorReasonMongoCollection)
	res = reasonCollection.FindOne(ctx, bson.M{"_id": params.Reason})
	if err := res.Err(); err != nil {
		if err == mongodriver.ErrNoDocuments {
			return nil, fmt.Errorf("pbehavior reason not exist: %q", params.Reason)
		} else {
			return nil, fmt.Errorf("cannot get pbehavior reason: %w", err)
		}
	}

	now := types.NewCpsTime()
	var start, stop types.CpsTime
	if params.Tstart != nil && params.Tstop != nil {
		start = types.NewCpsTime(*params.Tstart)
		stop = types.NewCpsTime(*params.Tstop)
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
		Author:  params.UserID, // since author now contains username, we should use user_id in author
		Enabled: true,
		Filter:  fmt.Sprintf(`{"_id": "%s"}`, entity.ID),
		Name:    params.Name,
		Reason:  params.Reason,
		RRule:   params.RRule,
		Start:   &start,
		Stop:    &stop,
		Type:    params.Type,
		Created: now,
		Updated: now,
	}

	collection := p.DbClient.Collection(mongo.PbehaviorMongoCollection)
	_, err := collection.InsertOne(ctx, pbehavior)
	if err != nil {
		return nil, fmt.Errorf("create new pbehavior failed: %w", err)
	}

	p.Logger.Info().Str("pbehavior", pbehavior.ID).Msg("create pbehavior")
	return &pbehavior, nil
}

func (p *createPbehaviorMessageProcessor) recomputePbehavior(ctx context.Context, pbehaviorID string) error {
	err := p.PbhService.RecomputeByID(ctx, pbehaviorID)

	if err != nil {
		return fmt.Errorf("pbehavior recompute failed: %w", err)
	}

	p.Logger.Debug().Str("pbehavior", pbehaviorID).Msg("pbehavior recomputed")

	return nil
}

func (p *createPbehaviorMessageProcessor) getResolveResult(ctx context.Context, entity *types.Entity) (libpbehavior.ResolveResult, error) {
	location := p.TimezoneConfigProvider.Get().Location
	now := time.Now().In(location)
	resolveResult, err := p.PbhService.Resolve(ctx, entity.ID, now)
	if err != nil {
		return libpbehavior.ResolveResult{}, fmt.Errorf("resolve an entity failed: %w", err)
	}

	return resolveResult, nil
}

func (p *rpcServerMessageProcessor) getErrRpcEvent(err error) []byte {
	msg, _ := p.getRpcEvent(types.RPCPBehaviorResultEvent{Error: &types.RPCError{Error: err}})
	return msg
}

func (p *rpcServerMessageProcessor) getRpcEvent(event types.RPCPBehaviorResultEvent) ([]byte, error) {
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
